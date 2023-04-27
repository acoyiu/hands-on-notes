# kustomize

## Project Folder Structure

```sh
├── base (not necessary to be named as "base")
│   ├── kustomization.yaml
│   └── **/*.yaml
├── overlay (not necessary to be named as "overlay" neither)
│   ├── kustomization.yaml
│   └── **/*.yaml
└── .gitignore
```

## [Available keys for "kustomization.yaml"](https://kubernetes.io/zh-cn/docs/tasks/manage-kubernetes-objects/kustomization/)

```yaml
namespace: default           # 設定 resource 的 namespace

resources:                   # 設定此 project 包含哪些設定檔 
  - deployment.yaml          # 注意，一個 yaml 内可多於一個 resource 聲明
  - service.yaml

namePrefix: base-            # 讓所有 project 内 resource 的名稱添加統一前綴
nameSuffix: "-001"           # 或統一後綴 （注意，如在 overlay 内用，會變成叠加，非 replace）

commonLabels:                # 統一 label
  exampleLabel: value 
commonAnnotations:           # 統一 annotation
  exampleAnno: value

images:                      # 在 "kustomization.yaml" 修改其他 yaml 内的 image   
  - name: namePlaceholder
    newName: nginx:alpine
    newTag: 1.4.0

                             # 注意 ：如需在 overlay 中需要覆蓋 base 中的 configmap，不要用 configmap 檔案，而要用 configMapGenerator
configMapGenerator:          # 在 "kustomization.yaml 聲明 ConfigMap, 效果類似領完單獨寫的 configmap.yaml
  - name: example-cm         # 用檔案生成 cm
    files:
      - ./files/exp.json
  - name: example-cm-1       # 用 .env (key=value pair) 檔案生成 cm
    env:
      - ./.env
  - name: example-cm-2       # 用 inline 文字生成 cm
    literals:
      - FOO=Bar

secretGenerator: {}          # 同上， 但 for secret
generatorOptions: {}         # https://github.com/kubernetes-sigs/kustomize/blob/master/api/types/generatoroptions.go#L7

bases:                       # 此 project 會繼承的 yaml（project）
  - ../base                  # 因此此 project 會擁有 ../base & ../rbac
  - ../rbac                  # 兩個 kustomization.yaml 裏面聲明的所有野
                             # 注意 ! ：resource 是不可以同名地 overlay 的！如只想修改 yaml 數值， 用 patch

patchesStrategicMerge:       # 會攞指定 yaml 裏的指定名字的 resource，patch 上新指定 yaml 裏的值
- deploy.yaml                # patchesJson6902 for JSON type

vars:                        # 聲明用於引用 resource 的 variable
- name: CREDENTIAL_CM        # 可以 pod 中以下列方式引用
  objref:                    # spec.containers.[]command: ["start", "--host", "$(MY_SERVICE_NAME)"]
    kind: ConfigMap          
    name: cm-credential
    apiVersion: v1           # version 要對應目標 resource yaml 裏面的 version
# 注意 1 ：在 pod 中引用 configmap 並不用使用 "vars", 即使 kustomize 會添加 cm 的後綴 hash 和 pre/suffix
# "kustomize" 懂得為 pod 中的引用自動添加 “namePrefix” & “nameSuffix”
```

## kustomize cmd

```sh
# Show built yaml
kubect kustomize ./project/

# Build from github
kubectl kustomize https://github.com/kubernetes-sigs/kustomize.git/examples/helloWorld?ref=v1.0.6

# Build to file
kubectl kustomize ./project/ > dev.yaml
```
