package mocks

import (
	_ "github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface" // used for testing
	_ "github.com/aws/aws-sdk-go/service/cloudtrail/cloudtrailiface"
	_ "github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	_ "github.com/aws/aws-sdk-go/service/eks/eksiface"
	_ "github.com/aws/aws-sdk-go/service/elb/elbiface"
	_ "github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	_ "github.com/aws/aws-sdk-go/service/iam/iamiface"
	_ "github.com/aws/aws-sdk-go/service/sts/stsiface"
	_ "github.com/vektra/mockery"
)

// Run make check-all-generated-files-up-to-date to generate the mocks
//go:generate "${GOBIN}/mockery" -tags netgo -dir=${AWS_SDK_GO_DIR}/service/cloudformation/cloudformationiface -name=CloudFormationAPI -output=./
//go:generate "${GOBIN}/mockery" -tags netgo -dir=${AWS_SDK_GO_DIR}/service/eks/eksiface -name=EKSAPI -output=./
//go:generate "${GOBIN}/mockery" -tags netgo -dir=${AWS_SDK_GO_DIR}/service/ec2/ec2iface -name=EC2API -output=./
//go:generate "${GOBIN}/mockery" -tags netgo -dir=${AWS_SDK_GO_DIR}/service/elb/elbiface -name=ELBAPI -output=./
//go:generate "${GOBIN}/mockery" -tags netgo -dir=${AWS_SDK_GO_DIR}/service/elbv2/elbv2iface -name=ELBV2API -output=./
//go:generate "${GOBIN}/mockery" -tags netgo -dir=${AWS_SDK_GO_DIR}/service/sts/stsiface -name=STSAPI -output=./
//go:generate "${GOBIN}/mockery" -tags netgo -dir=${AWS_SDK_GO_DIR}/service/iam/iamiface -name=IAMAPI -output=./
//go:generate "${GOBIN}/mockery" -tags netgo -dir=${AWS_SDK_GO_DIR}/service/cloudtrail/cloudtrailiface -name=CloudTrailAPI -output=./
//go:generate "${GOBIN}/mockery" -tags netgo -dir=${AWS_SDK_GO_DIR}/service/ssm/ssmiface -name=SSMAPI -output=./
