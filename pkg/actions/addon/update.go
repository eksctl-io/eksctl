package addon

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/google/uuid"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

// PodIdentityIAMUpdater creates or updates IAM resources for pod identity associations.
type PodIdentityIAMUpdater interface {
	// UpdateRole creates or updates IAM resources for podIdentityAssociations.
	UpdateRole(ctx context.Context, podIdentityAssociations []api.PodIdentityAssociation, addonName string, existingPodIdentityAssociations []PodIdentityAssociationSummary) ([]ekstypes.AddonPodIdentityAssociations, error)
	// DeleteRole deletes the IAM resources for the specified addon.
	DeleteRole(ctx context.Context, addonName, serviceAccountName string) (bool, error)
}

func (a *Manager) Update(ctx context.Context, addon *api.Addon, podIdentityIAMUpdater PodIdentityIAMUpdater, waitTimeout time.Duration) error {
	logger.Debug("addon: %v", addon)

	var configurationValues *string
	if addon.ConfigurationValues != "" {
		configurationValues = &addon.ConfigurationValues
	}
	updateAddonInput := &eks.UpdateAddonInput{
		AddonName:           &addon.Name,
		ClusterName:         &a.clusterConfig.Metadata.Name,
		ResolveConflicts:    addon.ResolveConflicts,
		ConfigurationValues: configurationValues,
	}

	if addon.Force {
		updateAddonInput.ResolveConflicts = ekstypes.ResolveConflictsOverwrite
	}

	logger.Debug("resolve conflicts set to %s", updateAddonInput.ResolveConflicts)

	summary, err := a.Get(ctx, addon)
	if err != nil {
		return err
	}

	var requiresIAMPermissions bool
	if addon.Version == "" {
		// preserve existing version
		// Might be redundant, does the API care?
		logger.Info("no new version provided, preserving existing version: %s", summary.Version)
		addon.Version = summary.Version
		_, requiresIAMPermissions, err = a.getLatestMatchingVersion(ctx, addon)
		if err != nil {
			var notFoundErr *versionNotFoundError
			if !errors.As(err, &notFoundErr) {
				return fmt.Errorf("failed to fetch addon version %s: %w", summary.Version, err)
			}
		}
		updateAddonInput.AddonVersion = &summary.Version
	} else {
		var latestVersion string
		latestVersion, requiresIAMPermissions, err = a.getLatestMatchingVersion(ctx, addon)
		if err != nil {
			return fmt.Errorf("failed to fetch addon version: %w", err)
		}
		if summary.Version != latestVersion {
			logger.Info("new version provided %s", latestVersion)
		}
		updateAddonInput.AddonVersion = &latestVersion
	}

	var deleteServiceAccountIAMResources []string
	if len(summary.PodIdentityAssociations) > 0 && !addon.UseDefaultPodIdentityAssociations && !a.clusterConfig.AddonsConfig.AutoApplyPodIdentityAssociations {
		if addon.PodIdentityAssociations == nil {
			return fmt.Errorf("addon %s has pod identity associations, to remove pod identity associations from an addon, "+
				"addon.podIdentityAssociations must be explicitly set to []; if the addon was migrated to use pod identity, "+
				"addon.podIdentityAssociations must be set to values obtained from `aws eks describe-pod-identity-association --cluster-name=%s --association-id=%s`",
				addon.Name, a.clusterConfig.Metadata.Name, summary.PodIdentityAssociations[0].AssociationID)
		}
		for _, pia := range summary.PodIdentityAssociations {
			if !slices.ContainsFunc(*addon.PodIdentityAssociations, func(addonPodIdentity api.PodIdentityAssociation) bool {
				return pia.ServiceAccount == addonPodIdentity.ServiceAccountName
			}) {
				deleteServiceAccountIAMResources = append(deleteServiceAccountIAMResources, pia.ServiceAccount)
			}
		}
		// to delete all pod IDs for the addon, explicitly set input.PodIdentityAssociations = []
		if len(*addon.PodIdentityAssociations) == 0 {
			logger.Info("addon.podIdentityAssociations is explicitly set to []; all pod identity associations corresponding to addon %s will be deleted", addon.Name)
			updateAddonInput.PodIdentityAssociations = []ekstypes.AddonPodIdentityAssociations{}
		}
	}

	if addon.HasPodIDsSet() || addon.UseDefaultPodIdentityAssociations {
		if requiresIAMPermissions {
			pidConfigList, supportsPodIdentity, err := a.getRecommendedPoliciesForPodID(ctx, addon)
			if err != nil {
				return fmt.Errorf("getting recommended policies for addon %s", addon.Name)
			}
			if !supportsPodIdentity {
				return &unsupportedPodIdentityErr{addonName: addon.Name}
			}
			var podIdentityAssociations []api.PodIdentityAssociation
			if addon.UseDefaultPodIdentityAssociations {
				for _, pidConfig := range pidConfigList {
					podIdentityAssociations = append(podIdentityAssociations, api.PodIdentityAssociation{
						ServiceAccountName:   *pidConfig.ServiceAccount,
						PermissionPolicyARNs: pidConfig.RecommendedManagedPolicies,
					})
				}
			} else {
				podIdentityAssociations = *addon.PodIdentityAssociations
			}
			addonPodIdentityAssociations, err := podIdentityIAMUpdater.UpdateRole(ctx, podIdentityAssociations, addon.Name, summary.PodIdentityAssociations)
			if err != nil {
				return fmt.Errorf("updating pod identity associations: %w", err)
			}
			updateAddonInput.PodIdentityAssociations = addonPodIdentityAssociations
		} else {
			logger.Warning(IAMPermissionsNotRequiredWarning(addon.Name))
		}
	} else if len(summary.PodIdentityAssociations) > 0 {
		if addon.HasIRSASet() {
			return fmt.Errorf("cannot set IRSA config (`addon.serviceAccountRoleARN`, `addon.attachPolicyARNs`, `addon.attachPolicy`, `addon.wellKnownPolicies`) "+
				"if addon has existing pod identity associations (addon: %s)", addon.Name)
		}
		if !a.clusterConfig.AddonsConfig.AutoApplyPodIdentityAssociations {
			logger.Warning("addon %s has existing pod identity associations but pod identity is not enabled in the config", addon.Name)
		}
	} else {
		// check if we have been provided a different set of policies/role
		if addon.ServiceAccountRoleARN != "" {
			updateAddonInput.ServiceAccountRoleArn = &addon.ServiceAccountRoleARN
		} else if addon.HasIRSAPoliciesSet() {
			serviceAccountRoleARN, err := a.updateWithNewPolicies(ctx, addon)
			if err != nil {
				return err
			}
			updateAddonInput.ServiceAccountRoleArn = &serviceAccountRoleARN
		} else if summary.IAMRole != "" { // Preserve current role.
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
	for _, serviceAccount := range deleteServiceAccountIAMResources {
		logger.Info("deleting IAM resources for pod identity service account %s", serviceAccount)
		deleted, err := podIdentityIAMUpdater.DeleteRole(ctx, addon.Name, serviceAccount)
		if err != nil {
			return fmt.Errorf("deleting IAM resources for addon %s: %w", addon.Name, err)
		}
		if deleted {
			logger.Info("deleted IAM resources for addon %s", addon.Name)
		}
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
		if !manager.IsStackDoesNotExistError(err) {
			return "", fmt.Errorf("failed to get stack: %w", err)
		}
	}

	namespace, serviceAccount := a.getKnownServiceAccountLocation(addon)

	if stack == nil {
		return a.createRoleForIRSA(ctx, addon, namespace, serviceAccount)
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
