package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSSSMResourceDataSync AWS CloudFormation Resource (AWS::SSM::ResourceDataSync)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ssm-resourcedatasync.html
type AWSSSMResourceDataSync struct {

	// BucketName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ssm-resourcedatasync.html#cfn-ssm-resourcedatasync-bucketname
	BucketName *StringIntrinsic `json:"BucketName,omitempty"`

	// BucketPrefix AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ssm-resourcedatasync.html#cfn-ssm-resourcedatasync-bucketprefix
	BucketPrefix *StringIntrinsic `json:"BucketPrefix,omitempty"`

	// BucketRegion AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ssm-resourcedatasync.html#cfn-ssm-resourcedatasync-bucketregion
	BucketRegion *StringIntrinsic `json:"BucketRegion,omitempty"`

	// KMSKeyArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ssm-resourcedatasync.html#cfn-ssm-resourcedatasync-kmskeyarn
	KMSKeyArn *StringIntrinsic `json:"KMSKeyArn,omitempty"`

	// SyncFormat AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ssm-resourcedatasync.html#cfn-ssm-resourcedatasync-syncformat
	SyncFormat *StringIntrinsic `json:"SyncFormat,omitempty"`

	// SyncName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ssm-resourcedatasync.html#cfn-ssm-resourcedatasync-syncname
	SyncName *StringIntrinsic `json:"SyncName,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSSMResourceDataSync) AWSCloudFormationType() string {
	return "AWS::SSM::ResourceDataSync"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSSSMResourceDataSync) MarshalJSON() ([]byte, error) {
	type Properties AWSSSMResourceDataSync
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
func (r *AWSSSMResourceDataSync) UnmarshalJSON(b []byte) error {
	type Properties AWSSSMResourceDataSync
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
		*r = AWSSSMResourceDataSync(*res.Properties)
	}

	return nil
}

// GetAllAWSSSMResourceDataSyncResources retrieves all AWSSSMResourceDataSync items from an AWS CloudFormation template
func (t *Template) GetAllAWSSSMResourceDataSyncResources() map[string]AWSSSMResourceDataSync {
	results := map[string]AWSSSMResourceDataSync{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSSSMResourceDataSync:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::SSM::ResourceDataSync" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						var result AWSSSMResourceDataSync
						if err := json.Unmarshal(b, &result); err == nil {
							results[name] = result
						}
					}
				}
			}
		}
	}
	return results
}

// GetAWSSSMResourceDataSyncWithName retrieves all AWSSSMResourceDataSync items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSSSMResourceDataSyncWithName(name string) (AWSSSMResourceDataSync, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSSSMResourceDataSync:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::SSM::ResourceDataSync" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						var result AWSSSMResourceDataSync
						if err := json.Unmarshal(b, &result); err == nil {
							return result, nil
						}
					}
				}
			}
		}
	}
	return AWSSSMResourceDataSync{}, errors.New("resource not found")
}
