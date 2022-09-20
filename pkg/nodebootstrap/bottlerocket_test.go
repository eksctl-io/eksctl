package nodebootstrap_test

import (
	"encoding/base64"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	toml "github.com/pelletier/go-toml"
	"github.com/stretchr/testify/require"
	"github.com/tj/assert"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

var _ = Describe("Bottlerocket", func() {
	var (
		clusterConfig *api.ClusterConfig
		ng            *api.NodeGroup
	)

	BeforeEach(func() {
		clusterConfig = api.NewClusterConfig()
		clusterConfig.Status = &api.ClusterStatus{
			Endpoint:                 "unit-test.example.com",
			CertificateAuthorityData: []byte(`CertificateAuthorityData`),
		}
		clusterConfig.Metadata = &api.ClusterMeta{
			Name: "unit-test",
		}

		ng = api.NewNodeGroup()
		ng.AMIFamily = "Bottlerocket"
		api.SetNodeGroupDefaults(ng, clusterConfig.Metadata, false)
	})

	Describe("with no user settings", func() {
		It("produces standard TOML userdata", func() {
			bootstrapper := newBootstrapper(clusterConfig, ng)
			userdata, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())
			Expect(userdata).NotTo(Equal(""))

			settings, parseErr := userdataTOML(userdata)
			Expect(parseErr).NotTo(HaveOccurred())
			Expect(settings.Has("settings.kubernetes.cluster-name")).To(BeTrue())
		})

		It("leaves settings.host-containers.admin.enabled commented", func() {
			bootstrapper := newBootstrapper(clusterConfig, ng)
			userdata, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())
			Expect(userdata).NotTo(Equal(""))

			settings, parseErr := userdataTOML(userdata)
			Expect(parseErr).NotTo(HaveOccurred())
			Expect(settings.Has("settings.host-containers.admin.enabled")).To(BeFalse())
			tomlStr, err := settings.ToTomlString()
			Expect(err).NotTo(HaveOccurred())
			// Generated TOML should contain a section (with the enabled
			// key=value) for the unset, commented out setting.
			Expect(tomlStr).To(ContainSubstring("[settings.host-containers.admin]"))
		})
	})

	Describe("with user settings", func() {
		BeforeEach(func() {
			ng.Bottlerocket = &api.NodeGroupBottlerocket{
				Settings: &api.InlineDocument{
					"host-containers": map[string]interface{}{},
				},
			}
		})

		It("produces TOML userdata with quoted keys", func() {
			keyName := "dotted.key.name"
			providedSettings := map[string]interface{}(*ng.Bottlerocket.Settings)
			providedSettings[keyName] = "value"
			keyPath := []string{"settings", keyName}
			splitKeyPath := append([]string{"settings"}, strings.Split(keyName, ".")...)

			bootstrapper := newBootstrapper(clusterConfig, ng)
			userdata, err := bootstrapper.UserData()
			Expect(err).NotTo(HaveOccurred())
			Expect(userdata).NotTo(Equal(""))

			tree, parseErr := userdataTOML(userdata)
			Expect(parseErr).NotTo(HaveOccurred())
			// Verify the keys made it to where they were
			// supposed to and that nothing happened to
			// make them appear as dotted/split key names.
			Expect(tree.HasPath(splitKeyPath)).To(BeFalse())
			Expect(tree.HasPath(keyPath)).To(BeTrue())
		})

		When("host containers are enabled", func() {
			BeforeEach(func() {
				ng.Bottlerocket.Settings = &api.InlineDocument{
					"host-containers": map[string]interface{}{
						"example": map[string]bool{
							"enabled": true,
						},
					},
				}
			})

			It("sets it on the userdata", func() {
				bootstrapper := newBootstrapper(clusterConfig, ng)
				userdata, err := bootstrapper.UserData()
				Expect(err).NotTo(HaveOccurred())
				Expect(userdata).NotTo(Equal(""))

				settings, parseErr := userdataTOML(userdata)
				Expect(parseErr).NotTo(HaveOccurred())
				Expect(settings.Has("settings.host-containers.example.enabled")).To(BeTrue())
			})
		})

		When("admin container is enabled", func() {
			BeforeEach(func() {
				ng.Bottlerocket.EnableAdminContainer = api.Enabled()
			})

			It("enables admin container on the userdata", func() {
				bootstrapper := newBootstrapper(clusterConfig, ng)
				userdata, err := bootstrapper.UserData()
				Expect(err).NotTo(HaveOccurred())
				Expect(userdata).NotTo(Equal(""))

				settings, parseErr := userdataTOML(userdata)
				Expect(parseErr).NotTo(HaveOccurred())
				Expect(settings.Has("settings.host-containers.admin.enabled")).To(BeTrue())
				val, ok := settings.Get("settings.host-containers.admin.enabled").(bool)
				Expect(ok).To(BeTrue())
				Expect(val).To(BeTrue())
			})

			It("retains user specified values", func() {
				// Set conflicting type and value to
				// otherwise managed key.
				providedSettings := map[string]interface{}(*ng.Bottlerocket.Settings)
				providedSettings["host-containers"].(map[string]interface{})["admin"] = map[string]string{"enabled": "user-val"}
				bootstrapper := newBootstrapper(clusterConfig, ng)
				userdata, err := bootstrapper.UserData()
				Expect(err).NotTo(HaveOccurred())
				Expect(userdata).NotTo(Equal(""))

				settings, parseErr := userdataTOML(userdata)
				Expect(parseErr).NotTo(HaveOccurred())
				// Check that the value specified in config is
				// set, not the higher level toggle.
				Expect(settings.Has("settings.host-containers.admin.enabled")).To(BeTrue())
				val, ok := settings.Get("settings.host-containers.admin.enabled").(string)
				Expect(ok).To(BeTrue())
				Expect(val).To(Equal("user-val"))
			})
		})
	})

	Describe("with NodeGroup settings", func() {
		var (
			maxPodsPath      = strings.Split("settings.kubernetes.max-pods", ".")
			labelsPath       = strings.Split("settings.kubernetes.node-labels", ".")
			taintsPath       = strings.Split("settings.kubernetes.node-taints", ".")
			clusterDNSIPPath = strings.Split("settings.kubernetes.cluster-dns-ip", ".")
		)

		When("labels are set on the node", func() {
			BeforeEach(func() {
				ng.Labels = map[string]string{"foo": "bar"}
			})

			It("adds the labels to the userdata", func() {
				bootstrapper := newBootstrapper(clusterConfig, ng)
				userdata, err := bootstrapper.UserData()
				Expect(err).NotTo(HaveOccurred())
				Expect(userdata).NotTo(Equal(""))

				tree, parseErr := userdataTOML(userdata)
				Expect(parseErr).NotTo(HaveOccurred())

				Expect(tree.HasPath(append(labelsPath, "foo"))).To(BeTrue())
				Expect(tree.GetPath(append(labelsPath, "foo"))).To(Equal("bar"))
			})
		})

		When("taints are set on the node", func() {
			BeforeEach(func() {
				ng.Taints = []api.NodeGroupTaint{
					{
						Key:    "foo",
						Value:  "bar",
						Effect: "NoExecute",
					},
				}
			})

			It("adds the taints to the userdata", func() {
				bootstrapper := newBootstrapper(clusterConfig, ng)
				userdata, err := bootstrapper.UserData()
				Expect(err).NotTo(HaveOccurred())
				Expect(userdata).NotTo(Equal(""))

				tree, parseErr := userdataTOML(userdata)
				Expect(parseErr).NotTo(HaveOccurred())

				Expect(tree.HasPath(append(taintsPath, "foo"))).To(BeTrue())
				Expect(tree.GetPath(append(taintsPath, "foo"))).To(Equal("bar:NoExecute"))
			})
		})

		When("clusterDNS is set", func() {
			It("adds clusterDNS to the userdata", func() {
				ng.ClusterDNS = "192.2.0.53"

				bootstrapper := newBootstrapper(clusterConfig, ng)
				userdata, err := bootstrapper.UserData()
				Expect(err).NotTo(HaveOccurred())
				Expect(userdata).NotTo(Equal(""))

				tree, parseErr := userdataTOML(userdata)
				Expect(parseErr).NotTo(HaveOccurred())

				Expect(tree.HasPath(clusterDNSIPPath)).To(BeTrue())
				Expect(tree.GetPath(clusterDNSIPPath)).To(Equal(ng.ClusterDNS))
			})
		})

		When("maxPods", func() {
			It("adds MaxPodsPerNode to userdata when set", func() {
				ng.MaxPodsPerNode = 32

				bootstrapper := newBootstrapper(clusterConfig, ng)
				userdata, err := bootstrapper.UserData()
				Expect(err).NotTo(HaveOccurred())
				Expect(userdata).NotTo(Equal(""))

				tree, parseErr := userdataTOML(userdata)
				Expect(parseErr).NotTo(HaveOccurred())

				Expect(tree.HasPath(maxPodsPath)).To(BeTrue())
				Expect(tree.GetPath(maxPodsPath)).To(Equal(int64(ng.MaxPodsPerNode)))
			})

			It("does not add MaxPodsPerNode when not set", func() {
				bootstrapper := newBootstrapper(clusterConfig, ng)
				userdata, err := bootstrapper.UserData()
				Expect(err).NotTo(HaveOccurred())
				Expect(userdata).NotTo(Equal(""))

				tree, parseErr := userdataTOML(userdata)
				Expect(parseErr).NotTo(HaveOccurred())

				Expect(tree.HasPath(maxPodsPath)).To(BeFalse())
			})
		})
	})
})

func userdataTOML(userdata string) (*toml.Tree, error) {
	data, err := base64.StdEncoding.DecodeString(userdata)
	if err != nil {
		return nil, err
	}
	return toml.LoadBytes(data)
}

// TODO sort this out. ginkgo or bust
// and make that func private
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

			nodebootstrap.ProtectTOMLKeys([]string{}, tree)

			for _, keyPath := range testcase.paths {
				assert.True(t, tree.HasPath(keyPath), "should have key at path %q", keyPath)
			}
			for _, keyPath := range testcase.notPaths {
				assert.False(t, tree.HasPath(keyPath), "should not have key at path %q", keyPath)
			}
		})
	}
}
