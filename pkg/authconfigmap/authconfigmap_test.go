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
	roleA    = "arn:aws:iam::122333:role/eksctl-cluster-5a-nodegroup-ng1-p-NodeInstanceRole-NNH3ISP12CX"
	roleB    = "arn:aws:iam::122333:role/eksctl-cluster-5a-nodegroup-ng1-p-NodeInstanceRole-ABCDEFGH"
	groupB   = "foo"
	accountA = "123"
	accountB = "789"
)

var (
	expectedA = makeExpectedRole(roleA, RoleNodeGroupGroups)
	expectedB = makeExpectedRole(roleB, []string{groupB})
)

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

func (c *mockClient) reset() {
	c.updated = nil
	c.created = nil
}

func makeExpectedRole(arn string, groups []string) string {
	return fmt.Sprintf(`- rolearn: %s
  username: system:node:{{EC2PrivateDNSName}}
  groups:
  - %s
`, arn, strings.Join(groups, "\n  - "))
}

func makeExpectedAccounts(accounts ...string) string {
	var y string
	for _, a := range accounts {
		// Having them quoted is important for the yaml parser to
		// recognize them as strings over numbers
		y += fmt.Sprintf("\n- %q", a)
	}

	return y
}

var _ = Describe("AuthConfigMap{}", func() {
	Describe("New()", func() {
		It("should create an empty configmap", func() {
			client := &mockClient{}
			acm := New(client, nil)
			err := acm.Save()
			Expect(err).NotTo(HaveOccurred())
			Expect(client.updated).To(BeNil())

			// Created!
			cm := client.created
			expected := ObjectMeta()
			expected.UID = cm.UID
			Expect(cm.ObjectMeta).To(Equal(expected))
			Expect(cm.Data["mapRoles"]).To(Equal(""))
		})
		It("should load an empty configmap", func() {
			empty := &corev1.ConfigMap{}

			client := &mockClient{}
			acm := New(client, empty)
			err := acm.Save()
			Expect(err).NotTo(HaveOccurred())
			Expect(client.updated).To(BeNil())

			// Created!
			cm := client.created
			expected := ObjectMeta()
			expected.UID = cm.UID
			Expect(cm.ObjectMeta).To(Equal(expected))
			Expect(cm.Data["mapRoles"]).To(Equal(""))
		})
		It("should load an existing configmap", func() {
			existing := &corev1.ConfigMap{
				ObjectMeta: ObjectMeta(),
				Data:       map[string]string{},
			}
			existing.ObjectMeta.UID = "123456"

			client := &mockClient{}
			acm := New(client, existing)
			err := acm.Save()
			Expect(err).NotTo(HaveOccurred())
			Expect(client.created).To(BeNil())

			// Updated!
			cm := client.updated
			Expect(cm.ObjectMeta.UID).To(Equal(types.UID("123456")))
		})
	})
	Describe("AddRole()", func() {
		existing := &corev1.ConfigMap{
			ObjectMeta: ObjectMeta(),
			Data:       map[string]string{},
		}
		existing.UID = "123456"
		client := &mockClient{}
		acm := New(client, existing)

		addAndSave := func(arn string, groups []string) *corev1.ConfigMap {
			client.reset()
			err := acm.AddRole(arn, RoleNodeGroupUsername, groups)
			Expect(err).NotTo(HaveOccurred())

			err = acm.Save()
			Expect(err).NotTo(HaveOccurred())
			Expect(client.created).To(BeNil())
			Expect(client.updated).NotTo(BeNil())

			return client.updated
		}

		It("should add a role", func() {
			cm := addAndSave(roleA, RoleNodeGroupGroups)
			Expect(cm.Data["mapRoles"]).To(MatchYAML(expectedA))
		})
		It("should append a second role", func() {
			cm := addAndSave(roleB, []string{groupB})
			Expect(cm.Data["mapRoles"]).To(MatchYAML(expectedA + expectedB))
		})
		It("should append a duplicate role", func() {
			cm := addAndSave(roleA, RoleNodeGroupGroups)
			expected := expectedA + expectedB + expectedA
			Expect(cm.Data["mapRoles"]).To(MatchYAML(expected))
		})
	})
	Describe("RemoveRole()", func() {
		existing := &corev1.ConfigMap{
			ObjectMeta: ObjectMeta(),
			Data:       map[string]string{"mapRoles": expectedA + expectedA + expectedB},
		}
		existing.UID = "123456"
		client := &mockClient{}
		acm := New(client, existing)

		removeAndSave := func(arn string) *corev1.ConfigMap {
			client.reset()
			err := acm.RemoveRole(arn, false)
			Expect(err).NotTo(HaveOccurred())

			err = acm.Save()
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
			err := acm.RemoveRole(roleA, false)
			Expect(err).To(HaveOccurred())
		})
		It("should remove all if specified", func() {
			err := acm.AddRole(roleA, RoleNodeGroupUsername, RoleNodeGroupGroups)
			Expect(err).NotTo(HaveOccurred())
			err = acm.AddRole(roleA, RoleNodeGroupUsername, RoleNodeGroupGroups)
			Expect(err).NotTo(HaveOccurred())
			Expect(client.updated.Data["mapRoles"]).To(Not(MatchYAML("[]")))

			err = acm.RemoveRole(roleA, true)
			Expect(err).NotTo(HaveOccurred())
			Expect(client.updated.Data["mapRoles"]).To(MatchYAML("[]"))
		})
	})
	Describe("AddAccount()", func() {
		existing := &corev1.ConfigMap{
			ObjectMeta: ObjectMeta(),
			Data:       map[string]string{},
		}
		existing.UID = "123456"
		client := &mockClient{}
		acm := New(client, existing)

		addAndSave := func(account string) *corev1.ConfigMap {
			client.reset()
			err := acm.AddAccount(account)
			Expect(err).NotTo(HaveOccurred())

			err = acm.Save()
			Expect(err).NotTo(HaveOccurred())
			Expect(client.created).To(BeNil())
			Expect(client.updated).NotTo(BeNil())

			return client.updated
		}

		It("should add an account", func() {
			cm := addAndSave(accountA)
			Expect(cm.Data["mapAccounts"]).To(MatchYAML(makeExpectedAccounts(accountA)))
		})
		It("should deduplicate when adding", func() {
			cm := addAndSave(accountA)
			Expect(cm.Data["mapAccounts"]).To(MatchYAML(makeExpectedAccounts(accountA)))
		})
		It("should add another account", func() {
			cm := addAndSave(accountB)
			Expect(cm.Data["mapAccounts"]).To(MatchYAML(makeExpectedAccounts(accountA) + makeExpectedAccounts(accountB)))
		})
	})
	Describe("RemoveAccount()", func() {
		existing := &corev1.ConfigMap{
			ObjectMeta: ObjectMeta(),
			Data:       map[string]string{"mapAccounts": makeExpectedAccounts(accountA) + makeExpectedAccounts(accountB)},
		}
		existing.UID = "123456"
		client := &mockClient{}
		acm := New(client, existing)

		removeAndSave := func(account string) *corev1.ConfigMap {
			client.reset()
			err := acm.RemoveAccount(account)
			Expect(err).NotTo(HaveOccurred())

			err = acm.Save()
			Expect(err).NotTo(HaveOccurred())
			Expect(client.created).To(BeNil())
			Expect(client.updated).NotTo(BeNil())

			return client.updated
		}

		It("should remove an account", func() {
			cm := removeAndSave(accountA)
			Expect(cm.Data["mapAccounts"]).To(MatchYAML(makeExpectedAccounts(accountB)))
		})
		It("should fail if account not found", func() {
			err := acm.RemoveAccount(accountA)
			Expect(err).To(HaveOccurred())
		})
	})
})
