# K8size

## Deploy to k8s

Assume already having the imagePullSecret set in ServiceAccount

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: ping-pusher
spec:
  restartPolicy: Never
  containers:
    - name: ping-pusher
      image: registry.greatics.net/ping-pusher
      args:
        - pushgateway-service.pushgateway:9091
        - instance_name
        - job_name
        - http://123.57.136.251
        - https://google.com
```
