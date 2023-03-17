package disassociate

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/weaveworks/eksctl/pkg/actions/identityproviders"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("cliToProviders", func() {
	It("works with CLI arguments", func() {
		cfg := api.NewClusterConfig()
		Expect(
			cliToProviders(cfg, cliProvidedIDP{Name: "idp", Type: "oidc"}),
		).To(Equal([]identityproviders.DisassociateIdentityProvider{
			{Name: "idp", Type: api.OIDCIdentityProviderType},
		}))
	})
	It("works with only config providers", func() {
		cfg := api.NewClusterConfig()
		cfg.IdentityProviders = []api.IdentityProvider{
			api.FromIdentityProvider(&api.OIDCIdentityProvider{
				Name:      "idp",
				IssuerURL: "http://url.com",
				ClientID:  "client-id",
			}),
		}
		Expect(
			cliToProviders(cfg, cliProvidedIDP{}),
		).To(Equal([]identityproviders.DisassociateIdentityProvider{
			{Name: "idp", Type: api.OIDCIdentityProviderType},
		}))
	})
	It("works with CLI override", func() {
		cfg := api.NewClusterConfig()
		cfg.IdentityProviders = []api.IdentityProvider{
			api.FromIdentityProvider(&api.OIDCIdentityProvider{
				Name:      "idp",
				IssuerURL: "http://url.com",
				ClientID:  "client-id",
			}),
		}
		Expect(
			cliToProviders(cfg, cliProvidedIDP{Name: "idp2", Type: "oidc"}),
		).To(Equal([]identityproviders.DisassociateIdentityProvider{
			{Name: "idp2", Type: api.OIDCIdentityProviderType},
		}))
	})
})
