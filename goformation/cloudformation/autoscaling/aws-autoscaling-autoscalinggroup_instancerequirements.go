package autoscaling

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// AutoScalingGroup_InstanceRequirements AWS CloudFormation Resource (AWS::AutoScaling::AutoScalingGroup.InstanceRequirements)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html
type AutoScalingGroup_InstanceRequirements struct {

	// AcceleratorCount AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html#cfn-autoscaling-autoscalinggroup-instancerequirements-acceleratorcount
	AcceleratorCount *AutoScalingGroup_AcceleratorCountRequest `json:"AcceleratorCount,omitempty"`

	// AcceleratorManufacturers AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html#cfn-autoscaling-autoscalinggroup-instancerequirements-acceleratormanufacturers
	AcceleratorManufacturers *types.Value `json:"AcceleratorManufacturers,omitempty"`

	// AcceleratorNames AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html#cfn-autoscaling-autoscalinggroup-instancerequirements-acceleratornames
	AcceleratorNames *types.Value `json:"AcceleratorNames,omitempty"`

	// AcceleratorTotalMemoryMiB AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html#cfn-autoscaling-autoscalinggroup-instancerequirements-acceleratortotalmemorymib
	AcceleratorTotalMemoryMiB *AutoScalingGroup_AcceleratorTotalMemoryMiBRequest `json:"AcceleratorTotalMemoryMiB,omitempty"`

	// AcceleratorTypes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html#cfn-autoscaling-autoscalinggroup-instancerequirements-acceleratortypes
	AcceleratorTypes *types.Value `json:"AcceleratorTypes,omitempty"`

	// BareMetal AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html#cfn-autoscaling-autoscalinggroup-instancerequirements-baremetal
	BareMetal *types.Value `json:"BareMetal,omitempty"`

	// BaselineEbsBandwidthMbps AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html#cfn-autoscaling-autoscalinggroup-instancerequirements-baselineebsbandwidthmbps
	BaselineEbsBandwidthMbps *AutoScalingGroup_BaselineEbsBandwidthMbpsRequest `json:"BaselineEbsBandwidthMbps,omitempty"`

	// BurstablePerformance AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html#cfn-autoscaling-autoscalinggroup-instancerequirements-burstableperformance
	BurstablePerformance *types.Value `json:"BurstablePerformance,omitempty"`

	// CpuManufacturers AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html#cfn-autoscaling-autoscalinggroup-instancerequirements-cpumanufacturers
	CpuManufacturers *types.Value `json:"CpuManufacturers,omitempty"`

	// ExcludedInstanceTypes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html#cfn-autoscaling-autoscalinggroup-instancerequirements-excludedinstancetypes
	ExcludedInstanceTypes *types.Value `json:"ExcludedInstanceTypes,omitempty"`

	// InstanceGenerations AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html#cfn-autoscaling-autoscalinggroup-instancerequirements-instancegenerations
	InstanceGenerations *types.Value `json:"InstanceGenerations,omitempty"`

	// LocalStorage AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html#cfn-autoscaling-autoscalinggroup-instancerequirements-localstorage
	LocalStorage *types.Value `json:"LocalStorage,omitempty"`

	// LocalStorageTypes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html#cfn-autoscaling-autoscalinggroup-instancerequirements-localstoragetypes
	LocalStorageTypes *types.Value `json:"LocalStorageTypes,omitempty"`

	// MemoryGiBPerVCpu AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html#cfn-autoscaling-autoscalinggroup-instancerequirements-memorygibpervcpu
	MemoryGiBPerVCpu *AutoScalingGroup_MemoryGiBPerVCpuRequest `json:"MemoryGiBPerVCpu,omitempty"`

	// MemoryMiB AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html#cfn-autoscaling-autoscalinggroup-instancerequirements-memorymib
	MemoryMiB *AutoScalingGroup_MemoryMiBRequest `json:"MemoryMiB,omitempty"`

	// NetworkInterfaceCount AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html#cfn-autoscaling-autoscalinggroup-instancerequirements-networkinterfacecount
	NetworkInterfaceCount *AutoScalingGroup_NetworkInterfaceCountRequest `json:"NetworkInterfaceCount,omitempty"`

	// OnDemandMaxPricePercentageOverLowestPrice AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html#cfn-autoscaling-autoscalinggroup-instancerequirements-ondemandmaxpricepercentageoverlowestprice
	OnDemandMaxPricePercentageOverLowestPrice *types.Value `json:"OnDemandMaxPricePercentageOverLowestPrice,omitempty"`

	// RequireHibernateSupport AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html#cfn-autoscaling-autoscalinggroup-instancerequirements-requirehibernatesupport
	RequireHibernateSupport *types.Value `json:"RequireHibernateSupport,omitempty"`

	// SpotMaxPricePercentageOverLowestPrice AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html#cfn-autoscaling-autoscalinggroup-instancerequirements-spotmaxpricepercentageoverlowestprice
	SpotMaxPricePercentageOverLowestPrice *types.Value `json:"SpotMaxPricePercentageOverLowestPrice,omitempty"`

	// TotalLocalStorageGB AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html#cfn-autoscaling-autoscalinggroup-instancerequirements-totallocalstoragegb
	TotalLocalStorageGB *AutoScalingGroup_TotalLocalStorageGBRequest `json:"TotalLocalStorageGB,omitempty"`

	// VCpuCount AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/cfn-as-mixedinstancespolicy-instancerequirements.html#cfn-autoscaling-autoscalinggroup-instancerequirements-vcpucount
	VCpuCount *AutoScalingGroup_VCpuCountRequest `json:"VCpuCount,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AutoScalingGroup_InstanceRequirements) AWSCloudFormationType() string {
	return "AWS::AutoScaling::AutoScalingGroup.InstanceRequirements"
}
