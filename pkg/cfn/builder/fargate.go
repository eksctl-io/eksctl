package builder

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
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

func (rs *FargateResourceSet) GetAllOutputs(stack types.Stack) error {
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

	// As per AWS docs, to avoid a confused deputy security problem, it's important that the role restricts access based on SourceArn.
	// See: https://docs.aws.amazon.com/eks/latest/userguide/pod-execution-role.html#check-pod-execution-role
	sourceArnCondition, err := makeSourceArnCondition(cfg)
	if err != nil {
		return fmt.Errorf("restricting access based on SourceArn: %w", err)
	}

	role := &gfniam.Role{
		AssumeRolePolicyDocument: cft.MakeAssumeRolePolicyDocumentForServicesWithConditions(
			sourceArnCondition,
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

func makeSourceArnCondition(cfg *api.ClusterConfig) (cft.MapOfInterfaces, error) {
	accountID, err := getAWSAccountID(cfg)
	if err != nil {
		return nil, err
	}
	return cft.MapOfInterfaces{
		"ArnLike": cft.MapOfInterfaces{
			"aws:SourceArn": fmt.Sprintf("arn:aws:eks:%s:%s:fargateprofile/%s/*", cfg.Metadata.Region, accountID, cfg.Metadata.Name),
		},
	}, nil
}

func getAWSAccountID(cfg *api.ClusterConfig) (string, error) {
	if cfg.Metadata.AccountID != "" {
		return cfg.Metadata.AccountID, nil
	}
	if cfg.Status != nil && cfg.Status.ARN != "" {
		parsedARN, err := arn.Parse(cfg.Status.ARN)
		if err != nil {
			return "", fmt.Errorf("error parsing cluster ARN: %v", err)
		}
		return parsedARN.AccountID, nil
	}
	return "", fmt.Errorf("failed to determine account ID")
}
