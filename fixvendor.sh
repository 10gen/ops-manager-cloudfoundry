#!/bin/bash

set -e

for dir in src/{smoke-tests,mongodb-{config-agent,service-adapter}}; do
    pushd $dir
    echo "Fixing vendor in $dir..."
    go mod tidy
    go mod vendor
    popd
done
