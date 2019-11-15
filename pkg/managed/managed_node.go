package managed

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

type Service struct {
	provider        v1alpha5.ClusterProvider
	clusterName     string
	stackCollection *manager.StackCollection
}

type HealthIssue struct {
	Message string
	Code    string
}

const (
	labelsPath = "Resources.ManagedNodeGroup.Properties.Labels"
)

func NewService(provider v1alpha5.ClusterProvider, stackCollection *manager.StackCollection, clusterName string) *Service {
	return &Service{provider: provider, stackCollection: stackCollection, clusterName: clusterName}
}

func (m *Service) GetHealth(nodeGroupName string) ([]HealthIssue, error) {
	input := &eks.DescribeNodegroupInput{
		ClusterName:   &m.clusterName,
		NodegroupName: &nodeGroupName,
	}

	output, err := m.provider.EKS().DescribeNodegroup(input)
	if err != nil {
		awsError, ok := err.(awserr.Error)
		if ok && awsError.Code() == eks.ErrCodeResourceNotFoundException {
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

func (m *Service) GetLabels(nodeGroupName string) (map[string]string, error) {
	template, err := m.stackCollection.GetManagedNodeGroupTemplate(nodeGroupName)
	if err != nil {
		return nil, err
	}
	return extractLabels(template)
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
