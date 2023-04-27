# kubeadm cli for (CKA)

## [All K8s Terms](https://kubernetes.io/zh-cn/docs/reference/glossary/?architecture=true&core-object=true&extension=true&fundamental=true&networking=true&operation=true&security=true&storage=true&tool=true&user-type=true&workload=true)

## Install of kubeadm

- If using your own equipment you will have to disable swap on every node
- and ensure there is only one network interface.

[installation process refer to ./_kubeadm_install](./_kubeadm_install/master/install_master.md)

---

## CP (Control Plane) 所需要的 4 個進程

| 進程               | 用途                                                   |
| :----------------- | :----------------------------------------------------- |
| api-server         | k8s 的 API server，接收所有 k8s 的操作（包括 kubectl） |
| scheduler          | 負責與 kubelet 協調如何生成 Pod                        |
| controller-manager | 負責監聽 k8s 資源變化，及執行後續操作                  |
| etcd               | 儲存 k8s 的狀態（State），需要 HA ~~化~~               |

## Worker Node (Worker) 共需要 2 個進程

| 進程       | 用途                                                   |
| :--------- | :----------------------------------------------------- |
| kubelet    | worker 的主要成分，worker node 同 CP 的 scheduler 溝通 |
| kube-proxy | 負責 worker node 的所有 network 的分發                 |

### kubelet 可以舊過 API service one minor version： the kubelet running 1.7.0 should be fully compatible with a 1.8.0 API server, but not vice versa.

<br/>

---

<br/>

# Index of Required Commands

