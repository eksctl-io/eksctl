package databrew

import (
	"bytes"
	"encoding/json"
	"fmt"

	"goformation/v4/cloudformation/types"

	"goformation/v4/cloudformation/cloudformation"
	"goformation/v4/cloudformation/policies"
)

// Job AWS CloudFormation Resource (AWS::DataBrew::Job)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-databrew-job.html
type Job struct {

	// DataCatalogOutputs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-databrew-job.html#cfn-databrew-job-datacatalogoutputs
	DataCatalogOutputs []Job_DataCatalogOutput `json:"DataCatalogOutputs,omitempty"`

	// DatabaseOutputs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-databrew-job.html#cfn-databrew-job-databaseoutputs
	DatabaseOutputs []Job_DatabaseOutput `json:"DatabaseOutputs,omitempty"`

	// DatasetName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-databrew-job.html#cfn-databrew-job-datasetname
	DatasetName *types.Value `json:"DatasetName,omitempty"`

	// EncryptionKeyArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-databrew-job.html#cfn-databrew-job-encryptionkeyarn
	EncryptionKeyArn *types.Value `json:"EncryptionKeyArn,omitempty"`

	// EncryptionMode AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-databrew-job.html#cfn-databrew-job-encryptionmode
	EncryptionMode *types.Value `json:"EncryptionMode,omitempty"`

	// JobSample AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-databrew-job.html#cfn-databrew-job-jobsample
	JobSample *Job_JobSample `json:"JobSample,omitempty"`

	// LogSubscription AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-databrew-job.html#cfn-databrew-job-logsubscription
	LogSubscription *types.Value `json:"LogSubscription,omitempty"`

	// MaxCapacity AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-databrew-job.html#cfn-databrew-job-maxcapacity
	MaxCapacity *types.Value `json:"MaxCapacity,omitempty"`

	// MaxRetries AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-databrew-job.html#cfn-databrew-job-maxretries
	MaxRetries *types.Value `json:"MaxRetries,omitempty"`

	// Name AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-databrew-job.html#cfn-databrew-job-name
	Name *types.Value `json:"Name,omitempty"`

	// OutputLocation AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-databrew-job.html#cfn-databrew-job-outputlocation
	OutputLocation *Job_OutputLocation `json:"OutputLocation,omitempty"`

	// Outputs AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-databrew-job.html#cfn-databrew-job-outputs
	Outputs []Job_Output `json:"Outputs,omitempty"`

	// ProfileConfiguration AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-databrew-job.html#cfn-databrew-job-profileconfiguration
	ProfileConfiguration *Job_ProfileConfiguration `json:"ProfileConfiguration,omitempty"`

	// ProjectName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-databrew-job.html#cfn-databrew-job-projectname
	ProjectName *types.Value `json:"ProjectName,omitempty"`

	// Recipe AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-databrew-job.html#cfn-databrew-job-recipe
	Recipe *Job_Recipe `json:"Recipe,omitempty"`

	// RoleArn AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-databrew-job.html#cfn-databrew-job-rolearn
	RoleArn *types.Value `json:"RoleArn,omitempty"`

	// Tags AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-databrew-job.html#cfn-databrew-job-tags
	Tags []cloudformation.Tag `json:"Tags,omitempty"`

	// Timeout AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-databrew-job.html#cfn-databrew-job-timeout
	Timeout *types.Value `json:"Timeout,omitempty"`

	// Type AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-databrew-job.html#cfn-databrew-job-type
	Type *types.Value `json:"Type,omitempty"`

	// ValidationConfigurations AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-databrew-job.html#cfn-databrew-job-validationconfigurations
	ValidationConfigurations []Job_ValidationConfiguration `json:"ValidationConfigurations,omitempty"`

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
func (r *Job) AWSCloudFormationType() string {
	return "AWS::DataBrew::Job"
}

// MarshalJSON is a custom JSON marshalling hook that embeds this object into
// an AWS CloudFormation JSON resource's 'Properties' field and adds a 'Type'.
func (r Job) MarshalJSON() ([]byte, error) {
	type Properties Job
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
func (r *Job) UnmarshalJSON(b []byte) error {
	type Properties Job
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
		*r = Job(*res.Properties)
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
