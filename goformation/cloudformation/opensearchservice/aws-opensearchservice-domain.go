package opensearchservice

import (
	"bytes"
	"encoding/json"
	"fmt"

	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/cloudformation"
	"goformation/v4/cloudformation/policies"
)

// Domain AWS CloudFormation Resource (AWS::OpenSearchService::Domain)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opensearchservice-domain.html
type Domain struct {

	// AccessPolicies AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opensearchservice-domain.html#cfn-opensearchservice-domain-accesspolicies
	AccessPolicies interface{} `json:"AccessPolicies,omitempty"`

	// AdvancedOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opensearchservice-domain.html#cfn-opensearchservice-domain-advancedoptions
	AdvancedOptions map[string]*types.Value `json:"AdvancedOptions,omitempty"`

	// AdvancedSecurityOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opensearchservice-domain.html#cfn-opensearchservice-domain-advancedsecurityoptions
	AdvancedSecurityOptions *Domain_AdvancedSecurityOptionsInput `json:"AdvancedSecurityOptions,omitempty"`

	// ClusterConfig AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opensearchservice-domain.html#cfn-opensearchservice-domain-clusterconfig
	ClusterConfig *Domain_ClusterConfig `json:"ClusterConfig,omitempty"`

	// CognitoOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opensearchservice-domain.html#cfn-opensearchservice-domain-cognitooptions
	CognitoOptions *Domain_CognitoOptions `json:"CognitoOptions,omitempty"`

	// DomainEndpointOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opensearchservice-domain.html#cfn-opensearchservice-domain-domainendpointoptions
	DomainEndpointOptions *Domain_DomainEndpointOptions `json:"DomainEndpointOptions,omitempty"`

	// DomainName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opensearchservice-domain.html#cfn-opensearchservice-domain-domainname
	DomainName *types.Value `json:"DomainName,omitempty"`

	// EBSOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opensearchservice-domain.html#cfn-opensearchservice-domain-ebsoptions
	EBSOptions *Domain_EBSOptions `json:"EBSOptions,omitempty"`

	// EncryptionAtRestOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opensearchservice-domain.html#cfn-opensearchservice-domain-encryptionatrestoptions
	EncryptionAtRestOptions *Domain_EncryptionAtRestOptions `json:"EncryptionAtRestOptions,omitempty"`

	// EngineVersion AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opensearchservice-domain.html#cfn-opensearchservice-domain-engineversion
	EngineVersion *types.Value `json:"EngineVersion,omitempty"`

	// LogPublishingOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opensearchservice-domain.html#cfn-opensearchservice-domain-logpublishingoptions
	LogPublishingOptions map[string]Domain_LogPublishingOption `json:"LogPublishingOptions,omitempty"`

	// NodeToNodeEncryptionOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opensearchservice-domain.html#cfn-opensearchservice-domain-nodetonodeencryptionoptions
	NodeToNodeEncryptionOptions *Domain_NodeToNodeEncryptionOptions `json:"NodeToNodeEncryptionOptions,omitempty"`

	// SnapshotOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opensearchservice-domain.html#cfn-opensearchservice-domain-snapshotoptions
	SnapshotOptions *Domain_SnapshotOptions `json:"SnapshotOptions,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opensearchservice-domain.html#cfn-opensearchservice-domain-tags
	Tags []cloudformation.Tag `json:"Tags,omitempty"`

	// VPCOptions AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-opensearchservice-domain.html#cfn-opensearchservice-domain-vpcoptions
	VPCOptions *Domain_VPCOptions `json:"VPCOptions,omitempty"`

	// AWSCloudFormationDeletionPolicy represents a CloudFormation DeletionPolicy
	AWSCloudFormationDeletionPolicy policies.DeletionPolicy `json:"-"`

	// AWSCloudFormationUpdateReplacePolicy represents a CloudFormation UpdateReplacePolicy
	AWSCloudFormationUpdateReplacePolicy policies.UpdateReplacePolicy `json:"-"`

	// AWSCloudFormationDependsOn stores the logical ID of the resources to be created before this resource
	AWSCloudFormationDependsOn []string `json:"-"`

	// AWSCloudFormationMetadata stores structured data associated with this resource
	AWSCloudFormationMetadata map[string]interface{} `json:"-"`

	// AWSCloudFormationCondition stores the logical ID of the condition that must be satisfied for this resource to be created
	AWSCloudFormationCondition string `json:"-"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *Domain) AWSCloudFormationType() string {
	return "AWS::OpenSearchService::Domain"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r Domain) MarshalJSON() ([]byte, error) {
	type Properties Domain
	return json.Marshal(&struct {
		Type                string
		Properties          Properties
		DependsOn           []string                     `json:"DependsOn,omitempty"`
		Metadata            map[string]interface{}       `json:"Metadata,omitempty"`
		DeletionPolicy      policies.DeletionPolicy      `json:"DeletionPolicy,omitempty"`
		UpdateReplacePolicy policies.UpdateReplacePolicy `json:"UpdateReplacePolicy,omitempty"`
		Condition           string                       `json:"Condition,omitempty"`
	}{
		Type:                r.AWSCloudFormationType(),
		Properties:          (Properties)(r),
		DependsOn:           r.AWSCloudFormationDependsOn,
		Metadata:            r.AWSCloudFormationMetadata,
		DeletionPolicy:      r.AWSCloudFormationDeletionPolicy,
		UpdateReplacePolicy: r.AWSCloudFormationUpdateReplacePolicy,
		Condition:           r.AWSCloudFormationCondition,
	})
}

// UnmarshalJSON is a custom JSON unmarshalling hook that strips the outer
// AWS CloudFormation resource object, and just keeps the 'Properties' field.
func (r *Domain) UnmarshalJSON(b []byte) error {
	type Properties Domain
	res := &struct {
		Type                string
		Properties          *Properties
		DependsOn           []string
		Metadata            map[string]interface{}
		DeletionPolicy      string
		UpdateReplacePolicy string
		Condition           string
	}{}

	dec := json.NewDecoder(bytes.NewReader(b))
	dec.DisallowUnknownFields() // Force error if unknown field is found

	if err := dec.Decode(&res); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return err
	}

	// If the resource has no Properties set, it could be nil
	if res.Properties != nil {
		*r = Domain(*res.Properties)
	}
	if res.DependsOn != nil {
		r.AWSCloudFormationDependsOn = res.DependsOn
	}
	if res.Metadata != nil {
		r.AWSCloudFormationMetadata = res.Metadata
	}
	if res.DeletionPolicy != "" {
		r.AWSCloudFormationDeletionPolicy = policies.DeletionPolicy(res.DeletionPolicy)
	}
	if res.UpdateReplacePolicy != "" {
		r.AWSCloudFormationUpdateReplacePolicy = policies.UpdateReplacePolicy(res.UpdateReplacePolicy)
	}
	if res.Condition != "" {
		r.AWSCloudFormationCondition = res.Condition
	}
	return nil
}
