#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

if [ "$#" -ne 1 ] ; then
  echo "Usage: ${0} <tag>"
  exit 1
fi

v="${1}"

export RELEASE_GIT_TAG="${v}"

go generate ./cmd/eksctl

git add ./cmd/eksctl/version_release.go

m="Tag ${v} release"

git commit --message "${m}"

git push git@github.com:weaveworks/eksctl master

git fetch --tags git@github.com:weaveworks/eksctl

git tag --annotate --message "${m}" "${v}"
git tag --delete latest_release

git push --tags git@github.com:weaveworks/eksctl
