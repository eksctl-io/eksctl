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


echo "NODE_IP=$(curl --silent http://169.254.169.254/latest/meta-data/local-ipv4)" > /etc/eksctl/kubelet.local.env
echo "INSTANCE_ID=$(curl --silent http://169.254.169.254/latest/meta-data/instance-id)" >> /etc/eksctl/kubelet.local.env
echo "INSTANCE_TYPE=$(curl -s http://169.254.169.254/latest/meta-data/instance-type)" >> /etc/eksctl/kubelet.local.env

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

  MAX_PODS=${MAX_PODS:-$(get_max_pods)}
  set -o nounset

  flags=(
    "address=0.0.0.0"
    "node-ip=${NODE_IP}"
    "cluster-dns=${CLUSTER_DNS}"
    "max-pods=${MAX_PODS}"
    "node-labels=${NODE_LABELS},alpha.eksctl.io/instance-id=${INSTANCE_ID}"
    "authentication-token-webhook=true"
    "authorization-mode=Webhook"
    "allow-privileged=true"
    "pod-infra-container-image=602401143452.dkr.ecr.${AWS_DEFAULT_REGION}.amazonaws.com/eks/pause-amd64:3.1"
    "cloud-provider=aws"
    "cluster-domain=cluster.local"
    "cni-bin-dir=/opt/cni/bin"
    "cni-conf-dir=/etc/cni/net.d"
    "container-runtime=docker"
    "network-plugin=cni"
    "cgroup-driver=cgroupfs"
    "register-node=true"
    "register-with-taints=${NODE_TAINTS}"
    "kubeconfig=/etc/eksctl/kubeconfig.yaml"
    "feature-gates=RotateKubeletServerCertificate=true"
    "anonymous-auth=false"
    "client-ca-file=/etc/eksctl/ca.crt"
  )

  snap set kubelet-eks "${flags[@]}"
)

snap start kubelet-eks
