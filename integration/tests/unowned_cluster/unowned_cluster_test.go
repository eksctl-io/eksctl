// +build integration

package unowned_clusters

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/weaveworks/eksctl/integration/utilities/unowned"

	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/weaveworks/eksctl/pkg/eks"

	"github.com/aws/aws-sdk-go/aws"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/pkg/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	// "unowned_clusters" lead to names longer than allowed for CF stacks
	params = tests.NewParams("uc")
}

func TestE2E(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) [non-eksctl cluster & nodegroup support]", func() {
	var (
		version, ng1, mng1, mng2 string
		ctl                      api.ClusterProvider
		configFile               *os.File
		cfg                      *api.ClusterConfig
		unownedCluster           *unowned.Cluster
	)

	BeforeSuite(func() {
		ng1 = "ng-1"
		mng1 = "mng-1"
		mng2 = "mng-2"
		version = "1.18"
		cfg = api.NewClusterConfig()
		cfg.Metadata = &api.ClusterMeta{
			Version: version,
			Name:    params.ClusterName,
			Region:  params.Region,
		}
		var err error
		configFile, err = ioutil.TempFile("", "")
		Expect(err).NotTo(HaveOccurred())

		if !params.SkipCreate {
			unownedCluster = unowned.NewCluster(cfg)
			cfg.VPC = unownedCluster.VPC
			unownedCluster.CreateNodegroups(mng1)

			clusterProvider, err := eks.New(&api.ProviderConfig{Region: params.Region}, cfg)
			Expect(err).NotTo(HaveOccurred())
			ctl = clusterProvider.Provider
		}
	})

	AfterSuite(func() {
		if !params.SkipCreate && !params.SkipDelete {
			unownedCluster.DeleteStack()
		}
		Expect(os.RemoveAll(configFile.Name())).To(Succeed())
	})

	It("supports creating nodegroups", func() {
		cfg.NodeGroups = []*api.NodeGroup{{
			NodeGroupBase: &api.NodeGroupBase{
				Name: ng1,
			}},
		}
		// write config file so that the nodegroup creates have access to the vpc spec
		cmd := params.EksctlCreateCmd.
			WithArgs(
				"nodegroup",
				"--config-file", "-",
				"--verbose", "4",
			).
			WithoutArg("--region", params.Region).
			WithStdinJSONContent(cfg)
		Expect(cmd).To(RunSuccessfully())
	})

	It("supports creating managed nodegroups", func() {
		cfg.ManagedNodeGroups = []*api.ManagedNodeGroup{{
			NodeGroupBase: &api.NodeGroupBase{
				Name: mng2,
			}},
		}
		// write config file so that the nodegroup creates have access to the vpc spec
		cmd := params.EksctlCreateCmd.
			WithArgs(
				"nodegroup",
				"--config-file", "-",
				"--verbose", "4",
			).
			WithoutArg("--region", params.Region).
			WithStdinJSONContent(cfg)
		Expect(cmd).To(RunSuccessfully())
	})

	It("supports getting non-eksctl resources", func() {
		By("getting clusters")
		cmd := params.EksctlGetCmd.
			WithArgs(
				"clusters",
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(ContainSubstring(params.ClusterName)),
		))

		By("getting nodegroups")
		cmd = params.EksctlGetCmd.
			WithArgs(
				"nodegroups",
				"--cluster", params.ClusterName,
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(ContainSubstring(ng1)),
		))
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(ContainSubstring(mng1)),
		))
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(ContainSubstring(mng2)),
		))
	})

	It("supports labels", func() {
		By("setting labels on a managed nodegroup")
		cmd := params.EksctlSetLabelsCmd.
			WithArgs(
				"--cluster", params.ClusterName,
				"--nodegroup", mng1,
				"--labels", "key=value",
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())

		By("getting labels for a managed nodegroup")
		cmd = params.EksctlGetCmd.
			WithArgs(
				"labels",
				"--cluster", params.ClusterName,
				"--nodegroup", mng1,
				"--verbose", "2",
			)
		// It sometimes takes forever for the above set to take effect
		Eventually(func() *gbytes.Buffer { return cmd.Run().Out }, time.Minute*2).Should(gbytes.Say("key=value"))

		By("unsetting labels on a managed nodegroup")
		cmd = params.EksctlUnsetLabelsCmd.
			WithArgs(
				"--cluster", params.ClusterName,
				"--nodegroup", mng1,
				"--labels", "key",
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())
	})

	It("supports IRSA", func() {
		By("enabling OIDC")
		cmd := params.EksctlUtilsCmd.
			WithArgs(
				"associate-iam-oidc-provider",
				"--name", params.ClusterName,
				"--approve",
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())

		By("creating an IAMServiceAccount")
		cmd = params.EksctlCreateCmd.
			WithArgs(
				"iamserviceaccount",
				"--cluster", params.ClusterName,
				"--name", "test-sa",
				"--namespace", "default",
				"--attach-policy-arn",
				"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
				"--approve",
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())

		By("getting IAMServiceAccounts")
		cmd = params.EksctlGetCmd.
			WithArgs(
				"iamserviceaccounts",
				"--cluster", params.ClusterName,
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(ContainSubstring("test-sa")),
		))
	})

	It("supports cluster upgrades", func() {
		By("upgrading the cluster")
		cmd := params.EksctlUpgradeCmd.
			WithArgs(
				"cluster",
				"--name", params.ClusterName,
				"--version", "1.19",
				"--approve",
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())
	})

	It("supports addons", func() {
		By("creating an addon")
		cmd := params.EksctlCreateCmd.
			WithArgs(
				"addon",
				"--cluster", params.ClusterName,
				"--name", "vpc-cni",
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())

		By("getting an addon")
		cmd = params.EksctlGetCmd.
			WithArgs(
				"addons",
				"--cluster", params.ClusterName,
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(ContainSubstring("vpc-cni")),
		))
	})

	It("supports fargate", func() {
		By("creating a fargate profile")
		cmd := params.EksctlCreateCmd.
			WithArgs(
				"fargateprofile",
				"--cluster", params.ClusterName,
				"--name", "fp-test",
				"--namespace", "default",
			)
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(SatisfyAll(ContainSubstring("created"), ContainSubstring("fp-test"))),
		))

		By("getting a fargate profile")
		cmd = params.EksctlGetCmd.
			WithArgs(
				"fargateprofile",
				"--cluster", params.ClusterName,
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(ContainSubstring("fp-test")),
		))

		By("deleting a fargate profile")
		cmd = params.EksctlDeleteCmd.
			WithArgs(
				"fargateprofile",
				"--cluster", params.ClusterName,
				"--name", "fp-test",
				"--wait",
			)
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(SatisfyAll(ContainSubstring("deleted"), ContainSubstring("fp-test"))),
		))
	})

	It("supports managed nodegroup upgrades", func() {
		cmd := params.EksctlUpgradeCmd.
			WithArgs(
				"nodegroup",
				"--name", mng1,
				"--cluster", params.ClusterName,
				"--kubernetes-version", "1.19",
				"--wait",
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())
	})

	It("supports draining and scaling nodegroups", func() {
		By("scaling a nodegroup")
		cmd := params.EksctlScaleNodeGroupCmd.
			WithArgs(
				"--name", mng1,
				"--nodes", "2",
				"--nodes-max", "3",
				"--cluster", params.ClusterName,
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())

		By("draining a nodegroup")
		cmd = params.EksctlDrainNodeGroupCmd.
			WithArgs(
				"--cluster", params.ClusterName,
				"--name", mng1,
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())
	})

	It("supports deleting nodegroups", func() {
		cmd := params.EksctlDeleteCmd.
			WithArgs(
				"nodegroup",
				"--name", mng1,
				"--cluster", params.ClusterName,
				"--verbose", "2",
			)
		Expect(cmd).To(RunSuccessfully())
	})

	Context("KMS", func() {
		var kmsKeyARN *string

		BeforeEach(func() {
			kmsClient := kms.New(ctl.ConfigProvider())
			output, err := kmsClient.CreateKey(&kms.CreateKeyInput{
				Description: aws.String("Key to test KMS encryption on EKS clusters"),
			})
			Expect(err).NotTo(HaveOccurred())
			kmsKeyARN = output.KeyMetadata.Arn
		})

		It("supports enabling KMS encryption", func() {
			enableEncryptionCMD := func() Cmd {
				return params.EksctlUtilsCmd.
					WithTimeout(1*time.Hour).
					WithArgs(
						"enable-secrets-encryption",
						"--cluster", params.ClusterName,
						"--key-arn", *kmsKeyARN,
					)
			}

			By("enabling KMS encryption on the cluster using the new key")
			cmd := enableEncryptionCMD()
			Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring(" initiated KMS encryption")),
				ContainElement(ContainSubstring("KMS encryption applied to all Secret resources")),
			))
		})

		AfterEach(func() {
			kmsClient := kms.New(ctl.ConfigProvider())
			_, err := kmsClient.ScheduleKeyDeletion(&kms.ScheduleKeyDeletionInput{
				KeyId:               kmsKeyARN,
				PendingWindowInDays: aws.Int64(7),
			})
			Expect(err).NotTo(HaveOccurred())
		})
	})

	It("supports deleting clusters", func() {
		if params.SkipDelete {
			Skip("params.SkipDelete is true")
		}
		By("deleting the cluster")
		cmd := params.EksctlDeleteCmd.
			WithArgs(
				"cluster",
				"--name", params.ClusterName,
				"--verbose", "3",
			)
		Expect(cmd).To(RunSuccessfully())
	})
})
