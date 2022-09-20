#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

source /var/lib/cloud/scripts/eksctl/bootstrap.helper.sh

echo "eksctl: running /etc/eks/bootstrap"
/etc/eks/bootstrap.sh "${CLUSTER_NAME}" \
  --apiserver-endpoint "${API_SERVER_URL}" \
  --b64-cluster-ca "${B64_CLUSTER_CA}" \
  --dns-cluster-ip "${CLUSTER_DNS}" \
  --kubelet-extra-args "${KUBELET_EXTRA_ARGS}" \
  --container-runtime "${CONTAINER_RUNTIME}" \
  ${ENABLE_LOCAL_OUTPOST:+--enable-local-outpost "${ENABLE_LOCAL_OUTPOST}"} \
  ${CLUSTER_ID:+--cluster-id "${CLUSTER_ID}"}


echo "eksctl: merging user options into kubelet-config.json"
trap 'rm -f ${TMP_KUBE_CONF}' EXIT
jq -s '.[0] * .[1]' "${KUBELET_CONFIG}" "${KUBELET_EXTRA_CONFIG}" > "${TMP_KUBE_CONF}"
mv "${TMP_KUBE_CONF}" "${KUBELET_CONFIG}"

systemctl daemon-reload
echo "eksctl: restarting kubelet-eks"
systemctl restart kubelet
echo "eksctl: done"
