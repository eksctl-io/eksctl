package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSIoTAnalyticsPipeline AWS CloudFormation Resource (AWS::IoTAnalytics::Pipeline)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iotanalytics-pipeline.html
type AWSIoTAnalyticsPipeline struct {

	// PipelineActivities AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iotanalytics-pipeline.html#cfn-iotanalytics-pipeline-pipelineactivities
	PipelineActivities []AWSIoTAnalyticsPipeline_Activity `json:"PipelineActivities,omitempty"`

	// PipelineName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iotanalytics-pipeline.html#cfn-iotanalytics-pipeline-pipelinename
	PipelineName *Value `json:"PipelineName,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-iotanalytics-pipeline.html#cfn-iotanalytics-pipeline-tags
	Tags []Tag `json:"Tags,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSIoTAnalyticsPipeline) AWSCloudFormationType() string {
	return "AWS::IoTAnalytics::Pipeline"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSIoTAnalyticsPipeline) MarshalJSON() ([]byte, error) {
	type Properties AWSIoTAnalyticsPipeline
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
func (r *AWSIoTAnalyticsPipeline) UnmarshalJSON(b []byte) error {
	type Properties AWSIoTAnalyticsPipeline
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
		*r = AWSIoTAnalyticsPipeline(*res.Properties)
	}

	return nil
}

// GetAllAWSIoTAnalyticsPipelineResources retrieves all AWSIoTAnalyticsPipeline items from an AWS CloudFormation template
func (t *Template) GetAllAWSIoTAnalyticsPipelineResources() map[string]AWSIoTAnalyticsPipeline {
	results := map[string]AWSIoTAnalyticsPipeline{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSIoTAnalyticsPipeline:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::IoTAnalytics::Pipeline" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSIoTAnalyticsPipeline{}
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

// GetAWSIoTAnalyticsPipelineWithName retrieves all AWSIoTAnalyticsPipeline items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSIoTAnalyticsPipelineWithName(name string) (AWSIoTAnalyticsPipeline, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSIoTAnalyticsPipeline:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::IoTAnalytics::Pipeline" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSIoTAnalyticsPipeline{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSIoTAnalyticsPipeline{}, errors.New("resource not found")
}
