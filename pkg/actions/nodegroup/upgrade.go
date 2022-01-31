package nodegroup

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/blang/semver"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/goformation/v4"
	"github.com/weaveworks/goformation/v4/cloudformation"

	"github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/managed"
	"github.com/weaveworks/eksctl/pkg/utils/waiters"
	"github.com/weaveworks/eksctl/pkg/version"
	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfneks "github.com/weaveworks/goformation/v4/cloudformation/eks"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

// UpgradeOptions contains options to configure nodegroup upgrades
type UpgradeOptions struct {
	// NodeGroupName nodegroup name
	NodegroupName string
	// KubernetesVersion EKS version
	KubernetesVersion string
	// LaunchTemplateVersion launch template version
	// valid only if a nodegroup was created with a launch template
	LaunchTemplateVersion string
	//ForceUpgrade enables force upgrade
	ForceUpgrade bool
	// ReleaseVersion AMI version of the EKS optimized AMI to use
	ReleaseVersion string
	// Wait for the upgrade to finish
	Wait bool
}

func (m *Manager) Upgrade(options UpgradeOptions) error {
	hasStacks, err := m.hasStacks(options.NodegroupName)
	if err != nil {
		return err
	}

	if options.KubernetesVersion != "" {
		if _, err := semver.ParseTolerant(options.KubernetesVersion); err != nil {
			return errors.Wrap(err, "invalid Kubernetes version")
		}
	}

	nodegroupOutput, err := m.ctl.Provider.EKS().DescribeNodegroup(&eks.DescribeNodegroupInput{
		ClusterName:   &m.cfg.Metadata.Name,
		NodegroupName: &options.NodegroupName,
	})

	if err != nil {
		if managed.IsNotFound(err) {
			return fmt.Errorf("upgrade is only supported for managed nodegroups; could not find one with name %q", options.NodegroupName)
		}
		return err
	}

	if hasStacks {
		return m.upgradeUsingStack(options, nodegroupOutput.Nodegroup)
	}

	return m.upgradeUsingAPI(options, nodegroupOutput.Nodegroup)
}

func (m *Manager) upgradeUsingAPI(options UpgradeOptions, nodegroup *eks.Nodegroup) error {
	input := &eks.UpdateNodegroupVersionInput{
		ClusterName:   &m.cfg.Metadata.Name,
		Force:         &options.ForceUpgrade,
		NodegroupName: &options.NodegroupName,
		Version:       &options.KubernetesVersion,
	}

	if options.LaunchTemplateVersion != "" {
		lt := nodegroup.LaunchTemplate
		if lt == nil || (lt.Id == nil && lt.Name == nil) {
			return errors.New("cannot update launch template version because the nodegroup is not configured to use one")
		}

		input.LaunchTemplate = &eks.LaunchTemplateSpecification{
			Version: &options.LaunchTemplateVersion,
		}

		if lt.Id != nil {
			input.LaunchTemplate.Id = nodegroup.LaunchTemplate.Id
		} else {
			input.LaunchTemplate.Name = nodegroup.LaunchTemplate.Name

		}
	}

	if options.KubernetesVersion == "" {
		// Use the current Kubernetes version
		version, err := semver.ParseTolerant(*nodegroup.Version)
		if err != nil {
			return errors.Wrapf(err, "unexpected error parsing Kubernetes version %q", *nodegroup.Version)
		}
		input.Version = aws.String(fmt.Sprintf("%v.%v", version.Major, version.Minor))
	}

	upgradeResponse, err := m.ctl.Provider.EKS().UpdateNodegroupVersion(input)

	if err != nil {
		return err
	}

	if upgradeResponse != nil {
		logger.Debug("upgrade response for %q: %s", options.NodegroupName, upgradeResponse.String())
	}

	logger.Info("upgrade of nodegroup %q in progress", options.NodegroupName)

	if options.Wait {
		return m.waitForUpgrade(options)
	}

	return nil
}

