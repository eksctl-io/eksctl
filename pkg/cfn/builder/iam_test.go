package builder_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	cft "github.com/weaveworks/eksctl/pkg/cfn/template"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"

	. "github.com/weaveworks/eksctl/pkg/cfn/template/matchers"
)

var _ = Describe("template builder for IAM", func() {
	Describe("IAMServiceAccount", func() {
		var (
			oidc *iamoidc.OpenIDConnectManager
			cfg  *api.ClusterConfig
			err  error
		)

		BeforeEach(func() {
			oidc, err = iamoidc.NewOpenIDConnectManager(nil, "456123987123", "https://oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E", "aws", nil)
			Expect(err).NotTo(HaveOccurred())

			oidc.ProviderARN = "arn:aws:iam::456123987123:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E"

			cfg = api.NewClusterConfig()

			cfg.IAM.WithOIDC = api.Enabled()
			cfg.IAM.ServiceAccounts = []*api.ClusterIAMServiceAccount{}
		})

		It("can construct an iamserviceaccount addon template with one managed policy", func() {
			serviceAccount := &api.ClusterIAMServiceAccount{}

			serviceAccount.Name = "sa-1"

			serviceAccount.AttachPolicyARNs = []string{"arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"}

			appendServiceAccountToClusterConfig(cfg, serviceAccount)

			rs := builder.NewIAMRoleResourceSetForServiceAccount(serviceAccount, oidc)

			templateBody := []byte{}

			Expect(rs).To(RenderWithoutErrors(&templateBody))

			t := cft.NewTemplate()

			Expect(t).To(LoadBytesWithoutErrors(templateBody))

			Expect(t.Description).To(Equal("IAM role for serviceaccount \"default/sa-1\" [created and managed by eksctl]"))

			Expect(t.Resources).To(HaveLen(1))
			Expect(t.Outputs).To(HaveLen(1))

			Expect(t).To(HaveResource(outputs.IAMServiceAccountRoleName, "AWS::IAM::Role"))

			Expect(t).To(HaveResourceWithPropertyValue(outputs.IAMServiceAccountRoleName, "AssumeRolePolicyDocument", expectedServiceAccountAssumeRolePolicyDocument))
			Expect(t).To(HaveResourceWithPropertyValue(outputs.IAMServiceAccountRoleName, "ManagedPolicyArns", `[
			"arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"
		]`))

			Expect(t).To(HaveOutputWithValue(outputs.IAMServiceAccountRoleName, `{ "Fn::GetAtt": "Role1.Arn" }`))
		})

		It("can construct an iamserviceaccount addon template with one inline policy", func() {
			serviceAccount := &api.ClusterIAMServiceAccount{}

			serviceAccount.Name = "sa-1"

			serviceAccount.AttachPolicy = cft.MakePolicyDocument(
				cft.MapOfInterfaces{
					"Effect": "Allow",
					"Action": []string{
						"s3:Get*",
					},
					"Resource": "*",
				},
			)

			appendServiceAccountToClusterConfig(cfg, serviceAccount)

			rs := builder.NewIAMRoleResourceSetForServiceAccount(serviceAccount, oidc)

			templateBody := []byte{}

			Expect(rs).To(RenderWithoutErrors(&templateBody))

			t := cft.NewTemplate()

			Expect(t).To(LoadBytesWithoutErrors(templateBody))

			Expect(t.Description).To(Equal("IAM role for serviceaccount \"default/sa-1\" [created and managed by eksctl]"))

			Expect(t.Resources).To(HaveLen(2))
			Expect(t.Outputs).To(HaveLen(1))

			Expect(t).To(HaveResource(outputs.IAMServiceAccountRoleName, "AWS::IAM::Role"))
			Expect(t).To(HaveResource("Policy1", "AWS::IAM::Policy"))

			Expect(t).NotTo(HaveResourceWithProperties(outputs.IAMServiceAccountRoleName, "ManagedPolicyArns"))

			Expect(t).To(HaveResourceWithPropertyValue(outputs.IAMServiceAccountRoleName, "AssumeRolePolicyDocument", expectedServiceAccountAssumeRolePolicyDocument))
			Expect(t).To(HaveResourceWithPropertyValue("Policy1", "PolicyName", `{ "Fn::Sub": "${AWS::StackName}-Policy1" }`))
			Expect(t).To(HaveResourceWithPropertyValue("Policy1", "PolicyDocument", `{
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Action": [
                        "s3:Get*"
                    ],
                    "Resource": "*"
                }
            ]
        }`))

			Expect(t).To(HaveOutputWithValue(outputs.IAMServiceAccountRoleName, `{ "Fn::GetAtt": "Role1.Arn" }`))
		})

		It("can construct an iamserviceaccount addon template with a custom role name", func() {
			serviceAccount := &api.ClusterIAMServiceAccount{}

			serviceAccount.Name = "sa-1"

			serviceAccount.RoleName = "custom-role-name"

			serviceAccount.AttachPolicy = cft.MakePolicyDocument(
				cft.MapOfInterfaces{
					"Effect": "Allow",
					"Action": []string{
						"s3:Get*",
					},
					"Resource": "*",
				},
			)

			rs := builder.NewIAMRoleResourceSetForServiceAccount(serviceAccount, oidc)

			templateBody := []byte{}

			Expect(rs).To(RenderWithoutErrors(&templateBody))

			t := cft.NewTemplate()

			Expect(t).To(LoadBytesWithoutErrors(templateBody))

			Expect(rs.WithNamedIAM()).To(Equal(true))
			Expect(t).To(HaveResource(outputs.IAMServiceAccountRoleName, "AWS::IAM::Role"))
			Expect(t).To(HaveResourceWithPropertyValue(outputs.IAMServiceAccountRoleName, "RoleName", `"custom-role-name"`))
		})

		It("can constuct an iamserviceaccount addon template with two managed policies and one inline policy", func() {
			serviceAccount := &api.ClusterIAMServiceAccount{}

			serviceAccount.Name = "sa-1"

			serviceAccount.AttachPolicyARNs = []string{
				"arn:aws:iam::aws:policy/AmazonDynamoDBReadOnlyAccess",
				"arn:aws:iam::aws:policy/AmazonElastiCacheFullAccess",
			}

			serviceAccount.AttachPolicy = cft.MakePolicyDocument(
				cft.MapOfInterfaces{
					"Effect": "Allow",
					"Action": []string{
						"s3:Get*",
					},
					"Resource": "*",
				},
			)

			appendServiceAccountToClusterConfig(cfg, serviceAccount)

			rs := builder.NewIAMRoleResourceSetForServiceAccount(serviceAccount, oidc)

			templateBody := []byte{}

			Expect(rs).To(RenderWithoutErrors(&templateBody))

			t := cft.NewTemplate()

			Expect(t).To(LoadBytesWithoutErrors(templateBody))

			Expect(t.Description).To(Equal("IAM role for serviceaccount \"default/sa-1\" [created and managed by eksctl]"))

			Expect(t.Resources).To(HaveLen(2))
			Expect(t.Outputs).To(HaveLen(1))

			Expect(t).To(HaveResource(outputs.IAMServiceAccountRoleName, "AWS::IAM::Role"))
			Expect(t).To(HaveResource("Policy1", "AWS::IAM::Policy"))

			Expect(t).To(HaveResourceWithPropertyValue(outputs.IAMServiceAccountRoleName, "ManagedPolicyArns", `[
			"arn:aws:iam::aws:policy/AmazonDynamoDBReadOnlyAccess",
			"arn:aws:iam::aws:policy/AmazonElastiCacheFullAccess"
		]`))

			Expect(t).To(HaveResourceWithPropertyValue(outputs.IAMServiceAccountRoleName, "AssumeRolePolicyDocument", expectedServiceAccountAssumeRolePolicyDocument))
			Expect(t).To(HaveResourceWithPropertyValue("Policy1", "PolicyName", `{ "Fn::Sub": "${AWS::StackName}-Policy1" }`))
			Expect(t).To(HaveResourceWithPropertyValue("Policy1", "PolicyDocument", `{
            "Version": "2012-10-17",
            "Statement": [
                {
                    "Effect": "Allow",
                    "Action": [
                        "s3:Get*"
                    ],
                    "Resource": "*"
                }
            ]
        }`))

			Expect(t).To(HaveOutputWithValue(outputs.IAMServiceAccountRoleName, `{ "Fn::GetAtt": "Role1.Arn" }`))
		})

		It("can construct an iamserviceaccount addon template with one managed policy and a permissions boundary", func() {
			serviceAccount := &api.ClusterIAMServiceAccount{}

			serviceAccount.Name = "sa-1"

			serviceAccount.AttachPolicyARNs = []string{"arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"}

			serviceAccount.PermissionsBoundary = "arn:aws:iam::aws:policy/policy/boundary"

			appendServiceAccountToClusterConfig(cfg, serviceAccount)

			rs := builder.NewIAMRoleResourceSetForServiceAccount(serviceAccount, oidc)

			templateBody := []byte{}

			Expect(rs).To(RenderWithoutErrors(&templateBody))

			t := cft.NewTemplate()

			Expect(t).To(LoadBytesWithoutErrors(templateBody))

			Expect(t.Description).To(Equal("IAM role for serviceaccount \"default/sa-1\" [created and managed by eksctl]"))

			Expect(t.Resources).To(HaveLen(1))
			Expect(t.Outputs).To(HaveLen(1))

			Expect(t).To(HaveResource(outputs.IAMServiceAccountRoleName, "AWS::IAM::Role"))

			Expect(t).To(HaveResourceWithPropertyValue(outputs.IAMServiceAccountRoleName, "AssumeRolePolicyDocument", expectedServiceAccountAssumeRolePolicyDocument))
			Expect(t).To(HaveResourceWithPropertyValue(outputs.IAMServiceAccountRoleName, "ManagedPolicyArns", `[
			"arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"
		]`))
			Expect(t).To(HaveResourceWithPropertyValue(outputs.IAMServiceAccountRoleName, "PermissionsBoundary", `"arn:aws:iam::aws:policy/policy/boundary"`))

			Expect(t).To(HaveOutputWithValue(outputs.IAMServiceAccountRoleName, `{ "Fn::GetAtt": "Role1.Arn" }`))
		})

		It("can construct an iamserviceaccount addon template with all the wellKnownPolicies", func() {
			serviceAccount := &api.ClusterIAMServiceAccount{}

			serviceAccount.Name = "sa-1"

			serviceAccount.WellKnownPolicies = api.WellKnownPolicies{
				ImageBuilder:              true,
				AutoScaler:                true,
				AWSLoadBalancerController: true,
				ExternalDNS:               true,
				CertManager:               true,
				EBSCSIController:          true,
				EFSCSIController:          true,
			}

			appendServiceAccountToClusterConfig(cfg, serviceAccount)

			rs := builder.NewIAMRoleResourceSetForServiceAccount(serviceAccount, oidc)

			templateBody := []byte{}

			Expect(rs).To(RenderWithoutErrors(&templateBody))

			t := cft.NewTemplate()

			Expect(t).To(LoadBytesWithoutErrors(templateBody))

			Expect(t.Description).To(Equal("IAM role for serviceaccount \"default/sa-1\" [created and managed by eksctl]"))

			Expect(t.Resources).To(HaveLen(9))
			Expect(t.Outputs).To(HaveLen(1))

			Expect(t).To(HaveResource(outputs.IAMServiceAccountRoleName, "AWS::IAM::Role"))

			Expect(t).To(HaveResourceWithPropertyValue(outputs.IAMServiceAccountRoleName, "AssumeRolePolicyDocument", expectedServiceAccountAssumeRolePolicyDocument))
			Expect(t).To(HaveResourceWithPropertyValue(outputs.IAMServiceAccountRoleName, "ManagedPolicyArns", `[
							{
                "Fn::Sub": "arn:${AWS::Partition}:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser"
							},
							{
								"Fn::Sub": "arn:${AWS::Partition}:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy"
		      		}
            ]`))
			Expect(t).To(HaveOutputWithValue(outputs.IAMServiceAccountRoleName, `{ "Fn::GetAtt": "Role1.Arn" }`))
			Expect(t).To(HaveResourceWithPropertyValue("PolicyAWSLoadBalancerController", "PolicyDocument", expectedAWSLoadBalancerControllerPolicyDocument))
		})

		It("can parse an iamserviceaccount addon template", func() {
			t := cft.NewTemplate()

			Expect(t).To(LoadFileWithoutErrors("../template/testdata/addon-example-1.json"))

			Expect(t.Description).To(Equal("IAM role for serviceaccount \"default/sa-1\" [created and managed by eksctl]"))

			Expect(t.Resources).To(HaveLen(1))
			Expect(t.Outputs).To(HaveLen(1))

			Expect(t).To(HaveResource(outputs.IAMServiceAccountRoleName, "AWS::IAM::Role"))

			Expect(t).To(HaveResourceWithPropertyValue(outputs.IAMServiceAccountRoleName, "ManagedPolicyArns", `[ "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess" ]`))

			Expect(t).To(HaveResourceWithPropertyValue(outputs.IAMServiceAccountRoleName, "AssumeRolePolicyDocument", expectedServiceAccountAssumeRolePolicyDocument))

			Expect(t).To(HaveOutputWithValue(outputs.IAMServiceAccountRoleName, `{ "Fn::GetAtt": "Role1.Arn" }`))
		})
	})

	Describe("IAMRole", func() {
		var (
			oidc *iamoidc.OpenIDConnectManager
			cfg  *api.ClusterConfig
			err  error
		)

		BeforeEach(func() {
			oidc, err = iamoidc.NewOpenIDConnectManager(nil, "456123987123", "https://oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E", "aws", nil)
			Expect(err).NotTo(HaveOccurred())

			oidc.ProviderARN = "arn:aws:iam::456123987123:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E"

			cfg = api.NewClusterConfig()

			cfg.IAM.WithOIDC = api.Enabled()
		})

		It("can construct an iamrole template with attachPolicyARNs", func() {
			arns := []string{"arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"}

			rs := builder.NewIAMRoleResourceSetWithAttachPolicyARNs("VPC-addon", "", "", "boundary-arn", arns, oidc)

			templateBody := []byte{}

			Expect(rs).To(RenderWithoutErrors(&templateBody))

			t := cft.NewTemplate()

			Expect(t).To(LoadBytesWithoutErrors(templateBody))

			Expect(t.Description).To(Equal("IAM role for \"VPC-addon\" [created and managed by eksctl]"))

			Expect(t.Resources).To(HaveLen(1))
			Expect(t.Outputs).To(HaveLen(1))

			Expect(t).To(HaveResource(outputs.IAMServiceAccountRoleName, "AWS::IAM::Role"))
			Expect(t).To(HaveResourceWithPropertyValue(outputs.IAMServiceAccountRoleName, "PermissionsBoundary", `"boundary-arn"`))

			Expect(t).To(HaveResourceWithPropertyValue(outputs.IAMServiceAccountRoleName, "AssumeRolePolicyDocument", expectedAssumeRolePolicyDocument))
			Expect(t).To(HaveResourceWithPropertyValue(outputs.IAMServiceAccountRoleName, "ManagedPolicyArns", `[
			"arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess"
		]`))

			Expect(t).To(HaveOutputWithValue(outputs.IAMServiceAccountRoleName, `{ "Fn::GetAtt": "Role1.Arn" }`))
		})

		It("can construct an iamrole template with attachPolicy", func() {
			attachPolicy := cft.MakePolicyDocument(
				cft.MapOfInterfaces{
					"Effect": "Allow",
					"Action": []string{
						"s3:Get*",
					},
					"Resource": "*",
				},
			)

			rs := builder.NewIAMRoleResourceSetWithAttachPolicy("VPC-addon", "", "", "boundary-arn", attachPolicy, oidc)

			templateBody := []byte{}

			Expect(rs).To(RenderWithoutErrors(&templateBody))

			t := cft.NewTemplate()

			Expect(t).To(LoadBytesWithoutErrors(templateBody))

			Expect(t.Description).To(Equal("IAM role for \"VPC-addon\" [created and managed by eksctl]"))

			Expect(t.Resources).To(HaveLen(2))
			Expect(t.Outputs).To(HaveLen(1))

			Expect(t).To(HaveResource(outputs.IAMServiceAccountRoleName, "AWS::IAM::Role"))
			Expect(t).To(HaveResourceWithPropertyValue(outputs.IAMServiceAccountRoleName, "PermissionsBoundary", `"boundary-arn"`))

			Expect(t).To(HaveResource("Policy1", "AWS::IAM::Policy"))

			Expect(t).NotTo(HaveResourceWithProperties(outputs.IAMServiceAccountRoleName, "ManagedPolicyArns"))

			Expect(t).To(HaveResourceWithPropertyValue(outputs.IAMServiceAccountRoleName, "AssumeRolePolicyDocument", expectedAssumeRolePolicyDocument))
			Expect(t).To(HaveResourceWithPropertyValue("Policy1", "PolicyName", `{ "Fn::Sub": "${AWS::StackName}-Policy1" }`))
			Expect(t).To(HaveResourceWithPropertyValue("Policy1", "PolicyDocument", `{
		   "Version": "2012-10-17",
		   "Statement": [
		       {
		           "Effect": "Allow",
		           "Action": [
		               "s3:Get*"
		           ],
		           "Resource": "*"
		       }
		   ]
		}`))

			Expect(t).To(HaveOutputWithValue(outputs.IAMServiceAccountRoleName, `{ "Fn::GetAtt": "Role1.Arn" }`))
		})
	})
})

