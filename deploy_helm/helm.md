## create project

```sh
helm create myapp
```

## run project

```sh
helm upgrade --install my-app-01 myapp
kubectl get pods
```

## expose of demo package

```sh
kubectl port-forward my-app-01-myapp-\*\*\* 9090:80
```

## Chart.yaml structure

<https://helm.sh/docs/topics/charts/>

| File                  | Optional? | Usage                                             |
| :-------------------- | --------: | :------------------------------------------------ |
| Chart.yaml            |      必需 | 定義 chart 資訊，包刮名稱、版本、敘述...          |
| values.yaml           |      必需 | 負責提供 yaml 需要的參數，在 templates 中可被引用 |
| templates/            |      必需 | 負責存放 + 定義 k8s yaml 檔案                     |
| LICENSE               |    可選擇 | 一份文檔紀錄 License 信息                         |
| README.md             |    可選擇 | 一份文檔紀錄介绍信息，跟 git 的 README.md 相同    |
| requirements.yaml     |    可選擇 | 定義 chart 依賴關係(方法一)                       |
| charts/               |    可選擇 | 定義 chart 依賴關係(方法二)                       |
| templates/helpers.tpl |    可選擇 | 可以將 chart.yml 或 value.yml 加工後變成新變數    |
| templates/NOTES.txt   |    可選擇 | 一份文檔，通常被用於顯示 install 後被帶入的參數值 |

<br/>

---

<br/>

## CLI

Pls refer to [kubectl.md](../deploy_k8s/kubectl.md)

<br/>

---

<br/>

## Anchor in Yaml

```yaml
app:
  replicaCount: &replicaNum 1

autoscaling:
  minReplicas: *replicaNum # this will reference to the above value "1"
```

<br/>

---

<br/>


## Values对象是从values.yaml文件和用户提供的文件传进模板的。默认为空

<br/>

---

<br/>


## Built-in Variable: Release

| Release           | Usage                                               |
| :---------------- | :-------------------------------------------------- |
| Release.Name      | release名称                                         |
| Release.Namespace | 版本中包含的命名空间(如果manifest没有覆盖的话)      |
| Release.IsUpgrade | 如果当前操作是升级或回滚的话，该值将被设置为true    |
| Release.IsInstall | 如果当前操作是安装的话，该值将被设置为true          |
| Release.Revision  | 此次修订的版本号。安装时是1，每次升级或回滚都会自增 |
| Release.Service   | 该service用来渲染当前模板。Helm里始终Helm           |

```sh
# example:
{{ .Release.Name }}
```

<br/>

---

<br/>


## Built-in Variable: Template

| Template：        | 包含当前被执行的当前模板信息                                        |
| :---------------- | :------------------------------------------------------------------ |
| Template.Name     | 当前模板的命名空间文件路径 (e.g. mychart/templates/mytemplate.yaml) |
| Template.BasePath | 当前chart模板目录的路径 (e.g. mychart/templates)                    |

<br/>

---

<br/>


## Built-in Variable: Chart

Chart.yaml文件是chart必需的。包含了以下字段：

从 v3.3.2，不再允许额外的字段。推荐的方法是在 annotations 中添加自定义元数据。

```yaml
apiVersion: chart API 版本 （必需）
name: chart名称 （必需）
version: 语义化2 版本（必需）
kubeVersion: 兼容Kubernetes版本的语义化版本（可选）
description: 一句话对这个项目的描述（可选）
type: chart类型 （可选）
keywords:
  - 关于项目的一组关键字（可选）
home: 项目home页面的URL （可选）
sources:
  - 项目源码的URL列表（可选）
dependencies: # chart 必要条件列表 （可选）
  - name: chart名称 (nginx)
    version: chart版本 ("1.2.3")
    repository: （可选）仓库URL ("https://example.com/charts") 或别名 ("@repo-name")
    condition: （可选） 解析为布尔值的yaml路径，用于启用/禁用chart (e.g. subchart1.enabled )
    tags: # （可选）
      - 用于一次启用/禁用 一组chart的tag
    import-values: # （可选）
      - ImportValue 保存源值到导入父键的映射。每项可以是字符串或者一对子/父列表项
    alias: （可选） chart中使用的别名。当你要多次添加相同的chart时会很有用
maintainers: # （可选）
  - name: 维护者名字 （每个维护者都需要）
    email: 维护者邮箱 （每个维护者可选）
    url: 维护者URL （每个维护者可选）
icon: 用做icon的SVG或PNG图片URL （可选）
appVersion: 包含的应用版本（可选）。不需要是语义化，建议使用引号
deprecated: 不被推荐的chart （可选，布尔值）
annotations:
  example: 按名称输入的批注列表 （可选）
```

