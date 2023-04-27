# Ingress Acess Auth

## 1: First create user with password

```sh
sudo apt install apache2-utils

# -c=創建，會刪除原有紀錄
htpasswd -c ./auth <username>

# 新增 user
htpasswd -B ./auth <username>
```

---

## 2: Create K8s resource "secret"

```sh
# kubectl create secret generic <name-of-secret> --from-file=<path-to-file>
kubectl create secret generic nginx-auth --from-file=./auth
# -n if needed
```

---

## 3: Create content for exposing to ingress

```sh
kubectl run sample --image=nginx:alpine
kubectl expose po/sample --name sample --type=ClusterIP --port=8080 --target-port=80
```

save below as yaml ing.yaml:

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: ingress-with-auth
  annotations:
    nginx.ingress.kubernetes.io/auth-type: basic                     # type of authentication
    nginx.ingress.kubernetes.io/auth-secret: nginx-auth              # name of the secret that contains the user/password definitions
    nginx.ingress.kubernetes.io/auth-realm: "Msg said why need auth" # message to display with an appropriate context why the authentication is required
spec:
  rules:
    - http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: sample
                port:
                  number: 8080
```

```sh
# create ingress resource
kubectl apply -f ./ing.yaml

# start tunnel if needed (minikube), need recreate ing.yaml
minikube addons enable ingress
sudo minikube tunnel
```

---

## 4: Try auth

### 1: Open in browser at [localhost:80](http://localhost)

### 2: Call by "curl"

```sh
curl -v http://localhost/ -u 'username:password'
```

---

## More kinds of [Authentication](https://kubernetes.github.io/ingress-nginx/examples/auth/client-certs/)
