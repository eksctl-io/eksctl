package addon

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/hashicorp/go-version"
	"github.com/kris-nova/logger"
	kubeclient "k8s.io/client-go/kubernetes"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
)

type Manager struct {
	clusterConfig *api.ClusterConfig
	eksAPI        awsapi.EKS
	withOIDC      bool
	oidcManager   *iamoidc.OpenIDConnectManager
	stackManager  manager.StackManager
	clientSet     kubeclient.Interface
}

func New(clusterConfig *api.ClusterConfig, eksAPI awsapi.EKS, stackManager manager.StackManager, withOIDC bool, oidcManager *iamoidc.OpenIDConnectManager, clientSet kubeclient.Interface) (*Manager, error) {
	return &Manager{
		clusterConfig: clusterConfig,
		eksAPI:        eksAPI,
		withOIDC:      withOIDC,
		oidcManager:   oidcManager,
		stackManager:  stackManager,
		clientSet:     clientSet,
	}, nil
}

func (a *Manager) waitForAddonToBeActive(ctx context.Context, addon *api.Addon, waitTimeout time.Duration) error {
	// Don't wait for coredns or aws-ebs-csi-driver if there are no nodegroups.
	// They will be in degraded state until nodegroups are added.
	if (addon.Name == api.CoreDNSAddon || addon.Name == api.AWSEBSCSIDriverAddon) && !a.clusterConfig.HasNodes() {
		return nil
	}
	activeWaiter := eks.NewAddonActiveWaiter(a.eksAPI)
	input := &eks.DescribeAddonInput{
		ClusterName: &a.clusterConfig.Metadata.Name,
		AddonName:   &addon.Name,
	}
	if err := activeWaiter.Wait(ctx, input, waitTimeout); err != nil {
		getAddonStatus := func() string {
			output, describeErr := a.eksAPI.DescribeAddon(ctx, input)
			if describeErr != nil {
				return describeErr.Error()
			}
			return string(output.Addon.Status)
		}

		switch {
		case strings.Contains(err.Error(), "exceeded max wait time"):
			return fmt.Errorf("timed out waiting for addon %q to become active, status: %q", addon.Name, getAddonStatus())
		case strings.Contains(err.Error(), "waiter state transitioned to Failure"):
			return fmt.Errorf("addon status transitioned to %q", getAddonStatus())
		default:
			return err
		}
	}
	logger.Info("addon %q active", addon.Name)
	return nil
}

func (a *Manager) getLatestMatchingVersion(ctx context.Context, addon *api.Addon) (string, error) {
	addonInfos, err := a.describeVersions(ctx, addon)
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
