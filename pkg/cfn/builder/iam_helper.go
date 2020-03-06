package builder

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/arn"
	gfn "github.com/awslabs/goformation/cloudformation"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	cft "github.com/weaveworks/eksctl/pkg/cfn/template"
	"k8s.io/apimachinery/pkg/util/sets"
)

type cfnTemplate interface {
	attachAllowPolicy(name string, refRole *gfn.Value, resources interface{}, actions []string)
	newResource(name string, resource interface{}) *gfn.Value
}

// createRole creates an IAM role with policies required for the worker nodes and addons
func createRole(cfnTemplate cfnTemplate, iamConfig *api.NodeGroupIAM, managed bool) error {
	managedPolicyARNs, err := makeManagedPolicies(iamConfig, managed)
	if err != nil {
		return err
	}
	role := gfn.AWSIAMRole{
		Path:                     gfn.NewString("/"),
		AssumeRolePolicyDocument: cft.MakeAssumeRolePolicyDocumentForServices(MakeServiceRef("EC2")),
		ManagedPolicyArns:        managedPolicyARNs,
	}

	if iamConfig.InstanceRoleName != "" {
		role.RoleName = gfn.NewString(iamConfig.InstanceRoleName)
	}

	if iamConfig.InstanceRolePermissionsBoundary != "" {
		role.PermissionsBoundary = gfn.NewString(iamConfig.InstanceRolePermissionsBoundary)
	}

	refIR := cfnTemplate.newResource(cfnIAMInstanceRoleName, &role)

	if api.IsEnabled(iamConfig.WithAddonPolicies.AutoScaler) {
		cfnTemplate.attachAllowPolicy("PolicyAutoScaling", refIR, "*",
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
		cfnTemplate.attachAllowPolicy("PolicyCertManagerChangeSet", refIR, addARNPartitionPrefix("route53:::hostedzone/*"),
			[]string{
				"route53:ChangeResourceRecordSets",
			},
		)
		cfnTemplate.attachAllowPolicy("PolicyCertManagerHostedZones", refIR, "*",
			[]string{
				"route53:ListHostedZones",
				"route53:ListResourceRecordSets",
				"route53:ListHostedZonesByName",
			},
		)
		cfnTemplate.attachAllowPolicy("PolicyCertManagerGetChange", refIR, addARNPartitionPrefix("route53:::change/*"),
			[]string{
				"route53:GetChange",
			},
		)
	} else if api.IsEnabled(iamConfig.WithAddonPolicies.ExternalDNS) {
		cfnTemplate.attachAllowPolicy("PolicyExternalDNSChangeSet", refIR, addARNPartitionPrefix("route53:::hostedzone/*"),
			[]string{
				"route53:ChangeResourceRecordSets",
			},
		)
		cfnTemplate.attachAllowPolicy("PolicyExternalDNSHostedZones", refIR, "*",
			[]string{
				"route53:ListHostedZones",
				"route53:ListResourceRecordSets",
			},
		)
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.AppMesh) {
		cfnTemplate.attachAllowPolicy("PolicyAppMesh", refIR, "*",
			[]string{
				"appmesh:*",
				"servicediscovery:CreateService",
				"servicediscovery:GetService",
				"servicediscovery:RegisterInstance",
				"servicediscovery:DeregisterInstance",
				"servicediscovery:ListInstances",
				"servicediscovery:ListNamespaces",
				"servicediscovery:ListServices",
				"route53:GetHealthCheck",
				"route53:CreateHealthCheck",
				"route53:UpdateHealthCheck",
				"route53:ChangeResourceRecordSets",
				"route53:DeleteHealthCheck",
			},
		)
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.EBS) {
		cfnTemplate.attachAllowPolicy("PolicyEBS", refIR, "*",
			[]string{
				"ec2:AttachVolume",
				"ec2:CreateSnapshot",
				"ec2:CreateTags",
				"ec2:CreateVolume",
				"ec2:DeleteSnapshot",
				"ec2:DeleteTags",
				"ec2:DeleteVolume",
				"ec2:DescribeAvailabilityZones",
				"ec2:DescribeInstances",
				"ec2:DescribeSnapshots",
				"ec2:DescribeTags",
				"ec2:DescribeVolumes",
				"ec2:DescribeVolumesModifications",
				"ec2:DetachVolume",
				"ec2:ModifyVolume",
			},
		)
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.FSX) {
		cfnTemplate.attachAllowPolicy("PolicyFSX", refIR, "*",
			[]string{
				"fsx:*",
			},
		)
		cfnTemplate.attachAllowPolicy("PolicyServiceLinkRole", refIR, addARNPartitionPrefix("iam::*:role/aws-service-role/*"),
			[]string{
				"iam:CreateServiceLinkedRole",
				"iam:AttachRolePolicy",
				"iam:PutRolePolicy",
			},
		)
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.EFS) {
		cfnTemplate.attachAllowPolicy("PolicyEFS", refIR, "*",
			[]string{
				"elasticfilesystem:*",
			},
		)
		cfnTemplate.attachAllowPolicy("PolicyEFSEC2", refIR, "*",
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
		cfnTemplate.attachAllowPolicy("PolicyALBIngress", refIR, "*",
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
		cfnTemplate.attachAllowPolicy("PolicyXRay", refIR, "*",
			[]string{
				"xray:PutTraceSegments",
				"xray:PutTelemetryRecords",
				"xray:GetSamplingRules",
				"xray:GetSamplingTargets",
				"xray:GetSamplingStatisticSummaries",
			},
		)
	}
	return nil
}

func makeManagedPolicies(iamConfig *api.NodeGroupIAM, managed bool) ([]*gfn.Value, error) {
	managedPolicyNames := sets.NewString()
	if len(iamConfig.AttachPolicyARNs) == 0 {
		managedPolicyNames.Insert(iamDefaultNodePolicies...)
		if managed {
			// The Managed Nodegroup API requires this managed policy to be present, even though
			// AmazonEC2ContainerRegistryPowerUser (attached if imageBuilder is enabled) contains a superset of the
			// actions allowed by this managed policy
			managedPolicyNames.Insert(iamPolicyAmazonEC2ContainerRegistryReadOnly)
		}
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.ImageBuilder) {
		managedPolicyNames.Insert(iamPolicyAmazonEC2ContainerRegistryPowerUser)
	} else if !managed {
		// attach this policy even if `AttachPolicyARNs` is specified to preserve existing behaviour for unmanaged
		// nodegroups
		managedPolicyNames.Insert(iamPolicyAmazonEC2ContainerRegistryReadOnly)
	}

	if api.IsEnabled(iamConfig.WithAddonPolicies.CloudWatch) {
		managedPolicyNames.Insert(iamPolicyCloudWatchAgentServerPolicy)
	}

	for _, policyARN := range iamConfig.AttachPolicyARNs {
		parsedARN, err := arn.Parse(policyARN)
		if err != nil {
			return nil, err
		}
		start := strings.IndexRune(parsedARN.Resource, '/')
		if start == -1 || start+1 == len(parsedARN.Resource) {
			return nil, fmt.Errorf("failed to find ARN resource name: %s", parsedARN.Resource)
		}
		resourceName := parsedARN.Resource[start+1:]
		managedPolicyNames.Delete(resourceName)
	}

	return append(
		makeStringSlice(iamConfig.AttachPolicyARNs...),
		makePolicyARNs(managedPolicyNames.List()...)...,
	), nil
}
