package manager

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

// GetStackTemplate gets the Cloudformation template for a stack
func (c *StackCollection) GetStackTemplate(stackName string) (string, error) {
	input := &cloudformation.GetTemplateInput{
		StackName: aws.String(stackName),
	}

	output, err := c.provider.CloudFormation().GetTemplate(input)
	if err != nil {
		return "", err
	}

	return *output.TemplateBody, nil
}
