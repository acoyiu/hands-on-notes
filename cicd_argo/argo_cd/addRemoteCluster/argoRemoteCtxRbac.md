# Adding remote cluster into ArgoCD

```sh
# when after install argo
kubectl -n argocd exec -it po/argocd-server-... -- bash

# Use the authed context from remote cluster to generate/export wanted mobile context
kubectl config view --minify --flatten

# save the credential inside the container
cat >> ./config.yaml << EOF
apiVersion: v1
kind: Config
clusters:
  - cluster:
      certificate-authority-data: ......
      server: ......
    name: target-cluster
contexts:
  - name: name-of-remote-cluster
    context:
      cluster: target-cluster
      user: service-account
users:
  - name: service-account
    user:
      token: ......
EOF

# add cluster by argocd cli
argocd cluster add --kubeconfig ./config.yaml <name-of-remote-cluster> --core
```

## Belows are what argo cli will create in target cluster

### Notice that direct create below will failed, most credentials are created and verified by domain/IP within the cert x/509 "subject" field

### which means the credentials can only be used in the environment (such as IP) where and when it be created.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: argocd-manager
  namespace: kube-system
secrets:
  - name: argocd-manager-token
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: argocd-manager-role
rules:
  - apiGroups:
      - "*"
    resources:
      - "*"
    verbs:
      - "*"
  - nonResourceURLs:
      - "*"
    verbs:
      - "*"
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: argocd-manager-role-binding
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: argocd-manager-role
subjects:
  - kind: ServiceAccount
    name: argocd-manager
    namespace: kube-system
---
apiVersion: v1
kind: Secret
type: kubernetes.io/service-account-token
metadata:
  name: argocd-manager-token
  namespace: kube-system
  annotations:
    kubernetes.io/service-account.name: argocd-manager
```