#!/bin/sh -ex

if [ -z "${CIRCLE_PROJECT_REPONAME}" ] ; then
  echo "Missing repo name, please set CIRCLE_PROJECT_REPONAME"
  exit 1
fi

if [ -z "${CIRCLE_PULL_REQUEST}" ] && [ -n "${CIRCLE_TAG}" ] && [ "${CIRCLE_PROJECT_USERNAME}" = "weaveworks" ] ; then
  export RELEASE_DESCRIPTION="${CIRCLE_TAG} (permalink)"
  RELEASE_NOTES_FILE="docs/release_notes/${CIRCLE_TAG}.md"

  if [[ ! -f "${RELEASE_NOTES_FILE}" ]]; then
    echo "Release notes ${RELEASE_NOTES_FILE} not found. Exiting..."
    exit 1
  fi

  cat ./.goreleaser.yml ./.goreleaser.brew.yml > .goreleaser.brew.combined.yml
  goreleaser release --rm-dist --timeout 60m --skip-validate --config=./.goreleaser.brew.combined.yml --release-notes="${RELEASE_NOTES_FILE}"

  sleep 90 # GitHub API resolves the time to the nearest minute, so in order to control the sorting oder we need this

  docker login --username weaveworkseksctlci --password "${DOCKER_HUB_PASSWORD}"
  EKSCTL_IMAGE_VERSION="${CIRCLE_TAG}" make -f Makefile.docker push-eksctl-image || true
else
  echo "Not a tag release, skip publish"
fi
