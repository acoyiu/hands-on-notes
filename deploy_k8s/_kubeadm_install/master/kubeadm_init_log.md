# Your Kubernetes control-plane has initialized successfully!

To start using your cluster, you need to run the following as a regular user:

```sh
mkdir -p $HOME/.kube
sudo cp -i /etc/kubernetes/admin.conf $HOME/.kube/config
sudo chown $(id -u):$(id -g) $HOME/.kube/config
```

Alternatively, if you are the root user, you can run:

```sh
export KUBECONFIG=/etc/kubernetes/admin.conf
```

You should now deploy a pod network to the cluster.
Run "kubectl apply -f [podnetwork].yaml" with one of the options listed at: https://kubernetes.io/docs/concepts/cluster-administration/addons/

- 添加主節點 control-plane node -
  You can now join any number of the control-plane node running the following command on each as root:

```sh
  kubeadm join k8scp:6443 --token o00pkq.jzyehath2lvqbaux \
        --discovery-token-ca-cert-hash sha256:f8adbdb365c72ec151c972a3a89ab002f00758f9e868d64a425588e2fb21e74f \
        --control-plane --certificate-key 9154285f5865122859b5e8c053ea73083d8e9577595dd0cacc0704d736bc94bd
```

Please note that the certificate-key gives access to cluster sensitive data, keep it secret!
As a safeguard, uploaded-certs will be deleted in two hours; If necessary, you can use

```sh
"kubeadm init phase upload-certs --upload-certs" to reload certs afterward.
```

- 添加工作節點 -
  Then you can join any number of worker nodes by running the following on each as root:

```sh
kubeadm join k8scp:6443 --token o00pkq.jzyehath2lvqbaux \
        --discovery-token-ca-cert-hash sha256:f8adbdb365c72ec151c972a3a89ab002f00758f9e868d64a425588e2fb21e74f
```
