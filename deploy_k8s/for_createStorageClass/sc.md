# Using K8s Sig subdir to create StorageClass

[Git repo](https://github.com/kubernetes-sigs/nfs-subdir-external-provisioner)

```sh
# Add helm repo
helm repo add nfs-subdir-external-provisioner https://kubernetes-sigs.github.io/nfs-subdir-external-provisioner/

# Check theble mount option (/exported/path in below) of a IP's NFS Specification
showmount -e <ip-if-nfs>

# Connect NFS(Linux)
helm install \
  nfs-subdir-external-provisioner \
  nfs-subdir-external-provisioner/nfs-subdir-external-provisioner \
  --set nfs.server=<ip-if-nfs> \
  --set nfs.path=/exported/path
```

## Testing PVC of StorageClass

```yaml
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: test-claim
spec:
  storageClassName: nfs-client
  accessModes:
    - ReadWriteMany
  resources:
    requests:
      storage: 1Gi
```

## Testing pod of StorageClass

```yaml
kind: Pod
apiVersion: v1
metadata:
  name: test-pod
spec:
  containers:
    - name: test-pod
      image: nginx:alpine
      volumeMounts:
        - name: nfs-pvc
          mountPath: /usr/share/nginx/html
  restartPolicy: "Always"
  volumes:
    - name: nfs-pvc
      persistentVolumeClaim:
        claimName: test-claim
```
