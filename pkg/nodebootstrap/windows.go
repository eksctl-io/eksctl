package nodebootstrap

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/utils"
)

type Windows struct {
	clusterConfig *api.ClusterConfig
	ng            *api.NodeGroup
}

type keyValue struct {
	key   string
	value string
}

func NewWindowsBootstrapper(clusterConfig *api.ClusterConfig, ng *api.NodeGroup) *Windows {
	return &Windows{
		clusterConfig: clusterConfig,
		ng:            ng,
	}
}

func (b *Windows) UserData() (string, error) {
	bootstrapCommands := []string{
		`<powershell>
[string]$EKSBootstrapScriptFile = "$env:ProgramFiles\Amazon\EKS\Start-EKSBootstrap.ps1"`,
	}

	bootstrapCommands = append(bootstrapCommands, b.ng.PreBootstrapCommands...)
	eksBootstrapCommand := fmt.Sprintf("& $EKSBootstrapScriptFile %s 3>&1 4>&1 5>&1 6>&1", b.makeBootstrapParams())
	bootstrapCommands = append(bootstrapCommands,
		eksBootstrapCommand,
		"</powershell>",
	)

	userData := base64.StdEncoding.EncodeToString([]byte(strings.Join(bootstrapCommands, "\n")))

	logger.Debug("user-data = %s", userData)
	return userData, nil
}

func (b *Windows) makeBootstrapParams() string {
	params := []keyValue{
		{
			key:   "EKSClusterName",
			value: b.clusterConfig.Metadata.Name,
		},
		{
			key:   "APIServerEndpoint",
			value: b.clusterConfig.Status.Endpoint,
		},
		{
			key:   "Base64ClusterCA",
			value: base64.StdEncoding.EncodeToString(b.clusterConfig.Status.CertificateAuthorityData),
		},
		{
			key:   "DNSClusterIP",
			value: b.ng.ClusterDNS,
		},
		{
			key:   "KubeletExtraArgs",
			value: b.makeKubeletOptions(),
		},
	}

	return formatWindowsParams(params)
}

func (b *Windows) makeKubeletOptions() string {
	kubeletOptions := []keyValue{
		{
			key:   "node-labels",
			value: formatLabels(b.ng.Labels),
		},
		{
			key:   "register-with-taints",
			value: utils.FormatTaints(b.ng.Taints),
		},
	}

	if b.ng.MaxPodsPerNode != 0 {
		kubeletOptions = append(kubeletOptions, keyValue{
			key:   "max-pods",
			value: strconv.Itoa(b.ng.MaxPodsPerNode),
		})
	}

	return toCLIArgs(kubeletOptions)
}

// formatWindowsParams formats params into `-key "value"`, ignoring keys with empty values
func formatWindowsParams(params []keyValue) string {
	var args []string
	for _, param := range params {
		if param.value != "" {
			args = append(args, fmt.Sprintf("-%s %q", param.key, param.value))
		}
	}
	return strings.Join(args, " ")
}

func toCLIArgs(params []keyValue) string {
	var args []string
	for _, param := range params {
		args = append(args, fmt.Sprintf("--%s=%s", param.key, param.value))
	}
	return strings.Join(args, " ")
}
