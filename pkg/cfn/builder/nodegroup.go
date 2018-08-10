package builder

import (
	"fmt"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	gfn "github.com/awslabs/goformation/cloudformation"

	"github.com/weaveworks/eksctl/pkg/eks/api"

	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

var (
	clusterOwnedTag = gfn.Tag{
		Key:   makeSub("kubernetes.io/cluster/${ClusterName}"),
		Value: gfn.NewString("owned"),
	}
)

type nodeGroupResourceSet struct {
	rs               *resourceSet
	id               int
	spec             *api.ClusterConfig
	clusterStackName string
	nodeGroupName    string
	instanceProfile  *gfn.StringIntrinsic
	securityGroups   []*gfn.StringIntrinsic
	vpc              *gfn.StringIntrinsic
	userData         *gfn.StringIntrinsic
	clusterOwnedTag  gfn.Tag
}

type awsCloudFormationResource struct {
	Type         string
	Properties   map[string]interface{}
	UpdatePolicy map[string]map[string]string
}

func NewNodeGroupResourceSet(spec *api.ClusterConfig, clusterStackName string, id int) *nodeGroupResourceSet {
	return &nodeGroupResourceSet{
		rs:               newResourceSet(),
		id:               id,
		clusterStackName: clusterStackName,
		nodeGroupName:    fmt.Sprintf("%s-%d", spec.ClusterName, id),
		spec:             spec,
	}
}

func (n *nodeGroupResourceSet) AddAllResources() error {
	n.rs.template.Description = nodeGroupTemplateDescription
	n.rs.template.Description += nodeGroupTemplateDescriptionDefaultFeatures
	n.rs.template.Description += templateDescriptionSuffix

	n.vpc = makeImportValue(n.clusterStackName, cfnOutputClusterVPC)

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

func (n *nodeGroupResourceSet) newResource(name string, resource interface{}) *gfn.StringIntrinsic {
	return n.rs.newResource(name, resource)
}

func (n *nodeGroupResourceSet) addResourcesForNodeGroup() {
	lc := &gfn.AWSAutoScalingLaunchConfiguration{
		AssociatePublicIpAddress: true,
		IamInstanceProfile:       n.instanceProfile,
		SecurityGroups:           n.securityGroups,

		ImageId:      gfn.NewString(n.spec.NodeAMI),
		InstanceType: gfn.NewString(n.spec.NodeType),
		UserData:     n.userData,
	}
	if n.spec.NodeSSH {
		lc.KeyName = gfn.NewString(n.spec.SSHPublicKeyName)
	}
	refLC := n.newResource("NodeLaunchConfig", lc)
	// currently goformation type system doesn't allow specifying `VPCZoneIdentifier: { "Fn::ImportValue": ... }`,
	// and tags don't have `PropagateAtLaunch` field, so we have a custom method here until this gets resolved
	n.newResource("NodeGroup", &awsCloudFormationResource{
		Type: "AWS::AutoScaling::AutoScalingGroup",
		Properties: map[string]interface{}{
			"LaunchConfigurationName": refLC,
			"DesiredCapacity":         fmt.Sprintf("%d", n.spec.Nodes),
			"MinSize":                 fmt.Sprintf("%d", n.spec.MinNodes),
			"MaxSize":                 fmt.Sprintf("%d", n.spec.MaxNodes),
			"VPCZoneIdentifier": map[string][]interface{}{
				fnSplit: []interface{}{
					",",
					makeImportValue(n.clusterStackName, cfnOutputClusterSubnets),
				},
			},
			"Tags": []map[string]interface{}{
				{"Key": "Name", "Value": fmt.Sprintf("%s-Node", n.nodeGroupName), "PropagateAtLaunch": "true"},
				{"Key": "kubernetes.io/cluster/" + n.spec.ClusterName, "Value": "owned", "PropagateAtLaunch": "true"},
			},
		},
		UpdatePolicy: map[string]map[string]string{
			"AutoScalingRollingUpdate": {
				"MinInstancesInService": "1",
				"MaxBatchSize":          "1",
			},
		},
	})
}

func (n *nodeGroupResourceSet) GetAllOutputs(stack cfn.Stack) error {
	return n.rs.GetAllOutputs(stack, n.spec)
}
