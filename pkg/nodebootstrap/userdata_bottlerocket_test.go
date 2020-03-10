package nodebootstrap

import (
	"encoding/base64"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/pelletier/go-toml"
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
		ng = &api.NodeGroup{
			// SetNodeGroupDefaults ensures this is non-nil for Bottlerocket nodegroups
			Bottlerocket: &api.NodeGroupBottlerocket{},
		}
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

		Describe("with NodeGroup settings", func() {
			var (
				maxPodsPath      = strings.Split("settings.kubernetes.max-pods", ".")
				labelsPath       = strings.Split("settings.kubernetes.node-labels", ".")
				taintsPath       = strings.Split("settings.kubernetes.node-taints", ".")
				clusterDNSIPPath = strings.Split("settings.kubernetes.cluster-dns-ip", ".")
			)
			BeforeEach(func() {
				ng.Labels = map[string]string{}
				ng.Taints = map[string]string{}

				api.SetNodeGroupDefaults(ng, clusterConfig.Metadata)
			})

			It("removes overlapping config", func() {
				checkKey := "expected-missing.example.com"
				checkKeyVal := "should-be-none"

				doc := api.InlineDocument(map[string]interface{}{
					"kubernetes": map[string]interface{}{
						// Neither of these should be present in the generated
						// TOML.
						"node-labels": map[string]string{
							checkKey: checkKeyVal,
						},
						"node-taints": map[string]string{
							checkKey: checkKeyVal,
						},
					},
				})

				ng.Bottlerocket.Settings = &doc

				userdata, err := NewUserDataForBottlerocket(clusterConfig, ng)
				Expect(err).ToNot(HaveOccurred())
				Expect(userdata).ToNot(Equal(""))

				tree, parseErr := userdataTOML(userdata)
				Expect(parseErr).ToNot(HaveOccurred())

				Expect(tree.HasPath(append(labelsPath, checkKey))).To(BeFalse())
				Expect(tree.HasPath(append(taintsPath, checkKey))).To(BeFalse())
			})

			It("uses MaxPodsPerNode", func() {
				ng.MaxPodsPerNode = 32

				userdata, err := NewUserDataForBottlerocket(clusterConfig, ng)
				Expect(err).ToNot(HaveOccurred())
				Expect(userdata).ToNot(Equal(""))

				tree, parseErr := userdataTOML(userdata)
				Expect(parseErr).ToNot(HaveOccurred())

				Expect(tree.HasPath(maxPodsPath)).To(BeTrue())
				Expect(tree.GetPath(maxPodsPath)).To(Equal(int64(ng.MaxPodsPerNode)))
			})

			It("handles unset MaxPodsPerNode", func() {
				userdata, err := NewUserDataForBottlerocket(clusterConfig, ng)
				Expect(err).ToNot(HaveOccurred())
				Expect(userdata).ToNot(Equal(""))

				tree, parseErr := userdataTOML(userdata)
				Expect(parseErr).ToNot(HaveOccurred())

				Expect(tree.HasPath(maxPodsPath)).To(BeFalse())
			})

			It("uses ClusterDNS", func() {
				ng.ClusterDNS = "192.2.0.53"

				userdata, err := NewUserDataForBottlerocket(clusterConfig, ng)
				Expect(err).ToNot(HaveOccurred())
				Expect(userdata).ToNot(Equal(""))

				tree, parseErr := userdataTOML(userdata)
				Expect(parseErr).ToNot(HaveOccurred())

				Expect(tree.HasPath(clusterDNSIPPath)).To(BeTrue())
				Expect(tree.GetPath(clusterDNSIPPath)).To(Equal(ng.ClusterDNS))
			})

			It("uses Taints", func() {
				taintName := "mytaint.example.com"
				taintVal := "00.00001"
				ng.Taints[taintName] = taintVal

				userdata, err := NewUserDataForBottlerocket(clusterConfig, ng)
				Expect(err).ToNot(HaveOccurred())
				Expect(userdata).ToNot(Equal(""))

				tree, parseErr := userdataTOML(userdata)
				Expect(parseErr).ToNot(HaveOccurred())

				for key, val := range ng.Taints {
					Expect(tree.HasPath(append(taintsPath, key))).To(BeTrue())
					Expect(tree.GetPath(append(taintsPath, key))).To(Equal(val))
				}
			})

			It("uses labels", func() {
				labelName := "mylabel.example.com"
				labelVal := "99.99999"
				ng.Labels[labelName] = labelVal

				userdata, err := NewUserDataForBottlerocket(clusterConfig, ng)
				Expect(err).ToNot(HaveOccurred())
				Expect(userdata).ToNot(Equal(""))

				tree, parseErr := userdataTOML(userdata)
				Expect(parseErr).ToNot(HaveOccurred())

				for key, val := range ng.Labels {
					Expect(tree.HasPath(append(labelsPath, key))).To(BeTrue())
					Expect(tree.GetPath(append(labelsPath, key))).To(Equal(val))
				}
			})
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
