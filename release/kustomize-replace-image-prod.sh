#!/bin/bash
set -eu -o pipefail
# Replaces the image in the for-production overlay
BUILD_ID=$1

(cd ./k8s-orchestration/overlays/for-production; kustomize edit set image azure-devops-exporter/image-placeholder=579478677147.dkr.ecr.eu-central-1.amazonaws.com/ded/azure-devops-exporter:${BUILD_ID}
)

sed -i -e 's/resources/bases/g' ./k8s-orchestration/overlays/for-production/kustomization.yaml