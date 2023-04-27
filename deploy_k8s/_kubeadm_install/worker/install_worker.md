# Install of kubeadm (kubelet)

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
kversion=1.24.6-00
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
cpip= ????
echo "$cpip k8scp" >> /etc/hosts
```

---

## Join to the CP node

**_By the command generated on Master(CP)_**

---

## Check cluster info in CP (Master) node

```sh
# RUN in MASTER(CP) node
kubectl get node

kubectl describe node cp
# Work line by line to view the resources and their current status !!!
```

- cp won’t allow non-infrastructure pods by default for security and resource contention reasons
- Taints: node-role.kubernetes.io/master:NoSchedule
- To allow: kubectl taint nodes --all node-role.kubernetes.io/master-

---

## Check the CNI is working

```sh
kubectl get pods --all-namespaces

# pods are stuck in ContainerCreating status you may have to delete them, causing new ones to be generated. Delete both pods and check to see they show a Running state. Your pod names will be different.
kubectl -n kube-system delete \
  pod coredns-576cbf47c7-vq5dz ...
```
