#!/usr/local/bin/dumb-init /bin/bash
set -euo pipefail
[ 'true' = "${DEBUG:-}" ] && set -x
base=$PWD

cf login -a $CF_APP_URL -u $CF_APP_USER -p $CF_APP_PASSWORD --skip-ssl-validation -o system -s system
cf create-service mongodb-odb $REPLICA_SET_PLAN mongodb-service-instance -c "{\"enable_backup\":\"$BACKUP_ENABLED\"}"
cf push app-ruby-sample -p $base/ops-manager-cloudfoundry/src/smoke-tests/assets/cf-mongo-example-app
service_status = "cf services | grep mongodb-service-instance | awk '{print $5" "$6" "$7}'"
until [ ${service_status} != "create in progress" ]; do
    echo "."
    sleep 3m
    service_status = "cf services | grep mongodb-service-instance | awk '{print $5" "$6" "$7}'"
done
if [ ${service_status} = "create succeeded" ]; then
    cf bind-service app-ruby-sample mongodb-service-instance --binding-name mongodb-service
    cf restage app-ruby-sample
  else
    echo "FAILED! wrong status: ${service_status}"
    exit 1
fi
