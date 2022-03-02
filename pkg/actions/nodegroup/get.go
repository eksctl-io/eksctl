package nodegroup

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/tidwall/gjson"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
)

const (
	imageIDPath       = "Resources.NodeGroupLaunchTemplate.Properties.LaunchTemplateData.ImageId"
	resourcesRootPath = "Resources"
)

// Summary represents a summary of a nodegroup stack
type Summary struct {
	StackName            string
	Cluster              string
	Name                 string
	Status               string
	MaxSize              int
	MinSize              int
	DesiredCapacity      int
	InstanceType         string
	ImageID              string
	CreationTime         time.Time
	NodeInstanceRoleARN  string
	AutoScalingGroupName string
	Version              string
	NodeGroupType        api.NodeGroupType `json:"Type"`
}

func (m *Manager) GetAll() ([]*Summary, error) {
	unmanagedSummaries, err := m.getUnmanagedSummaries()
	if err != nil {
		return nil, err
	}

	managedSummaries, err := m.getManagedSummaries()
	if err != nil {
		return nil, err
	}

	return append(unmanagedSummaries, managedSummaries...), nil
}

func (m *Manager) Get(name string) (*Summary, error) {
	summary, err := m.getUnmanagedSummary(name)
	if err != nil {
		return nil, fmt.Errorf("getting nodegroup stack summaries: %w", err)
	}

	if summary != nil {
		return summary, nil
	}

	return m.getManagedSummary(name)
}

