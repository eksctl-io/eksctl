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
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Name = "windohs"
		clusterConfig.Status = &api.ClusterStatus{
			Endpoint:                 "https://test.com",
			CertificateAuthorityData: []byte("test"),
		}

		ng := api.NewManagedNodeGroup()
		ng.AMIFamily = api.NodeImageFamilyWindowsServer2019CoreContainer
		if e.updateNodeGroup != nil {
			e.updateNodeGroup(ng)
		}

		bootstrapper := nodebootstrap.NewWindowsBootstrapper(clusterConfig, ng, "")
		userData, err := bootstrapper.UserData()
		Expect(err).NotTo(HaveOccurred())
		actual, err := base64.StdEncoding.DecodeString(userData)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(actual)).To(Equal(strings.TrimSpace(e.expectedUserData)))
	},
		Entry("mandatory bootstrap args in userdata", windowsEntry{

			expectedUserData: `
<powershell>
[string]$EKSBootstrapScriptFile = "$env:ProgramFiles\Amazon\EKS\Start-EKSBootstrap.ps1"
& $EKSBootstrapScriptFile -EKSClusterName "windohs" -APIServerEndpoint "https://test.com" -Base64ClusterCA "dGVzdA==" -KubeletExtraArgs "--node-labels= --register-with-taints=" 3>&1 4>&1 5>&1 6>&1
</powershell>
`,
		}),

		Entry("with labels", windowsEntry{
			updateNodeGroup: func(ng *api.ManagedNodeGroup) {
				ng.Labels = map[string]string{
					"foo": "bar",
				}
			},

			expectedUserData: `
<powershell>
[string]$EKSBootstrapScriptFile = "$env:ProgramFiles\Amazon\EKS\Start-EKSBootstrap.ps1"
& $EKSBootstrapScriptFile -EKSClusterName "windohs" -APIServerEndpoint "https://test.com" -Base64ClusterCA "dGVzdA==" -KubeletExtraArgs "--node-labels=foo=bar --register-with-taints=" 3>&1 4>&1 5>&1 6>&1
</powershell>
`,
		}),

		Entry("with taints", windowsEntry{
			updateNodeGroup: func(ng *api.ManagedNodeGroup) {
				ng.Taints = []api.NodeGroupTaint{
					{
						Key:    "foo",
						Value:  "bar",
						Effect: "NoSchedule",
					},
				}
			},

			expectedUserData: `
<powershell>
[string]$EKSBootstrapScriptFile = "$env:ProgramFiles\Amazon\EKS\Start-EKSBootstrap.ps1"
& $EKSBootstrapScriptFile -EKSClusterName "windohs" -APIServerEndpoint "https://test.com" -Base64ClusterCA "dGVzdA==" -KubeletExtraArgs "--node-labels= --register-with-taints=foo=bar:NoSchedule" 3>&1 4>&1 5>&1 6>&1
</powershell>
`,
		}),

		Entry("with a preBootstrapCommand", windowsEntry{
			updateNodeGroup: func(ng *api.ManagedNodeGroup) {
				ng.PreBootstrapCommands = []string{
					"wget -UseBasicParsing -O amazon-cloudwatch-agent.msi https://s3.amazonaws.com/amazoncloudwatch-agent/windows/amd64/latest/amazon-cloudwatch-agent.msi",
				}
			},

			expectedUserData: `
<powershell>
[string]$EKSBootstrapScriptFile = "$env:ProgramFiles\Amazon\EKS\Start-EKSBootstrap.ps1"
wget -UseBasicParsing -O amazon-cloudwatch-agent.msi https://s3.amazonaws.com/amazoncloudwatch-agent/windows/amd64/latest/amazon-cloudwatch-agent.msi
& $EKSBootstrapScriptFile -EKSClusterName "windohs" -APIServerEndpoint "https://test.com" -Base64ClusterCA "dGVzdA==" -KubeletExtraArgs "--node-labels= --register-with-taints=" 3>&1 4>&1 5>&1 6>&1
</powershell>
`,
		}),

		Entry("with several preBootstrapCommands", windowsEntry{
			updateNodeGroup: func(ng *api.ManagedNodeGroup) {
				ng.PreBootstrapCommands = []string{
					"wget -UseBasicParsing -O amazon-cloudwatch-agent.msi https://s3.amazonaws.com/amazoncloudwatch-agent/windows/amd64/latest/amazon-cloudwatch-agent.msi",
					"start /wait msiexec.exe /qb /i \"amazon-cloudwatch-agent.msi\"",
				}

			},

			expectedUserData: `
<powershell>
[string]$EKSBootstrapScriptFile = "$env:ProgramFiles\Amazon\EKS\Start-EKSBootstrap.ps1"
wget -UseBasicParsing -O amazon-cloudwatch-agent.msi https://s3.amazonaws.com/amazoncloudwatch-agent/windows/amd64/latest/amazon-cloudwatch-agent.msi
start /wait msiexec.exe /qb /i "amazon-cloudwatch-agent.msi"
& $EKSBootstrapScriptFile -EKSClusterName "windohs" -APIServerEndpoint "https://test.com" -Base64ClusterCA "dGVzdA==" -KubeletExtraArgs "--node-labels= --register-with-taints=" 3>&1 4>&1 5>&1 6>&1
</powershell>
`,
		}),
	)

})
