package cloudformation

import (
	"encoding/json"
)

// AWSEMRInstanceGroupConfig_EbsConfiguration AWS CloudFormation Resource (AWS::EMR::InstanceGroupConfig.EbsConfiguration)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-emr-ebsconfiguration.html
type AWSEMRInstanceGroupConfig_EbsConfiguration struct {

	// EbsBlockDeviceConfigs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-emr-ebsconfiguration.html#cfn-emr-ebsconfiguration-ebsblockdeviceconfigs
	EbsBlockDeviceConfigs []AWSEMRInstanceGroupConfig_EbsBlockDeviceConfig `json:"EbsBlockDeviceConfigs,omitempty"`

	// EbsOptimized AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-emr-ebsconfiguration.html#cfn-emr-ebsconfiguration-ebsoptimized
	EbsOptimized *Value `json:"EbsOptimized,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEMRInstanceGroupConfig_EbsConfiguration) AWSCloudFormationType() string {
	return "AWS::EMR::InstanceGroupConfig.EbsConfiguration"
}

func (r *AWSEMRInstanceGroupConfig_EbsConfiguration) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
