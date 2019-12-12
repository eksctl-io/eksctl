#!/bin/sh -eu

REQUIREMENTS_FILE=.requirements

if [ ! -f "$REQUIREMENTS_FILE" ]
then
	echo "Requirements file $REQUIREMENTS_FILE not found. Exiting..."
	exit 1
fi

# Versions of build tools. One in each line "<module> <version>"
while IFS= read -r req; do
  go list  -json "${req}"  | jq '.Module| "\(.Path) \(.Version)"'
done < ${REQUIREMENTS_FILE}

git ls-tree --full-tree @ -- Dockerfile
git ls-tree --full-tree @ -- install-build-deps.sh
