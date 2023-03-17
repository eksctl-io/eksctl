//go:build tools

package eksctl

// Mock imports to enforce their installation by `go mod`.
// See https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
import (
	_ "github.com/dave/dst"
	_ "github.com/dave/jennifer/jen"
	_ "github.com/github-release/github-release"
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/goreleaser/goreleaser"
	_ "github.com/maxbrunsfeld/counterfeiter/v6"
	_ "github.com/vburenin/ifacemaker"
	_ "github.com/vektra/mockery/cmd/mockery"
	_ "golang.org/x/tools/cmd/stringer"
	_ "k8s.io/code-generator/cmd/client-gen"
	_ "k8s.io/code-generator/cmd/deepcopy-gen"
	_ "k8s.io/code-generator/cmd/defaulter-gen"
	_ "k8s.io/code-generator/cmd/informer-gen"
	_ "k8s.io/code-generator/cmd/lister-gen"
	_ "sigs.k8s.io/mdtoc"
)
