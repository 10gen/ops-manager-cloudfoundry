#!/usr/local/bin/dumb-init /bin/bash
# shellcheck shell=bash

set -euo pipefail
[[ ${DEBUG:-} == true ]] && set -x

base=$PWD

# shellcheck source=ci/tasks/helpers/deploy.sh
source "$base/ops-manager-cloudfoundry/ci/tasks/helpers/deploy.sh"

VERSION=$(cat version/number)
if [ -z "${VERSION:-}" ]; then
	echo "missing version number"
	exit 1
fi

TILE_FILE=$(
	cd artifacts
	ls -- *-"${VERSION}".pivotal 2>/dev/null || true
)
if [ -z "${TILE_FILE}" ]; then
	echo "No files matching artifacts/*-${VERSION}.pivotal"
	ls -lR artifacts
	exit 1
fi

install_product "$VERSION" "$TILE_FILE"
