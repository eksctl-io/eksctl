package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSIoTAnalyticsChannel AWS CloudFormation Resource (AWS::IoTAnalytics::Channel)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iotanalytics-channel.html
type AWSIoTAnalyticsChannel struct {

	// ChannelName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iotanalytics-channel.html#cfn-iotanalytics-channel-channelname
	ChannelName *Value `json:"ChannelName,omitempty"`

	// RetentionPeriod AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iotanalytics-channel.html#cfn-iotanalytics-channel-retentionperiod
	RetentionPeriod *AWSIoTAnalyticsChannel_RetentionPeriod `json:"RetentionPeriod,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iotanalytics-channel.html#cfn-iotanalytics-channel-tags
	Tags []Tag `json:"Tags,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsChannel) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Channel"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSIoTAnalyticsChannel) MarshalJSON() ([]byte, error) {
	type Properties AWSIoTAnalyticsChannel
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
func (r *AWSIoTAnalyticsChannel) UnmarshalJSON(b []byte) error {
	type Properties AWSIoTAnalyticsChannel
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
		*r = AWSIoTAnalyticsChannel(*res.Properties)
	}

	return nil
}

// GetAllAWSIoTAnalyticsChannelResources retrieves all AWSIoTAnalyticsChannel items from an AWS CloudFormation template
func (t *Template) GetAllAWSIoTAnalyticsChannelResources() map[string]AWSIoTAnalyticsChannel {
	results := map[string]AWSIoTAnalyticsChannel{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSIoTAnalyticsChannel:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::IoTAnalytics::Channel" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSIoTAnalyticsChannel{}
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

// GetAWSIoTAnalyticsChannelWithName retrieves all AWSIoTAnalyticsChannel items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSIoTAnalyticsChannelWithName(name string) (AWSIoTAnalyticsChannel, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSIoTAnalyticsChannel:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::IoTAnalytics::Channel" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSIoTAnalyticsChannel{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSIoTAnalyticsChannel{}, errors.New("resource not found")
}
