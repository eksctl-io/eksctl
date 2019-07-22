#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

if [ "$#" -ne 1 ] ; then
  echo "Usage: ${0} <tag>"
  exit 1
fi

v="${1}"
candidate_for="${v/-rc.*}"
release_branch="release-${candidate_for%.*}"  # e.g.: 0.2.0-rc.0 -> release-0.2

if [ "${v}" = "${candidate_for}" ] ; then
  echo "Must provide release candidate tag, use './tag-release.sh ${v}' instead"
  exit 2
fi

if ! [[ "${release_branch}" =~ ^release-[0-9]+\.[0-9]+$ ]] ; then
  echo "Invalid release branch: ${release_branch}"
  exit 3
fi

function branch_exists() {
  git ls-remote --heads git@github.com:weaveworks/eksctl.git "${1}" | grep "${1}" >/dev/null
}

if ! branch_exists "${release_branch}" ; then
  git checkout master
  if [ ! "$(git rev-parse --abbrev-ref @)" = master ] ; then
    echo "Must be on master branch"
    exit 4
  fi
  # Ensure local master is up-to-date by pulling its latest version from origin
  # and fast-forwarding local master:
  git fetch origin master
  git merge --ff-only origin/master
  # Create the release branch:
  git push origin master:"${release_branch}"
fi

# Ensure local release branch is up-to-date by pulling its latest version from
# origin and fast-forwarding the local branch:
git fetch origin "${release_branch}"
git checkout "${release_branch}"
if [ ! "$(git rev-parse --abbrev-ref @)" = "${release_branch}" ] ; then
  echo "Must be on ${release_branch} branch"
  exit 5
fi
git merge --ff-only origin/"${release_branch}"

RELEASE_NOTES_FILE="docs/release_notes/${candidate_for}.md"

if [[ ! -f "${RELEASE_NOTES_FILE}" ]]; then
  echo "Must have release notes ${RELEASE_NOTES_FILE}"
  exit 6
fi

export RELEASE_GIT_TAG="${v}"

go generate ./pkg/version

git add ./pkg/version/release.go
git add ${RELEASE_NOTES_FILE}

m="Tag ${v} release candidate"

git commit --message "${m}"

git fetch --force --tags git@github.com:weaveworks/eksctl

git push git@github.com:weaveworks/eksctl master

git tag --annotate --message "${m}" "${v}"

git push --force --tags git@github.com:weaveworks/eksctl
