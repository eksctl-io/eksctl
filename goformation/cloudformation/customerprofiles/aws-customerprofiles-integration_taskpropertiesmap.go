package customerprofiles

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Integration_TaskPropertiesMap AWS CloudFormation Resource (AWS::CustomerProfiles::Integration.TaskPropertiesMap)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-taskpropertiesmap.html
type Integration_TaskPropertiesMap struct {

	// OperatorPropertyKey AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-taskpropertiesmap.html#cfn-customerprofiles-integration-taskpropertiesmap-operatorpropertykey
	OperatorPropertyKey *types.Value `json:"OperatorPropertyKey,omitempty"`

	// Property AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-taskpropertiesmap.html#cfn-customerprofiles-integration-taskpropertiesmap-property
	Property *types.Value `json:"Property,omitempty"`

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
func (r *Integration_TaskPropertiesMap) AWSCloudFormationType() string {
	return "AWS::CustomerProfiles::Integration.TaskPropertiesMap"
}
