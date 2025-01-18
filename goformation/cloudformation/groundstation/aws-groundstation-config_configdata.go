package groundstation

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Config_ConfigData AWS CloudFormation Resource (AWS::GroundStation::Config.ConfigData)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-groundstation-config-configdata.html
type Config_ConfigData struct {

	// AntennaDownlinkConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-groundstation-config-configdata.html#cfn-groundstation-config-configdata-antennadownlinkconfig
	AntennaDownlinkConfig *Config_AntennaDownlinkConfig `json:"AntennaDownlinkConfig,omitempty"`

	// AntennaDownlinkDemodDecodeConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-groundstation-config-configdata.html#cfn-groundstation-config-configdata-antennadownlinkdemoddecodeconfig
	AntennaDownlinkDemodDecodeConfig *Config_AntennaDownlinkDemodDecodeConfig `json:"AntennaDownlinkDemodDecodeConfig,omitempty"`

	// AntennaUplinkConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-groundstation-config-configdata.html#cfn-groundstation-config-configdata-antennauplinkconfig
	AntennaUplinkConfig *Config_AntennaUplinkConfig `json:"AntennaUplinkConfig,omitempty"`

	// DataflowEndpointConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-groundstation-config-configdata.html#cfn-groundstation-config-configdata-dataflowendpointconfig
	DataflowEndpointConfig *Config_DataflowEndpointConfig `json:"DataflowEndpointConfig,omitempty"`

	// S3RecordingConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-groundstation-config-configdata.html#cfn-groundstation-config-configdata-s3recordingconfig
	S3RecordingConfig *Config_S3RecordingConfig `json:"S3RecordingConfig,omitempty"`

	// TrackingConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-groundstation-config-configdata.html#cfn-groundstation-config-configdata-trackingconfig
	TrackingConfig *Config_TrackingConfig `json:"TrackingConfig,omitempty"`

	// UplinkEchoConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-groundstation-config-configdata.html#cfn-groundstation-config-configdata-uplinkechoconfig
	UplinkEchoConfig *Config_UplinkEchoConfig `json:"UplinkEchoConfig,omitempty"`

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
func (r *Config_ConfigData) AWSCloudFormationType() string {
	return "AWS::GroundStation::Config.ConfigData"
}
