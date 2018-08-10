#!/bin/bash

# Obtain cluster credentials and update `kubeconfig.yaml` with apiserver URL

set -o errexit
set -o pipefail
set -o nounset

source /etc/eksctl/metadata.env
source /etc/eksctl/kubelet.local.env

/opt/aws/bin/cfn-signal \
  --exit-code "$1" \
  --stack "${AWS_EKS_CLUSTER_NAME}" \
  --resource NodeGroup \
  --id "${NDOE_ID}"
  --region "${AWS_REGION}"
