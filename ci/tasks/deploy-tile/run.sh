#!/usr/local/bin/dumb-init /bin/bash
set -euo pipefail
[ 'true' = "${DEBUG:-}" ] && set -x

base=$PWD
PCF_URL="$PCF_URL"
PCF_USERNAME="$PCF_USERNAME"
PCF_PASSWORD="$PCF_PASSWORD"


VERSION=$(cat "$base"/version/number)
if [ -z "${VERSION:-}" ]; then
  echo "missing version number"
  exit 1
fi

TILE_FILE=`cd artifacts; ls *-${VERSION}.pivotal`
if [ -z "${TILE_FILE}" ]; then
	echo "No files matching artifacts/*.pivotal"
	ls -lR artifacts
	exit 1
fi

STEMCELL_FILE=`cd stemcell; ls *bosh-stemcell-*.tgz`
if [ -z "${STEMCELL_FILE}" ]; then
	echo "No files matching stemcell/*.tgz"
	ls -lR stemcell
	exit 1
fi

PRODUCT="$(cat $base/ops-manager-cloudfoundry/tile/tile.yml | grep '^name' | cut -d' ' -f 2)"
echo "Product " $PRODUCT
# if [ -z "${PRODUCT}" ]; then
# 	PRODUCT=mongodb-on-demand
# fi
om="om -t $PCF_URL -u $PCF_USERNAME -p $PCF_PASSWORD -k"

${om} upload-product --product "artifacts/$TILE_FILE"
${om} upload-stemcell --stemcell "stemcell/$STEMCELL_FILE"
${om} available-products
${om} stage-product --product-name "$PRODUCT" --product-version "$VERSION"


# if ${VERSION} = '1.0.5'; then
# 	${om} delete-product --product-name "$PRODUCT"
# 	${om} unstage-product --product-name "$PRODUCT"
# fi


	# echo ${om} configure-product --product-name "$PRODUCT" --config "$base/ops-manager-cloudfoundry/ci/tasks/deploy-tile/config"
	# ${om} configure-product  --config "$base/ops-manager-cloudfoundry/ci/tasks/deploy-tile/config"
cd ops-manager-cloudfoundry
cat > "$base/ops-manager-cloudfoundry/ci/tasks/deploy-tile/vars.yml" << EOF

OM_API_KEY: "$OM_API_KEY"
OM_API_USER: "$OM_API_USER"

EOF
# ${om} configure-product --product-name "$PRODUCT" --product-network "$network_config" --product-properties "$properties_config"
${om} configure-product  --config "$base/ops-manager-cloudfoundry/ci/tasks/deploy-tile/$CONFIG" -l "$base/ops-manager-cloudfoundry/ci/tasks/deploy-tile/vars.yml"

STAGED=$(${om} curl --path /api/v0/staged/products)
RESULT=$(echo "$STAGED" | jq --arg product_name "$PRODUCT" 'map(select(.type == $product_name)) | .[].guid')
DATA=$(echo '{"deploy_products": []}' | jq ".deploy_products += [$RESULT]")

${om} curl --path /api/v0/installations --request POST --data "$DATA"
${om} apply-changes -n cf,mongodb-on-demand
${om} delete-unused-products