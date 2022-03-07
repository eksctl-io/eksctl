//go:build integration
// +build integration

package kms

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/kms"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("kms")
}

func TestEKSkms(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) [EKS kms test]", func() {
	Context("Creating a cluster and enabling kms", func() {
		var (
			kmsKeyARN   *string
			clusterName string
			ctl         api.ClusterProvider
		)

		BeforeSuite(func() {
			clusterName = params.NewClusterName("kms")
			clusterConfig := api.NewClusterConfig()
			clusterConfig.Metadata.Name = clusterName
			clusterConfig.Metadata.Version = "latest"
			clusterConfig.Metadata.Region = params.Region

			data, err := json.Marshal(clusterConfig)
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

			clusterProvider, err := eks.New(&api.ProviderConfig{Region: params.Region}, clusterConfig)
			Expect(err).NotTo(HaveOccurred())
			ctl = clusterProvider.Provider

			kmsClient := kms.New(ctl.ConfigProvider())
			output, err := kmsClient.CreateKey(&kms.CreateKeyInput{
				Description: aws.String(fmt.Sprintf("Key to test KMS encryption on EKS cluster %s", clusterName)),
			})
			Expect(err).NotTo(HaveOccurred())
			kmsKeyARN = output.KeyMetadata.Arn
		})

		AfterSuite(func() {
			cmd := params.EksctlDeleteCmd.WithArgs(
				"cluster", clusterName,
				"--verbose", "2",
			)
			Expect(cmd).To(RunSuccessfully())

			kmsClient := kms.New(ctl.ConfigProvider())
			_, err := kmsClient.ScheduleKeyDeletion(&kms.ScheduleKeyDeletionInput{
				KeyId:               kmsKeyARN,
				PendingWindowInDays: aws.Int64(7),
			})
			Expect(err).NotTo(HaveOccurred())
		})

		It("supports enabling KMS encryption", func() {
			enableEncryptionCMD := func() Cmd {
				return params.EksctlUtilsCmd.
					WithTimeout(2*time.Hour).
					WithArgs(
						"enable-secrets-encryption",
						"--cluster", clusterName,
						"--key-arn", *kmsKeyARN,
					)
			}

			By(fmt.Sprintf("enabling KMS encryption on the cluster using key %q", *kmsKeyARN))
			cmd := enableEncryptionCMD()
			Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring("initiated KMS encryption")),
				ContainElement(ContainSubstring("KMS encryption applied to all Secret resources")),
			))

			By("ensuring `enable-secrets-encryption` works when KMS encryption is already enabled on the cluster")
			cmd = enableEncryptionCMD()
			Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
				ContainElement(ContainSubstring("KMS encryption is already enabled on the cluster")),
				ContainElement(ContainSubstring("KMS encryption applied to all Secret resources")),
			))
		})
	})
})
