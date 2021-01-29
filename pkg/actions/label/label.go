package label

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/managed"
)

//go:generate counterfeiter -o fakes/fake_managed_service.go . Service
type Service interface {
	GetLabels(nodeGroupName string) (map[string]string, error)
	UpdateLabels(nodeGroupName string, labelsToAdd map[string]string, labelsToRemove []string) error
}

type Manager struct {
	service     Service
	eksAPI      eksiface.EKSAPI
	clusterName string
}

func New(cfg *api.ClusterConfig, ctl *eks.ClusterProvider) *Manager {
	sc := manager.NewStackCollection(ctl.Provider, cfg)
	mc := managed.NewService(ctl.Provider, sc, cfg.Metadata.Name)
	return &Manager{
		service:     mc,
		eksAPI:      ctl.Provider.EKS(),
		clusterName: cfg.Metadata.Name,
	}
}

// If a ValidationError code is returned then an eksctl-marked stack was not
// found for that nodegroup so we can then try to call the EKS api directly.
func isValidationError(err error) bool {
	awsErr, ok := errors.Cause(err).(awserr.Error)
	if !ok {
		return false
	}
	return awsErr.Code() == "ValidationError"
}
