#!/usr/local/bin/dumb-init /bin/bash
set -eo pipefail
[[ ${DEBUG:-} = true ]] && set -x
base=$PWD

cf login -a $CF_APP_URL -u $CF_APP_USER -p $CF_APP_PASSWORD --skip-ssl-validation -o system -s system
service=$(cf services | awk '$1 ~ /test-mongodb-service/{print "exist"}')
if [[ $service == "exist" ]]; then
  cf purge-service-instance test-mongodb-service
fi
cf create-service mongodb-odb "$SET_PLAN" test-mongodb-service -c "{\"enable_backup\":\"$BACKUP_ENABLED\"}"
service_status=$(cf services | awk '$1 ~ /test-mongodb-service/{print $4" "$5" "$6}')
time=0
until [[ $service_status != "create in progress" ]] || [[ $time -gt $INSTALL_TIMEOUT ]]; do
  echo "...${service_status}"
  sleep 3m
  let "time=$time+3"
  service_status=$(cf services | awk '$1 ~ /test-mongodb-service/{print $4" "$5" "$6}') 
done
if [[ $service_status == "create succeeded mongodb-odb" ]]; then
  cf push app-ruby-sample -p $base/ops-manager-cloudfoundry/src/smoke-tests/assets/cf-mongo-example-app
  cf bind-service app-ruby-sample test-mongodb-service --binding-name mongodb-service
  cf restage app-ruby-sample
else
  echo "FAILED! wrong status: ${service_status}"
  exit 1
fi
cf logout
