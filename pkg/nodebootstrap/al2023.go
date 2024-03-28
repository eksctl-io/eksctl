package nodebootstrap

import (
	"bytes"
	"encoding/base64"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/assets"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/utils"
)

type AL2023 struct {
	cfg        *api.ClusterConfig
	ng         *api.NodeGroupBase
	taints     []api.NodeGroupTaint
	clusterDNS string

	scripts    []string
	cloudboot  []string
	nodeConfig *NodeConfig

	UserDataMimeBoundary string
}

func NewManagedAL2023Bootstrapper(cfg *api.ClusterConfig, mng *api.ManagedNodeGroup, clusterDNS string) *AL2023 {
	al2023 := newAL2023Bootstrapper(cfg, mng, clusterDNS)
	if api.IsEnabled(mng.EFAEnabled) {
		al2023.cloudboot = append(al2023.cloudboot, assets.EfaManagedAL2023Boothook)
	}
	if api.IsAMI(mng.AMI) {
		al2023.setNodeConfig()
	}
	return al2023
}

func NewAL2023Bootstrapper(cfg *api.ClusterConfig, ng *api.NodeGroup, clusterDNS string) *AL2023 {
	al2023 := newAL2023Bootstrapper(cfg, ng, clusterDNS)
	if api.IsEnabled(ng.EFAEnabled) {
		al2023.scripts = append(al2023.scripts, assets.EfaAl2023Sh)
	}
	al2023.setNodeConfig()
	return al2023
}

func newAL2023Bootstrapper(cfg *api.ClusterConfig, np api.NodePool, clusterDNS string) *AL2023 {
	return &AL2023{
		cfg:        cfg,
		ng:         np.BaseNodeGroup(),
		taints:     np.NGTaints(),
		clusterDNS: clusterDNS,
	}
}

func (m *AL2023) UserData() (string, error) {
	var buf bytes.Buffer

	if len(m.scripts) == 0 && len(m.cloudboot) == 0 && m.nodeConfig == nil {
		return "", nil
	}

	if err := createMimeMessage(&buf, m.scripts, m.cloudboot, m.nodeConfig, m.UserDataMimeBoundary); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func (m *AL2023) setNodeConfig() {
	m.nodeConfig = &NodeConfig{
		APIVersion: "node.eks.aws/v1alpha1",
		Kind:       "NodeConfig",
		Spec: NodeSpec{
			Cluster: ClusterSpec{
				Name:                 m.cfg.Metadata.Name,
				APIServerEndpoint:    m.cfg.Status.Endpoint,
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
		m.nodeConfig.Spec.Kubelet.Config.MaxPods = &m.ng.MaxPodsPerNode
	}
	if len(m.taints) > 0 {
		m.nodeConfig.Spec.Kubelet.Flags = append(m.nodeConfig.Spec.Kubelet.Flags, "--register-with-taints="+utils.FormatTaints(m.taints))
	}
}

type NodeConfig struct {
	APIVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	Spec       NodeSpec `yaml:"spec"`
}

type NodeSpec struct {
	Cluster ClusterSpec `yaml:"cluster"`
	Kubelet KubeletSpec `yaml:"kubelet"`
}

type ClusterSpec struct {
	APIServerEndpoint    string `yaml:"apiServerEndpoint"`
	CertificateAuthority string `yaml:"certificateAuthority"`
	CIDR                 string `yaml:"cidr"`
	Name                 string `yaml:"name"`
}

type KubeletSpec struct {
	Config KubeletConfig `yaml:"config"`
	Flags  []string      `yaml:"flags"`
}

type KubeletConfig struct {
	MaxPods    *int     `yaml:"maxPods,omitempty"`
	ClusterDNS []string `yaml:"clusterDNS"`
}
