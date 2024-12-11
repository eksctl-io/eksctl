package synthetics

import (
	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/policies"
)

// Canary_BaseScreenshot AWS CloudFormation Resource (AWS::Synthetics::Canary.BaseScreenshot)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-synthetics-canary-basescreenshot.html
type Canary_BaseScreenshot struct {

	// IgnoreCoordinates AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-synthetics-canary-basescreenshot.html#cfn-synthetics-canary-basescreenshot-ignorecoordinates
	IgnoreCoordinates *types.Value `json:"IgnoreCoordinates,omitempty"`

	// ScreenshotName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-synthetics-canary-basescreenshot.html#cfn-synthetics-canary-basescreenshot-screenshotname
	ScreenshotName *types.Value `json:"ScreenshotName,omitempty"`

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
func (r *Canary_BaseScreenshot) AWSCloudFormationType() string {
	return "AWS::Synthetics::Canary.BaseScreenshot"
}
