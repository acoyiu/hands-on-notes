source ./env.sh

# Add the CellInfo description for the cell.
# If the node already exists, it's fine, means we used existing data.

echo "add $cell CellInfo"

set +e

# vtctl VtctldCommand
# https://vitess.io/docs/15.0/reference/vtctldclient-transition/legacy_shim/#vtctlclient-vtctldcommand

# 用舊既 command 既 argument 去行新既 command

# AddCellInfo 會係 etcd 入邊寫入 cell 的資訊，但因有 role 原因，無 auth 既 etcdctl 唔會睇到唔會

vtctl $TOPOLOGY_FLAGS VtctldCommand AddCellInfo \
  --root /vitess/$cell \
  --server-address "${ETCD_SERVER}" \
  $cell

set -e