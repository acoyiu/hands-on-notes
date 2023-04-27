# ArgoCD

https://argo-cd.readthedocs.io/en/stable/getting_started/

```sh
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# expose service
kubectl -n argocd port-forward svc/argocd-server 8080:443 --address=0.0.0.0 &

# get login password
kubectl -n argocd get secret/argocd-initial-admin-secret -o jsonpath='{.data.password}' | base64 -d
```

then open [localhost](http://localhost:8080)

## Then make directory and file in Git project

**Need add below in argo server deployment for ingress exposure**

```yaml
- command:
    - argocd-server
    - --insecure
```

**Need remove the last slash"/" in the proxy pass in down stream reverse proxy**

```sh
proxy_pass  http://192.168.0.101:80; # <- the slash in the last mush me deleted
```

## 您可能希望防止對像被 prune

```yaml
metadata:
  annotations:
    argocd.argoproj.io/sync-options: Prune=false
```

## 打開選擇性同步選項，僅同步不同步的資源

Add ApplyOutOfSyncOnly=true in manifest

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
spec:
  syncPolicy:
    syncOptions:
      - ApplyOutOfSyncOnly=true
```

<br/>

---

<br/>

## [Create Local User](https://argo-cd.readthedocs.io/en/stable/operator-manual/user-management/#create-new-user)

- 創建本地用戶時，每個用戶都需要設置額外的 RBAC 規則
- 否則它們將回退到由 argocd-rbac-cm ConfigMap 的 policy.default 字段指定的默認策略

New users should be defined in argocd-cm ConfigMap:

```yaml
# Create new user
data:
  accounts.alice: apiKey, login # apiKey - allows generating API keys, login - allows to login using UI
  accounts.alice.enabled: "false" # User is enabled by default
```

### Change User's password (ArgoCD version)

```sh
# get user list
argocd --core account list

# get user details
argocd --core account get --account alice

# if you are managing users as the admin user, <current-user-password> should be the current admin password.
# argocd account update-password --account <name> --current-password <current-user-password> --new-password <new-user-password>
argocd --core account update-password --account alice

# if flag --account is omitted then Argo CD generates token for current user
argocd --core account generate-token --account <username>
```

### Change User's password (Port-forward version), using the "argocd" service account as credentials

```sh
# create authorization for the argocd service account
kubectl create clusterrolebinding argo-cluster --clusterrole admin --serviceaccount argocd:argocd-server --namespace argocd

# login argocd CLI
argocd --port-forward --port-forward-namespace argocd login # (--plaintext) <-- for use within the cluster

# can list all user
argocd --port-forward --port-forward-namespace argocd account list

# reset user password
argocd --port-forward --port-forward-namespace argocd account update-password --account <username>
```

## To change authorization [RBAC](https://argo-cd.readthedocs.io/en/stable/operator-manual/rbac/#tying-it-all-together)

```sh
# edit rbac configmap
kubectl -n argocd edit cm argocd-rbac-cm
```

- p = role 的授權 rule
- g = 定義 group/user ming

- p, <role>, <Resource>, <Action, CRUD>, <Project>, allow

```yaml
data:
  policy.default: role:readonly
  policy.csv: |
    p, role:test-role, clusters, create, *, allow
    g, alice, role:test-role
data:
  policy.default: role:readonly
  policy.csv: |
    p, role:test-role, clusters, get, *, allow
    p, role:test-role, applications, *, */*, allow
    p, role:test-role, repositories, *, *, allow
    p, role:test-role, logs, get, */*, allow
    p, role:test-role, exec, create, */*, allow
    g, alice, role:test-role
```

## Disable admin user, 一旦創建了其他用戶，建議禁用管理員用戶

```yaml
data:
  admin.enabled: "false"
```

<br/>

---

<br/>

```sh
# create argo application
argocd app create <app-name> \
  --repo <git-path>  \
  --path <directory-of-file> \
  --dest-server https://kubernetes.default.svc \
  --dest-namespace <namespace-name> \
  --revision <branch-name> \
  --sync-policy automated \
  --auto-prune \
  --self-heal
```

## register external cluster

```sh
argocd cluster add --core --kubeconfig ./config <context-name>
```

---

## Best to prepare credentials to call

```sh
cat >> ./client.crt << EOF
file string
EOF
```

---

## Argocd command

| common flag  | Usage                                             |
| :----------- | :------------------------------------------------ |
| --core       | use k8s api rather than argocd api                |
| --kubeconfig | use specific context, only works in (cluster) api |

| Command | Sub-command | Actions         | Example:                                                                                                               |
| :------ | :---------- | :-------------- | ---------------------------------------------------------------------------------------------------------------------- |
| argocd  |             |                 |                                                                                                                        |
| └─      | account     |                 |                                                                                                                        |
|         | └─          | list            | better use with "--core"                                                                                               |
|         | └─          | get             | argocd --core account get --account admin                                                                              |
|         | └─          | get-user-info   |                                                                                                                        |
|         | └─          | generate-token  |                                                                                                                        |
|         | └─          | delete-token    |                                                                                                                        |
|         | └─          | can-i           |                                                                                                                        |
|         | └─          | update-password |                                                                                                                        |
| └─      | app         |                 |                                                                                                                        |
|         | └─          | list            | argocd app list --core --server \$SERVER --server-crt ./ca.crt --client-crt ./client.crt --client-crt-key ./client.key |
| └─      | cluster     |                 |                                                                                                                        |
|         | └─          | add             | --core --kubeconfig ./config <context-name>                                                                            |

<br/>

---

<br/>

## Bug?

```yaml
# Sometime Argo cannot create namespace
# Add below to the Application resource
clusterResourceWhitelist:
- group: '*'
  kind: '*'
```
