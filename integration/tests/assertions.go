package tests

import (
	"context"
	"fmt"
	"strings"

	awsec2 "github.com/aws/aws-sdk-go-v2/service/ec2"
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

func AssertNodeVolumes(kubeConfig, region, nodeGroupName, volumeName string) {
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
