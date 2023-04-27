# 在集群内部的 DNS 服务里面添加上 hostname 的解析

```sh
$ kubectl edit configmap coredns -n kube-system
```

```conf
apiVersion: v1
data:
  Corefile: |
    .:53 {
        errors
        health
        hosts {  # 添加集群節點 hosts 影射信息
          10.151.30.11 ydzs-master
          10.151.30.57 ydzs-node3
          10.151.30.59 ydzs-node4
          10.151.30.22 ydzs-node1
          10.151.30.23 ydzs-node2
          fallthrough
        }
        kubernetes cluster.local in-addr.arpa ip6.arpa {
           pods insecure
           upstream
           fallthrough in-addr.arpa ip6.arpa
        }
        prometheus :9153
        proxy . /etc/resolv.conf
        cache 30
        reload
    }
kind: ConfigMap
metadata:
  creationTimestamp: 2019-05-18T11:07:46Z
  name: coredns
  namespace: kube-system
```