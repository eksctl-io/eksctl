package eks

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/ssh"
)

type NodeGroupService struct {
	cluster *v1alpha5.ClusterConfig
	ec2API  ec2iface.EC2API
}

func NewNodeGroupService(clusterConfig *v1alpha5.ClusterConfig, ec2API ec2iface.EC2API) *NodeGroupService {
	return &NodeGroupService{
		cluster: clusterConfig,
		ec2API:  ec2API,
	}
}

func (m *NodeGroupService) NormalizeManaged(nodeGroups []*v1alpha5.ManagedNodeGroup) error {
	for _, ng := range nodeGroups {
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