func (m *Manager) waitForUpgrade(options UpgradeOptions) error {

	newRequest := func() *request.Request {
		input := &eks.DescribeNodegroupInput{
			ClusterName:   &m.cfg.Metadata.Name,
			NodegroupName: &options.NodegroupName,
		}
		req, _ := m.ctl.Provider.EKS().DescribeNodegroupRequest(input)
		return req
	}

	msg := fmt.Sprintf("waiting for upgrade of nodegroup %q to complete", options.NodegroupName)

	acceptors := waiters.MakeAcceptors(
		"Nodegroup.Status",
		eks.NodegroupStatusActive,
		[]string{
			eks.NodegroupStatusDegraded,
		},
	)

	err := m.wait(options.NodegroupName, msg, acceptors, newRequest, m.ctl.Provider.WaitTimeout(), nil)
	if err != nil {
		return err
	}
	logger.Info("nodegroup successfully upgraded")
	return nil
}

// upgradeUsingStack upgrades nodegroup to the latest AMI release for the specified Kubernetes version, or
// the current Kubernetes version if the version isn't specified
// If options.LaunchTemplateVersion is set, it also upgrades the nodegroup to the specified launch template version
func (m *Manager) upgradeUsingStack(options UpgradeOptions, nodegroup *eks.Nodegroup) error {
	if options.KubernetesVersion != "" && options.ReleaseVersion != "" {
		return errors.New("only one of kubernetes-version or release-version can be specified")
	}

	template, err := m.stackManager.GetManagedNodeGroupTemplate(options.NodegroupName)
	if err != nil {
		return errors.Wrap(err, "error fetching nodegroup template")
	}

	stack, err := goformation.ParseJSON([]byte(template))
	if err != nil {
		return errors.Wrap(err, "unexpected error parsing nodegroup template")
	}

	ngResources := stack.GetAllEKSNodegroupResources()
	ngResource, ok := ngResources[builder.ManagedNodeGroupResourceName]
	if !ok {
		return errors.New("unexpected error: failed to find nodegroup resource in nodegroup stack")
	}

	updateStack := func(stack *cloudformation.Template, wait bool) error {
		bytes, err := stack.JSON()
		if err != nil {
			return err
		}

		if err := m.stackManager.UpdateNodeGroupStack(options.NodegroupName, string(bytes), true); err != nil {
			return errors.Wrap(err, "error updating nodegroup stack")
		}
		return nil
	}

	requiresUpdate, err := m.requiresStackUpdate(options.NodegroupName)
	if err != nil {
		return err
	}
	if requiresUpdate {
		logger.Info("updating nodegroup stack to a newer format before upgrading nodegroup version")
		// always wait for the main stack update
		if err := updateStack(stack, true); err != nil {
			return err
		}
	}

	if ngResource.ForceUpdateEnabled == nil || strings.ToLower(ngResource.ForceUpdateEnabled.String()) != strconv.FormatBool(options.ForceUpgrade) {
		ngResource.ForceUpdateEnabled = gfnt.NewBoolean(options.ForceUpgrade)
		logger.Info("setting ForceUpdateEnabled value to %t", options.ForceUpgrade)
		if err := updateStack(stack, true); err != nil {
			return err
		}
	}

	ltResources := stack.GetAllEC2LaunchTemplateResources()

	if options.LaunchTemplateVersion != "" {
		// TODO validate launch template version
		if len(ltResources) == 1 {
			return errors.New("launch-template-version is only valid if a nodegroup is using an explicit launch template")
		}
		if ngResource.LaunchTemplate == nil || ngResource.LaunchTemplate.Id == nil {
			return errors.New("nodegroup does not use a launch template")
		}
	}

	usesCustomAMI, err := m.usesCustomAMI(ltResources, ngResource)
	if err != nil {
		return err
	}

	if usesCustomAMI && (options.KubernetesVersion != "" || options.ReleaseVersion != "") {
		return errors.New("cannot specify kubernetes-version or release-version when using a custom AMI")
	}

	if options.ReleaseVersion != "" {
		ngResource.ReleaseVersion = gfnt.NewString(options.ReleaseVersion)
	} else if !usesCustomAMI {
		kubernetesVersion := options.KubernetesVersion
		if kubernetesVersion == "" {
			// Use the current Kubernetes version
			version, err := semver.ParseTolerant(*nodegroup.Version)
			if err != nil {
				return errors.Wrapf(err, "unexpected error parsing Kubernetes version %q", *nodegroup.Version)
			}
			kubernetesVersion = fmt.Sprintf("%v.%v", version.Major, version.Minor)
		}

		latestReleaseVersion, err := m.getLatestReleaseVersion(kubernetesVersion, nodegroup)
		if err != nil {
			return err
		}

		if latestReleaseVersion != "" {
			if err := m.updateReleaseVersion(latestReleaseVersion, options.LaunchTemplateVersion, nodegroup, ngResource); err != nil {
				return err
			}
		} else {
			ngResource.Version = gfnt.NewString(kubernetesVersion)
		}
	}
	if options.LaunchTemplateVersion != "" {
		ngResource.LaunchTemplate.Version = gfnt.NewString(options.LaunchTemplateVersion)
	}

	ngResource.ForceUpdateEnabled = gfnt.NewBoolean(options.ForceUpgrade)

	logger.Info("upgrading nodegroup version")
	if err := updateStack(stack, options.Wait); err != nil {
		return err
	}
	logger.Info("nodegroup successfully upgraded")
	return nil
}

