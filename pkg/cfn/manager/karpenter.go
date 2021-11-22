package manager

import (
	"context"
	"fmt"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/kris-nova/logger"

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
	return karpenterInstaller.InstallKarpenter(context.Background())
}

// GetKarpenterName will return karpenter name based on tags
func (*StackCollection) GetKarpenterName(s *Stack) string {
	return GetKarpenterTagName(s.Tags)
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
