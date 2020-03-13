#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

function get_max_pods() {
  while read instance_type pods; do
    if  [[ "${instance_type}" == "${1}" ]] && [[ "${pods}" =~ ^[0-9]+$ ]] ; then
      echo ${pods}
      return
    fi
  done < /etc/eksctl/max_pods.map
}

NODE_IP="$(curl --silent http://169.254.169.254/latest/meta-data/local-ipv4)"
INSTANCE_ID="$(curl --silent http://169.254.169.254/latest/meta-data/instance-id)"
INSTANCE_TYPE="$(curl --silent http://169.254.169.254/latest/meta-data/instance-type)"
MACHINE=$(uname -m)
AWS_SERVICES_DOMAIN="$(curl --silent http://169.254.169.254/latest/meta-data/services/domain)"

if [ "$MACHINE" == "x86_64" ]; then
    ARCH="amd64"
elif [ "$MACHINE" == "aarch64" ]; then
    ARCH="arm64"
else
    ARCH="amd64"
    echo "Set default architecture to '$ARCH'" >&2
fi

source /etc/eksctl/kubelet.env # this can override MAX_PODS

cat > /etc/eksctl/kubelet.local.env <<EOF
NODE_IP=${NODE_IP}
INSTANCE_ID=${INSTANCE_ID}
INSTANCE_TYPE=${INSTANCE_TYPE}
AWS_SERVICES_DOMAIN=${AWS_SERVICES_DOMAIN}
MAX_PODS=${MAX_PODS:-$(get_max_pods "${INSTANCE_TYPE}")}
ARCH=${ARCH}
EOF

systemctl daemon-reload
systemctl enable kubelet
systemctl start kubelet
