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

func (a *Manager) DescribeAllVersions(ctx context.Context) (string, error) {
	logger.Info("describing all addon versions")
	versions, err := a.describeVersions(ctx, &api.Addon{})
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

	output, err := a.eksAPI.DescribeAddonVersions(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to describe addon versions: %v", err)
	}

	return output, nil
}

func addonVersionsToString(output *eks.DescribeAddonVersionsOutput) (string, error) {
	data, err := json.Marshal(struct {
		Addons []ekstypes.AddonInfo `json:"Addons"`
	}{
		Addons: output.Addons,
	})
	if err != nil {
		return "", err
	}
	return string(data), nil
}
