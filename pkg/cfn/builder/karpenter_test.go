package builder_test

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

var _ = Describe("karpenter stack", func() {
	var (
		cfg *api.ClusterConfig
	)

	BeforeEach(func() {
		cfg = api.NewClusterConfig()
		cfg.Metadata = &api.ClusterMeta{
			Name: "test-karpenter",
		}
		cfg.Karpenter = &api.Karpenter{
			Version:              "0.4.3",
			CreateServiceAccount: api.Disabled(),
		}
	})

	Context("AddAllResources", func() {
		It("generates the correct CloudFormation template", func() {
			krs := builder.NewKarpenterResourceSet(cfg, "eksctl-KarpenterNodeInstanceProfile-test-karpenter")
			Expect(krs.AddAllResources()).To(Succeed())
			result, err := krs.RenderJSON()
			Expect(err).NotTo(HaveOccurred())
			Expect(string(result)).To(Equal(fmt.Sprintf(expectedTemplate, "eksctl-KarpenterNodeInstanceProfile-test-karpenter")))
		})
		When("defaultInstanceProfile is set", func() {
			It("generates the correct custom CloudFormation template", func() {
				krs := builder.NewKarpenterResourceSet(cfg, "KarpenterNodeInstanceProfile-custom")
				Expect(krs.AddAllResources()).To(Succeed())
				result, err := krs.RenderJSON()
				Expect(err).NotTo(HaveOccurred())
				Expect(string(result)).To(Equal(fmt.Sprintf(expectedTemplate, "KarpenterNodeInstanceProfile-custom")))
			})
		})
		When("ServiceRolePermissionsBoundary is set", func() {
			It("it applies the given permission boundary", func() {
				cfg.IAM.ServiceRolePermissionsBoundary = aws.String("arn:aws:iam::xxx:policy/PB-Standard")
				krs := builder.NewKarpenterResourceSet(cfg, "KarpenterNodeInstanceProfile-custom")
				Expect(krs.AddAllResources()).To(Succeed())
				result, err := krs.RenderJSON()
				Expect(err).NotTo(HaveOccurred())
				Expect(string(result)).To(Equal(expectedTemplateWithPermissionBoundary))
			})
		})
	})
})

var expectedTemplate = `{
  "AWSTemplateFormatVersion": "2010-09-09",
  "Description": "Karpenter Stack [created and managed by eksctl]",
  "Mappings": {
    "ServicePrincipalPartitionMap": {
      "aws": {
        "EC2": "ec2.amazonaws.com",
        "EKS": "eks.amazonaws.com",
        "EKSFargatePods": "eks-fargate-pods.amazonaws.com"
      },
      "aws-cn": {
        "EC2": "ec2.amazonaws.com.cn",
        "EKS": "eks.amazonaws.com",
        "EKSFargatePods": "eks-fargate-pods.amazonaws.com"
      },
      "aws-us-gov": {
        "EC2": "ec2.amazonaws.com",
        "EKS": "eks.amazonaws.com",
        "EKSFargatePods": "eks-fargate-pods.amazonaws.com"
      }
    }
  },
  "Resources": {
    "KarpenterControllerPolicy": {
      "Type": "AWS::IAM::ManagedPolicy",
      "Properties": {
        "ManagedPolicyName": "eksctl-KarpenterControllerPolicy-test-karpenter",
        "PolicyDocument": {
          "Statement": [
            {
              "Action": [
                "ec2:CreateFleet",
                "ec2:CreateLaunchTemplate",
                "ec2:CreateTags",
                "ec2:DescribeAvailabilityZones",
                "ec2:DescribeInstanceTypeOfferings",
                "ec2:DescribeInstanceTypes",
                "ec2:DescribeInstances",
                "ec2:DescribeLaunchTemplates",
                "ec2:DescribeSecurityGroups",
                "ec2:DescribeSubnets",
                "ec2:DeleteLaunchTemplate",
                "ec2:RunInstances",
                "ec2:TerminateInstances",
                "ec2:DescribeImages",
                "ec2:DescribeSpotPriceHistory",
                "iam:PassRole",
                "iam:CreateServiceLinkedRole",
                "ssm:GetParameter",
                "pricing:GetProducts"
              ],
              "Effect": "Allow",
              "Resource": "*"
            }
          ],
          "Version": "2012-10-17"
        }
      }
    },
    "KarpenterNodeInstanceProfile": {
      "Type": "AWS::IAM::InstanceProfile",
      "Properties": {
        "InstanceProfileName": "%s",
        "Path": "/",
        "Roles": [
          {
            "Ref": "KarpenterNodeRole"
          }
        ]
      }
    },
    "KarpenterNodeRole": {
      "Type": "AWS::IAM::Role",
      "Properties": {
        "AssumeRolePolicyDocument": {
          "Statement": [
            {
              "Action": [
                "sts:AssumeRole"
              ],
              "Effect": "Allow",
              "Principal": {
                "Service": [
                  {
                    "Fn::FindInMap": [
                      "ServicePrincipalPartitionMap",
                      {
                        "Ref": "AWS::Partition"
                      },
                      "EC2"
                    ]
                  }
                ]
              }
            }
          ],
          "Version": "2012-10-17"
        },
        "ManagedPolicyArns": [
          {
            "Fn::Sub": "arn:${AWS::Partition}:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
          },
          {
            "Fn::Sub": "arn:${AWS::Partition}:iam::aws:policy/AmazonEKSWorkerNodePolicy"
          },
          {
            "Fn::Sub": "arn:${AWS::Partition}:iam::aws:policy/AmazonEKS_CNI_Policy"
          },
          {
            "Fn::Sub": "arn:${AWS::Partition}:iam::aws:policy/AmazonSSMManagedInstanceCore"
          }
        ],
        "Path": "/",
        "RoleName": "eksctl-KarpenterNodeRole-test-karpenter",
        "Tags": [
          {
            "Key": "Name",
            "Value": {
              "Fn::Sub": "${AWS::StackName}/KarpenterNodeRole"
            }
          }
        ]
      }
    }
  }
}`

