//go:build integration

package karpenter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	cfn "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	awsec2 "github.com/aws/aws-sdk-go-v2/service/ec2"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	clusterutils "github.com/weaveworks/eksctl/integration/utilities/cluster"
	"github.com/weaveworks/eksctl/integration/utilities/kube"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/karpenter"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params

func init() {
	testing.Init()
	if err := api.Register(); err != nil {
		panic(fmt.Errorf("unexpected error registering API scheme: %w", err))
	}
	params = tests.NewParams("karpenter")
}

func TestKarpenter(t *testing.T) {
	testutils.RegisterAndRun(t)
}

// NewConfig creates an AWS config for the given region
func NewConfig(region string) aws.Config {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	Expect(err).NotTo(HaveOccurred())
	return cfg
}

var _ = Describe("(Integration) Karpenter", func() {
	var (
		clusterName string
	)
	BeforeEach(func() {
		// the randomly generated name we get usually makes one of the resources have a longer than 64 characters name
		// so create our own name here to avoid this error
		clusterName = fmt.Sprintf("it-karpenter-%d", time.Now().Unix())
	})
	AfterEach(func() {
		cmd := params.EksctlDeleteCmd.WithArgs(
			"cluster", clusterName,
			"--verbose", "4",
		)
		Expect(cmd).To(RunSuccessfully())
	})

	Context("Creating a cluster with Karpenter and security group tagging", func() {
		params.LogStacksEventsOnFailure()

		It("should deploy Karpenter successfully and tag security group with karpenter.sh/discovery", func() {
			By("creating a cluster with both Karpenter enabled and karpenter.sh/discovery in metadata.tags")

			clusterConfig := &api.ClusterConfig{
				TypeMeta: metav1.TypeMeta{
					Kind:       api.ClusterConfigKind,
					APIVersion: api.SchemeGroupVersion.String(),
				},
				Metadata: &api.ClusterMeta{
					Name:    clusterName,
					Region:  params.Region,
					Version: api.DefaultVersion,
					Tags: map[string]string{
						"karpenter.sh/discovery": clusterName,
					},
				},
				Karpenter: &api.Karpenter{
					Version: "1.6.2",
				},
				IAM: &api.ClusterIAM{
					WithOIDC: api.Enabled(),
				},
				ManagedNodeGroups: []*api.ManagedNodeGroup{
					{
						NodeGroupBase: &api.NodeGroupBase{
							Name: "managed-ng-1",
							ScalingConfig: &api.ScalingConfig{
								MinSize:         aws.Int(1),
								MaxSize:         aws.Int(2),
								DesiredCapacity: aws.Int(1),
							},
						},
					},
				},
			}

			data, err := json.Marshal(clusterConfig)
			Expect(err).NotTo(HaveOccurred())

			cmd := params.EksctlCreateCmd.
				WithArgs(
					"cluster",
					"--config-file=-",
					"--verbose=4",
					"--kubeconfig", params.KubeconfigPath,
				).
				WithoutArg("--region", params.Region).
				WithStdin(bytes.NewReader(data))
			Expect(cmd).To(RunSuccessfully())

			By("verifying Karpenter pods are healthy")
			kubeTest, err := kube.NewTest(params.KubeconfigPath)
			Expect(err).NotTo(HaveOccurred())
			
			// Check that Karpenter webhook pod is ready
			Expect(kubeTest.WaitForPodsReady(karpenter.DefaultNamespace, metav1.ListOptions{
				LabelSelector: "app.kubernetes.io/instance=karpenter",
			}, 1, 10*time.Minute)).To(Succeed())

			By("verifying the cluster shared node security group has karpenter.sh/discovery tags")
			config := NewConfig(params.Region)
			cfnSession := cfn.NewFromConfig(config)
			ec2Session := awsec2.NewFromConfig(config)

			// Get the cluster stack to find the node security group
			describeStackOut, err := cfnSession.DescribeStacks(context.Background(), &cfn.DescribeStacksInput{
				StackName: aws.String(fmt.Sprintf("eksctl-%s-cluster", clusterName)),
			})
			Expect(err).NotTo(HaveOccurred())

			var nodeSecurityGroupID string
			for _, output := range describeStackOut.Stacks[0].Outputs {
				if *output.OutputKey == outputs.ClusterSharedNodeSecurityGroup {
					nodeSecurityGroupID = *output.OutputValue
					break
				}
			}
			Expect(nodeSecurityGroupID).NotTo(BeEmpty(), "ClusterSharedNodeSecurityGroup should be found in stack outputs")

			// Verify the security group has the expected karpenter.sh/discovery tag
			sgOutput, err := ec2Session.DescribeSecurityGroups(context.Background(), &awsec2.DescribeSecurityGroupsInput{
				GroupIds: []string{nodeSecurityGroupID},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(sgOutput.SecurityGroups).To(HaveLen(1))

			securityGroup := sgOutput.SecurityGroups[0]
			var foundKarpenterTag bool
			var karpenterTagValue string

			for _, tag := range securityGroup.Tags {
				if *tag.Key == "karpenter.sh/discovery" {
					foundKarpenterTag = true
					karpenterTagValue = *tag.Value
					break
				}
			}

			Expect(foundKarpenterTag).To(BeTrue(), "Security group should have karpenter.sh/discovery tag")
			Expect(karpenterTagValue).To(Equal(clusterName), "karpenter.sh/discovery tag value should match cluster name")
		})


	})

	Context("Creating a cluster with Karpenter without any tag", func() {
		params.LogStacksEventsOnFailure()

		It("should support karpenter and verify security group is NOT tagged when metadata tag is missing", func() {
			cmd := params.EksctlCreateCmd.
				WithArgs(
					"cluster",
					"--config-file=-",
					"--verbose=4",
					"--kubeconfig", params.KubeconfigPath,
				).
				WithoutArg("--region", params.Region).
				WithStdin(clusterutils.ReaderFromFile(clusterName, params.Region, "testdata/cluster-config.yaml"))
			Expect(cmd).To(RunSuccessfully())

			By("verifying Karpenter pods are healthy")
			kubeTest, err := kube.NewTest(params.KubeconfigPath)
			Expect(err).NotTo(HaveOccurred())
			// Check webhook pod
			Expect(kubeTest.WaitForPodsReady(karpenter.DefaultNamespace, metav1.ListOptions{
				LabelSelector: "app.kubernetes.io/instance=karpenter",
			}, 1, 10*time.Minute)).To(Succeed())

			By("verifying the cluster shared node security group does NOT have karpenter.sh/discovery tags")
			config := NewConfig(params.Region)
			cfnSession := cfn.NewFromConfig(config)
			ec2Session := awsec2.NewFromConfig(config)

			// Get the cluster stack to find the node security group
			describeStackOut, err := cfnSession.DescribeStacks(context.Background(), &cfn.DescribeStacksInput{
				StackName: aws.String(fmt.Sprintf("eksctl-%s-cluster", clusterName)),
			})
			Expect(err).NotTo(HaveOccurred())

			var nodeSecurityGroupID string
			for _, output := range describeStackOut.Stacks[0].Outputs {
				if *output.OutputKey == outputs.ClusterSharedNodeSecurityGroup {
					nodeSecurityGroupID = *output.OutputValue
					break
				}
			}
			Expect(nodeSecurityGroupID).NotTo(BeEmpty(), "ClusterSharedNodeSecurityGroup should be found in stack outputs")

			// Verify the security group does NOT have the karpenter.sh/discovery tag
			sgOutput, err := ec2Session.DescribeSecurityGroups(context.Background(), &awsec2.DescribeSecurityGroupsInput{
				GroupIds: []string{nodeSecurityGroupID},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(sgOutput.SecurityGroups).To(HaveLen(1))

			securityGroup := sgOutput.SecurityGroups[0]
			var foundKarpenterTag bool

			for _, tag := range securityGroup.Tags {
				if *tag.Key == "karpenter.sh/discovery" {
					foundKarpenterTag = true
					break
				}
			}

			Expect(foundKarpenterTag).To(BeFalse(), "Security group should NOT have karpenter.sh/discovery tag when only Karpenter is enabled without metadata tag")

			GinkgoWriter.Printf("Successfully verified Karpenter deployment without security group tagging\n")
			GinkgoWriter.Printf("Karpenter pods are healthy but security group %s does NOT have karpenter.sh/discovery tag\n", nodeSecurityGroupID)
		})
	})
})
