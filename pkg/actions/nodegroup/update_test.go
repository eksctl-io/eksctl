package nodegroup

import (
	"context"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Update", func() {
	var (
		clusterName, ngName string
		p                   *mockprovider.MockProvider
		cfg                 *api.ClusterConfig
		m                   *Manager
	)

	BeforeEach(func() {
		clusterName = "my-cluster"
		ngName = "my-ng"
		p = mockprovider.NewMockProvider()
		cfg = api.NewClusterConfig()
		cfg.Metadata.Name = clusterName
		cfg.ManagedNodeGroups = []*api.ManagedNodeGroup{
			{
				NodeGroupBase: &api.NodeGroupBase{
					Name: ngName,
				},
			},
		}
	})

	It("fails for unmanaged nodegroups", func() {
		p.MockEKS().On("DescribeNodegroup", mock.Anything, &awseks.DescribeNodegroupInput{
			ClusterName:   &clusterName,
			NodegroupName: &ngName,
		}).Return(nil, &ekstypes.ResourceNotFoundException{})

		m = New(cfg, &eks.ClusterProvider{AWSProvider: p}, nil, nil)
		err := m.Update(context.Background(), false)
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring("could not find managed nodegroup with name \"my-ng\"")))
	})

	It("[happy path] successfully updates a nodegroup with updateConfig and maxUnavailable", func() {
		p.MockEKS().On("DescribeNodegroup", mock.Anything, &awseks.DescribeNodegroupInput{
			ClusterName:   &clusterName,
			NodegroupName: &ngName,
		}).Return(&awseks.DescribeNodegroupOutput{
			Nodegroup: &ekstypes.Nodegroup{
				UpdateConfig: &ekstypes.NodegroupUpdateConfig{
					MaxUnavailable: aws.Int32(4),
				},
			},
		}, nil)

		p.MockEKS().On("UpdateNodegroupConfig", mock.Anything, &awseks.UpdateNodegroupConfigInput{
			UpdateConfig: &ekstypes.NodegroupUpdateConfig{
				MaxUnavailable: aws.Int32(6),
			},
			ClusterName:   &clusterName,
			NodegroupName: &ngName,
		}).Return(nil, nil)

		cfg.ManagedNodeGroups[0].UpdateConfig = &api.NodeGroupUpdateConfig{
			MaxUnavailable: aws.Int(6),
		}

		m = New(cfg, &eks.ClusterProvider{AWSProvider: p}, nil, nil)
		err := m.Update(context.Background(), false)
		Expect(err).NotTo(HaveOccurred())
	})

	It("[happy path] successfully updates multiple nodegroups with updateConfig and maxUnavailable", func() {
		p.MockEKS().On("DescribeNodegroup", mock.Anything, &awseks.DescribeNodegroupInput{
			ClusterName:   &clusterName,
			NodegroupName: &ngName,
		}).Return(&awseks.DescribeNodegroupOutput{
			Nodegroup: &ekstypes.Nodegroup{
				UpdateConfig: &ekstypes.NodegroupUpdateConfig{
					MaxUnavailable: aws.Int32(4),
				},
			},
		}, nil)

		p.MockEKS().On("DescribeNodegroup", mock.Anything, &awseks.DescribeNodegroupInput{
			ClusterName:   &clusterName,
			NodegroupName: aws.String("ng-2"),
		}).Return(&awseks.DescribeNodegroupOutput{
			Nodegroup: &ekstypes.Nodegroup{
				UpdateConfig: &ekstypes.NodegroupUpdateConfig{
					MaxUnavailable: aws.Int32(4),
				},
			},
		}, nil)

		p.MockEKS().On("UpdateNodegroupConfig", mock.Anything, &awseks.UpdateNodegroupConfigInput{
			UpdateConfig: &ekstypes.NodegroupUpdateConfig{
				MaxUnavailable: aws.Int32(6),
			},
			ClusterName:   &clusterName,
			NodegroupName: &ngName,
		}).Return(nil, nil)

		p.MockEKS().On("UpdateNodegroupConfig", mock.Anything, &awseks.UpdateNodegroupConfigInput{
			UpdateConfig: &ekstypes.NodegroupUpdateConfig{
				MaxUnavailable: aws.Int32(6),
			},
			ClusterName:   &clusterName,
			NodegroupName: aws.String("ng-2"),
		}).Return(nil, nil)

		cfg.ManagedNodeGroups[0].UpdateConfig = &api.NodeGroupUpdateConfig{
			MaxUnavailable: aws.Int(6),
		}

		newNg := &api.ManagedNodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				Name: "ng-2",
			},
			UpdateConfig: &api.NodeGroupUpdateConfig{
				MaxUnavailable: aws.Int(6),
			},
		}

		cfg.ManagedNodeGroups = append(cfg.ManagedNodeGroups, newNg)

		m = New(cfg, &eks.ClusterProvider{AWSProvider: p}, nil, nil)
		err := m.Update(context.Background(), false)
		Expect(err).NotTo(HaveOccurred())
	})
})
