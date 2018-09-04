package cloudformation

import (
	"encoding/json"
	"errors"
	"fmt"
)

// AWSOpsWorksStack AWS CloudFormation Resource (AWS::OpsWorks::Stack)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html
type AWSOpsWorksStack struct {

	// AgentVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-agentversion
	AgentVersion *StringIntrinsic `json:"AgentVersion,omitempty"`

	// Attributes AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-attributes
	Attributes map[string]*StringIntrinsic `json:"Attributes,omitempty"`

	// ChefConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-chefconfiguration
	ChefConfiguration *AWSOpsWorksStack_ChefConfiguration `json:"ChefConfiguration,omitempty"`

	// CloneAppIds AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-cloneappids
	CloneAppIds []*StringIntrinsic `json:"CloneAppIds,omitempty"`

	// ClonePermissions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-clonepermissions
	ClonePermissions bool `json:"ClonePermissions,omitempty"`

	// ConfigurationManager AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-configmanager
	ConfigurationManager *AWSOpsWorksStack_StackConfigurationManager `json:"ConfigurationManager,omitempty"`

	// CustomCookbooksSource AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-custcookbooksource
	CustomCookbooksSource *AWSOpsWorksStack_Source `json:"CustomCookbooksSource,omitempty"`

	// CustomJson AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-custjson
	CustomJson interface{} `json:"CustomJson,omitempty"`

	// DefaultAvailabilityZone AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-defaultaz
	DefaultAvailabilityZone *StringIntrinsic `json:"DefaultAvailabilityZone,omitempty"`

	// DefaultInstanceProfileArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-defaultinstanceprof
	DefaultInstanceProfileArn *StringIntrinsic `json:"DefaultInstanceProfileArn,omitempty"`

	// DefaultOs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-defaultos
	DefaultOs *StringIntrinsic `json:"DefaultOs,omitempty"`

	// DefaultRootDeviceType AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-defaultrootdevicetype
	DefaultRootDeviceType *StringIntrinsic `json:"DefaultRootDeviceType,omitempty"`

	// DefaultSshKeyName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-defaultsshkeyname
	DefaultSshKeyName *StringIntrinsic `json:"DefaultSshKeyName,omitempty"`

	// DefaultSubnetId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#defaultsubnet
	DefaultSubnetId *StringIntrinsic `json:"DefaultSubnetId,omitempty"`

	// EcsClusterArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-ecsclusterarn
	EcsClusterArn *StringIntrinsic `json:"EcsClusterArn,omitempty"`

	// ElasticIps AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-elasticips
	ElasticIps []AWSOpsWorksStack_ElasticIp `json:"ElasticIps,omitempty"`

	// HostnameTheme AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-hostnametheme
	HostnameTheme *StringIntrinsic `json:"HostnameTheme,omitempty"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-name
	Name *StringIntrinsic `json:"Name,omitempty"`

	// RdsDbInstances AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-rdsdbinstances
	RdsDbInstances []AWSOpsWorksStack_RdsDbInstance `json:"RdsDbInstances,omitempty"`

	// ServiceRoleArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-servicerolearn
	ServiceRoleArn *StringIntrinsic `json:"ServiceRoleArn,omitempty"`

	// SourceStackId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-sourcestackid
	SourceStackId *StringIntrinsic `json:"SourceStackId,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-tags
	Tags []Tag `json:"Tags,omitempty"`

	// UseCustomCookbooks AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#usecustcookbooks
	UseCustomCookbooks bool `json:"UseCustomCookbooks,omitempty"`

	// UseOpsworksSecurityGroups AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-useopsworkssecuritygroups
	UseOpsworksSecurityGroups bool `json:"UseOpsworksSecurityGroups,omitempty"`

	// VpcId AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opsworks-stack.html#cfn-opsworks-stack-vpcid
	VpcId *StringIntrinsic `json:"VpcId,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSOpsWorksStack) AWSCloudFormationType() string {
	return "AWS::OpsWorks::Stack"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r *AWSOpsWorksStack) MarshalJSON() ([]byte, error) {
	type Properties AWSOpsWorksStack
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
func (r *AWSOpsWorksStack) UnmarshalJSON(b []byte) error {
	type Properties AWSOpsWorksStack
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
		*r = AWSOpsWorksStack(*res.Properties)
	}

	return nil
}

// GetAllAWSOpsWorksStackResources retrieves all AWSOpsWorksStack items from an AWS CloudFormation template
func (t *Template) GetAllAWSOpsWorksStackResources() map[string]AWSOpsWorksStack {
	results := map[string]AWSOpsWorksStack{}
	for name, untyped := range t.Resources {
		switch resource := untyped.(type) {
		case AWSOpsWorksStack:
			// We found a strongly typed resource of the correct type; use it
			results[name] = resource
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::OpsWorks::Stack" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						var result AWSOpsWorksStack
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

// GetAWSOpsWorksStackWithName retrieves all AWSOpsWorksStack items from an AWS CloudFormation template
// whose logical ID matches the provided name. Returns an error if not found.
func (t *Template) GetAWSOpsWorksStackWithName(name string) (AWSOpsWorksStack, error) {
	if untyped, ok := t.Resources[name]; ok {
		switch resource := untyped.(type) {
		case AWSOpsWorksStack:
			// We found a strongly typed resource of the correct type; use it
			return resource, nil
		case map[string]interface{}:
			// We found an untyped resource (likely from JSON) which *might* be
			// the correct type, but we need to check it's 'Type' field
			if resType, ok := resource["Type"]; ok {
				if resType == "AWS::OpsWorks::Stack" {
					// The resource is correct, unmarshal it into the results
					if b, err := json.Marshal(resource); err == nil {
						var result AWSOpsWorksStack
						if err := json.Unmarshal(b, &result); err == nil {
							return result, nil
						}
					}
				}
			}
		}
	}
	return AWSOpsWorksStack{}, errors.New("resource not found")
}
