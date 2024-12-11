package ssmcontacts

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Contact_ContactTargetInfo AWS CloudFormation Resource (AWS::SSMContacts::Contact.ContactTargetInfo)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssmcontacts-contact-contacttargetinfo.html
type Contact_ContactTargetInfo struct {

	// ContactId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssmcontacts-contact-contacttargetinfo.html#cfn-ssmcontacts-contact-contacttargetinfo-contactid
	ContactId *types.Value `json:"ContactId,omitempty"`

	// IsEssential AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-ssmcontacts-contact-contacttargetinfo.html#cfn-ssmcontacts-contact-contacttargetinfo-isessential
	IsEssential *types.Value `json:"IsEssential"`

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
func (r *Contact_ContactTargetInfo) AWSCloudFormationType() string {
	return "AWS::SSMContacts::Contact.ContactTargetInfo"
}
