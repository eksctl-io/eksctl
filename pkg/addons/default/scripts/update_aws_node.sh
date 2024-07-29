#!/bin/bash

# Set base URL for VPC CNI releases on GitHub
base_url="https://raw.githubusercontent.com/aws/amazon-vpc-cni-k8s/"

get_latest_release_tag() {
  curl -sL https://api.github.com/repos/aws/amazon-vpc-cni-k8s/releases/latest | jq -r '.tag_name'
}

latest_release_tag=$(get_latest_release_tag)

# Check if the latest release tag was found
if [ -z "$latest_release_tag" ]; then
  echo "Could not find the latest release tag."
  exit 1
fi

# If running in GitHub Actions, export the release tag for use in the workflow
if [ "$GITHUB_ACTIONS" = "true" ]; then
  echo "LATEST_RELEASE_TAG= to $latest_release_tag" >> $GITHUB_ENV
else
  echo "Found the latest release tag: $latest_release_tag"
fi

default_addons_dir="pkg/addons/default"

# Download the latest aws-k8s-cni.yaml file
curl -sL "$base_url$latest_release_tag/config/master/aws-k8s-cni.yaml?raw=1" --output "$default_addons_dir/assets/aws-node.yaml"

# Check if the download was successful
if [ $? -eq 0 ]; then
  echo "Downloaded the latest AWS Node manifest to $default_addons_dir/assets/aws-node.yaml"
else
  echo "Failed to download the latest AWS Node manifest."
  exit 1
fi

# Update the unit test file
sed -i "s/expectedVersion = \"\(.*\)\"/expectedVersion = \"$latest_release_tag\"/g" "$default_addons_dir/aws_node_test.go"
