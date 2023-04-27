# Make PushGateway

## Set k8s configuration

```sh
kubectl create ns pushgateway
kubectl config set-context --current --namespace pushgateway
```

##  Make deployment

```sh
cat <<EOF | kubectl apply -f -
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pushgateway-deployment
spec:
  replicas: 1
  selector:
    matchLabels:
      usage: pushgateway
  template:
    metadata:
      labels:
        usage: pushgateway
    spec:
      restartPolicy: Always
      containers:
        - image: prom/pushgateway
          imagePullPolicy: IfNotPresent
          name: pushgateway
          ports:
            - containerPort: 9091
---
apiVersion: v1
kind: Service
metadata:
  name: pushgateway-service
spec:
  type: ClusterIP
  selector:
    usage: pushgateway
  ports:
    - port: 9091
      targetPort: 9091
EOF

# Check is success
kubectl get pods -w

# Check gateway dashboard
kubectl port-forward svc/pushgateway-service 9091:9091 --address=0.0.0.0 &
```

If see things in [dashboard](http://localhost:9091), then can push metrics to:

pushgateway-service.pushgateway:9091/metrics/job/<job_name>
