# Run

```sh
# start deployment with
K8sNamespace=test
./run.sh $K8sNamespace

# get pods
kubectl -n $K8sNamespace get pods

# enter pods
kubectl -n $K8sNamespace exec -it $(kubectl get pods --sort-by '{.metadata.creationTimestamp}' | grep file-runner | grep Running | tail -n 1 | awk '{print $1}') -- sh -c 'cd /unzip/app && sh'
```

## Update "command" in runner.yaml
