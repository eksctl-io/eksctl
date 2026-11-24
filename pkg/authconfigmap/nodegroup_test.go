package authconfigmap_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	. "github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/iam"
)

// instanceIDNodeNameConfig returns a nodeadm NodeConfig that sets the InstanceIdNodeName feature gate.
func instanceIDNodeNameConfig(enabled bool) *string {
	config := fmt.Sprintf(`apiVersion: node.eks.aws/v1alpha1
kind: NodeConfig
spec:
  featureGates:
    InstanceIdNodeName: %t
`, enabled)
	return &config
}

var _ = DescribeTable("AddNodeGroup()", func(amiFamily string, overrideBootstrapCommand *string, expectedUsername string) {
	ng := api.NewNodeGroup()
	ng.Name = "ng"
	ng.AMIFamily = amiFamily
	ng.OverrideBootstrapCommand = overrideBootstrapCommand
	ng.IAM.InstanceRoleARN = roleA
	Expect(NodeGroupUsername(ng.NodeGroupBase)).To(Equal(expectedUsername))

	clientSet := fake.NewSimpleClientset()
	Expect(AddNodeGroup(clientSet, ng)).To(Succeed())

	acm, err := NewFromClientSet(clientSet)
	Expect(err).NotTo(HaveOccurred())
	identities, err := acm.GetIdentities()
	Expect(err).NotTo(HaveOccurred())
	Expect(identities).To(ConsistOf(iam.RoleIdentity{
		RoleARN: roleA,
		KubernetesIdentity: iam.KubernetesIdentity{
			KubernetesUsername: expectedUsername,
			KubernetesGroups:   RoleNodeGroupGroups,
		},
	}))
},
	Entry("AmazonLinux2027 by default", api.NodeImageFamilyAmazonLinux2027, (*string)(nil), "system:node:{{SessionName}}"),
	Entry("AmazonLinux2027 with InstanceIdNodeName off", api.NodeImageFamilyAmazonLinux2027, instanceIDNodeNameConfig(false), "system:node:{{EC2PrivateDNSName}}"),
	Entry("AmazonLinux2023 by default", api.NodeImageFamilyAmazonLinux2023, (*string)(nil), "system:node:{{EC2PrivateDNSName}}"),
	Entry("AmazonLinux2023 with InstanceIdNodeName on", api.NodeImageFamilyAmazonLinux2023, instanceIDNodeNameConfig(true), "system:node:{{SessionName}}"),
)

var _ = DescribeTable("AddNodeGroup() with a role already used by another nodegroup", func(amiFamily string, expectedErr string) {
	// Start with an existing auth ConfigMap. The fake client doesn't set a UID on create, so saving twice would fail.
	clientSet := fake.NewSimpleClientset(&corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{Name: ObjectName, Namespace: ObjectNamespace, UID: "123"},
	})
	existing := api.NewNodeGroup()
	existing.Name = "al2023"
	existing.AMIFamily = api.NodeImageFamilyAmazonLinux2023
	existing.IAM.InstanceRoleARN = roleA
	Expect(AddNodeGroup(clientSet, existing)).To(Succeed())

	ng := api.NewNodeGroup()
	ng.Name = "new"
	ng.AMIFamily = amiFamily
	ng.IAM.InstanceRoleARN = roleA
	err := AddNodeGroup(clientSet, ng)

	acm, aerr := NewFromClientSet(clientSet)
	Expect(aerr).NotTo(HaveOccurred())
	identities, aerr := acm.GetIdentities()
	Expect(aerr).NotTo(HaveOccurred())
	original := iam.RoleIdentity{
		RoleARN: roleA,
		KubernetesIdentity: iam.KubernetesIdentity{
			KubernetesUsername: RoleNodeGroupUsername,
			KubernetesGroups:   RoleNodeGroupGroups,
		},
	}
	if expectedErr != "" {
		Expect(err).To(MatchError(ContainSubstring(expectedErr)))
		Expect(identities).To(ConsistOf(original))
	} else {
		Expect(err).NotTo(HaveOccurred())
		Expect(identities).To(ConsistOf(original, original))
	}
},
	Entry("rejects a different node username and keeps the original mapping", api.NodeImageFamilyAmazonLinux2027, "use a separate instance role"),
	Entry("allows the same node username", api.NodeImageFamilyAmazonLinux2023, ""),
)
