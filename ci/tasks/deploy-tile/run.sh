#!/usr/local/bin/dumb-init /bin/bash
set -euo pipefail
[[ ${DEBUG:-} == true ]] && set -x

base=$PWD
. "$base/ops-manager-cloudfoundry/ci/tasks/helpers/tmp-helper.sh"
. "$base/ops-manager-cloudfoundry/ci/tasks/helpers/deploy.sh"

VERSION=$(cat "$base"/version/number)
if [ -z "${VERSION:-}" ]; then
	echo "missing version number"
	exit 1
fi

TILE_FILE=$(
	cd artifacts
	ls *-${VERSION}.pivotal
)
if [ -z "${TILE_FILE}" ]; then
	echo "No files matching artifacts/*.pivotal"
	ls -lR artifacts
	exit 1
fi

install_product $VERSION $TILE_FILE
