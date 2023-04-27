# Istio

## [Download istioctl](https://istio.io/latest/zh/docs/setup/getting-started/)

```sh
# should have more resource
minikube start --memory=7000 --cpus=4
```

<br/>

---

<br/>

## Install Istio

```sh
# preflight check like linkerd
istioctl x precheck

# show istio version
istioctl version

# list all available profile
istioctl profile list

# profile 有多種， e.g.： demo，default，minimal...
istioctl install --set profile=default -y

# verify is installation ok
istioctl verify-install

# uninstall
istioctl uninstall --set revision=default # Or
istioctl uninstall --purge
```

## Install Istio as Operator

```sh
# will create ns, istio control plane, and resource "iop,io"
istioctl operator init

# check with
kubectl get io

# Operator can be easily removed
istioctl operator init
```

<br/>

---

<br/>

## Install Addons

```sh
cd the/path/to/istio/folder

# install addons kiali, jaeger, prometheus, grafana, zipkin
kubectl apply -f samples/addons
kubectl apply -f samples/addons/extras # zipkin

# expose by direct service port-forward
kubectl -n istio-system port-forward svc/kiali 20001:20001 --address=0.0.0.0 &

# 或 by istio command
istioctl dashboard kiali
istioctl dashboard grafana
istioctl dashboard prometheus
istioctl dashboard jaeger
istioctl dashboard zipkin
```

<br/>

---

<br/>

## Inject to namespace

```sh
# kubectl label namespace <namespace-name> istio-injection=enabled
kubectl label namespace default istio-injection=enabled

# Or by yaml edit
kubectl edit ns/default

# with the following, linkerd use annotation, istio use labels
metadata:
  labels:
    istio-injection: enabled
```

### Example Application

```sh

```

<br/>

---

<br/>

## Available command

| Available Commands: | 描述                                                                       |
| :------------------ | :------------------------------------------------------------------------- |
| admin               | Manage control plane (istiod) configuration                                |
| analyze             | Analyze Istio configuration and print validation messages                  |
| bug-report          | Cluster information and log capture support tool.                          |
| completion          | Generate the autocompletion script for the specified shell                 |
| dashboard           | Access to Istio web UIs (Same as port-forward svc)                         |
| help                | Help about any command                                                     |
| install             | Applies an Istio manifest, installing or reconfiguring Istio on a cluster. |
| kube-inject         | Inject Istio sidecar into Kubernetes pod resources                         |
| manifest            | Commands related to Istio manifests                                        |
| operator            | Commands related to Istio operator controller.                             |
| profile             | Commands related to Istio configuration profiles                           |
| proxy-config        | Retrieve information about proxy configuration from Envoy [kube only]      |
| proxy-status        | Retrieves the synchronization status of each Envoy in the mesh [kube only] |
| tag                 | Command group used to interact with revision tags                          |
| uninstall           | Uninstall Istio from a cluster                                             |
| upgrade             | Upgrade Istio control plane in-place                                       |
| validate            | Validate Istio policy and rules files                                      |
| verify-install      | Verifies Istio Installation Status                                         |
| version             | Prints out build version information                                       |