```sh
# example:
{{ .Chart.Name }}-{{ .Chart.Version }}
```

<br/>

---

<br/>

## **Helm 會自動加載所有 .tpl 的檔案 ！！**

然後就可以在 template 中使用 "include"

```yaml
{{- define "ms.selectorLabels" -}}
app: 'mc2'
for-service: hello
{{- end }}
```

```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  selector:
    matchLabels:
      {{- include "ms.selectorLabels" . | nindent 6 }}
```

<br/>

---

<br/>

## (:=) assign value

```sh
# Declare
{{- $count := .Values.replicaCount }}

# using
replicas: {{ $count }}
```

<br/>

---

<br/>

## if/else， 用来创建条件语句

```yaml
# values.yaml
autoscaling:
  enabled: true
  replicaCount: 6
```

```yaml
# Deployment.yaml
apiVersion: apps/v1
kind: Deployment
spec:
  {{- if not .Values.autoscaling.enabled }}
  replicas: 1
  {{- else if gt (int .Values.autoscaling.replicaCount) 5 }}
  replicas: 10
  {{- else }}
  replicas: 3
  {{- end }}
```

<br/>

---

<br/>

## with， 攞整個 Key （with key child structure）

```yaml
# values.yaml
deployAnno:
  use: ms-abc
  example: hello
```

```yaml
# deploy.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  {{- with .Values.deployAnno }}
  annotations:
    {{- toYaml . | nindent 4 }}
  {{- end }}
```

<br/>

---

<br/>

## range， 提供"for each"类型的循环

```yaml
# values.yaml
containers:
  - resources: {}
    image: nginx
    imagePullPolicy: Always
    name: pod-{{ .Values.app.serviceName }}-1
    ports:
      - containerPort: 80
        protocol: tcp
  - resources: {}
    image: nginx
    imagePullPolicy: Always
    name: pod-{{ .Values.app.serviceName }}-2
    ports:
      - containerPort: 180
        protocol: udp
```

```yaml
# deploy.yaml
apiVersion: apps/v1
kind: Deployment

metadata:
  name: deploy-{{ .Chart.Name }}

spec:
  replicas: 1

  selector:
    matchLabels:
      app: hi

  template:
    metadata:
      labels:
        app: hi
    spec:
      restartPolicy: Always
      containers:
        {{- range .Values.containers }}
        - resources: {}
          image: {{ .image }}
          imagePullPolicy: {{ .imagePullPolicy }}
          name: {{ .name }}
          {{- with .ports }}
          ports:
            {{- toYaml . | nindent 12 }}
          {{- end }}
        {{- end }}
```

<br/>

---

<br/>

## Chart Hook 生命周期的某些点进行干预

| 可用的钩子    | usage                                        |
| :------------ | :------------------------------------------- |
| pre-install   | 模板渲染之后，Kubernetes资源创建之前执行     |
| post-install  | 所有资源加载到Kubernetes之后执行             |
| pre-delete    | Kubernetes删除之前，执行删除请求             |
| post-delete   | 所有的版本资源删除之后执行删除请求           |
| pre-upgrade   | 模板渲染之后，资源更新之前执行一个升级请求   |
| post-upgrade  | 所有资源升级之后执行一个升级请求             |
| pre-rollback  | 模板渲染之后，资源回滚之前，执行一个回滚请求 |
| post-rollback | 所有资源被修改之后执行一个回滚请求           |
| test          | 调用Helm test子命令时执行 ( test文档)        |

[詳細使用 How to use](https://helm.sh/zh/docs/topics/charts_hooks/)

<br/>

---

<br/>

## Attention

### Update of Configmap

***用 SHA256 hash 去查 configmap 有無變，有既就會導致 sum 值有變，引致 yaml 有變，重啓 pod***

**係 POD 既 annotations，唔係 Deployment 既 ！！！！！**

```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  template:
    metadata:
      annotations:
        checksum/config: {{ include (print $.Template.BasePath "/configmap.yaml") . | sha256sum }}
```

### Add timestamp to deployment to trigger redeploy everytime

***用 timestamp 去令 "pod" 的 yaml 字段有變化，helm 就會發佈新的 deployment***

```yaml
apiVersion: apps/v1
kind: Deployment
spec:
  ...
  template:
    metadata:
      annotations:
        timestamp: {{ now | quote }}
```

<br/>

---

<br/>

# Real Example

[Swagger Doc with chart](../doc_swaggerUi/by_helm/readme.md)
