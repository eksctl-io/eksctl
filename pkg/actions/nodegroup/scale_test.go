package nodegroup_test

import (
	"fmt"
	"time"

	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Scale", func() {
	When("the nodegroup was not created by eksctl", func() {
		var (
			clusterName, ngName string
			p                   *mockprovider.MockProvider
			cfg                 *api.ClusterConfig
			ng                  *api.NodeGroup
			manager             *nodegroup.Manager
		)
		BeforeEach(func() {
			clusterName = "my-cluster"
			ngName = "my-ng"
			p = mockprovider.NewMockProvider()
			cfg = api.NewClusterConfig()
			cfg.Metadata.Name = clusterName

			ng = &api.NodeGroup{
				NodeGroupBase: &api.NodeGroupBase{
					Name: ngName,
					ScalingConfig: &api.ScalingConfig{
						MinSize:         aws.Int(1),
						DesiredCapacity: aws.Int(3),
					},
				},
			}
			manager = nodegroup.New(cfg, &eks.ClusterProvider{Provider: p}, nil)
			p.MockCloudFormation().On("ListStacksPages", mock.Anything, mock.Anything).Return(nil, nil)
		})

		It("scales the nodegroup using the values provided", func() {
			p.MockEKS().On("UpdateNodegroupConfig", &awseks.UpdateNodegroupConfigInput{
				ScalingConfig: &awseks.NodegroupScalingConfig{
					MinSize:     aws.Int64(1),
					DesiredSize: aws.Int64(3),
				},
				ClusterName:   &clusterName,
				NodegroupName: &ngName,
			}).Return(nil, nil)

			p.MockEKS().On("DescribeNodegroupRequest", &awseks.DescribeNodegroupInput{
				ClusterName:   &clusterName,
				NodegroupName: &ngName,
			}).Return(&request.Request{}, nil)

			waitCallCount := 0
			manager.SetWaiter(func(name, msg string, acceptors []request.WaiterAcceptor, newRequest func() *request.Request, waitTimeout time.Duration, troubleshoot func(string) error) error {
				waitCallCount++
				return nil
			})

			err := manager.Scale(ng)

			Expect(err).NotTo(HaveOccurred())
			Expect(waitCallCount).To(Equal(1))
		})

		When("upgrade fails", func() {
			It("returns an error", func() {
				p.MockEKS().On("UpdateNodegroupConfig", &awseks.UpdateNodegroupConfigInput{
					ScalingConfig: &awseks.NodegroupScalingConfig{
						MinSize:     aws.Int64(1),
						DesiredSize: aws.Int64(3),
					},
					ClusterName:   &clusterName,
					NodegroupName: &ngName,
				}).Return(nil, fmt.Errorf("foo"))

				err := manager.Scale(ng)

				Expect(err).To(MatchError(fmt.Sprintf("failed to scale nodegroup for cluster %q, error: foo", clusterName)))
			})
		})
	})

})
