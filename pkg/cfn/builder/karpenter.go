package builder

import (
	"fmt"

	gfn "goformation/v4/cloudformation"
	gfnevents "goformation/v4/cloudformation/events"
	gfniam "goformation/v4/cloudformation/iam"
	gfnsqs "goformation/v4/cloudformation/sqs"
	gfnt "goformation/v4/cloudformation/types"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
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
	// KarpenterInterruptionQueue is the name of the interruption queue
	KarpenterInterruptionQueue = "KarpenterInterruptionQueue"
	// KarpenterInterruptionQueuePolicy interruption queue policy name
	KarpenterInterruptionQueuePolicy = "KarpenterInterruptionQueuePolicy"
	// KarpenterInterruptionQueueTarget interruption queue target ID
	KarpenterInterruptionQueueTarget = "KarpenterInterruptionQueueTarget"
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
	ec2DescribeImages                = "ec2:DescribeImages"
	ec2DescribeSpotPriceHistory      = "ec2:DescribeSpotPriceHistory"
	// IAM
	iamPassRole                 = "iam:PassRole"
	iamCreateServiceLinkedRole  = "iam:CreateServiceLinkedRole"
	iamGetInstanceProfile       = "iam:GetInstanceProfile"
	iamCreateInstanceProfile    = "iam:CreateInstanceProfile"
	iamDeleteInstanceProfile    = "iam:DeleteInstanceProfile"
	iamTagInstanceProfile       = "iam:TagInstanceProfile"
	iamAddRoleToInstanceProfile = "iam:AddRoleToInstanceProfile"
	// SSM
	ssmGetParameter = "ssm:GetParameter"
	// Pricing
	pricingGetProducts = "pricing:GetProducts"
	// SQS
	sqsDeleteMessage      = "sqs:DeleteMessage"
	sqsGetQueueAttributes = "sqs:GetQueueAttributes"
	sqsGetQueueURL        = "sqs:GetQueueUrl"
	sqsReceiveMessage     = "sqs:ReceiveMessage"
	sqsSendMessage        = "sqs:SendMessage"
)

const (
	ScheduledChangeRule     = "ScheduledChangeRule"
	SpotInterruptionRule    = "SpotInterruptionRule"
	RebalanceRule           = "RebalanceRule"
	InstanceStateChangeRule = "InstanceStateChangeRule"

	eventsService = "events.amazonaws.com"
	sqsService    = "sqs.amazonaws.com"

	awsHealth = "aws.health"
	awsEC2    = "aws.ec2"

	defaultMessageRetentionPeriod = 300
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
	managedPolicyNames := sets.New[string]()
	managedPolicyNames.Insert(iamPolicyAmazonEKSWorkerNodePolicy,
		iamPolicyAmazonEKSCNIPolicy,
		iamPolicyAmazonEC2ContainerRegistryReadOnly,
		iamPolicyAmazonSSMManagedInstanceCore,
	)
	k.Template().Mappings[servicePrincipalPartitionMapName] = api.Partitions.ServicePrincipalPartitionMappings()
	roleName := gfnt.NewString(fmt.Sprintf("eksctl-%s-%s", KarpenterNodeRoleName, k.clusterSpec.Metadata.Name))
	role := gfniam.Role{
		RoleName:                 roleName,
		Path:                     gfnt.NewString("/"),
		AssumeRolePolicyDocument: cft.MakeAssumeRolePolicyDocumentForServices(MakeServiceRef("EC2")),
		ManagedPolicyArns:        gfnt.NewSlice(makePolicyARNs(sets.List(managedPolicyNames)...)...),
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
	rolePolicyStatements := []cft.MapOfInterfaces{
		{
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
				ec2DescribeImages,
				ec2DescribeSpotPriceHistory,
				iamPassRole,
				iamCreateServiceLinkedRole,
				iamGetInstanceProfile,
				iamCreateInstanceProfile,
				iamDeleteInstanceProfile,
				iamTagInstanceProfile,
				iamAddRoleToInstanceProfile,
				ssmGetParameter,
				pricingGetProducts,
			},
		},
	}

	if api.IsEnabled(k.clusterSpec.Karpenter.WithSpotInterruptionQueue) {
		rolePolicyStatements = append(rolePolicyStatements, cft.MapOfInterfaces{
			"Effect":   effectAllow,
			"Resource": gfnt.MakeFnGetAtt(KarpenterInterruptionQueue, gfnt.NewString("Arn")),
			"Action": []string{
				sqsDeleteMessage,
				sqsGetQueueAttributes,
				sqsGetQueueURL,
				sqsReceiveMessage,
			},
		})
		k.addSpotInterruptionQueueWithRules()
	}

	managedPolicy := gfniam.ManagedPolicy{
		ManagedPolicyName: managedPolicyName,
		PolicyDocument:    cft.MakePolicyDocument(rolePolicyStatements...),
	}
	k.newResource(KarpenterManagedPolicy, &managedPolicy)

	return nil
}

