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
      git push origin "$(git branch --show-current)"
  fi
}

release_branch=$(release_branch)
candidate_for_version=$(release_generate print-version)
release_notes_file=$(ensure_release_notes "${candidate_for_version}")

check_prereqs
check_origin

git checkout "${default_branch}"
check_current_branch "${default_branch}"
ensure_up_to_date "${default_branch}"

git checkout -

create_branch_from_if_doesnt_exist "${release_branch}" "$(git branch --show-current)"

git checkout "${release_branch}"
check_current_branch "${release_branch}"
ensure_up_to_date "${release_branch}" || echo "${release_branch} not found in origin, will push new branch upstream."

# Update eksctl version to release-candidate
rc_version=$(release_generate release-candidate)
m="Tag ${rc_version} release candidate"

commit "${m}" "${release_notes_file}"

tag_version_and_latest "${m}" "${rc_version}"

# Make PR to release branch
make_pr "${release_branch}"

# Make PR to update default branch if necessary
git checkout "${default_branch}"
bump_version_if_not_at "${candidate_for_version}"
