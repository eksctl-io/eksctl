package ssmcontacts

import (
	"github.com/awslabs/goformation/v4/cloudformation/types"

	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// Contact_Stage AWS CloudFormation Resource (AWS::SSMContacts::Contact.Stage)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssmcontacts-contact-stage.html
type Contact_Stage struct {

	// DurationInMinutes AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssmcontacts-contact-stage.html#cfn-ssmcontacts-contact-stage-durationinminutes
	DurationInMinutes *types.Value `json:"DurationInMinutes"`

	// Targets AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssmcontacts-contact-stage.html#cfn-ssmcontacts-contact-stage-targets
	Targets []Contact_Targets `json:"Targets,omitempty"`

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
func (r *Contact_Stage) AWSCloudFormationType() string {
	return "AWS::SSMContacts::Contact.Stage"
}
