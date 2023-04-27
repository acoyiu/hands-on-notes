## Custom Ingress-Nginx configuration
ingress-nginx-controller deployment will load all the configmap in the "ingress-nginx" namespace

## aliyun tutor
https://help.aliyun.com/document_detail/86533.html


# create configmap
```sh
kubectl apply -f ./custom-headers.yaml
```

# list all configmap under "ingress-nginx namespace"
```sh
kubectl -n ingress-nginx get cm
```

# update deployment, as update configmap won't restart and update
```sh
kubectl -n ingress-nginx get deploy
kubectl -n ingress-nginx get pods
kubectl -n ingress-nginx rollout restart deploy/ingress-nginx-controller
```

# get all current pods and delete old pod for update
```sh
kubectl -n ingress-nginx delete pod/ingress-nginx-controller-[podHase]
kubectl -n ingress-nginx get pods
```

# log out nginx config for confirmation
```sh
kubectl exec ingress-nginx-controller-[podHash] -n ingress-nginx -- cat /etc/nginx/nginx.conf
```