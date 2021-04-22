package nodebootstrap_test

import (
	"encoding/base64"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

var _ = Describe("Windows", func() {
	var (
		clusterName string
		ng          *api.NodeGroup
	)

	BeforeEach(func() {
		clusterName = "windohs"
		ng = &api.NodeGroup{
			NodeGroupBase: &api.NodeGroupBase{
				AMIFamily: api.NodeImageFamilyWindowsServer2019CoreContainer,
			},
		}
	})

	It("produces correct standard userdata", func() {
		ng.PreBootstrapCommands = nil
		bootstrap := nodebootstrap.NewWindowsBootstrapper(clusterName, ng)
		userdata, err := bootstrap.UserData()
		Expect(err).ToNot(HaveOccurred())

		Expect(decodeData(userdata)).To(Equal(strings.TrimSpace(`
<powershell>
[string]$EKSBootstrapScriptFile = "$env:ProgramFiles\Amazon\EKS\Start-EKSBootstrap.ps1"
& $EKSBootstrapScriptFile -EKSClusterName "windohs" -KubeletExtraArgs "--node-labels= --register-with-taints=" 3>&1 4>&1 5>&1 6>&1
</powershell>
`)))
	})

	When("labels are set on the node", func() {
		It("adds them to the userdata", func() {
			ng.Labels = map[string]string{"foo": "bar"}
			bootstrap := nodebootstrap.NewWindowsBootstrapper(clusterName, ng)
			userdata, err := bootstrap.UserData()
			Expect(err).ToNot(HaveOccurred())

			Expect(decodeData(userdata)).To(Equal(strings.TrimSpace(`
<powershell>
[string]$EKSBootstrapScriptFile = "$env:ProgramFiles\Amazon\EKS\Start-EKSBootstrap.ps1"
& $EKSBootstrapScriptFile -EKSClusterName "windohs" -KubeletExtraArgs "--node-labels=foo=bar --register-with-taints=" 3>&1 4>&1 5>&1 6>&1
</powershell>
`)))
		})
	})

	When("taints are set on the node", func() {
		It("adds them to the userdata", func() {
			ng.Taints = map[string]string{"foo": "bar"}
			bootstrap := nodebootstrap.NewWindowsBootstrapper(clusterName, ng)
			userdata, err := bootstrap.UserData()
			Expect(err).ToNot(HaveOccurred())

			Expect(decodeData(userdata)).To(Equal(strings.TrimSpace(`
<powershell>
[string]$EKSBootstrapScriptFile = "$env:ProgramFiles\Amazon\EKS\Start-EKSBootstrap.ps1"
& $EKSBootstrapScriptFile -EKSClusterName "windohs" -KubeletExtraArgs "--node-labels= --register-with-taints=foo=bar" 3>&1 4>&1 5>&1 6>&1
</powershell>
`)))
		})
	})

	When("MaxPodsPerNode are set", func() {
		It("adds them to the userdata", func() {
			ng.MaxPodsPerNode = 100
			bootstrap := nodebootstrap.NewWindowsBootstrapper(clusterName, ng)
			userdata, err := bootstrap.UserData()
			Expect(err).ToNot(HaveOccurred())

			Expect(decodeData(userdata)).To(Equal(strings.TrimSpace(`
<powershell>
[string]$EKSBootstrapScriptFile = "$env:ProgramFiles\Amazon\EKS\Start-EKSBootstrap.ps1"
& $EKSBootstrapScriptFile -EKSClusterName "windohs" -KubeletExtraArgs "--max-pods=100 --node-labels= --register-with-taints=" 3>&1 4>&1 5>&1 6>&1
</powershell>
`)))
		})
	})

	When("a PreBootstrapCommands is set", func() {
		It("adds it to the userdata", func() {
			ng.PreBootstrapCommands = []string{
				"wget -UseBasicParsing -O amazon-cloudwatch-agent.msi https://s3.amazonaws.com/amazoncloudwatch-agent/windows/amd64/latest/amazon-cloudwatch-agent.msi",
			}
			bootstrap := nodebootstrap.NewWindowsBootstrapper(clusterName, ng)
			userdata, err := bootstrap.UserData()
			Expect(err).ToNot(HaveOccurred())

			Expect(decodeData(userdata)).To(Equal(strings.TrimSpace(`
<powershell>
[string]$EKSBootstrapScriptFile = "$env:ProgramFiles\Amazon\EKS\Start-EKSBootstrap.ps1"
wget -UseBasicParsing -O amazon-cloudwatch-agent.msi https://s3.amazonaws.com/amazoncloudwatch-agent/windows/amd64/latest/amazon-cloudwatch-agent.msi
& $EKSBootstrapScriptFile -EKSClusterName "windohs" -KubeletExtraArgs "--node-labels= --register-with-taints=" 3>&1 4>&1 5>&1 6>&1
</powershell>
`)))
		})
	})

	When("several PreBootstrapCommands are set", func() {
		It("adds them to the userdata", func() {
			ng.PreBootstrapCommands = []string{
				"wget -UseBasicParsing -O amazon-cloudwatch-agent.msi https://s3.amazonaws.com/amazoncloudwatch-agent/windows/amd64/latest/amazon-cloudwatch-agent.msi",
				"start /wait msiexec.exe /qb /i \"amazon-cloudwatch-agent.msi\"",
			}
			bootstrap := nodebootstrap.NewWindowsBootstrapper(clusterName, ng)
			userdata, err := bootstrap.UserData()
			Expect(err).ToNot(HaveOccurred())

			Expect(decodeData(userdata)).To(Equal(strings.TrimSpace(`
<powershell>
[string]$EKSBootstrapScriptFile = "$env:ProgramFiles\Amazon\EKS\Start-EKSBootstrap.ps1"
wget -UseBasicParsing -O amazon-cloudwatch-agent.msi https://s3.amazonaws.com/amazoncloudwatch-agent/windows/amd64/latest/amazon-cloudwatch-agent.msi
start /wait msiexec.exe /qb /i "amazon-cloudwatch-agent.msi"
& $EKSBootstrapScriptFile -EKSClusterName "windohs" -KubeletExtraArgs "--node-labels= --register-with-taints=" 3>&1 4>&1 5>&1 6>&1
</powershell>
`)))
		})
	})
})

func decodeData(userdata string) string {
	decodedBytes, err := base64.StdEncoding.DecodeString(userdata)
	Expect(err).ToNot(HaveOccurred())

	decodedString := string(decodedBytes)
	Expect(decodedString).ToNot(Equal(""))

	return decodedString
}
