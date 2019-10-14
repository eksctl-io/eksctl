#!/bin/sh -eux

# Make sure to run the following commands after changes to this file are made:
# `make -f Makefile.docker update-build-image-manifest && make -f Makefile.docker push-build-image`

if [ -z "${GOBIN+x}" ]; then
 GOBIN="$(go env GOPATH)/bin"
fi

if [ "$(uname)" = "Darwin" ] ; then
  OSARCH="darwin-amd64"
else
  OSARCH="linux-amd64"
fi

env CGO_ENABLED=1 go install -tags extended github.com/gohugoio/hugo

go install \
  github.com/goreleaser/goreleaser \
  github.com/kevinburke/go-bindata/go-bindata \
  github.com/vektra/mockery/cmd/mockery \
  github.com/weaveworks/github-release \
  golang.org/x/tools/cmd/stringer \
  k8s.io/code-generator/cmd/client-gen \
  k8s.io/code-generator/cmd/deepcopy-gen \
  k8s.io/code-generator/cmd/defaulter-gen \
  k8s.io/code-generator/cmd/informer-gen \
  k8s.io/code-generator/cmd/lister-gen

# TODO: metalinter is archived, we should switch to github.com/golangci/golangci-lint
# Install metalinter
# Managing all linters that gometalinter uses with dep is going to take
# a lot of work, so we install all of those from the release tarball
METALINTER_VERSION="3.0.0"
curl --silent --location "https://github.com/alecthomas/gometalinter/releases/download/v${METALINTER_VERSION}/gometalinter-${METALINTER_VERSION}-${OSARCH}.tar.gz" \
  | tar -x -z -C "${GOBIN}" --strip-components 1
