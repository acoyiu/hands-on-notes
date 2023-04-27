# Jaeger on Helm

## Install Helm Repo

```sh
# install repo
helm repo add jaegertracing https://jaegertracing.github.io/helm-charts

# check is install success
helm search repo jae
```

## Install Jaeger Helm Chart

```sh
kubectl create ns jaeger

# Unlimited resource ( Will use Up to 6 GB memory !!! )
helm -n jaeger install jaeger jaegertracing/jaeger

# for GitOps pattern
helm -n jaeger install jaeger jaegertracing/jaeger --dry-run --debug > ./jaeger.yaml
```

## Port-forward Jaeger

```sh
# forward dashboard
kubectl -n jaeger port-forward svc/jae-jaeger-query 20080:80 --address=0.0.0.0 &

# forward trace endpoint
kubectl -n jaeger port-forward svc/jae-jaeger-collector 20081:14268 --address=0.0.0.0 &
```

## Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: jaeger-ingress
  namespace: jaeger
spec:
  ingressClassName: public
  rules:
    - host: jae-jaeger-collector.jaeger
      http:
        paths:
          - backend:
              service:
                name: jae-jaeger-collector
                port:
                  number: 14268
            path: /
            pathType: Prefix
    - host: jaeger.dev.ppwi
      http:
        paths:
          - backend:
              service:
                name: jae-jaeger-query
                port:
                  number: 80
            path: /
            pathType: Prefix

```
