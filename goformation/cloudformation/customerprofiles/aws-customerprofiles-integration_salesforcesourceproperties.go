package customerprofiles

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Integration_SalesforceSourceProperties AWS CloudFormation Resource (AWS::CustomerProfiles::Integration.SalesforceSourceProperties)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-salesforcesourceproperties.html
type Integration_SalesforceSourceProperties struct {

	// EnableDynamicFieldUpdate AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-salesforcesourceproperties.html#cfn-customerprofiles-integration-salesforcesourceproperties-enabledynamicfieldupdate
	EnableDynamicFieldUpdate *types.Value `json:"EnableDynamicFieldUpdate,omitempty"`

	// IncludeDeletedRecords AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-salesforcesourceproperties.html#cfn-customerprofiles-integration-salesforcesourceproperties-includedeletedrecords
	IncludeDeletedRecords *types.Value `json:"IncludeDeletedRecords,omitempty"`

	// Object AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-integration-salesforcesourceproperties.html#cfn-customerprofiles-integration-salesforcesourceproperties-object
	Object *types.Value `json:"Object,omitempty"`

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
func (r *Integration_SalesforceSourceProperties) AWSCloudFormationType() string {
	return "AWS::CustomerProfiles::Integration.SalesforceSourceProperties"
}
