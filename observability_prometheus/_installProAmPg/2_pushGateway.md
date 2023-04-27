
# Helm PushGateway

## Create custom namespace for the pushgateway

```sh
kubectl create ns pushgateway
```

## create deployment and service for the pushgateway

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: pushgateway-deployment
  namespace: pushgateway
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
          name: pushgateway
          ports:
            - containerPort: 9091
---
apiVersion: v1
kind: Service
metadata:
  name: pushgateway-service
  namespace: pushgateway
spec:
  type: ClusterIP
  selector:
    usage: pushgateway
  ports:
    - port: 9091
      targetPort: 9091
```

### Endpoint of pushing metrics

| spec     | value                                          |
| :------- | :--------------------------------------------- |
| method   | PUT                                            |
| endpoint | http://pushgateway:9091/metrics/job/<job-name> |

<br/>

---

<br/>

## View the dashboard

```sh
kubectl -n pushgateway port-forward svc/pushgateway-service 9091:9091 &
```

### For pushing to Gateway, open below

[Go Push to Gateway](../others_customGoPushToGateway/run.md)