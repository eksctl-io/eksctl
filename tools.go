// +build tools

package eksctl

// Mock imports to enforce their installation by `go mod`.
// See https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
import (
	_ "github.com/dave/jennifer/jen"
	_ "github.com/jteeuwen/go-bindata/go-bindata"
	_ "github.com/mattn/goveralls"
	_ "github.com/vektra/mockery/cmd/mockery"
	_ "github.com/weaveworks/github-release"
	_ "golang.org/x/tools/cmd/stringer"
	_ "k8s.io/code-generator"
)
