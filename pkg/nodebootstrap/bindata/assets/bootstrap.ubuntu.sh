#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

source /var/lib/cloud/scripts/eksctl/bootstrap.helper.sh

echo "eksctl: running /etc/eks/bootstrap"
/etc/eks/bootstrap.sh "${CLUSTER_NAME}" \
  --dns-cluster-ip "${CLUSTER_DNS}" \
  --kubelet-extra-args "${KUBELET_EXTRA_ARGS}"

echo "eksctl: merging user options into kubelet-config.json"
trap 'rm -f ${TMP_KUBE_CONF}' EXIT
jq -s '.[0] * .[1]' "${KUBELET_CONFIG}" "${KUBELET_EXTRA_CONFIG}" > "${TMP_KUBE_CONF}"
mv "${TMP_KUBE_CONF}" "${KUBELET_CONFIG}"

echo "eksctl: restarting kubelet-eks"
snap restart kubelet-eks
echo "eksctl: done"
