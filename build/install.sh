#!/bin/sh

go install ./vendor/github.com/jteeuwen/go-bindata/go-bindata
go install ./vendor/github.com/weaveworks/github-release
go install ./vendor/golang.org/x/tools/cmd/stringer
go install ./vendor/github.com/mattn/goveralls
go install ./vendor/github.com/vektra/mockery/cmd/mockery
