#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

source /var/lib/cloud/scripts/eksctl/bootstrap.helper.sh

echo "eksctl: running /etc/eks/bootstrap"
/etc/eks/bootstrap.sh "${CLUSTER_NAME}" \
  --dns-cluster-ip "${CLUSTER_DNS}" \
  --kubelet-extra-args "--register-with-taints=${NODE_TAINTS} --node-labels=${NODE_LABELS}"

echo "eksctl: merging user options into kubelet-config.json"
trap 'rm -f ${TMP_KUBE_CONF}' EXIT
jq -s '.[0] * .[1]' "${KUBELET_CONFIG}" "${KUBELET_EXTRA_ARGS}" > "${TMP_KUBE_CONF}"
mv "${TMP_KUBE_CONF}" "${KUBELET_CONFIG}"

echo "eksctl: merging user options into docker daemon.json"
trap 'rm -f ${TMP_DOCKER_CONF}' EXIT
jq -s '.[0] * .[1]' "${DOCKER_CONFIG}" "${DOCKER_EXTRA_CONFIG}" > "${TMP_DOCKER_CONF}"
mv "${TMP_DOCKER_CONF}" "${DOCKER_CONFIG}"

systemctl daemon-reload
echo "eksctl: restarting docker daemon"
systemctl restart docker
echo "eksctl: restarting kubelet-eks"
systemctl restart kubelet
echo "eksctl: done"
