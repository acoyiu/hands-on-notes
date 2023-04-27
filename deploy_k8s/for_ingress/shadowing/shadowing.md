# Ingress Shadowing

## Start First Service

```sh
install ingress for minikube
minikube addons enable ingress

# start service 1
kubectl create ns nginx
kubectl config set-context --current --namespace nginx
kubectl create deploy nginx1 -r=1 --port=80 --image=nginx
kubectl expose deploy/nginx1 --type=ClusterIP --port=8080 --target-port=80
kubectl create ing nginx1 --rule=localhost/=nginx1:8080

# start tunneling for minikube
minikube tunnel

# test ingress ins working
curl localhost
```

## Start Second Service

```sh
kubectl create deploy nginx2 -r=1 --port=80 --image=nginx
kubectl expose deploy/nginx2 --type=ClusterIP --port=8080 --target-port=80
```

## Edit Configmap for ingress to inject code snippet

```sh
kubectl config set-context --current --namespace ingress-nginx

# the configmap for the nginx-ingress configuration may very in different distro
kubectl -n ingress-nginx edit cm/ingress-nginx-controller
```

### Add the following into config map

```yaml
data:
  http-snippet: |
    upstream mirror_target {
      server nginx2.nginx:8080;
    }

    split_clients "$date_gmt" $mirror_servers {
      100%    mirror_target;
    }
```

### Edit ingress to add split funtionality

```sh
kubectl delete ing/nginx1
```

#### Remember to Escape the \"$" sign if using bashscript

```sh
cat <<EOF | kubectl apply -f -
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: nginx1
  namespace: nginx
  annotations:
    nginx.ingress.kubernetes.io/configuration-snippet: |
      # You can repeat this setting N times to mirror the network traffic N times.
      mirror /mirror;

    nginx.ingress.kubernetes.io/server-snippet: |
      location = /mirror {

        # Make endpoint is not callable externally
        internal;

        # Does not print the log of mirrored requests
        access_log off;

        proxy_set_header Host \$mirror_servers;
        proxy_pass http://\$mirror_servers\$request_uri;
      }
spec:
  ingressClassName: nginx
  rules:
    - host: localhost
      http:
        paths:
          - backend:
              service:
                name: nginx1
                port:
                  number: 8080
            path: /
            pathType: Prefix
EOF
```