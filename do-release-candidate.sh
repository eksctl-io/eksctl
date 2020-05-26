#!/bin/sh -ex

if [ -z "${CIRCLE_PULL_REQUEST}" ] && [ -n "${CIRCLE_TAG}" ] && [ "${CIRCLE_PROJECT_USERNAME}" = "weaveworks" ] ; then
  export RELEASE_DESCRIPTION="${CIRCLE_TAG}"
  RELEASE_NOTES_FILE="docs/release_notes/${CIRCLE_TAG/-rc.*}.md"

  if [[ ! -f "${RELEASE_NOTES_FILE}" ]]; then
    echo "Release notes ${RELEASE_NOTES_FILE} not found. Exiting..."
    exit 1
  fi

  goreleaser release --skip-validate --config=./.goreleaser.yml --release-notes="${RELEASE_NOTES_FILE}"

  # By moving the latest_release tag to the latest release candidate we ensure the rc is accessible through
  # the `head` statement in the brew tap formula
  sleep 90 # GitHub API resolves the time to the nearest minute, so in order to control the sorting oder we need this

  git tag --delete "${CIRCLE_TAG}"
  git tag --force latest_release

  if github-release info --user weaveworks --repo "${CIRCLE_PROJECT_REPONAME}" --tag latest_release > /dev/null 2>&1 ; then
    github-release delete --user weaveworks --repo "${CIRCLE_PROJECT_REPONAME}" --tag latest_release
  fi

  export RELEASE_DESCRIPTION="${CIRCLE_TAG}"
  goreleaser release --skip-validate --rm-dist --config=./.goreleaser.yml --release-notes="${RELEASE_NOTES_FILE}"

else
  echo "Not a tag release, skip publish"
fi
