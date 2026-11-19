package authconfigmap_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
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
