//go:build integration
// +build integration

package kms

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("kms")
}

func TestEKSKMS(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var (
	kmsKeyARN *string
	ctl       api.ClusterProvider
)

var _ = BeforeSuite(func() {
	clusterConfig := api.NewClusterConfig()
	clusterConfig.Metadata.Name = params.ClusterName
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

	clusterProvider, err := eks.New(context.TODO(), &api.ProviderConfig{Region: params.Region}, clusterConfig)
	Expect(err).NotTo(HaveOccurred())
	ctl = clusterProvider.AWSProvider

	cfg := NewConfig(params.Region)
	kmsClient := kms.NewFromConfig(cfg)
	output, err := kmsClient.CreateKey(context.Background(), &kms.CreateKeyInput{
		Description: aws.String(fmt.Sprintf("Key to test KMS encryption on EKS cluster %s", params.ClusterName)),
	})
	Expect(err).NotTo(HaveOccurred())
	kmsKeyARN = output.KeyMetadata.Arn
})

var _ = Describe("(Integration) [EKS kms test]", func() {
	Context("Creating a cluster and enabling kms", func() {

		It("supports enabling KMS encryption", func() {
			enableEncryptionCMD := func() Cmd {
				return params.EksctlUtilsCmd.
					WithTimeout(2*time.Hour).
					WithArgs(
						"enable-secrets-encryption",
						"--cluster", params.ClusterName,
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

var _ = AfterSuite(func() {
	cmd := params.EksctlDeleteCmd.WithArgs(
		"cluster", params.ClusterName,
		"--verbose", "2",
	)
	Expect(cmd).To(RunSuccessfully())

	cfg := NewConfig(params.Region)
	kmsClient := kms.NewFromConfig(cfg)
	_, err := kmsClient.ScheduleKeyDeletion(context.Background(), &kms.ScheduleKeyDeletionInput{
		KeyId:               kmsKeyARN,
		PendingWindowInDays: aws.Int32(7),
	})
	Expect(err).NotTo(HaveOccurred())
})
