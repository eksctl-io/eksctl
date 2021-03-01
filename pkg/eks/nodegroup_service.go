package eks

import (
	"github.com/kris-nova/logger"
	"github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ssh"
)

// A NodeGroupService provides helpers for nodegroup creation
type NodeGroupService struct {
	cluster  *api.ClusterConfig
	provider api.ClusterProvider
}

// NewNodeGroupService creates a new NodeGroupService
func NewNodeGroupService(clusterConfig *api.ClusterConfig, provider api.ClusterProvider) *NodeGroupService {
	return &NodeGroupService{
		cluster:  clusterConfig,
		provider: provider,
	}
}

// Normalize normalizes nodegroups
func (m *NodeGroupService) Normalize(nodePools []api.NodePool) error {
	for _, np := range nodePools {
		switch ng := np.(type) {
		case *api.NodeGroup:
			// resolve AMI
			if !api.IsAMI(ng.AMI) {
				if err := ResolveAMI(m.provider, m.cluster.Metadata.Version, ng); err != nil {
					return err
				}
			}
			logger.Info("nodegroup %q will use %q [%s/%s]", ng.Name, ng.AMI, ng.AMIFamily, m.cluster.Metadata.Version)
		}

		ng := np.BaseNodeGroup()
		if ng.AMI != "" {
			if err := ami.Use(m.provider.EC2(), ng); err != nil {
				return err
			}
		}
		// load or use SSH key - name includes cluster name and the
		// fingerprint, so if unique keys are provided, each will get
		// loaded and used as intended and there is no need to have
		// nodegroup name in the key name
		publicKeyName, err := ssh.LoadKey(ng.SSH, m.cluster.Metadata.Name, ng.Name, m.provider.EC2())
		if err != nil {
			return err
		}
		if publicKeyName != "" {
			ng.SSH.PublicKeyName = &publicKeyName
		}
	}
	return nil
}
