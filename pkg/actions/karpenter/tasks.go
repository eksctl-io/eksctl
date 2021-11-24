package karpenter

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/karpenter"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

const (
	kubernetesTagFormat = "kubernetes.io/cluster/%s"
)

type karpenterTask struct {
	info               string
	stackManager       manager.StackManager
	cfg                *api.ClusterConfig
	ec2API             ec2iface.EC2API
	karpenterInstaller karpenter.InstallKarpenter
}

func (k *karpenterTask) Describe() string { return k.info }
func (k *karpenterTask) Do(errs chan error) error {
	return k.createKarpenterTask(errs)
}

// NewTasksToInstallKarpenter defines tasks required to create Karpenter
func NewTasksToInstallKarpenter(cfg *api.ClusterConfig, stackManager manager.StackManager, ec2API ec2iface.EC2API, karpenterInstaller karpenter.InstallKarpenter) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{Parallel: true}
	taskTree.Append(&karpenterTask{
		info:               fmt.Sprintf("create karpenter for stack %q", cfg.Metadata.Name),
		stackManager:       stackManager,
		cfg:                cfg,
		ec2API:             ec2API,
		karpenterInstaller: karpenterInstaller,
	})
	return taskTree
}

// createKarpenterTask creates Karpenter
func (k *karpenterTask) createKarpenterTask(errs chan error) error {
	name := k.makeKarpenterStackName()

	logger.Info("building nodegroup stack %q", name)
	stack := builder.NewKarpenterResourceSet(k.cfg)
	if err := stack.AddAllResources(); err != nil {
		return err
	}
	tags := map[string]string{
		api.KarpenterNameTag: name,
	}
	if err := k.stackManager.CreateStack(name, stack, tags, nil, errs); err != nil {
		return err
	}

	if err := k.maybeUpdateSubnetTags(); err != nil {
		return err
	}
	return k.karpenterInstaller.Install(context.Background())
}

// makeNodeGroupStackName generates the name of the Karpenter stack identified by its name, isolated by the cluster this StackCollection operates on
func (k *karpenterTask) makeKarpenterStackName() string {
	return fmt.Sprintf("eksctl-%s-karpenter", k.cfg.Metadata.Name)
}

// maybeUpdateSubnetTags will check if the kubernetes.io/cluster tag is present on the subnets.
// if not, it will create them.
func (k *karpenterTask) maybeUpdateSubnetTags() error {
	var ids []string
	for _, subnet := range k.cfg.VPC.Subnets.Private {
		ids = append(ids, subnet.ID)
	}
	for _, subnet := range k.cfg.VPC.Subnets.Public {
		ids = append(ids, subnet.ID)
	}
	input := &ec2.DescribeSubnetsInput{
		SubnetIds: aws.StringSlice(ids),
	}
	output, err := k.ec2API.DescribeSubnets(input)
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
		if _, err := k.ec2API.CreateTags(&ec2.CreateTagsInput{
			Resources: aws.StringSlice(ids),
			Tags: []*ec2.Tag{
				{
					Key:   aws.String(clusterTag),
					Value: aws.String(""),
				},
			},
		}); err != nil {
			return fmt.Errorf("failed to add tags for subnets: %w", err)
		}
	}
	return nil
}

// GetKarpenterName will return karpenter name based on tags
func (k *karpenterTask) GetKarpenterName(s *manager.Stack) string {
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
