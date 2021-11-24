package manager

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/karpenter"
	"github.com/weaveworks/eksctl/pkg/karpenter/providers/helm"
)

const (
	kubernetesTagFormat = "kubernetes.io/cluster/%s"
)

// makeNodeGroupStackName generates the name of the Karpenter stack identified by its name, isolated by the cluster this StackCollection operates on
func (c *StackCollection) makeKarpenterStackName() string {
	return fmt.Sprintf("eksctl-%s-karpenter", c.spec.Metadata.Name)
}

// createKarpenterTask creates Karpenter
func (c *StackCollection) createKarpenterTask(errs chan error) error {
	name := c.makeKarpenterStackName()

	logger.Info("building nodegroup stack %q", name)
	stack := builder.NewKarpenterResourceSet(c.spec)
	if err := stack.AddAllResources(); err != nil {
		return err
	}
	tags := map[string]string{
		api.KarpenterNameTag: name,
	}
	if err := c.CreateStack(name, stack, tags, nil, errs); err != nil {
		return err
	}

	if err := c.maybeUpdateSubnetTags(); err != nil {
		return err
	}

	helmInstaller, err := helm.NewInstaller(helm.Options{
		Namespace: karpenter.DefaultKarpenterNamespace,
	})
	if err != nil {
		return err
	}
	karpenterInstaller := karpenter.NewKarpenterInstaller(karpenter.Options{
		HelmInstaller:         helmInstaller,
		Namespace:             karpenter.DefaultKarpenterNamespace,
		ClusterName:           c.spec.Metadata.Name,
		AddDefaultProvisioner: api.IsEnabled(c.spec.Karpenter.AddDefaultProvisioner),
		CreateServiceAccount:  api.IsEnabled(c.spec.Karpenter.CreateServiceAccount),
		ClusterEndpoint:       c.spec.Status.Endpoint,
		Version:               c.spec.Karpenter.Version,
	})
	return karpenterInstaller.Install(context.Background())
}

// maybeUpdateSubnetTags will check if the kubernetes.io/cluster tag is present on the subnets.
// if not, it will create them.
func (c *StackCollection) maybeUpdateSubnetTags() error {
	var ids []string
	for _, subnet := range c.spec.VPC.Subnets.Private {
		ids = append(ids, subnet.ID)
	}
	for _, subnet := range c.spec.VPC.Subnets.Public {
		ids = append(ids, subnet.ID)
	}
	input := &ec2.DescribeSubnetsInput{
		SubnetIds: aws.StringSlice(ids),
	}
	output, err := c.ec2API.DescribeSubnets(input)
	if err != nil {
		return fmt.Errorf("failed to describe subnets: %w", err)
	}

	clusterTag := fmt.Sprintf(kubernetesTagFormat, c.spec.Metadata.Name)

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
		if _, err := c.ec2API.CreateTags(&ec2.CreateTagsInput{
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
func (*StackCollection) GetKarpenterName(s *Stack) string {
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
