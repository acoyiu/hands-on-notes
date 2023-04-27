# CKS Related

Topics:

- [CKS Related](#cks-related)
  - [kube-bench](#kube-bench)
    - [Manifest File for static pod](#manifest-file-for-static-pod)
  - [service-account](#service-account)
    - [Not auto mounting token inside /var/run/secrets/kubernetes.io/serviceaccount](#not-auto-mounting-token-inside-varrunsecretskubernetesioserviceaccount)
  - [Network-Policy](#network-policy)
  - [RBAC](#rbac)

<br/>

---

<br/>

## kube-bench

Running Kube-Bench on [CIS Benchmark](https://downloads.cisecurity.org/#/)

```sh
# Create a job to run kube-bench in k8s
k apply -f https://raw.githubusercontent.com/aquasecurity/kube-bench/main/job.yaml
```

### Manifest File for static pod

Path: /etc/kubernetes/manifests/

```sh
# Manifests yamls are in
ls -al /etc/kubernetes/manifests/
```

```yaml
# API Server - 所以 API server 其實係一個 static pod
apiVersion: v1
kind: Pod
metadata:
  annotations:
    kubeadm.kubernetes.io/kube-apiserver.advertise-address.endpoint: 192.168.49.2:8443
  labels:
    component: kube-apiserver
    tier: control-plane
  name: kube-apiserver
  namespace: kube-system
spec:
  containers:
    - command:
        - kube-apiserver
        - --advertise-address=192.168.49.2
        - --allow-privileged=true
        - --authorization-mode=Node,RBAC
        - --client-ca-file=/var/lib/minikube/certs/ca.crt
        - --enable-admission-plugins=NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,DefaultTolerationSeconds,NodeRestriction,MutatingAdmissionWebhook,ValidatingAdmissionWebhook,ResourceQuota
        - --enable-bootstrap-token-auth=true
        - --etcd-cafile=/var/lib/minikube/certs/etcd/ca.crt
        - --etcd-certfile=/var/lib/minikube/certs/apiserver-etcd-client.crt
        - --etcd-keyfile=/var/lib/minikube/certs/apiserver-etcd-client.key
        - --etcd-servers=https://127.0.0.1:2379
        - --kubelet-client-certificate=/var/lib/minikube/certs/apiserver-kubelet-client.crt
        - --kubelet-client-key=/var/lib/minikube/certs/apiserver-kubelet-client.key
        - --kubelet-preferred-address-types=InternalIP,ExternalIP,Hostname
        - --proxy-client-cert-file=/var/lib/minikube/certs/front-proxy-client.crt
        - --proxy-client-key-file=/var/lib/minikube/certs/front-proxy-client.key
        - --requestheader-allowed-names=front-proxy-client
        - --requestheader-client-ca-file=/var/lib/minikube/certs/front-proxy-ca.crt
        - --requestheader-extra-headers-prefix=X-Remote-Extra-
        - --requestheader-group-headers=X-Remote-Group
        - --requestheader-username-headers=X-Remote-User
        - --secure-port=8443
        - --service-account-issuer=https://kubernetes.default.svc.cluster.local
        - --service-account-key-file=/var/lib/minikube/certs/sa.pub
        - --service-account-signing-key-file=/var/lib/minikube/certs/sa.key
        - --service-cluster-ip-range=10.96.0.0/12
        - --tls-cert-file=/var/lib/minikube/certs/apiserver.crt
        - --tls-private-key-file=/var/lib/minikube/certs/apiserver.key
      image: k8s.gcr.io/kube-apiserver:v1.23.3
      imagePullPolicy: IfNotPresent
      livenessProbe:
        failureThreshold: 8
        httpGet:
          host: 192.168.49.2
          path: /livez
          port: 8443
          scheme: HTTPS
        initialDelaySeconds: 10
        periodSeconds: 10
        timeoutSeconds: 15
      name: kube-apiserver
      readinessProbe:
        failureThreshold: 3
        httpGet:
          host: 192.168.49.2
          path: /readyz
          port: 8443
          scheme: HTTPS
        periodSeconds: 1
        timeoutSeconds: 15
      resources:
        requests:
          cpu: 250m
      startupProbe:
        failureThreshold: 24
        httpGet:
          host: 192.168.49.2
          path: /livez
          port: 8443
          scheme: HTTPS
        initialDelaySeconds: 10
        periodSeconds: 10
        timeoutSeconds: 15
      volumeMounts:
        - mountPath: /etc/ssl/certs
          name: ca-certs
          readOnly: true
        - mountPath: /etc/ca-certificates
          name: etc-ca-certificates
          readOnly: true
        - mountPath: /var/lib/minikube/certs
          name: k8s-certs
          readOnly: true
        - mountPath: /usr/local/share/ca-certificates
          name: usr-local-share-ca-certificates
          readOnly: true
        - mountPath: /usr/share/ca-certificates
          name: usr-share-ca-certificates
          readOnly: true
  hostNetwork: true
  priorityClassName: system-node-critical
  securityContext:
    seccompProfile:
      type: RuntimeDefault
  volumes:
    - hostPath:
        path: /etc/ssl/certs
        type: DirectoryOrCreate
      name: ca-certs
    - hostPath:
        path: /etc/ca-certificates
        type: DirectoryOrCreate
      name: etc-ca-certificates
    - hostPath:
        path: /var/lib/minikube/certs
        type: DirectoryOrCreate
      name: k8s-certs
    - hostPath:
        path: /usr/local/share/ca-certificates
        type: DirectoryOrCreate
      name: usr-local-share-ca-certificates
    - hostPath:
        path: /usr/share/ca-certificates
        type: DirectoryOrCreate
      name: usr-share-ca-certificates
```

<br/>

---

<br/>

## service-account

### Not auto mounting token inside /var/run/secrets/kubernetes.io/serviceaccount

如果設定了，service account 創建時，並不會生成相關的 Secret

需要手動生成：

```yaml
# Create service account
apiVersion: v1
kind: ServiceAccount
metadata:
  name: build-robot
automountServiceAccountToken: false

# Create and bind secret to such service account
apiVersion: v1
kind: Secret
type: kubernetes.io/service-account-token
metadata:
  name: build-robot-token
  annotations:
    kubernetes.io/service-account.name: build-robot
```

<br/>

---

<br/>

## Network-Policy

Refer to [Kubectl](./kubectl.md#networkpolicies-netpol)

<br/>

---

<br/>

## RBAC

Refer to [Kubectl](./kubectl.md#rbac---role--rolebinding)
