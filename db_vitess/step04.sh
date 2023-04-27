source ./env.sh

echo "Starting Vttablet..."

# start vttablets for keyspace commerce
for i in 100 101 102; do
  TABLET_UID=$i ./_mysqlctl-up.sh
	TABLET_UID=$i KEYSPACE=${KEYSPACE} ./_vttablet-up.sh
done