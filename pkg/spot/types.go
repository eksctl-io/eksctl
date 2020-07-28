package spot

import (
	"encoding/json"

	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

type (
	Resource struct {
		ResourceCredentials

		ServiceToken *gfnt.Value `json:"ServiceToken,omitempty"`
		FeatureFlags *gfnt.Value `json:"featureFlags,omitempty"`
	}

	ResourceCredentials struct {
		Account *gfnt.Value `json:"accountId,omitempty"`
		Token   *gfnt.Value `json:"accessToken,omitempty"`
	}

	NodeGroupResource struct {
		Resource

		OceanCluster    *NodeGroupCluster    `json:"ocean,omitempty"`
		OceanLaunchSpec *NodeGroupLaunchSpec `json:"oceanLaunchSpec,omitempty"`

		// for internal use only; used by `eksctl get nodegroup` command.
		OceanSummary *NodeGroupSummary `json:"oceanSummary,omitempty"`
	}

	NodeGroupCluster struct {
		Name       *string              `json:"name,omitempty"`
		ClusterID  *string              `json:"controllerClusterId,omitempty"`
		Region     *gfnt.Value          `json:"region,omitempty"`
		Capacity   *NodeGroupCapacity   `json:"capacity,omitempty"`
		Strategy   *NodeGroupStrategy   `json:"strategy,omitempty"`
		Compute    *NodeGroupCompute    `json:"compute,omitempty"`
		Scheduling *NodeGroupScheduling `json:"scheduling,omitempty"`
		AutoScaler *NodeGroupAutoScaler `json:"autoScaler,omitempty"`
	}

	NodeGroupLaunchSpec struct {
		Name                     *string                  `json:"name,omitempty"`
		OceanID                  *gfnt.Value              `json:"oceanId,omitempty"`
		ImageID                  *gfnt.Value              `json:"imageId,omitempty"`
		UserData                 *gfnt.Value              `json:"userData,omitempty"`
		KeyPair                  *gfnt.Value              `json:"keyPair,omitempty"`
		AssociatePublicIPAddress *gfnt.Value              `json:"associatePublicIpAddress,omitempty"`
		VolumeSize               *int                     `json:"rootVolumeSize,omitempty"`
		EBSOptimized             *bool                    `json:"ebsOptimized,omitempty"`
		SubnetIDs                interface{}              `json:"subnetIds,omitempty"`
		InstanceTypes            []string                 `json:"instanceTypes,omitempty"`
		IAMInstanceProfile       map[string]*gfnt.Value   `json:"iamInstanceProfile,omitempty"`
		SecurityGroupIDs         *gfnt.Value              `json:"securityGroupIds,omitempty"`
		Tags                     []*NodeGroupTag          `json:"tags,omitempty"`
		LoadBalancers            []*NodeGroupLoadBalancer `json:"loadBalancers,omitempty"`
		Labels                   []*NodeGroupLabel        `json:"labels,omitempty"`
		Taints                   []*NodeGroupTaint        `json:"taints,omitempty"`
		AutoScaler               *NodeGroupAutoScaler     `json:"autoScale,omitempty"`
	}

	NodeGroupSummary struct {
		ImageID  *gfnt.Value        `json:"imageId,omitempty"`
		Capacity *NodeGroupCapacity `json:"capacity,omitempty"`
	}

	NodeGroupCapacity struct {
		Minimum *int `json:"minimum,omitempty"`
		Maximum *int `json:"maximum,omitempty"`
		Target  *int `json:"target,omitempty"`
	}

	NodeGroupStrategy struct {
		SpotPercentage           *int  `json:"spotPercentage,omitempty"`
		UtilizeReservedInstances *bool `json:"utilizeReservedInstances,omitempty"`
		FallbackToOnDemand       *bool `json:"fallbackToOd,omitempty"`
		DrainingTimeout          *int  `json:"drainingTimeout,omitempty"`
	}

	NodeGroupCompute struct {
		SubnetIDs           interface{}             `json:"subnetIds,omitempty"`
		InstanceTypes       *NodeGroupInstanceTypes `json:"instanceTypes,omitempty"`
		LaunchSpecification *NodeGroupLaunchSpec    `json:"launchSpecification,omitempty"`
	}

	NodeGroupInstanceTypes struct {
		Whitelist []string `json:"whitelist,omitempty"`
		Blacklist []string `json:"blacklist,omitempty"`
	}

	NodeGroupLoadBalancer struct {
		Type *string `json:"type,omitempty"`
		Arn  *string `json:"arn,omitempty"`
		Name *string `json:"name,omitempty"`
	}

	NodeGroupTag struct {
		Key   interface{} `json:"tagKey,omitempty"`
		Value interface{} `json:"tagValue,omitempty"`
	}

	NodeGroupScheduling struct {
		ShutdownHours *NodeGroupSchedulingShutdownHours `json:"shutdownHours,omitempty"`
		Tasks         []*NodeGroupSchedulingTask        `json:"tasks,omitempty"`
	}

	NodeGroupSchedulingShutdownHours struct {
		IsEnabled   *bool    `json:"isEnabled,omitempty"`
		TimeWindows []string `json:"timeWindows,omitempty"`
	}

	NodeGroupSchedulingTask struct {
		IsEnabled      *bool   `json:"isEnabled,omitempty"`
		Type           *string `json:"taskType,omitempty"`
		CronExpression *string `json:"cronExpression,omitempty"`
	}

	NodeGroupAutoScaler struct {
		IsEnabled    *bool                          `json:"isEnabled,omitempty"`
		IsAutoConfig *bool                          `json:"isAutoConfig,omitempty"`
		Cooldown     *int                           `json:"cooldown,omitempty"`
		Headroom     *NodeGroupAutoScalerHeadroom   `json:"headroom,omitempty"`  // cluster
		Headrooms    []*NodeGroupAutoScalerHeadroom `json:"headrooms,omitempty"` // launchspec
	}

	NodeGroupAutoScalerHeadroom struct {
		CPUPerUnit    *int `json:"cpuPerUnit,omitempty"`
		GPUPerUnit    *int `json:"gpuPerUnit,omitempty"`
		MemoryPerUnit *int `json:"memoryPerUnit,omitempty"`
		NumOfUnits    *int `json:"numOfUnits,omitempty"`
	}

	NodeGroupLabel struct {
		Key   *string `json:"key,omitempty"`
		Value *string `json:"value,omitempty"`
	}

	NodeGroupTaint struct {
		Key    *string `json:"key,omitempty"`
		Value  *string `json:"value,omitempty"`
		Effect *string `json:"effect,omitempty"`
	}
)

// MarshalJSON implements the json.Marshaler interface.
func (x *NodeGroupResource) MarshalJSON() ([]byte, error) {
	var typ string
	if x.OceanCluster != nil {
		typ = "Custom::ocean"
	} else if x.OceanLaunchSpec != nil {
		typ = "Custom::oceanLaunchSpec"
	}
	type Properties NodeGroupResource
	return json.Marshal(&struct {
		Type       string
		Properties Properties
	}{
		Type:       typ,
		Properties: Properties(*x),
	})
}
