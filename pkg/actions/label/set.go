package label

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
)

func (m *Manager) Set(nodeGroupName string, labels map[string]string) error {
	err := m.service.UpdateLabels(nodeGroupName, labels, nil)
	if err != nil {
		switch {
		case isValidationError(err):
			return m.setLabelsOnUnownedNodeGroup(nodeGroupName, labels)
		default:
			return err
		}
	}
	return nil
}

func (m *Manager) setLabelsOnUnownedNodeGroup(nodeGroupName string, labels map[string]string) error {
	pointyLabels := make(map[string]*string, len(labels))
	for k, v := range labels {
		pointyLabels[k] = &v
	}
	_, err := m.eksAPI.UpdateNodegroupConfig(&eks.UpdateNodegroupConfigInput{
		ClusterName:   aws.String(m.clusterName),
		NodegroupName: aws.String(nodeGroupName),
		Labels:        &eks.UpdateLabelsPayload{AddOrUpdateLabels: pointyLabels},
	})
	if err != nil {
		return err
	}

	return nil
}
