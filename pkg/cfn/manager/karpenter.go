package manager

import (
	"context"
	"fmt"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/karpenter"
	"github.com/weaveworks/eksctl/pkg/karpenter/providers/helm"
)

// KarpenterStack represents the Karpenter stack.
type KarpenterStack struct {
	KarpenterName string
}

// makeNodeGroupStackName generates the name of the Karpenter stack identified by its name, isolated by the cluster this StackCollection operates on
func (c *StackCollection) makeKarpenterStackName() string {
	return fmt.Sprintf("eksctl-%s-karpenter", c.spec.Metadata.Name)
}

// createKarpenterTask creates Karpenter
func (c *StackCollection) createKarpenterTask(errs chan error) error {
	name := c.makeKarpenterStackName()

	logger.Info("building nodegroup stack %q", name)
	stack := builder.NewKarpenterResourceSet(c.iamAPI, c.spec)
	if err := stack.AddAllResources(); err != nil {
		return err
	}
	tags := map[string]string{
		api.KarpenterNameTag: name,
	}
	if err := c.CreateStack(name, stack, tags, nil, errs); err != nil {
		return err
	}
	// Have to create these here, since the Helm Installer returns an error and I
	// don't want to change StackCollection's New*. But this will make
	// testing this function rather difficult.
	helmInstaller, err := helm.NewInstaller(helm.Options{
		Namespace: karpenter.KarpenterNamespace,
	})
	if err != nil {
		return err
	}
	karpenterInstaller := karpenter.NewKarpenterInstaller(karpenter.Options{
		HelmInstaller:         helmInstaller,
		Namespace:             karpenter.KarpenterNamespace,
		ClusterName:           c.spec.Metadata.Name,
		AddDefaultProvisioner: false, // make this configurable
		CreateServiceAccount:  true,  // make this configurable
		ClusterEndpoint:       c.spec.Status.Endpoint,
		Version:               c.spec.Karpenter.Version,
	})
	return karpenterInstaller.InstallKarpenter(context.Background())
}

// DescribeKarpenterStacks calls DescribeStacks and filters out karpenter resources
func (c *StackCollection) DescribeKarpenterStacks() ([]*Stack, error) {
	stacks, err := c.DescribeStacks()
	if err != nil {
		return nil, err
	}

	if len(stacks) == 0 {
		return nil, nil
	}

	var karpenterStacks []*Stack
	for _, s := range stacks {
		switch *s.StackStatus {
		case cfn.StackStatusDeleteComplete:
			continue
		case cfn.StackStatusDeleteFailed:
			logger.Warning("stack's status of karpenter named %s is %s", *s.StackName, *s.StackStatus)
			continue
		}
		if c.GetKarpenterName(s) != "" {
			karpenterStacks = append(karpenterStacks, s)
		}
	}
	logger.Debug("Karpenter stacks = %v", karpenterStacks)
	return karpenterStacks, nil
}

// ListKarpenterStacks returns a list of Karpenter Stacks.
func (c *StackCollection) ListKarpenterStacks() ([]KarpenterStack, error) {
	stacks, err := c.DescribeKarpenterStacks()
	if err != nil {
		return nil, err
	}
	var karpenterStacks []KarpenterStack
	for _, stack := range stacks {
		karpenterStacks = append(karpenterStacks, KarpenterStack{
			KarpenterName: c.GetKarpenterName(stack),
		})
	}
	return karpenterStacks, nil
}

// DescribeKarpenterStacksAndResources calls DescribeKarpenterStacks and fetches all resources,
// then returns it in a map by stack name
func (c *StackCollection) DescribeKarpenterStacksAndResources() (map[string]StackInfo, error) {
	stacks, err := c.DescribeKarpenterStacks()
	if err != nil {
		return nil, err
	}

	allResources := make(map[string]StackInfo)

	for _, s := range stacks {
		input := &cfn.DescribeStackResourcesInput{
			StackName: s.StackName,
		}
		template, err := c.GetStackTemplate(*s.StackName)
		if err != nil {
			return nil, errors.Wrapf(err, "getting template for %q stack", *s.StackName)
		}
		resources, err := c.cloudformationAPI.DescribeStackResources(input)
		if err != nil {
			return nil, errors.Wrapf(err, "getting all resources for %q stack", *s.StackName)
		}
		allResources[c.GetKarpenterName(s)] = StackInfo{
			Resources: resources.StackResources,
			Template:  &template,
			Stack:     s,
		}
	}

	return allResources, nil
}

// GetKarpenterName will return karpenter name based on tags
func (*StackCollection) GetKarpenterName(s *Stack) string {
	if tagName := GetKarpenterTagName(s.Tags); tagName != "" {
		return tagName
	}
	return ""
}

// GetKarpenterTagName returns the Karpenter name of a stack based on its tags.
func GetKarpenterTagName(tags []*cfn.Tag) string {
	for _, tag := range tags {
		switch *tag.Key {
		case api.KarpenterNameTag:
			return *tag.Value
		}
	}
	return ""
}
