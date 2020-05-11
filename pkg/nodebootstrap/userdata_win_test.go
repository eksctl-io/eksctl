package nodebootstrap

import (
	"encoding/base64"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("Windows", func() {
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
		ng = &api.NodeGroup{
			AMIFamily: api.NodeImageFamilyWindowsServer2019CoreContainer,
		}
	})

	Describe("with single pre bootstrap script", func() {
		It("produces correct userdata", func() {
			ng.PreBootstrapCommands = []string{
				"wget -UseBasicParsing -O amazon-cloudwatch-agent.msi https://s3.amazonaws.com/amazoncloudwatch-agent/windows/amd64/latest/amazon-cloudwatch-agent.msi",
			}
			userdata, err := NewUserDataForWindows(clusterConfig, ng)
			Expect(err).ToNot(HaveOccurred())

			decodedBytes, err := base64.StdEncoding.DecodeString(userdata)
			Expect(err).ToNot(HaveOccurred())

			decodedString := string(decodedBytes)

			Expect(decodedString).ToNot(Equal(""))
			Expect(decodedString).To(Equal(strings.TrimSpace(`
<powershell>
[string]$EKSBootstrapScriptFile = "$env:ProgramFiles\Amazon\EKS\Start-EKSBootstrap.ps1"
wget -UseBasicParsing -O amazon-cloudwatch-agent.msi https://s3.amazonaws.com/amazoncloudwatch-agent/windows/amd64/latest/amazon-cloudwatch-agent.msi
& $EKSBootstrapScriptFile -EKSClusterName "unit-test" -KubeletExtraArgs "--node-labels= --register-with-taints=" 3>&1 4>&1 5>&1 6>&1
</powershell>
`)))
		})
	})

	Describe("with multiple pre bootstrap script", func() {
		It("produces correct userdata", func() {
			ng.PreBootstrapCommands = []string{
				"wget -UseBasicParsing -O amazon-cloudwatch-agent.msi https://s3.amazonaws.com/amazoncloudwatch-agent/windows/amd64/latest/amazon-cloudwatch-agent.msi",
				"start /wait msiexec.exe /qb /i \"amazon-cloudwatch-agent.msi\"",
			}
			userdata, err := NewUserDataForWindows(clusterConfig, ng)
			Expect(err).ToNot(HaveOccurred())

			decodedBytes, err := base64.StdEncoding.DecodeString(userdata)
			Expect(err).ToNot(HaveOccurred())

			decodedString := string(decodedBytes)

			Expect(decodedString).ToNot(Equal(""))
			Expect(decodedString).To(Equal(strings.TrimSpace(`
<powershell>
[string]$EKSBootstrapScriptFile = "$env:ProgramFiles\Amazon\EKS\Start-EKSBootstrap.ps1"
wget -UseBasicParsing -O amazon-cloudwatch-agent.msi https://s3.amazonaws.com/amazoncloudwatch-agent/windows/amd64/latest/amazon-cloudwatch-agent.msi
start /wait msiexec.exe /qb /i "amazon-cloudwatch-agent.msi"
& $EKSBootstrapScriptFile -EKSClusterName "unit-test" -KubeletExtraArgs "--node-labels= --register-with-taints=" 3>&1 4>&1 5>&1 6>&1
</powershell>
`)))
		})
	})

	Describe("without pre bootstrap scripts", func() {
		It("produces correct userdata", func() {
			ng.PreBootstrapCommands = nil
			userdata, err := NewUserDataForWindows(clusterConfig, ng)
			Expect(err).ToNot(HaveOccurred())

			decodedBytes, err := base64.StdEncoding.DecodeString(userdata)
			Expect(err).ToNot(HaveOccurred())

			decodedString := string(decodedBytes)

			Expect(decodedString).ToNot(Equal(""))
			Expect(decodedString).To(Equal(strings.TrimSpace(`
<powershell>
[string]$EKSBootstrapScriptFile = "$env:ProgramFiles\Amazon\EKS\Start-EKSBootstrap.ps1"
& $EKSBootstrapScriptFile -EKSClusterName "unit-test" -KubeletExtraArgs "--node-labels= --register-with-taints=" 3>&1 4>&1 5>&1 6>&1
</powershell>
`)))
		})
	})
})
