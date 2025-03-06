package route53

import (
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/policies"
)

// RecordSet_GeoProximityLocation AWS CloudFormation Resource (AWS::Route53::RecordSet.GeoProximityLocation)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53-geoproximitylocation.html
type RecordSet_GeoProximityLocation struct {

	// AWSRegion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53-geoproximitylocation.html#cfn-route53-geoproximitylocation-awsregion
	AWSRegion *types.Value `json:"AWSRegion,omitempty"`

	// Bias AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53-geoproximitylocation.html#cfn-route53-geoproximitylocation-bias
	Bias *types.Value `json:"Bias,omitempty"`

	// Coordinates AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53-geoproximitylocation.html#cfn-route53-geoproximitylocation-coordinates
	Coordinates *RecordSet_Coordinates `json:"Coordinates,omitempty"`

	// LocalZoneGroup AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-route53-geoproximitylocation.html#cfn-route53-geoproximitylocation-LocalZoneGroup
	LocalZoneGroup *types.Value `json:"LocalZoneGroup,omitempty"`

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
func (r *RecordSet_GeoProximityLocation) AWSCloudFormationType() string {
	return "AWS::Route53::RecordSet.GeoProximityLocation"
}
