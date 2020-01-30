#!/usr/local/bin/dumb-init /bin/bash
# shellcheck shell=bash

set -euo pipefail
[[ ${DEBUG:-} == true ]] && set -x

VERSION=$(cat version/number)
if [ -z "${VERSION:-}" ]; then
  echo "missing version number"
  exit 1
fi

TILE_FILE=$(
  cd artifacts
  ls -- *-"$VERSION".pivotal
)
if [ -z "$TILE_FILE" ]; then
  echo "No files matching artifacts/$TILE_FILE.pivotal"
  ls -lR artifacts
  exit 1
fi

mkdir -p release
cp artifacts/"$TILE_FILE" release/

SHA256=$(sha256sum "artifacts/$TILE_FILE" | cut -d ' ' -f 1)
cat >release/notification <<EOF
<!here> New build v${VERSION} released!
You can download build at https://$RELEASE_BUCKET_NAME.s3.amazonaws.com/$TILE_FILE
Build SHA256: ${SHA256}
EOF
echo "Build URL: https://$RELEASE_BUCKET_NAME.s3.amazonaws.com/$TILE_FILE"
echo "Build SHA256: $SHA256"
