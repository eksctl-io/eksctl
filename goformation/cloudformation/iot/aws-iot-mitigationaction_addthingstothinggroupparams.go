package iot

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// MitigationAction_AddThingsToThingGroupParams AWS CloudFormation Resource (AWS::IoT::MitigationAction.AddThingsToThingGroupParams)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-mitigationaction-addthingstothinggroupparams.html
type MitigationAction_AddThingsToThingGroupParams struct {

	// OverrideDynamicGroups AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-mitigationaction-addthingstothinggroupparams.html#cfn-iot-mitigationaction-addthingstothinggroupparams-overridedynamicgroups
	OverrideDynamicGroups *types.Value `json:"OverrideDynamicGroups,omitempty"`

	// ThingGroupNames AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-mitigationaction-addthingstothinggroupparams.html#cfn-iot-mitigationaction-addthingstothinggroupparams-thinggroupnames
	ThingGroupNames *types.Value `json:"ThingGroupNames,omitempty"`

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
func (r *MitigationAction_AddThingsToThingGroupParams) AWSCloudFormationType() string {
	return "AWS::IoT::MitigationAction.AddThingsToThingGroupParams"
}
