# Helm RabbitMq

## Install Chart

```sh
# create ns
kubectl create ns rabbitmq

# get chart
helm repo add bitnami https://charts.bitnami.com/bitnami

# install
helm -n rabbitmq install rabbit bitnami/rabbitmq --set persistence.storageClass=default \
--set persistence.size=21Gi # <- this is for alicloud
```

**If can not Run, may due to the PV has no access right on host**

can describe pv to check details, if so, ssh host and chmod 777 for the mounted directory

## Expose MQ

```sh
# export message entry
kubectl -n rabbitmq port-forward svc/rabbit-rabbitmq 5672:5672

# export dashboard
kubectl -n rabbitmq port-forward svc/rabbit-rabbitmq 15672:15672

# username: user, password =
kubectl -n rabbitmq get secret rabbit-rabbitmq -o jsonpath="{.data.rabbitmq-password}" | base64 -d
```

## Expose by ingress (dashboard) and NodePort

```yaml
# expose dashboard
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ing-rabbitmq
spec:
  rules:
    - host: localhost
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: rabbit-rabbitmq
                port:
                  number: 15672
```

```yaml
# expose message input port
apiVersion: v1
kind: Service
metadata:
  name: rabbit-msg-input
spec:
  type: NodePort
  selector:
    app.kubernetes.io/instance: rabbit
    app.kubernetes.io/name: rabbitmq
  ports:
    - nodePort: 30123
      port: 5672
      targetPort: 5672
```
