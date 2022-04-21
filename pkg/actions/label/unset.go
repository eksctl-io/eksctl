package label

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

func (m *Manager) Unset(ctx context.Context, nodeGroupName string, labels []string) error {
	err := m.service.UpdateLabels(ctx, nodeGroupName, nil, labels)
	if err != nil {
		switch {
		case manager.IsStackDoesNotExistError(err):
			return m.unsetLabelsOnUnownedNodeGroup(nodeGroupName, labels)
		default:
			return err
		}
	}
	return nil
}

func (m *Manager) unsetLabelsOnUnownedNodeGroup(nodeGroupName string, labels []string) error {
	var pointyLabels []*string
	for _, v := range labels {
		pointyLabels = append(pointyLabels, &v)
	}
	_, err := m.eksAPI.UpdateNodegroupConfig(&eks.UpdateNodegroupConfigInput{
		ClusterName:   aws.String(m.clusterName),
		NodegroupName: aws.String(nodeGroupName),
		Labels:        &eks.UpdateLabelsPayload{RemoveLabels: pointyLabels},
	})
	if err != nil {
		return err
	}

	return nil
}
