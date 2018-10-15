package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awseks "github.com/aws/aws-sdk-go/service/eks"
)

// ClusterExists checks if an EKS cluster exists in AWS
func ClusterExists(clusterName string, session *session.Session) (bool, error) {
	eks := awseks.New(session)

	input := &awseks.DescribeClusterInput{
		Name: aws.String(clusterName),
	}
	_, err := eks.DescribeCluster(input)

	if err != nil {
		// Check if its a not found error: ResourceNotFoundException
		if !strings.Contains(err.Error(), awseks.ErrCodeResourceNotFoundException) {
			return false, err
		}

		return false, nil
	}

	return true, nil
}
