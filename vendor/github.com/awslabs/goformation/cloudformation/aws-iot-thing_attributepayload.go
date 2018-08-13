package cloudformation

import (
	"encoding/json"
)

// AWSIoTThing_AttributePayload AWS CloudFormation Resource (AWS::IoT::Thing.AttributePayload)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-thing-attributepayload.html
type AWSIoTThing_AttributePayload struct {

	// Attributes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-iot-thing-attributepayload.html#cfn-iot-thing-attributepayload-attributes
	Attributes map[string]*Value `json:"Attributes,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTThing_AttributePayload) AWSCloudFormationType() string {
	return "AWS::IoT::Thing.AttributePayload"
}

func (r *AWSIoTThing_AttributePayload) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
