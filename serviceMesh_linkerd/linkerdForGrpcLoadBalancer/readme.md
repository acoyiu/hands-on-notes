# gRPC need linkerd to as load balancer

as gRPC use http2, therefore connection level load balancing strategy is NOT working

---

## Run in dev mode

```sh
npm install
npm run server
npm run gateway # update the .env file to localhost, e.g.: localhost:8080
# the gateway will take the endpoint in .env file as gRPC connection endpoint

# then "curl localhost:7070"
```

---

## Run in Docker

in dev docker enviroment,
update .env file to local IP, e.g.: 192.168.0.170:8080
embed .env file inside docker image

```sh
# build server and run
docker build -t aco-grpc-server -f Dockerfile_server --no-cache .
docker run -dit -p 8080:8080 --name=aco-g-server aco-grpc-server

# build gateway and run
docker build -t aco-grpc-gateway -f Dockerfile_gateway --no-cache .
docker run -dit -p 7070:7070 --name=aco-g-gateway aco-grpc-gateway

# then "curl localhost:7070"

# remove container
docker kill aco-g-server aco-g-gateway
docker rm aco-g-server aco-g-gateway
```

---

## Run in K8s

in deploy enviroment, create configmap to insert into img's .env

```sh
# using minikube as example
minikube start

# install ingress, which will take few minutes
minikube addons enable ingress

# load image into minikube so that no need pull docker registry
minikube image load aco-grpc-server
minikube image load aco-grpc-gateway

# create gRPC server in K8s
kubectl apply -f ./_k8s/grpcServer.yaml

# make configmap to mount in the gateway as .env file
kubectl apply -f ./_k8s/gatewayConfigmap.yaml
kubectl apply -f ./_k8s/grpcGateway.yaml

# check is the configmap mount correctly
kubectl exec -it deployment-aco-gateway-[hash] -- cat /app/.env

# create ingress
kubectl apply -f ./_k8s/imgress.yaml

# show logs for 4 pods
kubectl logs -f [pod]

# expose of ingress
sudo minikube tunnel
```

## check the logs

```javascript
let i = 50;
while (i > 0) {
  i--;
  fetch("http://localhost:80");
}
```

### confirm load balancer of ingress to gateway is OKAY !

### but gateway to service gRPC is NOT WORKING !

---

## Install linkerd CLI

```sh
# install (ok to skip)
curl --proto '=https' --tlsv1.2 -sSfL https://run.linkerd.io/install | sh
# or by brew
brew install linkerd

# check is install ok
# "Server version: unavailable" is okay for the first time
linkerd version
```

## Try add Linkerd Service Mesh

```sh
# check version
kubectl version --short

# Validate your Kubernetes cluster
linkerd check --pre

# Install the control plane onto your cluster
linkerd install | kubectl apply -f -

# check install progress (a minute or two), pls check is control plane is up
linkerd check

## Install linkerd Dashboard, the on-cluster metrics stack
linkerd viz install | kubectl apply -f -
linkerd check
linkerd viz dashboard &
```

---

## redeploy Deployment OR inject linkerd to whole namspace

A: redeploy Deployment with config yaml edit

```sh
# uncomment spec.template.metadata.annotations

kubectl apply -f ./_k8s/grpcServer.yaml
kubectl apply -f ./_k8s/grpcGateway.yaml

kubectl get pods
kubectl logs -f [pod-hash]
kubectl get pods $PodName -o jsonpath='{.spec.containers[*].name}'
# will see more than one container inside one pod

# logs
kubectl logs -f --container pod-aco-gateway deployment-aco-gateway-[pod-hash]
kubectl logs -f --container pod-aco-server deployment-aco-server-[pod-hash]

# open ingress again and see logs
sudo minikube tunnel
```

B: OR inject linkerd to whole namspace

```sh
# kubectl edit namespace/[namespace_name]
kubectl edit ns/default
```

# add following into metadata of namespace yaml:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  annotations:
    linkerd.io/inject: enabled
```

---

# 玩完

```sh
kubectl delete -f ./_k8s/grpcServer.yaml &&
kubectl delete -f ./_k8s/gatewayConfigmap.yaml &&
kubectl delete -f ./_k8s/grpcGateway.yaml &&
kubectl delete -f ./_k8s/imgress.yaml
```
