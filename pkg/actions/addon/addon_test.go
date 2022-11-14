package addon_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/actions/addon"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("Addon", func() {
	var (
		clusterConfig *api.ClusterConfig
		mockProvider  *mockprovider.MockProvider
		manager       *addon.Manager
	)

	BeforeEach(func() {
		var err error
		clusterConfig = &api.ClusterConfig{Metadata: &api.ClusterMeta{
			Version: "1.18",
			Name:    "my-cluster",
		}}
		mockProvider = mockprovider.NewMockProvider()
		manager, err = addon.New(clusterConfig, mockProvider.EKS(), nil, false, nil, nil)
		Expect(err).NotTo(HaveOccurred())
	})

	When("addon is coredns with no nodegroup", func() {
		It("does not wait for addon to be active", func() {
			a := &api.Addon{
				Name: api.CoreDNSAddon,
			}
			err := addon.ExportWaitForAddonToBeActive(manager, context.Background(), a, 0)
			Expect(err).To(BeNil())
		})
	})

	When("addon is aws-ebs-csi-driver with no nodegroup", func() {
		It("does not wait for addon to be active", func() {
			a := &api.Addon{
				Name: api.AWSEBSCSIDriverAddon,
			}
			err := addon.ExportWaitForAddonToBeActive(manager, context.Background(), a, 0)
			Expect(err).To(BeNil())
		})
	})
})
