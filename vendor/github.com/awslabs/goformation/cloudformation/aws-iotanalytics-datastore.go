package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSIoTAnalyticsDatastore AWS CloudFormation Resource (AWS::IoTAnalytics::Datastore)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iotanalytics-datastore.html
type AWSIoTAnalyticsDatastore struct {

	// DatastoreName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iotanalytics-datastore.html#cfn-iotanalytics-datastore-datastorename
	DatastoreName *Value `json:"DatastoreName,omitempty"`

	// RetentionPeriod AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iotanalytics-datastore.html#cfn-iotanalytics-datastore-retentionperiod
	RetentionPeriod *AWSIoTAnalyticsDatastore_RetentionPeriod `json:"RetentionPeriod,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iotanalytics-datastore.html#cfn-iotanalytics-datastore-tags
	Tags []Tag `json:"Tags,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsDatastore) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Datastore"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSIoTAnalyticsDatastore) MarshalJSON() ([]byte, error) {
	type Properties AWSIoTAnalyticsDatastore
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
func (r *AWSIoTAnalyticsDatastore) UnmarshalJSON(b []byte) error {
	type Properties AWSIoTAnalyticsDatastore
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
		*r = AWSIoTAnalyticsDatastore(*res.Properties)
	}

	return nil
}

// GetAllAWSIoTAnalyticsDatastoreResources retrieves all AWSIoTAnalyticsDatastore items from an AWS CloudFormation template
func (t *Template) GetAllAWSIoTAnalyticsDatastoreResources() map[string]AWSIoTAnalyticsDatastore {
	results := map[string]AWSIoTAnalyticsDatastore{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSIoTAnalyticsDatastore:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::IoTAnalytics::Datastore" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSIoTAnalyticsDatastore{}
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

// GetAWSIoTAnalyticsDatastoreWithName retrieves all AWSIoTAnalyticsDatastore items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSIoTAnalyticsDatastoreWithName(name string) (AWSIoTAnalyticsDatastore, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSIoTAnalyticsDatastore:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::IoTAnalytics::Datastore" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSIoTAnalyticsDatastore{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSIoTAnalyticsDatastore{}, errors.New("resource not found")
}
