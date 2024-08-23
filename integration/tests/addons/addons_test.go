//go:build integration
// +build integration

package addons

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfn "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	k8sslices "k8s.io/utils/strings/slices"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/integration/matchers"
	"github.com/weaveworks/eksctl/integration/runner"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	"github.com/weaveworks/eksctl/pkg/actions/podidentityassociation"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var (
	params      *tests.Params
	rawClient   *kubewrapper.RawClient
	awsProvider api.ClusterProvider
)

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("addons")
}

func TestEKSAddons(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = BeforeSuite(func() {
	data, err := json.Marshal(getInitialClusterConfig())
	Expect(err).NotTo(HaveOccurred())

	cmd := params.EksctlCreateCmd.
		WithArgs(
			"cluster",
			"--config-file", "-",
			"--verbose", "4",
		).
		WithoutArg("--region", params.Region).
		WithStdin(bytes.NewReader(data))
	Expect(cmd).To(RunSuccessfully())

	rawClient = getRawClient(context.Background(), params.ClusterName)
	serverVersion, err := rawClient.ServerVersion()
	Expect(err).NotTo(HaveOccurred())
	Expect(serverVersion).To(HavePrefix(api.LatestVersion))

})

var _ = Describe("(Integration) [EKS Addons test]", func() {

	Context("Creating a cluster with addons", func() {
		clusterName := params.ClusterName

		It("should support addons", func() {
			By("Asserting the addon is listed in `get addons`")
			Eventually(func() runner.Cmd {
				cmd := params.EksctlGetCmd.
					WithArgs(
						"addons",
						"--cluster", clusterName,
						"--verbose", "2",
					)
				return cmd
			}, "5m", "30s").Should(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring("vpc-cni")),
			))

			By("Asserting the addons are healthy")
			Eventually(func() runner.Cmd {
				cmd := params.EksctlGetCmd.
					WithArgs(
						"addon",
						"--name", "vpc-cni",
						"--cluster", clusterName,
						"--verbose", "2",
					)
				return cmd
			}, "5m", "30s").Should(RunSuccessfullyWithOutputStringLines(ContainElement(ContainSubstring("ACTIVE"))))

			By("Deleting the kube-proxy addon")
			cmd := params.EksctlDeleteCmd.
				WithArgs(
					"addon",
					"--name", "kube-proxy",
					"--cluster", clusterName,
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfully())

			By("Deleting the vpc-cni addon with --preserve")
			cmd = params.EksctlDeleteCmd.
				WithArgs(
					"addon",
					"--name", "vpc-cni",
					"--preserve",
					"--cluster", clusterName,
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfully())
			_, err := rawClient.ClientSet().AppsV1().DaemonSets("kube-system").Get(context.Background(), "aws-node", metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should have full control over configMap when creating addons", func() {
			clusterConfig := getInitialClusterConfig()
			clusterConfig.Addons = []*api.Addon{
				{
					Name:    "coredns",
					Version: "latest",
				},
			}
			cmd := params.EksctlCreateCmd.
				WithArgs(
					"addon",
					"--config-file", "-",
				).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.Reader(clusterConfig))
			Expect(cmd).To(RunSuccessfully())

			By("deleting coredns but preserving its resources")
			cmd = params.EksctlDeleteCmd.
				WithArgs(
					"addon",
					"--cluster", clusterConfig.Metadata.Name,
					"--name", "coredns",
					"--verbose", "4",
					"--preserve",
					"--region", params.Region,
				)
			Expect(cmd).To(RunSuccessfully())

			Eventually(func() runner.Cmd {
				return params.EksctlGetCmd.
					WithArgs(
						"addon",
						"--name", "coredns",
						"--cluster", clusterName,
						"--verbose", "4",
					)
			}, "5m", "30s").ShouldNot(RunSuccessfully())

			configMap := getConfigMap(rawClient.ClientSet(), "coredns")
			oldCacheValue := getCacheValue(configMap)
			newCacheValue := addToString(oldCacheValue, 5)
			updateCacheValue(configMap, oldCacheValue, newCacheValue)
			updateConfigMap(rawClient.ClientSet(), configMap)

			By("erroring when there are config conflicts")
			clusterConfig = getInitialClusterConfig()
			clusterConfig.Addons = []*api.Addon{
				{
					Name:             "coredns",
					Version:          "latest",
					ResolveConflicts: ekstypes.ResolveConflictsNone,
				},
			}
			data, err := json.Marshal(clusterConfig)
			Expect(err).NotTo(HaveOccurred())

			cmd = params.EksctlCreateCmd.
				WithArgs(
					"addon",
					"--config-file", "-",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))
			Expect(cmd).NotTo(RunSuccessfully())

			Eventually(func() runner.Cmd {
				cmd := params.EksctlGetCmd.
					WithArgs(
						"addon",
						"--name", "coredns",
						"--cluster", clusterName,
						"--verbose", "2",
					)
				return cmd
			}, "5m", "30s").Should(RunSuccessfullyWithOutputStringLines(ContainElement(ContainSubstring("CREATE_FAILED"))))

			By("overwriting the configMap")
			clusterConfig = getInitialClusterConfig()
			clusterConfig.Addons = []*api.Addon{
				{
					Name:             "coredns",
					Version:          "latest",
					ResolveConflicts: ekstypes.ResolveConflictsOverwrite,
				},
			}
			data, err = json.Marshal(clusterConfig)
			Expect(err).NotTo(HaveOccurred())

			cmd = params.EksctlCreateCmd.
				WithArgs(
					"addon",
					"--config-file", "-",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))
			Expect(cmd).To(RunSuccessfully())

			Eventually(func() runner.Cmd {
				cmd := params.EksctlGetCmd.
					WithArgs(
						"addon",
						"--name", "coredns",
						"--cluster", clusterName,
						"--verbose", "2",
					)
				return cmd
			}, "5m", "30s").Should(RunSuccessfullyWithOutputStringLines(ContainElement(ContainSubstring("ACTIVE"))))

			Expect(getCacheValue(getConfigMap(rawClient.ClientSet(), "coredns"))).To(Equal(oldCacheValue))
		})

		It("should have full control over configMap when updating addons", func() {
			var (
				clusterConfig *api.ClusterConfig
				configMap     *corev1.ConfigMap
			)

			configMap = getConfigMap(rawClient.ClientSet(), "coredns")
			oldCacheValue := getCacheValue(configMap)
			newCacheValue := addToString(oldCacheValue, 5)
			updateCacheValue(configMap, oldCacheValue, newCacheValue)
			updateConfigMap(rawClient.ClientSet(), configMap)

			By("preserving the configMap")
			clusterConfig = getInitialClusterConfig()
			clusterConfig.Addons = []*api.Addon{
				{
					Name:             "coredns",
					Version:          "latest",
					ResolveConflicts: ekstypes.ResolveConflictsPreserve,
				},
			}

			data, err := json.Marshal(clusterConfig)
			Expect(err).NotTo(HaveOccurred())

			cmd := params.EksctlUpdateCmd.
				WithArgs(
					"addon",
					"--config-file", "-",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))
			Expect(cmd).To(RunSuccessfully())

			Eventually(func() runner.Cmd {
				cmd := params.EksctlGetCmd.
					WithArgs(
						"addon",
						"--name", "coredns",
						"--cluster", clusterName,
						"--verbose", "2",
					)
				return cmd
			}, "5m", "30s").Should(RunSuccessfullyWithOutputStringLines(ContainElement(ContainSubstring("ACTIVE"))))

			Expect(getCacheValue(getConfigMap(rawClient.ClientSet(), "coredns"))).To(Equal(newCacheValue))

			By("erroring when there are config conflicts")
			cmd = params.EksctlUpdateCmd.
				WithArgs(
					"addon",
					"--name", "coredns",
					"--cluster", clusterName,
					"--wait",
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfully())

			Eventually(func() runner.Cmd {
				cmd := params.EksctlGetCmd.
					WithArgs(
						"addon",
						"--name", "coredns",
						"--cluster", clusterName,
						"--verbose", "2",
					)
				return cmd
			}, "5m", "30s").Should(RunSuccessfullyWithOutputStringLines(ContainElement(ContainSubstring("UPDATE_FAILED"))))

			Expect(getCacheValue(getConfigMap(rawClient.ClientSet(), "coredns"))).To(Equal(newCacheValue))

			By("overwriting the configMap")
			clusterConfig = getInitialClusterConfig()
			clusterConfig.Addons = []*api.Addon{
				{
					Name:             "coredns",
					Version:          "latest",
					ResolveConflicts: ekstypes.ResolveConflictsOverwrite,
				},
			}

			data, err = json.Marshal(clusterConfig)
			Expect(err).NotTo(HaveOccurred())

			cmd = params.EksctlUpdateCmd.
				WithArgs(
					"addon",
					"--config-file", "-",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))
			Expect(cmd).To(RunSuccessfully())

			Eventually(func() runner.Cmd {
				cmd := params.EksctlGetCmd.
					WithArgs(
						"addon",
						"--name", "coredns",
						"--cluster", clusterName,
						"--verbose", "2",
					)
				return cmd
			}, "5m", "30s").Should(RunSuccessfullyWithOutputStringLines(ContainElement(ContainSubstring("ACTIVE"))))

			Expect(getCacheValue(getConfigMap(rawClient.ClientSet(), "coredns"))).To(Equal(oldCacheValue))
		})

		It("should support advanced addon configuration", func() {
			cmd := params.EksctlDeleteCmd.
				WithArgs(
					"addon",
					"--name", "coredns",
					"--cluster", clusterName,
					"--verbose", "2",
					"--region", params.Region,
				)
			Expect(cmd).To(RunSuccessfully())

			Eventually(func() runner.Cmd {
				cmd := params.EksctlGetCmd.
					WithArgs(
						"addon",
						"--name", "coredns",
						"--cluster", clusterName,
						"--verbose", "2",
						"--region", params.Region,
					)
				return cmd
			}, "5m", "30s").ShouldNot(RunSuccessfully())

			By("successfully creating an addon with configuration values")
			clusterConfig := getInitialClusterConfig()
			clusterConfig.Addons = []*api.Addon{
				{
					Name:                api.CoreDNSAddon,
					ConfigurationValues: "{\"replicaCount\":3}",
					ResolveConflicts:    ekstypes.ResolveConflictsOverwrite,
				},
			}
			data, err := json.Marshal(clusterConfig)

			Expect(err).NotTo(HaveOccurred())
			cmd = params.EksctlCreateCmd.
				WithArgs(
					"addon",
					"--config-file", "-",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))
			Expect(cmd).To(RunSuccessfully())

			Eventually(func() runner.Cmd {
				cmd := params.EksctlGetCmd.
					WithArgs(
						"addon",
						"--name", "coredns",
						"--cluster", clusterName,
						"--verbose", "2",
						"--region", params.Region,
					)
				return cmd
			}, "5m", "30s").Should(RunSuccessfullyWithOutputStringLines(ContainElement(ContainSubstring("{\"replicaCount\":3}"))))

			By("successfully updating the configuration values of the addon")
			clusterConfig = getInitialClusterConfig()
			clusterConfig.Addons = []*api.Addon{
				{
					Name:                api.CoreDNSAddon,
					ConfigurationValues: "{\"replicaCount\":3, \"computeType\":\"test\"}",
					ResolveConflicts:    ekstypes.ResolveConflictsOverwrite,
				},
			}
			data, err = json.Marshal(clusterConfig)

			Expect(err).NotTo(HaveOccurred())
			cmd = params.EksctlUpdateCmd.
				WithArgs(
					"addon",
					"--config-file", "-",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))
			Expect(cmd).To(RunSuccessfully())

			Eventually(func() runner.Cmd {
				cmd := params.EksctlGetCmd.
					WithArgs(
						"addon",
						"--name", "coredns",
						"--cluster", clusterName,
						"--verbose", "2",
						"--region", params.Region,
					)
				return cmd
			}, "5m", "30s").Should(RunSuccessfullyWithOutputStringLines(ContainElement(ContainSubstring("{\"replicaCount\":3, \"computeType\":\"test\"}"))))
		})
	})

	It("should describe addons", func() {
		cmd := params.EksctlUtilsCmd.
			WithArgs(
				"describe-addon-versions",
				"--kubernetes-version", api.LatestVersion,
			)
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(ContainSubstring("vpc-cni")),
		))
	})

	It("should output the configuration schema for addons", func() {
		addonWithSchema := "coredns"
		cfg := NewConfig(params.Region)
		eksAPI := awseks.NewFromConfig(cfg)
		By(fmt.Sprintf("listing available addon versions for %s", addonWithSchema))
		output, err := eksAPI.DescribeAddonVersions(context.Background(), &awseks.DescribeAddonVersionsInput{
			AddonName:         aws.String(addonWithSchema),
			KubernetesVersion: aws.String(api.LatestVersion),
		})
		Expect(err).NotTo(HaveOccurred(), "error describing addon versions")
		By(fmt.Sprintf("fetching the configuration schema for %s", addonWithSchema))
		Expect(output.Addons).NotTo(BeEmpty(), "expected to find addon versions for %s", addonWithSchema)
		addonVersions := output.Addons[0].AddonVersions
		Expect(addonVersions).NotTo(BeEmpty(), "expected to find at least one addon version")

		cmd := params.EksctlUtilsCmd.
			WithArgs(
				"describe-addon-configuration",
				"--name", addonWithSchema,
				"--version", *addonVersions[0].AddonVersion,
			)
		session := cmd.Run()
		Expect(session.ExitCode()).To(Equal(0))
		Expect(json.Valid(session.Buffer().Contents())).To(BeTrue(), "invalid JSON for configuration schema")
	})

	It("should describe addons when publisher, type and owner are supplied", func() {
		cmd := params.EksctlUtilsCmd.
			WithArgs(
				"describe-addon-versions",
				"--kubernetes-version", api.LatestVersion,
				"--types", "networking",
				"--owners", "aws",
				"--publishers", "eks",
			)
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(ContainSubstring("vpc-cni")),
			ContainElement(ContainSubstring("coredns")),
			ContainElement(ContainSubstring("kube-proxy")),
		))
	})

	It("should describe addons when multiple types are supplied", func() {
		cmd := params.EksctlUtilsCmd.
			WithArgs(
				"describe-addon-versions",
				"--kubernetes-version", api.LatestVersion,
				"--types", "networking, storage",
			)
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(ContainSubstring("vpc-cni")),
			ContainElement(ContainSubstring("aws-ebs-csi-driver")),
		))
	})

	Context("configure IAM permissions via pod identity associations or IRSA", Ordered, func() {
		const (
			pollInterval       = 10  //seconds
			timeOutSeconds     = 600 // 10 minutes
			awsNodeSA          = "aws-node"
			ebsCSIControllerSA = "ebs-csi-controller-sa"
			efsCSIControllerSA = "efs-csi-controller-sa"
		)

		clusterConfig := getInitialClusterConfig()
		cfg := NewConfig(clusterConfig.Metadata.Region)
		makeIRSAStackName := func(addonName string) string {
			return fmt.Sprintf("eksctl-%s-addon-%s", clusterConfig.Metadata.Name, addonName)
		}
		makePodIDStackName := func(addonName, serviceAccountName string) string {
			return podidentityassociation.MakeAddonPodIdentityStackName(clusterConfig.Metadata.Name, addonName, serviceAccountName)
		}
		makeCreateAddonCMD := func() runner.Cmd {
			return params.EksctlCreateCmd.
				WithArgs("addon").
				WithArgs("--config-file", "-").
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.Reader(clusterConfig))
		}
		makeUpdateAddonCMD := func(args ...string) runner.Cmd {
			cmd := params.EksctlUpdateCmd.
				WithArgs("addon").
				WithArgs("--config-file", "-").
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.Reader(clusterConfig))
			if len(args) > 0 {
				return cmd.WithArgs(args...)
			}
			return cmd
		}
		makeDeleteAddonCMD := func(addonName string, args ...string) runner.Cmd {
			cmd := params.EksctlDeleteCmd.WithArgs(
				"addon",
				"--cluster", clusterConfig.Metadata.Name,
				"--name", addonName,
			)
			if len(args) > 0 {
				return cmd.WithArgs(args...)
			}
			return cmd
		}
		assertStackExists := func(stackName string) {
			Expect(cfg).To(HaveExistingStack(stackName))
		}
		assertStackNotExists := func(stackName string) {
			Expect(cfg).NotTo(HaveExistingStack(stackName))
		}
		assertStackDeleted := func(stackName string) {
			Eventually(cfg, timeOutSeconds, pollInterval).ShouldNot(HaveExistingStack(stackName))
		}
		assertPodIDPresence := func(namespace, serviceAccountName string, expectPodIDExists bool) {
			cmd := params.EksctlGetCmd.WithArgs(
				"podidentityassociation",
				"--cluster", clusterConfig.Metadata.Name,
				"--namespace", namespace,
				"--service-account-name", serviceAccountName,
			)
			matcher := ContainElement("No podidentityassociations found")
			if expectPodIDExists {
				matcher = ContainElements(
					ContainSubstring(namespace),
					ContainSubstring(serviceAccountName),
				)
				cmd = cmd.WithArgs("--output", "json")
			}
			Expect(cmd).To(RunSuccessfullyWithOutputStringLines(matcher))
		}
		assertAddonHasPodIDs := func(addonName string, podIDsCount int) {
			addon, err := awsProvider.EKS().DescribeAddon(context.Background(), &awseks.DescribeAddonInput{
				AddonName:   aws.String(addonName),
				ClusterName: aws.String(clusterConfig.Metadata.Name),
			})
			ExpectWithOffset(1, err).NotTo(HaveOccurred())
			ExpectWithOffset(1, addon.Addon.PodIdentityAssociations).To(HaveLen(podIDsCount))
		}

		BeforeAll(func() {
			clusterConfig.Addons = []*api.Addon{{Name: api.PodIdentityAgentAddon}}
			Expect(makeCreateAddonCMD()).To(RunSuccessfully())
		})

		It("should provide IAM permissions when creating addons", func() {

			By("creating pod identity associations for addons when explicitly set by user")
			clusterConfig.Addons = []*api.Addon{
				{
					Name: api.VPCCNIAddon,
					PodIdentityAssociations: &[]api.PodIdentityAssociation{
						{
							Namespace:            "kube-system",
							ServiceAccountName:   awsNodeSA,
							PermissionPolicyARNs: []string{"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"},
						},
					},
				},
			}
			Expect(makeCreateAddonCMD()).To(RunSuccessfully())
			assertAddonHasPodIDs(api.VPCCNIAddon, 1)
			assertStackExists(makePodIDStackName(api.VPCCNIAddon, awsNodeSA))

			By("creating pod identity associations for addons when `autoApplyPodIdentityAssociations: true` and addon supports podIDs")
			clusterConfig.AddonsConfig.AutoApplyPodIdentityAssociations = true
			clusterConfig.Addons = []*api.Addon{{Name: api.AWSEBSCSIDriverAddon}}
			Expect(makeCreateAddonCMD()).To(RunSuccessfully())
			assertAddonHasPodIDs(api.AWSEBSCSIDriverAddon, 1)
			assertStackExists(makePodIDStackName(api.AWSEBSCSIDriverAddon, ebsCSIControllerSA))
			assertStackNotExists(makeIRSAStackName(api.AWSEBSCSIDriverAddon))

			By("falling back to IRSA when `autoApplyPodIdentityAssociations: true` but addon doesn't support podIDs")
			clusterConfig.Addons = []*api.Addon{{Name: api.AWSEFSCSIDriverAddon}}
			Expect(makeCreateAddonCMD()).To(RunSuccessfully())
			assertAddonHasPodIDs(api.AWSEFSCSIDriverAddon, 0)
			assertStackNotExists(makePodIDStackName(api.AWSEFSCSIDriverAddon, efsCSIControllerSA))
			assertStackExists(makeIRSAStackName(api.AWSEFSCSIDriverAddon))
		})

		It("should remove IAM permissions when deleting addons", func() {

			By("deleting pod identity associations and IAM role when deleting addon")
			Expect(makeDeleteAddonCMD(api.VPCCNIAddon)).To(RunSuccessfully())
			assertPodIDPresence("kube-system", awsNodeSA, false)
			assertStackDeleted(makeIRSAStackName(api.VPCCNIAddon))
			assertStackDeleted(makePodIDStackName(api.VPCCNIAddon, awsNodeSA))

			By("keeping pod identity associations and IAM role when deleting addon with preserve")
			Expect(makeDeleteAddonCMD(api.AWSEBSCSIDriverAddon, "--preserve")).To(RunSuccessfully())
			assertPodIDPresence("kube-system", ebsCSIControllerSA, true)
			assertStackExists(makePodIDStackName(api.AWSEBSCSIDriverAddon, ebsCSIControllerSA))

			By("cleaning up IAM role on subsequent deletion")
			Expect(makeDeleteAddonCMD(api.AWSEBSCSIDriverAddon)).To(RunSuccessfully())
			assertStackDeleted(makePodIDStackName(api.AWSEBSCSIDriverAddon, ebsCSIControllerSA))
			// now manually cleanup the pod ID to not conflict with subsequent tests
			Expect(params.EksctlDeleteCmd.
				WithArgs(
					"podidentityassociation",
					"--cluster", clusterConfig.Metadata.Name,
					"--namespace", "kube-system",
					"--service-account-name", ebsCSIControllerSA,
				)).
				To(RunSuccessfully())
			assertPodIDPresence("kube-system", ebsCSIControllerSA, false)

			By("deleting IRSA when deleting addons")
			Expect(makeDeleteAddonCMD(api.AWSEFSCSIDriverAddon)).To(RunSuccessfully())
			assertPodIDPresence("kube-system", efsCSIControllerSA, false)
			assertStackDeleted(makeIRSAStackName(api.AWSEFSCSIDriverAddon))
		})

		It("should update IAM permissions when updating or migrating addons", func() {
			clusterConfig.AddonsConfig.AutoApplyPodIdentityAssociations = false
			clusterConfig.Addons = []*api.Addon{
				{Name: api.VPCCNIAddon},
				{Name: api.AWSEBSCSIDriverAddon},
			}

			By("creating addons with IRSA")
			Expect(makeCreateAddonCMD()).To(RunSuccessfully())
			assertAddonHasPodIDs(api.VPCCNIAddon, 0)
			assertAddonHasPodIDs(api.AWSEBSCSIDriverAddon, 0)
			assertStackExists(makeIRSAStackName(api.VPCCNIAddon))
			assertStackExists(makeIRSAStackName(api.AWSEBSCSIDriverAddon))
			assertStackNotExists(makePodIDStackName(api.VPCCNIAddon, awsNodeSA))
			assertStackNotExists(makePodIDStackName(api.AWSEBSCSIDriverAddon, ebsCSIControllerSA))

			By("updating the addon to use pod identity")
			clusterConfig.Addons[1].PodIdentityAssociations = &[]api.PodIdentityAssociation{
				{
					Namespace:            "kube-system",
					ServiceAccountName:   ebsCSIControllerSA,
					PermissionPolicyARNs: []string{"arn:aws:iam::aws:policy/service-role/AmazonEBSCSIDriverPolicy"},
				},
			}
			Expect(makeUpdateAddonCMD()).To(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring("deleting old IRSA stack for addon %s", api.AWSEBSCSIDriverAddon)),
				ContainElement(ContainSubstring("will delete stack %q", makeIRSAStackName(api.AWSEBSCSIDriverAddon))),
			))
			assertAddonHasPodIDs(api.AWSEBSCSIDriverAddon, 1)
			assertStackNotExists(makeIRSAStackName(api.AWSEBSCSIDriverAddon))

			By("ensuring that updating addon again works")
			Expect(makeUpdateAddonCMD()).To(RunSuccessfullyWithOutputStringLines(
				Not(ContainElement(ContainSubstring("deleting old IRSA stack for addon %s", api.AWSEBSCSIDriverAddon))),
			))

			By("failing to update addon with pod identity if the field is unset")
			clusterConfig.Addons[1].PodIdentityAssociations = nil
			session := makeUpdateAddonCMD().Run()
			Expect(session.ExitCode()).To(Equal(1))
			Expect(string(session.Err.Contents())).To(ContainSubstring("addon %s has pod identity associations,"+
				" to remove pod identity associations from an addon, addon.podIdentityAssociations must be explicitly set to []; "+
				"if the addon was migrated to use pod identity, addon.podIdentityAssociations must be set to values obtained from "+
				"`aws eks describe-pod-identity-association --cluster-name=%s", api.AWSEBSCSIDriverAddon, clusterConfig.Metadata.Name))

			By("removing all pod identity associations owned by the addon")
			clusterConfig.Addons[1].PodIdentityAssociations = &[]api.PodIdentityAssociation{}
			Expect(makeUpdateAddonCMD()).To(RunSuccessfully())
			assertAddonHasPodIDs(api.AWSEBSCSIDriverAddon, 0)

			By("migrating an addon to pod identity using the utils command")
			Expect(params.EksctlUtilsCmd.
				WithArgs(
					"migrate-to-pod-identity",
					"--cluster", clusterConfig.Metadata.Name,
					"--approve",
				)).To(RunSuccessfully())
			assertAddonHasPodIDs(api.VPCCNIAddon, 1)
			// assert IRSA stack still exists, but tags reflect association with podID
			stackHasTag := func(stack cfntypes.Stack, tag string) bool {
				return slices.ContainsFunc(stack.Tags, func(t cfntypes.Tag) bool {
					return *t.Key == tag
				})
			}
			describeStackOutput, err := awsProvider.CloudFormation().DescribeStacks(context.Background(), &cfn.DescribeStacksInput{
				StackName: aws.String(makeIRSAStackName(api.VPCCNIAddon)),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(describeStackOutput.Stacks).To(HaveLen(1))
			Expect(stackHasTag(describeStackOutput.Stacks[0], api.IAMServiceAccountNameTag)).To(BeFalse())
			Expect(stackHasTag(describeStackOutput.Stacks[0], api.PodIdentityAssociationNameTag)).To(BeTrue())
		})
	})

	Context("addons in a cluster with no nodes", func() {
		var clusterConfig *api.ClusterConfig

		BeforeEach(func() {
			clusterConfig = api.NewClusterConfig()
			clusterConfig.Metadata.Name = params.NewClusterName("addons-wait")
			clusterConfig.Metadata.Region = params.Region
			clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name: "mng",
						ScalingConfig: &api.ScalingConfig{
							DesiredCapacity: aws.Int(0),
						},
					},
				},
			}
			clusterConfig.NodeGroups = []*api.NodeGroup{
				{
					NodeGroupBase: &api.NodeGroupBase{
						Name: "ng",
						ScalingConfig: &api.ScalingConfig{
							DesiredCapacity: aws.Int(0),
						},
					},
				},
			}
			clusterConfig.Addons = []*api.Addon{
				{
					Name:             "vpc-cni",
					AttachPolicyARNs: []string{"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"},
				},
				{
					Name: "coredns",
				},
				{
					Name: "kube-proxy",
				},
				{
					Name: "aws-ebs-csi-driver",
				},
			}

			By("creating a cluster with no worker nodes")
			cmd := params.EksctlCreateCmd.
				WithArgs(
					"cluster",
					"--config-file", "-",
					"--verbose", "4",
				).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.Reader(clusterConfig))
			Expect(cmd).To(RunSuccessfully())
		})

		It("should show some addons in a degraded state", func() {
			type addonStatus struct {
				Name   string `json:"Name"`
				Status string `json:"Status"`
			}

			Eventually(func() []addonStatus {
				cmd := params.EksctlGetCmd.
					WithArgs(
						"addons",
						"--cluster", clusterConfig.Metadata.Name,
						"-o", "json",
					)
				session := cmd.Run()
				Expect(session.ExitCode()).To(Equal(0))

				var as []addonStatus
				Expect(json.Unmarshal(session.Buffer().Contents(), &as)).To(Succeed())
				return as
			}, "5m", "10s").Should(ConsistOf(
				addonStatus{
					Name:   "aws-ebs-csi-driver",
					Status: "DEGRADED",
				},
				addonStatus{
					Name:   "coredns",
					Status: "DEGRADED",
				},
				addonStatus{
					Name:   "kube-proxy",
					Status: "ACTIVE",
				},
				addonStatus{
					Name:   "vpc-cni",
					Status: "ACTIVE",
				},
			))
		})

		AfterEach(func() {
			cmd := params.EksctlDeleteClusterCmd.
				WithArgs("--name", clusterConfig.Metadata.Name).
				WithArgs("--verbose", "2")
			Expect(cmd).To(RunSuccessfully())
		})
	})
})

