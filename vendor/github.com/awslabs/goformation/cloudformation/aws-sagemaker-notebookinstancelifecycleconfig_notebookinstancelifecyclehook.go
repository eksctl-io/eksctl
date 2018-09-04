package cloudformation

// AWSSageMakerNotebookInstanceLifecycleConfig_NotebookInstanceLifecycleHook AWS CloudFormation Resource (AWS::SageMaker::NotebookInstanceLifecycleConfig.NotebookInstanceLifecycleHook)
// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-notebookinstancelifecycleconfig-notebookinstancelifecyclehook.html
type AWSSageMakerNotebookInstanceLifecycleConfig_NotebookInstanceLifecycleHook struct {

	// Content AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-properties-sagemaker-notebookinstancelifecycleconfig-notebookinstancelifecyclehook.html#cfn-sagemaker-notebookinstancelifecycleconfig-notebookinstancelifecyclehook-content
	Content *StringIntrinsic `json:"Content,omitempty"`
}

// AWSCloudFormationType returns the AWS CloudFormation resource type
func (r *AWSSageMakerNotebookInstanceLifecycleConfig_NotebookInstanceLifecycleHook) AWSCloudFormationType() string {
	return "AWS::SageMaker::NotebookInstanceLifecycleConfig.NotebookInstanceLifecycleHook"
}
