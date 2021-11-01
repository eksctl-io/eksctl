package eks_test

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/stretchr/testify/mock"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	. "github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var _ = Describe("eks auth helpers", func() {
	var ctl *ClusterProviderImpl

	Describe("construct client configs", func() {
		Context("with a mock provider", func() {
			var awsProvider *mockprovider.MockAwsProvider
			var kubeProvider *AWSKubeProvider
			clusterName := "auth-test-cluster"

			BeforeEach(func() {
				awsProvider = mockprovider.NewMockAwsProvider()
				kubeProvider = &AWSKubeProvider{AWSProvider: awsProvider}

				ctl = NewWithMocks(awsProvider, kubeProvider)
				kubeProvider.ClusterProvider = ctl
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

				assertConfigValid := func(k *clientcmdapi.Config) {
					Expect(k).NotTo(BeNil())

					ctx := k.CurrentContext
					s := strings.Split(ctx, "@")
					Expect(s).To(HaveLen(2))
					cluster := s[1]
					Expect(cluster).To(Equal("auth-test-cluster.eu-west-3.eksctl.io"))

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

					Expect(k.Clusters).To(HaveKey(cluster))
					Expect(k.Clusters).To(HaveLen(1))

					Expect(k.Clusters[cluster].InsecureSkipTLSVerify).To(BeFalse())
					Expect(k.Clusters[cluster].Server).To(Equal(cfg.Status.Endpoint))
					Expect(k.Clusters[cluster].CertificateAuthorityData).To(Equal(cfg.Status.CertificateAuthorityData))
				}

				testAuthenticatorConfig := func(roleARN string) {
					k := kubeconfig.NewForKubectl(cfg, ctl.GetUsername(), roleARN, ctl.AWSProvider().Profile())
					ctx := k.CurrentContext

					// test shared expectations
					assertConfigValid(k)

					// test authenticator context
					username := strings.Split(ctx, "@")[0]
					Expect(username).To(Equal("iam-root-account"))
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
				}

				testClientConfig := func(roleARN string) {
					if roleARN != "" {
						awsProvider.MockSTS().On("GetCallerIdentity", mock.Anything).Return(&sts.GetCallerIdentityOutput{
							Arn: aws.String(roleARN),
						}, nil)
						Expect(ctl.CheckAuth()).To(Succeed()) // set roleARN
					}

					req := awsProvider.Client.NewRequest(&request.Operation{Name: "GetCallerIdentityRequest"}, nil, nil)
					awsProvider.MockSTS().On("GetCallerIdentityRequest", mock.Anything).Return(req, nil)

					client, err := ctl.NewClient(cfg)
					Expect(err).NotTo(HaveOccurred())

					// test shared expectations
					config := client.Config()
					assertConfigValid(config)

					// test embedded token
					ctx := config.CurrentContext
					username := strings.Split(ctx, "@")[0]
					if roleARN != "" {
						expectedUsername := strings.Split(roleARN, "/")[1]
						Expect(username).To(Equal(expectedUsername))
					} else {
						Expect(username).To(Equal("iam-root-account"))
					}
					Expect(config.AuthInfos[ctx].Token).ToNot(BeEmpty())
					Expect(config.AuthInfos[ctx].Exec).To(BeNil())
				}

				It("should create config with authenticator", func() {
					testAuthenticatorConfig("")
					testAuthenticatorConfig("arn:aws:iam::111111111111:role/eksctl")
				})

				It("should create config with embedded token", func() {
					testClientConfig("")
					testClientConfig("arn:aws:iam::111111111111:role/eksctl")
				})
			})
		})
	})
})
