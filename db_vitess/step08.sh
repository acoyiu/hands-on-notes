#!/bin/bash

source ./env.sh

web_port=15001
grpc_port=15991
mysql_server_port=15306
mysql_server_socket_path="/tmp/mysql.sock"

# Start vtgate
vtgate \
  $TOPOLOGY_FLAGS \
  --log_dir $VTDATAROOT/tmp \
  --log_queries_to_file $VTDATAROOT/tmp/vtgate_querylog.txt \
  --port $web_port \
  --grpc_port $grpc_port \
  --mysql_server_port $mysql_server_port \
  --mysql_server_socket_path $mysql_server_socket_path \
  --cell $cell \
  --cells_to_watch $cell \
  --tablet_types_to_wait PRIMARY,REPLICA \
  --service_map 'grpc-vtgateservice' \
  --pid_file $VTDATAROOT/tmp/vtgate.pid \
  --mysql_auth_server_impl none &

# Block waiting for vtgate to be listening
# Not the same as healthy

echo "Waiting for vtgate to be up..."
while true; do
  curl -I "http://$hostname:$web_port/debug/status" >/dev/null 2>&1 && break
  sleep 0.5
done;
echo "vtgate is up!"

echo "Access vtgate at http://$hostname:$web_port/debug/status"

# disown -a
