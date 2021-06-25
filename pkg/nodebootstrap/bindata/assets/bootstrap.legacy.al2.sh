#!/bin/bash
# TODO Deprecated, will be removed soon

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

# Use IMDSv2 to get metadata
TOKEN="$(curl --silent -X PUT -H "X-aws-ec2-metadata-token-ttl-seconds: 600" http://169.254.169.254/latest/api/token)"
function get_metadata() {
  curl --silent -H "X-aws-ec2-metadata-token: $TOKEN" "http://169.254.169.254/latest/meta-data/$1"
}

NODE_IP="$(get_metadata local-ipv4)"
INSTANCE_ID="$(get_metadata instance-id)"
INSTANCE_TYPE="$(get_metadata instance-type)"
AWS_SERVICES_DOMAIN="$(get_metadata services/domain)"


source /etc/eksctl/kubelet.env # this can override MAX_PODS

INSTANCE_LIFECYCLE="$(get_metadata instance-life-cycle)"
NODE_LABELS="${NODE_LABELS},node-lifecycle=${INSTANCE_LIFECYCLE}"


cat > /etc/eksctl/kubelet.local.env <<EOF
NODE_IP=${NODE_IP}
INSTANCE_ID=${INSTANCE_ID}
INSTANCE_TYPE=${INSTANCE_TYPE}
AWS_SERVICES_DOMAIN=${AWS_SERVICES_DOMAIN}
MAX_PODS=${MAX_PODS:-$(get_max_pods "${INSTANCE_TYPE}")}
NODE_LABELS=${NODE_LABELS}
EOF

systemctl daemon-reload
systemctl enable kubelet
systemctl start kubelet
