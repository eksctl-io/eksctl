#!/bin/bash

# Set base URL for VPC CNI releases on GitHub
base_url="https://raw.githubusercontent.com/aws/amazon-vpc-cni-k8s/"

get_latest_release_tag() {
  curl -sL https://api.github.com/repos/aws/amazon-vpc-cni-k8s/releases/latest | jq -r '.tag_name'
}

latest_release_tag=$(get_latest_release_tag)

default_addons_dir="pkg/addons/default"

# Download the latest aws-k8s-cni.yaml file
curl -sL "$base_url$latest_release_tag/config/master/aws-k8s-cni.yaml?raw=1" --output "$default_addons_dir/assets/aws-node.yaml"

echo "found latest release tag:" $latest_release_tag

# Update the unit test file
sed -i "s/expectedVersion = \"\(.*\)\"/expectedVersion = \"$latest_release_tag\"/g" "$default_addons_dir/aws_node_test.go"
