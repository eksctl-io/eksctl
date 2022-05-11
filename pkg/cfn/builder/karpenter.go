package builder

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	gfn "github.com/weaveworks/goformation/v4/cloudformation"
	gfniam "github.com/weaveworks/goformation/v4/cloudformation/iam"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	cft "github.com/weaveworks/eksctl/pkg/cfn/template"
)

const (
	// KarpenterNodeRoleName is the name of the role for nodes.
	KarpenterNodeRoleName = "KarpenterNodeRole"
	// KarpenterManagedPolicy managed policy name.
	KarpenterManagedPolicy = "KarpenterControllerPolicy"
	// KarpenterNodeInstanceProfile is the name of node instance profile.
	KarpenterNodeInstanceProfile = "KarpenterNodeInstanceProfile"
)

const (
	// actions
	// EC2
	ec2CreateFleet                   = "ec2:CreateFleet"
	ec2CreateLaunchTemplate          = "ec2:CreateLaunchTemplate"
	ec2CreateTags                    = "ec2:CreateTags"
	ec2DescribeAvailabilityZones     = "ec2:DescribeAvailabilityZones"
	ec2DescribeInstanceTypeOfferings = "ec2:DescribeInstanceTypeOfferings"
	ec2DescribeInstanceTypes         = "ec2:DescribeInstanceTypes"
	ec2DescribeInstances             = "ec2:DescribeInstances"
	ec2DescribeLaunchTemplates       = "ec2:DescribeLaunchTemplates"
	ec2DescribeSecurityGroups        = "ec2:DescribeSecurityGroups"
	ec2DescribeSubnets               = "ec2:DescribeSubnets"
	ec2DeleteLaunchTemplate          = "ec2:DeleteLaunchTemplate"
	ec2RunInstances                  = "ec2:RunInstances"
	ec2TerminateInstances            = "ec2:TerminateInstances"
	// IAM
	iamPassRole     = "iam:PassRole"
	ssmGetParameter = "ssm:GetParameter"
)

// KarpenterResourceSet stores the resource information of the Karpenter stack
type KarpenterResourceSet struct {
	rs                  *resourceSet
	clusterSpec         *api.ClusterConfig
	instanceProfileName string
}

// NewKarpenterResourceSet returns a resource set for a Karpenter embedded in a cluster config
func NewKarpenterResourceSet(spec *api.ClusterConfig, instanceProfileName string) *KarpenterResourceSet {
	return &KarpenterResourceSet{
		rs:                  newResourceSet(),
		clusterSpec:         spec,
		instanceProfileName: instanceProfileName,
	}
}

// AddAllResources adds all the information about Karpenter to the resource set
func (k *KarpenterResourceSet) AddAllResources() error {
	k.rs.template.Description = fmt.Sprintf("Karpenter Stack %s", templateDescriptionSuffix)
	return k.addResourcesForKarpenter()
}

// RenderJSON returns the rendered JSON
func (k *KarpenterResourceSet) RenderJSON() ([]byte, error) {
	return k.rs.renderJSON()
}

// Template returns the CloudFormation template
func (k *KarpenterResourceSet) Template() gfn.Template {
	return *k.rs.template
}

func (k *KarpenterResourceSet) newResource(name string, resource gfn.Resource) *gfnt.Value {
	return k.rs.newResource(name, resource)
}

func (k *KarpenterResourceSet) addResourcesForKarpenter() error {
	managedPolicyNames := sets.NewString()
	managedPolicyNames.Insert(iamPolicyAmazonEKSWorkerNodePolicy,
		iamPolicyAmazonEKSCNIPolicy,
		iamPolicyAmazonEC2ContainerRegistryReadOnly,
		iamPolicyAmazonSSMManagedInstanceCore,
	)
	k.Template().Mappings[servicePrincipalPartitionMapName] = servicePrincipalPartitionMappings
	roleName := gfnt.NewString(fmt.Sprintf("eksctl-%s-%s", KarpenterNodeRoleName, k.clusterSpec.Metadata.Name))
	role := gfniam.Role{
		RoleName:                 roleName,
		Path:                     gfnt.NewString("/"),
		AssumeRolePolicyDocument: cft.MakeAssumeRolePolicyDocumentForServices(MakeServiceRef("EC2")),
		ManagedPolicyArns:        gfnt.NewSlice(makePolicyARNs(managedPolicyNames.List()...)...),
	}

	if api.IsSetAndNonEmptyString(k.clusterSpec.IAM.ServiceRolePermissionsBoundary) {
		role.PermissionsBoundary = gfnt.NewString(*k.clusterSpec.IAM.ServiceRolePermissionsBoundary)
	}

	roleRef := k.newResource(KarpenterNodeRoleName, &role)

	instanceProfile := gfniam.InstanceProfile{
		InstanceProfileName: gfnt.NewString(k.instanceProfileName),
		Path:                gfnt.NewString("/"),
		Roles:               gfnt.NewSlice(roleRef),
	}
	k.newResource(KarpenterNodeInstanceProfile, &instanceProfile)

	managedPolicyName := gfnt.NewString(fmt.Sprintf("eksctl-%s-%s", KarpenterManagedPolicy, k.clusterSpec.Metadata.Name))
	statements := cft.MapOfInterfaces{
		"Effect":   effectAllow,
		"Resource": resourceAll,
		"Action": []string{
			ec2CreateFleet,
			ec2CreateLaunchTemplate,
			ec2CreateTags,
			ec2DescribeAvailabilityZones,
			ec2DescribeInstanceTypeOfferings,
			ec2DescribeInstanceTypes,
			ec2DescribeInstances,
			ec2DescribeLaunchTemplates,
			ec2DescribeSecurityGroups,
			ec2DescribeSubnets,
			ec2DeleteLaunchTemplate,
			ec2RunInstances,
			ec2TerminateInstances,
			iamPassRole,
			ssmGetParameter,
		},
	}
	managedPolicy := gfniam.ManagedPolicy{
		ManagedPolicyName: managedPolicyName,
		PolicyDocument:    cft.MakePolicyDocument(statements),
	}
	k.newResource(KarpenterManagedPolicy, &managedPolicy)
	return nil
}

// WithIAM implements the ResourceSet interface
func (k *KarpenterResourceSet) WithIAM() bool {
	// eksctl does not support passing pre-created IAM instance roles to Managed Nodes,
	// so the IAM capability is always required
	return true
}

// WithNamedIAM implements the ResourceSet interface
func (k *KarpenterResourceSet) WithNamedIAM() bool {
	return true
}

// GetAllOutputs collects all outputs of the nodegroup
func (k *KarpenterResourceSet) GetAllOutputs(stack types.Stack) error {
	return k.rs.GetAllOutputs(stack)
}
