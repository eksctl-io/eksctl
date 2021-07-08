#!/bin/sh -eu

REQUIREMENTS_FILE=.requirements

if [ ! -f "$REQUIREMENTS_FILE" ]
then
	echo "Requirements file $REQUIREMENTS_FILE not found. Exiting..."
	exit 1
fi

# Versions of build tools. One in each line "<module> <version>"
while IFS= read -r req; do
    # remove the @ as go list is not module aware unless using -m but then it
    # tries the local module in which we don't import it.
    req=${req%@*}
    go list -json "${req}"  | jq '.Module| "\(.Path) \(.Version)"'
done < ${REQUIREMENTS_FILE}

git ls-tree --full-tree @ -- build/docker/Dockerfile
git ls-tree --full-tree @ -- .requirements
git ls-tree --full-tree @ -- build/scripts/install-build-deps.sh
