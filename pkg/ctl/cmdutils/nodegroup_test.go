package cmdutils

import (
	"github.com/aws/aws-sdk-go/aws"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager/fakes"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("PopulateNodegroup", func() {
	var (
		ngName           string
		cfg              *api.ClusterConfig
		fakeStackManager *fakes.FakeStackManager
		mockProvider     *mockprovider.MockProvider
		err              error
	)

	BeforeEach(func() {
		fakeStackManager = new(fakes.FakeStackManager)
		ngName = "ng"
		cfg = api.NewClusterConfig()
		mockProvider = mockprovider.NewMockProvider()
	})

	Context("unmanaged nodegroup", func() {
		It("is added to the cfg", func() {
			fakeStackManager.GetNodeGroupStackTypeReturns(api.NodeGroupTypeUnmanaged, nil)
			err = PopulateNodegroup(fakeStackManager, ngName, cfg, mockProvider)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.NodeGroups[0].Name).To(Equal(ngName))
		})
	})

	Context("managed nodegroup", func() {
		It("is added to the cfg", func() {
			fakeStackManager.GetNodeGroupStackTypeReturns(api.NodeGroupTypeManaged, nil)
			err = PopulateNodegroup(fakeStackManager, ngName, cfg, mockProvider)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.ManagedNodeGroups[0].Name).To(Equal(ngName))
		})
	})

	Context("unowned nodegroup", func() {
		It("is added to the cfg", func() {
			clusterName := "cluster-name"
			cfg.Metadata.Name = clusterName
			fakeStackManager.GetNodeGroupStackTypeReturns("", errors.New(""))
			mockProvider.MockEKS().On("DescribeNodegroup", &awseks.DescribeNodegroupInput{
				ClusterName:   aws.String(clusterName),
				NodegroupName: aws.String(ngName),
			}).Return(&awseks.DescribeNodegroupOutput{Nodegroup: &awseks.Nodegroup{}}, nil)

			err = PopulateNodegroup(fakeStackManager, ngName, cfg, mockProvider)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.ManagedNodeGroups[0].Name).To(Equal(ngName))
		})
	})
})
