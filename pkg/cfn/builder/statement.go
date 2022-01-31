package builder

import (
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	cft "github.com/weaveworks/eksctl/pkg/cfn/template"
)

const (
	effectAllow = "Allow"
	resourceAll = "*"
)

func loadBalancerControllerStatements() []cft.MapOfInterfaces {
	return []cft.MapOfInterfaces{
		{
			"Effect":   effectAllow,
			"Resource": addARNPartitionPrefix("ec2:*:*:security-group/*"),
			"Action":   []string{"ec2:CreateTags"},
			"Condition": map[string]interface{}{
				"StringEquals": map[string]string{
					"ec2:CreateAction": "CreateSecurityGroup",
				},
				"Null": map[string]string{
					"aws:RequestTag/elbv2.k8s.aws/cluster": "false",
				},
			},
		},
		{
			"Effect":   effectAllow,
			"Resource": addARNPartitionPrefix("ec2:*:*:security-group/*"),
			"Action": []string{
				"ec2:CreateTags",
				"ec2:DeleteTags",
			},
			"Condition": map[string]interface{}{
				"Null": map[string]string{
					"aws:RequestTag/elbv2.k8s.aws/cluster":  "true",
					"aws:ResourceTag/elbv2.k8s.aws/cluster": "false",
				},
			},
		},
		{
			"Effect":   effectAllow,
			"Resource": resourceAll,
			"Action": []string{
				"elasticloadbalancing:CreateLoadBalancer",
				"elasticloadbalancing:CreateTargetGroup",
			},
			"Condition": map[string]interface{}{
				"Null": map[string]string{
					"aws:RequestTag/elbv2.k8s.aws/cluster": "false",
				},
			},
		},
		{
			"Effect": effectAllow,
			"Resource": []*gfnt.Value{
				addARNPartitionPrefix("elasticloadbalancing:*:*:targetgroup/*/*"),
				addARNPartitionPrefix("elasticloadbalancing:*:*:loadbalancer/net/*/*"),
				addARNPartitionPrefix("elasticloadbalancing:*:*:loadbalancer/app/*/*"),
			},
			"Action": []string{
				"elasticloadbalancing:AddTags",
				"elasticloadbalancing:RemoveTags",
			},
			"Condition": map[string]interface{}{
				"Null": map[string]string{
					"aws:RequestTag/elbv2.k8s.aws/cluster":  "true",
					"aws:ResourceTag/elbv2.k8s.aws/cluster": "false",
				},
			},
		},
		{
			"Effect": effectAllow,
			"Resource": []*gfnt.Value{
				addARNPartitionPrefix("elasticloadbalancing:*:*:listener/net/*/*/*"),
				addARNPartitionPrefix("elasticloadbalancing:*:*:listener/app/*/*/*"),
				addARNPartitionPrefix("elasticloadbalancing:*:*:listener-rule/net/*/*/*"),
				addARNPartitionPrefix("elasticloadbalancing:*:*:listener-rule/app/*/*/*"),
			},
			"Action": []string{
				"elasticloadbalancing:AddTags",
				"elasticloadbalancing:RemoveTags",
			},
		},
		{
			"Effect":   effectAllow,
			"Resource": resourceAll,
			"Action": []string{
				"ec2:AuthorizeSecurityGroupIngress",
				"ec2:RevokeSecurityGroupIngress",
				"ec2:DeleteSecurityGroup",
				"elasticloadbalancing:ModifyLoadBalancerAttributes",
				"elasticloadbalancing:SetIpAddressType",
				"elasticloadbalancing:SetSecurityGroups",
				"elasticloadbalancing:SetSubnets",
				"elasticloadbalancing:DeleteLoadBalancer",
				"elasticloadbalancing:ModifyTargetGroup",
				"elasticloadbalancing:ModifyTargetGroupAttributes",
				"elasticloadbalancing:DeleteTargetGroup",
			},
			"Condition": map[string]interface{}{
				"Null": map[string]string{
					"aws:ResourceTag/elbv2.k8s.aws/cluster": "false",
				},
			},
		},
		{
			"Effect":   effectAllow,
			"Resource": addARNPartitionPrefix("elasticloadbalancing:*:*:targetgroup/*/*"),
			"Action": []string{
				"elasticloadbalancing:RegisterTargets",
				"elasticloadbalancing:DeregisterTargets",
			},
		},
		{
			"Effect":   effectAllow,
			"Resource": resourceAll,
			"Action": []string{
				"iam:CreateServiceLinkedRole",
				"ec2:DescribeAccountAttributes",
				"ec2:DescribeAddresses",
				"ec2:DescribeAvailabilityZones",
				"ec2:DescribeInternetGateways",
				"ec2:DescribeVpcs",
				"ec2:DescribeSubnets",
				"ec2:DescribeSecurityGroups",
				"ec2:DescribeInstances",
				"ec2:DescribeNetworkInterfaces",
				"ec2:DescribeTags",
				"ec2:DescribeVpcPeeringConnections",
				"elasticloadbalancing:DescribeLoadBalancers",
				"elasticloadbalancing:DescribeLoadBalancerAttributes",
				"elasticloadbalancing:DescribeListeners",
				"elasticloadbalancing:DescribeListenerCertificates",
				"elasticloadbalancing:DescribeSSLPolicies",
				"elasticloadbalancing:DescribeRules",
				"elasticloadbalancing:DescribeTargetGroups",
				"elasticloadbalancing:DescribeTargetGroupAttributes",
				"elasticloadbalancing:DescribeTargetHealth",
				"elasticloadbalancing:DescribeTags",
				"cognito-idp:DescribeUserPoolClient",
				"acm:ListCertificates",
				"acm:DescribeCertificate",
				"iam:ListServerCertificates",
				"iam:GetServerCertificate",
				"waf-regional:GetWebACL",
				"waf-regional:GetWebACLForResource",
				"waf-regional:AssociateWebACL",
				"waf-regional:DisassociateWebACL",
				"wafv2:GetWebACL",
				"wafv2:GetWebACLForResource",
				"wafv2:AssociateWebACL",
				"wafv2:DisassociateWebACL",
				"shield:GetSubscriptionState",
				"shield:DescribeProtection",
				"shield:CreateProtection",
				"shield:DeleteProtection",
				"ec2:AuthorizeSecurityGroupIngress",
				"ec2:RevokeSecurityGroupIngress",
				"ec2:CreateSecurityGroup",
				"elasticloadbalancing:CreateListener",
				"elasticloadbalancing:DeleteListener",
				"elasticloadbalancing:CreateRule",
				"elasticloadbalancing:DeleteRule",
				"elasticloadbalancing:SetWebAcl",
				"elasticloadbalancing:ModifyListener",
				"elasticloadbalancing:AddListenerCertificates",
				"elasticloadbalancing:RemoveListenerCertificates",
				"elasticloadbalancing:ModifyRule",
			},
		},
	}
}

