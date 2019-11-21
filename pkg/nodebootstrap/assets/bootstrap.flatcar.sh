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

function install_kubernetes() {
  echo "Installing Kubernetes"

  cd
  mkdir -p /opt/cni/bin
  CNI_VERSION="v0.6.0"
  CNI_PLUGIN_VERSION="v0.7.5"
  ARCH="amd64"

  wget https://github.com/containernetworking/cni/releases/download/${CNI_VERSION}/cni-${ARCH}-${CNI_VERSION}.tgz
  wget https://github.com/containernetworking/cni/releases/download/${CNI_VERSION}/cni-${ARCH}-${CNI_VERSION}.tgz.sha512
  sha512sum -c cni-${ARCH}-${CNI_VERSION}.tgz.sha512
  tar -xvf cni-${ARCH}-${CNI_VERSION}.tgz -C /opt/cni/bin
  rm cni-${ARCH}-${CNI_VERSION}.tgz cni-${ARCH}-${CNI_VERSION}.tgz.sha512

  wget https://github.com/containernetworking/plugins/releases/download/${CNI_PLUGIN_VERSION}/cni-plugins-${ARCH}-${CNI_PLUGIN_VERSION}.tgz
  wget https://github.com/containernetworking/plugins/releases/download/${CNI_PLUGIN_VERSION}/cni-plugins-${ARCH}-${CNI_PLUGIN_VERSION}.tgz.sha512
  sha512sum -c cni-plugins-${ARCH}-${CNI_PLUGIN_VERSION}.tgz.sha512
  tar -xvf cni-plugins-${ARCH}-${CNI_PLUGIN_VERSION}.tgz -C /opt/cni/bin
  rm cni-plugins-${ARCH}-${CNI_PLUGIN_VERSION}.tgz cni-plugins-${ARCH}-${CNI_PLUGIN_VERSION}.tgz.sha512

  mkdir -p /opt/bin
  wget -O /opt/bin/aws-iam-authenticator https://github.com/kubernetes-sigs/aws-iam-authenticator/releases/download/v0.4.0/aws-iam-authenticator_0.4.0_linux_amd64
  chmod +x /opt/bin/aws-iam-authenticator
}

NODE_IP="$(curl --silent http://169.254.169.254/latest/meta-data/local-ipv4)"
INSTANCE_ID="$(curl --silent http://169.254.169.254/latest/meta-data/instance-id)"
INSTANCE_TYPE="$(curl --silent http://169.254.169.254/latest/meta-data/instance-type)"

source /etc/eksctl/kubelet.env # this can override MAX_PODS

cat > /etc/eksctl/kubelet.local.env <<EOF
NODE_IP=${NODE_IP}
INSTANCE_ID=${INSTANCE_ID}
INSTANCE_TYPE=${INSTANCE_TYPE}
MAX_PODS=${MAX_PODS:-$(get_max_pods "${INSTANCE_TYPE}")}
EOF

# Install Kubernetes
install_kubernetes

systemctl daemon-reload
systemctl enable kubelet
systemctl start kubelet
