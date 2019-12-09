#!/usr/local/bin/dumb-init /bin/bash
set -eo pipefail
[[ ${DEBUG:-} == true ]] && set -x
base=$PWD
. "$base/ops-manager-cloudfoundry/ci/tasks/helpers/cf-helper.sh"
. "$base/ops-manager-cloudfoundry/ci/tasks/helpers/tmp-helper.sh"
make_pcf_metadata
instance_name="test-mongodb-service"
app_name="app-ruby-sample"

cf_login
delete_service_app_if_exists $instance_name $app_name
create_service $instance_name
cf push $app_name -p $base/ops-manager-cloudfoundry/src/smoke-tests/assets/cf-mongo-example-app
cf bind-service $app_name $instance_name --binding-name mongodb-test-binding
cf restage $app_name
check_app_started $app_name
cf logout
