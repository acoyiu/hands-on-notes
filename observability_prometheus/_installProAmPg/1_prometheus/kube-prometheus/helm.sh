#!/bin/bash

kubectl create ns prometheus

helm -n prometheus upgrade -i prometheus bitnami/kube-prometheus -f ./value.yaml

# # Ingress if needed
# kubectl apply -f ./ing.yaml