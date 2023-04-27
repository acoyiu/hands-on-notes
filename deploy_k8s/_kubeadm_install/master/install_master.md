# Install of kubeadm

## Switch to Root

```sh
# if needed
sudo passwd ubuntu
sudo passwd root

# change to root user
sudo -i
```

## install openssh-server (if-needed)

```sh
sudo apt install openssh-server
```

---

## Set repository of apt for installing kubeadm

[release](https://kubernetes.io/releases/)

[kubeadm](https://kubernetes.io/docs/setup/production-environment/tools/kubeadm/install-kubeadm/#installing-kubeadm-kubelet-and-kubectl)

```sh
# Download the Google Cloud public signing key:
sudo apt-get update
sudo apt-get install -y apt-transport-https ca-certificates curl

# Add the Kubernetes apt repository:
sudo curl -fsSLo /usr/share/keyrings/kubernetes-archive-keyring.gpg https://packages.cloud.google.com/apt/doc/apt-key.gpg

# Update apt package index, install kubelet, kubeadm and kubectl, and pin their version:
echo "deb [signed-by=/usr/share/keyrings/kubernetes-archive-keyring.gpg] https://apt.kubernetes.io/ kubernetes-xenial main" | sudo tee /etc/apt/sources.list.d/kubernetes.list

# update apt
sudo apt-get update

# latest version
sudo apt-get install -y kubelet kubeadm kubectl
sudo apt-mark hold kubelet kubeadm kubectl

# previous version
kversion=1.23.8-00
sudo apt install -y \
  "kubelet=$kversion" \
  "kubeadm=$kversion" \
  "kubectl=$kversion"
sudo apt-mark hold kubelet kubeadm kubectl
```

---

## install docker as container engine

```sh
apt update && apt upgrade -y
apt install -y docker.io

# update docker daemon and config cluster yaml
cat >> /etc/docker/daemon.json << EOF
{ "exec-opts": ["native.cgroupdriver=systemd"] }
EOF

# for activating daemon.json
systemctl restart docker; sleep 20; systemctl status docker
```

---

## Add hostname for CP node

```sh
# get ip addr
ip addr
```

```sh
cpip= ?????
echo "$cpip k8scp" >> /etc/hosts

# calico yaml for later install in "kubeadm init"
wget https://docs.projectcalico.org/manifests/calico.yaml

# record down the value(192.168.0.0/16) for later use
cat calico.yaml | grep -A 1 CALICO_IPV4POOL_CIDR
```

---

## Start the cluster Way 1

```sh
# prefetch the images needed
kubeadm config images pull

kubeadm init \
  --control-plane-endpoint="k8scp:6443" \
  --service-cidr="192.168.0.0/16" \
  --upload-certs \
  | tee kubeadm-init.out

  # --kubernetes-version="1.23.8"
  # apt cli version relative to target k8s target can only be 1 in minor version
```

## Start the cluster Way 2

```sh
# Prepare cluster config yaml
cat >> kubeadm-config.yaml << EOF
apiVersion: kubeadm.k8s.io/v1beta2
kind: ClusterConfiguration
kubernetesVersion: 1.24.6           #<-- specify the k8s version to use
controlPlaneEndpoint: "k8scp:6443"  #<-- Remember to use the node alias not the IP
networking:
  podSubnet: 192.168.0.0/16         #<-- Match the IP range from the Calico config file
EOF

# init and save output for future review
kubeadm init --config=kubeadm-config.yaml --upload-certs | tee kubeadm-init.out
```

---

## Change back normal user to connect to the cluster and create CNI

```sh
logout

mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config

sudo cp /root/kubeadm-init.out ~/kubeadm-init.out

# copy back the calico yaml and apply
sudo cp /root/calico.yaml .
kubectl apply -f calico.yaml
```

---

## Print K8s init Config

```sh
sudo kubeadm config print init-defaults
```

---

## Get The token from Control Plane (Master) (with in 2hrs)

```sh
# In ssh of MASTER node
sudo kubeadm token list

# Prepare the worker node join command
sudo kubeadm token create --print-join-command
```

## Setup bash auto-completion and alias

```sh
sudo apt-get install bash-completion -y

cat >> ~/.bashrc << EOF

# k8s alias and bash auto-completion
alias k="kubectl"
source <(kubectl completion bash)
complete -o default -F __start_kubectl k
EOF
```
