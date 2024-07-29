#!/bin/bash

get_latest_release_tag() {
  curl -sL https://api.github.com/repos/NVIDIA/k8s-device-plugin/releases/latest | jq -r '.tag_name'
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

assets_addons_dir="pkg/addons/assets"

curl -sL "https://raw.githubusercontent.com/NVIDIA/k8s-device-plugin/$latest_release_tag/deployments/static/nvidia-device-plugin.yml" -o "$assets_addons_dir/nvidia-device-plugin.yaml"


# Check if the download was successful
if [ $? -eq 0 ]; then
  echo "Downloaded the latest NVIDIA device plugin manifest to $assets_addons_dir/nvidia-device-plugin.yaml"
else
  echo "Failed to download the NVIDIA device plugin manifest."
  exit 1
fi
