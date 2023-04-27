# Vitess Setup

## Install apt package

```sh
# other platform check: https://vitess.io/docs/16.0/get-started/local/
sudo apt install -y mysql-server etcd curl
```

## First start etcd

```sh
# start etcd process and create key/dir
./step01.sh

# check is process running
pgrep etcd

# check etcd has vitess keys (localhost)
etcdctl ls /vitess
```

## Step2: create Cell (Vitess Physical Division) Info

```sh
./step02.sh

# check is registered in local
./step02_check_local.sh
```

## Step3: Start vtctld process

```sh
./step03.sh

# check the process is up
lsof -i -n -P | grep -i listen | grep -e "\(my\)\|\(vt\)"
```

## Step4: start mysql and join as vttable

```sh
KEYSPACE=alpha ./step04.sh

# check processes are started
pgrep mysql
pgrep vttable

# check proceses with port
lsof -i -n -P | grep -i listen | grep -e "\(my\)\|\(vt\)"
```

## Step5: Update the namespace durability (Not sure is necessary)

```sh
vtctldclient --server localhost:15999 SetKeyspaceDurabilityPolicy --durability-policy=semi_sync alpha
```

## Step6: Start vtorc process (a vitess's self-trying-heal process)

```sh
./step06.sh
```

## Step7: Apply VSchema and Schame

```sh
./step07.sh

# when finish this step, then one of the replica will become master
```

## Step8: Start Vtgate process

```sh
./step08.sh

# check is vtgate up?
lsof -i -n -P | grep -i listen | grep -e "\(my\)\|\(vt\)"
```
