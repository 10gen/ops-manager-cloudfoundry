#!/usr/local/bin/dumb-init /bin/bash
# shellcheck shell=bash
# shellcheck source=ci/tasks/helpers/tmp-helper.sh
source "ops-manager-cloudfoundry/ci/tasks/helpers/tmp-helper.sh"

install_product() {
	VERSION=$1
	TILE_FILE=$2
	PCF_URL="$PCF_URL"
	PCF_USERNAME="$PCF_USERNAME"
	PCF_PASSWORD="$PCF_PASSWORD"
	UPDATE_PAS="$UPDATE_PAS"

	STEMCELL_FILE=$(
		cd stemcell
		ls -- *bosh-stemcell-*.tgz
	)
	if [ -z "${STEMCELL_FILE}" ]; then
		echo "No files matching stemcell/*bosh-stemcell-*.tgz"
		ls -lR stemcell
		exit 1
	fi

	PRODUCT="$(yq r ops-manager-cloudfoundry/tile/tile.yml name)"
	echo "Product $PRODUCT"

	om="om -t $PCF_URL -u $PCF_USERNAME -p $PCF_PASSWORD -k"
	cf_version=$(${om} available-products -f json | jq -j '.[] | select(.name == "cf").version')

	${om} upload-product --product "$TILE_FILE"
	${om} upload-stemcell --stemcell "stemcell/$STEMCELL_FILE"
	${om} available-products
	${om} stage-product --product-name "$PRODUCT" --product-version "$VERSION"
	${om} stage-product --product-name cf --product-version "$cf_version"

	config_path=ops-manager-cloudfoundry/ci/tasks/deploy-tile/config.pie
	make_env_config "$config_path"
	# replace "mongodb_broker" with "broker" for tile versions before 1.2
	if [ "$VERSION" = "$(echo -e "$VERSION\n1.2" | sort -V | head -n1)" ]; then
		config_tmp=$(mktemp)
		yq r -j $config_path | jq '."resource-config".broker = ."resource-config".mongodb_broker | del(."resource-config".mongodb_broker)' >$config_tmp
		mv $config_tmp $config_path
	fi
	OM_API_USER=$(yq r "$config_path" 'product-properties[.properties.username].value')
	OM_API_KEY=$(yq r "$config_path" 'product-properties[.properties.api_key].value.secret')
	export OM_API_KEY OM_API_USER
	${om} configure-product --config "$config_path" --vars-env OM_API

	STAGED=$(${om} curl --path /api/v0/staged/products)
	# get GUIDs of cf and $PRODUCT
	if [ "${UPDATE_PAS}" == "true" ]; then
		RESULT=$(echo "$STAGED" | jq --arg product_name "$PRODUCT" '.[] | select(.type == $product_name or .type == "cf") | .guid')
	else
		RESULT=$(echo "$STAGED" | jq --arg product_name "$PRODUCT" '.[] | select(.type == $product_name) | .guid')
	fi

	# merge GUIDs
	RESULT=$(echo "$RESULT" | jq --slurp)
	DATA=$(echo '{"deploy_products": []}' | jq ".deploy_products += $RESULT")

	${om} curl --path /api/v0/installations --request POST --data "$DATA"
	${om} apply-changes --reattach -n cf,"$PRODUCT"
	${om} delete-unused-products
}
