package manager

import (
	"context"
	"fmt"

	"goformation/v4"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

// GetStackTemplate gets the Cloudformation template for a stack
// and returns a json string representation
func (c *StackCollection) GetStackTemplate(ctx context.Context, stackName string) (string, error) {
	input := &cloudformation.GetTemplateInput{
		StackName: aws.String(stackName),
	}

	output, err := c.cloudformationAPI.GetTemplate(ctx, input)
	if err != nil {
		return "", err
	}

	return ensureJSONResponse([]byte(*output.TemplateBody))
}

func ensureJSONResponse(templateBody []byte) (string, error) {
	//since json is valid yaml we just need to check the response is valid yaml
	template, err := goformation.ParseYAML(templateBody)
	if err != nil {
		return "", fmt.Errorf("failed to parse GetStackTemplate response: %w", err)
	}
	bytes, err := template.JSON()
	return string(bytes), err
}
