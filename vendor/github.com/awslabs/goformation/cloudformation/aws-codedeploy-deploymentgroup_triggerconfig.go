package cloudformation

// AWSCodeDeployDeploymentGroup_TriggerConfig AWS CloudFormation Resource (AWS::CodeDeploy::DeploymentGroup.TriggerConfig)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-triggerconfig.html
type AWSCodeDeployDeploymentGroup_TriggerConfig struct {

	// TriggerEvents AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-triggerconfig.html#cfn-codedeploy-deploymentgroup-triggerconfig-triggerevents
	TriggerEvents []*StringIntrinsic `json:"TriggerEvents,omitempty"`

	// TriggerName AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-triggerconfig.html#cfn-codedeploy-deploymentgroup-triggerconfig-triggername
	TriggerName *StringIntrinsic `json:"TriggerName,omitempty"`

	// TriggerTargetArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-codedeploy-deploymentgroup-triggerconfig.html#cfn-codedeploy-deploymentgroup-triggerconfig-triggertargetarn
	TriggerTargetArn *StringIntrinsic `json:"TriggerTargetArn,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSCodeDeployDeploymentGroup_TriggerConfig) AWSCloudFormationType() string {
	return "AWS::CodeDeploy::DeploymentGroup.TriggerConfig"
}
