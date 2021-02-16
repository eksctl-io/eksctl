#!/bin/bash

set -o errexit
set -o pipefail
set -o nounset

export default_branch="main"

function branch_exists() {
  git ls-remote --heads origin "${1}" | grep -q "${1}"
}

function current_branch() {
  git rev-parse --abbrev-ref @
}

function release_generate() {
  go run pkg/version/generate/release_generate.go "${1}"
}

function check_origin() {
  if [[ ! "$(git remote get-url origin)" =~ weaveworks/eksctl(\-private)?(\.git)?$ ]] ; then
    echo "Invalid origin: $(git remote get-url origin)"
    exit 3
  fi
}

function release_branch() {
  echo "release-$(release_generate print-major-minor-version)"
}

function check_current_branch() {
  if [ ! "$(current_branch)" = "$1" ] ; then
    echo "Must be on $1 branch"
    exit 5
  fi
}

function ensure_up_to_date() {
  git pull --ff-only origin "$1"
}

function ensure_release_notes() {
  local release_notes_file="docs/release_notes/$1.md"
  if [[ ! -f "${release_notes_file}" ]]; then
    >&2 echo "Must have release notes ${release_notes_file}"
    exit 6
  fi
  echo "$release_notes_file"
}

function commit() {
  echo "Committing version changes"
  local commit_msg=$1
  local release_notes_file=$2
  git add ./pkg/version/release.go
  git add "${release_notes_file}"
  git commit --message "${commit_msg}"
}

function tag_version_and_latest() {
  echo "Tagging new version and latest_release"
  local commit_msg=$1
  local tag=$2
  git tag --annotate --message "${commit_msg}" --force "latest_release"
  git tag --annotate --message "${commit_msg}" "${tag}"
}

function prepare_for_next_version_if_at() {
  local dev_version
  dev_version=$(release_generate print-version)
  if [ "${dev_version}" != "$1" ]; then
    return 0
  fi
  echo "Preparing for next development iteration"
  release_generate development
  git add ./pkg/version/release.go
  git commit --message "Prepare for next development iteration"
}
