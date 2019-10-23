#!/usr/local/bin/dumb-init /bin/bash
set -euo pipefail
[ 'true' = "${DEBUG:-}" ] && set -x

base=$PWD

cat > metadata << EOF
---
opsmgr:
  url: "$PCF_URL"
  username: "$PCF_USERNAME"
  password: "$PCF_PASSWORD"
EOF

cat > config.json << EOF
{
  "service_name": "mongodb-odb",
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

export CONFIG_PATH="$base"/config.json

PACKAGE_NAME=github.com/10gen/ops-manager-cloudfoundry
PACKAGE_DIR=$GOPATH/src/$PACKAGE_NAME
mkdir -p $PACKAGE_DIR
cp -a ops-manager-cloudfoundry/* $PACKAGE_DIR

(
cd $PACKAGE_DIR/src/smoke-tests

./bin/test
)
