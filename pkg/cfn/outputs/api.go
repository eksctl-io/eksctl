package outputs

import (
	"fmt"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha3"
)

// Stack output names
const (
	// outputs from cluster stack
	ClusterVPC            = "VPC"
	ClusterSecurityGroup  = "SecurityGroup"
	ClusterSubnets        = "Subnets"
	ClusterSubnetsPrivate = string(ClusterSubnets + api.SubnetTopologyPrivate)
	ClusterSubnetsPublic  = string(ClusterSubnets + api.SubnetTopologyPublic)

	ClusterCertificateAuthorityData = "CertificateAuthorityData"
	ClusterEndpoint                 = "Endpoint"
	ClusterARN                      = "ARN"
	ClusterStackName                = "ClusterStackName"
	ClusterSharedNodeSecurityGroup  = "SharedNodeSecurityGroup"

	// outputs from nodegroup stack
	NodeGroupInstanceRoleARN = "InstanceRoleARN"
	// outputs to indicate configuration attributes that may have critical effect
	// on critical effect on forward-compatibility with respect to overal functionality
	// and integrity, e.g. networking
	NodeGroupFeaturePrivateNetworking   = "FeaturePrivateNetworking"
	NodeGroupFeatureSharedSecurityGroup = "FeatureSharedSecurityGroup"
)

// MustCollect will use each of the keys and attempt to find an output in the given
// stack, if any of the keys are not preset it will return an error
func MustCollect(stack cfn.Stack, keys []string, results map[string]string) error {
	for _, key := range keys {
		var value *string
		for _, x := range stack.Outputs {
			if *x.OutputKey == key {
				value = x.OutputValue
				break
			}
		}
		if value == nil {
			return fmt.Errorf("no ouput %q in stack %q", key, *stack.StackName)
		}
		results[key] = *value
	}
	return nil
}
