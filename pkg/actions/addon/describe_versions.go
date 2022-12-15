package addon

import (
	"context"
	"encoding/json"
	"fmt"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/aws/aws-sdk-go-v2/service/eks"
)

func (a *Manager) DescribeVersions(ctx context.Context, addon *api.Addon) (string, error) {
	logger.Debug("addon: %v", addon)
	logger.Info("describing addon versions for addon: %s", addon.Name)
	versions, err := a.describeVersions(ctx, addon)
	if err != nil {
		return "", err
	}
	return addonVersionsToString(versions)
}

func (a *Manager) DescribeAllVersions(ctx context.Context, addon *api.Addon) (string, error) {
	logger.Info("describing all addon versions")
	versions, err := a.describeVersions(ctx, addon)
	if err != nil {
		return "", err
	}
	return addonVersionsToString(versions)
}

func (a *Manager) describeVersions(ctx context.Context, addon *api.Addon) (*eks.DescribeAddonVersionsOutput, error) {
	input := &eks.DescribeAddonVersionsInput{
		KubernetesVersion: &a.clusterConfig.Metadata.Version,
	}

	if addon.Name != "" {
		input.AddonName = &addon.Name
	}
	if len(addon.Publishers) != 0 {
		input.Publishers = addon.Publishers
	}
	if len(addon.Types) != 0 {
		input.Types = addon.Types
	}
	if len(addon.Owners) != 0 {
		input.Owners = addon.Owners
	}

	output, err := a.eksAPI.DescribeAddonVersions(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe addon versions: %v", err)
	}

	return output, nil
}

func addonVersionsToString(output *eks.DescribeAddonVersionsOutput) (string, error) {
	data, err := json.MarshalIndent(struct {
		Addons []ekstypes.AddonInfo `json:"Addons"`
	}{
		Addons: output.Addons,
	}, "", "\t")
	if err != nil {
		return "", err
	}
	return string(data), nil
}
