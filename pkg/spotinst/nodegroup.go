package spotinst

import (
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
)

// UseNodeGroupSpotinstOceanId retrieves the Spotinst Ocean cluster identifier
// from an existing node group based on stack outputs.
func UseNodeGroupSpotinstOceanId(provider api.ClusterProvider, stack *cfn.Stack, ng *api.NodeGroup) error {
	if ng.Spotinst.Ocean == nil {
		ng.Spotinst.Ocean = &api.NodeGroupSpotinstOcean{}
	}

	requiredCollectors := map[string]outputs.Collector{
		outputs.NodeGroupSpotinstOceanID: func(v string) error {
			ng.Spotinst.Ocean.ID = &v
			return nil
		},
	}

	return outputs.Collect(*stack, requiredCollectors, nil)
}
