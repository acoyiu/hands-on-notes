docker network create cass

docker run -dit \
  --name cassandra1 \
  --network cass \
  --restart=always \
  -e CASSANDRA_SEEDS=cassandra1 \
  -e CASSANDRA_CLUSTER_NAME=alpha-cluster \
  -e CASSANDRA_ENDPOINT_SNITCH=GossipingPropertyFileSnitch \
  -e CASSANDRA_DC=hongkong \
  -e CASSANDRA_RACK=rack1 \
  -e CASSANDRA_LISTEN_ADDRESS=cassandra1 \
  -e CASSANDRA_RPC_ADDRESS=cassandra1 \
  cassandra:4.1

index=2
name="cassandra${index}"
docker run -dit \
  --name "cassandra${index}" \
  --network cass \
  --restart=always \
  -e CASSANDRA_SEEDS=cassandra1 \
  -e CASSANDRA_CLUSTER_NAME=alpha-cluster \
  -e CASSANDRA_ENDPOINT_SNITCH=GossipingPropertyFileSnitch \
  -e CASSANDRA_DC=hongkong \
  -e CASSANDRA_RACK=rack1 \
  -e CASSANDRA_LISTEN_ADDRESS=${name} \
  -e CASSANDRA_RPC_ADDRESS=${name} \
  cassandra:4.1

# # windows
# docker run -dit --name cassandra1 --network cass --restart=always -e CASSANDRA_SEEDS=cassandra1 -e CASSANDRA_CLUSTER_NAME=alpha-cluster -e CASSANDRA_ENDPOINT_SNITCH=GossipingPropertyFileSnitch -e CASSANDRA_DC=hongkong -e CASSANDRA_RACK=rack1 -e CASSANDRA_LISTEN_ADDRESS=cassandra1 -e CASSANDRA_RPC_ADDRESS=cassandra1 cassandra:4.1
# docker run -dit --name cassandra2 --network cass --restart=always -e CASSANDRA_SEEDS=cassandra1 -e CASSANDRA_CLUSTER_NAME=alpha-cluster -e CASSANDRA_ENDPOINT_SNITCH=GossipingPropertyFileSnitch -e CASSANDRA_DC=hongkong -e CASSANDRA_RACK=rack1 -e CASSANDRA_LISTEN_ADDRESS=cassandra2 -e CASSANDRA_RPC_ADDRESS=cassandra2 cassandra:4.1
# docker run -dit --name cassandra3 --network cass --restart=always -e CASSANDRA_SEEDS=cassandra1 -e CASSANDRA_CLUSTER_NAME=alpha-cluster -e CASSANDRA_ENDPOINT_SNITCH=GossipingPropertyFileSnitch -e CASSANDRA_DC=hongkong -e CASSANDRA_RACK=rack1 -e CASSANDRA_LISTEN_ADDRESS=cassandra3 -e CASSANDRA_RPC_ADDRESS=cassandra3 cassandra:4.1
# docker run -dit --name cassandra4 --network cass --restart=always -e CASSANDRA_SEEDS=cassandra1 -e CASSANDRA_CLUSTER_NAME=alpha-cluster -e CASSANDRA_ENDPOINT_SNITCH=GossipingPropertyFileSnitch -e CASSANDRA_DC=hongkong -e CASSANDRA_RACK=rack1 -e CASSANDRA_LISTEN_ADDRESS=cassandra4 -e CASSANDRA_RPC_ADDRESS=cassandra4 cassandra:4.1

# on Ubuntu 1
docker run -dit \
  --net=host \
  -p 7000:7000 \
  -p 9042:9042 \
  --name cassandra1 \
  -e CASSANDRA_CLUSTER_NAME=alpha-cluster \
  -e CASSANDRA_ENDPOINT_SNITCH=GossipingPropertyFileSnitch \
  -e CASSANDRA_DC=hongkong \
  -e CASSANDRA_RACK=rack1 \
  -e CASSANDRA_SEEDS=192.168.0.222:7000 \
  -e CASSANDRA_LISTEN_ADDRESS=192.168.0.222:7000 \
  cassandra:4.1

# on Ubuntu 2
docker run -dit \
  --net=host \
  --name cassandra1 \
  -e CASSANDRA_CLUSTER_NAME=alpha-cluster \
  -e CASSANDRA_ENDPOINT_SNITCH=GossipingPropertyFileSnitch \
  -e CASSANDRA_DC=beijing \
  -e CASSANDRA_RACK=rack1 \
  -e CASSANDRA_SEEDS=192.168.0.222:7000 \
  -e CASSANDRA_LISTEN_ADDRESS=192.168.0.15 \
  cassandra:4.1

# on Ubuntu 1 - sub
docker run -dit \
  --name cassandra2 \
  -e CASSANDRA_CLUSTER_NAME=alpha-cluster \
  -e CASSANDRA_ENDPOINT_SNITCH=GossipingPropertyFileSnitch \
  -e CASSANDRA_DC=hongkong \
  -e CASSANDRA_RACK=rack1 \
  -e CASSANDRA_SEEDS=192.168.0.222,192.168.0.15 \
  cassandra:4.1


docker run -dit \
  -p 7000:7000 \
  -p 9042:9042 \
  -v $PWD/cassandra.yaml:/etc/cassandra/cassandra.yaml \
  --name cassandra1 \
  -e CASSANDRA_ENDPOINT_SNITCH=GossipingPropertyFileSnitch \
  cassandra:4.1