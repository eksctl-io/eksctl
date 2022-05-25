package cmdutils

import (
	"context"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/stretchr/testify/mock"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	. "github.com/onsi/ginkgo/v2"
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
			err = PopulateNodegroup(context.Background(), fakeStackManager, ngName, cfg, mockProvider)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.NodeGroups[0].Name).To(Equal(ngName))
		})
	})

	Context("managed nodegroup", func() {
		It("is added to the cfg", func() {
			fakeStackManager.GetNodeGroupStackTypeReturns(api.NodeGroupTypeManaged, nil)
			err = PopulateNodegroup(context.Background(), fakeStackManager, ngName, cfg, mockProvider)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.ManagedNodeGroups[0].Name).To(Equal(ngName))
		})
	})

	Context("unowned nodegroup", func() {
		It("is added to the cfg", func() {
			clusterName := "cluster-name"
			cfg.Metadata.Name = clusterName
			fakeStackManager.GetNodeGroupStackTypeReturns("", errors.New(""))
			mockProvider.MockEKS().On("DescribeNodegroup", mock.Anything, &awseks.DescribeNodegroupInput{
				ClusterName:   aws.String(clusterName),
				NodegroupName: aws.String(ngName),
			}).Return(&awseks.DescribeNodegroupOutput{Nodegroup: &ekstypes.Nodegroup{}}, nil)

			err = PopulateNodegroup(context.Background(), fakeStackManager, ngName, cfg, mockProvider)
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.ManagedNodeGroups[0].Name).To(Equal(ngName))
		})
	})
})
