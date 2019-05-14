#!/bin/bash -e

function get_max_pods() {
  MAX_PODS_FILE="/etc/eksctl/max_pods_map.txt"
  grep "${INSTANCE_TYPE}" "${MAX_PODS_FILE}" | while read instance_type pods; do

    if  [[ "${pods}" =~ ^[0-9]+$ ]] ; then
      echo ${pods};
      return
    fi ;

    done
}

NODE_IP=$(curl --silent http://169.254.169.254/latest/meta-data/local-ipv4)
INSTANCE_ID=$(curl --silent http://169.254.169.254/latest/meta-data/instance-id)
INSTANCE_TYPE=$(curl --silent http://169.254.169.254/latest/meta-data/instance-type)

source /etc/eksctl/kubelet.env

MAX_PODS=${MAX_PODS:-$(get_max_pods)}

set -o nounset

echo "NODE_IP=${NODE_IP}" > /etc/eksctl/kubelet.local.env
echo "INSTANCE_ID=${INSTANCE_ID}" >> /etc/eksctl/kubelet.local.env
echo "INSTANCE_TYPE=${INSTANCE_TYPE}" >> /etc/eksctl/kubelet.local.env
echo "MAX_PODS=${MAX_PODS}" >> /etc/eksctl/kubelet.local.env

systemctl daemon-reload
systemctl enable kubelet
systemctl start kubelet
