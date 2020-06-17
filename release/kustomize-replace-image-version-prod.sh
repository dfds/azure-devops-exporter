#!/bin/bash
set -eu -o pipefail
# Replaces the image in the for-production overlay
BUILD_ID=$1

sed -i -e 's/tag_placeholder/'${BUILD_ID}'/g' ./k8s-orchestration/overlays/for-production/deployment-patch.yaml