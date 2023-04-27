# Linkerd is a set of control CLI like Helm

## step 0: check version

```sh
kubectl version --short
```

---

## step 1: install [linkerd CLI](https://linkerd.io/2.11/getting-started/#step-1-install-the-cli)

```sh
curl --proto '=https' --tlsv1.2 -sSfL https://run.linkerd.io/install | sh
# or by brew
brew install linkerd

# check is install ok
linkerd version
# "Server version: unavailable" is okay for the first time
```

---

## Step 2: Validate your Kubernetes cluster

```sh
linkerd check --pre
```

---

## Step 3: Install the control plane onto your cluster

```sh
# pre-install crd
linkerd install --crds | kubectl apply -f -

# install
linkerd install | kubectl apply -f -

# if with minikube Or if there are nodes using the docker container runtime and proxy-init container must run as root user.
linkerd install --set proxyInit.runAsRoot=true | kubectl apply -f -

# check install progress (a minute or two)
linkerd check
```

---

## Step 4: Install linkerd Dashboard

```sh
# install the on-cluster metrics stack
linkerd viz install | kubectl apply -f -
linkerd check

# for disable host check, to change the "-enforced-host" into '.*', or add the target host name
kubectl -n linkerd-viz edit deploy/web

# start linkerd dashboard
linkerd viz dashboard --address=0.0.0.0 &

# Or by port forward *(need disable host check like above)
kubectl -n linkerd-viz port-forward svc/web 8080:8084 --address=0.0.0.0
```

---

## Step 5: Inject Linkerd into namespace

start deploy deployments (If needed)

```sh
# add linkerd annotation in namespace yaml
k edit ns/default
```

### **Inject Linkerd into namespace**

```yaml
metadata:
  annotations:
    linkerd.io/inject: enabled
# Then restart a deployment
```

### **Inject Linkerd Requested CPU**

```yaml
metadata:
  annotations:
    config.linkerd.io/proxy-cpu-request: "0.025"
```

---

## Step 6: Exposing the linkerd monitoring via ingress

```yaml
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: web-ingress-auth
  namespace: linkerd-viz
data:
  auth: YWRtaW46JGFwcjEkbjdDdTZnSGwkRTQ3b2dmN0NPOE5SWWpFakJPa1dNLgoK
---
# apiVersion: networking.k8s.io/v1beta1 # for k8s < v1.19
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: web-ingress
  namespace: linkerd-viz
  annotations:
    nginx.ingress.kubernetes.io/upstream-vhost: $service_name.$namespace.svc.cluster.local:8084
    nginx.ingress.kubernetes.io/configuration-snippet: |
      proxy_set_header Origin "";
      proxy_hide_header l5d-remote-ip;
      proxy_hide_header l5d-server-id;
    nginx.ingress.kubernetes.io/auth-type: basic
    nginx.ingress.kubernetes.io/auth-secret: web-ingress-auth
    nginx.ingress.kubernetes.io/auth-realm: "Authentication Required"
    nginx.ingress.kubernetes.io/rewrite-target: /$1
spec:
  ingressClassName: <小心::ingressclass-name>
  rules:
    - host: localhost
      http:
        paths:
          - path: /prometheus(?:/|$)(.*)
            pathType: Prefix
            backend:
              service:
                name: prometheus
                port:
                  number: 9090
          - path: /(.*)
            pathType: Prefix
            backend:
              service:
                name: web
                port:
                  number: 8084
```

---

## Tap / stat

```sh
# Tap
linkerd viz tap deploy/deploy-name
linkerd viz tap deploy/deploy-name --to svc/service-name
linkerd viz tap ns/test
linkerd viz tap ns/test --to ns/prod

# Stat
linkerd viz stat deployments -n test
linkerd viz stat namespaces --from ns/namespace-name
```

---

## If Uninstall

```sh
# To remove Linkerd Viz
linkerd viz uninstall | kubectl delete -f -

# To remove the control plane, run:
linkerd uninstall | kubectl delete -f -
```

# Install the Linkerd-Jaeger extension

https://linkerd.io/2.12/tasks/distributed-tracing/

```sh
# install
linkerd jaeger install | kubectl apply -f -

# check
linkerd jaeger check

# Explore Jaeger
linkerd jaeger dashboard
# Or
kubectl -n linkerd-jaeger port-forward svc/jaeger 16686:16686 --address=0.0.0.0
```
