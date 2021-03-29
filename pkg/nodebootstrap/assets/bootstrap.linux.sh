#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

# Use IMDSv2 to get metadata
TOKEN="$(curl --silent -X PUT -H "X-aws-ec2-metadata-token-ttl-seconds: 600" http://169.254.169.254/latest/api/token)"
function get_metadata() {
  curl --silent -H "X-aws-ec2-metadata-token: $TOKEN" "http://169.254.169.254/latest/meta-data/$1"
}

source /etc/eksctl/kubelet.env # file written by bootstrapper

KUBELET_CONFIG='/etc/kubernetes/kubelet/kubelet-config.json'
KUBELET_EXTRA_ARGS='/etc/eksctl/kubelet-extra.json'
INSTANCE_ID="$(get_metadata instance-id)"
INSTANCE_LIFECYCLE="$(get_metadata instance-life-cycle)"
CLUSTER_DNS="${CLUSTER_DNS:-}"
NODE_TAINTS="${NODE_TAINTS:-}"
NODE_LABELS="${NODE_LABELS},node-lifecycle=${INSTANCE_LIFECYCLE},alpha.eksctl.io/instance-id=${INSTANCE_ID}"

echo "eksctl: running /etc/eks/bootstrap"
/etc/eks/bootstrap.sh "${CLUSTER_NAME}" \
  --dns-cluster-ip "${CLUSTER_DNS}" \
  --kubelet-extra-args "--register-with-taints=${NODE_TAINTS} --node-labels=${NODE_LABELS}"

echo "eksctl: merging user options into kubelet-config.json"
TMP_CONF='/tmp/kubelet-conf.json'
trap 'rm -f ${TMP_CONF}' EXIT
jq -s '.[0] * .[1]' "${KUBELET_CONFIG}" "${KUBELET_EXTRA_ARGS}" > "${TMP_CONF}"
mv "${TMP_CONF}" "${KUBELET_CONFIG}"
