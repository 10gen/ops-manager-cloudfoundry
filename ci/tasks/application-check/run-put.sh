#!/usr/local/bin/dumb-init /bin/bash
set -eo pipefail
[ 'true' = "${DEBUG:-}" ] && set -x

cf login -a $CF_APP_URL -u $CF_APP_USER -p $CF_APP_PASSWORD --skip-ssl-validation -o system -s system
host=$(echo $(cf apps | grep app-ruby-sample | awk '{print $6}'))
url="http://${host}/service/mongo/test3"
echo "send PUT to ${url} with {data:sometest130}"
curl -X PUT -d '{"data":"sometest130"}' ${url}
