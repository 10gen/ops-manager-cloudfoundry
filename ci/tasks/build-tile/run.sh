#!/usr/local/bin/dumb-init /bin/bash
# shellcheck shell=bash

set -euo pipefail
[[ ${DEBUG:-} == true ]] && set -x

VERSION=$(cat version/number)
if [ -z "${VERSION:-}" ]; then
  echo "missing version number"
  exit 1
fi

rm -r -f "ops-manager-cloudfoundry/.dev_builds"
rm -r -f "ops-manager-cloudfoundry/dev_releases"
rm -r -f "ops-manager-cloudfoundry/tile/releases/*"
rm -r -f "ops-manager-cloudfoundry/tile/resources/mongodb-*"
rm -r -f "artifacts/mongodb-on-demand-${VERSION}.pivotal"
mkdir -p ops-manager-cloudfoundry/src/mongodb
mkdir -p ops-manager-cloudfoundry/src/mongodb_versions

cp on-demand-service-broker-release/on-demand-service-broker-*.tgz ops-manager-cloudfoundry/tile/resources
cp syslog-migration-release/syslog-migration-*.tgz ops-manager-cloudfoundry/tile/resources
cp pcf-mongodb-helpers/pcf-mongodb-helpers-*.tgz ops-manager-cloudfoundry/tile/resources
cp bpm-release/bpm-release-*.tgz ops-manager-cloudfoundry/tile/resources
cp mongodb/mongodb-linux-x86_64-ubuntu1604-*.tgz ops-manager-cloudfoundry/src/mongodb
ls ops-manager-cloudfoundry/tile/resources
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

  wget https://opsmanager.mongodb.com/static/version_manifest/4.4.json
  mv 4.4.json src/mongodb_versions/versions.json

  tarball_path="tile/resources/mongodb-${VERSION}.tgz"
  mkdir -p "$(dirname "$tarball_path")"
  bosh -n create-release --sha2 --tarball="$tarball_path" --version="${VERSION}" --force
)

(
  cd ops-manager-cloudfoundry/tile

  yq r -j tile.yml >step1.json
  jq --arg v "${VERSION}" '(.. | objects | select(has("releases")).releases[] | select(.name == "mongodb").version) = $v' step1.json >step2.json
  jq --arg p "$(ls resources/mongodb-"${VERSION}".tgz)" '(.packages[] | select(.name == "mongodb").path) = $p' step2.json >step3.json

  mv step3.json tile.yml
  tile build "${VERSION}"
)

mkdir -p artifacts
cp ops-manager-cloudfoundry/tile/product/mongodb-on-demand-*.pivotal artifacts/
