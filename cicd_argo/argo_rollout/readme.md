# Argo Rollout

## [Installation](https://argoproj.github.io/argo-rollouts/installation/)

```sh
kubectl create namespace argo-rollouts
kubectl apply -n argo-rollouts -f https://github.com/argoproj/argo-rollouts/releases/latest/download/install.yaml

# install seperated CRD as it's not included
kubectl apply -k https://github.com/argoproj/argo-rollouts/manifests/crds\?ref\=stable
```

## [kubectl argo rollouts Plugin](https://argoproj.github.io/argo-rollouts/installation/#kubectl-plugin-installation)

```sh
# With Linux
curl -LO https://github.com/argoproj/argo-rollouts/releases/latest/download/kubectl-argo-rollouts-linux-amd64

chmod +x ./kubectl-argo-rollouts-linux-amd64

sudo mv ./kubectl-argo-rollouts-linux-amd64 /usr/local/bin/kubectl-argo-rollouts
```

## Sample Rollout

```sh
kubectl apply -f ./canary/app.yaml 
# https://raw.githubusercontent.com/argoproj/argo-rollouts/master/docs/getting-started/basic/rollout.yaml
# https://raw.githubusercontent.com/argoproj/argo-rollouts/master/docs/getting-started/basic/service.yaml

# watch the upgrade
kubectl argo rollouts get rollout rollouts-demo --watch
```

## Start Canary and Promoting

```sh
# upgrade the rollout by applying new yaml with image value changed :: rollouts-demo:blue -> rollouts-demo:yellow
kubectl apply -f ./canary/app.yaml

# expose svc
kubectl port-forward svc/rollouts-demo 9090:80

# Promoting = proceed to next step of upgrade
kubectl argo rollouts promote rollouts-demo

# The promote command also supports the ability to skip all remaining steps and analysis with the --full flag
```

## Abort rollout

```sh
# OR change partial pod to yellow image by argo rollouts command :: rollouts-demo:yellow -> rollouts-demo:red
kubectl argo rollouts set image rollouts-demo rollouts-demo=argoproj/rollouts-demo:red

# As first stage is made manual, can jam that awaiting for rollback
kubectl argo rollouts abort rollouts-demo

# Restore the rollout status to normal by applying the old image as latest state, the status will become "healthy" again
kubectl argo rollouts set image rollouts-demo rollouts-demo=argoproj/rollouts-demo:yellow
```

## Start with Blue-Green

```sh
kubectl apply -f ./blueGreen/app.yaml
kubectl argo rollouts get rollout rollout-bluegreen --watch

kubectl argo rollouts set image rollout-bluegreen rollouts-demo=argoproj/rollouts-demo:yellow

kubectl port-forward svc/rollout-bluegreen-active 9090:80
kubectl port-forward svc/rollout-bluegreen-preview 9091:80

kubectl argo rollouts promote rollout-bluegreen
```
