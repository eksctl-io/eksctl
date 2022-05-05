package mocks

import (
	_ "github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface" // used for testing
	_ "github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	_ "github.com/aws/aws-sdk-go/service/cloudtrail/cloudtrailiface"
	_ "github.com/aws/aws-sdk-go/service/eks/eksiface"
	_ "github.com/vektra/mockery"
)

// Run make check-all-generated-files-up-to-date to generate the mocks
//go:generate "${GOBIN}/mockery" --tags netgo --dir=${AWS_SDK_GO_DIR}/service/autoscaling/autoscalingiface --name=AutoScalingAPI --output=./
//go:generate "${GOBIN}/mockery" --tags netgo --dir=${AWS_SDK_GO_DIR}/service/cloudformation/cloudformationiface --name=CloudFormationAPI --output=./
//go:generate "${GOBIN}/mockery" --tags netgo --dir=${AWS_SDK_GO_DIR}/service/cloudtrail/cloudtrailiface --name=CloudTrailAPI --output=./
//go:generate "${GOBIN}/mockery" --tags netgo --dir=${AWS_SDK_GO_DIR}/service/cloudwatchlogs/cloudwatchlogsiface --name=CloudWatchLogsAPI --output=./
//go:generate "${GOBIN}/mockery" --tags netgo --dir=${AWS_SDK_GO_DIR}/aws/client --name=ConfigProvider --output=./
