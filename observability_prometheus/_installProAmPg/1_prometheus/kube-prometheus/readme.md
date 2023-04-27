# Helm install Prometheus (kube-prometheus)

```sh
# create namespaces and add repo
kubectl create ns prometheus

helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update bitnami
```

## Helm value.yaml

```yaml
global:
  storageClass: ""
prometheus:
  retention: 10d
  replicaCount: 2
  persistence:
    enabled: true
    size: 20Gi
```

## Create prometheus instance

```sh
# helm install prometheus in its namespace
helm -n prometheus upgrade -i prometheus bitnami/kube-prometheus

# port-forward prometheus
kubectl port-forward --namespace prometheus svc/prometheus-kube-prometheus-prometheus 9090:9090 --address='0.0.0.0' &
```

Please Read the Logs: Prometheus can be accessed via port "9090" on the following DNS name from within your cluster

### need update prometheus config in deployment yaml (linkerd) for more than 12 hrs retention time

```sh
--storage.tsdb.retention.time=6h
# -> 3d
```
