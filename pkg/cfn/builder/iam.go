package builder

import (
	gfn "github.com/awslabs/goformation/cloudformation"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/iam"
)

const (
	iamPolicyAmazonEKSServicePolicyARN = "arn:aws:iam::aws:policy/AmazonEKSServicePolicy"
	iamPolicyAmazonEKSClusterPolicyARN = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"

	iamPolicyAmazonEKSWorkerNodePolicyARN           = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
	iamPolicyAmazonEKSCNIPolicyARN                  = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
	iamPolicyAmazonEC2ContainerRegistryPowerUserARN = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser"
	iamPolicyAmazonEC2ContainerRegistryReadOnlyARN  = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
)

var (
	iamDefaultNodePolicyARNs = []string{
		iamPolicyAmazonEKSWorkerNodePolicyARN,
		iamPolicyAmazonEKSCNIPolicyARN,
	}
)

func makePolicyDocument(statement map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []interface{}{
			statement,
		},
	}
}

func makeAssumeRolePolicyDocument(service string) map[string]interface{} {
	return makePolicyDocument(map[string]interface{}{
		"Effect": "Allow",
		"Principal": map[string][]string{
			"Service": []string{service},
		},
		"Action": []string{"sts:AssumeRole"},
	})
}

func (c *resourceSet) attachAllowPolicy(name string, refRole *gfn.Value, resources interface{}, actions []string) {
	c.newResource(name, &gfn.AWSIAMPolicy{
		PolicyName: makeName(name),
		Roles:      makeSlice(refRole),
		PolicyDocument: makePolicyDocument(map[string]interface{}{
			"Effect":   "Allow",
			"Resource": resources,
			"Action":   actions,
		}),
	})
}

// WithIAM states, if IAM roles will be created or not
func (c *ClusterResourceSet) WithIAM() bool {
	return c.rs.withIAM
}

// WithNamedIAM states, if specifically named IAM roles will be created or not
func (c *ClusterResourceSet) WithNamedIAM() bool {
	return c.rs.withNamedIAM
}

func (c *ClusterResourceSet) addResourcesForIAM() {
	c.rs.withNamedIAM = false

	if c.spec.IAM.ServiceRoleARN != "" {
		c.rs.withIAM = false
		c.rs.defineOutputWithoutCollector(outputs.ClusterServiceRoleARN, c.spec.IAM.ServiceRoleARN, true)
		return
	}

	c.rs.withIAM = true

	refSR := c.newResource("ServiceRole", &gfn.AWSIAMRole{
		AssumeRolePolicyDocument: makeAssumeRolePolicyDocument("eks.amazonaws.com"),
		ManagedPolicyArns: makeStringSlice(
			iamPolicyAmazonEKSServicePolicyARN,
			iamPolicyAmazonEKSClusterPolicyARN,
		),
	})
	c.rs.attachAllowPolicy("PolicyNLB", refSR, "*", []string{
		"elasticloadbalancing:*",
		"ec2:CreateSecurityGroup",
		"ec2:Describe*",
	})
	c.rs.attachAllowPolicy("PolicyCloudWatchMetrics", refSR, "*", []string{
		"cloudwatch:PutMetricData",
	})
	c.rs.defineOutputFromAtt(outputs.ClusterServiceRoleARN, "ServiceRole.Arn", true, func(v string) error {
		c.spec.IAM.ServiceRoleARN = v
		return nil
	})
}

// WithIAM states, if IAM roles will be created or not
func (n *NodeGroupResourceSet) WithIAM() bool {
	return n.rs.withIAM
}

// WithNamedIAM states, if specifically named IAM roles will be created or not
func (n *NodeGroupResourceSet) WithNamedIAM() bool {
	return n.rs.withNamedIAM
}

