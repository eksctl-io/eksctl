package matchers

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
)

const (
	errorMessageTemplate = "Stack with id %s does not exist"
)

// stackExists checks to see if a CloudFormation stack exists
func stackExists(stackName string, cfg aws.Config) (bool, error) {
	cfn := cloudformation.NewFromConfig(cfg)

	input := &cloudformation.ListStackResourcesInput{
		StackName: aws.String(stackName),
	}
	_, err := cfn.ListStackResources(context.Background(), input)

	if err != nil {
		// Check if its a not found error
		errorMessage := fmt.Sprintf(errorMessageTemplate, stackName)
		if !strings.Contains(err.Error(), errorMessage) {
			return false, err
		}

		return false, nil
	}

	return true, nil
}
