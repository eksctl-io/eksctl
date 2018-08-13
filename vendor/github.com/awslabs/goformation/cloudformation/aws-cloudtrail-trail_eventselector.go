package cloudformation

import (
	"encoding/json"
)

// AWSCloudTrailTrail_EventSelector AWS CloudFormation Resource (AWS::CloudTrail::Trail.EventSelector)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudtrail-trail-eventselector.html
type AWSCloudTrailTrail_EventSelector struct {

	// DataResources AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudtrail-trail-eventselector.html#cfn-cloudtrail-trail-eventselector-dataresources
	DataResources []AWSCloudTrailTrail_DataResource `json:"DataResources,omitempty"`

	// IncludeManagementEvents AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudtrail-trail-eventselector.html#cfn-cloudtrail-trail-eventselector-includemanagementevents
	IncludeManagementEvents *Value `json:"IncludeManagementEvents,omitempty"`

	// ReadWriteType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-cloudtrail-trail-eventselector.html#cfn-cloudtrail-trail-eventselector-readwritetype
	ReadWriteType *Value `json:"ReadWriteType,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCloudTrailTrail_EventSelector) AWSCloudFormationType() string {
	return "AWS::CloudTrail::Trail.EventSelector"
}

func (r *AWSCloudTrailTrail_EventSelector) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
