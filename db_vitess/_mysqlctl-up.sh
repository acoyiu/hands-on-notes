source ./env.sh

mkdir -p $VTDATAROOT/backups

# mysql process id
uid=$TABLET_UID

# mysql port for vttablet
mysql_port=$[17000 + $uid]

# get formatted directory for current instance
printf -v tablet_dir 'vt_%010d' $uid

# reuse existing file for db if existed
action="init"
  if [ -d $VTDATAROOT/$tablet_dir ]; then
  echo "Resuming from existing vttablet dir: $VTDATAROOT/$tablet_dir"
  action='start'
fi

mysqlctl \
  --log_dir $VTDATAROOT/tmp \
  --tablet_uid $uid \
  --mysql_port $mysql_port \
  $action
