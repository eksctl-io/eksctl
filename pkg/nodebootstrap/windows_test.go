package nodebootstrap_test

import (
	"encoding/base64"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

var _ = Describe("Windows", func() {

	type windowsEntry struct {
		updateNodeGroup func(*api.NodeGroup)

		expectedUserData string
	}

	DescribeTable("Windows bootstrap", func(e windowsEntry) {
		clusterConfig := api.NewClusterConfig()
		clusterConfig.Metadata.Name = "windohs"
		clusterConfig.Status = &api.ClusterStatus{
			Endpoint:                 "https://test.com",
			CertificateAuthorityData: []byte("test"),
		}
		ng := &api.NodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				AMIFamily: api.NodeImageFamilyWindowsServer2019CoreContainer,
			},
		}
		if e.updateNodeGroup != nil {
			e.updateNodeGroup(ng)
		}

		bootstrapper := nodebootstrap.NewWindowsBootstrapper(clusterConfig, ng)
		userData, err := bootstrapper.UserData()
		Expect(err).NotTo(HaveOccurred())

		Expect(decodeData(userData)).To(Equal(strings.TrimSpace(e.expectedUserData)))
	},
		Entry("standard userdata", windowsEntry{

			expectedUserData: `
<powershell>
[string]$EKSBootstrapScriptFile = "$env:ProgramFiles\Amazon\EKS\Start-EKSBootstrap.ps1"
& $EKSBootstrapScriptFile -EKSClusterName "windohs" -APIServerEndpoint "https://test.com" -Base64ClusterCA "dGVzdA==" -KubeletExtraArgs "--node-labels= --register-with-taints=" 3>&1 4>&1 5>&1 6>&1
</powershell>
`,
		}),

		Entry("with labels", windowsEntry{
			updateNodeGroup: func(ng *api.NodeGroup) {
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
			updateNodeGroup: func(ng *api.NodeGroup) {
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

		Entry("with maxPods", windowsEntry{
			updateNodeGroup: func(ng *api.NodeGroup) {
				ng.MaxPodsPerNode = 100
			},

			expectedUserData: `
<powershell>
[string]$EKSBootstrapScriptFile = "$env:ProgramFiles\Amazon\EKS\Start-EKSBootstrap.ps1"
& $EKSBootstrapScriptFile -EKSClusterName "windohs" -APIServerEndpoint "https://test.com" -Base64ClusterCA "dGVzdA==" -KubeletExtraArgs "--node-labels= --register-with-taints= --max-pods=100" 3>&1 4>&1 5>&1 6>&1
</powershell>
`,
		}),

		Entry("with a preBootstrapCommand", windowsEntry{
			updateNodeGroup: func(ng *api.NodeGroup) {
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
			updateNodeGroup: func(ng *api.NodeGroup) {
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

func decodeData(userdata string) string {
	decodedBytes, err := base64.StdEncoding.DecodeString(userdata)
	Expect(err).NotTo(HaveOccurred())

	decodedString := string(decodedBytes)
	Expect(decodedString).NotTo(Equal(""))

	return decodedString
}
