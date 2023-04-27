# Name of SA
SA_NAME=build-robot

# Name of secret
SECRET_NAME=$SA_NAME-secret

# your server name goes here
server=https://192.168.0.57:16443

# 
# ---
# 

# Create Service Acc
kubectl apply -f - <<EOF
apiVersion: v1
kind: ServiceAccount
metadata:
  name: $SA_NAME
EOF

# Create of Secret
kubectl apply -f - <<EOF
apiVersion: v1
kind: Secret
type: kubernetes.io/service-account-token
metadata:
  name: $SECRET_NAME
  annotations:
    kubernetes.io/service-account.name: $SA_NAME
EOF

# 
# ---
# 

ca=$(kubectl get secret/$SECRET_NAME -o jsonpath='{.data.ca\.crt}')
token=$(kubectl get secret/$SECRET_NAME -o jsonpath='{.data.token}' | base64 --decode)
namespace=$(kubectl get secret/$SECRET_NAME -o jsonpath='{.data.namespace}' | base64 --decode)

echo "
apiVersion: v1
kind: Config
clusters:
- name: sa-cluster
  cluster:
    certificate-authority-data: ${ca}
    server: ${server}
contexts:
- name: sa-context
  context:
    cluster: sa-cluster
    user: sa-user
    namespace: $namespace
current-context: sa-context
users:
- name: sa-user
  user:
    token: ${token}
" > sa.kubeconfig