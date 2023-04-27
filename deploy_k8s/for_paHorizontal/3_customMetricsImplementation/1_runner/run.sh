# ==============================
# Make sure namespace is correct

NS=default

if [[ -n $1 ]]; then
  NS=$1
fi

echo $NS

# ===========================
# delete configmap if existed

kubectl -n $NS delete cm/appfiles || echo "OK: configmap not existed, delete action failed"

# ================
# create configmap

tar --exclude="node_modules" -cvf app.tar ./app

kubectl -n $NS create configmap appfiles --from-file=app.tar

rm -f app.tar

# ================
# create

kubectl -n $NS apply -f runner.yaml

kubectl -n $NS rollout restart deployment/httpserver