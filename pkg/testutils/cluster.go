package testutils

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
)

// NewFakeCluster creates a new fake cluster to be used in the tests
func NewFakeCluster(clusterName string, status ekstypes.ClusterStatus) *ekstypes.Cluster {
	created := &time.Time{}

	cluster := &ekstypes.Cluster{
		Name:      aws.String(clusterName),
		Status:    status,
		Arn:       aws.String("arn:aws:eks:us-west-2:12345:cluster/test-12345"),
		CreatedAt: created,
		ResourcesVpcConfig: &ekstypes.VpcConfigResponse{
			VpcId:     aws.String("vpc-1234"),
			SubnetIds: []string{"sub1", "sub2"},
		},
		CertificateAuthority: &ekstypes.Certificate{
			Data: aws.String("dGVzdAo="),
		},
		Endpoint: aws.String("https://localhost/"),
		Version:  aws.String("1.17"),
	}

	return cluster
}
