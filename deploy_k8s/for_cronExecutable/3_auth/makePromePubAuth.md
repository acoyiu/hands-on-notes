# Make Prometheus Public with Auth

```sh
# -c=創建，會刪除原有紀錄
htpasswd -c ./auth <username>

# 新增 user
htpasswd -B ./auth <username>

# kubectl create secret generic <name-of-secret> --from-file=<path-to-file>
kubectl create secret generic ing-prometheus-auth --from-file=./auth
# -n if needed

kubectl apply -f ./ing.yaml
```
