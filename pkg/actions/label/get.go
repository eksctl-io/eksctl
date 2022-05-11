package label

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

type Summary struct {
	Cluster   string
	NodeGroup string
	Labels    map[string]string
}

func (m *Manager) Get(ctx context.Context, nodeGroupName string) ([]Summary, error) {
	var (
		labels map[string]string
		err    error
	)

	labels, err = m.service.GetLabels(ctx, nodeGroupName)
	if err != nil {
		switch {
		case manager.IsStackDoesNotExistError(err):
			labels, err = m.getLabelsFromUnownedNodeGroup(ctx, nodeGroupName)
			if err != nil {
				return nil, err
			}
		default:
			return nil, err
		}
	}

	return []Summary{
		{
			Cluster:   m.clusterName,
			NodeGroup: nodeGroupName,
			Labels:    labels,
		},
	}, nil
}

func (m *Manager) getLabelsFromUnownedNodeGroup(ctx context.Context, nodeGroupName string) (map[string]string, error) {
	out, err := m.eksAPI.DescribeNodegroup(ctx, &eks.DescribeNodegroupInput{
		ClusterName:   aws.String(m.clusterName),
		NodegroupName: aws.String(nodeGroupName),
	})
	if err != nil {
		return nil, err
	}

	return out.Nodegroup.Labels, nil
}
