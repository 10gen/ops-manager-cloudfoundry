#!/usr/local/bin/dumb-init /bin/bash
set -euo pipefail
[[ ${DEBUG:-} == true ]] && set -x

base=$PWD
# PCF_URL="$PCF_URL"
# PCF_USERNAME="$PCF_USERNAME"
# PCF_PASSWORD="$PCF_PASSWORD"
. "$base/ops-manager-cloudfoundry/ci/tasks/helpers/tmp-helper.sh"
. "$base/ops-manager-cloudfoundry/ci/tasks/helpers/deploy.sh"

if [ -z "${VERSION:-}" ]; then
	echo "missing version number"
	exit 1
fi

if [ ! -z "ls tileold/*.pivotal" ]; then
	TILE_FILE=$(
		ls tileold/*.pivotal
	)
else
	echo "No files matching tileold/*.pivotal"
	ls -lR tileold
	exit 1
fi

install_product $VERSION $TILE_FILE
