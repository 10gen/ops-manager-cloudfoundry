#!/usr/local/bin/dumb-init /bin/bash
set -euo pipefail
[ 'true' = "${DEBUG:-}" ] && set -x

cf create-service mongodb-odb replica_set_small mongodb-service-instance -c '{ "backup_enabled" : "true"}'
cf push app-mongodb-on-demand
cf bind-service app-mongodb-on-demand mongodb-service-instance --binding-name mongodb-service
cf restage app-mongodb-on-demand