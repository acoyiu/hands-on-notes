# In view of Project Scope

The Project should contains:

```sh
/docker-compose.yml  : For local development
/Dockerfile          : For K8s deployment
/_helm/Chart.yaml    : For K8s deployment
```

<br/>

---

<br/>

## Single service run by docker-compose

```sh
docker-compose up -d
docker-compose down
```

<br/>

---

<br/>

## Deploy to "K8s" OR through "CICD"

The Host mahine should be able to:

- Login Docker Registry
- Build Docker Image & Push
- Kubectl Create Image Pull Secret
- Helm Inject k8s's secret into Charts

<br/>

### Step 1: Login Docker Registry in advance

```sh
docker login --username=[user_name] [registry]
# with password
```

### Step 2: Docker build Image & tag to new Registry & Push

```sh
# docker build -t [project_name]:[project_version] .
docker build -t project1:d-May-1 .

# docker build -t [project_name] [registry]/[project_name]:[project_version]
docker tag project1:d-May-1 192.168.0.170:34000/project1:d-May-1

# push the img to target registry
docker push 192.168.0.170:34000/project1:d-May-1
```

### Step 3: Kubectl create imagePullSecret

```sh
kubectl create secret docker-registry secret-temp-registry \
    --docker-server=192.168.0.170:34000 \
    --docker-username=user1 \
    --docker-password=user1 \
    --docker-email=aco@aco.com
```

### Step 4: Helm install/upgrade

```sh
# should inject image pull secret [if needed]
# --> spec.template.spec.imagePullSecrets
cd _helm
helm upgrade project1 . --install --set app.replicas=3 --set spec.imagePullSecrets="secret-temp-registry"
cd ../

# uninstall
cd _helm
helm uninstall project1
cd ../
```
