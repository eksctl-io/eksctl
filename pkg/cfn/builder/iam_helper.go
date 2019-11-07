package builder

import (
	gfn "github.com/awslabs/goformation/cloudformation"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	cft "github.com/weaveworks/eksctl/pkg/cfn/template"
)

type cfnTemplate interface {
	attachAllowPolicy(name string, refRole *gfn.Value, resources interface{}, actions []string)
	newResource(name string, resource interface{}) *gfn.Value
}

type IAMHelper struct {
	cfnTemplate cfnTemplate
	iamConfig   *api.NodeGroupIAM
}

func NewIAMHelper(cfnTemplate cfnTemplate, iamConfig *api.NodeGroupIAM) *IAMHelper {
	return &IAMHelper{cfnTemplate: cfnTemplate, iamConfig: iamConfig}
}

// CreateRole creates an IAM role with policies required for the worker nodes and addons
func (h *IAMHelper) CreateRole() {
	iamConfig := h.iamConfig
	attachPolicyARNs := iamConfig.AttachPolicyARNs
	if len(attachPolicyARNs) == 0 {
		attachPolicyARNs = iamDefaultNodePolicyARNs
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.ImageBuilder) {
		attachPolicyARNs = append(attachPolicyARNs, iamPolicyAmazonEC2ContainerRegistryPowerUserARN)
	} else {
		attachPolicyARNs = append(attachPolicyARNs, iamPolicyAmazonEC2ContainerRegistryReadOnlyARN)
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.CloudWatch) {
		attachPolicyARNs = append(attachPolicyARNs, iamPolicyCloudWatchAgentServerPolicyARN)
	}

	role := gfn.AWSIAMRole{
		Path: gfn.NewString("/"),
		AssumeRolePolicyDocument: cft.MakeAssumeRolePolicyDocumentForServices("ec2.amazonaws.com"),
		ManagedPolicyArns:        makeStringSlice(attachPolicyARNs...),
	}

	if iamConfig.InstanceRoleName != "" {
		role.RoleName = gfn.NewString(iamConfig.InstanceRoleName)
	}

	refIR := h.cfnTemplate.newResource(cfnIAMInstanceRoleName, &role)

	if api.IsEnabled(iamConfig.WithAddonPolicies.AutoScaler) {
		h.cfnTemplate.attachAllowPolicy("PolicyAutoScaling", refIR, "*",
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

	if api.IsEnabled(iamConfig.WithAddonPolicies.CertManager) {
		h.cfnTemplate.attachAllowPolicy("PolicyCertManagerChangeSet", refIR, "arn:aws:route53:::hostedzone/*",
			[]string{
				"route53:ChangeResourceRecordSets",
			},
		)
		h.cfnTemplate.attachAllowPolicy("PolicyCertManagerHostedZones", refIR, "*",
			[]string{
				"route53:ListHostedZones",
				"route53:ListResourceRecordSets",
				"route53:ListHostedZonesByName",
			},
		)
		h.cfnTemplate.attachAllowPolicy("PolicyCertManagerGetChange", refIR, "arn:aws:route53:::change/*",
			[]string{
				"route53:GetChange",
			},
		)
	} else if api.IsEnabled(iamConfig.WithAddonPolicies.ExternalDNS) {
		h.cfnTemplate.attachAllowPolicy("PolicyExternalDNSChangeSet", refIR, "arn:aws:route53:::hostedzone/*",
			[]string{
				"route53:ChangeResourceRecordSets",
			},
		)
		h.cfnTemplate.attachAllowPolicy("PolicyExternalDNSHostedZones", refIR, "*",
			[]string{
				"route53:ListHostedZones",
				"route53:ListResourceRecordSets",
			},
		)
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.AppMesh) {
		h.cfnTemplate.attachAllowPolicy("PolicyAppMesh", refIR, "*",
			[]string{
				"appmesh:*",
				"servicediscovery:CreateService",
				"servicediscovery:GetService",
				"servicediscovery:RegisterInstance",
				"servicediscovery:DeregisterInstance",
				"servicediscovery:ListInstances",
				"servicediscovery:ListNamespaces",
				"route53:GetHealthCheck",
				"route53:CreateHealthCheck",
				"route53:UpdateHealthCheck",
				"route53:ChangeResourceRecordSets",
				"route53:DeleteHealthCheck",
			},
		)
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.EBS) {
		h.cfnTemplate.attachAllowPolicy("PolicyEBS", refIR, "*",
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

	if api.IsEnabled(iamConfig.WithAddonPolicies.FSX) {
		h.cfnTemplate.attachAllowPolicy("PolicyFSX", refIR, "*",
			[]string{
				"fsx:*",
			},
		)
		h.cfnTemplate.attachAllowPolicy("PolicyServiceLinkRole", refIR, "arn:aws:iam::*:role/aws-service-role/*",
			[]string{
				"iam:CreateServiceLinkedRole",
				"iam:AttachRolePolicy",
				"iam:PutRolePolicy",
			},
		)
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.EFS) {
		h.cfnTemplate.attachAllowPolicy("PolicyEFS", refIR, "*",
			[]string{
				"elasticfilesystem:*",
			},
		)
		h.cfnTemplate.attachAllowPolicy("PolicyEFSEC2", refIR, "*",
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

	if api.IsEnabled(iamConfig.WithAddonPolicies.ALBIngress) {
		h.cfnTemplate.attachAllowPolicy("PolicyALBIngress", refIR, "*",
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

	if api.IsEnabled(iamConfig.WithAddonPolicies.XRay) {
		h.cfnTemplate.attachAllowPolicy("PolicyXRay", refIR, "*",
			[]string{
				"xray:PutTraceSegments",
				"xray:PutTelemetryRecords",
				"xray:GetSamplingRules",
				"xray:GetSamplingTargets",
				"xray:GetSamplingStatisticSummaries",
			},
		)
	}
}
