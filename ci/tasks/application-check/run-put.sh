#!/usr/local/bin/dumb-init /bin/bash
set -eo pipefail
[ 'true' = "${DEBUG:-}" ] && set -x
base=$PWD

host=$(echo $(cf apps | grep app-ruby-sample | awk '{print $6}'))
end_point="http://${host}/service/mongo/test3"
curl -X PUT -H "Content-Type: application/json" -d '{"data":"sometest130"}' ${end_point}
