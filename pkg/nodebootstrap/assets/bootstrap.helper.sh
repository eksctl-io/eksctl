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

INSTANCE_ID="$(get_metadata instance-id)"
INSTANCE_LIFECYCLE="$(get_metadata instance-life-cycle)"
CLUSTER_DNS="${CLUSTER_DNS:-}"
NODE_TAINTS="${NODE_TAINTS:-}"
NODE_LABELS="${NODE_LABELS},node-lifecycle=${INSTANCE_LIFECYCLE},alpha.eksctl.io/instance-id=${INSTANCE_ID}"
CLUSTER_NAME="${CLUSTER_NAME}"
KUBELET_CONFIG='/etc/kubernetes/kubelet/kubelet-config.json'
KUBELET_EXTRA_ARGS='/etc/eksctl/kubelet-extra.json'
DOCKER_CONFIG='/etc/docker/daemon.json'
DOCKER_EXTRA_CONFIG='/etc/eksctl/docker-extra.json'
TMP_KUBE_CONF='/tmp/kubelet-conf.json'
TMP_DOCKER_CONF='/tmp/docker-conf.json'
