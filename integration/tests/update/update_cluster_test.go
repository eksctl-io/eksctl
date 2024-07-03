//go:build integration
// +build integration

package update

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/go-version"

	"github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go/aws"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	"github.com/weaveworks/eksctl/pkg/addons"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils"
	"github.com/weaveworks/eksctl/pkg/utils/file"
)

const (
	k8sUpdatePollInterval = "30s"
	k8sUpdatePollTimeout  = "10m"
)

var (
	defaultCluster  string
	params          *tests.Params
	clusterProvider *eks.ClusterProvider
)

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("up")
	defaultCluster = params.ClusterName
}

func TestUpdate(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var (
	eksVersion     string
	nextEKSVersion string
)

const (
	initNG = "kp-ng-0"
	botNG  = "bot-ng-0"
)

var _ = BeforeSuite(func() {
	params.KubeconfigTemp = false
	if params.KubeconfigPath == "" {
		wd, _ := os.Getwd()
		f, _ := os.CreateTemp(wd, "kubeconfig-")
		params.KubeconfigPath = f.Name()
		params.KubeconfigTemp = true
	}

	if params.SkipCreate {
		fmt.Fprintf(GinkgoWriter, "will use existing cluster %s", defaultCluster)
		if !file.Exists(params.KubeconfigPath) {
			// Generate the Kubernetes configuration that eksctl create
			// would have generated otherwise:
			cmd := params.EksctlUtilsCmd.WithArgs(
				"write-kubeconfig",
				"--verbose", "4",
				"--cluster", defaultCluster,
				"--kubeconfig", params.KubeconfigPath,
			)
			Expect(cmd).To(RunSuccessfully())
		}
		return
	}

	fmt.Fprintf(GinkgoWriter, "Using kubeconfig: %s\n", params.KubeconfigPath)

	eksVersion, nextEKSVersion = clusterutils.GetCurrentAndNextVersionsForUpgrade(params.Version)

	clusterConfig := api.NewClusterConfig()
	clusterConfig.Metadata.Name = defaultCluster
	clusterConfig.Metadata.Region = params.Region
	clusterConfig.Metadata.Version = eksVersion
	clusterConfig.Metadata.Tags = map[string]string{
		"alpha.eksctl.io/description": "eksctl integration test",
	}
	clusterConfig.ManagedNodeGroups = []*api.ManagedNodeGroup{
		{
			NodeGroupBase: &api.NodeGroupBase{
				Name:         initNG,
				InstanceType: "t3.large",
				ScalingConfig: &api.ScalingConfig{
					DesiredCapacity: aws.Int(1),
				},
				Labels: map[string]string{
					"ng-name": initNG,
				},
			},
		},
		{
			NodeGroupBase: &api.NodeGroupBase{
				Name:         botNG,
				AMIFamily:    api.NodeImageFamilyBottlerocket,
				InstanceType: "t3.small",
				ScalingConfig: &api.ScalingConfig{
					DesiredCapacity: aws.Int(1),
				},
				Labels: map[string]string{
					"ng-name": botNG,
				},
			},
		},
	}

	cmd := params.EksctlCreateCmd.WithArgs(
		"cluster",
		"--config-file", "-",
		"--kubeconfig", params.KubeconfigPath,
		"--verbose", "4",
	).
		WithoutArg("--region", params.Region).
		WithStdin(clusterutils.Reader(clusterConfig))
	Expect(cmd).To(RunSuccessfully())

	var err error
	clusterProvider, err = newClusterProvider(context.Background())
	Expect(err).NotTo(HaveOccurred())
})

var _ = Describe("(Integration) Upgrading cluster", func() {

	Context("control plane", func() {
		It("should have created an EKS cluster and two CloudFormation stacks", func() {
			config := NewConfig(params.Region)

			Expect(config).To(HaveExistingCluster(params.ClusterName, string(types.ClusterStatusActive), eksVersion))
			Expect(config).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-cluster", params.ClusterName)))
			Expect(config).To(HaveExistingStack(fmt.Sprintf("eksctl-%s-nodegroup-%s", params.ClusterName, initNG)))
		})

		It("should upgrade the control plane to the next version", func() {
			cmd := params.EksctlUpgradeCmd.
				WithArgs(
					"cluster",
					"--verbose", "4",
					"--name", params.ClusterName,
					"--approve",
				)
			Expect(cmd).To(RunSuccessfully())

			By(fmt.Sprintf("checking that control plane is updated to %v", nextEKSVersion))
			config, err := clientcmd.BuildConfigFromFlags("", params.KubeconfigPath)
			Expect(err).NotTo(HaveOccurred())

			clientSet, err := kubernetes.NewForConfig(config)
			Expect(err).NotTo(HaveOccurred())

			Eventually(func() string {
				serverVersion, err := clientSet.ServerVersion()
				Expect(err).NotTo(HaveOccurred())
				return fmt.Sprintf("%s.%s", serverVersion.Major, strings.TrimSuffix(serverVersion.Minor, "+"))
			}, k8sUpdatePollTimeout, k8sUpdatePollInterval).Should(Equal(nextEKSVersion))
		})
	})

	Context("default networking addons", func() {
		defaultNetworkingAddons := []string{"vpc-cni", "kube-proxy", "coredns"}

		It("should suggest using `eksctl update addon` for updating default addons", func() {
			assertAddonError := func(updateAddonName, addonName string) {
				cmd := params.EksctlUtilsCmd.WithArgs(
					fmt.Sprintf("update-%s", updateAddonName),
					"--cluster", params.ClusterName,
					"--verbose", "4",
					"--approve",
				)
				session := cmd.Run()
				ExpectWithOffset(1, session.ExitCode()).NotTo(BeZero())
				ExpectWithOffset(1, string(session.Err.Contents())).To(ContainSubstring("Error: addon %s is installed as a managed EKS addon; "+
					"to update it, use `eksctl update addon` instead", addonName))
			}
			assertAddonError("aws-node", "vpc-cni")
			for _, addonName := range defaultNetworkingAddons {
				updateAddonName := addonName
				if addonName == "vpc-cni" {
					updateAddonName = "aws-node"
				}
				assertAddonError(updateAddonName, addonName)
			}
		})
	})

	Context("addons", func() {
		It("should upgrade kube-proxy", func() {
			cmd := params.EksctlUpdateCmd.
				WithArgs(
					"addon",
					"--name", "kube-proxy",
					"--cluster", params.ClusterName,
					"--version", "latest",
					"--wait",
					"--verbose", "4",
				)
			Expect(cmd).To(RunSuccessfully())

			rawClient := getRawClient(context.Background(), clusterProvider)
			Eventually(func() string {
				daemonSet, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), "kube-proxy", metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				kubeProxyVersion, err := addons.ImageTag(daemonSet.Spec.Template.Spec.Containers[0].Image)
				Expect(err).NotTo(HaveOccurred())
				v, err := version.NewVersion(kubeProxyVersion)
				Expect(err).NotTo(HaveOccurred())
				segments := v.Segments()
				Expect(len(segments)).To(BeNumerically(">=", 2))
				return fmt.Sprintf("%d.%d", segments[0], segments[1])
			}, k8sUpdatePollTimeout, k8sUpdatePollInterval).Should(Equal(nextEKSVersion))
		})

		It("should upgrade aws-node", func() {
			rawClient := getRawClient(context.Background(), clusterProvider)
			getAWSNodeVersion := func() string {
				awsNode, err := rawClient.ClientSet().AppsV1().DaemonSets(metav1.NamespaceSystem).Get(context.TODO(), "aws-node", metav1.GetOptions{})
				Expect(err).NotTo(HaveOccurred())
				imageTag, err := addons.ImageTag(awsNode.Spec.Template.Spec.Containers[0].Image)
				Expect(err).NotTo(HaveOccurred())
				return imageTag
			}
			preUpdateAWSNodeVersion := getAWSNodeVersion()

			cmd := params.EksctlUpdateCmd.
				WithArgs(
					"addon",
					"--name", "vpc-cni",
					"--cluster", params.ClusterName,
					"--version", "latest",
					"--wait",
					"--verbose", "4",
				)
			Expect(cmd).To(RunSuccessfully())
			Eventually(getAWSNodeVersion, k8sUpdatePollTimeout, k8sUpdatePollInterval).ShouldNot(Equal(preUpdateAWSNodeVersion))
		})

		It("should upgrade coredns", func() {
			cmd := params.EksctlUpdateCmd.
				WithArgs(
					"addon",
					"--name", "coredns",
					"--cluster", params.ClusterName,
					"--version", "latest",
					"--wait",
					"--verbose", "4",
				)
			Expect(cmd).To(RunSuccessfully())
		})

	})

	Context("nodegroup", func() {
		It("should upgrade the initial nodegroup to the next version", func() {
			cmd := params.EksctlUpgradeCmd.WithArgs(
				"nodegroup",
				"--verbose", "4",
				"--cluster", params.ClusterName,
				"--name", initNG,
				"--kubernetes-version", nextEKSVersion,
				"--timeout=60m", // wait for CF stacks to finish update
			)
			ExpectWithOffset(1, cmd).To(RunSuccessfullyWithOutputString(ContainSubstring("nodegroup successfully upgraded")))

			cmd = params.EksctlGetCmd.WithArgs(
				"nodegroup",
				"--cluster", params.ClusterName,
				"--name", initNG,
				"--output", "yaml",
			)
			ExpectWithOffset(1, cmd).To(RunSuccessfullyWithOutputString(ContainSubstring(fmt.Sprintf("Version: \"%s\"", nextEKSVersion))))
		})

		It("should upgrade the Bottlerocket nodegroup to the next version", func() {
			cmd := params.EksctlUpgradeCmd.WithArgs(
				"nodegroup",
				"--verbose", "4",
				"--cluster", params.ClusterName,
				"--name", botNG,
				"--kubernetes-version", nextEKSVersion,
				"--timeout=60m", // wait for CF stacks to finish update
			)
			ExpectWithOffset(1, cmd).To(RunSuccessfullyWithOutputString(ContainSubstring("nodegroup successfully upgraded")))

			cmd = params.EksctlGetCmd.WithArgs(
				"nodegroup",
				"--cluster", params.ClusterName,
				"--name", botNG,
				"--output", "yaml",
			)
			ExpectWithOffset(1, cmd).To(RunSuccessfullyWithOutputString(ContainSubstring(fmt.Sprintf("Version: \"%s\"", nextEKSVersion))))
		})
	})
})

