#!/bin/sh -eu

if [ -z "${GOBIN+x}" ]; then
 GOBIN="$(go env GOPATH)/bin"
fi

if [ "$(uname)" = "Darwin" ] ; then
  OSARCH="darwin-amd64"
else
  OSARCH="linux-amd64"
fi

go install github.com/jteeuwen/go-bindata/go-bindata
go install github.com/weaveworks/github-release
go install golang.org/x/tools/cmd/stringer
go install github.com/mattn/goveralls
go install github.com/vektra/mockery/cmd/mockery

# Install metalinter
# Managing all linters that gometalinter uses with dep is going to take
# a lot of work, so we install all of those from the release tarball
VERSION="2.0.11"
curl --silent --location "https://github.com/alecthomas/gometalinter/releases/download/v${VERSION}/gometalinter-${VERSION}-${OSARCH}.tar.gz" | \
 tar -x -z -C "${GOBIN}" --strip-components 1

