package addon

import (
	"fmt"

	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/aws/aws-sdk-go/service/eks"
)

func (a *Manager) DescribeVersions(addon *api.Addon) (string, error) {
	logger.Debug("addon: %v", addon)
	logger.Info("describing addon versions for addon: %s", addon.Name)
	versions, err := a.describeVersions(addon)
	if err != nil {
		return "", err
	}
	return versions.String(), nil
}

func (a *Manager) DescribeAllVersions() (string, error) {
	logger.Info("describing all addon versions")
	versions, err := a.describeVersions(&api.Addon{})
	if err != nil {
		return "", err
	}
	return versions.String(), nil
}

func (a *Manager) describeVersions(addon *api.Addon) (*eks.DescribeAddonVersionsOutput, error) {
	input := &eks.DescribeAddonVersionsInput{
		KubernetesVersion: &a.clusterConfig.Metadata.Version,
	}

	if addon.Name != "" {
		input.AddonName = &addon.Name
	}

	output, err := a.eksAPI.DescribeAddonVersions(input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe addon versions: %v", err)
	}

	return output, nil
}
