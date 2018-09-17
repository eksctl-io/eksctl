package cloudformation

import (
	"encoding/json"
)

// AWSEMRInstanceFleetConfig_SpotProvisioningSpecification AWS CloudFormation Resource (AWS::EMR::InstanceFleetConfig.SpotProvisioningSpecification)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-instancefleetconfig-spotprovisioningspecification.html
type AWSEMRInstanceFleetConfig_SpotProvisioningSpecification struct {

	// BlockDurationMinutes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-instancefleetconfig-spotprovisioningspecification.html#cfn-elasticmapreduce-instancefleetconfig-spotprovisioningspecification-blockdurationminutes
	BlockDurationMinutes *Value `json:"BlockDurationMinutes,omitempty"`

	// TimeoutAction AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-instancefleetconfig-spotprovisioningspecification.html#cfn-elasticmapreduce-instancefleetconfig-spotprovisioningspecification-timeoutaction
	TimeoutAction *Value `json:"TimeoutAction,omitempty"`

	// TimeoutDurationMinutes AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-elasticmapreduce-instancefleetconfig-spotprovisioningspecification.html#cfn-elasticmapreduce-instancefleetconfig-spotprovisioningspecification-timeoutdurationminutes
	TimeoutDurationMinutes *Value `json:"TimeoutDurationMinutes,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEMRInstanceFleetConfig_SpotProvisioningSpecification) AWSCloudFormationType() string {
	return "AWS::EMR::InstanceFleetConfig.SpotProvisioningSpecification"
}

func (r *AWSEMRInstanceFleetConfig_SpotProvisioningSpecification) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
