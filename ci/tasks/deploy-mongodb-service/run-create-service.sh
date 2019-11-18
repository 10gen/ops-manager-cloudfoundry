#!/usr/local/bin/dumb-init /bin/bash
set -eo pipefail
[[ ${DEBUG:-} = true ]] && set -x
base=$PWD
. "$base/ops-manager-cloudfoundry/ci/tasks/helpers/cf-helper.sh"
instance_name="test-mongodb-service"

cf login -a $CF_APP_URL -u $CF_APP_USER -p $CF_APP_PASSWORD --skip-ssl-validation -o system -s system
delete_service_if_exists $instance_name
cf create-service mongodb-odb "$SET_PLAN" $instance_name -c "{\"enable_backup\":\"$BACKUP_ENABLED\"}"
wait_service_status_change $instance_name "create in progress"
service_status=$(cf services | awk  '/$instance_name.*succeeded/{print "succeeded"}')
if [[ $service_status == "succeeded" ]]; then
  cf push app-ruby-sample -p $base/ops-manager-cloudfoundry/src/smoke-tests/assets/cf-mongo-example-app
  cf bind-service app-ruby-sample $instance_name --binding-name mongodb-service
  cf restage app-ruby-sample
  cf logout
else
  echo "FAILED! wrong status: ${cf service test-mongodb-service}"
  cf logout
  exit 1
fi
