package label

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

func (m *Manager) Unset(ctx context.Context, nodeGroupName string, labels []string) error {
	err := m.service.UpdateLabels(ctx, nodeGroupName, nil, labels)
	if err != nil {
		switch {
		case manager.IsStackDoesNotExistError(err):
			return m.unsetLabelsOnUnownedNodeGroup(ctx, nodeGroupName, labels)
		default:
			return err
		}
	}
	return nil
}

func (m *Manager) unsetLabelsOnUnownedNodeGroup(ctx context.Context, nodeGroupName string, labels []string) error {
	_, err := m.eksAPI.UpdateNodegroupConfig(ctx, &eks.UpdateNodegroupConfigInput{
		ClusterName:   aws.String(m.clusterName),
		NodegroupName: aws.String(nodeGroupName),
		Labels:        &ekstypes.UpdateLabelsPayload{RemoveLabels: labels},
	})
	if err != nil {
		return err
	}

	return nil
}
