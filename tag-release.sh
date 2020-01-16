#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

function branch_exists() {
  git ls-remote --heads origin "${1}" | grep -q "${1}"
}

function current_branch() {
  git rev-parse --abbrev-ref @
}

#if [[ ! "$(git remote get-url origin)" =~ ^git@github.com:weaveworks/eksctl(\.git)?$ ]] ; then
#  echo "Invalid origin: $(git remote get-url origin)"
#  exit 3
#fi

v=$(go run pkg/version/generate/release_generate.go print-version)

release_branch="release-${v}"  # e.g.: 0.2.0 -> release-0.2.0
if ! [[ "${release_branch}" =~ ^release-[0-9]+\.[0-9]+\.[0-9]+$ ]] ; then
  echo "Invalid release branch: ${release_branch}"
  exit 3
fi

if [ ! "$(current_branch)" = "${release_branch}" ] ; then
  echo "Must be on ${release_branch} branch"
  exit 5
fi

# Ensure local release branch is up-to-date by pulling its latest version from
# origin and fast-forwarding the local branch:
git pull --ff-only origin "${release_branch}" || echo "${release_branch} not found in origin. Will push new branch upstream"

RELEASE_NOTES_FILE="docs/release_notes/${v}.md"
if [[ ! -f "${RELEASE_NOTES_FILE}" ]]; then
  echo "Must have release notes ${RELEASE_NOTES_FILE}"
  exit 6
fi

export RELEASE_GIT_TAG="${v}"

# Update eksctl version by removing the pre-release id
go run pkg/version/generate/release_generate.go release
git add ./pkg/version/release.go
git add "${RELEASE_NOTES_FILE}"

m="Release ${v}"

git commit --message "${m}"
git push origin "${release_branch}"

# Create the release tag and push it to start release process
git tag --annotate --message "${m}" --force "latest_release"
git tag --annotate --message "${m}" "${v}"
git push origin "${v}"

# Update the site by putting everything from master into the docs branch
git push --force origin "${release_branch}":docs

### TODO if master is not dev then next dev iteration
git checkout master
if [ ! "$(current_branch)" = master ] ; then
  echo "Must be on master branch"
  exit 7
fi
git pull --ff-only origin master

master_version=$(go run pkg/version/generate/release_generate.go print-version)

# Increase next development iteration if needed
if [ "${master_version}" == "${v}" ]; then
  go run pkg/version/generate/release_generate.go print-version
  git add ./pkg/version/release.go
  git commit --message "Prepare for next development iteration"
  git push origin master:master
fi

