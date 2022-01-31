package addon

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"

	"github.com/kris-nova/logger"

	"github.com/aws/aws-sdk-go/service/eks"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func (a *Manager) DeleteWithPreserve(addon *api.Addon) error {
	logger.Info("deleting addon %q and preserving its resources", addon.Name)
	_, err := a.eksAPI.DeleteAddon(&eks.DeleteAddonInput{
		AddonName:   &addon.Name,
		ClusterName: &a.clusterConfig.Metadata.Name,
		Preserve:    aws.Bool(true),
	})

	if err != nil {
		if awsError, ok := err.(awserr.Error); ok && awsError.Code() == eks.ErrCodeResourceNotFoundException {
			logger.Info("addon %q does not exist", addon.Name)
		} else {
			return fmt.Errorf("failed to delete addon %q: %v", addon.Name, err)
		}
	}
	return nil
}

func (a *Manager) Delete(addon *api.Addon) error {
	addonExists := true
	logger.Debug("addon: %v", addon)
	logger.Info("deleting addon: %s", addon.Name)
	_, err := a.eksAPI.DeleteAddon(&eks.DeleteAddonInput{
		AddonName:   &addon.Name,
		ClusterName: &a.clusterConfig.Metadata.Name,
	})

	if err != nil {
		if awsError, ok := err.(awserr.Error); ok && awsError.Code() == eks.ErrCodeResourceNotFoundException {
			logger.Info("addon %q does not exist", addon.Name)
			addonExists = false
		} else {
			return fmt.Errorf("failed to delete addon %q: %v", addon.Name, err)
		}
	} else {
		logger.Info("deleted addon: %s", addon.Name)
	}

	stacks, err := a.stackManager.ListStacksMatching(a.makeAddonName(addon.Name))
	if err != nil {
		return fmt.Errorf("failed to list stacks: %v", err)
	}

	if len(stacks) != 0 {
		logger.Info("deleting associated IAM stacks")
		_, err = a.stackManager.DeleteStackByName(a.makeAddonName(addon.Name))
		if err != nil {
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
