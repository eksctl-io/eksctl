package karpenter

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

const (
	kubernetesTagFormat = "kubernetes.io/cluster/%s"
	karpenterTagFormat  = "karpenter.sh/discovery"
)

type karpenterIAMRolesTask struct {
	info                string
	stackManager        manager.StackManager
	cfg                 *api.ClusterConfig
	instanceProfileName string
	provider            api.ClusterProvider
}

func (k *karpenterIAMRolesTask) Describe() string { return k.info }
func (k *karpenterIAMRolesTask) Do(errs chan error) error {
	return k.createKarpenterIAMRolesTask(errs)
}

// newTasksToInstallKarpenterIAMRoles defines tasks required to create Karpenter IAM roles.
func newTasksToInstallKarpenterIAMRoles(cfg *api.ClusterConfig, stackManager manager.StackManager, provider api.ClusterProvider, instanceProfileName string) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{Parallel: true}
	taskTree.Append(&karpenterIAMRolesTask{
		info:                fmt.Sprintf("create karpenter for stack %q", cfg.Metadata.Name),
		stackManager:        stackManager,
		cfg:                 cfg,
		provider:            provider,
		instanceProfileName: instanceProfileName,
	})
	return taskTree
}

// createKarpenterIAMRolesTask creates Karpenter IAM Roles.
func (k *karpenterIAMRolesTask) createKarpenterIAMRolesTask(errs chan error) error {
	name := k.makeKarpenterStackName()

	logger.Info("building nodegroup stack %q", name)
	stack := builder.NewKarpenterResourceSet(k.cfg, k.instanceProfileName)
	if err := stack.AddAllResources(); err != nil {
		return err
	}
	tags := map[string]string{
		api.KarpenterNameTag:    name,
		api.KarpenterVersionTag: k.cfg.Karpenter.Version,
	}
	if err := k.stackManager.CreateStack(name, stack, tags, nil, errs); err != nil {
		return fmt.Errorf("failed to create stack: %w", err)
	}

	if err := k.ensureSubnetsHaveTags(); err != nil {
		return fmt.Errorf("failed to ensure tags on subnets")
	}

	return k.ensureSecurityGroupKarpenterTag()
}

// makeNodeGroupStackName generates the name of the Karpenter stack identified by its name, isolated by the cluster this StackCollection operates on
func (k *karpenterIAMRolesTask) makeKarpenterStackName() string {
	return fmt.Sprintf("eksctl-%s-karpenter", k.cfg.Metadata.Name)
}

// ensureSubnetsHaveTags will check if the kubernetes.io/cluster tag is present on the subnets.
// if not, it will create them.
func (k *karpenterIAMRolesTask) ensureSubnetsHaveTags() error {
	var ids []string
	for _, subnet := range k.cfg.VPC.Subnets.Private {
		ids = append(ids, subnet.ID)
	}
	for _, subnet := range k.cfg.VPC.Subnets.Public {
		ids = append(ids, subnet.ID)
	}
	sort.Strings(ids)
	input := &ec2.DescribeSubnetsInput{
		SubnetIds: aws.StringSlice(ids),
	}
	output, err := k.provider.EC2().DescribeSubnets(input)
	if err != nil {
		return fmt.Errorf("failed to describe subnets: %w", err)
	}

	clusterTag := fmt.Sprintf(kubernetesTagFormat, k.cfg.Metadata.Name)

	var updateSubnets []string
	for _, subnet := range output.Subnets {
		hasTag := false
		for _, tag := range subnet.Tags {
			if aws.StringValue(tag.Key) == clusterTag {
				hasTag = true
				break
			}
		}
		if !hasTag {
			updateSubnets = append(updateSubnets, *subnet.SubnetId)
		}
	}

	if len(updateSubnets) > 0 {
		creatTagsInput := &ec2.CreateTagsInput{
			Resources: aws.StringSlice(updateSubnets),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String(clusterTag),
					Value: aws.String(""),
				},
				{
					Key:   aws.String(karpenterTagFormat),
					Value: aws.String(k.cfg.Metadata.Name),
				},
			},
		}
		if _, err := k.provider.EC2().CreateTags(creatTagsInput); err != nil {
			return fmt.Errorf("failed to add tags for subnets: %w", err)
		}
	}
	return nil
}

// ensureSecurityGroupKarpenterTag tags all security groups with karpenter.sh/discovery tag.
func (k *karpenterIAMRolesTask) ensureSecurityGroupKarpenterTag() error {
	// Tag the cluster's security group.
	input := &awseks.DescribeClusterInput{
		Name: &k.cfg.Metadata.Name,
	}
	output, err := k.provider.EKS().DescribeCluster(input)
	if err != nil {
		return fmt.Errorf("failed to get cluster: %w", err)
	}
	var ids []string
	for _, id := range output.Cluster.ResourcesVpcConfig.SecurityGroupIds {
		ids = append(ids, aws.StringValue(id))
	}
	ids = append(ids, aws.StringValue(output.Cluster.ResourcesVpcConfig.ClusterSecurityGroupId))

	logger.Info("Attaching tag to the following SGs %q", strings.Join(ids, ", "))

	creatTagsInput := &ec2.CreateTagsInput{
		Resources: aws.StringSlice(ids),
		Tags: []*ec2.Tag{
			{
				Key:   aws.String(karpenterTagFormat),
				Value: aws.String(k.cfg.Metadata.Name),
			},
		},
	}
	if _, err := k.provider.EC2().CreateTags(creatTagsInput); err != nil {
		return fmt.Errorf("failed to add tags for security groups: %w", err)
	}
	return nil
}

// GetKarpenterName will return karpenter name based on tags
func (k *karpenterIAMRolesTask) GetKarpenterName(s *manager.Stack) string {
	return getKarpenterTagName(s.Tags)
}

// getKarpenterTagName returns the Karpenter name of a stack based on its tags.
func getKarpenterTagName(tags []*cfn.Tag) string {
	for _, tag := range tags {
		if *tag.Key == api.KarpenterNameTag {
			return *tag.Value
		}
	}
	return ""
}
