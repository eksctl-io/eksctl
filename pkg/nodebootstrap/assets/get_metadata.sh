#!/bin/bash

# Obtain metadata about the instance, parse eksctl tag value to and write out `metadata.env` and `kubelet.local.env`
# for other scripts and the systemd unit to use

set -o errexit
set -o pipefail
set -o nounset

# we could use the following, but we still need jq to get region and instance ID:
# aws --output text --region us-west-2 ec2 describe-tags --filters Name=resource-type,Values=instance Name=resource-id,Values=i-017e37452ee14d8d7 --query "Tags[?Key=='eksctl.cluster.k8s.io/v1alpha1/cluster-name'].Value"

# TODO: avoid having to do all this – https://github.com/weaveworks/eksctl/issues/157

instanceInfo="$(curl --silent http://169.254.169.254/latest/dynamic/instance-identity/document)"

instanceID="$(echo "${instanceInfo}" | jq -r .instanceId)"
instapceIP="$(echo "${instanceInfo}" | jq -r .privateIp)"

AWS_DEFAULT_REGION="$(echo "${instanceInfo}" | jq -r .region)"

export AWS_DEFAULT_REGION

# metadata file is intended for other bash scripts as well as systemd units, so it's in simple format
echo "AWS_DEFAULT_REGION=${AWS_DEFAULT_REGION}"  > /etc/eksctl/metadata.env

echo "NODE_IP=${instapceIP}" > /etc/eksctl/kubelet.local.env

get_tags() {
  aws ec2 describe-tags \
    --output json \
    --filters "Name=resource-type,Values=instance" "Name=resource-id,Values=${instanceID}"
}

get_cluster_name() {
    jq -r '.Tags[] | select(.Key == "eksctl.cluster.k8s.io/v1alpha1/cluster-name") | .Value'
}

check_cluster_owned() {
    tag="kubernetes.io/cluster/${1}"
    jq --arg tag "${tag}" -r '.Tags[] | select(.Key == $tag) | .Value'
}

tags="$(get_tags)"

cluster_name="$(echo "${tags}" | get_cluster_name)"

test "$(echo "${tags}" | check_cluster_owned "${cluster_name}")" = "owned"

echo "AWS_EKS_CLUSTER_NAME=${cluster_name}"  >> /etc/eksctl/metadata.env