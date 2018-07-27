package printers_test

import (
	"os"
	"testing"
	"time"
	//"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	. "github.com/weaveworks/eksctl/pkg/printers"
)

func TestTablePrinter(t *testing.T) {
	created := time.Now()
	test := &eks.DescribeClusterOutput{
		Cluster: &eks.Cluster{
			Name:      aws.String("test-cluster"),
			Status:    aws.String(eks.ClusterStatusActive),
			Arn:       aws.String("arn-12345678"),
			CreatedAt: &created,
			ResourcesVpcConfig: &eks.VpcConfigResponse {
				VpcId: aws.String("vpc-1234"),
				SubnetIds: []*string{aws.String("sub1"), aws.String("sub2")},
			},
		},
	}
	printer := NewTablePrinter()

	printer.(*TablePrinter).AddColumn("CLUSTERNAME", func(c *eks.Cluster) string {
		return *c.Name
	})
	printer.(*TablePrinter).AddColumn("ARN", func(c *eks.Cluster) string {
		return *c.Arn
	})

	printer.PrintObj([]*eks.Cluster{test.Cluster}, os.Stdout)
}


func TestJsonPrinter(t *testing.T) {
	created := time.Now()
	test := &eks.DescribeClusterOutput{
		Cluster: &eks.Cluster{
			Name:      aws.String("test-cluster"),
			Status:    aws.String(eks.ClusterStatusActive),
			Arn:       aws.String("arn-12345678"),
			CreatedAt: &created,
			ResourcesVpcConfig: &eks.VpcConfigResponse {
				VpcId: aws.String("vpc-1234"),
				SubnetIds: []*string{aws.String("sub1"), aws.String("sub2")},
			},
		},
	}

	printer := NewJSONPrinter()
	printer.PrintObj([]*eks.Cluster{test.Cluster}, os.Stdout)
}

func TestYamlPrinter(t *testing.T) {
	created := time.Now()
	test := &eks.DescribeClusterOutput{
		Cluster: &eks.Cluster{
			Name:      aws.String("test-cluster"),
			Status:    aws.String(eks.ClusterStatusActive),
			Arn:       aws.String("arn-12345678"),
			CreatedAt: &created,
			ResourcesVpcConfig: &eks.VpcConfigResponse {
				VpcId: aws.String("vpc-1234"),
				SubnetIds: []*string{aws.String("sub1"), aws.String("sub2")},
			},
		},
	}

	printer := NewYAMLPrinter()
	printer.PrintObj([]*eks.Cluster{test.Cluster}, os.Stdout)
}
