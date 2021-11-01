package nodegroup

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/fake"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Update", func() {
	var (
		clusterName, ngName string
		awsProvider         *mockprovider.MockAwsProvider
		kubeProvider        *mockprovider.MockKubeProvider
		cfg                 *api.ClusterConfig
		ctl                 *eks.ClusterProviderImpl
		m                   *Manager
	)

	BeforeEach(func() {
		clusterName = "my-cluster"
		ngName = "my-ng"
		clientSet := fake.NewSimpleClientset()
		awsProvider = mockprovider.NewMockAwsProvider()
		kubeProvider = mockprovider.NewMockKubeProvider(clientSet)
		ctl = eks.NewWithMocks(awsProvider, kubeProvider)

		cfg = api.NewClusterConfig()
		cfg.Metadata.Name = clusterName
		cfg.ManagedNodeGroups = []*api.ManagedNodeGroup{
			{
				NodeGroupBase: &api.NodeGroupBase{
					Name: ngName,
				},
			},
		}

		m = New(cfg, ctl, clientSet)
	})

	It("fails for unmanaged nodegroups", func() {
		awsProvider.MockEKS().On("DescribeNodegroup", &awseks.DescribeNodegroupInput{
			ClusterName:   &clusterName,
			NodegroupName: &ngName,
		}).Return(nil, awserr.New(awseks.ErrCodeResourceNotFoundException, "test-err", errors.New("err")))

		err := m.Update()
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring("could not find managed nodegroup with name \"my-ng\"")))
	})

	It("[happy path] successfully updates a nodegroup with updateConfig and maxUnavailable", func() {
		awsProvider.MockEKS().On("DescribeNodegroup", &awseks.DescribeNodegroupInput{
			ClusterName:   &clusterName,
			NodegroupName: &ngName,
		}).Return(&awseks.DescribeNodegroupOutput{
			Nodegroup: &awseks.Nodegroup{
				UpdateConfig: &awseks.NodegroupUpdateConfig{
					MaxUnavailable: aws.Int64(4),
				},
			},
		}, nil)

		awsProvider.MockEKS().On("UpdateNodegroupConfig", &awseks.UpdateNodegroupConfigInput{
			UpdateConfig: &awseks.NodegroupUpdateConfig{
				MaxUnavailable: aws.Int64(6),
			},
			ClusterName:   &clusterName,
			NodegroupName: &ngName,
		}).Return(nil, nil)

		cfg.ManagedNodeGroups[0].UpdateConfig = &api.NodeGroupUpdateConfig{
			MaxUnavailable: aws.Int(6),
		}

		err := m.Update()
		Expect(err).NotTo(HaveOccurred())
	})

	It("[happy path] successfully updates multiple nodegroups with updateConfig and maxUnavailable", func() {
		awsProvider.MockEKS().On("DescribeNodegroup", &awseks.DescribeNodegroupInput{
			ClusterName:   &clusterName,
			NodegroupName: &ngName,
		}).Return(&awseks.DescribeNodegroupOutput{
			Nodegroup: &awseks.Nodegroup{
				UpdateConfig: &awseks.NodegroupUpdateConfig{
					MaxUnavailable: aws.Int64(4),
				},
			},
		}, nil)

		awsProvider.MockEKS().On("DescribeNodegroup", &awseks.DescribeNodegroupInput{
			ClusterName:   &clusterName,
			NodegroupName: aws.String("ng-2"),
		}).Return(&awseks.DescribeNodegroupOutput{
			Nodegroup: &awseks.Nodegroup{
				UpdateConfig: &awseks.NodegroupUpdateConfig{
					MaxUnavailable: aws.Int64(4),
				},
			},
		}, nil)

		awsProvider.MockEKS().On("UpdateNodegroupConfig", &awseks.UpdateNodegroupConfigInput{
			UpdateConfig: &awseks.NodegroupUpdateConfig{
				MaxUnavailable: aws.Int64(6),
			},
			ClusterName:   &clusterName,
			NodegroupName: &ngName,
		}).Return(nil, nil)

		awsProvider.MockEKS().On("UpdateNodegroupConfig", &awseks.UpdateNodegroupConfigInput{
			UpdateConfig: &awseks.NodegroupUpdateConfig{
				MaxUnavailable: aws.Int64(6),
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

		err := m.Update()
		Expect(err).NotTo(HaveOccurred())
	})
})
