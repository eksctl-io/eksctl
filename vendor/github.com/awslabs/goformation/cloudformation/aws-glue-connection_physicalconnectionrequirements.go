package cloudformation

import (
	"encoding/json"
)

// AWSGlueConnection_PhysicalConnectionRequirements AWS CloudFormation Resource (AWS::Glue::Connection.PhysicalConnectionRequirements)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-connection-physicalconnectionrequirements.html
type AWSGlueConnection_PhysicalConnectionRequirements struct {

	// AvailabilityZone AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-connection-physicalconnectionrequirements.html#cfn-glue-connection-physicalconnectionrequirements-availabilityzone
	AvailabilityZone *Value `json:"AvailabilityZone,omitempty"`

	// SecurityGroupIdList AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-connection-physicalconnectionrequirements.html#cfn-glue-connection-physicalconnectionrequirements-securitygroupidlist
	SecurityGroupIdList []*Value `json:"SecurityGroupIdList,omitempty"`

	// SubnetId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-glue-connection-physicalconnectionrequirements.html#cfn-glue-connection-physicalconnectionrequirements-subnetid
	SubnetId *Value `json:"SubnetId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSGlueConnection_PhysicalConnectionRequirements) AWSCloudFormationType() string {
	return "AWS::Glue::Connection.PhysicalConnectionRequirements"
}

func (r *AWSGlueConnection_PhysicalConnectionRequirements) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
