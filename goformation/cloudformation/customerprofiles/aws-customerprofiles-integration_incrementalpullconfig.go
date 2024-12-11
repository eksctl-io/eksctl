package customerprofiles

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Integration_IncrementalPullConfig AWS CloudFormation Resource (AWS::CustomerProfiles::Integration.IncrementalPullConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-incrementalpullconfig.html
type Integration_IncrementalPullConfig struct {

	// DatetimeTypeFieldName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-incrementalpullconfig.html#cfn-customerprofiles-integration-incrementalpullconfig-datetimetypefieldname
	DatetimeTypeFieldName *types.Value `json:"DatetimeTypeFieldName,omitempty"`

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
func (r *Integration_IncrementalPullConfig) AWSCloudFormationType() string {
	return "AWS::CustomerProfiles::Integration.IncrementalPullConfig"
}
