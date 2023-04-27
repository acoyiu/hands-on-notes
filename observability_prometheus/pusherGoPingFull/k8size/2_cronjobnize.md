# Cronjobnize

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: ping-pusher
spec:
  schedule: "*/2 * * * *"
  successfulJobsHistoryLimit: 1
  failedJobsHistoryLimit: 1
  concurrencyPolicy: Replace
  jobTemplate:
    spec:
      parallelism: 1
      completions: 1
      backoffLimit: 3
      activeDeadlineSeconds: 20
      ttlSecondsAfterFinished: 600
      template:
        metadata:
          name: ping-pusher
        spec:
          restartPolicy: Never
          containers:
            - name: ping-pusher
              image: registry.greatics.net/ping-pusher
              args:
                - pushgateway-service.pushgateway:9091
                - dev_ping_pusher
                - dev_ping_to_production
                - https://api.funacademycn.com
                - https://highlights.milkcargo.cn
                - http://123.57.136.251
                - https://google.com
```
