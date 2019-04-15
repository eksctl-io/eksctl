package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSIoT1ClickDevice AWS CloudFormation Resource (AWS::IoT1Click::Device)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iot1click-device.html
type AWSIoT1ClickDevice struct {

	// DeviceId AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iot1click-device.html#cfn-iot1click-device-deviceid
	DeviceId *Value `json:"DeviceId,omitempty"`

	// Enabled AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iot1click-device.html#cfn-iot1click-device-enabled
	Enabled *Value `json:"Enabled,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoT1ClickDevice) AWSCloudFormationType() string {
	return "AWS::IoT1Click::Device"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSIoT1ClickDevice) MarshalJSON() ([]byte, error) {
	type Properties AWSIoT1ClickDevice
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
func (r *AWSIoT1ClickDevice) UnmarshalJSON(b []byte) error {
	type Properties AWSIoT1ClickDevice
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
		*r = AWSIoT1ClickDevice(*res.Properties)
	}

	return nil
}

// GetAllAWSIoT1ClickDeviceResources retrieves all AWSIoT1ClickDevice items from an AWS CloudFormation template
func (t *Template) GetAllAWSIoT1ClickDeviceResources() map[string]AWSIoT1ClickDevice {
	results := map[string]AWSIoT1ClickDevice{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSIoT1ClickDevice:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::IoT1Click::Device" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSIoT1ClickDevice{}
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

// GetAWSIoT1ClickDeviceWithName retrieves all AWSIoT1ClickDevice items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSIoT1ClickDeviceWithName(name string) (AWSIoT1ClickDevice, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSIoT1ClickDevice:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::IoT1Click::Device" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSIoT1ClickDevice{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSIoT1ClickDevice{}, errors.New("resource not found")
}