- [apt](#apt)
- [systemctl](#systemctl)
- [daemon.json](#daemonjson)
- [kubeadm](#kubeadm)
- [kubectl (CKA)](#kubectl-cka)
- [etcdctl](#etcdctl)
- kubelet
- [Process upgrape](#process-upgrade)

<br/>

---

<br/>

# apt (Advance Package Tool)

| CMD                                 | Usage                                                                                      |
| :---------------------------------- | :----------------------------------------------------------------------------------------- |
| apt update                          | 更新軟體庫清單                                                                             |
| apt upgrade                         | 升級系統軟體                                                                               |
| apt install <package-name>          | 安裝應用軟體                                                                               |
| apt remove <package-name>           | 移除應用程式                                                                               |
| apt purge <package-name>            | 移除應用程式及所有設定檔                                                                   |
| apt list -a                         | 列出所有可安裝版本                                                                         |
| apt list --installed                | 列出所有已安裝版本                                                                         |
| apt show <package-name>             | 列出應用程式的詳細資訊，最新版本，並包括相關軟體套件的相依項目                             |
| apt search                          | 搜尋特定的應用程式                                                                         |
| apt edit-sources                    | 編輯包含所有軟體庫的 /etc/apt/sources.list 文件設定檔。 實際執行會打開文字編輯器的快捷指令 |
| apt-mark hold/unhold <package-name> | 固定版本不 update/upgrade                                                                  |
| apt-key add -                       | 每个发布的deb包，都是通过密钥认证的，apt-key用来管理打開包的密钥。                         |

---

## apt repository list

[repo structure explained](https://www.cnblogs.com/kelamoyujuzhen/p/9728260.html)

```sh
# list available version
apt list kubelet -a | head -n 20

# install with specific version
apt install kubelet=1.23.12-00

# 通用的目錄
ls -al /etc/apt/

# 發行商的任意目錄，類似 nginx 的 conf.d
ls -al /etc/apt/sources.list.d
```

---

## apt edit-sources

repository structure:

```sh
deb http://archive.ubuntu.com/ubuntu/ focal main restricted
```

| Section 1: Type of src | Section 2: Url | Section 3: distro name | Section 4: component type                 |
| :--------------------- | :------------- | :--------------------- | :---------------------------------------- |
| deb or deb-src         | http://...     | e.g. focal, xenial     | main / restricted / universe / multiverse |

<br/>

---

<br/>

# systemctl

[systemctl vs init vs service](https://segmentfault.com/a/1190000038458363)

- initd 是最初的进程管理方式 (init system - PID 1)
- service 是 inits 的另一种实现
- systemd 则是一种取代 initd (/service) 的解决方案

Systemd 并不是一个命令，而是一组命令，涉及到系统管理的方方面面。

| ctl         | usage                          |
| :---------- | :----------------------------- |
| systemctl   | 用于管理系统, Systemd 的主命令 |
| hostnamectl | 用于查看当前主机的信息         |
| localectl   | 用于查看本地化设置             |
| timedatectl | 用于查看当前时区设置           |

```sh
# List all daemon
systemctl --type=service
```

| Other Commands             | Usage                                |
| :------------------------- | :----------------------------------- |
| start UNIT                 | Start (activate) one or more units   |
| stop UNIT                  | Stop (deactivate) one or more units  |
| reload UNIT                | Reload one or more units             |
| restart UNIT               | Start or restart one or more units   |
| enable [UNIT... / PATH...] | 開機自動啓動                         |
| disable UNIT               | 開機不自動啓動                       |
| daemon-reload              | reload systemd manager configuration |
| daemon-reexec              | restart systemd manager              |

Others pls reference to **systemctl -h**

<br/>

---

<br/>

# service

<br/>

---

<br/>

# daemon.json

## 簡單黎講：就係 Docker 既 config

- docker安装后默认没有daemon.json这个配置文件，需要进行手动创建
- 默认路径：/etc/docker/daemon.json

## 用途例如：配置 registry 私库相关的参数

```json
{
    "registry-mirrors": [
        "https://d8b3zdiw.mirror.aliyuncs.com"
    ],
    "insecure-registries": [
        "https://ower.site.com"
    ]
}
```

[reference](https://blog.csdn.net/u013948858/article/details/79974796)

[官方的配置地址](https://docs.docker.com/engine/reference/commandline/dockerd/#options)

## k8s docker runtime 需要修改 cgroup

[Kubernetes 推荐使用 systemd](https://blog.51cto.com/riverxyz/2537914) 来代替 cgroupfs因为systemd是Kubernetes自带的cgroup管理器, 负责为每个进程分配cgroups,但docker的cgroup driver默认是cgroupfs,这样就同时运行有两个cgroup控制管理器,当资源有压力的情况时,有可能出现不稳定的情况

```sh
# check docker cgroup engine
docker info

# on kubelet machine
cat >> /etc/docker/daemon.json << EOF
{
  "exec-opts": ["native.cgroupdriver=systemd"]
}
EOF
```

<br/>

---

<br/>

# kubeadm

## K8s file directory

- /etc/kubernetes
  - admin.conf
  - controller-manager.conf
  - kubelet.conf
  - scheduler.conf
  - addons
    - storage-provisioner.yaml
    - storageclass.yaml
  - manifests
    - etcd.yaml
    - kube-apiserver.yaml
    - kube-scheduler.yaml
    - kube-controller-manager.yaml

## ***manifests 裏面的其實就係 4 大核心模塊的 Pod 的 Yaml ！***

---

## kubeadm related command

| kubeadm Actions               | Usage                                                                                                              |
| :---------------------------- | :----------------------------------------------------------------------------------------------------------------- |
| [kubeadm init](#kubeadm-init) | 用于搭建控制平面节点 [steps of init](https://kubernetes.io/zh-cn/docs/reference/setup-tools/kubeadm/kubeadm-init/) |
| [kubeadm join](#kubeadm-join) | 用于搭建工作节点并将其加入到集群中                                                                                 |
| kubeadm upgrade               | 用于升级 Kubernetes 集群到新版本                                                                                   |
| kubeadm reset                 | 用于恢复通过 kubeadm init 或者 kubeadm join 命令对节点进行的任何变更                                               |

| kubeadm Resources                 | Usage                                                                                              |
| :-------------------------------- | :------------------------------------------------------------------------------------------------- |
| [kubeadm config](#kubeadm-config) | 如果你使用了 v1.7.x 或更低版本的 kubeadm 版本初始化你的集群，则使用 kubeadm upgrade 来配置你的集群 |
| kubeadm version                   | 用于打印 kubeadm 的版本信息                                                                        |
| kubeadm token                     | 用于管理 kubeadm join 使用的令牌                                                                   |
| kubeadm certs                     | 用于管理 Kubernetes 证书                                                                           |
| kubeadm kubeconfig                | 用于管理 kubeconfig 文件                                                                           |

- etcdctl
- systemctl
- service
- kubectl (Rest API Direct call)

---

kubeadm CMD example:

```yaml
init
  Start-Cluster:                   kubeadm init --config=kubeadm-config.yaml --upload-certs
join
  Join-Cluster: >-
                                   kubeadm join \
                                     k8scp:6443
                                     --token 4kp60v.t8aoqtzamenauxbu \
                                     --discovery-token-ca-cert-hash sha256:f8adbdb365c72ec151c972a3a89ab002f00758f9e868d64a425588e2fb21e74f
upgrade:
  cp:
    upgrade-preflight-check:       kubeadm upgrade plan
    actual-upgrade:                kubeadm upgrade apply <version> # v1.23.1
  worker:
    update-worker:                 kubeadm upgrade node
config:
  print-k8s-init-details:          kubeadm config print init-defaults
token:
  list-all-token:                  kubeadm token list
  create-token-and join-command:   kubeadm token create --print-join-command
certs:
  show-all-cert:                   kubeadm certs check-expiration
```

<br/>

---

<br/>

# kubectl (CKA)

CMD reference refer to [kubectl.md](./kubeadm.md)

## NGINX ingress controller setup

- [NGINX ingress controller](https://kubernetes.github.io/ingress-nginx/deploy/baremetal/)
- [See the available version](https://github.com/kubernetes/ingress-nginx)
- [List the available versions](https://github.com/kubernetes/ingress-nginx/releases)

```sh
# replace the version in the url to get the latest one
wget https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.3.0/deploy/static/provider/baremetal/deploy.yaml -O ingress.yaml

# create namespace and apply
kubectl create ns ingress-nginx
kubectl -n ingress-nginx apply -f ./ingress.yaml

# See the ingress's service NodePort
kubectl get services -n ingress-nginx

# example
ingress-nginx-controller             NodePort    10.102.207.233   <none>        80:31204/TCP,443:30437/TCP   11m
```

### 31204 and 30437 are the port need External Load Balancer to redirect to

---

## Rest API Direct Call <- NOT CKAD, but CKA

```sh
# prepare CA cert
ca=$(grep certificate-auth $HOME/.kube/config | awk '{print $2}')

# prepare client cert
cert=$(grep client-certificate $HOME/.kube/config | awk '{print $2}')

# prepare client key
key=$(grep client-key $HOME/.kube/config | awk '{print $2}')

# Make call to K8s API Server
curl \
  --key $key \
  --cert $cert \
  --cacert $ca \
  https://192.168.49.2:6443/api/v1/pods
  https://k8scp:6443/api/v1/pods
```

---

## Static Pod

- 不能通过 API 服务器来控制
- 不能引用其他 API 对象 （如：ServiceAccount、 ConfigMap、 Secret 等）
- kubelet 监视每个静态 Pod，those declared in node's "/etc/kubernetes/manifest/*.yaml"
- 静态static Pod 的 process（kernal）**會運作在 node（host）的 process pool 裏，而不是 container 裏**
- 即是以 host 的 process 來模擬 process isolation，所以 process 可以在 ps aux 裏 kill，然後透過 kubelet 重啟

**kubelet 會創建（/刪除）所有在或不在 "/etc/kubernetes/manifest/*.yaml" 這個 directory 内的資源**

```sh
# SSH login into node, and create yml in specific folder

cat <<EOF >/etc/kubernetes/manifests/static-web.yaml
apiVersion: v1
kind: Pod
metadata:
  name: static-web
  labels:
    role: myrole
spec:
  containers:
    - name: web
      image: nginx
      ports:
        - name: web
          containerPort: 80
          protocol: TCP
EOF
```

---

## unsafe TLS in metrics server

```sh
kubectl -n kube-system edit deployment metrics-serve
```

```yaml
spec:
  containers:
    - args:
    - --cert-dir=/tmp
    - --secure-port=4443
    - --kubelet-insecure-tls #<-- Add this line
    - --kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname #<--May be needed
  image: k8s.gcr.io/metrics-server/metrics-server:v0.3.7
```

<br/>

---

<br/>

# Upgrade K8s

### 用於調試 kubeadm 集群的命令

| command                            | Usage              |
| :--------------------------------- | :----------------- |
| kubeadm config print init-defaults | 檢查集群啟動配置   |
| systemctl status kubelet           | kubelet 進程的健康 |

### 需要記住的特殊路徑

| Path                       | Usage                                                                                                 |
| :------------------------- | :---------------------------------------------------------------------------------------------------- |
| /etc/kubernetes/admin.conf | 存取 cluster 的 admin（master） context                                                               |
| /etc/kubernetes/pki        | 所有用於 CA 的證書，需要多 Control Plane (Master Node) 的時候需複製此 directory                       |
| /etc/kubernetes/manifests  | (安裝 kubeadm 後就會有) 靜態資源（Static Resource）的 directory，K8s 會自動生成此 directory 下的 yaml |

## A: Files needed to backup for K8s

三樣 backup k8s 時必須保持的東西：

- 1: /var/lib/etcd/snapshot.db (需先用 etcdctl snapshot save)
- 2: /root/kubeadm-config.yaml
- 3: /etc/kubernetes/pki （各種 cert）

**kubeadm 一般會將 4 大 process 以 static pod 形式運行 (Host 的 process)**

## # Process Upgrade

- 1: Cordon + Drain target node
- 2：Backup ETCD
- 3: apt update tools
- 4: kubeadm upgrage apply (--etcd-upgrade)
- 5: Upgrade kubectl & kubelet
- 6: Uncordon node

<br/>

---

<br/>

## 1 -- Drain/Condon node

```sh
# Condon + Drain target node

kubectl drain <node-name> --force --ignore-daemonsets (--delete-emptydir-data 報錯才需要)
```

<br/>

---

<br/>

## 2 -- etcdctl 

```yaml
list-file-without-ls:    echo *
export-api-version:      export ETCDCTL_API=3
check-health:            etcdctl endpoint health
list-member:
  list:                  etcdctl member list
  table-format:          etcdctl member list -w table
snapshot:
  save:                  etcdctl snapshot save    ./snapshot.db
  restore:               etcdctl snapshot restore ./snapshot.db
  status:                etcdctl snapshot status  ./snapshot.db  # Gets backend snapshot status of a given file
```

## Run etcd under docker

```sh
docker run \
  -dit \
  -p 2379:2379 \
  -p 2380:2380 \
  --name etcd \
  quay.io/coreos/etcd:v3.4.0 \
  /usr/local/bin/etcd \
  --name s1 \
  --data-dir /etcd-data \
  --listen-client-urls http://0.0.0.0:2379 \
  --advertise-client-urls http://0.0.0.0:2379 \
  --listen-peer-urls http://0.0.0.0:2380 \
  --initial-advertise-peer-urls http://0.0.0.0:2380 \
  --initial-cluster s1=http://0.0.0.0:2380 \
  --initial-cluster-token tkn \
  --initial-cluster-state new \
  --log-level info \
  --logger zap \
  --log-outputs stderr
```

### 2.1, backup /var/lib/etcd/snapshot.db in host/static-pod

```sh
# list etcd pod
kubectl get pods -n kube-system

# enter etcd pod
kubectl exec -it kube-system <etcd-pod-name> -- sh

# cannot use ls, list directory with
echo *

# 使用etcd v3的版本時，需要設置環境變數ETCDCTL_API=3
export ETCDCTL_API=3

# backup etcd to directory
etcdctl snapshot save ./snapshot.db

# copy outside the container
kubectl cp <namespace>/<pod>:<remote-path> 
```

### 2.2, backup if outside pod

```sh
# https://ithelp.ithome.com.tw/articles/10240323

# get etcd-cli
apt install etcd-client
export ETCDCTL_API=3

# show credentials path
sudo cat /etc/kubernetes/manifests/etcd.yaml
```

| 需要的資訊              | 值                                  | 用途                |
| :---------------------- | :---------------------------------- | :------------------ |
| --advertise-client-urls | https://192.168.64.27:2379          | DB 的位置           |
| --cert-file             | /etc/kubernetes/pki/etcd/server.crt | ETCD 的 cert        |
| --key-file              | /etc/kubernetes/pki/etcd/server.key | ETCD 的 private key |
| --trusted-ca-file       | /etc/kubernetes/pki/etcd/ca.crt     | CA 的 public key    |

```sh
# copy credentials
sudo cp -r /etc/kubernetes/pki/etcd ~/etcd
sudo chown -R ubuntu:ubuntu ./etcd

# member list
param="--endpoints=192.168.64.27:2379 --cert=$HOME/etcd/server.crt --key=$HOME/etcd/server.key --cacert=$HOME/etcd/ca.crt"
# cat /etc/kubernetes/manifests/etcd.yaml | grep advertise | sed -n '2p' | grep -ro 'https.*' -

# 3 main api for backup
etcdctl endpoint health $param
etcdctl member list $param
etcdctl snapshot save ./snapshot.db $param
```

## 3: Update kube-* cli command by apt

```sh
# list all available version
apt list kubeadm -a | head -n 10

# install new version by apt
apt install kubeadm=1.24.2-00

# check current version
kubeadm version
```

## 4: kubeadm plan & upgrage apply

```sh
# show suggested upgrade command
kubeadm upgrade plan

kubeadm upgrade apply <version-number> --etcd-upgrade=false
# --etcd-upgrade=true , default 會一拼 upgrade 埋 etcd
```

## 5: kubectl & kubelet

```sh
apt install kubelet=1.24.2-00
apt install kubectl=1.24.2-00

# check version
kubelet --version
kubectl version --short

# check kubelet is running
systemctl start kubelet # debug 會考
systemctl enable kubelet # debug 會考
systemctl status kubelet # debug 會考

# show kubelet logs
journalctl -eu kubelet
```


## 6: Uncordon node

```sh
kubectl uncordon <node-name>
```
