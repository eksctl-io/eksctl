package aws

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awseks "github.com/aws/aws-sdk-go/service/eks"
)

// EksClusterExists checks if an EKS cluster exists in AWS
func EksClusterExists(clusterName string, session *session.Session) (bool, error) {
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

// EksClusterDelete deletes a EKS cluster with a given name
func EksClusterDelete(clusterName string, session *session.Session) error {
	eks := awseks.New(session)

	input := &awseks.DeleteClusterInput{
		Name: aws.String(clusterName),
	}

	_, err := eks.DeleteCluster(input)
	return err
}
