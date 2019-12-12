#!/bin/sh -eu

# Versions of build tools. One in each line "<module> <version>"
while IFS= read -r req; do
  go list  -json "${req}"  | jq '.Module| "\(.Path) \(.Version)"'
done < .requirements

git ls-tree --full-tree @ -- Dockerfile
git ls-tree --full-tree @ -- install-build-deps.sh
