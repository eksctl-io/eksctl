#!/bin/sh -ex


GIT_REPO_URL="git@github.com:weaveworks/eksctl"
RELEASE_URL_FORMAT="https://github.com/weaveworks/eksctl/releases/download/%s/eksctl_Linux_amd64.tar.gz"

download_previous_release() {
    previous_tag=$(git ls-remote --tags $GIT_REPO_URL | grep -E -v "latest_release|\-rc|\{\}" | cut -d/ -f3 | sort -Vr \
        | tail -n "+${GO_BACK_VERSIONS}" | head -1)

    download_url=$(printf "$RELEASE_URL_FORMAT" "$previous_tag")
    wget -q -O - "${download_url}" | tar -xz -C "${DOWNLOAD_DIR}"
}

download_previous_release
