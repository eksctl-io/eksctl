package builder

import (
	"fmt"

	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

var servicePrincipalPartitionMappings = map[string]map[string]string{
	"aws": {
		"EC2":            "ec2.amazonaws.com",
		"EKS":            "eks.amazonaws.com",
		"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
	},
	"aws-us-gov": {
		"EC2":            "ec2.amazonaws.com",
		"EKS":            "eks.amazonaws.com",
		"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
	},
	"aws-cn": {
		"EC2":            "ec2.amazonaws.com.cn",
		"EKS":            "eks.amazonaws.com",
		"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
	},
}

const servicePrincipalPartitionMapName = "ServicePrincipalPartitionMap"

// MakeServiceRef returns a reference to an intrinsic map function that looks up the servicePrincipalName
// in servicePrincipalPartitionMappings
func MakeServiceRef(servicePrincipalName string) *gfnt.Value {
	return gfnt.MakeFnFindInMap(
		gfnt.NewString(servicePrincipalPartitionMapName), gfnt.RefPartition, gfnt.NewString(servicePrincipalName),
	)
}

func makePolicyARN(policyName string) *gfnt.Value {
	return gfnt.MakeFnSubString(fmt.Sprintf("arn:${%s}:iam::aws:policy/%s", gfnt.Partition, policyName))
}

func makePolicyARNs(policyNames ...string) []*gfnt.Value {
	policyARNs := make([]*gfnt.Value, len(policyNames))
	for i, policy := range policyNames {
		policyARNs[i] = makePolicyARN(policy)
	}
	return policyARNs
}

func addARNPartitionPrefix(s string) *gfnt.Value {
	return gfnt.MakeFnSubString(fmt.Sprintf("arn:${%s}:%s", gfnt.Partition, s))
}
