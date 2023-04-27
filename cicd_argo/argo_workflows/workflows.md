# Argo Workflows

[installation Guide](https://argoproj.github.io/argo-workflows/quick-start/)

[installation CLI](https://github.com/argoproj/argo-workflows/releases/tag/v3.3.10)

<br/>

---

<br/>

## Install Argo Workflows

[Check release here](https://github.com/argoproj/argo-workflows/releases/tag/v3.3.10)

```sh
# create namespace
kubectl create namespace argo

# get version from list and update here
ARGO_WORKFLOWS_VERSION=3.3.10
kubectl apply -n argo -f https://github.com/argoproj/argo-workflows/releases/download/v${ARGO_WORKFLOWS_VERSION}/install.yaml

# Add RBAC auth to argo-server sa
kubectl create clusterrolebinding allow-argo-server --clusterrole cluster-admin --serviceaccount argo:argo-server

# Patch argo-server authentication for bypass the UI login
kubectl patch deployment \
  argo-server \
  --namespace argo \
  --type='json' \
  -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/args", "value": [
  "server",
  "--auth-mode=server"
]}]'

# Port-forward the UI
kubectl -n argo port-forward deployment/argo-server 2746:2746
```

remember to open with [https](https://localhost:2746/)

<br/>

---

<br/>

## Install Argo CLI

[Sample Page from above](https://github.com/argoproj/argo-workflows/releases/tag/v3.3.10)

- Argo command 同 helm 類似，會根據 k8s config 選擇 context
- workflow = Job
- workflow Template = CronJob —from

<br/>

---

<br/>

## Structure of Workflow

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Workflow
metadata:
  generateName: sample-wf-                  # Resource Prefix，所有這個 workflow 的資源都會 prefix
  labels:
    custom-labels: "shown-in-dashboard"     # custom labels，所有這個 workflow 的資源都會有這些 labels
spec:
  entrypoint: the-template-name             # 指定 "spec.template" 裏面哪一個 template 為第一個執行
  templates:
    - name: the-template-name               # template ~= drone's pipeline ~= K8s's pod
      container:
        image: alpine
        command: [sh, -c, 'sleep 120']      # can use "sleep infinity" to debug
```

<br/>

---

<br/>

## Submit workflow by yaml file

**Argo is like helm, workflows are namespace scoped**

```sh
# kn default if needed
kn default

# create workflow (job)
argo submit -n argo ./workflow_0_sample.yaml
```

see [workflow dashboard](https://localhost:2746/workflows/argo?limit=50)

## Basic Workflow

```sh
# create workflow in context default namespace
argo submit ./workflow_1_whalesay.yaml

# list all workflow in all namespace
argo list -A

# get the first workflow details
argo get hello-world-......

# get logs of the containers
argo logs hello-world-......
```

## Flow, input & param

```sh
# for the dag-B to works in "default" namespace
kubectl -n default create clusterrolebinding admin-nginx --clusterrole=admin --serviceaccount=default:default

# create workflow
argo submit ./workflow_2_multiple_wf.yaml -p 'workflow-param-1="abcd"'
```

## Artifacts

### Use Hardwired Artifacts, Otherwise, should be using PVC with storage class

```sh
# hard wired
argo submit ./workflow_3_1_artifactsHardwired.yaml

# volume
argo submit ./workflow_3_2_artifactsVolume.yaml
```

## WorkflowTemplate

```sh
# need to configure artifact storage
argo template create ./workflow_4_template.yaml

# Create workflow by CLI, can chec resource type by "k api-resources | grep argo"
argo submit --from wftmpl/workflow-template-example -p 'message="good"' --entrypoint steps-flow
```

## [Or create in dashboard](https://localhost:2746/workflow-templates/)

<br/>

---

<br/>

## Available Template

### template 就是最小的工作單元，template 内可以用 step/dig 再指定使用其他 template
### template 係 step/dig 就不可以有實際工作負載（container/pod）

---

| 可用 template 類型 | 用途                                    |
| :----------------- | :-------------------------------------- |
| steps              | 引用其他 template，決定 template 的順序 |
| dig                | 申明此 template 的運行 dependence       |
| container          | 用於生成 container，負責負載            |
| script             | 用於生成 container，負責負載            |
| resource           | 用於生成 k8s resource                   |
| suspend            | stop until time expire                  |
| http               | send http request                       |

<br/>

---

<br/>

## Argo CLI to remember

| argo CLI     | Usage                               |
| :----------- | :---------------------------------- |
| **General**  | ----------------------------------- |
| version      | Show Command Version                |
| auth         | manage authentication settings      |
| **Workflow** | ----------------------------------- |
| list         | list all workflow                   |
| get          | get specific workflow               |
| logs         | log a workflow                      |
| delete       | delete a workflow                   |
| archive      | manage the workflow archive         |
| cron         | manage cron workflow                |
| resubmit     | resubmit one or more workflows      |
| resume       | resume zero or more workflows       |
| retry        | retry zero or more workflows        |
| submit       | submit a workflow                   |
| suspend      | suspend zero or more workflow       |
| template     | manipulate workflow templates       |
| terminate    | terminate workflow(s) immediately   |
| version      | print version information           |
| wait         | waits for workflows to complete     |
| watch        | watch a workflow until it completes |
| **Template** | ----------------------------------- |
| create       | create a template                   |
| **Params**   | ----------------------------------- |
| -n           | = namespace                         |
| -A           | = all namespace                     |
| --watch      | watch resource                      |
| **Others**   | ----------------------------------- |
| lint         | validate manifests                  |

<br/>

---

<br/>

## [Argo Manifast Example](https://github.com/argoproj/argo-workflows/tree/master/examples)

## [Flow Loogs](https://argoproj.github.io/argo-workflows/walk-through/loops/)

## [Flow Conditions](https://argoproj.github.io/argo-workflows/walk-through/conditionals/)

## [Timeouts](https://argoproj.github.io/argo-workflows/walk-through/timeouts/)

## [Flow Retry Failed](https://argoproj.github.io/argo-workflows/walk-through/retrying-failed-or-errored-steps/)

## [Flow Exit handlers](https://argoproj.github.io/argo-workflows/walk-through/exit-handlers/)
