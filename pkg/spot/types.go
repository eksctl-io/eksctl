package spot

import (
	"encoding/json"

	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

type (
	Resource struct {
		ResourceCredentials

		ServiceToken *gfnt.Value        `json:"ServiceToken,omitempty"`
		FeatureFlags *gfnt.Value        `json:"featureFlags,omitempty"`
		Parameters   ResourceParameters `json:"parameters,omitempty"`
	}

	ResourceCredentials struct {
		Account *gfnt.Value `json:"accountId,omitempty"`
		Token   *gfnt.Value `json:"accessToken,omitempty"`
	}

	ResourceParameters struct {
		OnCreate map[string]interface{} `json:"create,omitempty"`
		OnUpdate map[string]interface{} `json:"update,omitempty"`
		OnDelete map[string]interface{} `json:"delete,omitempty"`
	}

	ResourceNodeGroup struct {
		Resource

		Cluster          *Cluster          `json:"ocean,omitempty"`
		VirtualNodeGroup *VirtualNodeGroup `json:"oceanLaunchSpec,omitempty"`
	}

	Cluster struct {
		Name       *string     `json:"name,omitempty"`
		ClusterID  *string     `json:"controllerClusterId,omitempty"`
		Region     *gfnt.Value `json:"region,omitempty"`
		Strategy   *Strategy   `json:"strategy,omitempty"`
		Compute    *Compute    `json:"compute,omitempty"`
		Scheduling *Scheduling `json:"scheduling,omitempty"`
		AutoScaler *AutoScaler `json:"autoScaler,omitempty"`
	}

	VirtualNodeGroup struct {
		Name                     *string                  `json:"name,omitempty"`
		OceanID                  *gfnt.Value              `json:"oceanId,omitempty"`
		ImageID                  *gfnt.Value              `json:"imageId,omitempty"`
		UserData                 *gfnt.Value              `json:"userData,omitempty"`
		KeyPair                  *gfnt.Value              `json:"keyPair,omitempty"`
		AssociatePublicIPAddress *gfnt.Value              `json:"associatePublicIpAddress,omitempty"`
		VolumeSize               *int                     `json:"rootVolumeSize,omitempty"`
		UseAsTemplateOnly        *bool                    `json:"useAsTemplateOnly,omitempty"`
		EBSOptimized             *gfnt.Value              `json:"ebsOptimized,omitempty"`
		SubnetIDs                interface{}              `json:"subnetIds,omitempty"`
		InstanceTypes            []string                 `json:"instanceTypes,omitempty"`
		InstanceMetadataOptions  *InstanceMetadataOptions `json:"instanceMetadataOptions,omitempty"`
		IAMInstanceProfile       map[string]*gfnt.Value   `json:"iamInstanceProfile,omitempty"`
		SecurityGroupIDs         *gfnt.Value              `json:"securityGroupIds,omitempty"`
		BlockDeviceMappings      []*BlockDevice           `json:"blockDeviceMappings,omitempty"`
		Tags                     []*Tag                   `json:"tags,omitempty"`
		LoadBalancers            []*LoadBalancer          `json:"loadBalancers,omitempty"`
		Labels                   []*Label                 `json:"labels,omitempty"`
		Taints                   []*Taint                 `json:"taints,omitempty"`
		AutoScaler               *AutoScaler              `json:"autoScale,omitempty"`
		Strategy                 *Strategy                `json:"strategy,omitempty"`
		Scheduling               *Scheduling              `json:"scheduling,omitempty"`
		ResourceLimits           *ResourceLimits          `json:"resourceLimits,omitempty"`
	}

	Strategy struct {
		SpotPercentage           *int  `json:"spotPercentage,omitempty"`
		UtilizeReservedInstances *bool `json:"utilizeReservedInstances,omitempty"`
		UtilizeCommitments       *bool `json:"utilizeCommitments,omitempty"`
		FallbackToOnDemand       *bool `json:"fallbackToOd,omitempty"`
		DrainingTimeout          *int  `json:"drainingTimeout,omitempty"`
	}

	Compute struct {
		SubnetIDs               interface{}              `json:"subnetIds,omitempty"`
		InstanceTypes           *InstanceTypes           `json:"instanceTypes,omitempty"`
		LaunchSpecification     *VirtualNodeGroup        `json:"launchSpecification,omitempty"`
		InstanceMetadataOptions *InstanceMetadataOptions `json:"instanceMetadataOptions,omitempty"`
	}

	InstanceTypes struct {
		Whitelist []string `json:"whitelist,omitempty"`
		Blacklist []string `json:"blacklist,omitempty"`
	}

	InstanceMetadataOptions struct {
		HttpPutResponseHopLimit *int    `json:"httpPutResponseHopLimit,omitempty"`
		HttpTokens              *string `json:"httpTokens,omitempty"`
	}

	LoadBalancer struct {
		Type *string `json:"type,omitempty"`
		ARN  *string `json:"arn,omitempty"`
		Name *string `json:"name,omitempty"`
	}

	BlockDevice struct {
		DeviceName *gfnt.Value     `json:"deviceName,omitempty"`
		EBS        *BlockDeviceEBS `json:"ebs,omitempty"`
	}

	BlockDeviceEBS struct {
		VolumeSize *gfnt.Value `json:"volumeSize,omitempty"`
		VolumeType *gfnt.Value `json:"volumeType,omitempty"`
		Encrypted  *gfnt.Value `json:"encrypted,omitempty"`
		KMSKeyID   *gfnt.Value `json:"kmsKeyId,omitempty"`
		IOPS       *gfnt.Value `json:"iops,omitempty"`
		Throughput *gfnt.Value `json:"throughput,omitempty"`
	}

	Tag struct {
		Key   interface{} `json:"tagKey,omitempty"`
		Value interface{} `json:"tagValue,omitempty"`
	}

	Scheduling struct {
		ShutdownHours *ShutdownHours `json:"shutdownHours,omitempty"`
		Tasks         []*Task        `json:"tasks,omitempty"`
	}

	ShutdownHours struct {
		IsEnabled   *bool    `json:"isEnabled,omitempty"`
		TimeWindows []string `json:"timeWindows,omitempty"`
	}

	Task struct {
		IsEnabled      *bool       `json:"isEnabled,omitempty"`
		Type           *string     `json:"taskType,omitempty"`
		CronExpression *string     `json:"cronExpression,omitempty"`
		Config         *TaskConfig `json:"config,omitempty"`
	}

	TaskConfig struct {
		Headrooms []*Headroom `json:"headrooms,omitempty"`
	}

	AutoScaler struct {
		IsEnabled      *bool           `json:"isEnabled,omitempty"`
		IsAutoConfig   *bool           `json:"isAutoConfig,omitempty"`
		Cooldown       *int            `json:"cooldown,omitempty"`
		ResourceLimits *ResourceLimits `json:"resourceLimits,omitempty"`
		Headroom       *Headroom       `json:"headroom,omitempty"`  // cluster
		Headrooms      []*Headroom     `json:"headrooms,omitempty"` // virtualnodegroup
	}

	Headroom struct {
		CPUPerUnit    *int `json:"cpuPerUnit,omitempty"`
		GPUPerUnit    *int `json:"gpuPerUnit,omitempty"`
		MemoryPerUnit *int `json:"memoryPerUnit,omitempty"`
		NumOfUnits    *int `json:"numOfUnits,omitempty"`
	}

	ResourceLimits struct {
		MaxVCPU          *int `json:"maxvCPU,omitempty"`
		MaxMemoryGiB     *int `json:"maxMemoryGib,omitempty"`
		MinInstanceCount *int `json:"minInstanceCount,omitempty"`
		MaxInstanceCount *int `json:"maxInstanceCount,omitempty"`
	}

	Label struct {
		Key   *string `json:"key,omitempty"`
		Value *string `json:"value,omitempty"`
	}

	Taint struct {
		Key    *string `json:"key,omitempty"`
		Value  *string `json:"value,omitempty"`
		Effect *string `json:"effect,omitempty"`
	}
)

// MarshalJSON implements the json.Marshaler interface.
func (x *ResourceNodeGroup) MarshalJSON() ([]byte, error) {
	var typ string
	if x.Cluster != nil {
		typ = "Custom::ocean"
	} else if x.VirtualNodeGroup != nil {
		typ = "Custom::oceanLaunchSpec"
	}
	type Properties ResourceNodeGroup
	return json.Marshal(&struct {
		Type       string
		Properties Properties
	}{
		Type:       typ,
		Properties: Properties(*x),
	})
}
