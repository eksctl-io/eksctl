package authconfigmap

import (
	"fmt"

	// go go:embed to work
	_ "embed"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/assetutil"
	"github.com/weaveworks/eksctl/pkg/iam"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
)

//go:embed assets/emr-containers-rbac.yaml
var emrContainersRbacYamlBytes []byte

type ServiceName string

const (
	emrContainers ServiceName = "emr-containers"
)

type serviceDetails struct {
	User        ServiceName
	IAMRoleName string
	Namespaced  bool
}

var (
	emrContainersService = serviceDetails{
		User:        emrContainers,
		IAMRoleName: "AWSServiceRoleForAmazonEMRContainers",
		Namespaced:  true,
	}
)

// Grants an AWS service access to an EKS cluster
type ServiceAccess struct {
	rawClient *kubernetes.RawClient
	acm       *AuthConfigMap
	accountID string
}

// NewServiceAccess creates a new ServiceAccess
func NewServiceAccess(rawClient *kubernetes.RawClient, acm *AuthConfigMap, accountID string) *ServiceAccess {
	return &ServiceAccess{
		rawClient: rawClient,
		acm:       acm,
		accountID: accountID,
	}
}

// Grant grants access to the specified service
func (s *ServiceAccess) Grant(serviceName, namespace string, partition string) error {
	resources, serviceDetails, err := lookupService(serviceName)
	if err != nil {
		return err
	}
	if serviceDetails.Namespaced && namespace == "" {
		return errors.Errorf("namespace is required for %s", serviceName)
	}
	if !serviceDetails.Namespaced && namespace != "" {
		return errors.Errorf("namespace is not valid for %s", serviceName)
	}

	list, err := kubernetes.NewList(resources)
	if err != nil {
		return err
	}

	for _, item := range list.Items {
		if err := s.applyResource(item.Object, namespace); err != nil {
			return err
		}
	}

	role := &iam.RoleIdentity{
		RoleARN: fmt.Sprintf("arn:%s:iam::%s:role/%s", partition, s.accountID, serviceDetails.IAMRoleName),
		KubernetesIdentity: iam.KubernetesIdentity{
			KubernetesUsername: string(serviceDetails.User),
		},
	}
	err = s.acm.AddIdentityIfNotPresent(role, func(identity iam.Identity) bool {
		return identity.ARN() == role.ARN() && identity.Username() == role.Username()
	})
	if err != nil {
		return err
	}

	if err := s.acm.Save(); err != nil {
		return errors.Wrap(err, "error applying service role")
	}
	return nil
}

func (s *ServiceAccess) applyResource(o runtime.Object, namespace string) error {
	if namespace != "" {
		metadataAccessor := meta.NewAccessor()
		if err := metadataAccessor.SetNamespace(o, namespace); err != nil {
			return errors.Wrap(err, "unexpected error setting namespace")
		}
	}
	r, err := s.rawClient.NewRawResource(o)
	if err != nil {
		return err
	}

	msg, err := r.CreateOrReplace(false)
	if err != nil {
		return errors.Wrap(err, "error applying resource")
	}
	logger.Info(msg)
	return nil
}

func lookupService(serviceName string) (resources []byte, sd serviceDetails, err error) {
	defer func() {
		if r := recover(); r != nil {
			if ae, ok := r.(*assetutil.Error); ok {
				err = ae
			} else {
				panic(r)
			}
		}
	}()

	switch ServiceName(serviceName) {
	case emrContainers:
		return emrContainersRbacYamlBytes, emrContainersService, nil
	default:
		return nil, sd, errors.Errorf("invalid service name %q", serviceName)
	}
}
