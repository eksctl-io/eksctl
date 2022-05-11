package generate

//go:generate ../../../build/scripts/generate-aws-interfaces.sh sts STS
//go:generate ../../../build/scripts/generate-aws-interfaces.sh autoscaling ASG
//go:generate ../../../build/scripts/generate-aws-interfaces.sh cloudwatchlogs CloudWatchLogs
//go:generate ../../../build/scripts/generate-aws-interfaces.sh cloudformation CloudFormation
//go:generate ../../../build/scripts/generate-aws-interfaces.sh cloudtrail CloudTrail
//go:generate ../../../build/scripts/generate-aws-interfaces.sh elasticloadbalancing ELB
//go:generate ../../../build/scripts/generate-aws-interfaces.sh elasticloadbalancingv2 ELBV2
//go:generate ../../../build/scripts/generate-aws-interfaces.sh ssm SSM
//go:generate ../../../build/scripts/generate-aws-interfaces.sh iam IAM
//go:generate ../../../build/scripts/generate-aws-interfaces.sh eks EKS