func (n *NodeGroupResourceSet) addResourcesForIAM() {
	if n.spec.IAM == nil {
		n.spec.IAM = &api.NodeGroupIAM{}
	}

	if n.spec.IAM.InstanceProfileARN != "" {
		n.rs.withIAM = false
		n.rs.withNamedIAM = false

		// if instance profile is given, as well as the role, we simply use both and export the role
		// (which is needed in order to authorise the nodegroup)
		n.instanceProfile = gfn.NewString(n.spec.IAM.InstanceProfileARN)
		if n.spec.IAM.InstanceRoleARN != "" {
			n.rs.defineOutputWithoutCollector(outputs.NodeGroupInstanceProfileARN, n.spec.IAM.InstanceProfileARN, true)
			n.rs.defineOutputWithoutCollector(outputs.NodeGroupInstanceRoleARN, n.spec.IAM.InstanceRoleARN, true)
			return
		}
		// if instance role is not given, export profile and use the getter to call importer function
		n.rs.defineOutput(outputs.NodeGroupInstanceProfileARN, n.spec.IAM.InstanceProfileARN, true, func(v string) error {
			return iam.ImportInstanceRoleFromProfileARN(n.provider, n.spec, v)
		})

		return
	}

	n.rs.withIAM = true

	if n.spec.IAM.InstanceRoleARN != "" {
		// if role is set, but profile isn't - create profile
		n.instanceProfile = n.newResource("NodeInstanceProfile", &gfn.AWSIAMInstanceProfile{
			Path:  gfn.NewString("/"),
			Roles: makeStringSlice(n.spec.IAM.InstanceRoleARN),
		})
		n.rs.defineOutputFromAtt(outputs.NodeGroupInstanceProfileARN, "NodeInstanceProfile.Arn", true, func(v string) error {
			n.spec.IAM.InstanceProfileARN = v
			return nil
		})
		n.rs.defineOutputWithoutCollector(outputs.NodeGroupInstanceRoleARN, n.spec.IAM.InstanceRoleARN, true)
		return
	}

	// if neither role nor profile are given - create both

	if n.spec.IAM.InstanceRoleName == "" {
		// setting role name requires additional capabilities
		n.rs.withNamedIAM = true
	}

	if len(n.spec.IAM.AttachPolicyARNs) == 0 {
		n.spec.IAM.AttachPolicyARNs = iamDefaultNodePolicyARNs
	}
	if v := n.spec.IAM.WithAddonPolicies.ImageBuilder; v != nil && *v {
		n.spec.IAM.AttachPolicyARNs = append(n.spec.IAM.AttachPolicyARNs, iamPolicyAmazonEC2ContainerRegistryPowerUserARN)
	} else {
		n.spec.IAM.AttachPolicyARNs = append(n.spec.IAM.AttachPolicyARNs, iamPolicyAmazonEC2ContainerRegistryReadOnlyARN)
	}

	role := gfn.AWSIAMRole{
		Path:                     gfn.NewString("/"),
		AssumeRolePolicyDocument: makeAssumeRolePolicyDocument("ec2.amazonaws.com"),
		ManagedPolicyArns:        makeStringSlice(n.spec.IAM.AttachPolicyARNs...),
	}

	if n.spec.IAM.InstanceRoleName != "" {
		role.RoleName = gfn.NewString(n.spec.IAM.InstanceRoleName)
	}

	refIR := n.newResource("NodeInstanceRole", &role)

	n.instanceProfile = n.newResource("NodeInstanceProfile", &gfn.AWSIAMInstanceProfile{
		Path:  gfn.NewString("/"),
		Roles: makeSlice(refIR),
	})

	if v := n.spec.IAM.WithAddonPolicies.AutoScaler; v != nil && *v {
		n.rs.attachAllowPolicy("PolicyAutoScaling", refIR, "*",
			[]string{
				"autoscaling:DescribeAutoScalingGroups",
				"autoscaling:DescribeAutoScalingInstances",
				"autoscaling:DescribeLaunchConfigurations",
				"autoscaling:DescribeTags",
				"autoscaling:SetDesiredCapacity",
				"autoscaling:TerminateInstanceInAutoScalingGroup",
			},
		)
	}

	if v := n.spec.IAM.WithAddonPolicies.ExternalDNS; v != nil && *v {
		n.rs.attachAllowPolicy("PolicyExternalDNS", refIR, "arn:aws:route53:::hostedzone/*",
			[]string{
				"route53:ChangeResourceRecordSets",
				"route53:ListHostedZones",
				"route53:ListResourceRecordSets",
			},
		)
	}

	n.rs.defineOutputFromAtt(outputs.NodeGroupInstanceProfileARN, "NodeInstanceProfile.Arn", true, func(v string) error {
		n.spec.IAM.InstanceProfileARN = v
		return nil
	})
	n.rs.defineOutputFromAtt(outputs.NodeGroupInstanceRoleARN, "NodeInstanceRole.Arn", true, func(v string) error {
		n.spec.IAM.InstanceRoleARN = v
		return nil
	})
}
