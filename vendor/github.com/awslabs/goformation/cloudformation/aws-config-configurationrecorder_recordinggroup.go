package cloudformation

import (
	"encoding/json"
)

// AWSConfigConfigurationRecorder_RecordingGroup AWS CloudFormation Resource (AWS::Config::ConfigurationRecorder.RecordingGroup)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configurationrecorder-recordinggroup.html
type AWSConfigConfigurationRecorder_RecordingGroup struct {

	// AllSupported AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configurationrecorder-recordinggroup.html#cfn-config-configurationrecorder-recordinggroup-allsupported
	AllSupported *Value `json:"AllSupported,omitempty"`

	// IncludeGlobalResourceTypes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configurationrecorder-recordinggroup.html#cfn-config-configurationrecorder-recordinggroup-includeglobalresourcetypes
	IncludeGlobalResourceTypes *Value `json:"IncludeGlobalResourceTypes,omitempty"`

	// ResourceTypes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-config-configurationrecorder-recordinggroup.html#cfn-config-configurationrecorder-recordinggroup-resourcetypes
	ResourceTypes []*Value `json:"ResourceTypes,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSConfigConfigurationRecorder_RecordingGroup) AWSCloudFormationType() string {
	return "AWS::Config::ConfigurationRecorder.RecordingGroup"
}

func (r *AWSConfigConfigurationRecorder_RecordingGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
