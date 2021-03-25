package builder

import (
	"fmt"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	cft "github.com/weaveworks/eksctl/pkg/cfn/template"
	gfniam "github.com/weaveworks/goformation/v4/cloudformation/iam"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

const (
	fargateTemplateDescription = "Fargate IAM"
	fargateRoleName            = "FargatePodExecutionRole"
	fargateRoleDescription     = "EKS Fargate pod execution IAM role [created by eksctl]"
)

// FargateResourceSet manages only fargate resources
type FargateResourceSet struct {
	rs   *resourceSet
	spec *api.ClusterConfig
}

// NewFargateResourceSet returns a resource set for managing fargate resources
func NewFargateResourceSet(spec *api.ClusterConfig) *FargateResourceSet {
	rs := newResourceSet()
	rs.withIAM = true
	rs.withNamedIAM = true
	return &FargateResourceSet{
		rs,
		spec,
	}
}

func (rs *FargateResourceSet) AddAllResources() error {
	rs.rs.template.Mappings[servicePrincipalPartitionMapName] = servicePrincipalPartitionMappings

	rs.rs.template.Description = fmt.Sprintf(
		"%s %s",
		fargateTemplateDescription,
		templateDescriptionSuffix,
	)
	return addResourcesForFargate(rs.rs, rs.spec)
}

func (rs *FargateResourceSet) WithIAM() bool {
	return true
}

func (rs *FargateResourceSet) WithNamedIAM() bool {
	return true
}

func (rs *FargateResourceSet) RenderJSON() ([]byte, error) {
	return rs.rs.renderJSON()
}

func (rs *FargateResourceSet) GetAllOutputs(stack cfn.Stack) error {
	return rs.rs.GetAllOutputs(stack)
}

// addResourcesForFargate adds resources for Fargate.
func addResourcesForFargate(rs *resourceSet, cfg *api.ClusterConfig) error {
	if api.IsSetAndNonEmptyString(cfg.IAM.FargatePodExecutionRoleARN) {
		rs.defineOutputWithoutCollector(outputs.FargatePodExecutionRoleARN, cfg.IAM.FargatePodExecutionRoleARN, true)
		return nil
	}
	// Create a role requires additional capabilities.
	// If not set to true, CloudFormation fails with:
	//   status code 400: InsufficientCapabilitiesException: Requires capabilities : [CAPABILITY_IAM]
	rs.withIAM = true

	rs.template.Description = fargateRoleDescription
	role := &gfniam.Role{
		AssumeRolePolicyDocument: cft.MakeAssumeRolePolicyDocumentForServices(
			MakeServiceRef("EKSFargatePods"), // Ensure that EKS can schedule pods onto Fargate.
		),
		ManagedPolicyArns: gfnt.NewSlice(makePolicyARNs(
			iamPolicyAmazonEKSFargatePodExecutionRolePolicy,
		)...),
	}

	if api.IsSetAndNonEmptyString(cfg.IAM.FargatePodExecutionRolePermissionsBoundary) {
		role.PermissionsBoundary = gfnt.NewString(*cfg.IAM.FargatePodExecutionRolePermissionsBoundary)
	}

	rs.newResource(fargateRoleName, role)
	rs.defineOutputFromAtt(outputs.FargatePodExecutionRoleARN, fargateRoleName, "Arn", true, func(v string) error {
		cfg.IAM.FargatePodExecutionRoleARN = &v
		return nil
	})
	return nil
}
