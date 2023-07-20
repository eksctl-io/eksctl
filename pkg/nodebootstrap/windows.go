package nodebootstrap

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/powershell"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/utils"
)

type Windows struct {
	clusterConfig *api.ClusterConfig
	np            api.NodePool
	clusterDNS    string
}

func NewWindowsBootstrapper(clusterConfig *api.ClusterConfig, np api.NodePool, clusterDNS string) *Windows {
	return &Windows{
		clusterConfig: clusterConfig,
		np:            np,
		clusterDNS:    clusterDNS,
	}
}

func (b *Windows) UserData() (string, error) {
	ng := b.np.BaseNodeGroup()
	bootstrapCommands := append([]string{
		`<powershell>
[string]$EKSBootstrapScriptFile = "$env:ProgramFiles\Amazon\EKS\Start-EKSBootstrap.ps1"`,
	}, ng.PreBootstrapCommands...)

	if ng.OverrideBootstrapCommand != nil {
		bootstrapCommands = append(bootstrapCommands,
			b.makeBootstrapParams(true),
			*ng.OverrideBootstrapCommand,
		)
	} else {
		bootstrapCommands = append(bootstrapCommands, fmt.Sprintf("& $EKSBootstrapScriptFile %s 3>&1 4>&1 5>&1 6>&1", b.makeBootstrapParams(false)))
	}

	bootstrapCommands = append(bootstrapCommands, "</powershell>")
	userData := base64.StdEncoding.EncodeToString([]byte(strings.Join(bootstrapCommands, "\n")))
	logger.Debug("user-data = %s", userData)
	return userData, nil
}

func (b *Windows) makeBootstrapParams(hasBootstrapCommand bool) string {
	params := []powershell.KeyValue{
		{
			Key:   "EKSClusterName",
			Value: b.clusterConfig.Metadata.Name,
		},
		{
			Key:   "APIServerEndpoint",
			Value: b.clusterConfig.Status.Endpoint,
		},
		{
			Key:   "Base64ClusterCA",
			Value: base64.StdEncoding.EncodeToString(b.clusterConfig.Status.CertificateAuthorityData),
		},
		{
			Key:   "ServiceCIDR",
			Value: b.clusterConfig.Status.KubernetesNetworkConfig.ServiceIPv4CIDR,
		},
	}
	if unmanaged, ok := b.np.(*api.NodeGroup); ok {
		// DNSClusterIP is only configurable for self-managed nodegroups.
		if b.clusterDNS != "" {
			params = append(params, powershell.KeyValue{
				Key:   "DNSClusterIP",
				Value: b.clusterDNS,
			})
		}
		// ContainerRuntime is only configurable for self-managed nodegroups.
		if unmanaged.ContainerRuntime != nil {
			params = append(params, powershell.KeyValue{
				Key:   "ContainerRuntime",
				Value: *unmanaged.ContainerRuntime,
			})
		}
	}

	kubeletOptions := b.makeKubeletOptions()

	params = append(params, powershell.KeyValue{
		Key:   "KubeletExtraArgs",
		Value: powershell.ToCLIArgs(kubeletOptions),
	})

	if hasBootstrapCommand {
		return powershell.JoinVariables(powershell.FormatStringVariables(params), powershell.FormatHashTable(kubeletOptions, "KubeletExtraArgsMap"))
	}
	return powershell.FormatParams(params)
}

func (b *Windows) makeKubeletOptions() []powershell.KeyValue {
	ng := b.np.BaseNodeGroup()

	kubeletOptions := []powershell.KeyValue{
		{
			Key:   "node-labels",
			Value: formatLabels(ng.Labels),
		},
		{
			Key:   "register-with-taints",
			Value: utils.FormatTaints(b.np.NGTaints()),
		},
	}

	if ng.MaxPodsPerNode != 0 {
		kubeletOptions = append(kubeletOptions, powershell.KeyValue{
			Key:   "max-pods",
			Value: strconv.Itoa(ng.MaxPodsPerNode),
		})
	}
	return kubeletOptions
}
