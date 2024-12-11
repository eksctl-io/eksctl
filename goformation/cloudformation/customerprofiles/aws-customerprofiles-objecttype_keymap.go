package customerprofiles

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// ObjectType_KeyMap AWS CloudFormation Resource (AWS::CustomerProfiles::ObjectType.KeyMap)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-objecttype-keymap.html
type ObjectType_KeyMap struct {

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-objecttype-keymap.html#cfn-customerprofiles-objecttype-keymap-name
	Name *types.Value `json:"Name,omitempty"`

	// ObjectTypeKeyList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-customerprofiles-objecttype-keymap.html#cfn-customerprofiles-objecttype-keymap-objecttypekeylist
	ObjectTypeKeyList []ObjectType_ObjectTypeKey `json:"ObjectTypeKeyList,omitempty"`

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
func (r *ObjectType_KeyMap) AWSCloudFormationType() string {
	return "AWS::CustomerProfiles::ObjectType.KeyMap"
}
