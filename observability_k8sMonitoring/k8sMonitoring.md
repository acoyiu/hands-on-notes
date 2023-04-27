# K8s's Metric-Server, kubernetes-dashboard & Portainer

```sh
# Create sample deployment
kubectl create deploy trial --image=nginx:alpine --replicas=8
```

---

## Install Metric-Server in K8s

### Metric-Server in minikube

```sh
# start kubernetes
minikube start --cni calico

# Metrics Server (if in Minikube)
minikube addons list
minikube addons enable metrics-server
```

### Metric-Server in Normal K8s

```sh
# install by YAML
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Or by helm

# add repo
helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server/
# install chart, [-n kube-system] optional
helm upgrade --install metrics-server metrics-server/metrics-server -n kube-system
```

---

## Install kubernetes-dashboard

```sh
# create namespaces
kubectl create namespace dashborad

# Add kubernetes-dashboard repository
helm repo add kubernetes-dashboard https://kubernetes.github.io/dashboard/

# Deploy a Helm Release named "kubernetes-dashboard" using the kubernetes-dashboard chart
helm -n dashborad install kubernetes-dashboard kubernetes-dashboard/kubernetes-dashboard

# Get the Kubernetes Dashboard URL by running:
export POD_NAME=$(kubectl get pods -n dashborad -l "app.kubernetes.io/name=kubernetes-dashboard,app.kubernetes.io/instance=kubernetes-dashboard" -o jsonpath="{.items[0].metadata.name}")
echo https://127.0.0.1:8443/
kubectl -n dashborad port-forward $POD_NAME 8443:8443
```

---

## Install Portainer by Helm

PersistentVolume and StorageClass **are required**

```sh
# create namespaces and add helm repo
kubectl create namespace portainer
helm repo add portainer https://portainer.github.io/k8s/

# install chart
helm upgrade -i -n portainer portainer portainer/portainer

# open dashboard
kubectl port-forward -n portainer svc/portainer 6060:9443 --address=0.0.0.0 &

# OR use dashboard by NodePort
export NODE_PORT=$(kubectl get --namespace portainer -o jsonpath="{.spec.ports[1].nodePort}" services portainer)
export NODE_IP=$(kubectl get nodes --namespace portainer -o jsonpath="{.items[0].status.addresses[0].address}")
echo https://$NODE_IP:$NODE_PORT
```
