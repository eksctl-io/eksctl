package builder

import (
	"fmt"

	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

const servicePrincipalPartitionMapName = "ServicePrincipalPartitionMap"

// MakeServiceRef returns a reference to an intrinsic map function that looks up the servicePrincipalName
// in ServicePrincipalPartitionMap.
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
