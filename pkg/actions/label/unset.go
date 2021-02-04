package label

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
)

func (m *Manager) Unset(nodeGroupName string, labels []string) error {
	err := m.service.UpdateLabels(nodeGroupName, nil, labels)
	if err != nil {
		switch {
		case isValidationError(err):
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