func elbStatements() []cft.MapOfInterfaces {
	return []cft.MapOfInterfaces{
		{
			"Effect":   effectAllow,
			"Resource": resourceAll,
			"Action": []string{
				"ec2:DescribeAccountAttributes",
				"ec2:DescribeAddresses",
				"ec2:DescribeInternetGateways",
			},
		},
	}
}

func cloudWatchMetricsStatements() []cft.MapOfInterfaces {
	return []cft.MapOfInterfaces{
		{
			"Effect":   effectAllow,
			"Resource": resourceAll,
			"Action": []string{
				"cloudwatch:PutMetricData",
			},
		},
	}
}

func certManagerHostedZonesStatements() []cft.MapOfInterfaces {
	return []cft.MapOfInterfaces{
		{
			"Effect":   effectAllow,
			"Resource": resourceAll,
			"Action": []string{
				"route53:ListResourceRecordSets",
				"route53:ListHostedZonesByName",
			},
		},
	}
}

func certManagerGetChangeStatements() []cft.MapOfInterfaces {
	return []cft.MapOfInterfaces{
		{
			"Effect":   effectAllow,
			"Resource": addARNPartitionPrefix("route53:::change/*"),
			"Action": []string{
				"route53:GetChange",
			},
		},
	}
}

func changeSetStatements() []cft.MapOfInterfaces {
	return []cft.MapOfInterfaces{
		{
			"Effect":   effectAllow,
			"Resource": addARNPartitionPrefix("route53:::hostedzone/*"),
			"Action": []string{
				"route53:ChangeResourceRecordSets",
			},
		},
	}
}

func externalDNSHostedZonesStatements() []cft.MapOfInterfaces {
	return []cft.MapOfInterfaces{
		{
			"Effect":   effectAllow,
			"Resource": resourceAll,
			"Action": []string{
				"route53:ListHostedZones",
				"route53:ListResourceRecordSets",
				"route53:ListTagsForResource",
			},
		},
	}
}

