package nodebootstrap

import (
	"bytes"
	"encoding/base64"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/assets"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/utils"
)

// AL2023 is a bootstrapper for both EKS-managed and self-managed AmazonLinux2023 nodegroups
type AL2023 struct {
	cfg                  *api.ClusterConfig
	ng                   *api.NodeGroupBase
	taints               []api.NodeGroupTaint
	clusterDNS           string
	UserDataMimeBoundary string
}

func NewAL2023Bootstrapper(cfg *api.ClusterConfig, np api.NodePool, clusterDNS string) *AL2023 {
	return &AL2023{
		cfg:        cfg,
		ng:         np.BaseNodeGroup(),
		taints:     np.NGTaints(),
		clusterDNS: clusterDNS,
	}
}

func (m *AL2023) UserData() (string, error) {
	var (
		buf       bytes.Buffer
		cloudboot []string
	)

	if api.IsEnabled(m.ng.EFAEnabled) {
		cloudboot = append(cloudboot, assets.EfaManagedBoothook)
	}

	if err := createMimeMessage(&buf, []string{}, cloudboot, m.makeNodeConfig(), m.UserDataMimeBoundary); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func (m *AL2023) makeNodeConfig() *NodeConfig {
	nodeConfig := &NodeConfig{
		ApiVersion: "node.eks.aws/v1alpha1",
		Kind:       "NodeConfig",
		Spec: NodeSpec{
			Cluster: ClusterSpec{
				Name:                 m.cfg.Metadata.Name,
				ApiServerEndpoint:    m.cfg.Status.Endpoint,
				CertificateAuthority: base64.StdEncoding.EncodeToString(m.cfg.Status.CertificateAuthorityData),
				CIDR:                 m.cfg.Status.KubernetesNetworkConfig.ServiceIPv4CIDR,
			},
			Kubelet: KubeletSpec{
				Config: KubeletConfig{
					ClusterDNS: []string{m.clusterDNS},
				},
				Flags: []string{"--node-labels=" + formatLabels(m.ng.Labels)},
			},
		},
	}

	if m.ng.MaxPodsPerNode > 0 {
		nodeConfig.Spec.Kubelet.Config.MaxPods = &m.ng.MaxPodsPerNode
	}

	if len(m.taints) > 0 {
		nodeConfig.Spec.Kubelet.Flags = append(nodeConfig.Spec.Kubelet.Flags, utils.FormatTaints(m.taints))
	}

	return nodeConfig
}

// NodeConfig represents the core EKS node configuration
type NodeConfig struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Spec       NodeSpec `yaml:"spec"`
}

// NodeSpec encapsulates the 'spec' section of the YAML
type NodeSpec struct {
	Cluster ClusterSpec `yaml:"cluster"`
	Kubelet KubeletSpec `yaml:"kubelet"`
}

// ClusterSpec holds cluster-related parameters
type ClusterSpec struct {
	ApiServerEndpoint    string `yaml:"apiServerEndpoint"`
	CertificateAuthority string `yaml:"certificateAuthority"`
	CIDR                 string `yaml:"cidr"`
	Name                 string `yaml:"name"`
}

// KubeletSpec captures Kubelet parameters and flags
type KubeletSpec struct {
	Config KubeletConfig `yaml:"config"`
	Flags  []string      `yaml:"flags"`
}

// KubeletConfig holds the 'config' section
type KubeletConfig struct {
	MaxPods    *int     `yaml:"maxPods,omitempty"`
	ClusterDNS []string `yaml:"clusterDNS"`
}
