# Weave Scope

## [Install](https://www.weave.works/docs/scope/latest/installing/#k8s)

```sh
# install
kubectl apply -f https://github.com/weaveworks/scope/releases/download/latest/k8s-scope.yaml

# Open Scope in Your Browser, The URL is: http://localhost:4040.
kubectl port-forward -n weave "$(kubectl get -n weave pod --selector=weave-scope-component=app -o jsonpath='{.items..metadata.name}')" 4040:4040 --address='0.0.0.0' &
```

## Uninstall

```sh
kubectl delete -f ./scope.yaml
```
