package addon

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"

	"github.com/blang/semver/v4"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type PodIdentityAssociationSummary struct {
	AssociationID  string
	Namespace      string
	ServiceAccount string
	RoleARN        string
}

type Summary struct {
	Name                    string
	Version                 string
	NewerVersion            string
	IAMRole                 string
	Status                  string
	ConfigurationValues     string
	Issues                  []Issue
	PodIdentityAssociations []PodIdentityAssociationSummary
}

type Issue struct {
	Code        string
	Message     string
	ResourceIDs []string
}

func (a *Manager) Get(ctx context.Context, addon *api.Addon) (Summary, error) {
	logger.Debug("addon: %v", addon)
	output, err := a.eksAPI.DescribeAddon(ctx, &eks.DescribeAddonInput{
		ClusterName: &a.clusterConfig.Metadata.Name,
		AddonName:   &addon.Name,
	})

	if err != nil {
		return Summary{}, fmt.Errorf("failed to get addon %q: %v", addon.Name, err)
	}

	var issues []Issue

	if output.Addon.Health != nil && output.Addon.Health.Issues != nil {
		for _, issue := range output.Addon.Health.Issues {
			issues = append(issues, Issue{
				Code:        string(issue.Code),
				Message:     aws.ToString(issue.Message),
				ResourceIDs: issue.ResourceIds,
			})
		}
	}
	serviceAccountRoleARN := ""
	if output.Addon.ServiceAccountRoleArn != nil {
		serviceAccountRoleARN = *output.Addon.ServiceAccountRoleArn
	}

	addonWithVersion := &api.Addon{
		Name:    addon.Name,
		Version: addon.Version,
	}
	if addonWithVersion.Version == "" {
		addonWithVersion.Version = *output.Addon.AddonVersion
	}

	newerVersion, err := a.findNewerVersions(ctx, addonWithVersion)
	if err != nil {
		return Summary{}, err
	}

	configurationValues := ""
	if output.Addon.ConfigurationValues != nil {
		configurationValues = *output.Addon.ConfigurationValues
	}
	var podIdentityAssociations []PodIdentityAssociationSummary
	podIdentityAssociationIDs, err := toPodIdentityAssociationIDs(output.Addon.PodIdentityAssociations)
	if err != nil {
		return Summary{}, err
	}
	for _, associationID := range podIdentityAssociationIDs {
		output, err := a.eksAPI.DescribePodIdentityAssociation(ctx, &eks.DescribePodIdentityAssociationInput{
			ClusterName:   aws.String(a.clusterConfig.Metadata.Name),
			AssociationId: aws.String(associationID),
		})
		if err != nil {
			return Summary{}, fmt.Errorf("describe pod identity association %q: %w", associationID, err)
		}
		association := output.Association
		podIdentityAssociations = append(podIdentityAssociations, PodIdentityAssociationSummary{
			Namespace:      *association.Namespace,
			ServiceAccount: *association.ServiceAccount,
			RoleARN:        *association.RoleArn,
			AssociationID:  *association.AssociationId,
		})
	}

	return Summary{
		Name:                    *output.Addon.AddonName,
		Version:                 *output.Addon.AddonVersion,
		IAMRole:                 serviceAccountRoleARN,
		Status:                  string(output.Addon.Status),
		NewerVersion:            newerVersion,
		ConfigurationValues:     configurationValues,
		PodIdentityAssociations: podIdentityAssociations,
		Issues:                  issues,
	}, nil
}

func (a *Manager) GetAll(ctx context.Context) ([]Summary, error) {
	logger.Info("getting all addons")
	output, err := a.eksAPI.ListAddons(ctx, &eks.ListAddonsInput{
		ClusterName: &a.clusterConfig.Metadata.Name,
	})
	if err != nil {
		return []Summary{}, fmt.Errorf("failed to list addons: %v", err)
	}

	var summaries []Summary
	for _, addon := range output.Addons {
		summary, err := a.Get(ctx, &api.Addon{Name: addon})
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}

func toPodIdentityAssociationIDs(podIdentityAssociationARNs []string) ([]string, error) {
	var piaIDs []string
	for _, podIdentityAssociationARN := range podIdentityAssociationARNs {
		piaID, err := api.ToPodIdentityAssociationID(podIdentityAssociationARN)
		if err != nil {
			return nil, err
		}
		piaIDs = append(piaIDs, piaID)
	}
	return piaIDs, nil
}

func (a *Manager) findNewerVersions(ctx context.Context, addon *api.Addon) (string, error) {
	var newerVersions []string
	currentVersion, err := semver.Parse(strings.TrimPrefix(addon.Version, "v"))
	if err != nil {
		logger.Debug("could not parse version %q, skipping finding newer versions: %v", addon.Version, err)
		return "-", nil
	}

	versions, err := a.describeVersions(ctx, addon)
	if err != nil {
		return "", err
	}

	if len(versions.Addons) == 0 {
		return "-", nil
	}

	for _, versionInfo := range versions.Addons[0].AddonVersions {
		version, err := semver.Parse(strings.TrimPrefix(*versionInfo.AddonVersion, "v"))
		if err != nil {
			logger.Debug("could not parse version %q, skipping version comparison: %v", addon.Version, err)
		} else {
			if currentVersion.LT(version) {
				newerVersions = append(newerVersions, *versionInfo.AddonVersion)
			}
		}
	}
	return strings.Join(newerVersions, ","), nil
}
