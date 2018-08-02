#!/bin/bash

# Obtain cluster credentials and update `kubeconfig.yaml` with apiserver URL

set -o errexit
set -o pipefail
set -o nounset

source /etc/eksctl/metadata.env

export AWS_DEFAULT_REGION
export KUBECONFIG="/etc/eksctl/kubeconfig.yaml"

describe_cluster() {
  aws eks describe-cluster \
    --output json \
    --name "${AWS_EKS_CLUSTER_NAME}"
}

clusterInfo="$(describe_cluster)"

endpoint="$(echo "${clusterInfo}" | jq -r .cluster.endpoint)"

echo "${clusterInfo}" | jq -r '.cluster.certificateAuthority.data' | base64 -d > /etc/eksctl/ca.crt

kubectl config set-cluster kubernetes --server "${endpoint}"

