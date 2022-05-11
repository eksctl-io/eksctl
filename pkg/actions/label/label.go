package label

import (
	"context"

	"github.com/weaveworks/eksctl/pkg/awsapi"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_managed_service.go . Service
type Service interface {
	GetLabels(ctx context.Context, nodeGroupName string) (map[string]string, error)
	UpdateLabels(ctx context.Context, nodeGroupName string, labelsToAdd map[string]string, labelsToRemove []string) error
}

type Manager struct {
	service     Service
	eksAPI      awsapi.EKS
	clusterName string
}

func New(clusterName string, service Service, eksAPI awsapi.EKS) *Manager {
	return &Manager{
		service:     service,
		eksAPI:      eksAPI,
		clusterName: clusterName,
	}
}
