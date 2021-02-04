package nodebootstrap

import (
	"encoding/base64"
	"fmt"
	"strconv"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/logger"
)

func NewUserDataForWindows(spec *api.ClusterConfig, ng *api.NodeGroup) (string, error) {
	bootstrapScript := `<powershell>
[string]$EKSBootstrapScriptFile = "$env:ProgramFiles\Amazon\EKS\Start-EKSBootstrap.ps1"
`
	for _, command := range ng.PreBootstrapCommands {
		bootstrapScript += fmt.Sprintf("%s\n", command)
	}

	kubeletOptions := map[string]string{
		"node-labels":          kvs(ng.Labels),
		"register-with-taints": kvs(ng.Taints),
	}
	if ng.MaxPodsPerNode != 0 {
		kubeletOptions["max-pods"] = strconv.Itoa(ng.MaxPodsPerNode)
	}

	kubeletArgs := toCLIArgs(kubeletOptions)
	bootstrapScript += fmt.Sprintf("& $EKSBootstrapScriptFile -EKSClusterName %q -KubeletExtraArgs %q 3>&1 4>&1 5>&1 6>&1\n</powershell>", spec.Metadata.Name, kubeletArgs)

	userData := base64.StdEncoding.EncodeToString([]byte(bootstrapScript))

	logger.Debug("user-data = %s", userData)
	return userData, nil
}
