#!/bin/sh -ex


GIT_REPO_URL="https://github.com/eksctl-io/eksctl"
RELEASE_URL_FORMAT="https://github.com/eksctl-io/eksctl/releases/download/v%s/eksctl_Linux_amd64.tar.gz"

download_previous_release() {
    previous_tag=$(git ls-remote --tags $GIT_REPO_URL | grep -E -v "(refs/tags/(latest_release|v))|\-rc|\{\}" | cut -d/ -f3 | sort -Vr \
        | tail -n "+${GO_BACK_VERSIONS}" | head -1)

    download_url=$(printf "$RELEASE_URL_FORMAT" "$previous_tag")
    wget -q -O - "${download_url}" | tar -xz -C "${DOWNLOAD_DIR}"
}

download_previous_release
