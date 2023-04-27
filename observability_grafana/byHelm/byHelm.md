# Helm install Grafana

## [Grafana CLI](https://grafana.com/docs/grafana/latest/administration/cli/)

## By official Helm chart

```sh
# create namespaces and add repo
helm repo add grafana https://grafana.github.io/helm-charts
helm search repo grafana

kubectl create ns grafana
```

# create value.yaml to overwrite default's to use persistent storage
```yaml
persistence:
  type: pvc
  enabled: true
  storageClassName: nfs-client
  size: 5Gi

## Use Replica or AutoScaling
# replicas: 2
autoscaling:
  enabled: true
  minReplicas: 2
  maxReplicas: 3
```

# install chart
helm -n grafana upgrade -i grafana grafana/grafana -f ./value.yaml --dry-run

# Get your 'admin' user password by running:
kubectl get secret --namespace grafana grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo

# Expose Service
kubectl -n grafana port-forward svc/grafana 8080:80 --address=0.0.0.0
```

## Or By binami/grafana (Not recommended)

```sh
kubectl create ns grafana
helm repo add bitnami https://charts.bitnami.com/bitnami
helm -n grafana upgrade -i grafana bitnami/grafana
```

### Please Read the Logs: Prometheus can be accessed via port "9090" on the following DNS name from within your cluster

```sh
echo "User: admin"
echo "Password: $(kubectl get secret grafana-admin --namespace grafana -o jsonpath="{.data.GF_SECURITY_ADMIN_PASSWORD}" | base64 -d)"
# or
echo "Password: $(kubectl get secret/grafana --namespace grafana -o jsonpath="{.data.admin-password}" | base64 -d)"
 
# port-forward grafana
kubectl -n grafana port-forward svc/grafana 3000:3000 --address='0.0.0.0' &
```

---

## The Data URL in Grafana

<http://prometheus-kube-prometheus-prometheus.prometheus:9090>

---

## Sample PromQL for K8s

### **Network**

```sh
# tcp in
irate( node_netstat_Tcp_InSegs[2m] ) / 2

# tcp out
irate( node_netstat_Tcp_OutSegs[1m] )
```

### **CPU**

[Reference](https://blog.csdn.net/shm19990131/article/details/107162470)

```sh
# K8s Average CPU Using by node
(
    1-(
        sum(
            increase(
                node_cpu_seconds_total{mode='idle'}[1m]
            )
        ) by (instance)
    )
    / 
    (
        sum(
            increase(
                node_cpu_seconds_total[1m]
            )
        ) by (instance)
    )
) * 100
```

### **Memory**

[Reference](https://cloud.tencent.com/developer/article/1644608)

| 名称                               | 类型    | 单位         | 说明                                                                                                                         |
| :--------------------------------- | :------ | :----------- | :--------------------------------------------------------------------------------------------------------------------------- |
| container_memory_rss               | gauge   | 字节数 bytes | 常驻内存集（Resident Set Size），分配给进程使用实际物理内存，而不是磁盘上缓存的虚拟内存                                      |
| container_memory_usage_bytes       | gauge   | 字节数 bytes | 当前使用的内存量，包括所有使用的内存，不管有没有被访问                                                                       |
| container_memory_max_usage_bytes   | gauge   | 字节数 bytes | 最大内存使用量的记录                                                                                                         |
| container_memory_cache             | gauge   | 字节数 bytes | 高速缓存（cache）的使用量                                                                                                    |
| container_memory_swap              | gauge   | 字节数 bytes | 虚拟内存（swap）是用磁盘模拟内存使用,当物理内存快用完达一定比例,把部分不用的内存数据交换到硬盘保存，需要使用时再调入物理内存 |
| container_memory_working_set_bytes | gauge   | 字节数 bytes | 当前内存工作集（working set）使用量。                                                                                        |
| container_memory_failcnt           | counter | 次           | 申请内存失败次数计数                                                                                                         |
| container_memory_failures_total    | counter | 次           | 累计的内存申请错误次数                                                                                                       |

### **Pod**

```sh
# Pod Restart count
rate(kube_pod_container_status_restarts_total{namespace=~'fa|mc2'}[15m])

# Pod total number
count(count by (pod) (container_memory_rss))

# Memory used by each pod
sum (container_memory_rss{pod!=''}) by (pod)
```

![stat](./_images/stat.jpg)

### **Storage**

```sh
# K8s nodes available storage
( node_filesystem_avail_bytes{mountpoint='/var'} / node_filesystem_size_bytes{mountpoint='/var'} ) * 100
```

![sto](./_images/sto.jpg)

### **Display Name**

${__field.labels.pod_name}