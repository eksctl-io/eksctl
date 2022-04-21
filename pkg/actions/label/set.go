package label

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/smithy-go"
	"github.com/pkg/errors"
)

func (m *Manager) Set(ctx context.Context, nodeGroupName string, labels map[string]string) error {
	err := m.service.UpdateLabels(ctx, nodeGroupName, labels, nil)
	if err != nil {
		fmt.Println("THIS IS THE ERROR: ", err)
		if awsErr, ok := errors.Cause(err).(smithy.APIError); ok {
			if awsErr.ErrorCode() == "ValidationError" {
				return m.setLabelsOnUnownedNodeGroup(nodeGroupName, labels)
			}
		}
	}
	return err
}

func (m *Manager) setLabelsOnUnownedNodeGroup(nodeGroupName string, labels map[string]string) error {
	pointyLabels := aws.StringMap(labels)
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
