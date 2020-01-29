#!/usr/local/bin/dumb-init /bin/bash
set -euo pipefail
[[ ${DEBUG:-} == true ]] && set -x

base=$PWD
PCF_URL="$PCF_URL"
PCF_USERNAME="$PCF_USERNAME"
PCF_PASSWORD="$PCF_PASSWORD"
. "$base/ops-manager-cloudfoundry/ci/tasks/helpers/tmp-helper.sh"

if [ -z "${VERSION:-}" ]; then
	echo "missing version number"
	exit 1
fi

ls "$base"
TILE_FILE=$(
	cd tileold
	ls *.pivotal
)
if [ -z "${TILE_FILE}" ]; then
	echo "No files matching tileold/*.pivotal"
	ls -lR artifacts
	exit 1
fi

STEMCELL_FILE=$(
	cd stemcell
	ls *bosh-stemcell-*.tgz
)
if [ -z "${STEMCELL_FILE}" ]; then
	echo "No files matching stemcell/*.tgz"
	ls -lR stemcell
	exit 1
fi

PRODUCT="$(yq r $base/ops-manager-cloudfoundry/tile/tile.yml name)"
echo "Product: " $PRODUCT

om="om -t $PCF_URL -u $PCF_USERNAME -p $PCF_PASSWORD -k"
cf_version=$(${om} available-products -f json | jq -j '.[] | select(.name == "cf").version')
echo ${om} $TILE_FILE
${om} upload-product --product "tileold/$TILE_FILE"
${om} upload-stemcell --stemcell "stemcell/$STEMCELL_FILE"
${om} available-products
${om} stage-product --product-name "$PRODUCT" --product-version "$VERSION"
${om} stage-product --product-name cf --product-version "$cf_version"

# if ${VERSION} = '1.0.5'; then
# 	${om} delete-product --product-name "$PRODUCT"
# 	${om} unstage-product --product-name "$PRODUCT"
# fi

# echo ${om} configure-product --product-name "$PRODUCT" --config "$base/ops-manager-cloudfoundry/ci/tasks/deploy-tile/config"
# ${om} configure-product  --config "$base/ops-manager-cloudfoundry/ci/tasks/deploy-tile/config"
cd ops-manager-cloudfoundry
# cat > "$base/ops-manager-cloudfoundry/ci/tasks/deploy-tile-old/vars.yml" << EOF

# OM_API_KEY: "$OM_API_KEY"
# OM_API_USER: "$OM_API_USER"

# EOF
# ${om} configure-product --product-name "$PRODUCT" --product-network "$network_config" --product-properties "$properties_config"
config_path=$base/ops-manager-cloudfoundry/ci/tasks/deploy-tile/config.pie
make_env_config $config_path

# replace "mongodb_broker" with "broker" for tile versions before 1.2
if [ "$VERSION" = "$(echo -e "$VERSION\n1.2" | sort -V | head -n1)" ]; then
	config_tmp=$(mktemp)
	yq r -j $config_path | jq '."resource-config".broker = ."resource-config".mongodb_broker | del(."resource-config".mongodb_broker)' >$config_tmp
	mv $config_tmp $config_path
fi

export OM_API_USER=$(yq r $config_path product-properties[.properties.username].value)
export OM_API_KEY=$(yq r $config_path product-properties[.properties.api_key].value.secret)
#${om} configure-product  --config "$config_path" -l "$base/ops-manager-cloudfoundry/ci/tasks/deploy-tile-old/vars.yml"
#TODO check if OM_API prefix real works as it says in docs
${om} configure-product --config "$config_path" --vars-env OM_API

STAGED=$(${om} curl --path /api/v0/staged/products)
# get GUIDs of cf and $PRODUCT
RESULT=$(echo "$STAGED" | jq --arg product_name "$PRODUCT" '.[] | select(.type == $product_name or .type == "cf") | .guid')
# merge GUIDs
RESULT=$(echo "$RESULT" | jq --slurp)
DATA=$(echo '{"deploy_products": []}' | jq ".deploy_products += $RESULT")

${om} curl --path /api/v0/installations --request POST --data "$DATA"
${om} apply-changes --reattach -n cf,"$PRODUCT"
${om} delete-unused-products
