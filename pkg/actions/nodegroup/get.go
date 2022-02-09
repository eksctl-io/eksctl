package nodegroup

import (
	"strconv"
	"strings"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
)

func (m *Manager) GetAll() ([]*manager.NodeGroupSummary, error) {
	summaries, err := m.stackManager.GetUnmanagedNodeGroupSummaries("")
	if err != nil {
		return nil, errors.Wrap(err, "getting nodegroup stack summaries")
	}

	for _, summary := range summaries {
		if summary.DesiredCapacity > 0 {
			summary.Version, err = kubewrapper.GetNodegroupKubernetesVersion(m.clientSet.CoreV1().Nodes(), summary.Name)
			if err != nil {
				return nil, errors.Wrap(err, "getting nodegroup's kubernetes version")
			}
		}
	}

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

		summary, err := m.makeManagedNGSummary(*ngName)
		if err != nil {
			return nil, err
		}
		summary.StackName = aws.StringValue(stack.StackName)
		summaries = append(summaries, summary)
	}

	return summaries, nil
}

func (m *Manager) Get(name string) (*manager.NodeGroupSummary, error) {
	summaries, err := m.stackManager.GetUnmanagedNodeGroupSummaries(name)
	if err != nil {
		return nil, errors.Wrap(err, "getting nodegroup stack summaries")
	}

	if len(summaries) > 0 {
		s := summaries[0]
		if s.DesiredCapacity > 0 {
			s.Version, err = kubewrapper.GetNodegroupKubernetesVersion(m.clientSet.CoreV1().Nodes(), s.Name)
			if err != nil {
				return nil, errors.Wrap(err, "getting nodegroup's kubernetes version")
			}
		}
		return s, nil
	}

	return m.makeManagedNGSummary(name)
}

func (m *Manager) makeManagedNGSummary(nodeGroupName string) (*manager.NodeGroupSummary, error) {
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

	return &manager.NodeGroupSummary{
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
