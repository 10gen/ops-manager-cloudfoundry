#!/usr/local/bin/dumb-init /bin/bash
set -eo pipefail
[[ ${DEBUG:-} = true ]] && set -x

cf login -a $CF_APP_URL -u $CF_APP_USER -p $CF_APP_PASSWORD --skip-ssl-validation -o system -s system
delete_service_app_if_exists "test-mongodb-service" "app-ruby-sample"
cf logout
