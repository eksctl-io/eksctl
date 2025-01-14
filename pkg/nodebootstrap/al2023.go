package nodebootstrap

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"sigs.k8s.io/yaml"

	nodeadmapi "github.com/awslabs/amazon-eks-ami/nodeadm/api"
	nodeadm "github.com/awslabs/amazon-eks-ami/nodeadm/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/assets"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap/utils"
)

type AL2023 struct {
	cfg        *api.ClusterConfig
	nodePool   api.NodePool
	clusterDNS string

	scripts     []string
	cloudboot   []string
	nodeConfigs []*nodeadm.NodeConfig

	UserDataMimeBoundary string
}

func NewManagedAL2023Bootstrapper(cfg *api.ClusterConfig, mng *api.ManagedNodeGroup, clusterDNS string) *AL2023 {
	al2023 := newAL2023Bootstrapper(cfg, mng, clusterDNS)
	return al2023
}

func NewAL2023Bootstrapper(cfg *api.ClusterConfig, ng *api.NodeGroup, clusterDNS string) *AL2023 {
	al2023 := newAL2023Bootstrapper(cfg, ng, clusterDNS)
	return al2023
}

func newAL2023Bootstrapper(cfg *api.ClusterConfig, np api.NodePool, clusterDNS string) *AL2023 {
	return &AL2023{
		cfg:        cfg,
		nodePool:   np,
		clusterDNS: clusterDNS,
		scripts:    []string{assets.AL2023XTablesLock},
	}
}

func (m *AL2023) UserData() (string, error) {
	ng := m.nodePool.BaseNodeGroup()

	minimalNodeConfig, err := m.createMinimalNodeConfig()
	if err != nil {
		return "", fmt.Errorf("generating minimal node config: %w", err)
	}
	if minimalNodeConfig != nil {
		m.nodeConfigs = append(m.nodeConfigs, minimalNodeConfig)
	}

	if ng.MaxPodsPerNode > 0 {
		nodeConfig := &nodeadm.NodeConfig{
			TypeMeta: metav1.TypeMeta{
				Kind:       nodeadmapi.KindNodeConfig,
				APIVersion: nodeadm.GroupVersion.String(),
			},
		}
		kubeletConfig, err := ToKubeletConfig(api.InlineDocument{"maxPods": ng.MaxPodsPerNode})
		if err != nil {
			return "", err
		}
		nodeConfig.Spec.Kubelet.Config = kubeletConfig
		m.nodeConfigs = append(m.nodeConfigs, nodeConfig)
	}

	for _, command := range m.nodePool.BaseNodeGroup().PreBootstrapCommands {
		m.scripts = append(m.scripts, "#!/bin/bash\n"+command)
	}

	if ng.OverrideBootstrapCommand != nil {
		nodeConfig, err := stringToNodeConfig(*ng.OverrideBootstrapCommand)
		if err != nil {
			return "", err
		}
		m.nodeConfigs = append(m.nodeConfigs, nodeConfig)
	}

	if len(m.scripts) == 0 && len(m.cloudboot) == 0 && len(m.nodeConfigs) == 0 {
		return "", nil
	}

	var buf bytes.Buffer
	if err := createMimeMessage(&buf, m.scripts, m.cloudboot, m.nodeConfigs, m.UserDataMimeBoundary); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func (m *AL2023) createMinimalNodeConfig() (*nodeadm.NodeConfig, error) {
	kubeletConfig := api.InlineDocument{}
	switch nodeGroup := m.nodePool.(type) {
	case *api.ManagedNodeGroup:
		if !api.IsAMI(nodeGroup.AMI) {
			return nil, nil
		}
	case *api.NodeGroup:
		if nodeGroup.KubeletExtraConfig != nil {
			kubeletConfig = *nodeGroup.KubeletExtraConfig.DeepCopy()
		}
	}

	kubeletConfig["clusterDNS"] = []string{m.clusterDNS}
	ng := m.nodePool.BaseNodeGroup()
	nodeKubeletConfig, err := ToKubeletConfig(kubeletConfig)
	if err != nil {
		return nil, err
	}

	kubeletOptions := nodeadm.KubeletOptions{
		Flags:  []string{"--node-labels=" + formatLabels(ng.Labels)},
		Config: nodeKubeletConfig,
	}
	if taints := m.nodePool.NGTaints(); len(taints) > 0 {
		kubeletOptions.Flags = append(kubeletOptions.Flags, "--register-with-taints="+utils.FormatTaints(taints))
	}

	clusterStatus := m.cfg.Status
	var serviceCIDR string
	if clusterStatus.KubernetesNetworkConfig.ServiceIPv6CIDR != "" {
		serviceCIDR = clusterStatus.KubernetesNetworkConfig.ServiceIPv6CIDR
	} else {
		serviceCIDR = clusterStatus.KubernetesNetworkConfig.ServiceIPv4CIDR
	}

	return &nodeadm.NodeConfig{
		TypeMeta: metav1.TypeMeta{
			Kind:       nodeadmapi.KindNodeConfig,
			APIVersion: nodeadm.GroupVersion.String(),
		},
		Spec: nodeadm.NodeConfigSpec{
			Cluster: nodeadm.ClusterDetails{
				Name:                 m.cfg.Metadata.Name,
				APIServerEndpoint:    clusterStatus.Endpoint,
				CertificateAuthority: clusterStatus.CertificateAuthorityData,
				CIDR:                 serviceCIDR,
			},
			Kubelet: kubeletOptions,
		},
	}, nil
}

// ToKubeletConfig generates a kubelet config that can be used with nodeadm.NodeConfig.
func ToKubeletConfig(kubeletExtraConfig api.InlineDocument) (map[string]runtime.RawExtension, error) {
	kubeletConfig := map[string]runtime.RawExtension{}
	for k, v := range kubeletExtraConfig {
		raw, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		kubeletConfig[k] = runtime.RawExtension{Raw: raw}
	}
	return kubeletConfig, nil
}

func stringToNodeConfig(overrideBootstrapCommand string) (*nodeadm.NodeConfig, error) {
	var config nodeadm.NodeConfig
	err := yaml.Unmarshal([]byte(overrideBootstrapCommand), &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling \"overrideBootstrapCommand\" into \"nodeadm.NodeConfig\": %w", err)
	}
	return &config, nil
}
