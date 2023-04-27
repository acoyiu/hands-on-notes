# Metrics

## Metrics export app

```sh
# Create metrics exporter app
kubectl -n default apply -f 1.app.yaml

# Check is metrics exporter runs normally
kubectl run temp --rm -i --restart=Never --image=nginx:alpine --command -- sh -c 'curl sample-app.default/metrics'
# Or
kubectl run temp --rm -i --restart=Never --image=nginx:alpine --command -- sh -c 'for i in $(seq 1 100); do sh -c "curl sample-app.default/metrics"; done;'

# Create HPA upon this metrics
kubectl -n default apply -f 2.hpa.yaml

# Get logs of controller-manager (if in minikube)
kubectl -n kube-system logs -f kube-controller-manager-minikube

# Check metrics
kubectl -n default get hpa
```

Result:
unable to fetch metrics from custom metrics API: no custom metrics API (custom.metrics.k8s.io) registered

## Way 1: By "Prometheus Operator", "Service monitor" & Prometheus Adaptor

"Service monitor" 可以在不修改 prometheus deployment 的情況下添加 scrape path

### Install Operator and CRD

```sh
# Install and update helm chart
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# Create Prometheus helm release
kubectl create ns monitoring
helm -n monitoring upgrade -i prometheus prometheus-community/kube-prometheus-stack -f ./3.prometheus.yaml

# Create monitor to the metric exporter app
kubectl -n default apply -f 4.monitor.yaml

# Check is metrics existed
kubectl -n monitoring port-forward svc/prometheus-kube-prometheus-prometheus 8080:9090
```

### Install Adaptor

Prometheus Adaptor 知道如何與 Kubernetes 和 Prometheus 通信，充當兩者之間的翻譯器，並且可以代替 metrics-server

```sh
# Update helm
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# Install Helm release
helm -n monitoring upgrade -i prometheus-adapter prometheus-community/prometheus-adapter -f 5.adaptorValues.yaml

# Check of K8s api service resource
kubectl get APIService -l "app.kubernetes.io/name=prometheus-adapter"

# Check what metrics are available
kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1"

# To get specific metrics
kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1/namespaces/default/pods/*/http_requests_total"
kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1/namespaces/default/pods/*/http_requests_total?selector=app%3Dsample-app"

# Check metrics is successfully read
kubectl -n default get hpa

# Add loading
kubectl run temp --image=nginx:alpine
kubectl exec -it temp -- sh
for i in $(seq 1 1000); do sh -c "curl sample-app.default/metrics"; done;
```