#!/usr/local/bin/dumb-init /bin/bash
# shellcheck shell=bash

set -euo pipefail
[[ ${DEBUG:-} == true ]] && set -x

# shellcheck source=ci/tasks/helpers/deploy.sh
source "ops-manager-cloudfoundry/ci/tasks/helpers/deploy.sh"

if [ -z "${VERSION:-}" ]; then
	echo "missing version number"
	exit 1
fi

TILE_FILE=$(
	cd tileold
	ls -- *-"${VERSION}".pivotal
)
if [ -z "${TILE_FILE}" ]; then
	echo "No files matching tileold/*-${VERSION}.pivotal"
	ls -lR tileold
	exit 1
fi

install_product "$VERSION" "$TILE_FILE"
