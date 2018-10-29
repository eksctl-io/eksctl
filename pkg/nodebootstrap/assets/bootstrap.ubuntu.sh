#!/bin/bash -eu

echo "NODE_IP=$(hostname -i)" > /etc/eksctl/kubelet.local.env

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
    "address=0.0.0.0"
    "node-ip=${NODE_IP}"
    "cluster-dns=${CLUSTER_DNS}"
    "max-pods=${MAX_PODS}"
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
    "kubeconfig=/etc/eksctl/kubeconfig.yaml"
    "feature-gates=RotateKubeletServerCertificate=true"
    "anonymous-auth=false"
    "client-ca-file=/etc/eksctl/ca.crt"
  )

  snap set kubelet-eks "${flags[@]}"
)

snap start kubelet-eks
