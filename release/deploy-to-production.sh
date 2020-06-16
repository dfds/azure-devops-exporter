#!/bin/bash
set -eu -o pipefail

BUILD_ID=$1

(cd ../k8s-orchestration/overlays/for-production; kustomize edit set image azure-devops-exporter/image-placeholder=579478677147.dkr.ecr.eu-central-1.amazonaws.com/ded/kafka-janitor:${BUILD_ID}
)

kustomize build ../k8s-orchestration/overlays/for-production | \
tee /dev/tty | \
kubectl apply -f -
