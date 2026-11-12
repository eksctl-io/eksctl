package authconfigmap_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/fake"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	. "github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/iam"
)

var _ = DescribeTable("AddNodeGroup()", func(amiFamily, expectedUsername string) {
	Expect(NodeGroupUsername(amiFamily)).To(Equal(expectedUsername))

	ng := api.NewNodeGroup()
	ng.Name = "ng"
	ng.AMIFamily = amiFamily
	ng.IAM.InstanceRoleARN = roleA

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
	Entry(api.NodeImageFamilyAmazonLinux2027, api.NodeImageFamilyAmazonLinux2027, "system:node:{{SessionName}}"),
	Entry(api.NodeImageFamilyAmazonLinux2023, api.NodeImageFamilyAmazonLinux2023, "system:node:{{EC2PrivateDNSName}}"),
)
