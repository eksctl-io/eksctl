#!/bin/sh -eux

# Make sure to run the following commands after changes to this file are made:
# `make -f Makefile.docker check-build-image-manifest-up-to-date && make -f Makefile.docker push-build-image`

if [ -z "${GOBIN+x}" ]; then
 GOBIN="${GOPATH%%:*}/bin"
fi

if [ "$(uname)" = "Darwin" ] ; then
  OSARCH="darwin-amd64"
else
  OSARCH="linux-amd64"
fi

REQUIREMENTS_FILE=.requirements

if [ ! -f "$REQUIREMENTS_FILE" ]
then
	echo "Requirements file $REQUIREMENTS_FILE not found. Exiting..."
	exit 1
fi

# Install all other Go build requirements
while IFS= read -r req; do
  go install "${req}"
done < ${REQUIREMENTS_FILE}
