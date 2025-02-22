//go:build tools

package eksctl

// Mock imports to enforce their installation by `go mod`.
// See https://github.com/golang/go/wiki/Modules#how-can-i-track-tool-dependencies-for-a-module
import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
	_ "github.com/maxbrunsfeld/counterfeiter/v6"
	_ "github.com/vburenin/ifacemaker"
	_ "github.com/vektra/mockery/v2"
)
