package mocks

import (
	_ "github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface" // used for testing
	_ "github.com/aws/aws-sdk-go/service/cloudtrail/cloudtrailiface"
	_ "github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	_ "github.com/aws/aws-sdk-go/service/eks/eksiface"
	_ "github.com/aws/aws-sdk-go/service/iam/iamiface"
	_ "github.com/aws/aws-sdk-go/service/sts/stsiface"
	_ "github.com/vektra/mockery"
)

// "../../../vendor/aws-sdk-go" is artificially created in the main Makefile to make mockery work properly with go modules
//go:generate "${GOBIN}/mockery" -tags netgo -dir=../../../vendor/github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface -name=CloudFormationAPI -output=./
//go:generate "${GOBIN}/mockery" -tags netgo -dir=../../../vendor/github.com/aws/aws-sdk-go/service/eks/eksiface -name=EKSAPI -output=./
//go:generate "${GOBIN}/mockery" -tags netgo -dir=../../../vendor/github.com/aws/aws-sdk-go/service/ec2/ec2iface -name=EC2API -output=./
//go:generate "${GOBIN}/mockery" -tags netgo -dir=../../../vendor/github.com/aws/aws-sdk-go/service/sts/stsiface -name=STSAPI -output=./
//go:generate "${GOBIN}/mockery" -tags netgo -dir=../../../vendor/github.com/aws/aws-sdk-go/service/iam/iamiface -name=IAMAPI -output=./
//go:generate "${GOBIN}/mockery" -tags netgo -dir=../../../vendor/github.com/aws/aws-sdk-go/service/cloudtrail/cloudtrailiface -name=CloudTrailAPI -output=./
