#!/bin/sh -eu

if [ -z "${GOBIN+x}" ]; then
 GOBIN="$(go env GOPATH)/bin"
fi

if [ "$(uname)" = "Darwin" ] ; then
  OSARCH="darwin-amd64"
else
  OSARCH="linux-amd64"
fi

go mod download
go install github.com/jteeuwen/go-bindata/go-bindata
go install github.com/weaveworks/github-release
go install golang.org/x/tools/cmd/stringer
go install github.com/mattn/goveralls
go install github.com/vektra/mockery/cmd/mockery


# TODO: metalinter is archived, we should switch to github.com/golangci/golangci-lint
# Install metalinter
# Managing all linters that gometalinter uses with dep is going to take
# a lot of work, so we install all of those from the release tarball
VERSION="3.0.0"
curl --silent --location "https://github.com/alecthomas/gometalinter/releases/download/v${VERSION}/gometalinter-${VERSION}-${OSARCH}.tar.gz" | \
 tar -x -z -C "${GOBIN}" --strip-components 1

