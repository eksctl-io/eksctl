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

go generate ./pkg/version

git add ./pkg/version/release.go

m="Tag ${v} release"

git commit --message "${m}"

git push git@github.com:weaveworks/eksctl master

git fetch --tags git@github.com:weaveworks/eksctl

git tag --annotate --message "${m}" "${v}"
git tag --annotate --message "${m}" --force "latest_release" "${v}"

git push --force --tags git@github.com:weaveworks/eksctl
