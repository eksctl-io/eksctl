// +build tools

package eksctl

// Mock imports to enforce their installation by `go mod`.
// See https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
import (
	// These dependencies are only included to enforce aws-iam-authenticator to use a specific versions.
	// aws-iam-authenticator v0.4.0 uses `dep`, causing `go mod` to ignore its dependency files and just to pull
	// master versions. Otherwise, build to breaks and is unpredictable
	// TODO(fons): get rid of this dependencies once we bump the version of aws-iam-authenticator,
	//             which now uses go modules
	_ "github.com/christopherhein/go-version"
	_ "github.com/spf13/viper"

	_ "github.com/dave/jennifer/jen"
	_ "github.com/gohugoio/hugo"
	_ "github.com/goreleaser/goreleaser"
	_ "github.com/kevinburke/go-bindata/go-bindata"
	_ "github.com/kubernetes-sigs/aws-iam-authenticator/cmd/aws-iam-authenticator"
	_ "github.com/vektra/mockery/cmd/mockery"
	_ "github.com/weaveworks/github-release"
	_ "golang.org/x/tools/cmd/stringer"
	_ "k8s.io/code-generator/cmd/client-gen"
	_ "k8s.io/code-generator/cmd/deepcopy-gen"
	_ "k8s.io/code-generator/cmd/defaulter-gen"
	_ "k8s.io/code-generator/cmd/informer-gen"
	_ "k8s.io/code-generator/cmd/lister-gen"
)
