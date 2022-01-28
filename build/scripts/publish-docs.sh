#!/bin/bash -ex

DIR="${BASH_SOURCE%/*}"

if [ -z "${NETLIFY_BUILD_HOOK_URL}" ] ; then
  echo "NETLIFY_BUILD_HOOK_URL is required"
  exit 1
fi

# shellcheck source=tag-common.sh
. "${DIR}/tag-common.sh"

release_branch=$(release_branch)

curl -X POST -d "trigger_branch=${release_branch}" -d "trigger_title=Triggered+by+Release+action" "${NETLIFY_BUILD_HOOK_URL}"
