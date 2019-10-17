// +build integration

package integration_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	harness "github.com/dlespiau/kube-test-harness"
	"github.com/kubicorn/kubicorn/pkg/namer"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/iam"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/utils/file"
)

var _ = Describe("(Integration) Create, Get, Scale & Delete", func() {

	const (
		initNG = "ng-0"
		testNG = "ng-1"
	)

	commonTimeout := 10 * time.Minute

	BeforeSuite(func() {
		kubeconfigTemp = false
		if kubeconfigPath == "" {
			wd, _ := os.Getwd()
			f, _ := ioutil.TempFile(wd, "kubeconfig-")
			kubeconfigPath = f.Name()
			kubeconfigTemp = true
		}
	})

	AfterSuite(func() {
		gexec.KillAndWait()
		if kubeconfigTemp {
			os.Remove(kubeconfigPath)
		}
		os.RemoveAll(testDirectory)
	})

	Describe("when creating a cluster with 1 node", func() {
		It("should not return an error", func() {
			if !doCreate {
				fmt.Fprintf(GinkgoWriter, "will use existing cluster %s", clusterName)
				if !file.Exists(kubeconfigPath) {
					// Generate the Kubernetes configuration that eksctl create
					// would have generated otherwise:
					cmd := eksctlUtilsCmd.WithArgs(
						"write-kubeconfig",
						"--verbose", "4",
						"--cluster", clusterName,
						"--kubeconfig", kubeconfigPath,
					)
					Expect(cmd).To(RunSuccessfully())
				}
				return
			}

			fmt.Fprintf(GinkgoWriter, "Using kubeconfig: %s\n", kubeconfigPath)

			if clusterName == "" {
				clusterName = cmdutils.ClusterName("", "")
			}

			cmd := eksctlCreateCmd.WithArgs(
				"cluster",
				"--verbose", "4",
				"--name", clusterName,
				"--tags", "alpha.eksctl.io/description=eksctl integration test",
				"--nodegroup-name", initNG,
				"--node-labels", "ng-name="+initNG,
				"--nodes", "1",
				"--version", version,
				"--kubeconfig", kubeconfigPath,
			)
			Expect(cmd).To(RunSuccessfully())
		})

		It("should have created an EKS cluster and two CloudFormation stacks", func() {
			awsSession := NewSession(region)

			Expect(awsSession).To(HaveExistingCluster(clusterName, awseks.ClusterStatusActive, version))

			Expect(awsSession).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-cluster", clusterName)))
			Expect(awsSession).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%s", clusterName, initNG)))
		})

		It("should have created a valid kubectl config file", func() {
			config, err := clientcmd.LoadFromFile(kubeconfigPath)
			Expect(err).ShouldNot(HaveOccurred())

			err = clientcmd.ConfirmUsable(*config, "")
			Expect(err).ShouldNot(HaveOccurred())

			Expect(config.CurrentContext).To(ContainSubstring("eksctl"))
			Expect(config.CurrentContext).To(ContainSubstring(clusterName))
			Expect(config.CurrentContext).To(ContainSubstring(region))
		})

		Context("and listing clusters", func() {
			It("should return the previously created cluster", func() {
				cmd := eksctlGetCmd.WithArgs("clusters", "--all-regions")
				Expect(cmd).To(RunSuccessfullyWithOutputString(ContainSubstring(clusterName)))
			})
		})

		Context("and enabling gitops on the configured repo", func() {
			It("should not return an error", func() {
				// Use a random branch to ensure test runs don't step on each others.
				branch := namer.RandomName()
				cloneDir, err := createBranch(branch)
				Expect(err).ShouldNot(HaveOccurred())
				defer deleteBranch(branch, cloneDir)

				assertFluxManifestsAbsentInGit(branch)
				assertFluxPodsAbsentInKubernetes(kubeconfigPath)

				cmd := eksctlExperimentalCmd.WithArgs(
					"enable", "repo",
					"--git-url", Repository,
					"--git-email", Email,
					"--git-private-ssh-key-path", privateSSHKeyPath,
					"--git-branch", branch,
					"--cluster", clusterName,
				)
				Expect(cmd).To(RunSuccessfully())

				assertFluxManifestsPresentInGit(branch)
				assertFluxPodsPresentInKubernetes(kubeconfigPath)
			})
		})

		Context("enable profile", func() {
			It("should add quickstart to the repo and the cluster", func() {
				// Use a random branch to ensure test runs don't step on each others.
				branch := namer.RandomName()
				cloneDir, err := createBranch(branch)
				Expect(err).ShouldNot(HaveOccurred())
				defer deleteBranch(branch, cloneDir)

				assertFluxManifestsAbsentInGit(branch)

				deleteFluxInstallation(kubeconfigPath)
				assertFluxPodsAbsentInKubernetes(kubeconfigPath)

				cmd := eksctlExperimentalCmd.WithArgs(
					"enable", "profile",
					"--git-url", Repository,
					"--git-email", Email,
					"--git-branch", branch,
					"--git-private-ssh-key-path", privateSSHKeyPath,
					"--cluster", clusterName,
					"app-dev",
				)
				Expect(cmd).To(RunSuccessfully())

				assertQuickStartComponentsPresentInGit(branch)
				assertFluxManifestsPresentInGit(branch)
				assertFluxPodsPresentInKubernetes(kubeconfigPath)
			})
		})

		Context("and scale the initial nodegroup", func() {
			It("should not return an error", func() {
				cmd := eksctlScaleNodeGroupCmd.WithArgs(
					"--cluster", clusterName,
					"--nodes", "4",
					"--name", initNG,
				)
				Expect(cmd).To(RunSuccessfully())
			})

			It("{FLAKY: https://github.com/weaveworks/eksctl/issues/717} should make it 4 nodes total", func() {
				test, err := newKubeTest()
				Expect(err).ShouldNot(HaveOccurred())
				defer test.Close()

				test.WaitForNodesReady(4, commonTimeout)

				nodes := test.ListNodes((metav1.ListOptions{
					LabelSelector: api.NodeGroupNameLabel + "=" + initNG,
				}))

				Expect(len(nodes.Items)).To(Equal(4))
			})
		})

		Context("and add the second nodegroup", func() {
			It("should not return an error", func() {
				cmd := eksctlCreateCmd.WithArgs(
					"nodegroup",
					"--cluster", clusterName,
					"--nodes", "4",
					"--node-private-networking",
					testNG,
				)
				Expect(cmd).To(RunSuccessfully())
			})

			It("should be able to list nodegroups", func() {
				{
					cmd := eksctlGetCmd.WithArgs(
						"nodegroup",
						"--cluster", clusterName,
						initNG,
					)
					Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
						ContainElement(ContainSubstring(initNG)),
						Not(ContainElement(ContainSubstring(testNG))),
					))
				}
				{
					cmd := eksctlGetCmd.WithArgs(
						"nodegroup",
						"--cluster", clusterName,
						testNG,
					)
					Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
						ContainElement(ContainSubstring(testNG)),
						Not(ContainElement(ContainSubstring(initNG))),
					))
				}
				{
					cmd := eksctlGetCmd.WithArgs(
						"nodegroup",
						"--cluster", clusterName,
					)
					Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
						ContainElement(ContainSubstring(testNG)),
						ContainElement(ContainSubstring(initNG)),
					))
				}
			})

			It("{FLAKY: https://github.com/weaveworks/eksctl/issues/717} should make it 8 nodes total", func() {
				test, err := newKubeTest()
				Expect(err).ShouldNot(HaveOccurred())
				defer test.Close()

				test.WaitForNodesReady(8, commonTimeout)

				nodes := test.ListNodes(metav1.ListOptions{})

				Expect(len(nodes.Items)).To(Equal(8))
			})

			Context("toggle CloudWatch logging", func() {
				var (
					cfg *api.ClusterConfig
					ctl *eks.ClusterProvider
				)

				BeforeEach(func() {
					cfg = &api.ClusterConfig{
						Metadata: &api.ClusterMeta{
							Name:   clusterName,
							Region: region,
						},
					}
					ctl = eks.New(&api.ProviderConfig{Region: region}, cfg)
				})

				It("should have all types disabled by default", func() {
					enabled, disable, err := ctl.GetCurrentClusterConfigForLogging(cfg)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(enabled.List()).To(HaveLen(0))
					Expect(disable.List()).To(HaveLen(5))
				})

				It("should plan to enable two of the types using flags", func() {
					cmd := eksctlUtilsCmd.WithArgs(
						"update-cluster-logging",
						"--cluster", clusterName,
						"--enable-types", "api,controllerManager",
					)
					Expect(cmd).To(RunSuccessfully())

					enabled, disable, err := ctl.GetCurrentClusterConfigForLogging(cfg)
					Expect(err).ShouldNot(HaveOccurred())
					Expect(enabled.List()).To(HaveLen(0))
					Expect(disable.List()).To(HaveLen(5))
				})

				It("should enable two of the types using flags", func() {
					cmd := eksctlUtilsCmd.WithArgs(
						"update-cluster-logging",
						"--cluster", clusterName,
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
					cmd := eksctlUtilsCmd.WithArgs(
						"update-cluster-logging",
						"--cluster", clusterName,
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
					cmd := eksctlUtilsCmd.WithArgs(
						"update-cluster-logging",
						"--cluster", clusterName,
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
					cmd := eksctlUtilsCmd.WithArgs(
						"update-cluster-logging",
						"--cluster", clusterName,
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
					cmd := eksctlUtilsCmd.WithArgs(
						"update-cluster-logging",
						"--cluster", clusterName,
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
							Name:   clusterName,
							Region: region,
						},
					}
					ctl = eks.New(&api.ProviderConfig{Region: region}, cfg)
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

					setupCmd := eksctlUtilsCmd.WithArgs(
						"associate-iam-oidc-provider",
						"--cluster", clusterName,
						"--approve",
					)
					Expect(setupCmd).To(RunSuccessfully())

					{
						exists, err := oidc.CheckProviderExists()
						Expect(err).ShouldNot(HaveOccurred())
						Expect(exists).To(BeTrue())
					}

					cmds := []Cmd{
						eksctlCreateCmd.WithArgs(
							"iamserviceaccount",
							"--cluster", clusterName,
							"--name", "app-cache-access",
							"--namespace", "app1",
							"--attach-policy-arn", "arn:aws:iam::aws:policy/AmazonDynamoDBReadOnlyAccess",
							"--attach-policy-arn", "arn:aws:iam::aws:policy/AmazonElastiCacheFullAccess",
							"--approve",
						),
						eksctlCreateCmd.WithArgs(
							"iamserviceaccount",
							"--cluster", clusterName,
							"--name", "s3-read-only",
							"--attach-policy-arn", "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess",
							"--approve",
						),
					}

					Expect(cmds).To(RunSuccessfully())

					awsSession := NewSession(region)

					stackNamePrefix := fmt.Sprintf("eksctl-%s-addon-iamserviceaccount-", clusterName)

					Expect(awsSession).To(HaveExistingStack(stackNamePrefix + "default-s3-read-only"))
					Expect(awsSession).To(HaveExistingStack(stackNamePrefix + "app1-app-cache-access"))

					clientSet, err := ctl.NewStdClientSet(cfg)
					Expect(err).ShouldNot(HaveOccurred())

					{
						sa, err := clientSet.CoreV1().ServiceAccounts(metav1.NamespaceDefault).Get("s3-read-only", metav1.GetOptions{})
						Expect(err).ShouldNot(HaveOccurred())

						Expect(sa.Annotations).To(HaveLen(1))
						Expect(sa.Annotations).To(HaveKey(api.AnnotationEKSRoleARN))
						Expect(sa.Annotations[api.AnnotationEKSRoleARN]).To(MatchRegexp("^arn:aws:iam::.*:role/eksctl-" + clusterName + ".*$"))
					}

					{
						sa, err := clientSet.CoreV1().ServiceAccounts("app1").Get("app-cache-access", metav1.GetOptions{})
						Expect(err).ShouldNot(HaveOccurred())

						Expect(sa.Annotations).To(HaveLen(1))
						Expect(sa.Annotations).To(HaveKey(api.AnnotationEKSRoleARN))
						Expect(sa.Annotations[api.AnnotationEKSRoleARN]).To(MatchRegexp("^arn:aws:iam::.*:role/eksctl-" + clusterName + ".*$"))
					}
				})

				It("should list both iamserviceaccounts", func() {
					cmd := eksctlGetCmd.WithArgs(
						"iamserviceaccount",
						"--cluster", clusterName,
					)

					Expect(cmd).To(RunSuccessfullyWithOutputString(MatchRegexp(
						`(?m:^NAMESPACE\s+NAME\s+ROLE\sARN$)` +
							`|(?m:^app1\s+app-cache-access\s+arn:aws:iam::.*$)` +
							`|(?m:^default\s+s3-read-only\s+arn:aws:iam::.*$)`,
					)))
				})

				It("delete both iamserviceaccounts", func() {
					cmds := []Cmd{
						eksctlDeleteCmd.WithArgs(
							"iamserviceaccount",
							"--cluster", clusterName,
							"--name", "s3-read-only",
							"--wait",
						),
						eksctlDeleteCmd.WithArgs(
							"iamserviceaccount",
							"--cluster", clusterName,
							"--name", "app-cache-access",
							"--namespace", "app1",
							"--wait",
						),
					}
					Expect(cmds).To(RunSuccessfully())

					awsSession := NewSession(region)

					stackNamePrefix := fmt.Sprintf("eksctl-%s-addon-iamserviceaccount-", clusterName)

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
					test, err = newKubeTest()
					Expect(err).ShouldNot(HaveOccurred())
				})

				AfterEach(func() {
					test.Close()
					Eventually(func() int {
						return len(test.ListPods(test.Namespace, metav1.ListOptions{}).Items)
					}, "3m", "1s").Should(BeZero())
				})

				It("should deploy podinfo service to the cluster and access it via proxy", func() {
					d := test.CreateDeploymentFromFile(test.Namespace, "podinfo.yaml")
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
					d := test.CreateDaemonSetFromFile(test.Namespace, "test-dns.yaml")

					test.WaitForDaemonSetReady(d, commonTimeout)

					{
						ds, err := test.GetDaemonSet(test.Namespace, d.Name)
						Expect(err).ShouldNot(HaveOccurred())
						fmt.Fprintf(GinkgoWriter, "ds.Status = %#v", ds.Status)
					}
				})

				It("should have access to HTTP(S) sites", func() {
					d := test.CreateDaemonSetFromFile(test.Namespace, "test-http.yaml")

					test.WaitForDaemonSetReady(d, commonTimeout)

					{
						ds, err := test.GetDaemonSet(test.Namespace, d.Name)
						Expect(err).ShouldNot(HaveOccurred())
						fmt.Fprintf(GinkgoWriter, "ds.Status = %#v", ds.Status)
					}
				})

				It("should be able to run pods with an iamserviceaccount", func() {
					createCmd := eksctlCreateCmd.WithArgs(
						"iamserviceaccount",
						"--cluster", clusterName,
						"--name", "s3-reader",
						"--namespace", test.Namespace,
						"--attach-policy-arn", "arn:aws:iam::aws:policy/AmazonS3ReadOnlyAccess",
						"--approve",
					)

					Expect(createCmd).To(RunSuccessfully())

					d := test.CreateDeploymentFromFile(test.Namespace, "iamserviceaccount-checker.yaml")
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

						Expect(so.AssumedRoleUser.Arn).To(MatchRegexp("^arn:aws:sts::.*:assumed-role/eksctl-" + clusterName + "-.*/integration-test$"))

						Expect(so.Audience).To(Equal("sts.amazonaws.com"))

						Expect(so.Provider).To(MatchRegexp("^arn:aws:iam::.*:oidc-provider/oidc.eks." + region + ".amazonaws.com/id/.*$"))

						Expect(so.SubjectFromWebIdentityToken).To(Equal("system:serviceaccount:" + test.Namespace + ":s3-reader"))

						Expect(so.Credentials.SecretAccessKey).ToNot(BeEmpty())
						Expect(so.Credentials.SessionToken).ToNot(BeEmpty())
						Expect(so.Credentials.Expiration).ToNot(BeEmpty())
						Expect(so.Credentials.AccessKeyId).ToNot(BeEmpty())
					}

					deleteCmd := eksctlDeleteCmd.WithArgs(
						"iamserviceaccount",
						"--cluster", clusterName,
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
					cmd := eksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", clusterName,
						"--arn", "arn:aws:iam::123456:role/idontexist",
						"-o", "yaml",
					)
					Expect(cmd).ToNot(RunSuccessfully())
				})
				It("fails getting unknown user mapping", func() {
					cmd := eksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", clusterName,
						"--arn", "arn:aws:iam::123456:user/bob",
						"-o", "yaml",
					)
					Expect(cmd).ToNot(RunSuccessfully())
				})
				It("creates role mapping", func() {
					create := eksctlCreateCmd.WithArgs(
						"iamidentitymapping",
						"--name", clusterName,
						"--arn", role0.ARN(),
						"--username", role0.Username(),
						"--group", role0.Groups()[0],
						"--group", role0.Groups()[1],
					)
					Expect(create).To(RunSuccessfully())

					get := eksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--name", clusterName,
						"--arn", role0.ARN(),
						"-o", "yaml",
					)
					Expect(get).To(RunSuccessfullyWithOutputString(MatchYAML(expR0)))
				})
				It("creates user mapping", func() {
					create := eksctlCreateCmd.WithArgs(
						"iamidentitymapping",
						"--name", clusterName,
						"--arn", user0.ARN(),
						"--username", user0.Username(),
						"--group", user0.Groups()[0],
						"--group", user0.Groups()[1],
					)
					Expect(create).To(RunSuccessfully())

					get := eksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", clusterName,
						"--arn", user0.ARN(),
						"-o", "yaml",
					)
					Expect(get).To(RunSuccessfullyWithOutputString(MatchYAML(expU0)))
				})
				It("creates a duplicate role mapping", func() {
					createRole := eksctlCreateCmd.WithArgs(
						"iamidentitymapping",
						"--name", clusterName,
						"--arn", role0.ARN(),
						"--username", role0.Username(),
						"--group", role0.Groups()[0],
						"--group", role0.Groups()[1],
					)
					Expect(createRole).To(RunSuccessfully())

					get := eksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--name", clusterName,
						"--arn", role0.ARN(),
						"-o", "yaml",
					)
					Expect(get).To(RunSuccessfullyWithOutputString(MatchYAML(expR0 + expR0)))
				})
				It("creates a duplicate user mapping", func() {
					createCmd := eksctlCreateCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", clusterName,
						"--arn", user0.ARN(),
						"--username", user0.Username(),
						"--group", user0.Groups()[0],
						"--group", user0.Groups()[1],
					)
					Expect(createCmd).To(RunSuccessfully())

					getCmd := eksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", clusterName,
						"--arn", user0.ARN(),
						"-o", "yaml",
					)
					Expect(getCmd).To(RunSuccessfullyWithOutputString(MatchYAML(expU0 + expU0)))
				})
				It("creates a duplicate role mapping with different identity", func() {
					createCmd := eksctlCreateCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", clusterName,
						"--arn", role1.ARN(),
						"--group", role1.Groups()[0],
					)
					Expect(createCmd).To(RunSuccessfully())

					getCmd := eksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", clusterName,
						"--arn", role1.ARN(),
						"-o", "yaml",
					)
					Expect(getCmd).To(RunSuccessfullyWithOutputString(MatchYAML(expR0 + expR0 + expR1)))
				})
				It("deletes a single role mapping fifo", func() {
					deleteCmd := eksctlDeleteCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", clusterName,
						"--arn", role1.ARN(),
					)
					Expect(deleteCmd).To(RunSuccessfully())

					getCmd := eksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", clusterName,
						"--arn", role1.ARN(),
						"-o", "yaml",
					)
					Expect(getCmd).To(RunSuccessfullyWithOutputString(MatchYAML(expR0 + expR1)))
				})
				It("fails when deleting unknown mapping", func() {
					deleteCmd := eksctlDeleteCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", clusterName,
						"--arn", "arn:aws:iam::123456:role/idontexist",
					)
					Expect(deleteCmd).ToNot(RunSuccessfully())
				})
				It("deletes duplicate role mappings with --all", func() {
					deleteCmd := eksctlDeleteCmd.WithArgs(
						"iamidentitymapping",
						"--name", clusterName,
						"--arn", role1.ARN(),
						"--all",
					)
					Expect(deleteCmd).To(RunSuccessfully())

					getCmd := eksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--name", clusterName,
						"--arn", role1.ARN(),
						"-o", "yaml",
					)
					Expect(getCmd).ToNot(RunSuccessfully())
				})
				It("deletes duplicate user mappings with --all", func() {
					deleteCmd := eksctlDeleteCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", clusterName,
						"--arn", user0.ARN(),
						"--all",
					)
					Expect(deleteCmd).To(RunSuccessfully())

					getCmd := eksctlGetCmd.WithArgs(
						"iamidentitymapping",
						"--cluster", clusterName,
						"--arn", user0.ARN(),
						"-o", "yaml",
					)
					Expect(getCmd).ToNot(RunSuccessfully())
				})
			})

			Context("and delete the second nodegroup", func() {
				It("should not return an error", func() {
					cmd := eksctlDeleteCmd.WithArgs(
						"nodegroup",
						"--verbose", "4",
						"--cluster", clusterName,
						testNG,
					)
					Expect(cmd).To(RunSuccessfully())
				})

				It("{FLAKY: https://github.com/weaveworks/eksctl/issues/717} should make it 4 nodes total", func() {
					test, err := newKubeTest()
					Expect(err).ShouldNot(HaveOccurred())
					defer test.Close()

					test.WaitForNodesReady(4, commonTimeout)

					nodes := test.ListNodes(metav1.ListOptions{
						LabelSelector: api.NodeGroupNameLabel + "=" + initNG,
					})
					allNodes := test.ListNodes(metav1.ListOptions{})
					Expect(len(nodes.Items)).To(Equal(4))
					Expect(len(allNodes.Items)).To(Equal(4))
				})
			})
		})

		Context("and scale the initial nodegroup back to 1 node", func() {
			It("should not return an error", func() {
				cmd := eksctlScaleNodeGroupCmd.WithArgs(
					"--cluster", clusterName,
					"--nodes", "1",
					"--name", initNG,
				)
				Expect(cmd).To(RunSuccessfully())
			})

			It("{FLAKY: https://github.com/weaveworks/eksctl/issues/717} should make it 1 nodes total", func() {
				test, err := newKubeTest()
				Expect(err).ShouldNot(HaveOccurred())
				defer test.Close()

				test.WaitForNodesReady(1, commonTimeout)

				nodes := test.ListNodes(metav1.ListOptions{
					LabelSelector: api.NodeGroupNameLabel + "=" + initNG,
				})

				Expect(len(nodes.Items)).To(Equal(1))
			})
		})

		Context("and deleting the cluster", func() {

			It("{FLAKY: https://github.com/weaveworks/eksctl/issues/536} should not return an error", func() {
				if !doDelete {
					Skip("will not delete cluster " + clusterName)
				}

				cmd := eksctlDeleteClusterCmd.WithArgs(
					"--name", clusterName,
					"--wait",
				)
				Expect(cmd).To(RunSuccessfully())
			})

			It("{FLAKY: https://github.com/weaveworks/eksctl/issues/536} should have deleted the EKS cluster and both CloudFormation stacks", func() {
				if !doDelete {
					Skip("will not delete cluster " + clusterName)
				}

				awsSession := NewSession(region)

				Expect(awsSession).ToNot(HaveExistingCluster(clusterName, awseks.ClusterStatusActive, version))

				Expect(awsSession).ToNot(HaveExistingStack(fmt.Sprintf("eksctl-%s-cluster", clusterName)))
				Expect(awsSession).ToNot(HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-ng-%d", clusterName, 0)))
			})
		})
	})
})
