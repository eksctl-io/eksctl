#!/bin/sh -eu

go install ./vendor/github.com/jteeuwen/go-bindata/go-bindata
go install ./vendor/github.com/weaveworks/github-release
go install ./vendor/golang.org/x/tools/cmd/stringer
go install ./vendor/github.com/mattn/goveralls
go install ./vendor/github.com/vektra/mockery/cmd/mockery

curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $GOPATH/bin v1.10.2
