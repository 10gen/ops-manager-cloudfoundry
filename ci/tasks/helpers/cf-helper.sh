#!/usr/local/bin/dumb-init /bin/bash

#cf.helper. wait for particular service status 
wait_service_status_change() {
  local instance_name=$1
  local status=$2
  local time=0
  echo "checking " $instance_name $status
  local verify_status=$(cf services | awk  '/'"$instance_name"'[ ].*'"$status"'/{print "'"$status"'"}')
  while [[ $verify_status == $status ]] && [[ $time -lt $INSTALL_TIMEOUT ]]; do
    echo "...${verify_status}"
    sleep 3m
    let "time=$time+3"
    verify_status=$(cf services | awk  '/'"$instance_name"'[ ].*'"$status"'/{print "'"$status"'"}')
  done
}

delete_service_app_if_exists() {
  local instance_name=$1
  local app_name=$2
  local service=$(cf services | awk '/'"$instance_name"'[ $]/{print "exist"}')
  if [[ $service == "exist" ]]; then
    echo "check if $app_name is exist"
    local app=$(cf apps | awk '/'"$app_name"'[ $]/{print "exist"}')
    if [[ $app == "exist" ]]; then
      cf unbind-service $app_name $instance_name
      check_app_unbinding $app_name
      cf delete $app_name -f
    fi
    cf delete-service $instance_name -f
    wait_service_status_change $instance_name "delete in progress"
    service_status=$(cf services | awk  '/'"$instance_name"'[ $].*failed/{print "failed"}')
    if [[ $service_status == "failed" ]]; then
      cf purge-service-instance $instance_name -f
    fi
  fi
}

create_service() {
  local instance_name=$1
  cf create-service mongodb-odb "$SET_PLAN" $instance_name -c "{\"enable_backup\":\"$BACKUP_ENABLED\"}"
  wait_service_status_change $instance_name "create in progress"
  service_status=$(cf services | awk  '/'"$instance_name"'[ ].*succeeded/{print "succeeded"}')
  if [[ $service_status != "succeeded" ]]; then
    echo "FAILED! wrong status: $(cf service $instance_name)"
    cf logout
    exit 1
  fi
}

check_app_unbinding() {
  local instance_name=$1
  local app_name=$2
  local app_binding=$(cf services | grep "$instance_name " | awk '!/'"$app_name"'/{print "not bounded"}')
  local try=10
  until [[ $app_binding == "not bounded" ]]; do
    app_binding=$(cf services | grep "$instance_name " | awk '!/'"$app_name"'/{print "not bounded"}')
    if [[ $try -lt 0 ]]; then
      echo "ERROR: unbinding is getting too long"
      exit 1
    fi
    let "try--"
    echo "checking unbinding ($try)"
  done
}

check_app_started() {
  local app_name=$1
  local app=$(cf apps | grep "$app_name " | awk  '/started/{print "started"}')
  local try=10
  until [[ $app == "started" ]]; do
    app=$(cf apps | grep "$app_name " | awk  '/started/{print "started"}')
    if [[ $try -lt 0 ]]; then
      echo "ERROR: unbinding is getting too long"
      exit 1
    fi
    let "try--"
    echo "checking application status ($try)"
  done
  echo "Application started"
}