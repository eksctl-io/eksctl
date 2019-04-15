package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSKinesisStreamConsumer AWS CloudFormation Resource (AWS::Kinesis::StreamConsumer)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-kinesis-streamconsumer.html
type AWSKinesisStreamConsumer struct {

	// ConsumerName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-kinesis-streamconsumer.html#cfn-kinesis-streamconsumer-consumername
	ConsumerName *Value `json:"ConsumerName,omitempty"`

	// StreamARN AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-kinesis-streamconsumer.html#cfn-kinesis-streamconsumer-streamarn
	StreamARN *Value `json:"StreamARN,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSKinesisStreamConsumer) AWSCloudFormationType() string {
	return "AWS::Kinesis::StreamConsumer"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSKinesisStreamConsumer) MarshalJSON() ([]byte, error) {
	type Properties AWSKinesisStreamConsumer
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
func (r *AWSKinesisStreamConsumer) UnmarshalJSON(b []byte) error {
	type Properties AWSKinesisStreamConsumer
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
		*r = AWSKinesisStreamConsumer(*res.Properties)
	}

	return nil
}

// GetAllAWSKinesisStreamConsumerResources retrieves all AWSKinesisStreamConsumer items from an AWS CloudFormation template
func (t *Template) GetAllAWSKinesisStreamConsumerResources() map[string]AWSKinesisStreamConsumer {
	results := map[string]AWSKinesisStreamConsumer{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSKinesisStreamConsumer:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Kinesis::StreamConsumer" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSKinesisStreamConsumer{}
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

// GetAWSKinesisStreamConsumerWithName retrieves all AWSKinesisStreamConsumer items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSKinesisStreamConsumerWithName(name string) (AWSKinesisStreamConsumer, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSKinesisStreamConsumer:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::Kinesis::StreamConsumer" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSKinesisStreamConsumer{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSKinesisStreamConsumer{}, errors.New("resource not found")
}
