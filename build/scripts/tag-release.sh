#!/bin/bash

DIR="${BASH_SOURCE%/*}"

# shellcheck source=tag-common.sh
source "${DIR}/tag-common.sh"

release_branch=$(release_branch)

check_prereqs
check_origin

git checkout "${default_branch}"
check_current_branch "${default_branch}"
ensure_up_to_date "${default_branch}"

git checkout "${release_branch}"
check_current_branch "${release_branch}"
ensure_up_to_date "${release_branch}"

# Update eksctl version by removing the pre-release id
release_version=$(release_generate release)
release_notes_file=$(ensure_release_notes "${release_version}")

m="Release ${release_version}"

commit "${m}" "${release_notes_file}"

tag_version_and_latest "${m}" "${release_version}"

make_pr "${release_branch}"

# Make PR to update default branch if necessary
git checkout "${default_branch}"
bump_version_if_not_at "${release_version}"
