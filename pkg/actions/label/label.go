package label

import (
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/pkg/errors"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_managed_service.go . Service
type Service interface {
	GetLabels(nodeGroupName string) (map[string]string, error)
	UpdateLabels(nodeGroupName string, labelsToAdd map[string]string, labelsToRemove []string) error
}

type Manager struct {
	service     Service
	eksAPI      eksiface.EKSAPI
	clusterName string
}

func New(clusterName string, service Service, eksAPI eksiface.EKSAPI) *Manager {
	return &Manager{
		service:     service,
		eksAPI:      eksAPI,
		clusterName: clusterName,
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
