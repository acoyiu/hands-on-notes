source ./env.sh

echo "Starting vtctld..."

vtctld \
  $TOPOLOGY_FLAGS \
  --cell $cell \
  --service_map 'grpc-vtctl,grpc-vtctld' \
  --backup_storage_implementation file \
  --file_backup_storage_root $VTDATAROOT/backups \
  --log_dir $VTDATAROOT/tmp \
  --port $vtctld_web_port \
  --grpc_port $grpc_port \
  --pid_file $VTDATAROOT/tmp/vtctld.pid &

echo "vtctld ran"