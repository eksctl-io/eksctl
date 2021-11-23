package builder_test

import (
	. "github.com/onsi/ginkgo"
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
		cfg.Karpenter = &api.Karpenter{
			Version:               "0.4.3",
			AddDefaultProvisioner: api.Disabled(),
			CreateServiceAccount:  api.Disabled(),
		}
	})

	Context("AddAllResources", func() {
		It("generates the correct CloudFormation template", func() {
			krs := builder.NewKarpenterResourceSet(cfg)
			Expect(krs.AddAllResources()).To(Succeed())
			result, err := krs.RenderJSON()
			Expect(err).NotTo(HaveOccurred())
			Expect(string(result)).To(Equal(expectedTemplate))
		})
	})
})

var expectedTemplate = `{
  "AWSTemplateFormatVersion": "2010-09-09",
  "Description": "Karpenter Stack [created and managed by eksctl]",
  "Resources": {
    "KarpenterControllerPolicy": {
      "Type": "AWS::IAM::ManagedPolicy",
      "Properties": {
        "ManagedPolicyName": {
          "Fn::Sub": "KarpenterControllerPolicy-${AWS::StackName}"
        },
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
                "ec2:RunInstances",
                "ec2:TerminateInstances",
                "iam:PassRole",
                "ssm:GetParameter"
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
        "InstanceProfileName": {
          "Fn::Sub": "KarpenterNodeInstanceProfile-${AWS::StackName}"
        },
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
          "Action": [
            "sts:AssumeRole"
          ],
          "Effect": "Allow",
          "Principal": {
            "Service": [
              "ec2.amazonaws.com"
            ]
          }
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
        "RoleName": {
          "Fn::Sub": "KarpenterNodeRole-${AWS::StackName}"
        },
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
