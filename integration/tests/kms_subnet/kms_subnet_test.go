//go:build integration
// +build integration

//revive:disable Not changing package name
package kms_subnet

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	corev1 "k8s.io/api/core/v1"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("kms-subnet")
	if err := api.Register(); err != nil {
		panic("unexpected error registering API scheme")
	}
}

func TestKMSSubnet(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var (
	kmsKeyARN *string
	awsConfig aws.Config
)

var _ = BeforeSuite(func() {
	awsConfig = NewConfig(params.Region)

	By("enabling resource-based hostname for subnets")
	cmd := params.EksctlCreateCmd.
		WithArgs(
			"cluster",
			"--config-file", "-",
			"--verbose", "4",
		).
		WithoutArg("--region", params.Region).
		WithStdin(clusterutils.ReaderFromFile(params.ClusterName, params.Region, "testdata/kms-subnet-cluster.yaml"))
	Expect(cmd).To(RunSuccessfully())

	kmsClient := kms.NewFromConfig(awsConfig)
	output, err := kmsClient.CreateKey(context.Background(), &kms.CreateKeyInput{
		Description: aws.String(fmt.Sprintf("Key to test KMS encryption on EKS cluster %s", params.ClusterName)),
	})
	Expect(err).NotTo(HaveOccurred())
	kmsKeyARN = output.KeyMetadata.Arn
})

var _ = Describe("(Integration) [EKS KMS and subnet test]", func() {
	Context("creating a cluster and enabling KMS", func() {
		params.LogStacksEventsOnFailure()

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

	Context("nodes launched in the VPC's subnets", func() {
		It("should have a resource-based hostname", func() {
			clientSet := makeClientSet()
			nodeList := tests.ListNodes(clientSet, "mng-1")
			Expect(nodeList.Items).To(HaveLen(2))
			for _, node := range nodeList.Items {
				instanceID, err := getInstanceID(node)
				Expect(err).NotTo(HaveOccurred())
				Expect(node.Name).To(Equal(fmt.Sprintf("%s.%s.compute.internal", instanceID, params.Region)))
			}
		})
	})
})

var _ = AfterSuite(func() {
	cmd := params.EksctlDeleteCmd.WithArgs(
		"cluster", params.ClusterName,
		"--verbose", "2",
	)
	Expect(cmd).To(RunSuccessfully())

	kmsClient := kms.NewFromConfig(awsConfig)
	_, err := kmsClient.ScheduleKeyDeletion(context.Background(), &kms.ScheduleKeyDeletionInput{
		KeyId:               kmsKeyARN,
		PendingWindowInDays: aws.Int32(7),
	})
	Expect(err).NotTo(HaveOccurred())
})

func getInstanceID(node corev1.Node) (string, error) {
	providerID := node.Spec.ProviderID
	idx := strings.LastIndex(providerID, "/")
	if idx == -1 || idx == len(providerID)-1 {
		return "", fmt.Errorf("unexpected format for node.spec.providerID %q", providerID)
	}
	return providerID[idx+1:], nil
}

func makeClientSet() kubernetes.Interface {
	clusterConfig := &api.ClusterConfig{
		Metadata: &api.ClusterMeta{
			Name:   params.ClusterName,
			Region: params.Region,
		},
	}
	ctx := context.Background()
	clusterProvider, err := eks.New(ctx, &api.ProviderConfig{Region: params.Region}, clusterConfig)
	Expect(err).NotTo(HaveOccurred())
	Expect(clusterProvider.RefreshClusterStatus(ctx, clusterConfig)).To(Succeed())
	clientSet, err := clusterProvider.NewStdClientSet(clusterConfig)
	Expect(err).NotTo(HaveOccurred())
	return clientSet
}
