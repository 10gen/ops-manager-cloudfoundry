#!/usr/local/bin/dumb-init /bin/bash
set -eo pipefail
[[ ${DEBUG:-} = true ]] && set -x

cf login -a $CF_APP_URL -u $CF_APP_USER -p $CF_APP_PASSWORD --skip-ssl-validation -o system -s system
cf unbind-service app-ruby-sample test-mongodb-service
cf delete-service test-mongodb-service -f
cf delete app-ruby-sample -f
cf logout
