package addon

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

func (a *Manager) DeleteWithPreserve(ctx context.Context, addon *api.Addon) error {
	logger.Info("deleting addon %q and preserving its resources", addon.Name)
	_, err := a.deleteAddon(ctx, addon, true)
	return err
}

func (a *Manager) Delete(ctx context.Context, addon *api.Addon) error {
	logger.Debug("addon: %v", addon)
	logger.Info("deleting addon: %s", addon.Name)
	addonExists, err := a.deleteAddon(ctx, addon, false)
	if err != nil {
		return err
	}
	if addonExists {
		logger.Info("deleted addon: %s", addon.Name)
	}

	stack, err := a.stackManager.DescribeStack(ctx, &manager.Stack{StackName: aws.String(a.makeAddonName(addon.Name))})
	if err != nil {
		if !manager.IsStackDoesNotExistError(err) {
			return fmt.Errorf("failed to get stack: %w", err)
		}
	}
	if stack != nil {
		logger.Info("deleting associated IAM stacks")
		if _, err = a.stackManager.DeleteStackBySpec(ctx, stack); err != nil {
			return fmt.Errorf("failed to delete cloudformation stack %q: %v", a.makeAddonName(addon.Name), err)
		}
	} else {
		if addonExists {
			logger.Info("no associated IAM stacks found")
		} else {
			return errors.New("could not find addon or associated IAM stack to delete")
		}
	}
	return nil
}

func (a *Manager) deleteAddon(ctx context.Context, addon *api.Addon, preserve bool) (addonExists bool, err error) {
	_, err = a.eksAPI.DeleteAddon(ctx, &eks.DeleteAddonInput{
		AddonName:   &addon.Name,
		ClusterName: &a.clusterConfig.Metadata.Name,
		Preserve:    preserve,
	})

	if err != nil {
		var notFoundErr *ekstypes.ResourceNotFoundException
		if errors.As(err, &notFoundErr) {
			logger.Info("addon %q does not exist", addon.Name)
			return false, nil
		}
		return true, fmt.Errorf("failed to delete addon %q: %v", addon.Name, err)
	}
	return true, nil
}
