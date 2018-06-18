#!/bin/sh

if [ -z "${CIRCLE_PULL_REQUEST}" ] && [ "${CIRCLE_TAG}" ] && [ "${CIRCLE_PROJECT_USERNAME}" = "weaveworks" ] ; then
  export RELEASE_DESCRIPTION="${CIRCLE_TAG}"
  goreleaser release --skip-validate --config=./.goreleaser.yml

  git tag -d "${CIRCLE_TAG}"
  gig tag latest_release

  if github-release info --user weaveworks --repo eksctl --tag latest_release > /dev/null 2>&1 ; then
    github-release delete --user weaveworks --repo eksctl --tag latest_release
  fi

  export RELEASE_DESCRIPTION="${RELEASE_DESCRIPTION} (permalink)"
  goreleaser release --config=./.goreleaser.yml
else
  echo "Not a tag release, skip publish"
fi
