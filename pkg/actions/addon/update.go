package addon

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/google/uuid"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

func (a *Manager) Update(ctx context.Context, addon *api.Addon, waitTimeout time.Duration) error {
	logger.Debug("addon: %v", addon)

	updateAddonInput := &eks.UpdateAddonInput{
		AddonName:   &addon.Name,
		ClusterName: &a.clusterConfig.Metadata.Name,
	}

	if addon.Force {
		updateAddonInput.ResolveConflicts = ekstypes.ResolveConflictsOverwrite
		logger.Debug("setting resolve conflicts to overwrite")

	}

	summary, err := a.Get(ctx, addon)
	if err != nil {
		return err
	}

	if addon.Version == "" {
		// preserve existing version
		// Might be redundant, does the API care?
		logger.Info("no new version provided, preserving existing version: %s", summary.Version)

		updateAddonInput.AddonVersion = &summary.Version
	} else {
		version, err := a.getLatestMatchingVersion(ctx, addon)
		if err != nil {
			return fmt.Errorf("failed to fetch addon version: %w", err)
		}

		if summary.Version != version {
			logger.Info("new version provided %s", version)
		}

		updateAddonInput.AddonVersion = &version
	}

	//check if we have been provided a different set of policies/role
	if addon.ServiceAccountRoleARN != "" {
		updateAddonInput.ServiceAccountRoleArn = &addon.ServiceAccountRoleARN
	} else if hasPoliciesSet(addon) {
		serviceAccountRoleARN, err := a.updateWithNewPolicies(ctx, addon)
		if err != nil {
			return err
		}
		updateAddonInput.ServiceAccountRoleArn = &serviceAccountRoleARN
	} else {
		//preserve current role
		if summary.IAMRole != "" {
			updateAddonInput.ServiceAccountRoleArn = &summary.IAMRole
		}
	}

	logger.Info("updating addon")
	logger.Debug("%+v", updateAddonInput)

	output, err := a.eksAPI.UpdateAddon(ctx, updateAddonInput)
	if err != nil {
		return fmt.Errorf("failed to update addon %q: %v", addon.Name, err)
	}
	if output != nil {
		logger.Debug("%+v", output.Update)
	}
	if waitTimeout > 0 {
		return a.waitForAddonToBeActive(ctx, addon, waitTimeout)
	}
	return nil
}

func (a *Manager) updateWithNewPolicies(ctx context.Context, addon *api.Addon) (string, error) {
	stackName := a.makeAddonName(addon.Name)
	stack, err := a.stackManager.DescribeStack(ctx, &manager.Stack{StackName: aws.String(stackName)})
	if err != nil {
		if manager.IsStackDoesNotExistError(err) {
			return "", fmt.Errorf("failed to get stack: %w", err)
		}
	}

	namespace, serviceAccount := a.getKnownServiceAccountLocation(addon)

	if stack == nil {
		return a.createRole(ctx, addon, namespace, serviceAccount)
	}

	createNewTemplate, err := a.createNewTemplate(addon, namespace, serviceAccount)
	if err != nil {
		return "", err
	}
	var templateBody manager.TemplateBody = createNewTemplate
	err = a.stackManager.UpdateStack(ctx, manager.UpdateStackOptions{
		Stack:         stack,
		ChangeSetName: fmt.Sprintf("updating-policy-%s", uuid.NewString()),
		Description:   "updating policies",
		TemplateData:  templateBody,
		Wait:          true,
	})
	if err != nil {
		return "", err
	}

	stack, err = a.stackManager.DescribeStack(ctx, &manager.Stack{StackName: aws.String(stackName)})
	if err != nil {
		return "", err
	}
	return *stack.Outputs[0].OutputValue, nil
}

func (a *Manager) createNewTemplate(addon *api.Addon, namespace, serviceAccount string) ([]byte, error) {
	resourceSet, err := a.createRoleResourceSet(addon, namespace, serviceAccount)
	if err != nil {
		return nil, err
	}
	return resourceSet.RenderJSON()
}
