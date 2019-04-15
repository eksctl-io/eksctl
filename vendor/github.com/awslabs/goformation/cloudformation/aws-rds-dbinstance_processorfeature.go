package cloudformation

import (
	"encoding/json"
)

// AWSRDSDBInstance_ProcessorFeature AWS CloudFormation Resource (AWS::RDS::DBInstance.ProcessorFeature)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-rds-dbinstance-processorfeature.html
type AWSRDSDBInstance_ProcessorFeature struct {

	// Name AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-rds-dbinstance-processorfeature.html#cfn-rds-dbinstance-processorfeature-name
	Name *Value `json:"Name,omitempty"`

	// Value AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-rds-dbinstance-processorfeature.html#cfn-rds-dbinstance-processorfeature-value
	Value *Value `json:"Value,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSRDSDBInstance_ProcessorFeature) AWSCloudFormationType() string {
	return "AWS::RDS::DBInstance.ProcessorFeature"
}

func (r *AWSRDSDBInstance_ProcessorFeature) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
