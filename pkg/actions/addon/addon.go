package addon

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/hashicorp/go-version"
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

func New(clusterConfig *api.ClusterConfig, eksAPI eksiface.EKSAPI, stackManager manager.StackManager, withOIDC bool, oidcManager *iamoidc.OpenIDConnectManager, clientSet kubeclient.Interface, timeout time.Duration) (*Manager, error) {
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
		timeout:       timeout,
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

func (a *Manager) getLatestMatchingVersion(addon *api.Addon) (string, error) {
	addonInfos, err := a.describeVersions(addon)
	if err != nil {
		return "", err
	}
	if len(addonInfos.Addons) == 0 || len(addonInfos.Addons[0].AddonVersions) == 0 {
		return "", fmt.Errorf("no versions available for %q", addon.Name)
	}

	addonVersion := addon.Version
	var versions []*version.Version
	for _, addonVersionInfo := range addonInfos.Addons[0].AddonVersions {
		v, err := a.parseVersion(*addonVersionInfo.AddonVersion)
		if err != nil {
			return "", err
		}

		if addonVersion == "latest" || strings.Contains(*addonVersionInfo.AddonVersion, addonVersion) {
			versions = append(versions, v)
		}
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("no version(s) found matching %q for %q", addonVersion, addon.Name)
	}

	sort.SliceStable(versions, func(i, j int) bool {
		return versions[j].LessThan(versions[i])
	})
	return versions[0].Original(), nil
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

func (a *Manager) parseVersion(v string) (*version.Version, error) {
	version, err := version.NewVersion(v)
	if err != nil {
		return nil, fmt.Errorf("failed to parse version %q: %w", v, err)
	}
	return version, nil
}
