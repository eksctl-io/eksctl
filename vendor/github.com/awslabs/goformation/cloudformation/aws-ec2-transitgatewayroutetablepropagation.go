package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSEC2TransitGatewayRouteTablePropagation AWS CloudFormation Resource (AWS::EC2::TransitGatewayRouteTablePropagation)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-transitgatewayroutetablepropagation.html
type AWSEC2TransitGatewayRouteTablePropagation struct {

	// TransitGatewayAttachmentId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-transitgatewayroutetablepropagation.html#cfn-ec2-transitgatewayroutetablepropagation-transitgatewayattachmentid
	TransitGatewayAttachmentId *Value `json:"TransitGatewayAttachmentId,omitempty"`

	// TransitGatewayRouteTableId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-transitgatewayroutetablepropagation.html#cfn-ec2-transitgatewayroutetablepropagation-transitgatewayroutetableid
	TransitGatewayRouteTableId *Value `json:"TransitGatewayRouteTableId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2TransitGatewayRouteTablePropagation) AWSCloudFormationType() string {
	return "AWS::EC2::TransitGatewayRouteTablePropagation"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSEC2TransitGatewayRouteTablePropagation) MarshalJSON() ([]byte, error) {
	type Properties AWSEC2TransitGatewayRouteTablePropagation
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
func (r *AWSEC2TransitGatewayRouteTablePropagation) UnmarshalJSON(b []byte) error {
	type Properties AWSEC2TransitGatewayRouteTablePropagation
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
		*r = AWSEC2TransitGatewayRouteTablePropagation(*res.Properties)
	}

	return nil
}

// GetAllAWSEC2TransitGatewayRouteTablePropagationResources retrieves all AWSEC2TransitGatewayRouteTablePropagation items from an AWS CloudFormation template
func (t *Template) GetAllAWSEC2TransitGatewayRouteTablePropagationResources() map[string]AWSEC2TransitGatewayRouteTablePropagation {
	results := map[string]AWSEC2TransitGatewayRouteTablePropagation{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSEC2TransitGatewayRouteTablePropagation:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::EC2::TransitGatewayRouteTablePropagation" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSEC2TransitGatewayRouteTablePropagation{}
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

// GetAWSEC2TransitGatewayRouteTablePropagationWithName retrieves all AWSEC2TransitGatewayRouteTablePropagation items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSEC2TransitGatewayRouteTablePropagationWithName(name string) (AWSEC2TransitGatewayRouteTablePropagation, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSEC2TransitGatewayRouteTablePropagation:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::EC2::TransitGatewayRouteTablePropagation" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSEC2TransitGatewayRouteTablePropagation{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSEC2TransitGatewayRouteTablePropagation{}, errors.New("resource not found")
}
