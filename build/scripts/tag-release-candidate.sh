#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

default_branch="master"

function branch_exists() {
  git ls-remote --heads origin "${1}" | grep -q "${1}"
}

function current_branch() {
  git rev-parse --abbrev-ref @
}

function release_generate() {
  go run pkg/version/generate/release_generate.go ${1}
}

if [[ ! "$(git remote get-url origin)" =~ weaveworks/eksctl(\-private)?(\.git)?$ ]] ; then
  echo "Invalid origin: $(git remote get-url origin)"
  exit 3
fi

candidate_for=$(release_generate print-version)
release_branch="release-$(release_generate print-major-minor-version)"

if ! [[ "${release_branch}" =~ ^release-[0-9]+\.[0-9]+$ ]] ; then
  echo "Invalid release branch: ${release_branch}"
  exit 3
fi

if [ ! "$(current_branch)" = "${release_branch}" ] ; then
    echo "Creating ${release_branch} from ${default_branch}"
    if [ ! "$(current_branch)" = "${default_branch}" ] ; then
      echo "Must be on ${default_branch} branch"
      exit 7
    fi
    git checkout -b "${release_branch}"
fi

if [ ! "$(current_branch)" = "${release_branch}" ] ; then
  echo "Must be on ${release_branch} branch"
  exit 5
fi

if ! git show-ref --verify --quiet "refs/heads/${default_branch}"; then
  echo "The ${default_branch} branch must exist"
  exit 7
fi

# Ensure local release branch is up-to-date by pulling its latest version from
# origin and fast-forwarding the local branch:
git pull --ff-only origin "${release_branch}" || echo "${release_branch} not found in origin. Will push new branch upstream"

RELEASE_NOTES_FILE="docs/release_notes/${candidate_for}.md"

if [[ ! -f "${RELEASE_NOTES_FILE}" ]]; then
  echo "Must have release notes ${RELEASE_NOTES_FILE}"
  exit 6
fi


# Update eksctl version
full_version=$(release_generate release-candidate)
export RELEASE_GIT_TAG="${full_version}"
git add ./pkg/version/release.go
git add "${RELEASE_NOTES_FILE}"

m="Tag ${full_version} release candidate"

# Push branch
git commit --message "${m}"
git push origin "${release_branch}"

# Push tags
git tag --annotate --message "${m}" --force "latest_release"
git tag --annotate --message "${m}" "${full_version}"
git push origin "${full_version}"

# Check if we need to bump version in the default branch
git checkout "${default_branch}"
if [ ! "$(current_branch)" = "${default_branch}" ] ; then
  echo "Must be on ${default_branch} branch"
  exit 7
fi
git pull --ff-only origin "${default_branch}"

dev_version=$(release_generate print-version)

# Increase next development iteration if needed
if [ "${dev_version}" == "${candidate_for}" ]; then
  echo "Preparing for next development iteration"
  release_generate development
  git add ./pkg/version/release.go
  git commit --message "Prepare for next development iteration"
  git push origin "${default_branch}":"${default_branch}" || gh pr create --fill --label "skip-release-notes"
fi
