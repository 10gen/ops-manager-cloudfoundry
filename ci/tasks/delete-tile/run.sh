#!/usr/local/bin/dumb-init /bin/bash
set -euo pipefail
[[ ${DEBUG:-} == true ]] && set -x

base=$PWD

PRODUCT="$(yq r $base/ops-manager-cloudfoundry/tile/tile.yml name)"
om --version
om="om -t $PCF_URL -u $PCF_USERNAME -p $PCF_PASSWORD -k"

echo "Retrieving current staged version of ${PRODUCT}"

product_version=$(${om} staged-products -f json | jq -r --arg product_name $PRODUCT '.[] | select(.name == $product_name) | .version')

if [ ! -z "${product_version}" ]; then
    echo "Deleting product [${PRODUCT}], version [${product_version}] , from ${PCF_URL}"
    ${om} unstage-product --product-name "$PRODUCT"
    #${om} apply-changes --product-name "$PRODUCT" --ignore-warnings true --reattach
    ${om} apply-changes --reattach --ignore-warnings true --skip-unchanged-products true
else
    echo "Check product [${PRODUCT}] - probably already unstaged"
fi

echo "Finish deleting ${PRODUCT}"
