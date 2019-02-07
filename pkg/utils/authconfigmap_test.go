package utils_test

import (
	"log"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	. "github.com/weaveworks/eksctl/pkg/utils"
)

var _ = Describe("Auth ConfigMap", func() {
	Describe("create new ConfigMap", func() {

		cm, err := NewAuthConfigMap("arn:aws:iam::122333:role/eksctl-cluster-5a-nodegroup-ng1-p-NodeInstanceRole-NNH3ISP12CX")

		expected := []string{
			`- username: 'system:node:{{EC2PrivateDNSName}}'`,
			`  groups: [ 'system:bootstrappers', 'system:nodes' ]`,
			`  rolearn: 'arn:aws:iam::122333:role/eksctl-cluster-5a-nodegroup-ng1-p-NodeInstanceRole-NNH3ISP12CX'`,
		}

		It("should create correct configuration for a new nodegroup", func() {
			Expect(err).To(Not(HaveOccurred()))
			Expect(cm).To(Not(BeNil()))

			Expect(cm.ObjectMeta).To(Equal(metav1.ObjectMeta{
				Name:      AuthConfigMapName,
				Namespace: AuthConfigMapNamespace,
			}))

			Expect(cm.Data).To(HaveKey("mapRoles"))
			Expect(cm.Data["mapRoles"]).To(MatchYAML(strings.Join(expected, "\n")))
		})
		It("should add a new node group ARN to the configmap", func() {
			err = UpdateAuthConfigMap(cm, "arn:aws:iam::122333:role/eksctl-cluster-5a-nodegroup-ng2-p-NodeInstanceRole-1L35GCVYSTW4E")

			Expect(err).To(Not(HaveOccurred()))
			Expect(cm).To(Not(BeNil()))

			expected = append(expected,
				`- username: 'system:node:{{EC2PrivateDNSName}}'`,
				`  groups: [ 'system:bootstrappers', 'system:nodes' ]`,
				`  rolearn: 'arn:aws:iam::122333:role/eksctl-cluster-5a-nodegroup-ng2-p-NodeInstanceRole-1L35GCVYSTW4E'`,
			)

			Expect(cm.ObjectMeta).To(Equal(metav1.ObjectMeta{
				Name:      AuthConfigMapName,
				Namespace: AuthConfigMapNamespace,
			}))

			Expect(cm.Data).To(HaveKey("mapRoles"))
			Expect(cm.Data["mapRoles"]).To(MatchYAML(strings.Join(expected, "\n")))
		})

		It("should remove arn:aws:iam::122333:role/eksctl-cluster-5a-nodegroup-ng2-p-NodeInstanceRole-1L35GCVYSTW4E from ConfigMap", func() {
			err = RemoveNodeGroupFromAuthConfigMap(cm, "arn:aws:iam::122333:role/eksctl-cluster-5a-nodegroup-ng2-p-NodeInstanceRole-1L35GCVYSTW4E")

			Expect(err).To(Not(HaveOccurred()))
			Expect(cm).To(Not(BeNil()))

			expected = []string{
				`- username: 'system:node:{{EC2PrivateDNSName}}'`,
				`  groups: [ 'system:bootstrappers', 'system:nodes' ]`,
				`  rolearn: 'arn:aws:iam::122333:role/eksctl-cluster-5a-nodegroup-ng1-p-NodeInstanceRole-NNH3ISP12CX'`,
			}

			Expect(cm.ObjectMeta).To(Equal(metav1.ObjectMeta{
				Name:      AuthConfigMapName,
				Namespace: AuthConfigMapNamespace,
			}))

			Expect(cm.Data).To(HaveKey("mapRoles"))
			Expect(cm.Data["mapRoles"]).To(MatchYAML(strings.Join(expected, "\n")))
		})

		It("should fail if an role ARN is is not in the config map", func() {
			err = RemoveNodeGroupFromAuthConfigMap(cm, "arn:aws:iam::122333:role/eksctl-cluster-5a-nodegroup-ng1-p-NodeInstanceRole-ABCDEFGH'")

			Expect(err).To(HaveOccurred())
			Expect(cm.Data["mapRoles"]).To(MatchYAML(strings.Join(expected, "\n")))
		})

		It("should remove arn:aws:iam::122333:role/eksctl-cluster-5a-nodegroup-ng1-p-NodeInstanceRole-NNH3ISP12CX and make mapRoles be []", func() {
			err = RemoveNodeGroupFromAuthConfigMap(cm, "arn:aws:iam::122333:role/eksctl-cluster-5a-nodegroup-ng1-p-NodeInstanceRole-NNH3ISP12CX")

			Expect(err).To(Not(HaveOccurred()))
			Expect(cm).To(Not(BeNil()))

			Expect(cm.ObjectMeta).To(Equal(metav1.ObjectMeta{
				Name:      AuthConfigMapName,
				Namespace: AuthConfigMapNamespace,
			}))
			log.Println("post ", cm.Data["mapRoles"])
			Expect(cm.Data).To(HaveKey("mapRoles"))
			Expect(cm.Data["mapRoles"]).To(MatchYAML("[]"))
		})

		It("should fail if you try removing a role when the mapRole is empty", func() {
			err = RemoveNodeGroupFromAuthConfigMap(cm, "arn:aws:iam::122333:role/eksctl-cluster-5a-nodegroup-ng1-p-NodeInstanceRole-ABCDEFGH'")

			Expect(err).To(HaveOccurred())
			Expect(cm.Data).To(HaveKey("mapRoles"))
			Expect(cm.Data["mapRoles"]).To(MatchYAML("[]"))
		})
	})
})
