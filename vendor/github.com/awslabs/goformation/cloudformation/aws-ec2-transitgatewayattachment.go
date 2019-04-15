package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSEC2TransitGatewayAttachment AWS CloudFormation Resource (AWS::EC2::TransitGatewayAttachment)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-transitgatewayattachment.html
type AWSEC2TransitGatewayAttachment struct {

	// SubnetIds AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-transitgatewayattachment.html#cfn-ec2-transitgatewayattachment-subnetids
	SubnetIds []*Value `json:"SubnetIds,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-transitgatewayattachment.html#cfn-ec2-transitgatewayattachment-tags
	Tags []Tag `json:"Tags,omitempty"`

	// TransitGatewayId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-transitgatewayattachment.html#cfn-ec2-transitgatewayattachment-transitgatewayid
	TransitGatewayId *Value `json:"TransitGatewayId,omitempty"`

	// VpcId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-transitgatewayattachment.html#cfn-ec2-transitgatewayattachment-vpcid
	VpcId *Value `json:"VpcId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2TransitGatewayAttachment) AWSCloudFormationType() string {
	return "AWS::EC2::TransitGatewayAttachment"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSEC2TransitGatewayAttachment) MarshalJSON() ([]byte, error) {
	type Properties AWSEC2TransitGatewayAttachment
	return json.Marshal(&struct {
		Type       string
		Properties Properties
	}{
		Type:       r.AWSCloudFormationType(),
		Properties: (Properties)(*r),
	})
}

// UnmarshalJSON is a custom JSON unmarshalling hook that strips the outer
// AWS CloudFormation resource object, and just keeps the 'Properties' field.
func (r *AWSEC2TransitGatewayAttachment) UnmarshalJSON(b []byte) error {
	type Properties AWSEC2TransitGatewayAttachment
	res := &struct {
		Type       string
		Properties *Properties
	}{}
	if err := json.Unmarshal(b, &res); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return err
	}

	// If the resource has no Properties set, it could be nil
	if res.Properties != nil {
		*r = AWSEC2TransitGatewayAttachment(*res.Properties)
	}

	return nil
}

// GetAllAWSEC2TransitGatewayAttachmentResources retrieves all AWSEC2TransitGatewayAttachment items from an AWS CloudFormation template
func (t *Template) GetAllAWSEC2TransitGatewayAttachmentResources() map[string]AWSEC2TransitGatewayAttachment {
	results := map[string]AWSEC2TransitGatewayAttachment{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSEC2TransitGatewayAttachment:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::EC2::TransitGatewayAttachment" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSEC2TransitGatewayAttachment{}
						if err := result.UnmarshalJSON(b); err == nil {
							results[name] = *result
						}
					}
				}
			}
		}
	}
	return results
}

// GetAWSEC2TransitGatewayAttachmentWithName retrieves all AWSEC2TransitGatewayAttachment items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSEC2TransitGatewayAttachmentWithName(name string) (AWSEC2TransitGatewayAttachment, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSEC2TransitGatewayAttachment:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::EC2::TransitGatewayAttachment" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSEC2TransitGatewayAttachment{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSEC2TransitGatewayAttachment{}, errors.New("resource not found")
}
