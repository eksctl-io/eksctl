package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSKinesisAnalyticsV2ApplicationReferenceDataSource AWS CloudFormation Resource (AWS::KinesisAnalyticsV2::ApplicationReferenceDataSource)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-kinesisanalyticsv2-applicationreferencedatasource.html
type AWSKinesisAnalyticsV2ApplicationReferenceDataSource struct {

	// ApplicationName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-kinesisanalyticsv2-applicationreferencedatasource.html#cfn-kinesisanalyticsv2-applicationreferencedatasource-applicationname
	ApplicationName *Value `json:"ApplicationName,omitempty"`

	// ReferenceDataSource AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-kinesisanalyticsv2-applicationreferencedatasource.html#cfn-kinesisanalyticsv2-applicationreferencedatasource-referencedatasource
	ReferenceDataSource *AWSKinesisAnalyticsV2ApplicationReferenceDataSource_ReferenceDataSource `json:"ReferenceDataSource,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisAnalyticsV2ApplicationReferenceDataSource) AWSCloudFormationType() string {
	return "AWS::KinesisAnalyticsV2::ApplicationReferenceDataSource"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSKinesisAnalyticsV2ApplicationReferenceDataSource) MarshalJSON() ([]byte, error) {
	type Properties AWSKinesisAnalyticsV2ApplicationReferenceDataSource
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
func (r *AWSKinesisAnalyticsV2ApplicationReferenceDataSource) UnmarshalJSON(b []byte) error {
	type Properties AWSKinesisAnalyticsV2ApplicationReferenceDataSource
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
		*r = AWSKinesisAnalyticsV2ApplicationReferenceDataSource(*res.Properties)
	}

	return nil
}

// GetAllAWSKinesisAnalyticsV2ApplicationReferenceDataSourceResources retrieves all AWSKinesisAnalyticsV2ApplicationReferenceDataSource items from an AWS CloudFormation template
func (t *Template) GetAllAWSKinesisAnalyticsV2ApplicationReferenceDataSourceResources() map[string]AWSKinesisAnalyticsV2ApplicationReferenceDataSource {
	results := map[string]AWSKinesisAnalyticsV2ApplicationReferenceDataSource{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSKinesisAnalyticsV2ApplicationReferenceDataSource:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::KinesisAnalyticsV2::ApplicationReferenceDataSource" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSKinesisAnalyticsV2ApplicationReferenceDataSource{}
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

// GetAWSKinesisAnalyticsV2ApplicationReferenceDataSourceWithName retrieves all AWSKinesisAnalyticsV2ApplicationReferenceDataSource items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSKinesisAnalyticsV2ApplicationReferenceDataSourceWithName(name string) (AWSKinesisAnalyticsV2ApplicationReferenceDataSource, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSKinesisAnalyticsV2ApplicationReferenceDataSource:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::KinesisAnalyticsV2::ApplicationReferenceDataSource" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSKinesisAnalyticsV2ApplicationReferenceDataSource{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSKinesisAnalyticsV2ApplicationReferenceDataSource{}, errors.New("resource not found")
}
