package nodebootstrap

import (
	"encoding/base64"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	toml "github.com/pelletier/go-toml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("Bottlerocket", func() {
	var (
		clusterConfig *api.ClusterConfig
		ng            *api.NodeGroup
	)

	userdataTOML := func(userdata string) (*toml.Tree, error) {
		data, err := base64.StdEncoding.DecodeString(userdata)
		if err != nil {
			return nil, err
		}
		return toml.LoadBytes(data)
	}

	BeforeEach(func() {
		clusterConfig = api.NewClusterConfig()
		clusterConfig.Status = &api.ClusterStatus{
			Endpoint:                 "unit-test.example.com",
			CertificateAuthorityData: []byte(`CertificateAuthorityData`),
		}
		clusterConfig.Metadata = &api.ClusterMeta{
			Name: "unit-test",
		}
		ng = &api.NodeGroup{}
	})

	Describe("with no user settings", func() {
		It("produces TOML userdata", func() {
			userdata, err := NewUserDataForBottlerocket(clusterConfig, ng)
			Expect(err).ToNot(HaveOccurred())
			Expect(userdata).ToNot(Equal(""))

			settings, parseErr := userdataTOML(userdata)
			Expect(parseErr).ToNot(HaveOccurred())
			Expect(settings.Has("settings.kubernetes.cluster-name")).To(BeTrue())
		})

		It("leaves settings.host-containers.admin.enabled commented", func() {
			userdata, err := NewUserDataForBottlerocket(clusterConfig, ng)
			Expect(err).ToNot(HaveOccurred())
			Expect(userdata).ToNot(Equal(""))

			settings, parseErr := userdataTOML(userdata)
			Expect(parseErr).ToNot(HaveOccurred())
			Expect(settings.Has("settings.host-containers.admin.enabled")).To(BeFalse())
			tomlStr, err := settings.ToTomlString()
			Expect(err).ToNot(HaveOccurred())
			// Generated TOML should contain a section (with the enabled
			// key=value) for the unset, commented out setting.
			const settingSection = "[settings.host-containers.admin]"
			Expect(tomlStr).To(ContainSubstring(settingSection))
		})
	})

	Describe("with user settings", func() {
		BeforeEach(func() {
			ng.Bottlerocket = &api.NodeGroupBottlerocket{
				Settings: &api.InlineDocument{
					"host-containers": map[string]interface{}{
						"example": map[string]bool{
							"enabled": true,
						},
					},
				},
			}
		})

		It("produces TOML userdata with provided settings", func() {
			userdata, err := NewUserDataForBottlerocket(clusterConfig, ng)
			Expect(err).ToNot(HaveOccurred())
			Expect(userdata).ToNot(Equal(""))

			settings, parseErr := userdataTOML(userdata)
			Expect(parseErr).ToNot(HaveOccurred())
			Expect(settings.Has("settings.kubernetes.cluster-name")).To(BeTrue())
			Expect(settings.Has("settings.host-containers.example.enabled")).To(BeTrue())
		})

		It("enables admin container", func() {
			ng.Bottlerocket.EnableAdminContainer = api.Enabled()
			userdata, err := NewUserDataForBottlerocket(clusterConfig, ng)
			Expect(err).ToNot(HaveOccurred())
			Expect(userdata).ToNot(Equal(""))

			settings, parseErr := userdataTOML(userdata)
			Expect(parseErr).ToNot(HaveOccurred())
			Expect(settings.Has("settings.host-containers.admin.enabled")).To(BeTrue())
			val, ok := settings.Get("settings.host-containers.admin.enabled").(bool)
			Expect(ok).To(BeTrue())
			Expect(val).To(BeTrue())
		})

		It("retains user specified values", func() {
			// Enable in Bottlerocket's top level config.
			ng.Bottlerocket.EnableAdminContainer = api.Enabled()
			// But set conflicting type and value to
			// otherwise managed key.
			providedSettings := map[string]interface{}(*ng.Bottlerocket.Settings)
			providedSettings["host-containers"].(map[string]interface{})["admin"] = map[string]string{"enabled": "user-val"}
			userdata, err := NewUserDataForBottlerocket(clusterConfig, ng)
			Expect(err).ToNot(HaveOccurred())
			Expect(userdata).ToNot(Equal(""))

			settings, parseErr := userdataTOML(userdata)
			Expect(parseErr).ToNot(HaveOccurred())
			// Check that the value specified in config is
			// set, not the higher level toggle.
			Expect(settings.Has("settings.host-containers.admin.enabled")).To(BeTrue())
			val, ok := settings.Get("settings.host-containers.admin.enabled").(string)
			Expect(ok).To(BeTrue())
			Expect(val).To(Equal("user-val"))
		})

		It("produces TOML userdata with quoted keys", func() {
			keyName := "dotted.key.name"
			providedSettings := map[string]interface{}(*ng.Bottlerocket.Settings)
			providedSettings[keyName] = "value"
			keyPath := []string{"settings", keyName}
			splitKeyPath := append([]string{"settings"}, strings.Split(keyName, ".")...)

			userdata, err := NewUserDataForBottlerocket(clusterConfig, ng)
			Expect(err).ToNot(HaveOccurred())
			Expect(userdata).ToNot(Equal(""))

			tree, parseErr := userdataTOML(userdata)
			Expect(parseErr).ToNot(HaveOccurred())
			// Verify the keys made it to where they were
			// supposed to and that nothing happened to
			// make them appear as dotted/split key names.
			Expect(tree.HasPath(splitKeyPath)).To(BeFalse())
			Expect(tree.HasPath(keyPath)).To(BeTrue())
		})
	})

})