func (k *KarpenterResourceSet) addSpotInterruptionQueueWithRules() {
	interruptionQueue := gfnsqs.Queue{
		QueueName:              gfnt.NewString(k.clusterSpec.Metadata.Name),
		MessageRetentionPeriod: gfnt.NewInteger(defaultMessageRetentionPeriod),
	}
	queueRef := k.newResource(KarpenterInterruptionQueue, &interruptionQueue)

	queuePolicyStatements := cft.MapOfInterfaces{
		"Effect": effectAllow,
		"Principal": cft.MapOfInterfaces{
			"Service": cft.SliceOfInterfaces{
				eventsService,
				sqsService,
			},
		},
		"Resource": gfnt.MakeFnGetAtt(KarpenterInterruptionQueue, gfnt.NewString("Arn")),
		"Action": []string{
			sqsSendMessage,
		},
	}
	interruptionQueuePolicy := gfnsqs.QueuePolicy{
		Queues:         gfnt.NewSlice(queueRef),
		PolicyDocument: cft.MakePolicyDocument(queuePolicyStatements),
	}
	k.newResource(KarpenterInterruptionQueuePolicy, &interruptionQueuePolicy)

	scheduledChangeRule := gfnevents.Rule{
		EventPattern: cft.MapOfInterfaces{
			"source":      gfnt.NewSlice(gfnt.NewString(awsHealth)),
			"detail-type": gfnt.NewSlice(gfnt.NewString("AWS Health Event")),
		},
		Targets: []gfnevents.Rule_Target{
			{
				Id:  gfnt.NewString(KarpenterInterruptionQueueTarget),
				Arn: gfnt.MakeFnGetAtt(KarpenterInterruptionQueue, gfnt.NewString("Arn")),
			},
		},
	}
	k.newResource(ScheduledChangeRule, &scheduledChangeRule)

	spotInterruptionRule := gfnevents.Rule{
		EventPattern: cft.MapOfInterfaces{
			"source":      gfnt.NewSlice(gfnt.NewString(awsEC2)),
			"detail-type": gfnt.NewSlice(gfnt.NewString("EC2 Spot Instance Interruption Warning")),
		},
		Targets: []gfnevents.Rule_Target{
			{
				Id:  gfnt.NewString(KarpenterInterruptionQueueTarget),
				Arn: gfnt.MakeFnGetAtt(KarpenterInterruptionQueue, gfnt.NewString("Arn")),
			},
		},
	}
	k.newResource(SpotInterruptionRule, &spotInterruptionRule)

	rebalanceRule := gfnevents.Rule{
		EventPattern: cft.MapOfInterfaces{
			"source":      gfnt.NewSlice(gfnt.NewString(awsEC2)),
			"detail-type": gfnt.NewSlice(gfnt.NewString("EC2 Instance Rebalance Recommendation")),
		},
		Targets: []gfnevents.Rule_Target{
			{
				Id:  gfnt.NewString(KarpenterInterruptionQueueTarget),
				Arn: gfnt.MakeFnGetAtt(KarpenterInterruptionQueue, gfnt.NewString("Arn")),
			},
		},
	}
	k.newResource(RebalanceRule, &rebalanceRule)

	instanceStateChangeRule := gfnevents.Rule{
		EventPattern: cft.MapOfInterfaces{
			"source":      gfnt.NewSlice(gfnt.NewString(awsEC2)),
			"detail-type": gfnt.NewSlice(gfnt.NewString("EC2 Instance State-change Notification")),
		},
		Targets: []gfnevents.Rule_Target{
			{
				Id:  gfnt.NewString(KarpenterInterruptionQueueTarget),
				Arn: gfnt.MakeFnGetAtt(KarpenterInterruptionQueue, gfnt.NewString("Arn")),
			},
		},
	}
	k.newResource(InstanceStateChangeRule, &instanceStateChangeRule)
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
