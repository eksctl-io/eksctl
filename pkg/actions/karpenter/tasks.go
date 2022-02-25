package karpenter

import (
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go/aws"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

const (
	kubernetesTagFormat = "kubernetes.io/cluster/%s"
)

type karpenterIAMRolesTask struct {
	info                string
	stackManager        manager.StackManager
	cfg                 *api.ClusterConfig
	ec2API              ec2iface.EC2API
	instanceProfileName string
}

func (k *karpenterIAMRolesTask) Describe() string { return k.info }
func (k *karpenterIAMRolesTask) Do(errs chan error) error {
	return k.createKarpenterIAMRolesTask(errs)
}

// newTasksToInstallKarpenterIAMRoles defines tasks required to create Karpenter IAM roles.
func newTasksToInstallKarpenterIAMRoles(cfg *api.ClusterConfig, stackManager manager.StackManager, ec2API ec2iface.EC2API, instanceProfileName string) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{Parallel: true}
	taskTree.Append(&karpenterIAMRolesTask{
		info:                fmt.Sprintf("create karpenter for stack %q", cfg.Metadata.Name),
		stackManager:        stackManager,
		cfg:                 cfg,
		ec2API:              ec2API,
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

	return k.ensureSubnetsHaveTags()
}

// makeNodeGroupStackName generates the name of the Karpenter stack identified by its name, isolated by the cluster this StackCollection operates on
func (k *karpenterIAMRolesTask) makeKarpenterStackName() string {
	return fmt.Sprintf("eksctl-%s-karpenter", k.cfg.Metadata.Name)
}

// ensureSubnetsHaveTags sets of overwrites kubernetes.io/cluster/<name> tags on subnets with the current value.
func (k *karpenterIAMRolesTask) ensureSubnetsHaveTags() error {
	var ids []string
	for _, subnet := range k.cfg.VPC.Subnets.Private {
		ids = append(ids, subnet.ID)
	}
	for _, subnet := range k.cfg.VPC.Subnets.Public {
		ids = append(ids, subnet.ID)
	}
	sort.Strings(ids)
	clusterTag := fmt.Sprintf(kubernetesTagFormat, k.cfg.Metadata.Name)
	creatTagsInput := &ec2.CreateTagsInput{
		Resources: aws.StringSlice(ids),
		Tags: []*ec2.Tag{
			{
				Key:   aws.String(clusterTag),
				Value: aws.String(""),
			},
		},
	}
	if _, err := k.ec2API.CreateTags(creatTagsInput); err != nil {
		return fmt.Errorf("failed to add tags for subnets: %w", err)
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
