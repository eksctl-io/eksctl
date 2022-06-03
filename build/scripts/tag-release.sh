#!/bin/bash

DIR="${BASH_SOURCE%/*}"

# shellcheck source=tag-common.sh
source "${DIR}/tag-common.sh"

check_prereqs
check_origin

release_version=$(release_generate print-version)
release_notes_file=$(ensure_release_notes "${release_version}")

msg="Release ${release_version}"
tag_and_push_release "${release_version}" "${msg}"

# Make PR to update default branch if necessary
git checkout "${default_branch}"
bump_version_if_not_at "${release_version}"