var _ = AfterSuite(func() {
	cmd := params.EksctlDeleteCmd.WithArgs(
		"cluster", params.ClusterName,
		"--disable-nodegroup-eviction",
		"--verbose", "2",
	)
	Expect(cmd).To(RunSuccessfully())
})

func getRawClient(ctx context.Context, clusterName string) *kubewrapper.RawClient {
	cfg := &api.ClusterConfig{
		Metadata: &api.ClusterMeta{
			Name:   clusterName,
			Region: params.Region,
		},
	}
	ctl, err := eks.New(context.TODO(), &api.ProviderConfig{Region: params.Region}, cfg)
	Expect(err).NotTo(HaveOccurred())

	err = ctl.RefreshClusterStatus(ctx, cfg)
	Expect(err).ShouldNot(HaveOccurred())
	rawClient, err := ctl.NewRawClient(cfg)
	Expect(err).NotTo(HaveOccurred())
	awsProvider = ctl.AWSProvider
	return rawClient
}

func getInitialClusterConfig() *api.ClusterConfig {
	clusterConfig := api.NewClusterConfig()
	clusterConfig.Metadata.Name = params.ClusterName
	clusterConfig.Metadata.Version = api.LatestVersion
	clusterConfig.Metadata.Region = params.Region
	clusterConfig.IAM.WithOIDC = api.Enabled()
	clusterConfig.Addons = []*api.Addon{
		{
			Name:             "vpc-cni",
			AttachPolicyARNs: []string{"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"},
		},
		{
			Name: "kube-proxy",
		},
	}
	clusterConfig.AddonsConfig.DisableDefaultAddons = true

	ng := &api.ManagedNodeGroup{
		NodeGroupBase: &api.NodeGroupBase{
			Name: "ng",
		},
	}
	clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{ng}

	return clusterConfig
}

