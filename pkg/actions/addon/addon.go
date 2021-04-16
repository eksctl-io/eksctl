package addon

import (
	"context"
	"fmt"
	"time"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/cfn/waiter"
	kubeclient "k8s.io/client-go/kubernetes"

	"github.com/weaveworks/eksctl/pkg/utils"

	"github.com/weaveworks/eksctl/pkg/cfn/manager"

	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type Manager struct {
	clusterConfig *api.ClusterConfig
	eksAPI        eksiface.EKSAPI
	withOIDC      bool
	oidcManager   *iamoidc.OpenIDConnectManager
	stackManager  manager.StackManager
	clientSet     kubeclient.Interface
	timeout       time.Duration
}

func New(clusterConfig *api.ClusterConfig, eksAPI eksiface.EKSAPI, stackManager manager.StackManager, withOIDC bool, oidcManager *iamoidc.OpenIDConnectManager, clientSet kubeclient.Interface) (*Manager, error) {
	if err := supportedVersion(clusterConfig.Metadata.Version); err != nil {
		return nil, err
	}

	return &Manager{
		clusterConfig: clusterConfig,
		eksAPI:        eksAPI,
		withOIDC:      withOIDC,
		oidcManager:   oidcManager,
		stackManager:  stackManager,
		clientSet:     clientSet,
		timeout:       5 * time.Minute,
	}, nil
}

func (a *Manager) waitForAddonToBeActive(addon *api.Addon) error {
	var out *awseks.DescribeAddonOutput
	operation := func() (bool, error) {
		var err error
		out, err = a.eksAPI.DescribeAddon(&awseks.DescribeAddonInput{
			ClusterName: &a.clusterConfig.Metadata.Name,
			AddonName:   &addon.Name,
		})
		if err != nil {
			return false, err
		}
		if *out.Addon.Status == awseks.AddonStatusActive {
			return true, nil
		}
		return false, nil
	}

	w := waiter.Waiter{
		Operation: operation,
		NextDelay: func(_ int) time.Duration {
			return a.timeout / 10
		},
	}

	err := w.WaitWithTimeout(a.timeout)
	if err != nil {
		if err == context.DeadlineExceeded {
			return errors.Errorf("timed out waiting for addon %q to become active, status: %q", addon.Name, *out.Addon.Status)
		}
		return err
	}
	logger.Info("addon %q active", addon.Name)
	return nil
}

func supportedVersion(version string) error {
	supported, err := utils.IsMinVersion(api.Version1_18, version)
	if err != nil {
		return err
	}
	switch supported {
	case true:
		return nil
	default:
		return fmt.Errorf("addons not supported on %s. Must be using %s or newer", version, api.Version1_18)
	}
}

func (a *Manager) makeAddonName(name string) string {
	return fmt.Sprintf("eksctl-%s-addon-%s", a.clusterConfig.Metadata.Name, name)
}
