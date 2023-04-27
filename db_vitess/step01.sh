source ./env.sh

ETCDCTL_API=2

# Check that etcd is not already running
curl "http://${ETCD_SERVER}" > /dev/null 2>&1 && echo "etcd is already running. Exiting." && exit

# Start the etcd process
etcd --enable-v2=true --data-dir "${VTDATAROOT}/etcd/" --listen-client-urls "http://${ETCD_SERVER}" --advertise-client-urls "http://${ETCD_SERVER}" > "${VTDATAROOT}"/tmp/etcd.out 2>&1 &
PID=$!
echo $PID
echo $PID > "${VTDATAROOT}/tmp/etcd.pid"

# wait for etcd up
sleep 5

echo "add key '/vitess/global'"
etcdctl --endpoints "http://${ETCD_SERVER}" mkdir /vitess/global &

echo "add key /vitess/$cell"
etcdctl --endpoints "http://${ETCD_SERVER}" mkdir /vitess/$cell &

echo "etcd start done..."