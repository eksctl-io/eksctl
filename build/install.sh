#!/bin/sh -eu

go install ./vendor/github.com/jteeuwen/go-bindata/go-bindata
go install ./vendor/github.com/weaveworks/github-release
go install ./vendor/golang.org/x/tools/cmd/stringer
go install ./vendor/github.com/mattn/goveralls
go install ./vendor/github.com/vektra/mockery/cmd/mockery

# managing all linters that gometalinter uses with dep is going to take
# a lot of work, so we install all of those from the release tarball
install_gometalinter() {
  version="${1}"
  prefix="https://github.com/alecthomas/gometalinter/releases/download"
  if [ "$(uname)" = "Darwin" ] ; then
    suffix="darwin-amd64"
  else
    suffix="linux-amd64"
  fi
  basename="gometalinter-${version}-${suffix}"
  url="${prefix}/v${version}/${basename}.tar.gz"
  cd "${GOPATH}/bin/"
  curl --silent --location "${url}" | tar xz
  (cd "./${basename}/" ; mv ./* ../)
  rmdir "./${basename}"
  unset version prefix suffix basename url
}

install_golangci_lint() {
  version="${1}"
  curl --silent --fail --location \
    "https://install.goreleaser.com/github.com/golangci/golangci-lint.sh" \
    | sh -s -- -b $GOPATH/bin "${version}"
  unset version
}

install_gometalinter "2.0.11"
install_golangci_lint "v1.10.2"