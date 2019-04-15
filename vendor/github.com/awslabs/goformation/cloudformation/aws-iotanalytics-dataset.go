package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSIoTAnalyticsDataset AWS CloudFormation Resource (AWS::IoTAnalytics::Dataset)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iotanalytics-dataset.html
type AWSIoTAnalyticsDataset struct {

	// Actions AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iotanalytics-dataset.html#cfn-iotanalytics-dataset-actions
	Actions []AWSIoTAnalyticsDataset_Action `json:"Actions,omitempty"`

	// DatasetName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iotanalytics-dataset.html#cfn-iotanalytics-dataset-datasetname
	DatasetName *Value `json:"DatasetName,omitempty"`

	// RetentionPeriod AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iotanalytics-dataset.html#cfn-iotanalytics-dataset-retentionperiod
	RetentionPeriod *AWSIoTAnalyticsDataset_RetentionPeriod `json:"RetentionPeriod,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iotanalytics-dataset.html#cfn-iotanalytics-dataset-tags
	Tags []Tag `json:"Tags,omitempty"`

	// Triggers AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iotanalytics-dataset.html#cfn-iotanalytics-dataset-triggers
	Triggers []AWSIoTAnalyticsDataset_Trigger `json:"Triggers,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsDataset) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Dataset"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSIoTAnalyticsDataset) MarshalJSON() ([]byte, error) {
	type Properties AWSIoTAnalyticsDataset
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
func (r *AWSIoTAnalyticsDataset) UnmarshalJSON(b []byte) error {
	type Properties AWSIoTAnalyticsDataset
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
		*r = AWSIoTAnalyticsDataset(*res.Properties)
	}

	return nil
}

// GetAllAWSIoTAnalyticsDatasetResources retrieves all AWSIoTAnalyticsDataset items from an AWS CloudFormation template
func (t *Template) GetAllAWSIoTAnalyticsDatasetResources() map[string]AWSIoTAnalyticsDataset {
	results := map[string]AWSIoTAnalyticsDataset{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSIoTAnalyticsDataset:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::IoTAnalytics::Dataset" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSIoTAnalyticsDataset{}
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

// GetAWSIoTAnalyticsDatasetWithName retrieves all AWSIoTAnalyticsDataset items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSIoTAnalyticsDatasetWithName(name string) (AWSIoTAnalyticsDataset, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSIoTAnalyticsDataset:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::IoTAnalytics::Dataset" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSIoTAnalyticsDataset{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSIoTAnalyticsDataset{}, errors.New("resource not found")
}
