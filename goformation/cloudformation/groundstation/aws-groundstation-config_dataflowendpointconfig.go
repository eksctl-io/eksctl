package groundstation

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Config_DataflowEndpointConfig AWS CloudFormation Resource (AWS::GroundStation::Config.DataflowEndpointConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-groundstation-config-dataflowendpointconfig.html
type Config_DataflowEndpointConfig struct {

	// DataflowEndpointName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-groundstation-config-dataflowendpointconfig.html#cfn-groundstation-config-dataflowendpointconfig-dataflowendpointname
	DataflowEndpointName *types.Value `json:"DataflowEndpointName,omitempty"`

	// DataflowEndpointRegion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-groundstation-config-dataflowendpointconfig.html#cfn-groundstation-config-dataflowendpointconfig-dataflowendpointregion
	DataflowEndpointRegion *types.Value `json:"DataflowEndpointRegion,omitempty"`

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
func (r *Config_DataflowEndpointConfig) AWSCloudFormationType() string {
	return "AWS::GroundStation::Config.DataflowEndpointConfig"
}
