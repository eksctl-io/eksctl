package mediapackage

import (
	"github.com/awslabs/goformation/v4/cloudformation/policies"
)

// OriginEndpoint_AdTriggers AWS CloudFormation Resource (AWS::MediaPackage::OriginEndpoint.AdTriggers)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-adtriggers.html
type OriginEndpoint_AdTriggers struct {

	// AdTriggers AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-mediapackage-originendpoint-adtriggers.html#cfn-mediapackage-originendpoint-adtriggers-adtriggers
	AdTriggers []string `json:"AdTriggers,omitempty"`

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
func (r *OriginEndpoint_AdTriggers) AWSCloudFormationType() string {
	return "AWS::MediaPackage::OriginEndpoint.AdTriggers"
}
