package cloudformation

import (
	"encoding/json"
)

// AWSOpsWorksLayer_Recipes AWS CloudFormation Resource (AWS::OpsWorks::Layer.Recipes)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-layer-recipes.html
type AWSOpsWorksLayer_Recipes struct {

	// Configure AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-layer-recipes.html#cfn-opsworks-layer-customrecipes-configure
	Configure []*Value `json:"Configure,omitempty"`

	// Deploy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-layer-recipes.html#cfn-opsworks-layer-customrecipes-deploy
	Deploy []*Value `json:"Deploy,omitempty"`

	// Setup AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-layer-recipes.html#cfn-opsworks-layer-customrecipes-setup
	Setup []*Value `json:"Setup,omitempty"`

	// Shutdown AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-layer-recipes.html#cfn-opsworks-layer-customrecipes-shutdown
	Shutdown []*Value `json:"Shutdown,omitempty"`

	// Undeploy AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-opsworks-layer-recipes.html#cfn-opsworks-layer-customrecipes-undeploy
	Undeploy []*Value `json:"Undeploy,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSOpsWorksLayer_Recipes) AWSCloudFormationType() string {
	return "AWS::OpsWorks::Layer.Recipes"
}

func (r *AWSOpsWorksLayer_Recipes) MarshalJSON() ([]byte, error) {
	return json.Marshal(*r)
}
