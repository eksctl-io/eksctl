package nodebootstrap_test

import (
	"encoding/base64"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

var _ = Describe("Managed Bottlerocket", func() {

	type bottlerocketEntry struct {
		setFields func(group *api.ManagedNodeGroup)

		expectedErr      string
		expectedUserData string
	}

	DescribeTable("User data", func(e bottlerocketEntry) {
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Name = "managed-bottlerocket"
		clusterConfig.Status = &api.ClusterStatus{
			Endpoint:                 "bottlerocket.cluster.com",
			CertificateAuthorityData: []byte("CertificateAuthorityData"),
		}

		ng := api.NewManagedNodeGroup()
		ng.AMIFamily = api.NodeImageFamilyBottlerocket
		api.SetManagedNodeGroupDefaults(ng, clusterConfig.Metadata, false)
		if e.setFields != nil {
			e.setFields(ng)
		}

		bootstrapper, err := nodebootstrap.NewManagedBootstrapper(clusterConfig, ng)
		Expect(err).NotTo(HaveOccurred())
		userData, err := bootstrapper.UserData()
		if e.expectedErr != "" {
			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError(ContainSubstring(e.expectedErr)))
			return
		}

		Expect(err).NotTo(HaveOccurred())
		actual, err := base64.StdEncoding.DecodeString(userData)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(actual)).To(Equal(e.expectedUserData))
	},
		Entry("no settings", bottlerocketEntry{
			expectedUserData: `
[settings]

  [settings.kubernetes]
`,
		}),
		Entry("maxPods set", bottlerocketEntry{
			setFields: func(ng *api.ManagedNodeGroup) {
				ng.MaxPodsPerNode = 44
			},
			expectedUserData: `
[settings]

  [settings.kubernetes]
    max-pods = 44
`,
		}),

		Entry("enableAdminContainer set", bottlerocketEntry{
			setFields: func(ng *api.ManagedNodeGroup) {
				ng.Bottlerocket.EnableAdminContainer = api.Enabled()
			},
			expectedUserData: `
[settings]

  [settings.host-containers]

    [settings.host-containers.admin]
      enabled = true

  [settings.kubernetes]
`,
		}),

		Entry("host containers enabled", bottlerocketEntry{
			setFields: func(ng *api.ManagedNodeGroup) {
				ng.Bottlerocket.Settings = &api.InlineDocument{
					"host-containers": api.InlineDocument{
						"example": map[string]bool{
							"enabled": true,
						},
					},
				}
			},
			expectedUserData: `
[settings]

  [settings.host-containers]

    [settings.host-containers.example]
      enabled = true

  [settings.kubernetes]
`,
		}),

		Entry("retain user-specified admin container setting", bottlerocketEntry{
			setFields: func(ng *api.ManagedNodeGroup) {
				ng.Bottlerocket.Settings = &api.InlineDocument{
					"host-containers": map[string]interface{}{
						"admin": map[string]interface{}{
							"enabled": true,
						},
					},
				}
			},
			expectedUserData: `
[settings]

  [settings.host-containers]

    [settings.host-containers.admin]
      enabled = true

  [settings.kubernetes]
`,
		}),

		Entry("labels and taints set", bottlerocketEntry{
			setFields: func(ng *api.ManagedNodeGroup) {
				ng.Labels = map[string]string{
					"key": "value",
				}
				ng.Taints = []api.NodeGroupTaint{
					{
						Key:    "foo",
						Value:  "bar",
						Effect: "NoExecute",
					},
				}
			},

			expectedUserData: `
[settings]

  [settings.kubernetes]
`,
		}),

		Entry("preserve dotted keys", bottlerocketEntry{
			setFields: func(ng *api.ManagedNodeGroup) {
				ng.Bottlerocket.Settings = &api.InlineDocument{
					"a.b.c": "value",
				}
			},

			expectedUserData: `
[settings]
  "a.b.c" = "value"

  [settings.kubernetes]
`,
		}),

		Entry("cluster bootstrap settings set", bottlerocketEntry{
			setFields: func(ng *api.ManagedNodeGroup) {
				ng.Bottlerocket.Settings = &api.InlineDocument{
					"kubernetes": map[string]interface{}{
						"cluster-certificate": "test",
						"api-server":          "https://test.com",
						"cluster-name":        "test",
						"cluster-dns-ip":      "192.2.0.53",
					},
				}
			},

			expectedErr: "EKS automatically injects cluster bootstrapping fields into user-data",
		}),

		Entry("labels and taints in Bottlerocket settings", bottlerocketEntry{
			setFields: func(ng *api.ManagedNodeGroup) {
				ng.Bottlerocket.Settings = &api.InlineDocument{
					"kubernetes": map[string]interface{}{
						"node-labels": map[string]string{
							"key": "value",
						},
						"node-taints": map[string]string{
							"foo": "bar:NoExecute",
						},
					},
				}
			},

			expectedErr: "cannot set settings.kubernetes.node-labels; labels and taints should be set on the managedNodeGroup object",
		}),

		Entry("conflicting settings", bottlerocketEntry{
			setFields: func(ng *api.ManagedNodeGroup) {
				ng.Bottlerocket.EnableAdminContainer = api.Enabled()
				ng.Bottlerocket.Settings = &api.InlineDocument{
					"host-containers": map[string]interface{}{
						"admin": map[string]interface{}{
							"enabled": false,
						},
					},
				}
			},

			expectedErr: "cannot set both bottlerocket.enableAdminContainer and settings.host-containers.admin.enabled",
		}),
	)
})
