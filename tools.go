// +build tools

package eksctl

// Mock imports to enforce their installation by `go mod`.
// See https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
import (
	// This go-version dependency only included to enforce aws-iam-authenticator to use a specific version of go-version.
	// aws-iam-authenticator v0.4.0 uses `dep`, causing `go mod` to ignore its dependency files and just to pull
	// github.com/christopherhein/go-version:master. This causes the build to break because go-version:master has recently
	// switched to module name go.hein.dev/go-version, causing a module name mismatch.
	// TODO(fons): get rid of this dependency once we bump the version of aws-iam-authenticator
	//             which now uses go modules
	_ "github.com/christopherhein/go-version"
	_ "github.com/dave/jennifer/jen"
	_ "github.com/goreleaser/goreleaser"
	_ "github.com/jteeuwen/go-bindata/go-bindata"

	_ "github.com/kubernetes-sigs/aws-iam-authenticator/cmd/aws-iam-authenticator"
	_ "github.com/mattn/goveralls"
	_ "github.com/vektra/mockery/cmd/mockery"
	_ "github.com/weaveworks/github-release"
	_ "golang.org/x/tools/cmd/stringer"
	_ "k8s.io/code-generator"
)
