package mocks

import _ "github.com/vektra/mockery" // used for testing

//go:generate ${GOPATH}/bin/mockery -tags netgo -dir=../../../vendor/github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface -name=CloudFormationAPI -output=./
//go:generate ${GOPATH}/bin/mockery -tags netgo -dir=../../../vendor/github.com/aws/aws-sdk-go/service/eks/eksiface -name=EKSAPI -output=./
//go:generate ${GOPATH}/bin/mockery -tags netgo -dir=../../../vendor/github.com/aws/aws-sdk-go/service/ec2/ec2iface -name=EC2API -output=./
//go:generate ${GOPATH}/bin/mockery -tags netgo -dir=../../../vendor/github.com/aws/aws-sdk-go/service/sts/stsiface -name=STSAPI -output=./
