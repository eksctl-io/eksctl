#!/bin/bash

DIR="${BASH_SOURCE%/*}"

# shellcheck source=tag-common.sh
source "${DIR}/tag-common.sh"

function create_branch_from_if_doesnt_exist() {
  wanted_branch="$1"
  source_branch="$2"
  if ! git checkout "${wanted_branch}" >/dev/null; then
      git checkout "${source_branch}"
      echo "Creating ${wanted_branch} from ${source_branch}"
      git checkout -b "${wanted_branch}"
  fi
}

release_branch=$(release_branch)
candidate_for_version=$(release_generate print-version)
release_notes_file=$(ensure_release_notes "${candidate_for_version}")

# Check all conditions
check_origin

git checkout "${default_branch}"
check_current_branch "${default_branch}"
ensure_up_to_date "${default_branch}"

create_branch_from_if_doesnt_exist "${release_branch}" "${default_branch}"

git checkout "${release_branch}"
check_current_branch "${release_branch}"
ensure_up_to_date "${release_branch}" || echo "${release_branch} not found in origin, will push new branch upstream."

# Update eksctl version to release-candidate
rc_version=$(release_generate release-candidate)
m="Tag ${rc_version} release candidate"

commit "${m}" "${release_notes_file}"

tag_version_and_latest "${m}" "v${rc_version}"

# Check if we need to bump version in the default branch
git checkout "${default_branch}"
prepare_for_next_version_if_at "${candidate_for_version}"

git push origin "${release_branch}:${release_branch}"
git push origin "v${rc_version}"
git push origin "${default_branch}:${default_branch}" || gh pr create --fill --label "skip-release-notes"