func (m *Manager) updateReleaseVersion(latestReleaseVersion, launchTemplateVersion string, nodegroup *eks.Nodegroup, ngResource *gfneks.Nodegroup) error {
	latest, err := ParseReleaseVersion(latestReleaseVersion)
	if err != nil {
		return err
	}
	current, err := ParseReleaseVersion(*nodegroup.ReleaseVersion)
	if err != nil {
		return err
	}

	if latest.LTE(current) && launchTemplateVersion == "" {
		logger.Info("nodegroup %q is already up-to-date", *nodegroup.NodegroupName)
		return nil
	}
	if latest.GTE(current) {
		ngResource.ReleaseVersion = gfnt.NewString(latestReleaseVersion)
	}
	return nil
}

func (m *Manager) requiresStackUpdate(nodeGroupName string) (bool, error) {
	ngStack, err := m.stackManager.DescribeNodeGroupStack(nodeGroupName)
	if err != nil {
		return false, err
	}

	ver, found, err := manager.GetEksctlVersionFromTags(ngStack.Tags)
	if err != nil {
		return false, err
	}
	if !found {
		return true, nil
	}

	curVer, err := version.ParseEksctlVersion(version.GetVersion())
	if err != nil {
		return false, errors.Wrap(err, "unexpected error parsing current eksctl version")
	}
	return !ver.EQ(curVer), nil
}

func (m *Manager) getLatestReleaseVersion(kubernetesVersion string, nodeGroup *eks.Nodegroup) (string, error) {
	ssmParameterName, err := ami.MakeManagedSSMParameterName(kubernetesVersion, *nodeGroup.AmiType)
	if err != nil {
		return "", err
	}

	if ssmParameterName == "" {
		return "", nil
	}

	ssmOutput, err := m.ctl.Provider.SSM().GetParameter(&ssm.GetParameterInput{
		Name: &ssmParameterName,
	})
	if err != nil {
		return "", err
	}
	return *ssmOutput.Parameter.Value, nil
}

func (m *Manager) usesCustomAMI(ltResources map[string]*gfnec2.LaunchTemplate, ng *gfneks.Nodegroup) (bool, error) {
	if lt, ok := ltResources["LaunchTemplate"]; ok {
		return lt.LaunchTemplateData.ImageId != nil, nil
	}

	if ng.LaunchTemplate == nil || ng.LaunchTemplate.Id == nil {
		return false, nil
	}

	lt := &api.LaunchTemplate{
		ID: ng.LaunchTemplate.Id.String(),
	}
	if version := ng.LaunchTemplate.Version; version != nil {
		lt.Version = aws.String(version.String())
	}

	customLaunchTemplate, err := m.launchTemplateFetcher.Fetch(lt)
	if err != nil {
		return false, errors.Wrap(err, "error fetching launch template data")
	}
	return customLaunchTemplate.ImageId != nil, nil
}
