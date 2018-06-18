#!/bin/sh

if [[ -z "${CIRCLE_PULL_REQUEST}" ]] && [[ "${CIRCLE_TAG}" ]] && [[ "${CIRCLE_PROJECT_USERNAME}" = "weaveworks" ]] ; then
  make release
else
  echo "Not a tag release, skip publish"
fi
