# ArgoCD Add Dashboard Support

https://argo-cd.readthedocs.io/en/stable/operator-manual/web_based_terminal/

## 1: Set the exec.enabled key to "true" on the argocd-cm ConfigMap

```sh
# change to argocd namespace
kubectl config set-context --current --namespace argocd

# edit cm/argocd-cm
kubectl edit cm/argocd-cm
```

### Add the below to "data"

```yaml
data:
  exec.enabled: "true"
```

## 2: Patch the argocd-server Role/ClusterRole (clustere wide) to allow argocd-server to exec into pods

```sh
# edit the clusterrole for cluster-wide operation, use role for specific ns
kubectl edit clusterrole argocd-server
```

```yaml
# add the following authoriztion
- apiGroups:
  - ""
  resources:
  - pods/exec
  verbs:
  - create
```

## 3: Need added Nginx websocket allow upgrade config if need to expose by ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  annotations:
    nginx.ingress.kubernetes.io/configuration-snippet: |
      proxy_set_header Upgrade "websocket";
      proxy_set_header Connection "Upgrade";
```

## 3: Add RBAC rules to allow your users to create the exec resource, i.e.

### [Example:](https://argo-cd.readthedocs.io/en/stable/operator-manual/rbac/#tying-it-all-together)

| RBAC Resources and Actions: |                                 |
| :-------------------------- | :------------------------------ |
| Resources:                  | clusters                        |
|                             | projects                        |
|                             | applications                    |
|                             | repositories                    |
|                             | certificates                    |
|                             | accounts                        |
|                             | gpgkeys                         |
|                             | logs                            |
|                             | exec                            |
| Actions:                    |                                 |
|                             | get                             |
|                             | create                          |
|                             | update                          |
|                             | delete                          |
|                             | sync                            |
|                             | override                        |
|                             | action/<group/kind/action-name> |

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: argocd-rbac-cm
  namespace: argocd
data:
  policy.default: role:readonly
  policy.csv: |
    p, role:org-admin, applications, *, */*, allow
    p, role:org-admin, clusters, get, *, allow
    p, role:org-admin, repositories, get, *, allow
    p, role:org-admin, repositories, create, *, allow
    p, role:org-admin, repositories, update, *, allow
    p, role:org-admin, repositories, delete, *, allow
    p, role:org-admin, logs, get, *, allow
    p, role:org-admin, exec, create, */*, allow
```
