#!/bin/bash -ex

DIR="${BASH_SOURCE%/*}"

# shellcheck source=tag-common.sh
. "${DIR}/tag-common.sh"


check_origin

release_branch=$(release_branch)

git switch docs
git merge --no-edit "origin/${release_branch}"
git push origin docs
