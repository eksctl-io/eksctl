#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

function get_max_pods() {
  MAX_PODS_FILE="/etc/eksctl/max_pods_map.txt"
  while read instance_type pods; do

    if  [[ "${instance_type}" == "${1}" ]] && [[ "${pods}" =~ ^[0-9]+$ ]] ; then
      echo ${pods};
      return
    fi ;

  done < "${MAX_PODS_FILE}"
}

NODE_IP=$(curl --silent http://169.254.169.254/latest/meta-data/local-ipv4)
INSTANCE_ID=$(curl --silent http://169.254.169.254/latest/meta-data/instance-id)
INSTANCE_TYPE=$(curl --silent http://169.254.169.254/latest/meta-data/instance-type)

source /etc/eksctl/kubelet.env

if [[ -z "${MAX_PODS+x}" ]];
  then export MAX_PODS=$(get_max_pods ${INSTANCE_TYPE});
fi

echo "NODE_IP=${NODE_IP}" > /etc/eksctl/kubelet.local.env
echo "INSTANCE_ID=${INSTANCE_ID}" >> /etc/eksctl/kubelet.local.env
echo "INSTANCE_TYPE=${INSTANCE_TYPE}" >> /etc/eksctl/kubelet.local.env
echo "MAX_PODS=${MAX_PODS}" >> /etc/eksctl/kubelet.local.env

systemctl daemon-reload
systemctl enable kubelet
systemctl start kubelet
