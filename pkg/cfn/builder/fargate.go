package builder

import (
	"fmt"

	gfn "github.com/awslabs/goformation/cloudformation"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	cft "github.com/weaveworks/eksctl/pkg/cfn/template"
)

const (
	fargateRoleName        = "FargatePodExecutionRole"
	fargateRoleDescription = "EKS Fargate pod execution IAM role [created by eksctl]"
)

// AddResourcesForFargate adds resources for Fargate.
func AddResourcesForFargate(rs *resourceSet, cfg *api.ClusterConfig) error {
	if api.IsSetAndNonEmptyString(cfg.IAM.FargatePodExecutionRoleARN) {
		rs.defineOutputWithoutCollector(outputs.FargatePodExecutionRoleARN, cfg.IAM.FargatePodExecutionRoleARN, true)
		return nil
	}
	// Create a role requires additional capabilities.
	// If not set to true, CloudFormation fails with:
	//   status code 400: InsufficientCapabilitiesException: Requires capabilities : [CAPABILITY_IAM]
	rs.withIAM = true

	rs.template.Description = fargateRoleDescription
	rs.newResource(fargateRoleName, &gfn.AWSIAMRole{
		AssumeRolePolicyDocument: cft.MakeAssumeRolePolicyDocumentForServices(
			"eks.amazonaws.com",
			"eks-fargate-pods.amazonaws.com", // Ensure that EKS can schedule pods onto Fargate.
		),
		ManagedPolicyArns: makeStringSlice(
			iamPolicyAmazonEKSFargatePodExecutionRolePolicyARN,
		),
	})
	rs.defineOutputFromAtt(outputs.FargatePodExecutionRoleARN, fmt.Sprintf("%s.Arn", fargateRoleName), true, func(v string) error {
		cfg.IAM.FargatePodExecutionRoleARN = &v
		return nil
	})
	return nil
}
