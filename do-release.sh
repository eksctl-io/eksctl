#!/bin/sh -ex

if [ -z "${CIRCLE_PULL_REQUEST}" ] && [ -n "${CIRCLE_TAG}" ] && [ "${CIRCLE_PROJECT_USERNAME}" = "weaveworks" ] ; then
  export RELEASE_DESCRIPTION="${CIRCLE_TAG} (permalink)"
  cat ./.goreleaser.yml ./.goreleaser.brew.yml > .goreleaser.brew.combined.yml
  goreleaser release --skip-validate --config=./.goreleaser.brew.combined.yml

  sleep 90 # GitHub API resolves the time to the nearest minute, so in order to control the sorting oder we need this

  git tag --delete "${CIRCLE_TAG}"
  git tag --force latest_release

  if github-release info --user weaveworks --repo eksctl --tag latest_release > /dev/null 2>&1 ; then
    github-release delete --user weaveworks --repo eksctl --tag latest_release
  fi

  export RELEASE_DESCRIPTION="${CIRCLE_TAG}"
  goreleaser release --skip-validate --rm-dist --config=./.goreleaser.yml

else
  echo "Not a tag release, skip publish"
fi
