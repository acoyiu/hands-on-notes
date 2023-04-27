#!/bin/bash

# helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
# helm repo update

kubectl create ns prometheus

helm -n prometheus upgrade -i prometheus prometheus-community/prometheus -f ./value.yaml

# # Ingress if needed
# kubectl apply -f ./ing.yaml