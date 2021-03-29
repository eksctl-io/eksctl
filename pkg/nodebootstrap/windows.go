package nodebootstrap

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type Windows struct {
	clusterName string
	ng          *api.NodeGroup
}

func NewWindowsBootstrapper(clusterName string, ng *api.NodeGroup) *Windows {
	return &Windows{
		clusterName: clusterName,
		ng:          ng,
	}
}

func (b *Windows) UserData() (string, error) {
	bootstrapScript := `<powershell>
[string]$EKSBootstrapScriptFile = "$env:ProgramFiles\Amazon\EKS\Start-EKSBootstrap.ps1"
`
	for _, command := range b.ng.PreBootstrapCommands {
		bootstrapScript += fmt.Sprintf("%s\n", command)
	}

	kubeletOptions := map[string]string{
		"node-labels":          kvs(b.ng.Labels),
		"register-with-taints": kvs(b.ng.Taints),
	}
	if b.ng.MaxPodsPerNode != 0 {
		kubeletOptions["max-pods"] = strconv.Itoa(b.ng.MaxPodsPerNode)
	}

	kubeletArgs := toCLIArgs(kubeletOptions)
	bootstrapScript += fmt.Sprintf("& $EKSBootstrapScriptFile -EKSClusterName %q -KubeletExtraArgs %q 3>&1 4>&1 5>&1 6>&1\n</powershell>", b.clusterName, kubeletArgs)

	userData := base64.StdEncoding.EncodeToString([]byte(bootstrapScript))

	logger.Debug("user-data = %s", userData)
	return userData, nil
}

func toCLIArgs(values map[string]string) string {
	var args []string
	for k, v := range values {
		args = append(args, fmt.Sprintf("--%s=%s", k, v))
	}
	sort.Strings(args)
	return strings.Join(args, " ")
}
