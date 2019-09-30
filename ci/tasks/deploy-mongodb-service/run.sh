#!/usr/local/bin/dumb-init /bin/bash
set -euo pipefail
[ 'true' = "${DEBUG:-}" ] && set -x
cf login -a $CF_APP_URL -u $CF_APP_USER -p $CF_APP_PASSWORD --skip-ssl-validation -o system -s system
cf create-service mongodb-odb $REPLICA_SET_PLAN mongodb-service-instance -c "{\"backup_enabled\":\"$CONFIG\"}"
cf push app-mongodb-on-demand
cf bind-service app-mongodb-on-demand mongodb-service-instance --binding-name mongodb-service
cf restage app-mongodb-on-demand