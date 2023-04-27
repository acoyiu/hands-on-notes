# Install by Helm

https://github.com/drone/charts

```sh
helm repo add drone https://charts.drone.io
helm repo update
```

https://github.com/drone/charts/blob/master/charts/drone/docs/install.md

```sh
vim ./drone-values.yaml
```

```yaml
env:
  DRONE_SERVER_PROTO: http
  DRONE_SERVER_HOST: localhost
  DRONE_RPC_SECRET: 81e04d83a6054b464f5c5b13365578fd
  DRONE_GITEA_CLIENT_ID: b138b14a-872c-41ac-9667-0f26406c5b21
  DRONE_GITEA_CLIENT_SECRET: V7iuU4QEbGglQW1xKbBwdA1Rx0KsE5UF4fKcoY9ZVziQ
  DRONE_GITEA_SERVER: https://gitea.greatics.net
  DRONE_GIT_ALWAYS_AUTH: true
  # If using github
  DRONE_GITHUB_CLIENT_ID: 7d281e31c86af925c3b8
  DRONE_GITHUB_CLIENT_SECRET: 587726b9d598c4bccd1wef04eb8f772982a0db5e
```

```sh
kubectl create ns drone
helm install --namespace drone drone drone/drone -f drone-values.yaml

# port forward service
kubectl port-forward svc/drone 4000:8080 --address=0.0.0.0 &
```

## install Runner

https://docs.drone.io/runner/kubernetes/installation/

### Create role and binding if needed

```sh
k create ns drone-runner
k config set-context --current --namespace drone-runner
```

```yaml
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  namespace: default
  name: drone
rules:
  - apiGroups:
      - ""
    resources:
      - secrets
    verbs:
      - create
      - delete
  - apiGroups:
      - ""
    resources:
      - pods
      - pods/log
    verbs:
      - get
      - create
      - delete
      - list
      - watch
      - update

---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: drone
  namespace: default
subjects:
  - kind: ServiceAccount
    name: default
    namespace: default
roleRef:
  kind: Role
  name: drone
  apiGroup: rbac.authorization.k8s.io
```

### install runner

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: drone
  namespace: drone-runner
  labels:
    app.kubernetes.io/name: drone
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: drone
  template:
    metadata:
      labels:
        app.kubernetes.io/name: drone
    spec:
      containers:
        - name: runner
          image: drone/drone-runner-kube:latest
          ports:
            - containerPort: 3000
          env:
            - name: DRONE_RPC_HOST
              value: 192.168.0.15:28000
            - name: DRONE_RPC_PROTO
              value: http
            - name: DRONE_RPC_SECRET
              value: 81e04d83a6054b464f5c5b13365578fd
```
