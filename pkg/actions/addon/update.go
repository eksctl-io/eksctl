package addon

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

func (a *Manager) Update(addon *api.Addon, wait bool) error {
	logger.Debug("addon: %v", addon)

	updateAddonInput := &eks.UpdateAddonInput{
		AddonName:   &addon.Name,
		ClusterName: &a.clusterConfig.Metadata.Name,
		//ResolveConflicts: 		"enum":["OVERWRITE","NONE"]
	}

	if addon.Force {
		updateAddonInput.ResolveConflicts = aws.String("overwrite")
		logger.Debug("setting resolve conflicts to overwrite")

	}

	summary, err := a.Get(addon)
	if err != nil {
		return err
	}

	if addon.Version == "" {
		// preserve existing version
		// Might be redundant, does the API care?
		logger.Info("no new version provided, preserving existing version: %s", summary.Version)

		updateAddonInput.AddonVersion = &summary.Version
	} else {
		version, err := a.getLatestMatchingVersion(addon)
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
		serviceAccountRoleARN, err := a.updateWithNewPolicies(addon)
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
	logger.Debug(updateAddonInput.String())

	output, err := a.eksAPI.UpdateAddon(updateAddonInput)
	if err != nil {
		return fmt.Errorf("failed to update addon %q: %v", addon.Name, err)
	}
	if output != nil {
		logger.Debug(output.String())
	}
	if wait {
		return a.waitForAddonToBeActive(addon)
	}
	return nil
}

func (a *Manager) updateWithNewPolicies(addon *api.Addon) (string, error) {
	stackName := a.makeAddonName(addon.Name)
	existingStacks, err := a.stackManager.ListStacksMatching(stackName)
	if err != nil {
		return "", err
	}

	namespace, serviceAccount := a.getKnownServiceAccountLocation(addon)

	if len(existingStacks) == 0 {
		return a.createRole(addon, namespace, serviceAccount)
	}

	createNewTemplate, err := a.createNewTemplate(addon, namespace, serviceAccount)
	if err != nil {
		return "", err
	}
	var templateBody manager.TemplateBody = createNewTemplate
	err = a.stackManager.UpdateStack(stackName, fmt.Sprintf("updating-policy-%s", uuid.NewString()), "updating policies", templateBody, nil)
	if err != nil {
		return "", err
	}

	existingStacks, err = a.stackManager.ListStacksMatching(stackName)
	if err != nil {
		return "", err
	}

	return *existingStacks[0].Outputs[0].OutputValue, nil
}

func (a *Manager) createNewTemplate(addon *api.Addon, namespace, serviceAccount string) ([]byte, error) {
	resourceSet, err := a.createRoleResourceSet(addon, namespace, serviceAccount)
	if err != nil {
		return nil, err
	}
	return resourceSet.RenderJSON()
}
