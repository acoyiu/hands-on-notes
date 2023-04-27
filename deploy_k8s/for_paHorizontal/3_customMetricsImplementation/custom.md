# Custom Metric Self-Implemented
 
## 1 -- 啟動 export metrics 的 golang server

```sh
# Create namespace
kubectl create ns httpserver

# Create Deployment
cd ./1_runner
./run.sh httpserver
cd ../

# Confirm can get metrics
kubectl run temp --rm -i --restart=Never --image=nginx:alpine --command -- sh -c 'curl httpserver.httpserver/metrics'
kubectl run temp --rm -i --restart=Never --image=nginx:alpine --command -- sh -c 'for i in $(seq 1 100); do sh -c "curl httpserver.httpserver/metrics"; done;'

# Check "http_requests_total"
```

<br/>

---

<br/>

## 2 -- Install Prometheus if there is not installed

```sh
# Update Chart
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

# Create Prometheus
kubectl create ns prometheus
helm -n prometheus upgrade -i prometheus prometheus-community/prometheus -f ./2_prometheus/values.yaml
```

<br/>

---

<br/>

## 3 -- Modify the prometheus to fetch metrics

```sh
# Edit Prometheus Scrape Endpoint
kubectl -n prometheus edit cm/prometheus-server
```

### [Prometheus Scrape Confif List](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#scrape_config)

```yaml
    # Scrape all the k8s's "pod" resources in namespace "httpserver"
    - job_name: httpserver
      scrape_interval: 5s
      kubernetes_sd_configs:
      - role: pod
        namespaces:
          names:
          - httpserver
```

```sh
# restart deploy
kubectl -n prometheus rollout restart deploy/prometheus-server

# Re-port-forward
kubectl -n prometheus port-forward svc/prometheus-server 8080:80
```

<br/>

---

<br/>

## 4 -- [Prometheus Adapter](https://artifacthub.io/packages/helm/prometheus-community/prometheus-adapter) Install 

Prometheus Adaptor 知道如何與 Kubernetes 和 Prometheus 通信，充當兩者之間的翻譯器，並且可以代替 metrics-server

[Prometheus Adapter for Kubernetes Metrics APIs](https://github.com/kubernetes-sigs/prometheus-adapter)

### Install by Helm

```sh
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update

helm -n prometheus upgrade -i adaptor prometheus-community/prometheus-adapter -f ./3_adaptor/values.yaml

# Check is API Aggregation OK
kubectl get APIService/v1beta1.custom.metrics.k8s.io
```


### Test the metrics

```sh
# Check all metrics
kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1"

# Get target metrics
kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1/namespaces/httpserver/pods/*/http_requests_rate"

# If for non-namespaced resource
kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta1/{object-type}/{object-name}/{metric-name...}"
```

<br/>

---

<br/>

## Create HPA

```sh
# Create HPA upon this metrics
kubectl -n httpserver apply -f ./4_hpa/hpa.yaml

# Check metrics
kubectl -n httpserver get hpa
```