var _ = AfterSuite(func() {
	params.DeleteClusters()
	gexec.KillAndWait()
	if params.KubeconfigTemp {
		os.Remove(params.KubeconfigPath)
	}
	os.RemoveAll(params.TestDirectory)
})

func newClusterProvider(ctx context.Context) (*eks.ClusterProvider, error) {
	cfg := &api.ClusterConfig{
		Metadata: &api.ClusterMeta{
			Name:   params.ClusterName,
			Region: params.Region,
		},
	}
	ctl, err := eks.New(ctx, &api.ProviderConfig{Region: params.Region}, cfg)
	if err != nil {
		return nil, err
	}
	if err := ctl.RefreshClusterStatus(ctx, cfg); err != nil {
		return nil, err
	}
	return ctl, nil
}

func defaultClusterConfig() *api.ClusterConfig {
	return &api.ClusterConfig{
		Metadata: &api.ClusterMeta{
			Name:   params.ClusterName,
			Region: params.Region,
		},
	}
}

func getRawClient(ctx context.Context, ctl *eks.ClusterProvider) *kubewrapper.RawClient {
	clusterConfig := defaultClusterConfig()
	Expect(ctl.RefreshClusterStatus(ctx, clusterConfig)).To(Succeed())
	rawClient, err := ctl.NewRawClient(clusterConfig)
	Expect(err).NotTo(HaveOccurred())
	return rawClient
}
