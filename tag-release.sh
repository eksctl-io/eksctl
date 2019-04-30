#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

if [ "$#" -ne 1 ] ; then
  echo "Usage: ${0} <tag>"
  exit 1
fi

if [ ! "$(git rev-parse --abbrev-ref @)" = master ] ; then
  echo "Must be on master branch"
  exit 2
fi

v="${1}"

RELEASE_NOTES_FILE="docs/release_notes/${v}.md"

if [[ ! -f "${RELEASE_NOTES_FILE}" ]]; then
  echo "Release notes ${RELEASE_NOTES_FILE} not found. Exiting..."
  return
fi

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