func getConfigMap(clientset kubernetes.Interface, addonName string) *corev1.ConfigMap {
	configMap, err := clientset.CoreV1().ConfigMaps("kube-system").Get(context.Background(), addonName, metav1.GetOptions{})
	Expect(err).NotTo(HaveOccurred())
	return configMap
}

func updateConfigMap(clientset kubernetes.Interface, configMap *corev1.ConfigMap) *corev1.ConfigMap {
	configMap, err := clientset.CoreV1().ConfigMaps("kube-system").Update(context.Background(), configMap, metav1.UpdateOptions{})
	Expect(err).NotTo(HaveOccurred())
	return configMap
}

func getCacheValue(configMap *corev1.ConfigMap) string {
	coreFile, ok := configMap.Data["Corefile"]
	Expect(ok).To(BeTrue())

	coreFileValues := strings.Fields(strings.Replace(coreFile, "\n", " ", -1))
	return coreFileValues[k8sslices.Index(coreFileValues, "cache")+1]
}

func updateCacheValue(configMap *corev1.ConfigMap, currentValue string, newValue string) {
	coreFile, ok := configMap.Data["Corefile"]
	Expect(ok).To(BeTrue())

	configMap.Data["Corefile"] = strings.Replace(coreFile, "cache "+currentValue, "cache "+newValue, -1)
}

func addToString(s string, n int) string {
	i, err := strconv.Atoi(s)
	Expect(err).NotTo(HaveOccurred())
	return strconv.Itoa(i + n)
}
