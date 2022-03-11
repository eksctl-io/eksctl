package generate

import (
	_ "github.com/aws/aws-sdk-go-v2/service/autoscaling"
	_ "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	_ "github.com/aws/aws-sdk-go-v2/service/sts"
)

//go:generate ../../../build/scripts/generate-aws-interfaces.sh sts STS
//go:generate ../../../build/scripts/generate-aws-interfaces.sh autoscaling ASG
//go:generate ../../../build/scripts/generate-aws-interfaces.sh cloudformation CloudFormation