func TestProtectTOMLKeys(t *testing.T) {
	testcases := []struct {
		name     string
		data     map[string]interface{}
		paths    [][]string
		notPaths [][]string
	}{
		{
			// Essential traversal.
			name: "shallow",
			data: map[string]interface{}{
				"key": "val",
			},
			paths: [][]string{
				{"key"},
			},
		},
		{
			// Essential traversal.
			name: "empty",
			data: map[string]interface{}{},
			paths: [][]string{
				{}, // Roughly: "the current Tree node", so checking this "key" should return true.
			},
		},

		{
			// Nested tree traversal.
			name: "nested",
			data: map[string]interface{}{
				"nested": map[string]interface{}{
					"nestedKey": "val",
				},
			},
			paths: [][]string{
				{"nested", "nestedKey"},
			},
		},
		{
			// Traversal and transformation with targeted dotted key naming.
			name: "dotted-shallow",
			data: map[string]interface{}{
				"dotted.key": "val",
			},
			paths: [][]string{
				{`"dotted.key"`},
			},
			notPaths: [][]string{
				{"dotted", "key"},
				{"dotted.key"},
			},
		},
		{
			// Traversal and transformation with targeted dotted key naming
			// within nested trees.
			name: "dotted-nested",
			data: map[string]interface{}{
				"nested": map[string]interface{}{
					"dotted.key": "val",
				},
				"dotted.nested": map[string]interface{}{
					"key": "val",
				},
			},
			paths: [][]string{
				{"nested", `"dotted.key"`},

				{`"dotted.nested"`, "key"},
			},
			notPaths: [][]string{
				{"nested", "dotted.key"},
				{"nested", "dotted", "key"},

				{"dotted.nested", "key"},
				{"dotted", "nested", "key"},
			},
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.name, func(t *testing.T) {
			tree, err := toml.TreeFromMap(testcase.data)
			require.NoError(t, err)

			protectTOMLKeys([]string{}, tree)

			for _, keyPath := range testcase.paths {
				assert.True(t, tree.HasPath(keyPath), "should have key at path %q", keyPath)
			}
			for _, keyPath := range testcase.notPaths {
				assert.False(t, tree.HasPath(keyPath), "should not have key at path %q", keyPath)
			}
		})
	}
}
