# Grafana Alert

| Grafana has 3 Element in Alerting | Description                                                   |
| :-------------------------------- | :------------------------------------------------------------ |
| Alert Rule                        | 用於設定 Alert 的規則                                         |
| Contact Point                     | 用於設定發送的對象方式，如：Slack                             |
| Notification Policy               | 用於設定 Root policy，選擇 Contact Point，Group By，Timing 等 |

## Alert Rule

![a](./_img/grafana%20alert%20rule.jpg)

## Contact Point

![b](_img/grafana%20contact%20poing%20(send%20out%20method).jpg)

## Notification Policy

![c](./_img/grafana%20alerting%20policy.jpg)

## Alert Template Data

![d](./_img/grafana%20label%20getter.jpg)

https://grafana.com/docs/grafana/next/alerting/contact-points/message-templating/template-data/

```txt
$value = [
    var='B0' 
    metric='{
        instance="go-exporter.custom-metric-export:80",
        job="custom-exporter",
        namespace="hw",
        pod_name="hiw-api-develop-deployment-bc6db44bd-p525j",
        under_deployment="hiw-api-develop-deployment-bc6db44bd"
    }' 
    labels={
        instance=go-exporter.custom-metric-export:80, 
        job=custom-exporter, 
        namespace=hw, 
        pod_name=hiw-api-develop-deployment-bc6db44bd-p525j, 
        under_deployment=hiw-api-develop-deployment-bc6db44bd
    } 
    value=0.002127655047542452 
]
```
