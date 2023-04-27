# Sealed Secret

## [Install](https://github.com/bitnami-labs/sealed-secrets#installation)

```sh
# add to chart
helm repo add sealed-secrets https://bitnami-labs.github.io/sealed-secrets

# install helm chart
helm install sealed-secrets -n kube-system --set-string fullnameOverride=sealed-secrets-controller sealed-secrets/sealed-secrets
```

## [Install Kubeseal](https://github.com/bitnami-labs/sealed-secrets#homebrew)

## Create Sealed Secret

```sh
kubectl create secret generic secret-name --dry-run=client --from-file=secret.json -o=json | \
kubeseal \
  --controller-name=sealed-secrets-controller \
  --controller-namespace=kube-system \
  --format yaml > sealedsecret.yaml
```

```yaml
# the sealed secret will look something like this
apiVersion: bitnami.com/v1alpha1
kind: SealedSecret
metadata:
  creationTimestamp: null
  name: secret-name
  namespace: default
spec:
  encryptedData:
    secret.json: AgBtHkqoT1OTP20hJpAnc6E.....
  template:
    metadata:
      creationTimestamp: null
      name: secret-name
      namespace: default
```

## To Use the secret

```sh
# create the secret by the above yaml
kubectl apply -f sealedsecret.yaml

# Profit!
kubectl get secret <secret-name>
```
