package managed

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/blang/semver"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/weaveworks/eksctl/pkg/ami"
	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

// A Service provides methods for managing managed nodegroups
type Service struct {
	provider        v1alpha5.ClusterProvider
	clusterName     string
	stackCollection *manager.StackCollection
}

// HealthIssue represents a health issue with a managed nodegroup
type HealthIssue struct {
	Message string
	Code    string
}

// TODO use goformation types
const (
	labelsPath         = "Resources.ManagedNodeGroup.Properties.Labels"
	releaseVersionPath = "Resources.ManagedNodeGroup.Properties.ReleaseVersion"
	versionPath        = "Resources.ManagedNodeGroup.Properties.Version"
)

// NewService creates a new Service
func NewService(provider v1alpha5.ClusterProvider, stackCollection *manager.StackCollection, clusterName string) *Service {
	return &Service{provider: provider, stackCollection: stackCollection, clusterName: clusterName}
}

// GetHealth fetches the health status for a nodegroup
func (m *Service) GetHealth(nodeGroupName string) ([]HealthIssue, error) {
	input := &eks.DescribeNodegroupInput{
		ClusterName:   &m.clusterName,
		NodegroupName: &nodeGroupName,
	}

	output, err := m.provider.EKS().DescribeNodegroup(input)
	if err != nil {
		if isNotFound(err) {
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

	return m.stackCollection.UpdateNodeGroupStack(nodeGroupName, template)
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
func (m *Service) UpgradeNodeGroup(nodeGroupName, kubernetesVersion string) error {
	// Use the latest AMI release version
	output, err := m.provider.EKS().DescribeNodegroup(&eks.DescribeNodegroupInput{
		ClusterName:   &m.clusterName,
		NodegroupName: &nodeGroupName,
	})

	if err != nil {
		if isNotFound(err) {
			return fmt.Errorf("upgrade is only supported for managed nodegroups; could not find one with name %q",
				nodeGroupName)
		}
		return err
	}

	nodeGroup := output.Nodegroup

	if kubernetesVersion == "" {
		// Use the current Kubernetes version
		kubernetesVersion = *nodeGroup.Version
	} else if _, err := semver.ParseTolerant(kubernetesVersion); err != nil {
		return errors.Wrap(err, "invalid Kubernetes version")
	}

	// Upgrade only Version for kubernetes, by default CF will use the latest working AMI for ReleaseVersion
	// Docs: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-eks-nodegroup.html#cfn-eks-nodegroup-releaseversion
	if kubernetesVersion != *nodeGroup.Version {
		return m.updateNodeGroupVersion(nodeGroupName, kubernetesVersion)
	}

	instanceType := nodeGroup.InstanceTypes[0]
	ssmParameterName, err := ami.MakeSSMParameterName(kubernetesVersion, *instanceType, v1alpha5.NodeImageFamilyAmazonLinux2)
	if err != nil {
		return err
	}

	ssmOutput, err := m.provider.SSM().GetParameter(&ssm.GetParameterInput{
		Name: &ssmParameterName,
	})
	if err != nil {
		return err
	}

	imageID := *ssmOutput.Parameter.Value

	// To get the Kubernetes patch version, as it's not reported in the SSM GetParameter call
	imagesOutput, err := m.provider.EC2().DescribeImages(&ec2.DescribeImagesInput{
		ImageIds: aws.StringSlice([]string{imageID}),
	})

	if err != nil {
		return err
	}

	if len(imagesOutput.Images) != 1 {
		return fmt.Errorf("expected to find exactly 1 image; got %d", len(imagesOutput.Images))
	}

	image := *imagesOutput.Images[0]
	amiReleaseVersion, err := extractAMIReleaseVersion(*image.Name)
	if err != nil {
		return errors.Wrap(err, "error extracting AMI release version")
	}

	kubernetesVersion, err = extractKubeVersion(*image.Description)
	if err != nil {
		return errors.Wrap(err, "error extracting Kubernetes version")
	}
	releaseVersion := makeReleaseVersion(kubernetesVersion, amiReleaseVersion)
	if releaseVersion == *nodeGroup.ReleaseVersion {
		logger.Info("nodegroup %q is already up-to-date", nodeGroupName)
		return nil
	}
	return m.updateNodeGroupReleaseVersion(nodeGroupName, releaseVersion)
}

func (m *Service) updateNodeGroupVersion(nodeGroupName, kubernetesVersion string) error {
	template, err := m.stackCollection.GetManagedNodeGroupTemplate(nodeGroupName)
	if err != nil {
		return err
	}

	template, err = sjson.Set(template, versionPath, kubernetesVersion)
	if err != nil {
		return err
	}

	return m.stackCollection.UpdateNodeGroupStack(nodeGroupName, template)
}

func (m *Service) updateNodeGroupReleaseVersion(nodeGroupName, releaseVersion string) error {
	template, err := m.stackCollection.GetManagedNodeGroupTemplate(nodeGroupName)
	if err != nil {
		return err
	}

	template, err = sjson.Set(template, releaseVersionPath, releaseVersion)
	if err != nil {
		return err
	}

	return m.stackCollection.UpdateNodeGroupStack(nodeGroupName, template)
}

func isNotFound(err error) bool {
	awsError, ok := err.(awserr.Error)
	return ok && awsError.Code() == eks.ErrCodeResourceNotFoundException
}

var (
	kubeVersionRegex = regexp.MustCompile(`\(k8s:\s([\d.]+),`)
	amiVersionRegex  = regexp.MustCompile(`v(\d+)$`)
)

// extractKubeVersion extracts the full Kubernetes version (including patch number) from the image description
// format: "EKS Kubernetes Worker AMI with AmazonLinux2 image, (k8s: 1.13.11, docker:18.06)"
func extractKubeVersion(description string) (string, error) {
	match := kubeVersionRegex.FindStringSubmatch(description)
	if len(match) != 2 {
		return "", fmt.Errorf("expected 2 matching items; got %d", len(match))
	}
	return match[1], nil
}

// extractAMIReleaseVersion extracts the AMI release version from the image name
// the format of the image name is amazon-eks-node-1.14-v20190927
func extractAMIReleaseVersion(imageName string) (string, error) {
	match := amiVersionRegex.FindStringSubmatch(imageName)
	if len(match) != 2 {
		return "", fmt.Errorf("expected 2 matching items; got %d", len(match))
	}
	return match[1], nil
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

func makeReleaseVersion(kubernetesVersion, amiVersion string) string {
	return fmt.Sprintf("%s-%s", kubernetesVersion, amiVersion)
}
