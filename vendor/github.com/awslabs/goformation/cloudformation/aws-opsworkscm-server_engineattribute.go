package cloudformation

import (
	"encoding/json"
)

// AWSOpsWorksCMServer_EngineAttribute AWS CloudFormation Resource (AWS::OpsWorksCM::Server.EngineAttribute)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworkscm-server-engineattribute.html
type AWSOpsWorksCMServer_EngineAttribute struct {

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworkscm-server-engineattribute.html#cfn-opsworkscm-server-engineattribute-name
	Name *Value `json:"Name,omitempty"`

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworkscm-server-engineattribute.html#cfn-opsworkscm-server-engineattribute-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSOpsWorksCMServer_EngineAttribute) AWSCloudFormationType() string {
	return "AWS::OpsWorksCM::Server.EngineAttribute"
}

func (r *AWSOpsWorksCMServer_EngineAttribute) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
