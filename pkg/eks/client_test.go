package eks_test

import (
	"fmt"
	"strings"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5/fakes"
	. "github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

var _ = Describe("eks auth helpers", func() {
	var (
		ctl          *ClusterProvider
		kubeProvider *KubernetesProvider
	)

	Describe("construct client configs", func() {
		Context("with a mock provider", func() {
			clusterName := "auth-test-cluster"

			BeforeEach(func() {

				p := mockprovider.NewMockProvider()
				kubeProvider = &KubernetesProvider{
					WaitTimeout: 1 * time.Second,
					RoleARN:     "role",
					Signer:      p.STSPresigner(),
				}

				ctl = &ClusterProvider{
					AWSProvider: p,
					Status: &ProviderStatus{
						IAMRoleARN: "eksctl",
					},
					KubeProvider: kubeProvider,
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

				testAuthenticatorConfig := func(roleARN string) {
					clientConfig := kubeconfig.NewForKubectl(cfg, GetUsername(ctl.Status.IAMRoleARN), roleARN, ctl.AWSProvider.Profile())
					Expect(clientConfig).To(Not(BeNil()))
					ctx := clientConfig.CurrentContext
					cluster := strings.Split(ctx, "@")[1]
					Expect(ctx).To(Equal("iam-root-account@auth-test-cluster.eu-west-3.eksctl.io"))

					k := clientConfig

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

					// TODO: This test depends on which authenticator(s) is(are) installed and
					// the code deciding which one should be picked up. Ideally we'd like to
					// test all combinations, probably best done with a unit test.
					Expect(k.AuthInfos[ctx].Exec.Command).To(MatchRegexp("(heptio-authenticator-aws|aws-iam-authenticator|aws)"))

					var expectedArgs, roleARNArg string
					switch k.AuthInfos[ctx].Exec.Command {
					case "aws":
						expectedArgs = "eks get-token --cluster-name auth-test-cluster --region eu-west-3"
						roleARNArg = "--role-arn"
					case "heptio-authenticator-aws":
						fallthrough
					case "aws-iam-authenticator":
						expectedArgs = "token -i auth-test-cluster"
						roleARNArg = "-r"
					}
					if roleARN != "" {
						expectedArgs += fmt.Sprintf(" %s %s", roleARNArg, roleARN)
					}
					Expect(strings.Join(k.AuthInfos[ctx].Exec.Args, " ")).To(Equal(expectedArgs))

					Expect(k.Clusters).To(HaveKey(cluster))
					Expect(k.Clusters).To(HaveLen(1))

					Expect(k.Clusters[cluster].InsecureSkipTLSVerify).To(BeFalse())
					Expect(k.Clusters[cluster].Server).To(Equal(cfg.Status.Endpoint))
					Expect(k.Clusters[cluster].CertificateAuthorityData).To(Equal(cfg.Status.CertificateAuthorityData))
				}

				It("should create config with authenticator", func() {
					testAuthenticatorConfig("")
					testAuthenticatorConfig("arn:aws:iam::111111111111:role/eksctl")
				})

				It("should create config with embedded token", func() {
					mockPresigner := ctl.AWSProvider.STSPresigner().(*fakes.FakeSTSPresigner)
					mockPresigner.PresignGetCallerIdentityReturns(&v4.PresignedHTTPRequest{
						URL: "https://example.com",
					}, nil)
					client, err := kubeProvider.NewClient(cfg)
					Expect(err).NotTo(HaveOccurred())
					Expect(client.Config.AuthInfos[client.Config.CurrentContext].Token).To(HavePrefix("k8s-aws-v1."))
				})
			})
		})
	})
})
