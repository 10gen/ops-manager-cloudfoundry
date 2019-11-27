#!/usr/local/bin/dumb-init /bin/bash
set -euo pipefail
[[ ${DEBUG:-} == true ]] && set -x

base=$PWD
. "$base/ops-manager-cloudfoundry/ci/tasks/helpers/tmp-helper.sh"
config_path="$base/ops-manager-cloudfoundry/ci/tasks/deploy-tile/config.pie"
make_env_config $config_path

cat >metadata <<EOF
---
opsmgr:
  url: "$PCF_URL"
  username: "$PCF_USERNAME"
  password: "$PCF_PASSWORD"
EOF

cat >config.json <<EOF
{
  "service_name": "mongodb-odb",
  "mongodb_version": $MONGO_VERSION,
  "plan_names": $PLAN_NAMES,
  "backup_enabled": $BACKUP_ENABLED,
  "ssl_enabled": $SSL_ENABLED,
  "retry": {
    "max_attempts": 10,
    "backoff": "linear",
    "baseline_interval_milliseconds": 1000
  },
  "apps_domain": "$(pcf cf-info | grep apps_domain | cut -d' ' -f 3)",
  "system_domain": "$(pcf cf-info | grep system_domain | cut -d' ' -f 3)",
  "api": "api.$(pcf cf-info | grep system_domain | cut -d' ' -f 3)",
  "admin_user": "$(pcf cf-info | grep admin_username | cut -d' ' -f 3)",
  "admin_password": "$(pcf cf-info | grep admin_password | cut -d' ' -f 3)",
  "skip_ssl_validation": true,
  "create_permissive_security_group": false
}
EOF

cat >mongo-ops.json <<EOF
{
    "url": "$(yq r $config_path product-properties[.properties.url].value)",
    "username": "$(yq r $config_path product-properties[.properties.username].value)",
    "api_key": "$(yq r $config_path product-properties[.properties.api_key].value.secret)",
    "group": "5d68ef6a6a94055d183032e5",
    "auth_key": "5d68fcd06a94055d18308834f769efdc2edb26f530db6269411aceb8",
    "nodes": "localhost",
    "org": "5d68ef6a6a94055d183032e5"
}
EOF

export CONFIG_PATH="$base"/config.json

PACKAGE_NAME=github.com/10gen/ops-manager-cloudfoundry
PACKAGE_DIR=$GOPATH/src/$PACKAGE_NAME
mkdir -p $PACKAGE_DIR
cp -a ops-manager-cloudfoundry/* $PACKAGE_DIR

(
  cd $PACKAGE_DIR/src/smoke-tests

  ./bin/test
)
