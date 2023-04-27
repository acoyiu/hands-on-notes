# External Prometheus's AlertManager (Example of Linkerd-viz Prometheus)

## Expose Linkerd's Prometheus

```sh
kubectl -n linkerd-viz port-forward service/prometheus 9090:9090 --address=0.0.0.0 &
```

<br/>

---

<br/>

## create configmap for the alert rules

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-alert
  namespace: linkerd-viz
data:
  alert.yml: |-
    groups:
      - name: ping_to_service
        rules:
          - alert: InstanceDown
            expr: ping_time_to_service{endpoint!=""} > 80
            for: 10s
            labels:
              severity: urgent
            annotations:
              summary: >
                1: {{ $labels.instance }}
                2: {{ $labels.job }}
                3: {{ $labels.endpoint }}
```

<br/>

---

<br/>

## step 2: adding alert rule by editing prometheus deployment

```sh
kubectl -n linkerd-viz edit deploy/prometheus
```

Add belows:

```yaml
spec:
  volumes:
    ... # import volume as below
    - name: prometheus-alert
      configMap:
        name: prometheus-alert
  containers:
    ... # mount configmap as file as below
    volumeMounts:
      - mountPath: /etc/prometheus/alert_rules.yml
        name: prometheus-alert
        subPath: alert.yml
```

<br/>

---

<br/>

## Step 3: Create AlertManager

```sh
# create new ns
kubectl create ns alertmanager
```

```yaml
# alertmanager.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: alertmanager
  name: alertmanager
  namespace: alertmanager
spec:
  replicas: 1
  selector:
    matchLabels:
      app: alertmanager
  template:
    metadata:
      labels:
        app: alertmanager
    spec:
      containers:
        - image: prom/alertmanager
          name: alertmanager
          resources: {}
---
apiVersion: v1
kind: Service
metadata:
  name: alertmanager
  namespace: alertmanager
  labels:
    app: alertmanager
spec:
  type: ClusterIP
  selector:
    app: alertmanager
  ports:
    - port: 9093
      protocol: TCP
      targetPort: 9093
```

```sh
# apply alertmanager
kubectl apply -f ./alertmanager.yaml

# expose (port-forward) service in background
kubectl -n alertmanager port-forward svc/alertmanager 9093:9093 --address=0.0.0.0 &
```

<br/>

---

<br/>

## Step 4: Connect Prometheus and Alert Manager

```sh
kubectl -n linkerd-viz edit cm/prometheus-config
```

```yaml
# Add below section in to Prometheus config's configmap

# Alerting
alerting:
  alertmanagers:
  - static_configs:
    - targets:
      - alertmanager.alertmanager:9093
```

```sh
# restart prometheus deployment
kubectl -n linkerd-viz rollout restart deploy/prometheus
```

then you should see the alerts being triggered will show in alertmanager
