export ETCD_SERVER="localhost:2379"

export hostname=$(hostname -f)

export TOPOLOGY_FLAGS="--topo_implementation etcd2 --topo_global_server_address $ETCD_SERVER --topo_global_root /vitess/global"
export vtctld_web_port=15000
export VTDATAROOT="${PWD}/vtdataroot"
export grpc_port=15999
export cell="zone1"

mkdir -p "${VTDATAROOT}"
mkdir -p "${VTDATAROOT}/tmp"
mkdir -p "${VTDATAROOT}/etcd"
mkdir -p "${VTDATAROOT}/backups"