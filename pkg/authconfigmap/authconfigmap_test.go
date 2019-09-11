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
	"github.com/weaveworks/eksctl/pkg/iam"
)

const (
	roleA = "arn:aws:iam::122333:role/eksctl-cluster-5a-nodegroup-ng1-p-NodeInstanceRole-NNH3ISP12CX"
	roleB = "arn:aws:iam::122333:role/eksctl-cluster-5a-nodegroup-ng1-p-NodeInstanceRole-ABCDEFGH"

	userA         = "arn:aws:iam::122333:user/alice"
	userAUsername = "alice"

	userB         = "arn:aws:iam::122333:user/bob"
	userBUsername = "bob"

	groupB   = "foo"
	accountA = "123"
	accountB = "789"
)

var (
	userAGroups = []string{"cryptographers", "tin-foil-hat-wearers"}
	userBGroups = []string{"cryptographers", "private-messages-authors", "dislikers-of-eve"}

	expectedRoleA = makeExpectedRole(roleA, RoleNodeGroupGroups)
	expectedRoleB = makeExpectedRole(roleB, []string{groupB})

	expectedUserA = makeExpectedUser(userA, "alice", "cryptographers", "tin-foil-hat-wearers")
	expectedUserB = makeExpectedUser(userB, "bob", "cryptographers", "private-messages-authors", "dislikers-of-eve")
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

func makeExpectedUser(arn, user string, groups ...string) string {
	return fmt.Sprintf(`- userarn: %s
  username: %s
  groups:
  - %s
`, arn, user, strings.Join(groups, "\n  - "))

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

func mustIdentity(arn iam.ARN, username string, groups []string) iam.Identity {
	id, err := iam.NewIdentity(arn, username, groups)
	Expect(err).ToNot(HaveOccurred())
	if id != nil {
		return *id
	}
	return iam.Identity{}
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

		addAndSave := func(canonicalArn string, groups []string) *corev1.ConfigMap {
			client.reset()

			arn, err := iam.Parse(canonicalArn)
			Expect(err).NotTo(HaveOccurred())

			err = acm.AddIdentity(mustIdentity(arn, RoleNodeGroupUsername, groups))
			Expect(err).NotTo(HaveOccurred())

			err = acm.Save()
			Expect(err).NotTo(HaveOccurred())
			Expect(client.created).To(BeNil())
			Expect(client.updated).NotTo(BeNil())

			return client.updated
		}

		It("should add a role", func() {
			cm := addAndSave(roleA, RoleNodeGroupGroups)
			Expect(cm.Data["mapRoles"]).To(MatchYAML(expectedRoleA))
		})
		It("should append a second role", func() {
			cm := addAndSave(roleB, []string{groupB})
			Expect(cm.Data["mapRoles"]).To(MatchYAML(expectedRoleA + expectedRoleB))
		})
		It("should append a duplicate role", func() {
			cm := addAndSave(roleA, RoleNodeGroupGroups)
			expected := expectedRoleA + expectedRoleB + expectedRoleA
			Expect(cm.Data["mapRoles"]).To(MatchYAML(expected))
		})
	})
	Describe("RemoveRole()", func() {
		existing := &corev1.ConfigMap{
			ObjectMeta: ObjectMeta(),
			Data:       map[string]string{"mapRoles": expectedRoleA + expectedRoleA + expectedRoleB},
		}
		existing.UID = "123456"
		client := &mockClient{}
		acm := New(client, existing)

		removeAndSave := func(canonicalArn string) *corev1.ConfigMap {
			client.reset()

			arn, err := iam.Parse(canonicalArn)
			Expect(err).NotTo(HaveOccurred())

			err = acm.RemoveIdentity(arn, false)
			Expect(err).NotTo(HaveOccurred())

			err = acm.Save()
			Expect(err).NotTo(HaveOccurred())

			return client.updated
		}

		roleAArn, err := iam.Parse(roleA)
		if err != nil {
			panic(err)
		}

		It("should remove role", func() {
			cm := removeAndSave(roleB)
			Expect(cm.Data["mapRoles"]).To(MatchYAML(expectedRoleA + expectedRoleA))
		})
		It("should remove one role for duplicates", func() {
			cm := removeAndSave(roleA)
			Expect(cm.Data["mapRoles"]).To(MatchYAML(expectedRoleA))
		})
		It("should remove last role", func() {
			cm := removeAndSave(roleA)
			Expect(cm.Data["mapRoles"]).To(MatchYAML("[]"))
		})
		It("should fail if role not found", func() {
			err := acm.RemoveIdentity(roleAArn, false)
			Expect(err).To(HaveOccurred())
		})
		It("should remove all if specified", func() {
			err := acm.AddIdentity(mustIdentity(roleAArn, RoleNodeGroupUsername, RoleNodeGroupGroups))
			Expect(err).NotTo(HaveOccurred())
			err = acm.AddIdentity(mustIdentity(roleAArn, RoleNodeGroupUsername, RoleNodeGroupGroups))
			Expect(err).NotTo(HaveOccurred())
			Expect(client.updated.Data["mapRoles"]).To(Not(MatchYAML("[]")))

			err = acm.RemoveIdentity(roleAArn, true)
			Expect(err).NotTo(HaveOccurred())
			Expect(client.updated.Data["mapRoles"]).To(MatchYAML("[]"))
		})
	})
	Describe("AddUser()", func() {
		existing := &corev1.ConfigMap{
			ObjectMeta: ObjectMeta(),
			Data:       map[string]string{},
		}
		existing.UID = "123456"
		client := &mockClient{}
		acm := New(client, existing)

		addAndSave := func(canonicalArn, user string, groups []string) *corev1.ConfigMap {
			client.reset()

			arn, err := iam.Parse(canonicalArn)
			Expect(err).NotTo(HaveOccurred())

			err = acm.AddIdentity(mustIdentity(arn, user, groups))
			Expect(err).NotTo(HaveOccurred())

			err = acm.Save()
			Expect(err).NotTo(HaveOccurred())
			Expect(client.created).To(BeNil())
			Expect(client.updated).NotTo(BeNil())

			return client.updated
		}

		It("should add a user", func() {
			cm := addAndSave(userA, userAUsername, userAGroups)
			Expect(cm.Data["mapUsers"]).To(MatchYAML(expectedUserA))
		})
		It("should append a second user", func() {
			cm := addAndSave(userB, userBUsername, userBGroups)
			Expect(cm.Data["mapUsers"]).To(MatchYAML(expectedUserA + expectedUserB))
		})
		It("should append a duplicate user", func() {
			cm := addAndSave(userA, userAUsername, userAGroups)
			expected := expectedUserA + expectedUserB + expectedUserA
			Expect(cm.Data["mapUsers"]).To(MatchYAML(expected))
		})
	})
	Describe("RemoveUser()", func() {
		existing := &corev1.ConfigMap{
			ObjectMeta: ObjectMeta(),
			Data:       map[string]string{"mapUsers": expectedUserA + expectedUserA + expectedUserB},
		}
		existing.UID = "123456"
		client := &mockClient{}
		acm := New(client, existing)

		removeAndSave := func(canonicalArn string) *corev1.ConfigMap {
			client.reset()

			arn, err := iam.Parse(canonicalArn)
			Expect(err).NotTo(HaveOccurred())

			err = acm.RemoveIdentity(arn, false)
			Expect(err).NotTo(HaveOccurred())

			err = acm.Save()
			Expect(err).NotTo(HaveOccurred())

			return client.updated
		}

		userAArn, err := iam.Parse(userA)
		if err != nil {
			panic(err)
		}

		It("should remove user", func() {
			cm := removeAndSave(userB)
			Expect(cm.Data["mapUsers"]).To(MatchYAML(expectedUserA + expectedUserA))
		})
		It("should remove one user for duplicates", func() {
			cm := removeAndSave(userA)
			Expect(cm.Data["mapUsers"]).To(MatchYAML(expectedUserA))
		})
		It("should remove last user", func() {
			cm := removeAndSave(userA)
			Expect(cm.Data["mapUsers"]).To(MatchYAML("[]"))
		})
		It("should fail if user not found", func() {
			err := acm.RemoveIdentity(userAArn, false)
			Expect(err).To(HaveOccurred())
		})
		It("should remove all if specified", func() {
			err := acm.AddIdentity(mustIdentity(userAArn, userAUsername, userAGroups))
			Expect(err).NotTo(HaveOccurred())
			err = acm.AddIdentity(mustIdentity(userAArn, userAUsername, userAGroups))
			Expect(err).NotTo(HaveOccurred())
			Expect(client.updated.Data["mapUsers"]).To(Not(MatchYAML("[]")))

			err = acm.RemoveIdentity(userAArn, true)
			Expect(err).NotTo(HaveOccurred())
			Expect(client.updated.Data["mapUsers"]).To(MatchYAML("[]"))
		})
	})
	Describe("AddIdentity()", func() {
		existing := &corev1.ConfigMap{
			ObjectMeta: ObjectMeta(),
			Data:       map[string]string{},
		}
		existing.UID = "123456"
		client := &mockClient{}
		acm := New(client, existing)

		addAndSave := func(canonicalArn, user string, groups []string) *corev1.ConfigMap {
			client.reset()

			arn, err := iam.Parse(canonicalArn)
			Expect(err).NotTo(HaveOccurred())

			err = acm.AddIdentity(mustIdentity(arn, user, groups))
			Expect(err).NotTo(HaveOccurred())

			err = acm.Save()
			Expect(err).NotTo(HaveOccurred())
			Expect(client.created).To(BeNil())
			Expect(client.updated).NotTo(BeNil())

			return client.updated
		}

		It("should add a role and a user", func() {
			cm := addAndSave(roleA, RoleNodeGroupUsername, RoleNodeGroupGroups)
			cm = addAndSave(userA, userAUsername, userAGroups)
			Expect(cm.Data["mapRoles"]).To(MatchYAML(expectedRoleA))
			Expect(cm.Data["mapUsers"]).To(MatchYAML(expectedUserA))
		})
		It("should append a second role and user", func() {
			cm := addAndSave(roleB, RoleNodeGroupUsername, []string{groupB})
			cm = addAndSave(userB, userBUsername, userBGroups)
			Expect(cm.Data["mapRoles"]).To(MatchYAML(expectedRoleA + expectedRoleB))
			Expect(cm.Data["mapUsers"]).To(MatchYAML(expectedUserA + expectedUserB))
		})
		It("should append a duplicate role", func() {
			cm := addAndSave(roleA, RoleNodeGroupUsername, RoleNodeGroupGroups)
			expectedRoles := expectedRoleA + expectedRoleB + expectedRoleA

			cm = addAndSave(userA, userAUsername, userAGroups)
			expectedUsers := expectedUserA + expectedUserB + expectedUserA

			Expect(cm.Data["mapRoles"]).To(MatchYAML(expectedRoles))
			Expect(cm.Data["mapUsers"]).To(MatchYAML(expectedUsers))
		})
	})
	Describe("RemoveIdentity()", func() {
		existing := &corev1.ConfigMap{
			ObjectMeta: ObjectMeta(),
			Data: map[string]string{
				"mapRoles": expectedRoleA + expectedRoleA + expectedRoleB,
				"mapUsers": expectedUserA + expectedUserA + expectedUserB,
			},
		}
		existing.UID = "123456"
		client := &mockClient{}
		acm := New(client, existing)

		removeAndSave := func(canonicalArn string) *corev1.ConfigMap {
			client.reset()

			arn, err := iam.Parse(canonicalArn)
			Expect(err).NotTo(HaveOccurred())

			err = acm.RemoveIdentity(arn, false)
			Expect(err).NotTo(HaveOccurred())

			err = acm.Save()
			Expect(err).NotTo(HaveOccurred())

			return client.updated
		}

		roleAArn, err := iam.Parse(roleA)
		if err != nil {
			panic(err)
		}

		userAArn, err := iam.Parse(userA)
		if err != nil {
			panic(err)
		}

		It("should remove role and user", func() {
			cm := removeAndSave(roleB)
			cm = removeAndSave(userB)
			Expect(cm.Data["mapRoles"]).To(MatchYAML(expectedRoleA + expectedRoleA))
			Expect(cm.Data["mapUsers"]).To(MatchYAML(expectedUserA + expectedUserA))
		})
		It("should remove one role and one user for duplicates", func() {
			cm := removeAndSave(roleA)
			cm = removeAndSave(userA)
			Expect(cm.Data["mapRoles"]).To(MatchYAML(expectedRoleA))
			Expect(cm.Data["mapUsers"]).To(MatchYAML(expectedUserA))
		})
		It("should remove last role and last user", func() {
			cm := removeAndSave(roleA)
			cm = removeAndSave(userA)
			Expect(cm.Data["mapRoles"]).To(MatchYAML("[]"))
			Expect(cm.Data["mapUsers"]).To(MatchYAML("[]"))
		})
		It("should fail if role or user not found", func() {
			err := acm.RemoveIdentity(roleAArn, false)
			Expect(err).To(HaveOccurred())
			err = acm.RemoveIdentity(userAArn, false)
			Expect(err).To(HaveOccurred())
		})
		It("should remove all if specified", func() {
			err := acm.AddIdentity(mustIdentity(roleAArn, RoleNodeGroupUsername, RoleNodeGroupGroups))
			Expect(err).NotTo(HaveOccurred())
			err = acm.AddIdentity(mustIdentity(roleAArn, RoleNodeGroupUsername, RoleNodeGroupGroups))
			Expect(err).NotTo(HaveOccurred())
			Expect(client.updated.Data["mapRoles"]).To(Not(MatchYAML("[]")))

			err = acm.AddIdentity(mustIdentity(userAArn, userAUsername, userAGroups))
			Expect(err).NotTo(HaveOccurred())
			err = acm.AddIdentity(mustIdentity(userAArn, userAUsername, userAGroups))
			Expect(err).NotTo(HaveOccurred())
			Expect(client.updated.Data["mapUsers"]).To(Not(MatchYAML("[]")))

			err = acm.RemoveIdentity(roleAArn, true)
			Expect(err).NotTo(HaveOccurred())
			Expect(client.updated.Data["mapRoles"]).To(MatchYAML("[]"))
			Expect(client.updated.Data["mapUsers"]).To(Not(MatchYAML("[]")))

			err = acm.RemoveIdentity(userAArn, true)
			Expect(err).NotTo(HaveOccurred())
			Expect(client.updated.Data["mapRoles"]).To(MatchYAML("[]"))
			Expect(client.updated.Data["mapUsers"]).To(MatchYAML("[]"))
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
