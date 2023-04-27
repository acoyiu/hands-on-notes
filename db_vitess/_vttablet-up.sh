#!/bin/bash

source ./env.sh

# process config
uid=$TABLET_UID
mysql_port=$[17000 + $uid]
port=$[15000 + $uid]
grpc_port=$[16000 + $uid]
echo "================"
echo $uid
echo $mysql_port
echo $port
echo $grpc_port

# vitess config
keyspace=${KEYSPACE}
printf -v instance_name_alias '%s-%010d' $cell $uid
printf -v tablet_dir 'vt_%010d' $uid
printf -v tablet_logfile 'vttablet_%010d_querylog.txt' $uid
echo "================"
echo $keyspace
echo $instance_name_alias
echo $tablet_dir
echo $tablet_logfile

# vttablet should under which shard
shard='0'

# vttablet config
# if the last character minus one is greater than 1
# which means 100,101 will be replica
# 102 will be read only
tablet_type=replica
if [[ "${uid: -1}" -gt 1 ]]; then
  tablet_type=rdonly
fi
echo $tablet_type

echo "Starting vttablet for $instance_name_alias..."

vttablet \
  $TOPOLOGY_FLAGS \
  --log_dir $VTDATAROOT/tmp \
  --log_queries_to_file $VTDATAROOT/tmp/$tablet_logfile \
  --tablet-path $instance_name_alias \
  --tablet_hostname "" \
  --init_keyspace $keyspace \
  --init_shard $shard \
  --init_tablet_type $tablet_type \
  --health_check_interval 5s \
  --enable_replication_reporter \
  --backup_storage_implementation file \
  --file_backup_storage_root $VTDATAROOT/backups \
  --restore_from_backup \
  --port $port \
  --grpc_port $grpc_port \
  --service_map 'grpc-queryservice,grpc-tabletmanager,grpc-updatestream' \
  --pid_file $VTDATAROOT/$tablet_dir/vttablet.pid \
  --vtctld_addr http://$hostname:$vtctld_web_port/ \
  --disable_active_reparents &

echo "${instance_name_alias} ran"