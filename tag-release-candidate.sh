#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

if [ "$#" -ne 0 ] ; then
  echo "Usage: ${0}. Got at least one extra argument"
  exit 1
fi

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

#if [[ ! "$(git remote get-url origin)" =~ ^git@github.com:weaveworks/eksctl(\.git)?$ ]] ; then
#  echo "Invalid origin: $(git remote get-url origin)"
#  exit 3
#fi

function branch_exists() {
  git ls-remote --heads origin "${1}" | grep -q "${1}"
}

function current_branch() {
  git rev-parse --abbrev-ref @
}

if ! branch_exists "${release_branch}" ; then
  git checkout master
  if [ ! "$(current_branch)" = master ] ; then
    echo "Must be on master branch"
    exit 4
  fi
  # Ensure local master is up-to-date by pulling its latest version from origin
  # and fast-forwarding local master:
  git fetch origin master
  git merge --ff-only origin/master
  # Create the release branch:
  git push origin master:"${release_branch}"

  # Prepare next development iteration in master
  go run pkg/version/generate/release_generate.go development
  git add ./pkg/version/release.go
  git commit --message "Prepare for next development iteration"
  git push origin master:master
fi

# Ensure local release branch is up-to-date by pulling its latest version from
# origin and fast-forwarding the local branch:
git fetch origin "${release_branch}"
git checkout "${release_branch}"
if [ ! "$(current_branch)" = "${release_branch}" ] ; then
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

# Update eksctl version
go run pkg/version/generate/release_generate.go release-candidate
v=$(go run pkg/version/generate/release_generate.go print-version)
git add ./pkg/version/release.go
git add "${RELEASE_NOTES_FILE}"

m="Tag ${v} release candidate"

git commit --message "${m}"

git fetch --force --tags origin

git push origin "${release_branch}"

git tag --annotate --message "${m}" "${v}"

git push --force --tags origin
