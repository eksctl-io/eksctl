package v1alpha5_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

const identityProviders = `addonsConfig: {}
apiVersion: eksctl.io/v1alpha5
identityProviders:
- clientID: client
  issuerURL: example.com
  name: name
  type: oidc
  usernameClaim: email
kind: ClusterConfig
metadata:
  name: ip
  region: us-west-2
`

var _ = Describe("IdentityProvider", func() {
	It("can be Unmarshaled and marshaled", func() {
		Expect(api.Register()).To(Succeed())
		cfg, err := eks.ParseConfig([]byte(identityProviders))
		Expect(err).NotTo(HaveOccurred())
		Expect(*cfg).To(Equal(api.ClusterConfig{
			Metadata: &api.ClusterMeta{
				Name:   "ip",
				Region: "us-west-2",
			},
			TypeMeta: metav1.TypeMeta{
				APIVersion: "eksctl.io/v1alpha5",
				Kind:       "ClusterConfig",
			},
			IdentityProviders: []api.IdentityProvider{
				{
					Inner: &api.OIDCIdentityProvider{
						ClientID:      "client",
						IssuerURL:     "example.com",
						Name:          "name",
						UsernameClaim: "email",
					},
				},
			},
		}))

		data, err := yaml.Marshal(cfg)
		Expect(err).NotTo(HaveOccurred())
		fmt.Println(string(data))
		Expect(string(data)).To(Equal(identityProviders))
	})
})
