// +build integration

package crud

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/integration/utilities/kube"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/iam"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/utils/file"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	harness "github.com/dlespiau/kube-test-harness"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("crud")
}

func TestSuite(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) Create, Get, Scale & Delete", func() {

	const (
		initNG = "ng-0"
		testNG = "ng-1"
	)

	commonTimeout := 10 * time.Minute

	BeforeSuite(func() {
		params.KubeconfigTemp = false
		if params.KubeconfigPath == "" {
			wd, _ := os.Getwd()
			f, _ := ioutil.TempFile(wd, "kubeconfig-")
			params.KubeconfigPath = f.Name()
			params.KubeconfigTemp = true
		}
	})

	AfterSuite(func() {
		params.DeleteClusters()
		gexec.KillAndWait()
		if params.KubeconfigTemp {
			os.Remove(params.KubeconfigPath)
		}
		os.RemoveAll(params.TestDirectory)
	})

	Describe("when creating a cluster with 1 node", func() {
		It("should not return an error", func() {
			if params.SkipCreate {
				fmt.Fprintf(GinkgoWriter, "will use existing cluster %s", params.ClusterName)
				if !file.Exists(params.KubeconfigPath) {
					// Generate the Kubernetes configuration that eksctl create
					// would have generated otherwise:
					cmd := params.EksctlUtilsCmd.WithArgs(
						"write-kubeconfig",
						"--verbose", "4",
						"--cluster", params.ClusterName,
						"--kubeconfig", params.KubeconfigPath,
					)
					Expect(cmd).To(RunSuccessfully())
				}
				return
			}

			fmt.Fprintf(GinkgoWriter, "Using kubeconfig: %s\n", params.KubeconfigPath)

			cmd := params.EksctlCreateCmd.WithArgs(
				"cluster",
				"--verbose", "4",
				"--name", params.ClusterName,
				"--tags", "alpha.eksctl.io/description=eksctl integration test",
				"--nodegroup-name", initNG,
				"--node-labels", "ng-name="+initNG,
				"--nodes", "1",
				"--version", params.Version,
				"--kubeconfig", params.KubeconfigPath,
			)
			Expect(cmd).To(RunSuccessfully())
		})

		It("should have created an EKS cluster and two CloudFormation stacks", func() {
			awsSession := NewSession(params.Region)

			Expect(awsSession).To(HaveExistingCluster(params.ClusterName, awseks.ClusterStatusActive, params.Version))

			Expect(awsSession).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-cluster", params.ClusterName)))
			Expect(awsSession).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%s", params.ClusterName, initNG)))
		})

		It("should have created a valid kubectl config file", func() {
			config, err := clientcmd.LoadFromFile(params.KubeconfigPath)
			Expect(err).ShouldNot(HaveOccurred())

			err = clientcmd.ConfirmUsable(*config, "")
			Expect(err).ShouldNot(HaveOccurred())

			Expect(config.CurrentContext).To(ContainSubstring("eksctl"))
			Expect(config.CurrentContext).To(ContainSubstring(params.ClusterName))
			Expect(config.CurrentContext).To(ContainSubstring(params.Region))
		})

		Context("and listing clusters", func() {
			It("should return the previously created cluster", func() {
				cmd := params.EksctlGetCmd.WithArgs("clusters", "--all-regions")
				Expect(cmd).To(RunSuccessfullyWithOutputString(ContainSubstring(params.ClusterName)))
			})
		})

		Context("and scale the initial nodegroup", func() {
			It("should not return an error", func() {
				cmd := params.EksctlScaleNodeGroupCmd.WithArgs(
					"--cluster", params.ClusterName,
					"--nodes", "4",
					"--name", initNG,
				)
				Expect(cmd).To(RunSuccessfully())
			})
		})

		Context("and add the second nodegroup", func() {
			It("should not return an error", func() {
				cmd := params.EksctlCreateCmd.WithArgs(
					"nodegroup",
					"--cluster", params.ClusterName,
					"--nodes", "4",
					"--node-private-networking",
					testNG,
				)
				Expect(cmd).To(RunSuccessfully())
			})

			It("should be able to list nodegroups", func() {
				{
					cmd := params.EksctlGetCmd.WithArgs(
						"nodegroup",
						"-o", "json",
						"--cluster", params.ClusterName,
						initNG,
					)
					Expect(cmd).To(RunSuccessfullyWithOutputString(BeNodeGroupsWithNamesWhich(
						HaveLen(1),
						ContainElement(initNG),
						Not(ContainElement(testNG)),
					)))
				}
				{
					cmd := params.EksctlGetCmd.WithArgs(
						"nodegroup",
						"-o", "json",
						"--cluster", params.ClusterName,
						testNG,
					)
					Expect(cmd).To(RunSuccessfullyWithOutputString(BeNodeGroupsWithNamesWhich(
						HaveLen(1),
						ContainElement(testNG),
						Not(ContainElement(initNG)),
					)))
				}
				{
					cmd := params.EksctlGetCmd.WithArgs(
						"nodegroup",
						"-o", "json",
						"--cluster", params.ClusterName,
					)
					Expect(cmd).To(RunSuccessfullyWithOutputString(BeNodeGroupsWithNamesWhich(
						HaveLen(2),
						ContainElement(initNG),
						ContainElement(testNG),
					)))
				}
			})

			Context("toggle CloudWatch logging", func() {
				var (
					cfg *api.ClusterConfig
					ctl *eks.ClusterProvider
				)

				BeforeEach(func() {
					cfg = &api.ClusterConfig{
						Metadata: &api.ClusterMeta{
							Name:   params.ClusterName,
							Region: params.Region,
						},
					}
					ctl = eks.New(&api.ProviderConfig{Region: params.Region}, cfg)
				})

				It("should have all types disabled by default", func() {
					enabled, disable, err := ctl.GetCurrentClusterConfigForLogging(cfg)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(enabled.List()).To(HaveLen(0))
					Expect(disable.List()).To(HaveLen(5))
				})

				It("should plan to enable two of the types using flags", func() {
					cmd := params.EksctlUtilsCmd.WithArgs(
						"update-cluster-logging",
						"--cluster", params.ClusterName,
						"--enable-types", "api,controllerManager",
					)
					Expect(cmd).To(RunSuccessfully())

					enabled, disable, err := ctl.GetCurrentClusterConfigForLogging(cfg)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(enabled.List()).To(HaveLen(0))
					Expect(disable.List()).To(HaveLen(5))
				})

				It("should enable two of the types using flags", func() {
					cmd := params.EksctlUtilsCmd.WithArgs(
						"update-cluster-logging",
						"--cluster", params.ClusterName,
						"--approve",
						"--enable-types", "api,controllerManager",
					)
					Expect(cmd).To(RunSuccessfully())

					enabled, disable, err := ctl.GetCurrentClusterConfigForLogging(cfg)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(enabled.List()).To(HaveLen(2))
					Expect(disable.List()).To(HaveLen(3))
					Expect(enabled.List()).To(ConsistOf("api", "controllerManager"))
					Expect(disable.List()).To(ConsistOf("audit", "authenticator", "scheduler"))
				})

				It("should enable all of the types with --enable-types=all", func() {
					cmd := params.EksctlUtilsCmd.WithArgs(
						"update-cluster-logging",
						"--cluster", params.ClusterName,
						"--approve",
						"--enable-types", "all",
					)
					Expect(cmd).To(RunSuccessfully())

					enabled, disable, err := ctl.GetCurrentClusterConfigForLogging(cfg)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(enabled.List()).To(HaveLen(5))
					Expect(disable.List()).To(HaveLen(0))
				})

				It("should enable all but one type", func() {
					cmd := params.EksctlUtilsCmd.WithArgs(
						"update-cluster-logging",
						"--cluster", params.ClusterName,
						"--approve",
						"--enable-types", "all",
						"--disable-types", "controllerManager",
					)
					Expect(cmd).To(RunSuccessfully())

					enabled, disable, err := ctl.GetCurrentClusterConfigForLogging(cfg)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(enabled.List()).To(HaveLen(4))
					Expect(disable.List()).To(HaveLen(1))
					Expect(enabled.List()).To(ConsistOf("api", "audit", "authenticator", "scheduler"))
					Expect(disable.List()).To(ConsistOf("controllerManager"))
				})

				It("should disable all but one type", func() {
					cmd := params.EksctlUtilsCmd.WithArgs(
						"update-cluster-logging",
						"--cluster", params.ClusterName,
						"--approve",
						"--disable-types", "all",
						"--enable-types", "controllerManager",
					)
					Expect(cmd).To(RunSuccessfully())

					enabled, disable, err := ctl.GetCurrentClusterConfigForLogging(cfg)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(disable.List()).To(HaveLen(4))
					Expect(enabled.List()).To(HaveLen(1))
					Expect(disable.List()).To(ConsistOf("api", "audit", "authenticator", "scheduler"))
					Expect(enabled.List()).To(ConsistOf("controllerManager"))
				})

				It("should disable all of the types with --disable-types=all", func() {
					cmd := params.EksctlUtilsCmd.WithArgs(
						"update-cluster-logging",
						"--cluster", params.ClusterName,
						"--approve",
						"--disable-types", "all",
					)
					Expect(cmd).To(RunSuccessfully())

					enabled, disable, err := ctl.GetCurrentClusterConfigForLogging(cfg)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(enabled.List()).To(HaveLen(0))
					Expect(disable.List()).To(HaveLen(5))
					Expect(disable.HasAll(api.SupportedCloudWatchClusterLogTypes()...)).To(BeTrue())
				})
			})

			Context("create and delete iamserviceaccounts", func() {
				var (
					cfg  *api.ClusterConfig
					ctl  *eks.ClusterProvider
					oidc *iamoidc.OpenIDConnectManager
					err  error
				)

				BeforeEach(func() {
					cfg = &api.ClusterConfig{
						Metadata: &api.ClusterMeta{
							Name:   params.ClusterName,
							Region: params.Region,
						},
					}
					ctl = eks.New(&api.ProviderConfig{Region: params.Region}, cfg)
					err = ctl.RefreshClusterStatus(cfg)
					Expect(err).ShouldNot(HaveOccurred())
					oidc, err = ctl.NewOpenIDConnectManager(cfg)
					Expect(err).ShouldNot(HaveOccurred())
				})

				It("should enable OIDC and create two iamserviceaccounts", func() {
					{
						exists, err := oidc.CheckProviderExists()
						Expect(err).ShouldNot(HaveOccurred())
						Expect(exists).To(BeFalse())
					}

					setupCmd := params.EksctlUtilsCmd.WithArgs(
						"associate-iam-oidc-provider",
						"--cluster", params.ClusterName,
						"--approve",
					)
					Expect(setupCmd).To(RunSuccessfully())

					{
						exists, err := oidc.CheckProviderExists()
						Expect(err).ShouldNot(HaveOccurred())
						Expect(exists).To(BeTrue())
					}

					cmds := []Cmd{
						params.EksctlCreateCmd.WithArgs(
							"iamserviceaccount",
							"--cluster", params.ClusterName,
							"--name", "app-cache-access",
							"--namespace", "app1",
							"--attach-policy-arn", "arn:aws:iam::aws:policy/AmazonDynamoDBReadOnlyAccess",
							"--attach-policy-arn", "arn:aws:iam::aws:policy/AmazonElastiCacheFullAccess",
							"--approve",
						),
						params.EksctlCreateCmd.WithArgs(
							"iamserviceaccount",
							"--cluster", params.ClusterName,
							"--name", "s3-read-only",
							"--attach-policy-arn", "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess",
							"--approve",
						),
					}

					Expect(cmds).To(RunSuccessfully())

					awsSession := NewSession(params.Region)

					stackNamePrefix := fmt.Sprintf("eksctl-%s-addon-iamserviceaccount-", params.ClusterName)

					Expect(awsSession).To(HaveExistingStack(stackNamePrefix + "default-s3-read-only"))
					Expect(awsSession).To(HaveExistingStack(stackNamePrefix + "app1-app-cache-access"))

					clientSet, err := ctl.NewStdClientSet(cfg)
					Expect(err).ShouldNot(HaveOccurred())

					{
						sa, err := clientSet.CoreV1().ServiceAccounts(metav1.NamespaceDefault).Get("s3-read-only", metav1.GetOptions{})
						Expect(err).ShouldNot(HaveOccurred())

						Expect(sa.Annotations).To(HaveLen(1))
						Expect(sa.Annotations).To(HaveKey(api.AnnotationEKSRoleARN))
						Expect(sa.Annotations[api.AnnotationEKSRoleARN]).To(MatchRegexp("^arn:aws:iam::.*:role/eksctl-" + truncate(params.ClusterName) + ".*$"))
					}

					{
						sa, err := clientSet.CoreV1().ServiceAccounts("app1").Get("app-cache-access", metav1.GetOptions{})
						Expect(err).ShouldNot(HaveOccurred())

						Expect(sa.Annotations).To(HaveLen(1))
						Expect(sa.Annotations).To(HaveKey(api.AnnotationEKSRoleARN))
						Expect(sa.Annotations[api.AnnotationEKSRoleARN]).To(MatchRegexp("^arn:aws:iam::.*:role/eksctl-" + truncate(params.ClusterName) + ".*$"))
					}
				})

				It("should list both iamserviceaccounts", func() {
					cmd := params.EksctlGetCmd.WithArgs(
						"iamserviceaccount",
						"--cluster", params.ClusterName,
					)

					Expect(cmd).To(RunSuccessfullyWithOutputString(MatchRegexp(
						`(?m:^NAMESPACE\s+NAME\s+ROLE\sARN$)` +
							`|(?m:^app1\s+app-cache-access\s+arn:aws:iam::.*$)` +
							`|(?m:^default\s+s3-read-only\s+arn:aws:iam::.*$)`,
					)))
				})

				It("delete both iamserviceaccounts", func() {
					cmds := []Cmd{
						params.EksctlDeleteCmd.WithArgs(
							"iamserviceaccount",
							"--cluster", params.ClusterName,
							"--name", "s3-read-only",
							"--wait",
						),
						params.EksctlDeleteCmd.WithArgs(
							"iamserviceaccount",
							"--cluster", params.ClusterName,
							"--name", "app-cache-access",
							"--namespace", "app1",
							"--wait",
						),
					}
					Expect(cmds).To(RunSuccessfully())

					awsSession := NewSession(params.Region)

					stackNamePrefix := fmt.Sprintf("eksctl-%s-addon-iamserviceaccount-", params.ClusterName)

					Expect(awsSession).ToNot(HaveExistingStack(stackNamePrefix + "default-s3-read-only"))
					Expect(awsSession).ToNot(HaveExistingStack(stackNamePrefix + "app1-app-cache-access"))
				})
			})

			Context("create test workloads", func() {
				var (
					err  error
					test *harness.Test
				)

				BeforeEach(func() {
					test, err = kube.NewTest(params.KubeconfigPath)
					Expect(err).ShouldNot(HaveOccurred())
				})

				AfterEach(func() {
					test.Close()
					Eventually(func() int {
						return len(test.ListPods(test.Namespace, metav1.ListOptions{}).Items)
					}, "3m", "1s").Should(BeZero())
				})

				It("should deploy podinfo service to the cluster and access it via proxy", func() {
					d := test.CreateDeploymentFromFile(test.Namespace, "../../data/podinfo.yaml")
					test.WaitForDeploymentReady(d, commonTimeout)

					pods := test.ListPodsFromDeployment(d)
					Expect(len(pods.Items)).To(Equal(2))

					// For each pod of the Deployment, check we receive a sensible response to a
					// GET request on /version.
					for _, pod := range pods.Items {
						Expect(pod.Namespace).To(Equal(test.Namespace))

						req := test.PodProxyGet(&pod, "", "/version")
						fmt.Fprintf(GinkgoWriter, "url = %#v", req.URL())

						var js interface{}
						test.PodProxyGetJSON(&pod, "", "/version", &js)

						Expect(js.(map[string]interface{})).To(HaveKeyWithValue("version", "1.5.1"))
					}
				})

				It("should have functional DNS", func() {
					d := test.CreateDaemonSetFromFile(test.Namespace, "../../data/test-dns.yaml")

					test.WaitForDaemonSetReady(d, commonTimeout)

					{
						ds, err := test.GetDaemonSet(test.Namespace, d.Name)
						Expect(err).ShouldNot(HaveOccurred())
						fmt.Fprintf(GinkgoWriter, "ds.Status = %#v", ds.Status)
					}
				})

				It("should have access to HTTP(S) sites", func() {
					d := test.CreateDaemonSetFromFile(test.Namespace, "../../data/test-http.yaml")

					test.WaitForDaemonSetReady(d, commonTimeout)

					{
						ds, err := test.GetDaemonSet(test.Namespace, d.Name)
						Expect(err).ShouldNot(HaveOccurred())
						fmt.Fprintf(GinkgoWriter, "ds.Status = %#v", ds.Status)
					}
				})

				It("should be able to run pods with an iamserviceaccount", func() {
					createCmd := params.EksctlCreateCmd.WithArgs(
						"iamserviceaccount",
						"--cluster", params.ClusterName,
						"--name", "s3-reader",
						"--namespace", test.Namespace,
						"--attach-policy-arn", "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess",
						"--approve",
					)

					Expect(createCmd).To(RunSuccessfully())

					d := test.CreateDeploymentFromFile(test.Namespace, "../../data/iamserviceaccount-checker.yaml")
					test.WaitForDeploymentReady(d, commonTimeout)

					pods := test.ListPodsFromDeployment(d)
					Expect(len(pods.Items)).To(Equal(2))

					// For each pod of the Deployment, check we get expected environment variables
					// via a GET request on /env.
					type sessionObject struct {
						AssumedRoleUser struct {
							AssumedRoleId, Arn string
						}
						Audience, Provider, SubjectFromWebIdentityToken string
						Credentials                                     struct {
							SecretAccessKey, SessionToken, Expiration, AccessKeyId string
						}
					}

					for _, pod := range pods.Items {
						Expect(pod.Namespace).To(Equal(test.Namespace))

						so := sessionObject{}

						var js []string
						test.PodProxyGetJSON(&pod, "", "/env", &js)

						Expect(js).To(ContainElement(HavePrefix("AWS_ROLE_ARN=arn:aws:iam::")))
						Expect(js).To(ContainElement("AWS_WEB_IDENTITY_TOKEN_FILE=/var/run/secrets/eks.amazonaws.com/serviceaccount/token"))
						Expect(js).To(ContainElement(HavePrefix("AWS_SESSION_OBJECT=")))

						for _, envVar := range js {
							if strings.HasPrefix(envVar, "AWS_SESSION_OBJECT=") {
								err := json.Unmarshal([]byte(strings.TrimPrefix(envVar, "AWS_SESSION_OBJECT=")), &so)
								Expect(err).ShouldNot(HaveOccurred())
							}
						}

						Expect(so.AssumedRoleUser.AssumedRoleId).To(HaveSuffix(":integration-test"))

						Expect(so.AssumedRoleUser.Arn).To(MatchRegexp("^arn:aws:sts::.*:assumed-role/eksctl-" + truncate(params.ClusterName) + "-.*/integration-test$"))

						Expect(so.Audience).To(Equal("sts.amazonaws.com"))

						Expect(so.Provider).To(MatchRegexp("^arn:aws:iam::.*:oidc-provider/oidc.eks." + params.Region + ".amazonaws.com/id/.*$"))

						Expect(so.SubjectFromWebIdentityToken).To(Equal("system:serviceaccount:" + test.Namespace + ":s3-reader"))

						Expect(so.Credentials.SecretAccessKey).ToNot(BeEmpty())
						Expect(so.Credentials.SessionToken).ToNot(BeEmpty())
						Expect(so.Credentials.Expiration).ToNot(BeEmpty())
						Expect(so.Credentials.AccessKeyId).ToNot(BeEmpty())
					}

					deleteCmd := params.EksctlDeleteCmd.WithArgs(
						"iamserviceaccount",
						"--cluster", params.ClusterName,
						"--name", "s3-reader",
						"--namespace", test.Namespace,
					)

					Expect(deleteCmd).To(RunSuccessfully())
				})
			})

			Context("and manipulating iam identity mappings", func() {
				var (
					expR0, expR1, expU0 string
					role0, role1        iam.Identity
					user0               iam.Identity
					admin               = "admin"
					alice               = "alice"
				)

				BeforeEach(func() {
					roleCanonicalArn := "arn:aws:iam::123456:role/eksctl-testing-XYZ"
					var err error
					role0 = iam.RoleIdentity{
						RoleARN: roleCanonicalArn,
						KubernetesIdentity: iam.KubernetesIdentity{
							KubernetesUsername: admin,
							KubernetesGroups:   []string{"system:masters", "system:nodes"},
						},
					}
					role1 = iam.RoleIdentity{
						RoleARN: roleCanonicalArn,
						KubernetesIdentity: iam.KubernetesIdentity{
							KubernetesGroups: []string{"system:something"},
						},
					}

					userCanonicalArn := "arn:aws:iam::123456:user/alice"

					user0 = iam.UserIdentity{
						UserARN: userCanonicalArn,
						KubernetesIdentity: iam.KubernetesIdentity{
							KubernetesUsername: alice,
							KubernetesGroups:   []string{"system:masters", "cryptographers"},
						},
					}

					bs, err := yaml.Marshal([]iam.Identity{role0})
					Expect(err).ShouldNot(HaveOccurred())
					expR0 = string(bs)

					bs, err = yaml.Marshal([]iam.Identity{role1})
					Expect(err).ShouldNot(HaveOccurred())
					expR1 = string(bs)

					bs, err = yaml.Marshal([]iam.Identity{user0})
					Expect(err).ShouldNot(HaveOccurred())
					expU0 = string(bs)
				})

				It("fails getting unknown role mapping", func() {
					cmd := params.EksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", "arn:aws:iam::123456:role/idontexist",
						"-o", "yaml",
					)
					Expect(cmd).ToNot(RunSuccessfully())
				})
				It("fails getting unknown user mapping", func() {
					cmd := params.EksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", "arn:aws:iam::123456:user/bob",
						"-o", "yaml",
					)
					Expect(cmd).ToNot(RunSuccessfully())
				})
				It("creates role mapping", func() {
					create := params.EksctlCreateCmd.WithArgs(
						"iamidentitymapping",
						"--name", params.ClusterName,
						"--arn", role0.ARN(),
						"--username", role0.Username(),
						"--group", role0.Groups()[0],
						"--group", role0.Groups()[1],
					)
					Expect(create).To(RunSuccessfully())

					get := params.EksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--name", params.ClusterName,
						"--arn", role0.ARN(),
						"-o", "yaml",
					)
					Expect(get).To(RunSuccessfullyWithOutputString(MatchYAML(expR0)))
				})
				It("creates user mapping", func() {
					create := params.EksctlCreateCmd.WithArgs(
						"iamidentitymapping",
						"--name", params.ClusterName,
						"--arn", user0.ARN(),
						"--username", user0.Username(),
						"--group", user0.Groups()[0],
						"--group", user0.Groups()[1],
					)
					Expect(create).To(RunSuccessfully())

					get := params.EksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", user0.ARN(),
						"-o", "yaml",
					)
					Expect(get).To(RunSuccessfullyWithOutputString(MatchYAML(expU0)))
				})
				It("creates a duplicate role mapping", func() {
					createRole := params.EksctlCreateCmd.WithArgs(
						"iamidentitymapping",
						"--name", params.ClusterName,
						"--arn", role0.ARN(),
						"--username", role0.Username(),
						"--group", role0.Groups()[0],
						"--group", role0.Groups()[1],
					)
					Expect(createRole).To(RunSuccessfully())

					get := params.EksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--name", params.ClusterName,
						"--arn", role0.ARN(),
						"-o", "yaml",
					)
					Expect(get).To(RunSuccessfullyWithOutputString(MatchYAML(expR0 + expR0)))
				})
				It("creates a duplicate user mapping", func() {
					createCmd := params.EksctlCreateCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", user0.ARN(),
						"--username", user0.Username(),
						"--group", user0.Groups()[0],
						"--group", user0.Groups()[1],
					)
					Expect(createCmd).To(RunSuccessfully())

					getCmd := params.EksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", user0.ARN(),
						"-o", "yaml",
					)
					Expect(getCmd).To(RunSuccessfullyWithOutputString(MatchYAML(expU0 + expU0)))
				})
				It("creates a duplicate role mapping with different identity", func() {
					createCmd := params.EksctlCreateCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", role1.ARN(),
						"--group", role1.Groups()[0],
					)
					Expect(createCmd).To(RunSuccessfully())

					getCmd := params.EksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", role1.ARN(),
						"-o", "yaml",
					)
					Expect(getCmd).To(RunSuccessfullyWithOutputString(MatchYAML(expR0 + expR0 + expR1)))
				})
				It("deletes a single role mapping fifo", func() {
					deleteCmd := params.EksctlDeleteCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", role1.ARN(),
					)
					Expect(deleteCmd).To(RunSuccessfully())

					getCmd := params.EksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", role1.ARN(),
						"-o", "yaml",
					)
					Expect(getCmd).To(RunSuccessfullyWithOutputString(MatchYAML(expR0 + expR1)))
				})
				It("fails when deleting unknown mapping", func() {
					deleteCmd := params.EksctlDeleteCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", "arn:aws:iam::123456:role/idontexist",
					)
					Expect(deleteCmd).ToNot(RunSuccessfully())
				})
				It("deletes duplicate role mappings with --all", func() {
					deleteCmd := params.EksctlDeleteCmd.WithArgs(
						"iamidentitymapping",
						"--name", params.ClusterName,
						"--arn", role1.ARN(),
						"--all",
					)
					Expect(deleteCmd).To(RunSuccessfully())

					getCmd := params.EksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--name", params.ClusterName,
						"--arn", role1.ARN(),
						"-o", "yaml",
					)
					Expect(getCmd).ToNot(RunSuccessfully())
				})
				It("deletes duplicate user mappings with --all", func() {
					deleteCmd := params.EksctlDeleteCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", user0.ARN(),
						"--all",
					)
					Expect(deleteCmd).To(RunSuccessfully())

					getCmd := params.EksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", params.ClusterName,
						"--arn", user0.ARN(),
						"-o", "yaml",
					)
					Expect(getCmd).ToNot(RunSuccessfully())
				})
			})

			Context("and delete the second nodegroup", func() {
				It("should not return an error", func() {
					cmd := params.EksctlDeleteCmd.WithArgs(
						"nodegroup",
						"--verbose", "4",
						"--cluster", params.ClusterName,
						testNG,
					)
					Expect(cmd).To(RunSuccessfully())
				})
			})
		})

		Context("and scale the initial nodegroup back to 1 node", func() {
			It("should not return an error", func() {
				cmd := params.EksctlScaleNodeGroupCmd.WithArgs(
					"--cluster", params.ClusterName,
					"--nodes", "1",
					"--name", initNG,
				)
				Expect(cmd).To(RunSuccessfully())
			})
		})

		Context("and deleting the cluster", func() {
			It("should not return an error", func() {
				if params.SkipDelete {
					Skip("will not delete cluster " + params.ClusterName)
				}

				cmd := params.EksctlDeleteClusterCmd.WithArgs(
					"--name", params.ClusterName,
				)
				Expect(cmd).To(RunSuccessfully())
			})
		})
	})
})

func truncate(clusterName string) string {
	// CloudFormation seems to truncate long cluster names at 37 characters:
	if len(clusterName) > 37 {
		return clusterName[:37]
	}
	return clusterName
}
