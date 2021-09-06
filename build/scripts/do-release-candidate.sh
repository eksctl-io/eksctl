#!/bin/sh -ex

if [ -z "${CIRCLE_PULL_REQUEST}" ] && [ -n "${CIRCLE_TAG}" ] && [ "${CIRCLE_PROJECT_USERNAME}" = "weaveworks" ] ; then
  export RELEASE_DESCRIPTION="${CIRCLE_TAG}"
  RELEASE_NOTES_FILE="docs/release_notes/${CIRCLE_TAG/-rc.*}.md"

  if [[ ! -f "${RELEASE_NOTES_FILE}" ]]; then
    echo "Release notes ${RELEASE_NOTES_FILE} not found. Exiting..."
    exit 1
  fi

  goreleaser release --rm-dist --timeout 60m --skip-validate --config=./.goreleaser.yml --release-notes="${RELEASE_NOTES_FILE}"

  sleep 90 # GitHub API resolves the time to the nearest minute, so in order to control the sorting oder we need this

  docker login --username weaveworkseksctlci --password "${DOCKER_HUB_PASSWORD}"
  EKSCTL_IMAGE_VERSION="${CIRCLE_TAG}" make -f Makefile.docker push-eksctl-image || true
else
  echo "Not a tag release, skip publish"
fi
