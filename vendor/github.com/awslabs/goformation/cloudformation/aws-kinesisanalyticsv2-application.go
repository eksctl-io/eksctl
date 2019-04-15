package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSKinesisAnalyticsV2Application AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::Application)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-kinesisanalyticsv2-application.html
type AWSKinesisAnalyticsV2Application struct {

	// ApplicationConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-kinesisanalyticsv2-application.html#cfn-kinesisanalyticsv2-application-applicationconfiguration
	ApplicationConfiguration *AWSKinesisAnalyticsV2Application_ApplicationConfiguration `json:"ApplicationConfiguration,omitempty"`

	// ApplicationDescription AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-kinesisanalyticsv2-application.html#cfn-kinesisanalyticsv2-application-applicationdescription
	ApplicationDescription *Value `json:"ApplicationDescription,omitempty"`

	// ApplicationName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-kinesisanalyticsv2-application.html#cfn-kinesisanalyticsv2-application-applicationname
	ApplicationName *Value `json:"ApplicationName,omitempty"`

	// RuntimeEnvironment AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-kinesisanalyticsv2-application.html#cfn-kinesisanalyticsv2-application-runtimeenvironment
	RuntimeEnvironment *Value `json:"RuntimeEnvironment,omitempty"`

	// ServiceExecutionRole AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-kinesisanalyticsv2-application.html#cfn-kinesisanalyticsv2-application-serviceexecutionrole
	ServiceExecutionRole *Value `json:"ServiceExecutionRole,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2Application) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::Application"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSKinesisAnalyticsV2Application) MarshalJSON() ([]byte, error) {
	type Properties AWSKinesisAnalyticsV2Application
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
func (r *AWSKinesisAnalyticsV2Application) UnmarshalJSON(b []byte) error {
	type Properties AWSKinesisAnalyticsV2Application
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
		*r = AWSKinesisAnalyticsV2Application(*res.Properties)
	}

	return nil
}

// GetAllAWSKinesisAnalyticsV2ApplicationResources retrieves all AWSKinesisAnalyticsV2Application items from an AWS CloudFormation template
func (t *Template) GetAllAWSKinesisAnalyticsV2ApplicationResources() map[string]AWSKinesisAnalyticsV2Application {
	results := map[string]AWSKinesisAnalyticsV2Application{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSKinesisAnalyticsV2Application:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::KinesisAnalyticsV2::Application" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSKinesisAnalyticsV2Application{}
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

// GetAWSKinesisAnalyticsV2ApplicationWithName retrieves all AWSKinesisAnalyticsV2Application items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSKinesisAnalyticsV2ApplicationWithName(name string) (AWSKinesisAnalyticsV2Application, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSKinesisAnalyticsV2Application:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::KinesisAnalyticsV2::Application" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSKinesisAnalyticsV2Application{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSKinesisAnalyticsV2Application{}, errors.New("resource not found")
}
