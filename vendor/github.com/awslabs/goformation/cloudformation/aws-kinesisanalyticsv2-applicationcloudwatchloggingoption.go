package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::ApplicationCloudWatchLoggingOption)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-kinesisanalyticsv2-applicationcloudwatchloggingoption.html
type AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption struct {

	// ApplicationName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-kinesisanalyticsv2-applicationcloudwatchloggingoption.html#cfn-kinesisanalyticsv2-applicationcloudwatchloggingoption-applicationname
	ApplicationName *Value `json:"ApplicationName,omitempty"`

	// CloudWatchLoggingOption AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-kinesisanalyticsv2-applicationcloudwatchloggingoption.html#cfn-kinesisanalyticsv2-applicationcloudwatchloggingoption-cloudwatchloggingoption
	CloudWatchLoggingOption *AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption_CloudWatchLoggingOption `json:"CloudWatchLoggingOption,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::ApplicationCloudWatchLoggingOption"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption) MarshalJSON() ([]byte, error) {
	type Properties AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption
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
func (r *AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption) UnmarshalJSON(b []byte) error {
	type Properties AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption
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
		*r = AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption(*res.Properties)
	}

	return nil
}

// GetAllAWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOptionResources retrieves all AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption items from an AWS CloudFormation template
func (t *Template) GetAllAWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOptionResources() map[string]AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption {
	results := map[string]AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::KinesisAnalyticsV2::ApplicationCloudWatchLoggingOption" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption{}
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

// GetAWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOptionWithName retrieves all AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOptionWithName(name string) (AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::KinesisAnalyticsV2::ApplicationCloudWatchLoggingOption" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSKinesisAnalyticsV2ApplicationCloudWatchLoggingOption{}, errors.New("resource not found")
}
