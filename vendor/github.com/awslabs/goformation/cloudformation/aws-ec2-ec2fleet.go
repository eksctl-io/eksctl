package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSEC2EC2Fleet AWS CloudFormation Resource (AWS::EC2::EC2Fleet)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-ec2fleet.html
type AWSEC2EC2Fleet struct {

	// ExcessCapacityTerminationPolicy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-ec2fleet.html#cfn-ec2-ec2fleet-excesscapacityterminationpolicy
	ExcessCapacityTerminationPolicy *Value `json:"ExcessCapacityTerminationPolicy,omitempty"`

	// LaunchTemplateConfigs AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-ec2fleet.html#cfn-ec2-ec2fleet-launchtemplateconfigs
	LaunchTemplateConfigs []AWSEC2EC2Fleet_FleetLaunchTemplateConfigRequest `json:"LaunchTemplateConfigs,omitempty"`

	// OnDemandOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-ec2fleet.html#cfn-ec2-ec2fleet-ondemandoptions
	OnDemandOptions *AWSEC2EC2Fleet_OnDemandOptionsRequest `json:"OnDemandOptions,omitempty"`

	// ReplaceUnhealthyInstances AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-ec2fleet.html#cfn-ec2-ec2fleet-replaceunhealthyinstances
	ReplaceUnhealthyInstances *Value `json:"ReplaceUnhealthyInstances,omitempty"`

	// SpotOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-ec2fleet.html#cfn-ec2-ec2fleet-spotoptions
	SpotOptions *AWSEC2EC2Fleet_SpotOptionsRequest `json:"SpotOptions,omitempty"`

	// TagSpecifications AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-ec2fleet.html#cfn-ec2-ec2fleet-tagspecifications
	TagSpecifications []AWSEC2EC2Fleet_TagSpecification `json:"TagSpecifications,omitempty"`

	// TargetCapacitySpecification AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-ec2fleet.html#cfn-ec2-ec2fleet-targetcapacityspecification
	TargetCapacitySpecification *AWSEC2EC2Fleet_TargetCapacitySpecificationRequest `json:"TargetCapacitySpecification,omitempty"`

	// TerminateInstancesWithExpiration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-ec2fleet.html#cfn-ec2-ec2fleet-terminateinstanceswithexpiration
	TerminateInstancesWithExpiration *Value `json:"TerminateInstancesWithExpiration,omitempty"`

	// Type AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-ec2fleet.html#cfn-ec2-ec2fleet-type
	Type *Value `json:"Type,omitempty"`

	// ValidFrom AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-ec2fleet.html#cfn-ec2-ec2fleet-validfrom
	ValidFrom *Value `json:"ValidFrom,omitempty"`

	// ValidUntil AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-ec2-ec2fleet.html#cfn-ec2-ec2fleet-validuntil
	ValidUntil *Value `json:"ValidUntil,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSEC2EC2Fleet) AWSCloudFormationType() string {
	return "AWS::EC2::EC2Fleet"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSEC2EC2Fleet) MarshalJSON() ([]byte, error) {
	type Properties AWSEC2EC2Fleet
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
func (r *AWSEC2EC2Fleet) UnmarshalJSON(b []byte) error {
	type Properties AWSEC2EC2Fleet
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
		*r = AWSEC2EC2Fleet(*res.Properties)
	}

	return nil
}

// GetAllAWSEC2EC2FleetResources retrieves all AWSEC2EC2Fleet items from an AWS CloudFormation template
func (t *Template) GetAllAWSEC2EC2FleetResources() map[string]AWSEC2EC2Fleet {
	results := map[string]AWSEC2EC2Fleet{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSEC2EC2Fleet:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::EC2::EC2Fleet" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSEC2EC2Fleet{}
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

// GetAWSEC2EC2FleetWithName retrieves all AWSEC2EC2Fleet items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSEC2EC2FleetWithName(name string) (AWSEC2EC2Fleet, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSEC2EC2Fleet:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::EC2::EC2Fleet" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						result := &AWSEC2EC2Fleet{}
						if err := result.UnmarshalJSON(b); err == nil {
							return *result, nil
						}
					}
				}
			}
		}
	}
	return AWSEC2EC2Fleet{}, errors.New("resource not found")
}
