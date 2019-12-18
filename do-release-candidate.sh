#!/bin/sh -ex

if [ -z "${CIRCLE_PULL_REQUEST}" ] && [ -n "${CIRCLE_TAG}" ] && [ "${CIRCLE_PROJECT_USERNAME}" = "weaveworks" ] ; then
  export RELEASE_DESCRIPTION="${CIRCLE_TAG}"
  RELEASE_NOTES_FILE="docs/release_notes/${CIRCLE_TAG/-rc.*}.md"

  if [[ ! -f "${RELEASE_NOTES_FILE}" ]]; then
    echo "Release notes ${RELEASE_NOTES_FILE} not found. Exiting..."
    exit 1
  fi

  goreleaser release --skip-validate --config=./.goreleaser.yml --release-notes="${RELEASE_NOTES_FILE}"

else
  echo "Not a tag release, skip publish"
fi
