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
	_, err := a.eksAPI.DeleteAddon(ctx, &eks.DeleteAddonInput{
		AddonName:   &addon.Name,
		ClusterName: &a.clusterConfig.Metadata.Name,
		Preserve:    true,
	})

	if err != nil {
		var notFoundErr *ekstypes.ResourceNotFoundException
		if errors.As(err, &notFoundErr) {
			logger.Info("addon %q does not exist", addon.Name)
		} else {
			return fmt.Errorf("failed to delete addon %q: %v", addon.Name, err)
		}
	}
	return nil
}

func (a *Manager) Delete(ctx context.Context, addon *api.Addon) error {
	addonExists := true
	logger.Debug("addon: %v", addon)
	logger.Info("deleting addon: %s", addon.Name)
	_, err := a.eksAPI.DeleteAddon(ctx, &eks.DeleteAddonInput{
		AddonName:   &addon.Name,
		ClusterName: &a.clusterConfig.Metadata.Name,
	})

	if err != nil {
		var notFoundErr *ekstypes.ResourceNotFoundException
		if errors.As(err, &notFoundErr) {
			logger.Info("addon %q does not exist", addon.Name)
			addonExists = false
		} else {
			return fmt.Errorf("failed to delete addon %q: %v", addon.Name, err)
		}
	} else {
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
