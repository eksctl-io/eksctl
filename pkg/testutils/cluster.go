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
		Arn:       aws.String("arn-12345678"),
		CreatedAt: created,
		ResourcesVpcConfig: &awseks.VpcConfigResponse{
			VpcId:     aws.String("vpc-1234"),
			SubnetIds: []*string{aws.String("sub1"), aws.String("sub2")},
		},
	}

	return cluster
}
