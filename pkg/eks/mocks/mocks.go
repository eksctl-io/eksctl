package mocks

import (
	_ "github.com/vektra/mockery"
)

// Run make check-all-generated-files-up-to-date to generate the mocks
//go:generate "${GOBIN}/mockery" --tags netgo --dir=${AWS_SDK_GO_DIR}/aws/client --name=ConfigProvider --output=./
