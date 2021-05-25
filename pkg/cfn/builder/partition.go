package builder

import (
	"fmt"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

var servicePrincipalPartitionMappings = map[string]map[string]string{
	api.PartitionAWS: {
		"EC2":            "ec2.amazonaws.com",
		"EKS":            "eks.amazonaws.com",
		"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
	},
	api.PartitionUSGov: {
		"EC2":            "ec2.amazonaws.com",
		"EKS":            "eks.amazonaws.com",
		"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
	},
	api.PartitionChina: {
		"EC2":            "ec2.amazonaws.com.cn",
		"EKS":            "eks.amazonaws.com.cn",
		"EKSFargatePods": "eks-fargate-pods.amazonaws.com",
	},
	api.PartitionUSIso: {
		"EC2":            "ec2.c2s.ic.gov",
		"EKS":            "eks.c2s.ic.gov",
		"EKSFargatePods": "eks-fargate-pods.c2s.ic.gov",
	},
	api.PartitionUSIsob: {
		"EC2":            "ec2.sc2s.sgov.gov",
		"EKS":            "eks.sc2s.sgov.gov",
		"EKSFargatePods": "eks-fargate-pods.sc2s.sgov.gov",
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
