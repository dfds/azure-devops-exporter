#!/bin/bash
CONTINER_IMAGE=${1}

# kustomize edit set image a=b/c:2
kustomize build ../k8s-orchestration/overlays/for-development | \
tee /dev/tty | \
kubectl apply -f -
