package evidently

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Project_DataDeliveryObject AWS CloudFormation Resource (AWS::Evidently::Project.DataDeliveryObject)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-project-datadeliveryobject.html
type Project_DataDeliveryObject struct {

	// LogGroup AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-project-datadeliveryobject.html#cfn-evidently-project-datadeliveryobject-loggroup
	LogGroup *types.Value `json:"LogGroup,omitempty"`

	// S3 AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-evidently-project-datadeliveryobject.html#cfn-evidently-project-datadeliveryobject-s3
	S3 interface{} `json:"S3,omitempty"`

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
func (r *Project_DataDeliveryObject) AWSCloudFormationType() string {
	return "AWS::Evidently::Project.DataDeliveryObject"
}
