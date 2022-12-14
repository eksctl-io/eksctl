package nodebootstrap_test

import (
	"encoding/base64"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

var _ = Describe("Managed Windows UserData", func() {

	type windowsEntry struct {
		updateNodeGroup func(*api.ManagedNodeGroup)

		expectedUserData string
	}

	DescribeTable("Managed Windows bootstrap", func(e windowsEntry) {
		ng := api.NewManagedNodeGroup()
		ng.AMIFamily = api.NodeImageFamilyWindowsServer2019CoreContainer
		if e.updateNodeGroup != nil {
			e.updateNodeGroup(ng)
		}

		bootstrapper := nodebootstrap.ManagedWindows{NodeGroup: ng}
		userData, err := bootstrapper.UserData()
		Expect(err).NotTo(HaveOccurred())
		actual, err := base64.StdEncoding.DecodeString(userData)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(actual)).To(Equal(strings.TrimSpace(e.expectedUserData)))
	},
		Entry("should have no bootstrap args in userdata", windowsEntry{
			expectedUserData: "",
		}),

		Entry("should not have labels here", windowsEntry{
			updateNodeGroup: func(ng *api.ManagedNodeGroup) {
				ng.Labels = map[string]string{
					"foo": "bar",
				}
			},

			expectedUserData: "",
		}),

		Entry("should not have taints here", windowsEntry{
			updateNodeGroup: func(ng *api.ManagedNodeGroup) {
				ng.Taints = []api.NodeGroupTaint{
					{
						Key:    "foo",
						Value:  "bar",
						Effect: "NoSchedule",
					},
				}
			},

			expectedUserData: "",
		}),

		Entry("should only have a preBootstrapCommand", windowsEntry{
			updateNodeGroup: func(ng *api.ManagedNodeGroup) {
				ng.PreBootstrapCommands = []string{
					"wget -UseBasicParsing -O amazon-cloudwatch-agent.msi https://s3.amazonaws.com/amazoncloudwatch-agent/windows/amd64/latest/amazon-cloudwatch-agent.msi",
				}
			},

			expectedUserData: `
<powershell>
wget -UseBasicParsing -O amazon-cloudwatch-agent.msi https://s3.amazonaws.com/amazoncloudwatch-agent/windows/amd64/latest/amazon-cloudwatch-agent.msi
</powershell>
`,
		}),

		Entry("should have several preBootstrapCommands", windowsEntry{
			updateNodeGroup: func(ng *api.ManagedNodeGroup) {
				ng.PreBootstrapCommands = []string{
					"wget -UseBasicParsing -O amazon-cloudwatch-agent.msi https://s3.amazonaws.com/amazoncloudwatch-agent/windows/amd64/latest/amazon-cloudwatch-agent.msi",
					"start /wait msiexec.exe /qb /i \"amazon-cloudwatch-agent.msi\"",
				}

			},

			expectedUserData: `
<powershell>
wget -UseBasicParsing -O amazon-cloudwatch-agent.msi https://s3.amazonaws.com/amazoncloudwatch-agent/windows/amd64/latest/amazon-cloudwatch-agent.msi
start /wait msiexec.exe /qb /i "amazon-cloudwatch-agent.msi"
</powershell>
`,
		}),
	)

})