func appendServiceAccountToClusterConfig(cfg *api.ClusterConfig, serviceAccount *api.ClusterIAMServiceAccount) {
	cfg.IAM.ServiceAccounts = append(cfg.IAM.ServiceAccounts, serviceAccount)

	api.SetClusterConfigDefaults(cfg)
	err := api.ValidateClusterConfig(cfg)
	Expect(err).NotTo(HaveOccurred())
}

const expectedServiceAccountAssumeRolePolicyDocument = `{
	"Statement": [
	  {
		"Action": [
		  "sts:AssumeRoleWithWebIdentity"
		],
		"Condition": {
		  "StringEquals": {
			"oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E:aud": "sts.amazonaws.com",
			"oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E:sub": "system:serviceaccount:default:sa-1"
		  }
		},
		"Effect": "Allow",
		"Principal": {
		  "Federated": "arn:aws:iam::456123987123:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E"
		}
	  }
	],
	"Version": "2012-10-17"
}`

const expectedAssumeRolePolicyDocument = `{
	"Statement": [
	  {
		"Action": [
		  "sts:AssumeRoleWithWebIdentity"
		],
		"Condition": {
		  "StringEquals": {
			"oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E:aud": "sts.amazonaws.com"
		  }
		},
		"Effect": "Allow",
		"Principal": {
		  "Federated": "arn:aws:iam::456123987123:oidc-provider/oidc.eks.us-west-2.amazonaws.com/id/A39A2842863C47208955D753DE205E6E"
		}
	  }
	],
	"Version": "2012-10-17"
}`

