#!/usr/bin/env bash
set -euo pipefail
[ 'true' = "${DEBUG:-}" ] && set -x

base=$PWD

PRODUCT="$(yq r $base/mongodb-on-demand-release/tile/tile.yml name)"

om="om -t $PCF_URL -u $PCF_USERNAME -p $PCF_PASSWORD -k"

echo "Retrieving current staged version of ${PRODUCT}"

product_version=$(${om} -f json deployed-products | jq -r --arg product_name $PRODUCT '.[] | select(.name == $product_name) | .version')

echo "Deleting product [${PRODUCT}], version [${product_version}] , from ${PCF_URL}"

${om} unstage-product --product-name "$PRODUCT"

${om} apply-changes --ignore-warnings true
