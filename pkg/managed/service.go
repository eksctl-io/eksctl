package managed

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

// A Service provides methods for managing managed nodegroups
type Service struct {
	eksAPI                awsapi.EKS
	launchTemplateFetcher *builder.LaunchTemplateFetcher
	clusterName           string
	stackCollection       manager.StackManager
}

// HealthIssue represents a health issue with a managed nodegroup
type HealthIssue struct {
	Message string
	Code    string
}

// TODO use goformation types
const (
	labelsPath = "Resources.ManagedNodeGroup.Properties.Labels"
)

func NewService(eksAPI awsapi.EKS, ec2API awsapi.EC2,
	stackCollection manager.StackManager, clusterName string) *Service {
	return &Service{
		eksAPI:                eksAPI,
		stackCollection:       stackCollection,
		launchTemplateFetcher: builder.NewLaunchTemplateFetcher(ec2API),
		clusterName:           clusterName,
	}
}

// GetHealth fetches the health status for a nodegroup
func (m *Service) GetHealth(ctx context.Context, nodeGroupName string) ([]HealthIssue, error) {
	input := &eks.DescribeNodegroupInput{
		ClusterName:   &m.clusterName,
		NodegroupName: &nodeGroupName,
	}

	output, err := m.eksAPI.DescribeNodegroup(ctx, input)
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
			Code:    string(issue.Code),
		})
	}

	return healthIssues, nil
}

// UpdateLabels adds or removes labels for a nodegroup
func (m *Service) UpdateLabels(ctx context.Context, nodeGroupName string, labelsToAdd map[string]string, labelsToRemove []string) error {
	template, err := m.stackCollection.GetManagedNodeGroupTemplate(ctx, manager.GetNodegroupOption{
		NodeGroupName: nodeGroupName,
	})
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

	return m.stackCollection.UpdateNodeGroupStack(ctx, nodeGroupName, template, true)
}

// GetLabels fetches the labels for a nodegroup
func (m *Service) GetLabels(ctx context.Context, nodeGroupName string) (map[string]string, error) {
	template, err := m.stackCollection.GetManagedNodeGroupTemplate(ctx, manager.GetNodegroupOption{
		NodeGroupName: nodeGroupName,
	})
	if err != nil {
		return nil, err
	}
	return extractLabels(template)
}

func IsNotFound(err error) bool {
	var notFoundErr *ekstypes.ResourceNotFoundException
	return errors.As(err, &notFoundErr)
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
