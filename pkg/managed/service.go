package managed

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"

	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"

	"github.com/aws/aws-sdk-go/service/eks/eksiface"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/blang/semver"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/weaveworks/goformation/v4/cloudformation"

	"github.com/weaveworks/eksctl/pkg/version"

	"github.com/weaveworks/goformation/v4"
	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfneks "github.com/weaveworks/goformation/v4/cloudformation/eks"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

// A Service provides methods for managing managed nodegroups
type Service struct {
	eksAPI                eksiface.EKSAPI
	ssmAPI                ssmiface.SSMAPI
	launchTemplateFetcher *builder.LaunchTemplateFetcher
	clusterName           string
	stackCollection       manager.StackManager
}

// HealthIssue represents a health issue with a managed nodegroup
type HealthIssue struct {
	Message string
	Code    string
}

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

// TODO use goformation types
const (
	labelsPath = "Resources.ManagedNodeGroup.Properties.Labels"
)

func NewService(eksAPI eksiface.EKSAPI, ssmAPI ssmiface.SSMAPI, ec2API ec2iface.EC2API,
	stackCollection manager.StackManager, clusterName string) *Service {
	return &Service{
		eksAPI:                eksAPI,
		ssmAPI:                ssmAPI,
		stackCollection:       stackCollection,
		launchTemplateFetcher: builder.NewLaunchTemplateFetcher(ec2API),
		clusterName:           clusterName,
	}
}

// GetHealth fetches the health status for a nodegroup
func (m *Service) GetHealth(nodeGroupName string) ([]HealthIssue, error) {
	input := &eks.DescribeNodegroupInput{
		ClusterName:   &m.clusterName,
		NodegroupName: &nodeGroupName,
	}

	output, err := m.eksAPI.DescribeNodegroup(input)
	if err != nil {
		if IsNotFound(err) {
			return nil, errors.Wrapf(err, "could not find a managed nodegroup with name %q", nodeGroupName)
		}
		return nil, err
	}

	health := output.Nodegroup.Health
	if health == nil || len(health.Issues) == 0 {
		// No health issues
		return nil, nil
	}

	var healthIssues []HealthIssue
	for _, issue := range health.Issues {
		healthIssues = append(healthIssues, HealthIssue{
			Message: *issue.Message,
			Code:    *issue.Code,
		})
	}

	return healthIssues, nil
}

// UpdateLabels adds or removes labels for a nodegroup
func (m *Service) UpdateLabels(nodeGroupName string, labelsToAdd map[string]string, labelsToRemove []string) error {
	template, err := m.stackCollection.GetManagedNodeGroupTemplate(nodeGroupName)
	if err != nil {
		return err
	}

	newLabels, err := extractLabels(template)
	if err != nil {
		return err
	}

	for k, v := range labelsToAdd {
		newLabels[k] = v
	}

	for _, k := range labelsToRemove {
		delete(newLabels, k)
	}

	template, err = sjson.Set(template, labelsPath, newLabels)
	if err != nil {
		return err
	}

	return m.stackCollection.UpdateNodeGroupStack(nodeGroupName, template, true)
}

// GetLabels fetches the labels for a nodegroup
func (m *Service) GetLabels(nodeGroupName string) (map[string]string, error) {
	template, err := m.stackCollection.GetManagedNodeGroupTemplate(nodeGroupName)
	if err != nil {
		return nil, err
	}
	return extractLabels(template)
}

