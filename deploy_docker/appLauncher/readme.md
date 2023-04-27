# App/MS Launcher

### Required project(ms) structure

can reference to the "\_template" directory

```sh
/project1
  /Dockerfile             # (for Docker/K8s Deployment)
  /_dev                   # (for localhost Docker Dev Environment)
    docker-compose.yml
    prebuild.sh           # prebuild.sh is Optional
    postbuild.sh           # postbuild.sh is Optional
  /_helm                  # (for Docker/K8s Deployment)
    /Chart.yaml
    /values.yaml          # Remeber to contain "spec.imagePullSecrets" property
    /templates            # k8s resource yaml file
```

<br/>

---

<br/>

## To run with examples, copy projects from "\_example" to project root

Note that "prebuild.sh" in every "\_dev" directories will be ran BEFORE docker script

Note that "postbuild.sh" in every "\_dev" directories will be ran AFTER docker script

(remember to change the sh file being executable in advance)

(\_runMode/local restart will not re-run the prebuild.sh and postbuild.sh)

```
Directory "dependencies" will be ran all before all projects in the directory,
also with the sequence of prebuild -> docker-compose -> postbuild
```

<br/>

---

<br/>

## Local development via Docker/Docker-Compose

```sh
_runMode/local up
# which will start all in dev for docker by
# loop through all the directories to
# start run the docker-compose.yml inside the project

_runMode/local restart
# restart all docker-compose

_runMode/local restart [...project_name_to_restart]
# restart specific compose only

_runMode/local down
# shutdown all docker
```

<br/>

---

<br/>

# K8s Deployment with Helm Chart in current k8s context

Loop build image, push to registry and upgrade helm chart

Prerequisite:

- "docker login": for later docker image building, tagging and pushing
- [kubectl](https://kubernetes.io/zh/docs/tasks/tools/): for Creating Image Pull Secret
- [helm](https://github.com/helm/helm/releases): for Injecting k8s's secret into Charts

<br/>

---

<br/>

### 1: docker login

```sh
docker login --username=[user_name] [registry]
```

<br/>

### 2: K8s create "secret" for "imagePullSecrets"

```sh
kubectl create secret docker-registry [k8s-secret-name] \
    --docker-server=[registry_url] \
    --docker-username=[user_name] \
    --docker-password=[password] \
    --docker-email=[optional_email]
```

<br/>

### 3: Check "helm" command is installed

```sh
command -v helm

# if no, install it
wget -c https://get.helm.sh/helm-v3.8.2-linux-amd64.tar.gz
tar xf helm-v3......
```

<br/>
 
## Preworks

```yaml
# "spec.imagePullSecrets" property should be existed in values.yaml
app:
  ...
spec:
  ...
  imagePullSecrets: k8s-secret-name # <- will be replaced if in launcher's command param
  fromRegistry: 192.168.0.170:34000 # <- will be replaced if in launcher's command param
  image: project1
  tag: d-April-30                   # <- will be replaced if in launcher's command param
```

## Commands

```sh
./_runMode/deploy [registry_url] [k8s-secret-name] [tagName?] [nocache?] [imagePullPolicy?]

# example:

# deploy with target registry and k8s's secret name, tag will as date string like "d-May-04"
./_runMode/deploy 192.168.0.170:34000 k8s-secret-name

# deploy with specify tag name
./_runMode/deploy 192.168.0.170:34000 k8s-secret-name 'monday-build'

# build with no docker cache mode
./_runMode/deploy 192.168.0.170:34000 k8s-secret-name '' nocache

# for case like minikube, which need to use
minikube image load [image_name]
./_runMode/deploy 192.168.0.170:34000 k8s-secret-name '' '' 'IfNotPresent'
```