func (m *Manager) getManagedSummaries() ([]*Summary, error) {
	var summaries []*Summary
	managedNodeGroups, err := m.ctl.Provider.EKS().ListNodegroups(&eks.ListNodegroupsInput{
		ClusterName: aws.String(m.cfg.Metadata.Name),
	})
	if err != nil {
		return nil, err
	}

	for _, ngName := range managedNodeGroups.Nodegroups {
		var stack *cloudformation.Stack
		stack, err = m.stackManager.DescribeNodeGroupStack(*ngName)
		if err != nil {
			stack = &cloudformation.Stack{}
		}

		summary, err := m.getManagedSummary(*ngName)
		if err != nil {
			return nil, err
		}
		summary.StackName = aws.StringValue(stack.StackName)
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func (m *Manager) getUnmanagedSummaries() ([]*Summary, error) {
	stacks, err := m.stackManager.DescribeNodeGroupStacks()
	if err != nil {
		return nil, fmt.Errorf("getting nodegroup stacks: %w", err)
	}

	// Create an empty array here so that an object is returned rather than null
	summaries := make([]*Summary, 0)
	for _, s := range stacks {
		summary, err := m.unmanagedStackToSummary(s)
		if err != nil {
			return nil, err
		}
		if summary != nil {
			summaries = append(summaries, summary)
		}
	}

	return summaries, nil
}

func (m *Manager) getUnmanagedSummary(name string) (*Summary, error) {
	stack, err := m.stackManager.DescribeNodeGroupStack(name)
	if err != nil {
		return nil, err
	}

	return m.unmanagedStackToSummary(stack)
}

func (m *Manager) unmanagedStackToSummary(s *manager.Stack) (*Summary, error) {
	nodeGroupType, err := manager.GetNodeGroupType(s.Tags)
	if err != nil {
		return nil, err
	}

	if nodeGroupType != api.NodeGroupTypeUnmanaged {
		return nil, nil
	}

	ngPaths, err := getNodeGroupPaths(s.Tags)
	if err != nil {
		return nil, err
	}

	summary, err := m.mapStackToNodeGroupSummary(s, ngPaths)

	if err != nil {
		return nil, fmt.Errorf("mapping stack to nodegroup summary: %w", err)
	}
	summary.NodeGroupType = api.NodeGroupTypeUnmanaged

	asgName, err := m.stackManager.GetUnmanagedNodeGroupAutoScalingGroupName(s)
	if err != nil {
		return nil, fmt.Errorf("getting autoscalinggroupname: %w", err)
	}

	summary.AutoScalingGroupName = asgName

	scalingGroup, err := m.stackManager.GetAutoScalingGroupDesiredCapacity(asgName)
	if err != nil {
		return nil, fmt.Errorf("getting autoscalinggroup desired capacity: %w", err)
	}
	summary.DesiredCapacity = int(*scalingGroup.DesiredCapacity)
	summary.MinSize = int(*scalingGroup.MinSize)
	summary.MaxSize = int(*scalingGroup.MaxSize)

	if summary.DesiredCapacity > 0 {
		summary.Version, err = kubewrapper.GetNodegroupKubernetesVersion(m.clientSet.CoreV1().Nodes(), summary.Name)
		if err != nil {
			return nil, fmt.Errorf("getting nodegroup's kubernetes version: %w", err)
		}
	}

	return summary, nil
}

func getNodeGroupPaths(tags []*cfn.Tag) (*nodeGroupPaths, error) {
	nodeGroupType, err := manager.GetNodeGroupType(tags)
	if err != nil {
		return nil, err
	}

	switch nodeGroupType {
	case api.NodeGroupTypeManaged:
		makePath := func(fieldPath string) string {
			return fmt.Sprintf("%s.ManagedNodeGroup.Properties.%s", resourcesRootPath, fieldPath)
		}
		makeScalingPath := func(field string) string {
			return makePath(fmt.Sprintf("ScalingConfig.%s", field))
		}
		return &nodeGroupPaths{
			InstanceType:    makePath("InstanceTypes.0"),
			DesiredCapacity: makeScalingPath("DesiredSize"),
			MinSize:         makeScalingPath("MinSize"),
			MaxSize:         makeScalingPath("MaxSize"),
		}, nil

		// Tag may not exist for existing nodegroups
	case api.NodeGroupTypeUnmanaged, "":
		makePath := func(field string) string {
			return fmt.Sprintf("%s.NodeGroup.Properties.%s", resourcesRootPath, field)
		}
		return &nodeGroupPaths{
			InstanceType:    resourcesRootPath + ".NodeGroupLaunchTemplate.Properties.LaunchTemplateData.InstanceType",
			DesiredCapacity: makePath("DesiredCapacity"),
			MinSize:         makePath("MinSize"),
			MaxSize:         makePath("MaxSize"),
		}, nil

	default:
		return nil, fmt.Errorf("unexpected nodegroup type tag: %q", nodeGroupType)
	}

}

type nodeGroupPaths struct {
	InstanceType    string
	DesiredCapacity string
	MinSize         string
	MaxSize         string
}

func (m *Manager) mapStackToNodeGroupSummary(stack *manager.Stack, ngPaths *nodeGroupPaths) (*Summary, error) {
	template, err := m.stackManager.GetStackTemplate(*stack.StackName)
	if err != nil {
		return nil, fmt.Errorf("error getting CloudFormation template for stack %s: %w", *stack.StackName, err)
	}

	summary := &Summary{
		StackName:       *stack.StackName,
		Cluster:         getClusterNameTag(stack),
		Name:            m.stackManager.GetNodeGroupName(stack),
		Status:          *stack.StackStatus,
		MaxSize:         int(gjson.Get(template, ngPaths.MaxSize).Int()),
		MinSize:         int(gjson.Get(template, ngPaths.MinSize).Int()),
		DesiredCapacity: int(gjson.Get(template, ngPaths.DesiredCapacity).Int()),
		InstanceType:    gjson.Get(template, ngPaths.InstanceType).String(),
		ImageID:         gjson.Get(template, imageIDPath).String(),
		CreationTime:    *stack.CreationTime,
	}

	nodeGroupType, err := manager.GetNodeGroupType(stack.Tags)
	if err != nil {
		return nil, err
	}

	var nodeInstanceRoleARN string
	if nodeGroupType == api.NodeGroupTypeUnmanaged {
		nodeInstanceRoleARNCollector := func(s string) error {
			nodeInstanceRoleARN = s
			return nil
		}
		collectors := map[string]outputs.Collector{
			outputs.NodeGroupInstanceRoleARN: nodeInstanceRoleARNCollector,
		}
		collectorSet := outputs.NewCollectorSet(collectors)
		if err := collectorSet.MustCollect(*stack); err != nil {
			logger.Warning(fmt.Errorf("error collecting Cloudformation outputs for stack %s: %w", *stack.StackName, err).Error())
		}
	}

	summary.NodeInstanceRoleARN = nodeInstanceRoleARN

	return summary, nil
}

func getClusterNameTag(s *manager.Stack) string {
	for _, tag := range s.Tags {
		if *tag.Key == api.ClusterNameTag || *tag.Key == api.OldClusterNameTag {
			return *tag.Value
		}
	}
	return ""
}

func (m *Manager) getManagedSummary(nodeGroupName string) (*Summary, error) {
	describeOutput, err := m.ctl.Provider.EKS().DescribeNodegroup(&eks.DescribeNodegroupInput{
		ClusterName:   aws.String(m.cfg.Metadata.Name),
		NodegroupName: aws.String(nodeGroupName),
	})

	if err != nil {
		return nil, err
	}

	ng := describeOutput.Nodegroup

	var asgs []string
	if ng.Resources != nil {
		for _, asg := range ng.Resources.AutoScalingGroups {
			asgs = append(asgs, aws.StringValue(asg.Name))
		}
	}

	var imageID string
	if *ng.AmiType == awseks.AMITypesCustom {
		// ReleaseVersion contains the AMI ID for custom AMIs.
		imageID = *ng.ReleaseVersion
	} else {
		imageID = *ng.AmiType
	}

	return &Summary{
		Name:                 *ng.NodegroupName,
		Cluster:              *ng.ClusterName,
		Status:               *ng.Status,
		MaxSize:              int(*ng.ScalingConfig.MaxSize),
		MinSize:              int(*ng.ScalingConfig.MinSize),
		DesiredCapacity:      int(*ng.ScalingConfig.DesiredSize),
		InstanceType:         m.getInstanceTypes(ng),
		ImageID:              imageID,
		CreationTime:         *ng.CreatedAt,
		NodeInstanceRoleARN:  *ng.NodeRole,
		AutoScalingGroupName: strings.Join(asgs, ","),
		Version:              getOptionalValue(ng.Version),
		NodeGroupType:        api.NodeGroupTypeManaged,
	}, nil
}

func (m *Manager) getInstanceTypes(ng *awseks.Nodegroup) string {
	if len(ng.InstanceTypes) > 0 {
		return strings.Join(aws.StringValueSlice(ng.InstanceTypes), ",")
	}

	if ng.LaunchTemplate == nil {
		logger.Info("no instance type found for nodegroup %q", *ng.NodegroupName)
		return "-"
	}

	resp, err := m.ctl.Provider.EC2().DescribeLaunchTemplateVersions(&ec2.DescribeLaunchTemplateVersionsInput{
		LaunchTemplateId: ng.LaunchTemplate.Id,
	})
	if err != nil {
		return "-"
	}

	for _, template := range resp.LaunchTemplateVersions {
		if strconv.Itoa(int(*template.VersionNumber)) == *ng.LaunchTemplate.Version {
			return *template.LaunchTemplateData.InstanceType
		}
	}

	logger.Info("no instance type found for nodegroup %q", *ng.NodegroupName)
	return "-"
}

func getOptionalValue(v *string) string {
	if v == nil {
		return "-"
	}
	return *v
}
