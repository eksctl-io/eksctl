#!/bin/bash -ex

if [ -z "${GITHUB_REF_NAME}" ] || [ "${GITHUB_REF_TYPE}" != "tag" ] ; then
  echo "Expected a tag push event, skipping release workflow"
  exit 1
fi

tag="${GITHUB_REF_NAME}"
RELEASE_NOTES_FILE="docs/release_notes/${tag/-rc.*}.md"

if [ ! -f "${RELEASE_NOTES_FILE}" ]; then
    echo "Release notes file ${RELEASE_NOTES_FILE} does not exist. Exiting..."
    exit 1
fi

export RELEASE_DESCRIPTION="${tag}"

GORELEASER_CURRENT_TAG=v${tag} goreleaser release --rm-dist --timeout 60m --skip-validate --config=./.goreleaser.yml --release-notes="${RELEASE_NOTES_FILE}"
