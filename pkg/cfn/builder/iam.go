package builder

import (
	gfn "github.com/awslabs/goformation/cloudformation"
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

func (c *ClusterResourceSet) addResourcesForIAM() {
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
}

// WithIAM states, if IAM roles will be created or not
func (n *NodeGroupResourceSet) WithIAM() bool {
	return n.rs.withIAM
}

func (n *NodeGroupResourceSet) addResourcesForIAM() {
	n.rs.withIAM = true

	if len(n.spec.PolicyARNs) == 0 {
		n.spec.PolicyARNs = iamDefaultNodePolicyARNs
	}
	if n.clusterSpec.Addons.WithIAM.PolicyAmazonEC2ContainerRegistryPowerUser {
		n.spec.PolicyARNs = append(n.spec.PolicyARNs, iamPolicyAmazonEC2ContainerRegistryPowerUserARN)
	} else {
		n.spec.PolicyARNs = append(n.spec.PolicyARNs, iamPolicyAmazonEC2ContainerRegistryReadOnlyARN)
	}

	refIR := n.newResource("NodeInstanceRole", &gfn.AWSIAMRole{
		Path:                     gfn.NewString("/"),
		AssumeRolePolicyDocument: makeAssumeRolePolicyDocument("ec2.amazonaws.com"),
		ManagedPolicyArns:        makeStringSlice(n.spec.PolicyARNs...),
	})

	n.instanceProfile = n.newResource("NodeInstanceProfile", &gfn.AWSIAMInstanceProfile{
		Path:  gfn.NewString("/"),
		Roles: makeSlice(refIR),
	})
	n.rs.attachAllowPolicy("PolicyTagDiscovery", refIR, "*", []string{
		"ec2:DescribeTags",
	})
	n.rs.attachAllowPolicy("PolicyStackSignal", refIR,
		map[string]interface{}{
			gfn.FnJoin: []interface{}{
				":",
				[]interface{}{
					"arn:aws:cloudformation",
					map[string]string{"Ref": gfn.Region},
					map[string]string{"Ref": gfn.AccountID},
					map[string]interface{}{
						gfn.FnJoin: []interface{}{
							"/",
							[]interface{}{
								"stack",
								map[string]string{"Ref": gfn.StackName},
								"*",
							},
						},
					},
				},
			},
		},
		[]string{
			"cloudformation:SignalResource",
		},
	)

	if n.clusterSpec.Addons.WithIAM.PolicyAutoScaling {
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

	if n.clusterSpec.Addons.WithIAM.PolicyExternalDNS {
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

	n.rs.newOutputFromAtt(cfnOutputInstanceRoleARN, "NodeInstanceRole.Arn", true)
}
