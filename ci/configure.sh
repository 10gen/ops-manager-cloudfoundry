#!/usr/bin/env bash

set -eu

fly -t altoros set-pipeline -n \
 -p ops-manager-cloudfoundry \
 -c ./pipeline.yml \
 -l <(lpass show --note "pcf:ops-manager-cloudfoundry")