var expectedTemplateWithPermissionBoundary = `{
  "AWSTemplateFormatVersion": "2010-09-09",
  "Description": "Karpenter Stack [created and managed by eksctl]",
  "Mappings": {
    "ServicePrincipalPartitionMap": {
      "aws": {
        "EC2": "ec2.amazonaws.com",
        "EKS": "eks.amazonaws.com",
        "EKSFargatePods": "eks-fargate-pods.amazonaws.com"
      },
      "aws-cn": {
        "EC2": "ec2.amazonaws.com.cn",
        "EKS": "eks.amazonaws.com",
        "EKSFargatePods": "eks-fargate-pods.amazonaws.com"
      },
      "aws-us-gov": {
        "EC2": "ec2.amazonaws.com",
        "EKS": "eks.amazonaws.com",
        "EKSFargatePods": "eks-fargate-pods.amazonaws.com"
      }
    }
  },
  "Resources": {
    "KarpenterControllerPolicy": {
      "Type": "AWS::IAM::ManagedPolicy",
      "Properties": {
        "ManagedPolicyName": "eksctl-KarpenterControllerPolicy-test-karpenter",
        "PolicyDocument": {
          "Statement": [
            {
              "Action": [
                "ec2:CreateFleet",
                "ec2:CreateLaunchTemplate",
                "ec2:CreateTags",
                "ec2:DescribeAvailabilityZones",
                "ec2:DescribeInstanceTypeOfferings",
                "ec2:DescribeInstanceTypes",
                "ec2:DescribeInstances",
                "ec2:DescribeLaunchTemplates",
                "ec2:DescribeSecurityGroups",
                "ec2:DescribeSubnets",
                "ec2:DeleteLaunchTemplate",
                "ec2:RunInstances",
                "ec2:TerminateInstances",
                "ec2:DescribeImages",
                "ec2:DescribeSpotPriceHistory",
                "iam:PassRole",
                "iam:CreateServiceLinkedRole",
                "ssm:GetParameter",
                "pricing:GetProducts"
              ],
              "Effect": "Allow",
              "Resource": "*"
            }
          ],
          "Version": "2012-10-17"
        }
      }
    },
    "KarpenterNodeInstanceProfile": {
      "Type": "AWS::IAM::InstanceProfile",
      "Properties": {
        "InstanceProfileName": "KarpenterNodeInstanceProfile-custom",
        "Path": "/",
        "Roles": [
          {
            "Ref": "KarpenterNodeRole"
          }
        ]
      }
    },
    "KarpenterNodeRole": {
      "Type": "AWS::IAM::Role",
      "Properties": {
        "AssumeRolePolicyDocument": {
          "Statement": [
            {
              "Action": [
                "sts:AssumeRole"
              ],
              "Effect": "Allow",
              "Principal": {
                "Service": [
                  {
                    "Fn::FindInMap": [
                      "ServicePrincipalPartitionMap",
                      {
                        "Ref": "AWS::Partition"
                      },
                      "EC2"
                    ]
                  }
                ]
              }
            }
          ],
          "Version": "2012-10-17"
        },
        "ManagedPolicyArns": [
          {
            "Fn::Sub": "arn:${AWS::Partition}:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
          },
          {
            "Fn::Sub": "arn:${AWS::Partition}:iam::aws:policy/AmazonEKSWorkerNodePolicy"
          },
          {
            "Fn::Sub": "arn:${AWS::Partition}:iam::aws:policy/AmazonEKS_CNI_Policy"
          },
          {
            "Fn::Sub": "arn:${AWS::Partition}:iam::aws:policy/AmazonSSMManagedInstanceCore"
          }
        ],
        "Path": "/",
        "PermissionsBoundary": "arn:aws:iam::xxx:policy/PB-Standard",
        "RoleName": "eksctl-KarpenterNodeRole-test-karpenter",
        "Tags": [
          {
            "Key": "Name",
            "Value": {
              "Fn::Sub": "${AWS::StackName}/KarpenterNodeRole"
            }
          }
        ]
      }
    }
  }
}`
