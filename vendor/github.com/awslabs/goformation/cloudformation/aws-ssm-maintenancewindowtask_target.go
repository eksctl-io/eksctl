package cloudformation

import (
	"encoding/json"
)

// AWSSSMMaintenanceWindowTask_Target AWS CloudFormation Resource (AWS::SSM::MaintenanceWindowTask.Target)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-maintenancewindowtask-target.html
type AWSSSMMaintenanceWindowTask_Target struct {

	// Key AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-maintenancewindowtask-target.html#cfn-ssm-maintenancewindowtask-target-key
	Key *Value `json:"Key,omitempty"`

	// Values AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-maintenancewindowtask-target.html#cfn-ssm-maintenancewindowtask-target-values
	Values []*Value `json:"Values,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSSMMaintenanceWindowTask_Target) AWSCloudFormationType() string {
	return "AWS::SSM::MaintenanceWindowTask.Target"
}

func (r *AWSSSMMaintenanceWindowTask_Target) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
