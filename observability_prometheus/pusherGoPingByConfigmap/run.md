# Run

```sh
ns=default
```

```sh
# Only run this also if is for restart
kubectl -n $ns \
  delete \
  cm/gofiles \
  svc/svc-go-runner \
  deploy/deploy-go-runner --force && \

# => Start Here <= Need update below if directory structure changed
kubectl -n $ns \
  create cm gofiles \
  --from-file=./go.mod \
  --from-file=./main.go && \

kubectl -n $ns \
  get cm gofiles -o yaml > ./_gofiles.cm.yaml && \
  
sed -i '/namespace/d' ./_gofiles.cm.yaml && \
sed -i '/creationTimestamp/d' ./_gofiles.cm.yaml && \
sed -i '/resourceVersion/d' ./_gofiles.cm.yaml && \
sed -i '/uid/d' ./_gofiles.cm.yaml && \

kubectl -n $ns \
  apply -f ./
```

```sh
# logs
kubectl -n $ns logs -f po/deploy-go-runner-    (tap)

# check service
kubectl -n $ns run tmp --image=nginx:alpine --restart=Never --rm -i --command -- curl svc-go-runner:80/metrics
```

```sh
# delete
kubectl -n $ns delete -f ./
```
