package cloudformation

import (
	"encoding/json"
)

// AWSSSMMaintenanceWindowTask_LoggingInfo AWS CloudFormation Resource (AWS::SSM::MaintenanceWindowTask.LoggingInfo)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-maintenancewindowtask-logginginfo.html
type AWSSSMMaintenanceWindowTask_LoggingInfo struct {

	// Region AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-maintenancewindowtask-logginginfo.html#cfn-ssm-maintenancewindowtask-logginginfo-region
	Region *Value `json:"Region,omitempty"`

	// S3Bucket AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-maintenancewindowtask-logginginfo.html#cfn-ssm-maintenancewindowtask-logginginfo-s3bucket
	S3Bucket *Value `json:"S3Bucket,omitempty"`

	// S3Prefix AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssm-maintenancewindowtask-logginginfo.html#cfn-ssm-maintenancewindowtask-logginginfo-s3prefix
	S3Prefix *Value `json:"S3Prefix,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSSMMaintenanceWindowTask_LoggingInfo) AWSCloudFormationType() string {
	return "AWS::SSM::MaintenanceWindowTask.LoggingInfo"
}

func (r *AWSSSMMaintenanceWindowTask_LoggingInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
