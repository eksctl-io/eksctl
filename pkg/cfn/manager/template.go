package manager

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/pkg/errors"
	"github.com/weaveworks/goformation/v4"
)

// GetStackTemplate gets the Cloudformation template for a stack
// and returns a json string representation
func (c *StackCollection) GetStackTemplate(stackName string) (string, error) {
	input := &cloudformation.GetTemplateInput{
		StackName: aws.String(stackName),
	}

	output, err := c.provider.CloudFormation().GetTemplate(input)
	if err != nil {
		return "", err
	}

	return ensureJSONResponse([]byte(*output.TemplateBody))
}

func ensureJSONResponse(templateBody []byte) (string, error) {
	//since json is valid yaml we just need to check the response is valid yaml
	template, err := goformation.ParseYAML(templateBody)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse GetStackTemplate response")
	}
	bytes, err := template.JSON()
	return string(bytes), err
}
