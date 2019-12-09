#!/usr/local/bin/dumb-init /bin/bash
set -eo pipefail
[[ ${DEBUG:-} = true ]] && set -x

cf_login
host=$(echo $(cf apps | grep app-ruby-sample | awk '{print $6}'))
url="http://${host}/service/mongo/test3"
echo "send PUT to ${url} with {data:sometest130}"
status=$(curl -X PUT -d '{"data":"sometest130"}' ${url})
if [[ $status != "success" ]]; then
    echo "Error: check application"
    exit 1
fi
cf logout
