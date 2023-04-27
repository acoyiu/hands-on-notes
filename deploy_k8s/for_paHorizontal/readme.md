# Metrics & Autoscale

## 基于自定义指标的 HPA

https://imroc.cc/k8s/best-practice/custom-metrics-hpa/

Kubernetes 默認提供 CPU 和內存作為 HPA 彈性伸縮的指標，如果有更複雜的場景需求，比如基於業務單副本 QPS 大小來進行自動擴縮容，可以考慮自行安裝 prometheus-adapter 來實現基於自定義指標的 Pod 彈性伸縮。

Kubernetes 提供了 Custom Metrics API 與 External Metrics API 來對 HPA 的指標進行擴展，讓用戶能夠根據實際需求進行自定義。

prometheus-adapter 對這兩種 API 都有支持，通常使用 Custom Metrics API 就夠了，本文也主要針對此 API 來實現使用自定義指標進行彈性伸縮。

- 需要一个部署并配置了 Metrics Server 的集群 (可以用 kubectl top)
- Kubernetes Metrics Server 从集群中的 kubelets 收集资源指标，
- 並通過Metrics API，將它們暴露在 Kubernetes apiserver 中，從而公开这些指标
- 使用 APIService 添加代表指标读数的新资源

APIService 就是暴露 API-Server 的 Endpoint，透過它可以取得 Kubelet 的資訊

## 指標服務器提供

- 適用於大多數集群的單一部署
- 快速自動縮放，每 15 秒收集一次指標。
- 資源效率，為集群中的每個節點使用 1 mili CPU 內核和 2 MB 內存。
- 可擴展支持多達 5,000 個節點的集群。

## 您可以將 Metrics Server 用於：

- （HorizontalAutoScaler）基於 CPU/內存的水平自動縮放
- （VerticalAutoScaler）自動調整/建議容器所需的資源

==================================================

## Prometheus metrics need to associate with K8s resource (the resource to export metrics) to start working

- Group = api-versions (not include version)
- Resource = Name in api-resources

## Prometheus-Adaptor 不是直接存取 Prometheus 的 metrics，而是從 Prometheus 發現有什麽 metrics，然後自己去 metrics-exporter instance 取出 metrics

https://github.com/kubernetes-sigs/prometheus-adapter/blob/master/docs/config.md

https://github.com/kubernetes-sigs/prometheus-adapter/blob/master/docs/config-walkthrough.md