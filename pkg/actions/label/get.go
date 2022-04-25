package label

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"

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
			labels, err = m.getLabelsFromUnownedNodeGroup(nodeGroupName)
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

func (m *Manager) getLabelsFromUnownedNodeGroup(nodeGroupName string) (map[string]string, error) {
	out, err := m.eksAPI.DescribeNodegroup(&eks.DescribeNodegroupInput{
		ClusterName:   aws.String(m.clusterName),
		NodegroupName: aws.String(nodeGroupName),
	})
	if err != nil {
		return nil, err
	}

	labels := make(map[string]string, len(out.Nodegroup.Labels))
	for k, v := range out.Nodegroup.Labels {
		labels[k] = *v
	}

	return labels, nil
}
