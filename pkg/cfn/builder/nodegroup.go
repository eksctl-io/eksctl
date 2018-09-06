package builder

import (
	"fmt"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/awslabs/goformation"
	gfn "github.com/awslabs/goformation/cloudformation"
	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/eks/api"

	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

type nodeGroupResourceSet struct {
	rs               *resourceSet
	id               int
	spec             *api.ClusterConfig
	clusterStackName string
	nodeGroupName    string
	instanceProfile  *gfn.Value
	securityGroups   []*gfn.Value
	vpc              *gfn.Value
	userData         *gfn.Value
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
	n.userData = gfn.NewString(userData)

	if n.spec.MinNodes == 0 && n.spec.MaxNodes == 0 {
		n.spec.MinNodes = n.spec.Nodes
		n.spec.MaxNodes = n.spec.Nodes
	}

	n.addResourcesForIAM()
	n.addResourcesForSecurityGroups()
	n.addResourcesForNodeGroup()

	return nil
}

func (n *nodeGroupResourceSet) AddResourcesForScaling(stackTemplate string) error {
	n.rs.template.Description = nodeGroupTemplateDescription
	n.rs.template.Description += nodeGroupTemplateDescriptionDefaultFeatures
	n.rs.template.Description += templateDescriptionSuffix

	n.rs.newStringParameter(ParamClusterName, "")
	n.rs.newStringParameter(ParamClusterStackName, "")
	n.rs.newNumberParameter(ParamNodeGroupID, "")

	if n.spec.MinNodes == 0 && n.spec.MaxNodes == 0 {
		n.spec.MinNodes = n.spec.Nodes
		n.spec.MaxNodes = n.spec.Nodes
	}

	asg, err := n.getCurrentNodeGroup(stackTemplate)
	if err != nil {
		return err
	}

	n.addResourcesForNodeGroupScaling(asg)

	return nil
}

func (n *nodeGroupResourceSet) RenderJSON() ([]byte, error) {
	return n.rs.renderJSON()
}

func (n *nodeGroupResourceSet) Template() gfn.Template {
	return *n.rs.template
}

func (n *nodeGroupResourceSet) newResource(name string, resource interface{}) *gfn.Value {
	return n.rs.newResource(name, resource)
}

func (n *nodeGroupResourceSet) addResourcesForNodeGroup() {
	lc := &gfn.AWSAutoScalingLaunchConfiguration{
		AssociatePublicIpAddress: gfn.True(),
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
				gfn.FnSplit: []interface{}{
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

func (n *nodeGroupResourceSet) addResourcesForNodeGroupScaling(asg *gfn.AWSAutoScalingAutoScalingGroup) {
	asg.DesiredCapacity = gfn.NewStringRef(fmt.Sprintf("%d", n.spec.Nodes))
	asg.MinSize = gfn.NewStringRef(fmt.Sprintf("%d", n.spec.MinNodes))
	asg.MaxSize = gfn.NewStringRef(fmt.Sprintf("%d", n.spec.MaxNodes))

	n.newResource("NodeGroup", asg)
}

func (n *nodeGroupResourceSet) GetAllOutputs(stack cfn.Stack) error {
	return n.rs.GetAllOutputs(stack, n.spec)
}

func (c *nodeGroupResourceSet) getCurrentNodeGroup(templateBody string) (*gfn.AWSAutoScalingAutoScalingGroup, error) {
	template, err := goformation.ParseYAML([]byte(templateBody))
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse CloudFormation template")
	}

	asg, err := template.GetAWSAutoScalingAutoScalingGroupWithName("NodeGroup")

	if err == nil {
		return nil, fmt.Errorf("Unable to find NodeGroup in existing template")
	}

	return &asg, nil
}
