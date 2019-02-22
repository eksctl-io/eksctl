package eks_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	. "github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var _ = Describe("eks auth helpers", func() {
	var ctl *ClusterProvider

	Describe("constuct client configs", func() {
		Context("with a mock provider", func() {
			clusterName := "auth-test-cluster"

			BeforeEach(func() {

				p := mockprovider.NewMockProvider()

				ctl = &ClusterProvider{
					Provider: p,
					Status:   &ProviderStatus{},
				}

			})

			Context("for a cluster", func() {
				cfg := &api.ClusterConfig{
					Metadata: &api.ClusterMeta{
						Name:   clusterName,
						Region: "eu-west-3",
					},
					Status: &api.ClusterStatus{
						Endpoint:                 "https://TEST.aws",
						CertificateAuthorityData: []byte("123"),
					},
				}

				It("should create config with authenticator", func() {
					clientConfig, err := ctl.NewClientConfig(cfg, false)

					Expect(err).To(Not(HaveOccurred()))

					Expect(clientConfig).To(Not(BeNil()))
					ctx := clientConfig.ContextName
					cluster := strings.Split(ctx, "@")[1]
					Expect(ctx).To(Equal("iam-root-account@auth-test-cluster.eu-west-3.eksctl.io"))

					k := clientConfig.Client

					Expect(k.CurrentContext).To(Equal(ctx))

					Expect(k.Contexts).To(HaveKey(ctx))
					Expect(k.Contexts).To(HaveLen(1))

					Expect(k.Contexts[ctx].Cluster).To(Equal(cluster))
					Expect(k.Contexts[ctx].AuthInfo).To(Equal(ctx))

					Expect(k.Contexts[ctx].LocationOfOrigin).To(BeEmpty())
					Expect(k.Contexts[ctx].Namespace).To(BeEmpty())
					Expect(k.Contexts[ctx].Extensions).To(BeNil())

					Expect(k.AuthInfos).To(HaveKey(ctx))
					Expect(k.AuthInfos).To(HaveLen(1))

					Expect(k.AuthInfos[ctx].Token).To(BeEmpty())
					Expect(k.AuthInfos[ctx].Exec).To(Not(BeNil()))

					Expect(k.AuthInfos[ctx].Exec.Command).To(MatchRegexp("(heptio-authenticator-aws|aws-iam-authenticator)"))

					Expect(strings.Join(k.AuthInfos[ctx].Exec.Args, " ")).To(Equal("token -i auth-test-cluster"))

					Expect(k.Clusters).To(HaveKey(cluster))
					Expect(k.Clusters).To(HaveLen(1))

					Expect(k.Clusters[cluster].InsecureSkipTLSVerify).To(BeFalse())
					Expect(k.Clusters[cluster].Server).To(Equal(cfg.Status.Endpoint))
					Expect(k.Clusters[cluster].CertificateAuthorityData).To(Equal(cfg.Status.CertificateAuthorityData))
				})

				It("should create config with embedded token", func() {
					// TODO: cannot test this, as token generator uses STS directly, we cannot pass the interface
					// we can probably fix the package itself
				})

				It("should create clientset", func() {
					clientConfig, err := ctl.NewClientConfig(cfg, false)

					Expect(err).To(Not(HaveOccurred()))
					Expect(clientConfig).To(Not(BeNil()))

					clientSet, err := clientConfig.NewClientSet()

					Expect(err).To(Not(HaveOccurred()))
					Expect(clientSet).To(Not(BeNil()))
				})
			})
		})
	})
})
