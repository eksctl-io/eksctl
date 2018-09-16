package builder

import (
	"fmt"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfn "github.com/awslabs/goformation/cloudformation"

	"github.com/weaveworks/eksctl/pkg/eks/api"

	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

const (
	nodeGroupNameFmt = "${ClusterName}-${NodeGroupID}"
)

var (
	clusterOwnedTag = gfn.Tag{
		Key:   gfn.Sub("kubernetes.io/cluster/${ClusterName}"),
		Value: "owned",
	}
)

type nodeGroupResourceSet struct {
	rs               *resourceSet
	spec             *api.ClusterConfig
	clusterStackName string
	instanceProfile  string
	securityGroups   []string
	vpc              string
	userData         string
}

type awsCloudFormationResource struct {
	Type         string
	Properties   map[string]interface{}
	UpdatePolicy map[string]map[string]string
}

func NewNodeGroupResourceSet(spec *api.ClusterConfig) *nodeGroupResourceSet {
	return &nodeGroupResourceSet{
		rs:   newResourceSet(),
		spec: spec,
	}
}

func (n *nodeGroupResourceSet) AddAllResources() error {
	n.rs.template.Description = nodeGroupTemplateDescription
	n.rs.template.Description += nodeGroupTemplateDescriptionDefaultFeatures
	n.rs.template.Description += templateDescriptionSuffix

	n.vpc = makeImportValue(ParamClusterStackName, cfnOutputClusterVPC)

	userData, err := nodebootstrap.NewUserDataForAmazonLinux2(n.spec)
	if err != nil {
		return err
	}
	n.userData = userData

	n.rs.newStringParameter(ParamClusterName, "")
	n.rs.newStringParameter(ParamClusterStackName, "")
	n.rs.newNumberParameter(ParamNodeGroupID, "")

	if n.spec.MinNodes == 0 && n.spec.MaxNodes == 0 {
		n.spec.MinNodes = n.spec.Nodes
		n.spec.MaxNodes = n.spec.Nodes
	}

	n.addResourcesForIAM()
	n.addResourcesForSecurityGroups()
	n.addResourcesForNodeGroup()

	return nil
}

func (n *nodeGroupResourceSet) RenderJSON() ([]byte, error) {
	return n.rs.renderJSON()
}

func (n *nodeGroupResourceSet) newResource(name string, resource interface{}) string {
	return n.rs.newResource(name, resource)
}

func (n *nodeGroupResourceSet) addResourcesForNodeGroup() {
	lc := &gfn.AWSAutoScalingLaunchConfiguration{
		AssociatePublicIpAddress: true,
		IamInstanceProfile:       n.instanceProfile,
		SecurityGroups:           n.securityGroups,

		ImageId:      n.spec.NodeAMI,
		InstanceType: n.spec.NodeType,
		UserData:     n.userData,
	}
	if n.spec.NodeSSH {
		lc.KeyName = n.spec.SSHPublicKeyName
	}
	refLC := n.newResource("NodeLaunchConfig", lc)

	nodegroup := &gfn.AWSAutoScalingAutoScalingGroup{
		LaunchConfigurationName: refLC,
		DesiredCapacity:         fmt.Sprintf("%d", n.spec.Nodes),
		MinSize:                 fmt.Sprintf("%d", n.spec.MinNodes),
		MaxSize:                 fmt.Sprintf("%d", n.spec.MaxNodes),
		VPCZoneIdentifier:       []string{makeImportValue(ParamClusterStackName, cfnOutputClusterSubnets)},
		Tags: []gfn.AWSAutoScalingAutoScalingGroup_TagProperty{
			gfn.AWSAutoScalingAutoScalingGroup_TagProperty{
				Key:               "Name",
				Value:             gfn.Sub(nodeGroupNameFmt + "-Node"),
				PropagateAtLaunch: true,
			},
			gfn.AWSAutoScalingAutoScalingGroup_TagProperty{
				Key:               gfn.Sub("kubernetes.io/cluster/${ClusterName}"),
				Value:             "owned",
				PropagateAtLaunch: true,
			},
		},
	}

	nodegroup.SetUpdatePolicy(&gfn.UpdatePolicy{
		AutoScalingRollingUpdate: &gfn.AutoScalingRollingUpdate{
			MinInstancesInService: 1,
			MaxBatchSize:          1,
		},
	})

	n.newResource("NodeGroup", nodegroup)

}

func (n *nodeGroupResourceSet) GetAllOutputs(stack cfn.Stack) error {
	return n.rs.GetAllOutputs(stack, n.spec)
}
