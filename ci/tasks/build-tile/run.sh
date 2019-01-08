#!/usr/local/bin/dumb-init /bin/bash
set -euo pipefail
[ 'true' = "${DEBUG:-}" ] && set -x

base=$PWD

VERSION=$(cat "$base"/ops-manager-cloudfoundry/version/number)
if [ -z "${VERSION:-}" ]; then
  echo "missing version number"
  exit 1
fi

(
cd ops-manager-cloudfoundry

rm -r -f dev_releases
rm -r -f tile/product/*
rm -r -f tile/resources/mongodb-*

tarball_path="$base/ops-manager-cloudfoundry/tile/resources/mongodb-${VERSION}.tgz"
mkdir -p "$(dirname "$tarball_path")"
bosh -n create-release --sha2 --tarball="$tarball_path" --version="${VERSION}" --force
)

(
cd ops-manager-cloudfoundry/tile

yq w -i tile.yml packages.[4].path "$(ls resources/mongodb-*.tgz)"
yq w -i tile.yml packages.[4].jobs[0].properties.service_deployment.releases[0].version "${VERSION}"
yq w -i tile.yml runtime_configs[0].runtime_config.releases[0].version "${VERSION}"
tile build "${VERSION}"
)
mkdir -p "$base"/artifacts
cp "$base"/ops-manager-cloudfoundry/tile/product/mongodb-on-demand-*.pivotal "$base"/artifacts/