#!/bin/bash -ex

if [ -z "${GITHUB_REF_NAME}" ] || [ "${GITHUB_REF_TYPE}" != "tag" ] ; then
    echo "Expected a tag push event, skipping release workflow"
    exit 1
fi

tag="${GITHUB_REF_NAME}"

export RELEASE_DESCRIPTION="${tag} (permalink)"
RELEASE_NOTES_FILE="docs/release_notes/${tag}.md"

if [ ! -f "${RELEASE_NOTES_FILE}" ]; then
    echo "Release notes ${RELEASE_NOTES_FILE} not found. Exiting..."
    exit 1
fi

cat ./.goreleaser.yml ./.goreleaser.brew.yml > .goreleaser.brew.combined.yml
GORELEASER_CURRENT_TAG=v${tag} goreleaser release --rm-dist --timeout 60m --skip-validate --config=./.goreleaser.brew.combined.yml --release-notes="${RELEASE_NOTES_FILE}"
