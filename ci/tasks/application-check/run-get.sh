#!/usr/local/bin/dumb-init /bin/bash
set -eo pipefail
[ 'true' = "${DEBUG:-}" ] && set -x
base=$PWD

host=$(echo $(cf apps | grep app-ruby-sample | awk '{print $6}'))
end_point="http://${host}/service/mongo/test3"
result=$(echo $(curl -X GET -H "Content-Type: application/json" ${end_point}))
if [ "${result}" = '{"data":"sometest130"}' ]; then
        echo "Application is working"
        curl -X DELETE -H "Content-Type: application/json" ${end_point}
    else
        echo "GET ${end-point} = ${result}"
        echo "FAILED. Application doesn't work"
        exit 1
fi