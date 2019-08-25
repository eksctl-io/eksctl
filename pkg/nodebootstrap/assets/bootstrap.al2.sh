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

EC2_AVAIL_ZONE=$(curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone)
AWS_DEFAULT_REGION="`echo \"$EC2_AVAIL_ZONE\" | sed 's/[a-z]$//'`"
NODE_IP="$(curl --silent http://169.254.169.254/latest/meta-data/local-ipv4)"
INSTANCE_ID="$(curl --silent http://169.254.169.254/latest/meta-data/instance-id)"
INSTANCE_TYPE="$(curl --silent http://169.254.169.254/latest/meta-data/instance-type)"

source /etc/eksctl/kubelet.env # this can override MAX_PODS

#FIXME: Ideally this script should not call & depends on the AWS APIs
#however it is the only way at the moment to obtain the instanceLifecycle value (Spot/on-demand)

INSTANCE_LIFECYCLE=$(aws ec2 describe-instances --region $AWS_DEFAULT_REGION --instance-ids ${INSTANCE_ID} \
--query 'Reservations[0].Instances[0].InstanceLifecycle' --output text)

if [ "$INSTANCE_LIFECYCLE" == "spot" ] && [ "$SPOT_NODE_LABELS" != ""  ]; then
  if [ "$NODE_LABELS" == "" ];then
     NODE_LABELS="${SPOT_NODE_LABELS}"
  else
     NODE_LABELS="${NODE_LABELS},${SPOT_NODE_LABELS}"
  fi
fi

if [ "$INSTANCE_LIFECYCLE" == "spot" ] && [ "$SPOT_NODE_TAINTS" != "" ]; then
  if [ "$NODE_TAINTS" == "" ];then
     NODE_TAINTS="${SPOT_NODE_TAINTS}"
  else
     NODE_TAINTS="${NODE_TAINTS},${SPOT_NODE_TAINTS}"
  fi
fi

cat > /etc/eksctl/kubelet.local.env <<EOF
NODE_IP=${NODE_IP}
INSTANCE_ID=${INSTANCE_ID}
INSTANCE_TYPE=${INSTANCE_TYPE}
MAX_PODS=${MAX_PODS:-$(get_max_pods "${INSTANCE_TYPE}")}
NODE_LABELS=${NODE_LABELS}
NODE_TAINTS=${NODE_TAINTS}
EOF

systemctl daemon-reload
systemctl enable kubelet
systemctl start kubelet
