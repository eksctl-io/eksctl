package builder_test

import (
	"context"
	"os"
	"path"

	. "github.com/benjamintf1/unmarshalledmatchers"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"

	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	bootstrapfakes "github.com/weaveworks/eksctl/pkg/nodebootstrap/fakes"
	vpcfakes "github.com/weaveworks/eksctl/pkg/vpc/fakes"
)

var _ = Describe("Nodegroup Builder", func() {
	type ngResourceTest struct {
		disableAccessEntryCreation bool
		resourceFilename           string
	}

	DescribeTable("AddAllResources", func(t ngResourceTest) {
		provider := mockprovider.NewMockProvider()
		fakeVPCImporter := &vpcfakes.FakeImporter{}
		fakeBootstrapper := &bootstrapfakes.FakeBootstrapper{}
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Name = "cluster"
		mockSubnetsAndAZInstanceSupport(clusterConfig, provider, []string{"us-west-1a"}, nil, nil)

		resourceSet := builder.NewNodeGroupResourceSet(provider.EC2(), provider.IAM(), builder.NodeGroupOptions{
			ClusterConfig:              clusterConfig,
			NodeGroup:                  api.NewNodeGroup(),
			Bootstrapper:               fakeBootstrapper,
			ForceAddCNIPolicy:          false,
			VPCImporter:                fakeVPCImporter,
			SkipEgressRules:            false,
			DisableAccessEntryCreation: t.disableAccessEntryCreation,
		})
		Expect(resourceSet.AddAllResources(context.Background())).To(Succeed())
		actual, err := resourceSet.RenderJSON()
		Expect(err).NotTo(HaveOccurred())
		expected, err := os.ReadFile(path.Join("testdata", "nodegroup_access_entry", t.resourceFilename))
		Expect(err).NotTo(HaveOccurred())
		Expect(actual).To(MatchOrderedJSON(expected, WithUnorderedListKeys("Tags")))
	},
		Entry("with access entry", ngResourceTest{
			resourceFilename: "1.json",
		}),
		Entry("without access entry", ngResourceTest{
			disableAccessEntryCreation: true,
			resourceFilename:           "2.json",
		}),
	)
})
