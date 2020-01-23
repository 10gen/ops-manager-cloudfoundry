#!/usr/local/bin/dumb-init /bin/bash
set -euo pipefail
[[ ${DEBUG:-} == true ]] && set -x

base=$PWD

VERSION=$(cat "$base"/version/number)

if [ -z "${VERSION:-}" ]; then
  echo "missing version number"
  exit 1
fi
rm -r -f "$base/ops-manager-cloudfoundry/.dev_builds"
rm -r -f "$base/ops-manager-cloudfoundry/dev_releases"
rm -r -f "$base/ops-manager-cloudfoundry/tile/releases/*"
rm -r -f "$base/ops-manager-cloudfoundry/tile/resources/mongodb-*"
rm -r -f "$base/artefacts/mongodb-on-demand-${VERSION}.pivotal"
mkdir -p "$base"/ops-manager-cloudfoundry/src/mongodb

cp "$base"/on-demand-service-broker-release/on-demand-service-broker-*.tgz "$base"/ops-manager-cloudfoundry/tile/resources
cp "$base"/syslog-migration-release/syslog-migration-*.tgz "$base"/ops-manager-cloudfoundry/tile/resources
cp "$base"/pcf-mongodb-helpers/pcf-mongodb-helpers-*.tgz "$base"/ops-manager-cloudfoundry/tile/resources
cp "$base"/bpm-release/bpm-release-*.tgz "$base"/ops-manager-cloudfoundry/tile/resources
cp "$base"/mongodb/mongodb-linux-x86_64-ubuntu1604-*.tgz "$base"/ops-manager-cloudfoundry/src/mongodb
ls "$base"/ops-manager-cloudfoundry/tile/resources
(
  cd ops-manager-cloudfoundry
  cat >config/private.yml <<EOF
---
blobstore:
  options:
    access_key_id: "$AWS_KEY"
    secret_access_key: "$AWS_SECRET_KEY"
EOF
  rm -r -f dev_releases
  rm -r -f tile/product/*
  rm -r -f tile/resources/mongodb-*

  tarball_path="$base/ops-manager-cloudfoundry/tile/resources/mongodb-${VERSION}.tgz"
  mkdir -p "$(dirname "$tarball_path")"
  bosh -n create-release --sha2 --tarball="$tarball_path" --version="${VERSION}" --force
)

(
  cd ops-manager-cloudfoundry/tile

  yq r -j tile.yml >step1.json
  jq --arg v "${VERSION}" '(.. | objects | select(has("releases")).releases[] | select(.name == "mongodb").version) = $v' step1.json >step2.json
  jq --arg p "$(ls resources/mongodb-${VERSION}.tgz)" '(.packages[] | select(.name == "mongodb").path) = $p' step2.json >step3.json

  mv step3.json tile.yml
  tile build "${VERSION}"
)

mkdir -p "$base"/artifacts
cp "$base"/ops-manager-cloudfoundry/tile/product/mongodb-on-demand-*.pivotal "$base"/artifacts/