func autoScalerStatements() []cft.MapOfInterfaces {
	return []cft.MapOfInterfaces{
		{
			"Effect":   effectAllow,
			"Resource": resourceAll,
			"Action": []string{
				"autoscaling:DescribeAutoScalingGroups",
				"autoscaling:DescribeAutoScalingInstances",
				"autoscaling:DescribeLaunchConfigurations",
				"autoscaling:DescribeTags",
				"autoscaling:SetDesiredCapacity",
				"autoscaling:TerminateInstanceInAutoScalingGroup",
				"ec2:DescribeInstanceTypes",
				"ec2:DescribeLaunchTemplateVersions",
			},
		},
	}
}

func appMeshStatements(appendAction string) []cft.MapOfInterfaces {
	return []cft.MapOfInterfaces{
		{
			"Effect":   effectAllow,
			"Resource": resourceAll,
			"Action": []string{
				"servicediscovery:CreateService",
				"servicediscovery:DeleteService",
				"servicediscovery:GetService",
				"servicediscovery:GetInstance",
				"servicediscovery:RegisterInstance",
				"servicediscovery:DeregisterInstance",
				"servicediscovery:ListInstances",
				"servicediscovery:ListNamespaces",
				"servicediscovery:ListServices",
				"servicediscovery:GetInstancesHealthStatus",
				"servicediscovery:UpdateInstanceCustomHealthStatus",
				"servicediscovery:GetOperation",
				"route53:GetHealthCheck",
				"route53:CreateHealthCheck",
				"route53:UpdateHealthCheck",
				"route53:ChangeResourceRecordSets",
				"route53:DeleteHealthCheck",
				appendAction,
			},
		},
	}
}

func ebsStatements() []cft.MapOfInterfaces {
	return []cft.MapOfInterfaces{
		{
			"Effect": "Allow",
			"Action": []string{
				"ec2:CreateSnapshot",
				"ec2:AttachVolume",
				"ec2:DetachVolume",
				"ec2:ModifyVolume",
				"ec2:DescribeAvailabilityZones",
				"ec2:DescribeInstances",
				"ec2:DescribeSnapshots",
				"ec2:DescribeTags",
				"ec2:DescribeVolumes",
				"ec2:DescribeVolumesModifications",
			},
			"Resource": "*",
		},
		{
			"Effect": "Allow",
			"Action": []string{
				"ec2:CreateTags",
			},
			"Resource": []*gfnt.Value{
				addARNPartitionPrefix("ec2:*:*:volume/*"),
				addARNPartitionPrefix("ec2:*:*:snapshot/*"),
			},
			"Condition": cft.MapOfInterfaces{
				"StringEquals": cft.MapOfInterfaces{
					"ec2:CreateAction": []string{
						"CreateVolume",
						"CreateSnapshot",
					},
				},
			},
		},
		{
			"Effect": "Allow",
			"Action": []string{
				"ec2:DeleteTags",
			},
			"Resource": []*gfnt.Value{
				addARNPartitionPrefix("ec2:*:*:volume/*"),
				addARNPartitionPrefix("ec2:*:*:snapshot/*"),
			},
		},
		{
			"Effect": "Allow",

			"Action": []string{

				"ec2:CreateVolume",
			},
			"Resource": "*",
			"Condition": cft.MapOfInterfaces{
				"StringLike": cft.MapOfInterfaces{
					"aws:RequestTag/ebs.csi.aws.com/cluster": "true",
				},
			},
		},
		{
			"Effect": "Allow",
			"Action": []string{
				"ec2:CreateVolume",
			},
			"Resource": "*",
			"Condition": cft.MapOfInterfaces{
				"StringLike": cft.MapOfInterfaces{
					"aws:RequestTag/CSIVolumeName": "*",
				},
			},
		},
		{
			"Effect": "Allow",
			"Action": []string{
				"ec2:CreateVolume",
			},
			"Resource": "*",
			"Condition": cft.MapOfInterfaces{
				"StringLike": cft.MapOfInterfaces{
					"aws:RequestTag/kubernetes.io/cluster/*": "owned",
				},
			},
		},
		{
			"Effect": "Allow",
			"Action": []string{
				"ec2:DeleteVolume",
			},
			"Resource": "*",
			"Condition": cft.MapOfInterfaces{
				"StringLike": cft.MapOfInterfaces{
					"ec2:ResourceTag/ebs.csi.aws.com/cluster": "true",
				},
			},
		},
		{
			"Effect": "Allow",

			"Action": []string{
				"ec2:DeleteVolume",
			},
			"Resource": "*",
			"Condition": cft.MapOfInterfaces{
				"StringLike": cft.MapOfInterfaces{
					"ec2:ResourceTag/CSIVolumeName": "*",
				},
			},
		},
		{
			"Effect": "Allow",
			"Action": []string{
				"ec2:DeleteVolume",
			},
			"Resource": "*",
			"Condition": cft.MapOfInterfaces{
				"StringLike": cft.MapOfInterfaces{
					"ec2:ResourceTag/kubernetes.io/cluster/*": "owned",
				},
			},
		},
		{
			"Effect": "Allow",
			"Action": []string{
				"ec2:DeleteSnapshot",
			},
			"Resource": "*",
			"Condition": cft.MapOfInterfaces{
				"StringLike": cft.MapOfInterfaces{
					"ec2:ResourceTag/CSIVolumeSnapshotName": "*",
				},
			},
		},
		{
			"Effect": "Allow",
			"Action": []string{
				"ec2:DeleteSnapshot",
			},
			"Resource": "*",
			"Condition": cft.MapOfInterfaces{
				"StringLike": cft.MapOfInterfaces{
					"ec2:ResourceTag/ebs.csi.aws.com/cluster": "true",
				},
			},
		},
	}
}

