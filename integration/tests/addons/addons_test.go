//go:build integration
// +build integration

package addons

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/utils/strings/slices"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/weaveworks/eksctl/integration/runner"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	params    *tests.Params
	rawClient *kubewrapper.RawClient
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

			By("successfully creating the aws-ebs-csi-driver addon via config file")
			// setup config file
			clusterConfig := getInitialClusterConfig()
			clusterConfig.Addons = append(clusterConfig.Addons, &api.Addon{
				Name: api.AWSEBSCSIDriverAddon,
			})
			data, err := json.Marshal(clusterConfig)

			Expect(err).NotTo(HaveOccurred())
			cmd := params.EksctlCreateCmd.
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
						"--name", api.AWSEBSCSIDriverAddon,
						"--cluster", clusterName,
						"--verbose", "2",
					)
				return cmd
			}, "5m", "30s").Should(RunSuccessfullyWithOutputStringLines(ContainElement(ContainSubstring("ACTIVE"))))

			By("successfully creating the kube-proxy addon")
			cmd = params.EksctlCreateCmd.
				WithArgs(
					"addon",
					"--name", "kube-proxy",
					"--cluster", clusterName,
					"--force",
					"--wait",
					"--verbose", "2",
				)
			Expect(cmd).To(RunSuccessfully())

			Eventually(func() runner.Cmd {
				cmd := params.EksctlGetCmd.
					WithArgs(
						"addon",
						"--name", "kube-proxy",
						"--cluster", clusterName,
						"--verbose", "2",
					)
				return cmd
			}, "5m", "30s").Should(RunSuccessfullyWithOutputStringLines(ContainElement(ContainSubstring("ACTIVE"))))

			By("Deleting the kube-proxy addon")
			cmd = params.EksctlDeleteCmd.
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

			_, err = rawClient.ClientSet().AppsV1().DaemonSets("kube-system").Get(context.Background(), "aws-node", metav1.GetOptions{})
			Expect(err).NotTo(HaveOccurred())
		})

		It("should have full control over configMap when creating addons", func() {
			var (
				clusterConfig *api.ClusterConfig
				configMap     *corev1.ConfigMap
			)

			configMap = getConfigMap(rawClient.ClientSet(), "coredns")
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

			cmd := params.EksctlCreateCmd.
				WithArgs(
					"addon",
					"--config-file", "-",
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))
			Expect(cmd).ShouldNot(RunSuccessfully())

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
})

var _ = AfterSuite(func() {
	cmd := params.EksctlDeleteCmd.WithArgs(
		"cluster", params.ClusterName,
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
	}

	ng := &api.ManagedNodeGroup{
		NodeGroupBase: &api.NodeGroupBase{
			Name: "ng",
		},
	}
	clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{ng}

	return clusterConfig
}

func getConfigMap(clientset kubernetes.Interface, name string) *corev1.ConfigMap {
	configMap, err := clientset.CoreV1().ConfigMaps("kube-system").Get(context.Background(), "coredns", metav1.GetOptions{})
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
	return coreFileValues[slices.Index(coreFileValues, "cache")+1]
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
