package builder

import (
	gfn "github.com/awslabs/goformation/cloudformation"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
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
	iamPolicyCloudWatchAgentServerPolicyARN         = "arn:aws:iam::aws:policy/CloudWatchAgentServerPolicy"
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
		n.instanceProfileARN = gfn.NewString(n.spec.IAM.InstanceProfileARN)
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
		n.newResource("NodeInstanceProfile", &gfn.AWSIAMInstanceProfile{
			Path:  gfn.NewString("/"),
			Roles: makeStringSlice(n.spec.IAM.InstanceRoleARN),
		})
		n.instanceProfileARN = gfn.MakeFnGetAttString("NodeInstanceProfile.Arn")
		n.rs.defineOutputFromAtt(outputs.NodeGroupInstanceProfileARN, "NodeInstanceProfile.Arn", true, func(v string) error {
			n.spec.IAM.InstanceProfileARN = v
			return nil
		})
		n.rs.defineOutputWithoutCollector(outputs.NodeGroupInstanceRoleARN, n.spec.IAM.InstanceRoleARN, true)
		return
	}

	// if neither role nor profile are given - create both

	if n.spec.IAM.InstanceRoleName != "" {
		// setting role name requires additional capabilities
		n.rs.withNamedIAM = true
	}

	if len(n.spec.IAM.AttachPolicyARNs) == 0 {
		n.spec.IAM.AttachPolicyARNs = iamDefaultNodePolicyARNs
	}
	if api.IsEnabled(n.spec.IAM.WithAddonPolicies.ImageBuilder) {
		n.spec.IAM.AttachPolicyARNs = append(n.spec.IAM.AttachPolicyARNs, iamPolicyAmazonEC2ContainerRegistryPowerUserARN)
	} else {
		n.spec.IAM.AttachPolicyARNs = append(n.spec.IAM.AttachPolicyARNs, iamPolicyAmazonEC2ContainerRegistryReadOnlyARN)
	}

	if api.IsEnabled(n.spec.IAM.WithAddonPolicies.CloudWatch) {
		n.spec.IAM.AttachPolicyARNs = append(n.spec.IAM.AttachPolicyARNs, iamPolicyCloudWatchAgentServerPolicyARN)
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

	n.newResource("NodeInstanceProfile", &gfn.AWSIAMInstanceProfile{
		Path:  gfn.NewString("/"),
		Roles: makeSlice(refIR),
	})
	n.instanceProfileARN = gfn.MakeFnGetAttString("NodeInstanceProfile.Arn")

	if api.IsEnabled(n.spec.IAM.WithAddonPolicies.AutoScaler) {
		n.rs.attachAllowPolicy("PolicyAutoScaling", refIR, "*",
			[]string{
				"autoscaling:DescribeAutoScalingGroups",
				"autoscaling:DescribeAutoScalingInstances",
				"autoscaling:DescribeLaunchConfigurations",
				"autoscaling:DescribeTags",
				"autoscaling:SetDesiredCapacity",
				"autoscaling:TerminateInstanceInAutoScalingGroup",
				"ec2:DescribeLaunchTemplateVersions",
			},
		)
	}

	if api.IsEnabled(n.spec.IAM.WithAddonPolicies.CertManager) {
		n.rs.attachAllowPolicy("PolicyCertManagerChangeSet", refIR, "arn:aws:route53:::hostedzone/*",
			[]string{
				"route53:ChangeResourceRecordSets",
			},
		)
		n.rs.attachAllowPolicy("PolicyCertManagerHostedZones", refIR, "*",
			[]string{
				"route53:ListHostedZones",
				"route53:ListResourceRecordSets",
				"route53:ListHostedZonesByName",
			},
		)
		n.rs.attachAllowPolicy("PolicyCertManagerGetChange", refIR, "arn:aws:route53:::change/*",
			[]string{
				"route53:GetChange",
			},
		)
	} else if api.IsEnabled(n.spec.IAM.WithAddonPolicies.ExternalDNS) {
		n.rs.attachAllowPolicy("PolicyExternalDNSChangeSet", refIR, "arn:aws:route53:::hostedzone/*",
			[]string{
				"route53:ChangeResourceRecordSets",
			},
		)
		n.rs.attachAllowPolicy("PolicyExternalDNSHostedZones", refIR, "*",
			[]string{
				"route53:ListHostedZones",
				"route53:ListResourceRecordSets",
			},
		)
	}


	if api.IsEnabled(n.spec.IAM.WithAddonPolicies.AppMesh) {
		n.rs.attachAllowPolicy("PolicyAppMesh", refIR, "*",
			[]string{
				"appmesh:*",
			},
		)
	}

	if api.IsEnabled(n.spec.IAM.WithAddonPolicies.EBS) {
		n.rs.attachAllowPolicy("PolicyEBS", refIR, "*",
			[]string{
				"ec2:AttachVolume",
				"ec2:CreateSnapshot",
				"ec2:CreateTags",
				"ec2:CreateVolume",
				"ec2:DeleteSnapshot",
				"ec2:DeleteTags",
				"ec2:DeleteVolume",
				"ec2:DescribeInstances",
				"ec2:DescribeSnapshots",
				"ec2:DescribeTags",
				"ec2:DescribeVolumes",
				"ec2:DetachVolume",
			},
		)
	}

	if api.IsEnabled(n.spec.IAM.WithAddonPolicies.FSX) {
		n.rs.attachAllowPolicy("PolicyFSX", refIR, "*",
			[]string{
				"fsx:*",
			},
		)
		n.rs.attachAllowPolicy("PolicyServiceLinkRole", refIR, "arn:aws:iam::*:role/aws-service-role/*",
			[]string{
				"iam:CreateServiceLinkedRole",
				"iam:AttachRolePolicy",
				"iam:PutRolePolicy",
			},
		)
	}

	if api.IsEnabled(n.spec.IAM.WithAddonPolicies.EFS) {
		n.rs.attachAllowPolicy("PolicyEFS", refIR, "*",
			[]string{
				"elasticfilesystem:*",
			},
		)
		n.rs.attachAllowPolicy("PolicyEFSEC2", refIR, "*",
			[]string{
				"ec2:DescribeSubnets",
				"ec2:CreateNetworkInterface",
				"ec2:DescribeNetworkInterfaces",
				"ec2:DeleteNetworkInterface",
				"ec2:ModifyNetworkInterfaceAttribute",
				"ec2:DescribeNetworkInterfaceAttribute",
			},
		)
	}

	if api.IsEnabled(n.spec.IAM.WithAddonPolicies.ALBIngress) {
		n.rs.attachAllowPolicy("PolicyALBIngress", refIR, "*",
			[]string{
				"acm:DescribeCertificate",
				"acm:ListCertificates",
				"acm:GetCertificate",
				"ec2:AuthorizeSecurityGroupIngress",
				"ec2:CreateSecurityGroup",
				"ec2:CreateTags",
				"ec2:DeleteTags",
				"ec2:DeleteSecurityGroup",
				"ec2:DescribeAccountAttributes",
				"ec2:DescribeAddresses",
				"ec2:DescribeInstances",
				"ec2:DescribeInstanceStatus",
				"ec2:DescribeInternetGateways",
				"ec2:DescribeNetworkInterfaces",
				"ec2:DescribeSecurityGroups",
				"ec2:DescribeSubnets",
				"ec2:DescribeTags",
				"ec2:DescribeVpcs",
				"ec2:ModifyInstanceAttribute",
				"ec2:ModifyNetworkInterfaceAttribute",
				"ec2:RevokeSecurityGroupIngress",
				"elasticloadbalancing:AddListenerCertificates",
				"elasticloadbalancing:AddTags",
				"elasticloadbalancing:CreateListener",
				"elasticloadbalancing:CreateLoadBalancer",
				"elasticloadbalancing:CreateRule",
				"elasticloadbalancing:CreateTargetGroup",
				"elasticloadbalancing:DeleteListener",
				"elasticloadbalancing:DeleteLoadBalancer",
				"elasticloadbalancing:DeleteRule",
				"elasticloadbalancing:DeleteTargetGroup",
				"elasticloadbalancing:DeregisterTargets",
				"elasticloadbalancing:DescribeListenerCertificates",
				"elasticloadbalancing:DescribeListeners",
				"elasticloadbalancing:DescribeLoadBalancers",
				"elasticloadbalancing:DescribeLoadBalancerAttributes",
				"elasticloadbalancing:DescribeRules",
				"elasticloadbalancing:DescribeSSLPolicies",
				"elasticloadbalancing:DescribeTags",
				"elasticloadbalancing:DescribeTargetGroups",
				"elasticloadbalancing:DescribeTargetGroupAttributes",
				"elasticloadbalancing:DescribeTargetHealth",
				"elasticloadbalancing:ModifyListener",
				"elasticloadbalancing:ModifyLoadBalancerAttributes",
				"elasticloadbalancing:ModifyRule",
				"elasticloadbalancing:ModifyTargetGroup",
				"elasticloadbalancing:ModifyTargetGroupAttributes",
				"elasticloadbalancing:RegisterTargets",
				"elasticloadbalancing:RemoveListenerCertificates",
				"elasticloadbalancing:RemoveTags",
				"elasticloadbalancing:SetIpAddressType",
				"elasticloadbalancing:SetSecurityGroups",
				"elasticloadbalancing:SetSubnets",
				"elasticloadbalancing:SetWebACL",
				"iam:CreateServiceLinkedRole",
				"iam:GetServerCertificate",
				"iam:ListServerCertificates",
				"waf-regional:GetWebACLForResource",
				"waf-regional:GetWebACL",
				"waf-regional:AssociateWebACL",
				"waf-regional:DisassociateWebACL",
				"tag:GetResources",
				"tag:TagResources",
				"waf:GetWebACL",
			},
		)
	}

	if api.IsEnabled(n.spec.IAM.WithAddonPolicies.XRay) {
		n.rs.attachAllowPolicy("PolicyXRay", refIR, "*",
			[]string{
				"xray:PutTraceSegments",
				"xray:PutTelemetryRecords",
				"xray:GetSamplingRules",
				"xray:GetSamplingTargets",
				"xray:GetSamplingStatisticSummaries",
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
