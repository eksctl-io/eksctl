package testutils

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	awseks "github.com/aws/aws-sdk-go/service/eks"
)

// NewFakeCluster creates a new fake cluster to be used in the tests
func NewFakeCluster(clusterName string, status string) *awseks.Cluster {
	created := &time.Time{}

	cluster := &awseks.Cluster{
		Name:      aws.String(clusterName),
		Status:    aws.String(status),
		Arn:       aws.String("arn:aws:eks:us-west-2:12345:cluster/test-12345"),
		CreatedAt: created,
		ResourcesVpcConfig: &awseks.VpcConfigResponse{
			VpcId:     aws.String("vpc-1234"),
			SubnetIds: aws.StringSlice([]string{"sub1", "sub2"}),
		},
		CertificateAuthority: &awseks.Certificate{
			Data: aws.String("dGVzdAo="),
		},
		Endpoint: aws.String("https://localhost/"),
	}

	return cluster
}
