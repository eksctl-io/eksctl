package tests

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/weaveworks/eksctl/integration/matchers"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func AssertNodeTaints(nodeList *corev1.NodeList, expectedTaints []corev1.Taint) {
	//unset the time so the structs can be compared
	for _, node := range nodeList.Items {
		for _, t := range node.Spec.Taints {
			t.TimeAdded = nil
		}
	}

	for _, node := range nodeList.Items {
		for _, taint := range expectedTaints {
			Expect(node.Spec.Taints).To(ContainElement(taint))
		}
	}
}

func ListNodes(clientset kubernetes.Interface, nodeGroupName string) *corev1.NodeList {
	nodeList, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", api.NodeGroupNameLabel, nodeGroupName),
	})
	Expect(err).NotTo(HaveOccurred())
	return nodeList
}

// NodeInstanceIDs returns the EC2 instance IDs of a nodegroup's nodes, read from their
// Kubernetes provider IDs.
func NodeInstanceIDs(kubeConfig, nodeGroupName string) []string {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfig)
	Expect(err).NotTo(HaveOccurred())
	clientSet, err := kubernetes.NewForConfig(config)
	Expect(err).NotTo(HaveOccurred())
	nodes, err := clientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("%s=%s", api.NodeGroupNameLabel, nodeGroupName),
	})
	Expect(err).NotTo(HaveOccurred())
	var instanceIDs []string
	for _, node := range nodes.Items {
		// aws:///us-west-2c/i-00bb587a7011eb63c
		split := strings.Split(node.Spec.ProviderID, "/")
		id := split[len(split)-1]
		Expect(id).To(
			HavePrefix("i"),
			fmt.Sprintf("provider ID %q should have instance ID format aws:///us-west-2c/i-00bb587a7011eb63c", node.Spec.ProviderID),
		)
		instanceIDs = append(instanceIDs, id)
	}
	return instanceIDs
}

// NodeInstances returns the EC2 instances backing a nodegroup's nodes.
func NodeInstances(kubeConfig, region, nodeGroupName string) []ec2types.Instance {
	instanceIDs := NodeInstanceIDs(kubeConfig, nodeGroupName)
	Expect(instanceIDs).NotTo(BeEmpty(), fmt.Sprintf("nodegroup %q should have joined nodes", nodeGroupName))
	ec2API := awsec2.NewFromConfig(matchers.NewConfig(region))
	output, err := ec2API.DescribeInstances(context.Background(), &awsec2.DescribeInstancesInput{
		InstanceIds: instanceIDs,
	})
	Expect(err).NotTo(HaveOccurred())
	var instances []ec2types.Instance
	for _, reservation := range output.Reservations {
		instances = append(instances, reservation.Instances...)
	}
	Expect(instances).To(HaveLen(len(instanceIDs)))
	return instances
}

// PrimaryNetworkInterface returns an instance's primary network interface, the one at device
// index 0 on network card 0. That is the interface the VPC CNI reads connection tracking
// settings from before replicating them onto the interfaces it creates for pods.
func PrimaryNetworkInterface(instance ec2types.Instance) *ec2types.InstanceNetworkInterface {
	for i, networkInterface := range instance.NetworkInterfaces {
		attachment := networkInterface.Attachment
		if attachment == nil {
			continue
		}
		if aws.ToInt32(attachment.DeviceIndex) == 0 && aws.ToInt32(attachment.NetworkCardIndex) == 0 {
			return &instance.NetworkInterfaces[i]
		}
	}
	return nil
}

func AssertNodeVolumes(kubeConfig, region, nodeGroupName, volumeName string) {
	instanceIDs := NodeInstanceIDs(kubeConfig, nodeGroupName)
	cfg := matchers.NewConfig(region)
	ec2 := awsec2.NewFromConfig(cfg)
	instances, err := ec2.DescribeInstances(context.Background(), &awsec2.DescribeInstancesInput{
		InstanceIds: instanceIDs,
	})
	Expect(err).NotTo(HaveOccurred())
	for _, res := range instances.Reservations {
		var deviceNames []string
		for _, instance := range res.Instances {
			for _, mapping := range instance.BlockDeviceMappings {
				deviceNames = append(deviceNames, *mapping.DeviceName)
			}
			Expect(deviceNames).To(ContainElement(volumeName))
		}
	}
}
