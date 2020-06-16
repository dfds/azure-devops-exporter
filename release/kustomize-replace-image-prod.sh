#!/bin/bash
set -eu -o pipefail
# Replaces the image in the for-production overlay
BUILD_ID=$1

(cd ./k8s-orchestration/overlays/for-production; kustomize edit set image azure-devops-exporter/image-placeholder=579478677147.dkr.ecr.eu-central-1.amazonaws.com/ded/kafka-janitor:${BUILD_ID}
)