# Argo Events

[installation Guide](https://argoproj.github.io/argo-events/installation/)

Argo Workflow Dashboard can also view argo event resources

<br/>

---

<br/>

## Install Operator by k8s yaml

```sh
# Create the namespace
kubectl create namespace argo-events

# install the Operator
kubectl apply -f https://raw.githubusercontent.com/argoproj/argo-events/stable/manifests/install.yaml

# install validating admission controller
kubectl apply -f https://raw.githubusercontent.com/argoproj/argo-events/stable/manifests/install-validating-webhook.yaml

# Deploy the eventbus
kubectl apply -n argo-events -f https://raw.githubusercontent.com/argoproj/argo-events/stable/examples/eventbus/native.yaml

# the sensor will create a pod for the K8s api calling in order to archieve the handler
# therefore need authrize it to do so by RBAC
kubectl create clusterrolebinding allow-webhook-sensor --clusterrole cluster-admin --serviceaccount argo-events:default
```

[Can also install to specific namespace](https://argoproj.github.io/argo-events/installation/#namespace-installation)

<br/>

---

<br/>

## Install Operator by Helm

```sh
# Add the helm repo
helm repo add argo https://argoproj.github.io/argo-helm

# Helm Chart may not update as fast as official repository
helm install argo-events argo/argo-events
```

<br/>

---

<br/>

## Argo EventSource

[Argo Documentations](https://argoproj.github.io/argo-events/concepts/architecture/)

[argo structure](!./../argo-events-architecture.png)

```sh
# create hook & sensor (handler)
kubectl -n argo-events apply -f ./1_eventSource_webhook.yaml
kubectl -n argo-events apply -f ./2_source_webhook.yaml
kubectl -n argo-events apply -f ./3_source_webhook_to_wf.yaml

# Try call webhook
SVC="webhook-event-source-eventsource-svc"
kubectl run caller --image=nginx:alpine --restart=Never --rm -i --command -- curl $SVC.argo-events:11000/endpoint
```

```sh
kubectl -n argo-events delete -f ./1_eventSource_webhook.yaml
kubectl -n argo-events delete -f ./2_source_webhook.yaml
kubectl -n argo-events delete -f ./3_source_webhook_to_wf.yaml

k -n default delete po/nginx --force
```

## delete installation

```sh
kubectl delete -f https://raw.githubusercontent.com/argoproj/argo-events/stable/manifests/install.yaml

kubectl delete -f https://raw.githubusercontent.com/argoproj/argo-events/stable/manifests/install-validating-webhook.yaml

kubectl delete -n argo-events -f https://raw.githubusercontent.com/argoproj/argo-events/stable/examples/eventbus/native.yaml

```
