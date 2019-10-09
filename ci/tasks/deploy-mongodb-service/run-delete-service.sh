#!/usr/local/bin/dumb-init /bin/bash
set -eo pipefail
[ 'true' = "${DEBUG:-}" ] && set -x

cf login -a $CF_APP_URL -u $CF_APP_USER -p $CF_APP_PASSWORD --skip-ssl-validation -o system -s system
cf unbind-service app-ruby-sample mongodb-service-instance
cf delete-service mongodb-service-instance -f
cf delete app-ruby-sample -f
