package backup

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Framework_FrameworkControl AWS CloudFormation Resource (AWS::Backup::Framework.FrameworkControl)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-backup-framework-frameworkcontrol.html
type Framework_FrameworkControl struct {

	// ControlInputParameters AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-backup-framework-frameworkcontrol.html#cfn-backup-framework-frameworkcontrol-controlinputparameters
	ControlInputParameters []Framework_ControlInputParameter `json:"ControlInputParameters,omitempty"`

	// ControlName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-backup-framework-frameworkcontrol.html#cfn-backup-framework-frameworkcontrol-controlname
	ControlName *types.Value `json:"ControlName,omitempty"`

	// ControlScope AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-backup-framework-frameworkcontrol.html#cfn-backup-framework-frameworkcontrol-controlscope
	ControlScope interface{} `json:"ControlScope,omitempty"`

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
func (r *Framework_FrameworkControl) AWSCloudFormationType() string {
	return "AWS::Backup::Framework.FrameworkControl"
}
