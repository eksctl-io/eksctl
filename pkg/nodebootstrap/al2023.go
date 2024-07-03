package nodebootstrap

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"

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

	scripts   []string
	cloudboot []string

	UserDataMimeBoundary string
}

func NewManagedAL2023Bootstrapper(cfg *api.ClusterConfig, mng *api.ManagedNodeGroup, clusterDNS string) *AL2023 {
	al2023 := newAL2023Bootstrapper(cfg, mng, clusterDNS)
	if api.IsEnabled(mng.EFAEnabled) {
		al2023.cloudboot = append(al2023.cloudboot, assets.EfaManagedAL2023Boothook)
	}
	return al2023
}

func NewAL2023Bootstrapper(cfg *api.ClusterConfig, ng *api.NodeGroup, clusterDNS string) *AL2023 {
	al2023 := newAL2023Bootstrapper(cfg, ng, clusterDNS)
	if api.IsEnabled(ng.EFAEnabled) {
		al2023.scripts = append(al2023.scripts, assets.EfaAl2023Sh)
	}
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
	nodeConfig, err := m.createNodeConfig()
	if err != nil {
		return "", fmt.Errorf("generating node config: %w", err)
	}
	if len(m.scripts) == 0 && len(m.cloudboot) == 0 && nodeConfig == nil {
		return "", nil
	}

	var buf bytes.Buffer
	if err := createMimeMessage(&buf, m.scripts, m.cloudboot, nodeConfig, m.UserDataMimeBoundary); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func (m *AL2023) createNodeConfig() (*nodeadm.NodeConfig, error) {
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
	if ng.MaxPodsPerNode > 0 {
		kubeletConfig["maxPods"] = strconv.Itoa(ng.MaxPodsPerNode)
	}
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
				CIDR:                 clusterStatus.KubernetesNetworkConfig.ServiceIPv4CIDR,
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