const expectedAWSLoadBalancerControllerPolicyDocument = `{
  "Statement": [
    {
      "Action": [
        "iam:CreateServiceLinkedRole"
      ],
      "Condition": {
        "StringEquals": {
          "iam:AWSServiceName": "elasticloadbalancing.amazonaws.com"
        }
      },
      "Effect": "Allow",
      "Resource": "*"
    },
    {
      "Action": [
        "ec2:DescribeAccountAttributes",
        "ec2:DescribeAddresses",
        "ec2:DescribeAvailabilityZones",
        "ec2:DescribeInternetGateways",
        "ec2:DescribeVpcs",
        "ec2:DescribeVpcPeeringConnections",
        "ec2:DescribeSubnets",
        "ec2:DescribeSecurityGroups",
        "ec2:DescribeInstances",
        "ec2:DescribeNetworkInterfaces",
        "ec2:DescribeTags",
        "ec2:GetCoipPoolUsage",
        "ec2:DescribeCoipPools",
        "ec2:GetSecurityGroupsForVpc",
        "ec2:DescribeIpamPools",
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
        "elasticloadbalancing:DescribeTrustStores",
        "elasticloadbalancing:DescribeListenerAttributes",
        "elasticloadbalancing:DescribeCapacityReservation"
      ],
      "Effect": "Allow",
      "Resource": "*"
    },
    {
      "Action": [
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
        "shield:DeleteProtection"
      ],
      "Effect": "Allow",
      "Resource": "*"
    },
    {
      "Action": [
        "ec2:AuthorizeSecurityGroupIngress",
        "ec2:RevokeSecurityGroupIngress"
      ],
      "Effect": "Allow",
      "Resource": "*"
    },
    {
      "Action": [
        "ec2:CreateSecurityGroup"
      ],
      "Effect": "Allow",
      "Resource": "*"
    },
    {
      "Action": [
        "ec2:CreateTags"
      ],
      "Condition": {
        "Null": {
          "aws:RequestTag/elbv2.k8s.aws/cluster": "false"
        },
        "StringEquals": {
          "ec2:CreateAction": "CreateSecurityGroup"
        }
      },
      "Effect": "Allow",
      "Resource": {
        "Fn::Sub": "arn:${AWS::Partition}:ec2:*:*:security-group/*"
      }
    },
    {
      "Action": [
        "ec2:CreateTags",
        "ec2:DeleteTags"
      ],
      "Condition": {
        "Null": {
          "aws:RequestTag/elbv2.k8s.aws/cluster": "true",
          "aws:ResourceTag/elbv2.k8s.aws/cluster": "false"
        }
      },
      "Effect": "Allow",
      "Resource": {
        "Fn::Sub": "arn:${AWS::Partition}:ec2:*:*:security-group/*"
      }
    },
    {
      "Action": [
        "ec2:AuthorizeSecurityGroupIngress",
        "ec2:RevokeSecurityGroupIngress",
        "ec2:DeleteSecurityGroup"
      ],
      "Condition": {
        "Null": {
          "aws:ResourceTag/elbv2.k8s.aws/cluster": "false"
        }
      },
      "Effect": "Allow",
      "Resource": "*"
    },
    {
      "Action": [
        "elasticloadbalancing:CreateLoadBalancer",
        "elasticloadbalancing:CreateTargetGroup"
      ],
      "Condition": {
        "Null": {
          "aws:RequestTag/elbv2.k8s.aws/cluster": "false"
        }
      },
      "Effect": "Allow",
      "Resource": "*"
    },
    {
      "Action": [
        "elasticloadbalancing:CreateListener",
        "elasticloadbalancing:DeleteListener",
        "elasticloadbalancing:CreateRule",
        "elasticloadbalancing:DeleteRule"
      ],
      "Effect": "Allow",
      "Resource": "*"
    },
    {
      "Action": [
        "elasticloadbalancing:AddTags",
        "elasticloadbalancing:RemoveTags"
      ],
      "Condition": {
        "Null": {
          "aws:RequestTag/elbv2.k8s.aws/cluster": "true",
          "aws:ResourceTag/elbv2.k8s.aws/cluster": "false"
        }
      },
      "Effect": "Allow",
      "Resource": [
        {
          "Fn::Sub": "arn:${AWS::Partition}:elasticloadbalancing:*:*:targetgroup/*/*"
        },
        {
          "Fn::Sub": "arn:${AWS::Partition}:elasticloadbalancing:*:*:loadbalancer/net/*/*"
        },
        {
          "Fn::Sub": "arn:${AWS::Partition}:elasticloadbalancing:*:*:loadbalancer/app/*/*"
        }
      ]
    },
    {
      "Action": [
        "elasticloadbalancing:AddTags",
        "elasticloadbalancing:RemoveTags"
      ],
      "Effect": "Allow",
      "Resource": [
        {
          "Fn::Sub": "arn:${AWS::Partition}:elasticloadbalancing:*:*:listener/net/*/*/*"
        },
        {
          "Fn::Sub": "arn:${AWS::Partition}:elasticloadbalancing:*:*:listener/app/*/*/*"
        },
        {
          "Fn::Sub": "arn:${AWS::Partition}:elasticloadbalancing:*:*:listener-rule/net/*/*/*"
        },
        {
          "Fn::Sub": "arn:${AWS::Partition}:elasticloadbalancing:*:*:listener-rule/app/*/*/*"
        }
      ]
    },
    {
      "Action": [
        "elasticloadbalancing:ModifyLoadBalancerAttributes",
        "elasticloadbalancing:SetIpAddressType",
        "elasticloadbalancing:SetSecurityGroups",
        "elasticloadbalancing:SetSubnets",
        "elasticloadbalancing:DeleteLoadBalancer",
        "elasticloadbalancing:ModifyTargetGroup",
        "elasticloadbalancing:ModifyTargetGroupAttributes",
        "elasticloadbalancing:DeleteTargetGroup",
        "elasticloadbalancing:ModifyListenerAttributes",
        "elasticloadbalancing:ModifyCapacityReservation",
        "elasticloadbalancing:ModifyIpPools"
      ],
      "Condition": {
        "Null": {
          "aws:ResourceTag/elbv2.k8s.aws/cluster": "false"
        }
      },
      "Effect": "Allow",
      "Resource": "*"
    },
    {
      "Action": [
        "elasticloadbalancing:AddTags"
      ],
      "Condition": {
        "Null": {
          "aws:RequestTag/elbv2.k8s.aws/cluster": "false"
        },
        "StringEquals": {
          "elasticloadbalancing:CreateAction": [
            "CreateTargetGroup",
            "CreateLoadBalancer"
          ]
        }
      },
      "Effect": "Allow",
      "Resource": [
        {
          "Fn::Sub": "arn:${AWS::Partition}:elasticloadbalancing:*:*:targetgroup/*/*"
        },
        {
          "Fn::Sub": "arn:${AWS::Partition}:elasticloadbalancing:*:*:loadbalancer/net/*/*"
        },
        {
          "Fn::Sub": "arn:${AWS::Partition}:elasticloadbalancing:*:*:loadbalancer/app/*/*"
        }
      ]
    },
    {
      "Action": [
        "elasticloadbalancing:RegisterTargets",
        "elasticloadbalancing:DeregisterTargets"
      ],
      "Effect": "Allow",
      "Resource": {
        "Fn::Sub": "arn:${AWS::Partition}:elasticloadbalancing:*:*:targetgroup/*/*"
      }
    },
    {
      "Action": [
        "elasticloadbalancing:SetWebAcl",
        "elasticloadbalancing:ModifyListener",
        "elasticloadbalancing:AddListenerCertificates",
        "elasticloadbalancing:RemoveListenerCertificates",
        "elasticloadbalancing:ModifyRule",
        "elasticloadbalancing:SetRulePriorities"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ],
  "Version": "2012-10-17"
}`
