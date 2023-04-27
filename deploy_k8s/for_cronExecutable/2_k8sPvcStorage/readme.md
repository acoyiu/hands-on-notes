# Usage of project

```sh
k create ns generic-app-storage

chmod -R 744 ./executables
k -n generic-app-storage cp ./executables generic-app-storage/generic-app-pvc-reader:/storage
k -n generic-app-storage exec -it po/generic-app-pvc-reader -- sh -c "ls -al /storage/executables"
```

```sh
# cmd <pushgateway-location> <instance_name> <job_name> [<url> <url> ...]
./ping_metrics_to_gateway http://localhost:9091 instance_name job_name http://123.57.136.251 https://google.com
./ping_metrics_to_gateway http://pushgateway-service.pushgateway:9091 dev_pushgateway ping_production_link http://123.57.136.251 https://google.com

./ping_metrics_to_gateway http://localhost:9091 dev_pushgateway ping_production_link https://api.funacademycn.com https://highlights.milkcargo.cn
./ping_metrics_to_gateway http://pushgateway-service.pushgateway:9091 dev_pushgateway ping_production_link https://api.funacademycn.com https://highlights.milkcargo.cn
```

```sh
./api-test-fa \
  -gateway http://localhost:9091/metrics/job/api_test_fa/instance/dev_pushgateway/usage/monitoring \
  -address https://api.funacademycn.com/v1

./api-test-fa \
  -gateway http://pushgateway-service.pushgateway:9091/metrics/job/api_test_fa/instance/dev_pushgateway/usage/monitoring \
  -address http://service-of-grpc-gateway.fa:8080/v1
```

```sh
./api-test-mc2 \
  -gateway http://pushgateway-service.pushgateway:9091/metrics/job/api_test_mc2/instance/dev_pushgateway/usage/monitoring \
  -address http://service-of-grpc-gateway.app-mc2:8080/v1
```

## Common port-forward command

```sh
k -n linkerd-viz port-forward svc/prometheus 9090:9090 --address=0.0.0.0 &
k -n pushgateway port-forward svc/pushgateway-service 9091:9091 --address=0.0.0.0 &
```