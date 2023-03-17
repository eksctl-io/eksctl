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
	np            api.NodePool
	clusterDNS    string
}

type keyValue struct {
	key   string
	value string
}

func NewWindowsBootstrapper(clusterConfig *api.ClusterConfig, np api.NodePool, clusterDNS string) *Windows {
	return &Windows{
		clusterConfig: clusterConfig,
		np:            np,
		clusterDNS:    clusterDNS,
	}
}

func (b *Windows) UserData() (string, error) {
	bootstrapCommands := []string{
		`<powershell>
[string]$EKSBootstrapScriptFile = "$env:ProgramFiles\Amazon\EKS\Start-EKSBootstrap.ps1"`,
	}
	ng := b.np.BaseNodeGroup()
	bootstrapCommands = append(bootstrapCommands, ng.PreBootstrapCommands...)
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
			key:   "ServiceCIDR",
			value: b.clusterConfig.Status.KubernetesNetworkConfig.ServiceIPv4CIDR,
		},
	}
	if unmanaged, ok := b.np.(*api.NodeGroup); ok {
		// DNSClusterIP is only configurable for self-managed nodegroups.
		if b.clusterDNS != "" {
			params = append(params, keyValue{
				key:   "DNSClusterIP",
				value: b.clusterDNS,
			})
		}
		// ContainerRuntime is only configurable for self-managed nodegroups.
		if unmanaged.ContainerRuntime != nil {
			params = append(params, keyValue{
				key:   "ContainerRuntime",
				value: *unmanaged.ContainerRuntime,
			})
		}
	}

	params = append(params, keyValue{
		key:   "KubeletExtraArgs",
		value: b.makeKubeletOptions(),
	})
	return formatWindowsParams(params)
}

func (b *Windows) makeKubeletOptions() string {
	ng := b.np.BaseNodeGroup()

	kubeletOptions := []keyValue{
		{
			key:   "node-labels",
			value: formatLabels(ng.Labels),
		},
		{
			key:   "register-with-taints",
			value: utils.FormatTaints(b.np.NGTaints()),
		},
	}

	if ng.MaxPodsPerNode != 0 {
		kubeletOptions = append(kubeletOptions, keyValue{
			key:   "max-pods",
			value: strconv.Itoa(ng.MaxPodsPerNode),
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
