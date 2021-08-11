#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

source /etc/eksctl/kubelet.env # file written by bootstrapper

# Use IMDSv2 to get metadata
TOKEN="$(curl --silent -X PUT -H "X-aws-ec2-metadata-token-ttl-seconds: 600" http://169.254.169.254/latest/api/token)"
function get_metadata() {
  curl --silent -H "X-aws-ec2-metadata-token: $TOKEN" "http://169.254.169.254/latest/meta-data/$1"
}

API_SERVER_URL="${API_SERVER_URL}"
B64_CLUSTER_CA="${B64_CLUSTER_CA}"
INSTANCE_ID="$(get_metadata instance-id)"
INSTANCE_LIFECYCLE="$(get_metadata instance-life-cycle)"
CLUSTER_DNS="${CLUSTER_DNS:-}"
NODE_TAINTS="${NODE_TAINTS:-}"
MAX_PODS="${MAX_PODS:-}"
NODE_LABELS="${NODE_LABELS},node-lifecycle=${INSTANCE_LIFECYCLE},alpha.eksctl.io/instance-id=${INSTANCE_ID}"

KUBELET_ARGS=("--node-labels=${NODE_LABELS}")
[[ -n "${NODE_TAINTS}" ]] && KUBELET_ARGS+=("--register-with-taints=${NODE_TAINTS}")
# --max-pods as a CLI argument is deprecated, this is a workaround until we deprecate support for maxPodsPerNode
[[ -n "${MAX_PODS}" ]] && KUBELET_ARGS+=("--max-pods=${MAX_PODS}")
KUBELET_EXTRA_ARGS="${KUBELET_ARGS[@]}"

CLUSTER_NAME="${CLUSTER_NAME}"
KUBELET_CONFIG='/etc/kubernetes/kubelet/kubelet-config.json'
KUBELET_EXTRA_CONFIG='/etc/eksctl/kubelet-extra.json'
TMP_KUBE_CONF='/tmp/kubelet-conf.json'
CONTAINER_RUNTIME="${CONTAINER_RUNTIME:-dockerd}" # default for al2 just in case, not used in ubuntu