// UpgradeNodeGroup upgrades nodegroup to the latest AMI release for the specified Kubernetes version, or
// the current Kubernetes version if the version isn't specified
// If options.LaunchTemplateVersion is set, it also upgrades the nodegroup to the specified launch template version
func (m *Service) UpgradeNodeGroup(options UpgradeOptions) error {
	output, err := m.eksAPI.DescribeNodegroup(&eks.DescribeNodegroupInput{
		ClusterName:   &m.clusterName,
		NodegroupName: &options.NodegroupName,
	})

	if err != nil {
		if IsNotFound(err) {
			return fmt.Errorf("upgrade is only supported for managed nodegroups; could not find one with name %q", options.NodegroupName)
		}
		return err
	}

	if options.KubernetesVersion != "" && options.ReleaseVersion != "" {
		return errors.New("only one of kubernetes-version or release-version can be specified")
	}

	nodeGroup := output.Nodegroup

	template, err := m.stackCollection.GetManagedNodeGroupTemplate(options.NodegroupName)
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
		if err := m.stackCollection.UpdateNodeGroupStack(options.NodegroupName, string(bytes), true); err != nil {
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
			version, err := semver.ParseTolerant(*nodeGroup.Version)
			if err != nil {
				return errors.Wrapf(err, "unexpected error parsing Kubernetes version %q", *nodeGroup.Version)
			}
			kubernetesVersion = fmt.Sprintf("%v.%v", version.Major, version.Minor)
		}
		latestReleaseVersion, err := m.getLatestReleaseVersion(kubernetesVersion, nodeGroup)
		if err != nil {
			return err
		}
		latest, err := parseReleaseVersion(latestReleaseVersion)
		if err != nil {
			return err
		}
		current, err := parseReleaseVersion(*nodeGroup.ReleaseVersion)
		if err != nil {
			return err
		}

		if latest.LTE(current) && options.LaunchTemplateVersion == "" {
			logger.Info("nodegroup %q is already up-to-date", *nodeGroup.NodegroupName)
			return nil
		}
		if latest.GTE(current) {
			ngResource.ReleaseVersion = gfnt.NewString(latestReleaseVersion)
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

// parseReleaseVersion parses an AMI release version string that's in the format `1.18.8-20201007`
func parseReleaseVersion(releaseVersion string) (amiReleaseVersion, error) {
	parts := strings.Split(releaseVersion, "-")
	if len(parts) != 2 {
		return amiReleaseVersion{}, errors.Errorf("unexpected format for release version: %q", releaseVersion)
	}
	v, err := semver.ParseTolerant(parts[0])
	if err != nil {
		return amiReleaseVersion{}, errors.Wrap(err, "invalid SemVer version")
	}
	return amiReleaseVersion{
		Version: v,
		Date:    parts[1],
	}, nil
}

type amiReleaseVersion struct {
	Version semver.Version
	Date    string
}

// LTE checks if a is less than or equal to b.
func (a amiReleaseVersion) LTE(b amiReleaseVersion) bool {
	return a.Compare(b) <= 0
}

// GTE checks if a is greater than or equal to b.
func (a amiReleaseVersion) GTE(b amiReleaseVersion) bool {
	return a.Compare(b) >= 0
}

// Compare returns 0 if a==b, -1 if a < b, and +1 if a > b.
func (a amiReleaseVersion) Compare(b amiReleaseVersion) int {
	cmp := a.Version.Compare(b.Version)
	if cmp == 0 {
		return strings.Compare(a.Date, b.Date)
	}
	return cmp
}

func (m *Service) requiresStackUpdate(nodeGroupName string) (bool, error) {
	ngStack, err := m.stackCollection.DescribeNodeGroupStack(nodeGroupName)
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

func (m *Service) getLatestReleaseVersion(kubernetesVersion string, nodeGroup *eks.Nodegroup) (string, error) {
	ssmParameterName, err := ami.MakeManagedSSMParameterName(kubernetesVersion, api.NodeImageFamilyAmazonLinux2, *nodeGroup.AmiType)
	if err != nil {
		return "", err
	}

	ssmOutput, err := m.ssmAPI.GetParameter(&ssm.GetParameterInput{
		Name: &ssmParameterName,
	})
	if err != nil {
		return "", err
	}
	return *ssmOutput.Parameter.Value, nil
}

func (m *Service) usesCustomAMI(ltResources map[string]*gfnec2.LaunchTemplate, ng *gfneks.Nodegroup) (bool, error) {
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

func IsNotFound(err error) bool {
	awsError, ok := err.(awserr.Error)
	return ok && awsError.Code() == eks.ErrCodeResourceNotFoundException
}

// TODO switch to using goformation types
func extractLabels(template string) (map[string]string, error) {
	labelsValue := gjson.Get(template, labelsPath)
	if !labelsValue.Exists() {
		return nil, errors.New("failed to find labels")
	}
	values, ok := labelsValue.Value().(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected type for labels: %T", labelsValue.Value())
	}

	labels := make(map[string]string)
	for k, v := range values {
		value, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("unexpected type for label value: %T", value)
		}
		labels[k] = value
	}

	return labels, nil
}
