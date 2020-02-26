package builder

import (
	"fmt"

	gfn "github.com/awslabs/goformation/cloudformation"
)

var servicePrincipalPartitionMappings = map[string]map[string]string{
	"aws": {
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

func makeFnFindInMap(mapName string, args ...*gfn.Value) *gfn.Value {
	return gfn.MakeIntrinsic(gfn.FnFindInMap, append([]*gfn.Value{gfn.NewString(mapName)}, args...))
}

func makeServiceRef(servicePrincipalName string) *gfn.Value {
	return makeFnFindInMap(servicePrincipalPartitionMapName, gfn.RefPartition, gfn.NewString(servicePrincipalName))
}

func addARNPrefix(s string) *gfn.Value {
	return gfn.MakeFnSubString(fmt.Sprintf("arn:${%s}:%s", gfn.Partition, s))
}
