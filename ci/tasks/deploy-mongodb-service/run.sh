#!/usr/local/bin/dumb-init /bin/bash
set -euo pipefail
[ 'true' = "${DEBUG:-}" ] && set -x
cf login -a $CF_APP_URL -u $CF_APP_USER -p $CF_APP_PASSWORD --skip-ssl-validation -o system -s system
cf create-service mongodb-odb $REPLICA_SET_PLAN mongodb-service-instance -c "{\"enable_backup\":\"$BACKUP_ENABLED\"}"
cf push app-ruby-sample -p $base/ops-manager-cloudfoundry/src/smoke-tests/assets/cf-mongo-example-app
cf bind-service app-ruby-sample mongodb-service-instance --binding-name mongodb-service
cf restage app-ruby-sample