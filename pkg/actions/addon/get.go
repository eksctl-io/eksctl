package addon

import (
	"fmt"
	"strings"

	"github.com/blang/semver"

	"github.com/kris-nova/logger"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"

	"github.com/aws/aws-sdk-go/service/eks"
)

type Summary struct {
	Name         string
	Version      string
	NewerVersion string
	IAMRole      string
	Status       string
	Issues       []string
}

func (a *Manager) Get(addon *api.Addon) (Summary, error) {
	logger.Debug("addon: %v", addon)
	output, err := a.eksAPI.DescribeAddon(&eks.DescribeAddonInput{
		ClusterName: &a.clusterConfig.Metadata.Name,
		AddonName:   &addon.Name,
	})

	if err != nil {
		return Summary{}, fmt.Errorf("failed to get addon %q: %v", addon.Name, err)
	}

	var issues []string

	if output.Addon.Health != nil && output.Addon.Health.Issues != nil {
		for _, issue := range output.Addon.Health.Issues {
			issues = append(issues, issue.String())
		}
	}
	serviceAccountRoleARN := ""
	if output.Addon.ServiceAccountRoleArn != nil {
		serviceAccountRoleARN = *output.Addon.ServiceAccountRoleArn
	}

	if addon.Version == "" {
		addon.Version = *output.Addon.AddonVersion
	}

	newerVersion, err := a.findNewerVersions(addon)
	if err != nil {
		return Summary{}, err
	}

	return Summary{
		Name:         *output.Addon.AddonName,
		Version:      *output.Addon.AddonVersion,
		IAMRole:      serviceAccountRoleARN,
		Status:       *output.Addon.Status,
		NewerVersion: newerVersion,
		Issues:       issues,
	}, nil
}

func (a *Manager) GetAll() ([]Summary, error) {
	logger.Info("getting all addons")
	output, err := a.eksAPI.ListAddons(&eks.ListAddonsInput{
		ClusterName: &a.clusterConfig.Metadata.Name,
	})
	if err != nil {
		return []Summary{}, fmt.Errorf("failed to list addons: %v", err)
	}

	var summaries []Summary
	for _, addon := range output.Addons {
		summary, err := a.Get(&api.Addon{Name: *addon})
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, summary)
	}
	return summaries, nil
}

func (a *Manager) findNewerVersions(addon *api.Addon) (string, error) {
	var newerVersions []string
	currentVersion, err := semver.Parse(strings.TrimPrefix(addon.Version, "v"))
	if err != nil {
		logger.Debug("could not parse version %q, skipping finding newer versions: %v", addon.Version, err)
		return "-", nil
	}
	//trim off anything after x.y.z so its not used in comparison, e.g. 1.7.5-eksbuild.1 > 1.7.5
	currentVersion.Build = []string{}
	currentVersion.Pre = []semver.PRVersion{}

	versions, err := a.describeVersions(addon)
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
			//trim off anything after x.y.z and don't use in comparison, e.g. v1.7.5-eksbuild.1 > v1.7.5
			version.Build = []string{}
			version.Pre = []semver.PRVersion{}
			if currentVersion.LT(version) {
				newerVersions = append(newerVersions, *versionInfo.AddonVersion)
			}
		}
	}
	return strings.Join(newerVersions, ","), nil
}
