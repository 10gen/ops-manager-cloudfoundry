om_api_call() {
  method=$1
  path=$2
  data=$3
  res=$(curl -sS -H "Content-Type: application/json" -X $method --fail --digest  -u "$MONGO_OM_USER:$MONGO_OM_API_KEY" ${MONGO_OM_URL}/api/public/v1.0/groups/$MONGO_OM_GROUP_ID/$path -d "$data")
  echo "$res"
}

wait_for_service() {
  elapsed=0
  until [ $elapsed -ge 900 ]
  do
    mongo \
      $1 admin --eval 'quit(db.runCommand({ping: 1}).ok ? 0 : 1)' \
      --quiet &> /dev/null && break
    elapsed=$[$elapsed+5]
    sleep 5
  done

  if [ "$elapsed" -ge "900" ]; then
     echo "ERROR:  Cannot connect to MongoDB. Exiting..."
     exit 1
  fi
}
