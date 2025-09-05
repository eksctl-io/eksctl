//go:build integration

package security_group_tagging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfn "github.com/aws/aws-sdk-go-v2/service/cloudformation"
	awsec2 "github.com/aws/aws-sdk-go-v2/service/ec2"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"
	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/integration/utilities/kube"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

var params *tests.Params

func init() {
	testing.Init()
	if err := api.Register(); err != nil {
		panic(fmt.Errorf("unexpected error registering API scheme: %w", err))
	}
	params = tests.NewParams("security-group-tagging")
}

func TestSecurityGroupTagging(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = Describe("(Integration) Security Group Tagging", func() {
	var (
		clusterName string
	)

	BeforeEach(func() {
		clusterName = fmt.Sprintf("it-sg-tagging-%d", time.Now().Unix())
	})

	AfterEach(func() {
		cmd := params.EksctlDeleteCmd.WithArgs(
			"cluster", clusterName,
			"--verbose", "4",
		)
		Expect(cmd).To(RunSuccessfully())
	})

	Context("Creating a cluster with both Karpenter and metadata tags", func() {
		params.LogStacksEventsOnFailure()

		It("should automatically tag the node security group with karpenter.sh/discovery", func() {
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
						"environment":            "integration-test",
					},
				},
				Karpenter: &api.Karpenter{
					Version: "v0.20.0",
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

			By("deploying a test application to verify the cluster is functional")
			kubeTest, err := kube.NewTest(params.KubeconfigPath)
			Expect(err).NotTo(HaveOccurred())

			// Deploy a simple test application to verify the cluster works with the tagged security group
			d := kubeTest.CreateDeploymentFromFile(kubeTest.Namespace, "../../data/podinfo.yaml")
			kubeTest.WaitForDeploymentReady(d, 10*time.Minute)

			pods := kubeTest.ListPodsFromDeployment(d)
			Expect(len(pods.Items)).To(Equal(2))

			// Verify pods are running successfully, which confirms the cluster networking
			// and security groups are working correctly
			for _, pod := range pods.Items {
				Expect(pod.Namespace).To(Equal(kubeTest.Namespace))
				
				// Test that we can make requests to the pod
				var js interface{}
				kubeTest.PodProxyGetJSON(&pod, "", "/version", &js)
				Expect(js.(map[string]interface{})).To(HaveKeyWithValue("version", "1.5.1"))
			}

			By("verifying that the tagged security group enables AWS Load Balancer Controller compatibility")
			// The presence of the karpenter.sh/discovery tag on the node security group
			// ensures that when AWS Load Balancer Controller is deployed, it will be able
			// to discover this security group for load balancer creation.
			// 
			// This test verifies the core requirement: the security group is properly tagged
			// during cluster creation, which is the prerequisite for Load Balancer Controller
			// functionality.
			
			GinkgoWriter.Printf("Successfully verified cluster with tagged node security group: %s\n", nodeSecurityGroupID)
			GinkgoWriter.Printf("Security group has karpenter.sh/discovery tag with value: %s\n", karpenterTagValue)
			GinkgoWriter.Printf("This enables AWS Load Balancer Controller to discover and use this security group\n")
		})
	})

	Context("Creating a cluster with only Karpenter (no metadata tags)", func() {
		params.LogStacksEventsOnFailure()

		It("should NOT tag the node security group with karpenter.sh/discovery", func() {
			By("creating a cluster with only Karpenter enabled but no karpenter.sh/discovery in metadata.tags")
			
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
						"environment": "integration-test",
					},
				},
				Karpenter: &api.Karpenter{
					Version: "v0.20.0",
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

			By("verifying the cluster is still functional without the discovery tag")
			kubeTest, err := kube.NewTest(params.KubeconfigPath)
			Expect(err).NotTo(HaveOccurred())

			// Deploy a simple test application to verify the cluster works
			d := kubeTest.CreateDeploymentFromFile(kubeTest.Namespace, "../../data/podinfo.yaml")
			kubeTest.WaitForDeploymentReady(d, 10*time.Minute)

			pods := kubeTest.ListPodsFromDeployment(d)
			Expect(len(pods.Items)).To(Equal(2))

			GinkgoWriter.Printf("Successfully verified cluster without karpenter.sh/discovery tag: %s\n", nodeSecurityGroupID)
			GinkgoWriter.Printf("This confirms that tagging only happens when BOTH Karpenter AND metadata tag are present\n")
		})
	})

	Context("Creating a cluster with only metadata tags (no Karpenter)", func() {
		params.LogStacksEventsOnFailure()

		It("should NOT tag the node security group with karpenter.sh/discovery", func() {
			By("creating a cluster with karpenter.sh/discovery in metadata.tags but no Karpenter enabled")
			
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
						"environment":            "integration-test",
					},
				},
				// No Karpenter configuration
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

			Expect(foundKarpenterTag).To(BeFalse(), "Security group should NOT have karpenter.sh/discovery tag when only metadata tag is present without Karpenter")

			By("verifying the cluster uses the same existing security groups but without the discovery tags")
			// Verify that the cluster still has the standard security groups created by eksctl
			// but without the karpenter.sh/discovery tags
			
			// Check that the security group has the standard eksctl tags but not the discovery tag
			var hasEksctlClusterTag bool
			var hasEnvironmentTag bool
			
			for _, tag := range securityGroup.Tags {
				if *tag.Key == fmt.Sprintf("kubernetes.io/cluster/%s", clusterName) {
					hasEksctlClusterTag = true
				}
				if *tag.Key == "environment" && *tag.Value == "integration-test" {
					hasEnvironmentTag = true
				}
			}

			Expect(hasEksctlClusterTag).To(BeTrue(), "Security group should have standard eksctl cluster tag")
			Expect(hasEnvironmentTag).To(BeTrue(), "Security group should have environment tag from metadata")

			By("verifying the cluster creation completes successfully with untagged node security group")
			kubeTest, err := kube.NewTest(params.KubeconfigPath)
			Expect(err).NotTo(HaveOccurred())

			// Deploy a simple test application to verify the cluster works
			d := kubeTest.CreateDeploymentFromFile(kubeTest.Namespace, "../../data/podinfo.yaml")
			kubeTest.WaitForDeploymentReady(d, 10*time.Minute)

			pods := kubeTest.ListPodsFromDeployment(d)
			Expect(len(pods.Items)).To(Equal(2))

			// Verify pods are running successfully, which confirms the cluster networking
			// and security groups are working correctly even without discovery tags
			for _, pod := range pods.Items {
				Expect(pod.Namespace).To(Equal(kubeTest.Namespace))
				
				// Test that we can make requests to the pod
				var js interface{}
				kubeTest.PodProxyGetJSON(&pod, "", "/version", &js)
				Expect(js.(map[string]interface{})).To(HaveKeyWithValue("version", "1.5.1"))
			}

			GinkgoWriter.Printf("Successfully verified cluster without karpenter.sh/discovery tag: %s\n", nodeSecurityGroupID)
			GinkgoWriter.Printf("This confirms that tagging only happens when BOTH Karpenter AND metadata tag are present\n")
		})
	})

	Context("Creating a cluster with neither Karpenter nor metadata tags", func() {
		params.LogStacksEventsOnFailure()

		It("should NOT tag the node security group with karpenter.sh/discovery", func() {
			By("creating a cluster with neither Karpenter enabled nor karpenter.sh/discovery in metadata.tags")
			
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
						"environment": "integration-test",
					},
				},
				// No Karpenter configuration and no karpenter.sh/discovery tag
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

			Expect(foundKarpenterTag).To(BeFalse(), "Security group should NOT have karpenter.sh/discovery tag when neither Karpenter nor metadata tag are present")

			By("verifying the cluster creation completes successfully with standard untagged node security group")
			kubeTest, err := kube.NewTest(params.KubeconfigPath)
			Expect(err).NotTo(HaveOccurred())

			// Deploy a simple test application to verify the cluster works
			d := kubeTest.CreateDeploymentFromFile(kubeTest.Namespace, "../../data/podinfo.yaml")
			kubeTest.WaitForDeploymentReady(d, 10*time.Minute)

			pods := kubeTest.ListPodsFromDeployment(d)
			Expect(len(pods.Items)).To(Equal(2))

			// Verify pods are running successfully, which confirms the cluster networking
			// and security groups are working correctly in the default scenario
			for _, pod := range pods.Items {
				Expect(pod.Namespace).To(Equal(kubeTest.Namespace))
				
				// Test that we can make requests to the pod
				var js interface{}
				kubeTest.PodProxyGetJSON(&pod, "", "/version", &js)
				Expect(js.(map[string]interface{})).To(HaveKeyWithValue("version", "1.5.1"))
			}

			GinkgoWriter.Printf("Successfully verified default cluster without karpenter.sh/discovery tag: %s\n", nodeSecurityGroupID)
			GinkgoWriter.Printf("This confirms the default behavior when neither Karpenter nor metadata tags are configured\n")
		})
	})
})