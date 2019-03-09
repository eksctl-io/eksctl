package authconfigmap_test

import (
	"fmt"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/typed/core/v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/weaveworks/eksctl/pkg/authconfigmap"
)

const (
	roleA  = "arn:aws:iam::122333:role/eksctl-cluster-5a-nodegroup-ng1-p-NodeInstanceRole-NNH3ISP12CX"
	roleB  = "arn:aws:iam::122333:role/eksctl-cluster-5a-nodegroup-ng1-p-NodeInstanceRole-ABCDEFGH"
	groupB = "foo"
)

var (
	expectedA = makeExpectedRole(roleA, DefaultNodeGroups)
	expectedB = makeExpectedRole(roleB, []string{groupB})
)

func makeExpectedRole(arn string, groups []string) string {
	return fmt.Sprintf(`- rolearn: %s
  username: system:node:{{EC2PrivateDNSName}}
  groups:
  - %s
`, arn, strings.Join(groups, "\n  - "))
}

// mockClient implements v1.ConfigMapInterface
type mockClient struct {
	v1.ConfigMapInterface
	created *corev1.ConfigMap
	updated *corev1.ConfigMap
}

func (c *mockClient) Create(cm *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	cm.ObjectMeta.UID = "18b9e60c-2057-11e7-8868-0eba8ef9df1a"
	c.created = cm
	return cm, nil
}

func (c *mockClient) Update(cm *corev1.ConfigMap) (*corev1.ConfigMap, error) {
	c.updated = cm
	return cm, nil
}

var _ = Describe("AuthConfigMap{}", func() {
	Describe("New()", func() {
		It("should create an empty configmap", func() {
			acm := New(nil)
			client := &mockClient{}
			err := acm.Save(client)
			Expect(err).NotTo(HaveOccurred())
			Expect(client.updated).To(BeNil())

			// Created!
			cm := client.created
			om := ObjectMeta()
			om.UID = cm.UID
			Expect(cm.ObjectMeta).To(Equal(om))
			Expect(cm.Data["mapRoles"]).To(Equal(""))
		})
		It("should load an existing configmap", func() {
			cm := &corev1.ConfigMap{
				ObjectMeta: ObjectMeta(),
				Data:       map[string]string{},
			}
			cm.ObjectMeta.UID = "123456"

			acm := New(cm)
			client := &mockClient{}
			err := acm.Save(client)
			Expect(err).NotTo(HaveOccurred())
			Expect(client.created).To(BeNil())

			// Updated!
			cm = client.updated
			Expect(cm.ObjectMeta.UID).To(Equal(types.UID("123456")))
		})
	})
	Describe("AddRole()", func() {
		cm := &corev1.ConfigMap{
			ObjectMeta: ObjectMeta(),
			Data:       map[string]string{},
		}
		cm.UID = "123456"
		acm := New(cm)

		addAndSave := func(arn string, groups []string) *corev1.ConfigMap {
			err := acm.AddRole(arn, groups)
			Expect(err).NotTo(HaveOccurred())

			client := &mockClient{}
			err = acm.Save(client)
			Expect(err).NotTo(HaveOccurred())
			Expect(client.created).To(BeNil())
			Expect(client.updated).NotTo(BeNil())

			return client.updated
		}

		It("should add a role", func() {
			cm := addAndSave(roleA, DefaultNodeGroups)
			Expect(cm.Data["mapRoles"]).To(MatchYAML(expectedA))
		})
		It("should append a second role", func() {
			cm := addAndSave(roleB, []string{groupB})
			Expect(cm.Data["mapRoles"]).To(MatchYAML(expectedA + expectedB))
		})
		It("should append a duplicate role", func() {
			cm := addAndSave(roleA, DefaultNodeGroups)
			expected := expectedA + expectedB + expectedA
			Expect(cm.Data["mapRoles"]).To(MatchYAML(expected))
		})
	})
	Describe("RemoveRole()", func() {
		cm := &corev1.ConfigMap{
			ObjectMeta: ObjectMeta(),
			Data:       map[string]string{"mapRoles": expectedA + expectedA + expectedB},
		}
		cm.UID = "123456"
		acm := New(cm)

		removeAndSave := func(arn string) *corev1.ConfigMap {
			err := acm.RemoveRole(arn)
			Expect(err).NotTo(HaveOccurred())

			client := &mockClient{}
			err = acm.Save(client)
			Expect(err).NotTo(HaveOccurred())

			return client.updated
		}

		It("should remove role", func() {
			cm := removeAndSave(roleB)
			Expect(cm.Data["mapRoles"]).To(MatchYAML(expectedA + expectedA))
		})
		It("should remove one role for duplicates", func() {
			cm := removeAndSave(roleA)
			Expect(cm.Data["mapRoles"]).To(MatchYAML(expectedA))
		})
		It("should remove last role", func() {
			cm := removeAndSave(roleA)
			Expect(cm.Data["mapRoles"]).To(MatchYAML("[]"))
		})
		It("should fail if role not found", func() {
			err := acm.RemoveRole(roleA)
			Expect(err).To(HaveOccurred())
		})
	})
})
