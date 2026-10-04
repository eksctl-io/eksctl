package builder_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

var _ = Describe("Capability", func() {
	Describe("CapabilityResourceSet", func() {
		var capability api.Capability

		BeforeEach(func() {
			capability = api.Capability{
				Name: "my-capability",
				Type: "ARGOCD",
				Configuration: &api.CapabilityConfiguration{
					ArgoCD: &api.ArgoCDConfiguration{
						Namespace: "argocd",
					},
				},
			}
		})

		It("renders AWSIDC configuration when provided", func() {
			capability.Configuration.ArgoCD.AWSIDC = &api.ArgoCDAWSIDC{
				IDCInstanceARN: "arn:aws:sso:::instance/ssoins-123",
				IDCRegion:      "us-west-2",
			}
			rs := builder.NewCapabilityResourceSet("my-cluster", capability)
			Expect(rs.AddAllResources()).To(Succeed())
		})

		It("returns a clean error instead of panicking when AWSIDC is missing", func() {
			rs := builder.NewCapabilityResourceSet("my-cluster", capability)
			err := rs.AddAllResources()
			Expect(err).To(MatchError(ContainSubstring("awsIdc configuration is required for ARGOCD capability")))
		})
	})
})
