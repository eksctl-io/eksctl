package nodegroup

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws/awserr"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
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

		m = New(cfg, &eks.ClusterProvider{Provider: p}, nil)
	})

	It("fails for unmanaged nodegroups", func() {
		p.MockEKS().On("DescribeNodegroup", &awseks.DescribeNodegroupInput{
			ClusterName:   &m.cfg.Metadata.Name,
			NodegroupName: &m.cfg.ManagedNodeGroups[0].Name,
		}).Return(nil, awserr.New(awseks.ErrCodeResourceNotFoundException, "test-err", errors.New("err")))

		err := m.Update()
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(ContainSubstring("update is only supported for managed nodegroups; could not find one with name \"my-ng\"")))
	})

	It("successfully updates nodegroup", func() {
		p.MockEKS().On("DescribeNodegroup", &awseks.DescribeNodegroupInput{
			ClusterName:   &m.cfg.Metadata.Name,
			NodegroupName: &m.cfg.ManagedNodeGroups[0].Name,
		}).Return(nil, nil)

		err := m.Update()
		Expect(err).NotTo(HaveOccurred())
	})
})
