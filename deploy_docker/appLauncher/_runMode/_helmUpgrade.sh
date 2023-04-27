#!/bin/bash

ProjectDirectory=$1
InstallmentName=$2
k8sSecret=$3
Registry=$4
imageTag=$5
imagePullPolicy=$6

pullPolicy="Always"
if [[ -n "$imagePullPolicy" ]]; then
  pullPolicy=$imagePullPolicy
fi

cd "${ProjectDirectory}/_helm"

helm upgrade \
    $InstallmentName \
    . \
    --install \
    --set spec.imagePullSecrets="${k8sSecret}" \
    --set spec.fromRegistry="${Registry}" \
    --set spec.tag="${imageTag}" \
    --set spec.imagePullPolicy="${pullPolicy}"
    # --dry-run

cd ../../
