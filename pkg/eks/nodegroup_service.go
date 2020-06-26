package eks

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ssh"
)

// A NodeGroupService provides helpers for nodegroup creation
type NodeGroupService struct {
	cluster *api.ClusterConfig
	ec2API  ec2iface.EC2API
}

// NewNodeGroupService creates a new NodeGroupService
func NewNodeGroupService(clusterConfig *api.ClusterConfig, ec2API ec2iface.EC2API) *NodeGroupService {
	return &NodeGroupService{
		cluster: clusterConfig,
		ec2API:  ec2API,
	}
}

// Normalize normalizes nodegroups
func (m *NodeGroupService) Normalize(nodeGroups []*api.NodeGroupBase) error {
	for _, ng := range nodeGroups {
		// load or use SSH key - name includes cluster name and the
		// fingerprint, so if unique keys are provided, each will get
		// loaded and used as intended and there is no need to have
		// nodegroup name in the key name
		publicKeyName, err := ssh.LoadKey(ng.SSH, m.cluster.Metadata.Name, ng.Name, m.ec2API)
		if err != nil {
			return err
		}
		if publicKeyName != "" {
			ng.SSH.PublicKeyName = &publicKeyName
		}
	}
	return nil
}
