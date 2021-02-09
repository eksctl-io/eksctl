#!/bin/bash

DIR="${BASH_SOURCE%/*}"

# shellcheck source=tag-common.sh
source "${DIR}/tag-common.sh"

release_branch=$(release_branch)
release_version=$(release_generate release)
release_notes_file=$(ensure_release_notes "${release_version}")

check_origin

git checkout "${default_branch}"
check_current_branch "${default_branch}"
ensure_up_to_date "${default_branch}"

git checkout "${release_branch}"
check_current_branch "${release_branch}"
ensure_up_to_date "${release_branch}"

# Update eksctl version by removing the pre-release id

m="Release ${release_version}"

commit "${m}" "${release_notes_file}"


tag_version_and_latest "${m}" "v${release_version}"

# Update the site by putting everything from the release into the docs branch
git push --force origin "${release_branch}":docs

git checkout "${default_branch}"
git pull --ff-only origin "${default_branch}"

prepare_for_next_version_if_at "${release_version}"

git push origin "${release_branch}:${release_branch}"
git push origin "v${release_version}"
git push origin "${default_branch}:${default_branch}" || gh pr create --fill --label "skip-release-notes"
