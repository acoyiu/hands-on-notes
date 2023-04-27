# Create service
kubectl apply -f ./php.yaml

# Create AutoScaler
kubectl apply -f ./hpa.yaml

# Create Loading
kubectl run -i --tty load-generator --rm --image=busybox:1.28 --restart=Never -- /bin/sh -c "while sleep 0.01; do wget -q -O- http://php-apache; done"

# Watch scaling
kubectl get hpa php-apache --watch