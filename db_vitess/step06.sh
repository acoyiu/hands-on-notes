#!/bin/bash

source ./env.sh

# set log directory and port
log_dir="${VTDATAROOT}/tmp"
port=16000

vtorc \
  $TOPOLOGY_FLAGS \
  --logtostderr \
  --alsologtostderr \
  --config="./_vtorc-config.json" \
  --port $port \
  > "${log_dir}/vtorc.out" 2>&1 &

# get the last process ID and save
vtorc_pid=$!
echo ${vtorc_pid} > "${log_dir}/vtorc.pid"

echo "\
vtorc is running!
  - UI: http://localhost:${port}
  - Logs: ${log_dir}/vtorc.out
  - PID: ${vtorc_pid}
" > ~/temp/out.txt
