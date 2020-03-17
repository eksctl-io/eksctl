#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

function get_max_pods() {
  while read instance_type pods; do
    if  [[ "${instance_type}" == "${1}" ]] && [[ "${pods}" =~ ^[0-9]+$ ]] ; then
      echo ${pods};
      return
    fi
  done < /etc/eksctl/max_pods.map
}

NODE_IP="$(curl --silent http://169.254.169.254/latest/meta-data/local-ipv4)"
INSTANCE_ID="$(curl --silent http://169.254.169.254/latest/meta-data/instance-id)"
INSTANCE_TYPE="$(curl --silent http://169.254.169.254/latest/meta-data/instance-type)"
AWS_SERVICES_DOMAIN="$(curl --silent http://169.254.169.254/latest/meta-data/services/domain)"

source /etc/eksctl/kubelet.env # this can override MAX_PODS

cat > /etc/eksctl/kubelet.local.env <<EOF
NODE_IP=${NODE_IP}
INSTANCE_ID=${INSTANCE_ID}
INSTANCE_TYPE=${INSTANCE_TYPE}
MAX_PODS=${MAX_PODS:-$(get_max_pods "${INSTANCE_TYPE}")}
EOF

snap alias kubelet-eks.kubelet kubelet
snap alias kubectl-eks.kubectl kubectl
snap stop kubelet-eks
systemctl reset-failed

(
  # TODO: these should be looked at every time kubelet starts up,
  # which is what we do in AL2 (which is based on plain systemd,
  # and meant to be portable to most systemd distros), but it's
  # not clear how to load these from kubelet snap without having
  # to customise the snap itself
  source /etc/eksctl/kubelet.local.env
  source /etc/eksctl/kubelet.env
  source /etc/eksctl/metadata.env

  flags=(
    "node-ip=${NODE_IP}"
    "max-pods=${MAX_PODS}"
    "node-labels=${NODE_LABELS},alpha.eksctl.io/instance-id=${INSTANCE_ID}"
    "pod-infra-container-image=${AWS_EKS_ECR_ACCOUNT}.dkr.ecr.${AWS_DEFAULT_REGION}.${AWS_SERVICES_DOMAIN}/eks/pause-amd64:3.1"
    "cloud-provider=aws"
    "cni-bin-dir=/opt/cni/bin"
    "cni-conf-dir=/etc/cni/net.d"
    "container-runtime=docker"
    "network-plugin=cni"
    "register-node=true"
    "register-with-taints=${NODE_TAINTS}"
    "kubeconfig=/etc/eksctl/kubeconfig.yaml"
    "config=/etc/eksctl/kubelet.yaml"
  )

  snap set kubelet-eks "${flags[@]}"
)

snap start kubelet-eks
