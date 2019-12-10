#!/usr/local/bin/dumb-init /bin/bash
set -eo pipefail
[[ ${DEBUG:-} = true ]] && set -x
base=$PWD
app_name="app-ruby-sample"
. "$base/ops-manager-cloudfoundry/ci/tasks/helpers/cf-helper.sh"
. "$base/ops-manager-cloudfoundry/ci/tasks/helpers/tmp-helper.sh"
make_pcf_metadata

cf_login
check_app_started $instance_name $app_name
host=$(echo $(cf apps | grep $app_name | awk '{print $6}'))
url="http://${host}/service/mongo/test3"
result=$(echo $(curl -X GET ${url}))
if [ "${result}" = '{"data":"sometest130"}' ]; then
    echo "Application is working"
    echo "Cleaning data.."
    curl -X DELETE ${url}
else
    echo "GET ${url} finished with result: ${result}"
    echo "FAILED. Application doesn't work"
    exit 1
fi
cf logout