func serviceLinkRoleStatements() []cft.MapOfInterfaces {
	return []cft.MapOfInterfaces{
		{
			"Effect":   effectAllow,
			"Resource": addARNPartitionPrefix("iam::*:role/aws-service-role/*"),
			"Action": []string{
				"iam:CreateServiceLinkedRole",
				"iam:AttachRolePolicy",
				"iam:PutRolePolicy",
			},
		},
	}
}

func fsxStatements() []cft.MapOfInterfaces {
	return []cft.MapOfInterfaces{
		{
			"Effect":   effectAllow,
			"Resource": resourceAll,
			"Action": []string{
				"fsx:*",
			},
		},
	}
}

func xRayStatements() []cft.MapOfInterfaces {
	return []cft.MapOfInterfaces{
		{
			"Effect":   effectAllow,
			"Resource": resourceAll,
			"Action": []string{
				"xray:PutTraceSegments",
				"xray:PutTelemetryRecords",
				"xray:GetSamplingRules",
				"xray:GetSamplingTargets",
				"xray:GetSamplingStatisticSummaries",
			},
		},
	}
}

func efsStatements() []cft.MapOfInterfaces {
	return []cft.MapOfInterfaces{
		{
			"Effect":   effectAllow,
			"Resource": resourceAll,
			"Action": []string{
				"elasticfilesystem:*",
			},
		},
	}
}

func efsEc2Statements() []cft.MapOfInterfaces {
	return []cft.MapOfInterfaces{
		{
			"Effect":   effectAllow,
			"Resource": resourceAll,
			"Action": []string{
				"ec2:DescribeSubnets",
				"ec2:CreateNetworkInterface",
				"ec2:DescribeNetworkInterfaces",
				"ec2:DeleteNetworkInterface",
				"ec2:ModifyNetworkInterfaceAttribute",
				"ec2:DescribeNetworkInterfaceAttribute",
			},
		},
	}
}

func efsCSIControllerStatements() []cft.MapOfInterfaces {
	return []cft.MapOfInterfaces{
		{
			"Effect":   effectAllow,
			"Resource": resourceAll,
			"Action": []string{
				"elasticfilesystem:DescribeAccessPoints",
				"elasticfilesystem:DescribeFileSystems",
			},
		},
		{
			"Effect":   effectAllow,
			"Resource": resourceAll,
			"Action":   []string{"elasticfilesystem:CreateAccessPoint"},
			"Condition": map[string]interface{}{
				"StringLike": map[string]string{
					"aws:RequestTag/efs.csi.aws.com/cluster": "true",
				},
			},
		},
		{
			"Effect":   effectAllow,
			"Resource": resourceAll,
			"Action":   []string{"elasticfilesystem:DeleteAccessPoint"},
			"Condition": map[string]interface{}{
				"StringLike": map[string]string{
					"aws:ResourceTag/efs.csi.aws.com/cluster": "true",
				},
			},
		},
	}
}
