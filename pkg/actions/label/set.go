package label

import (
	"context"

	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

func (m *Manager) Set(ctx context.Context, nodeGroupName string, labels map[string]string) error {
	err := m.service.UpdateLabels(ctx, nodeGroupName, labels, nil)
	if manager.IsStackDoesNotExistError(err) {
		return m.setLabelsOnUnownedNodeGroup(ctx, nodeGroupName, labels)
	}
	return err
}

func (m *Manager) setLabelsOnUnownedNodeGroup(ctx context.Context, nodeGroupName string, labels map[string]string) error {
	_, err := m.eksAPI.UpdateNodegroupConfig(ctx, &eks.UpdateNodegroupConfigInput{
		ClusterName:   aws.String(m.clusterName),
		NodegroupName: aws.String(nodeGroupName),
		Labels:        &ekstypes.UpdateLabelsPayload{AddOrUpdateLabels: labels},
	})
	return err
}
