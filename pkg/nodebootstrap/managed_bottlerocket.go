package nodebootstrap

import (
	"encoding/base64"

	toml "github.com/pelletier/go-toml"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type ManagedBottlerocket struct {
	clusterConfig *api.ClusterConfig
	ng            *api.ManagedNodeGroup
}

// NewManagedBottlerocketBootstrapper returns a new bootstrapper for managed Bottlerocket.
func NewManagedBottlerocketBootstrapper(clusterConfig *api.ClusterConfig, ng *api.ManagedNodeGroup) *ManagedBottlerocket {
	return &ManagedBottlerocket{
		clusterConfig: clusterConfig,
		ng:            ng,
	}
}

// UserData generates TOML userdata for bootstrapping a Bottlerocket node.
func (b *ManagedBottlerocket) UserData() (string, error) {
	if err := b.setDerivedSettings(); err != nil {
		return "", err
	}

	settings, err := toml.TreeFromMap(map[string]interface{}{
		"settings": *b.ng.Bottlerocket.Settings,
	})
	if err != nil {
		return "", errors.Wrap(err, "error loading user-provided Bottlerocket settings")
	}

	// Check and protect all input key names against TOML's dotted key semantics.
	ProtectTOMLKeys([]string{"settings"}, settings)

	if enableAdminContainer := b.ng.Bottlerocket.EnableAdminContainer; enableAdminContainer != nil {
		const adminContainerEnabledKey = "settings.host-containers.admin.enabled"
		if settings.Has(adminContainerEnabledKey) {
			return "", errors.Errorf("cannot set both bottlerocket.enableAdminContainer and %s", adminContainerEnabledKey)
		}
		settings.Set(adminContainerEnabledKey, *enableAdminContainer)
	}

	userData := settings.String()
	if userData == "" {
		return "", errors.New("generated unexpected empty TOML user-data from input")
	}

	return base64.StdEncoding.EncodeToString([]byte(userData)), nil
}

// setDerivedSettings configures settings that are derived from top-level nodegroup config
// as opposed to settings configured in `bottlerocket.settings`.
func (b *ManagedBottlerocket) setDerivedSettings() error {
	kubernetesSettings, err := extractKubernetesSettings(b.ng)
	if err != nil {
		return err
	}
	if err := validateBottlerocketSettings(kubernetesSettings); err != nil {
		return err
	}

	if b.ng.MaxPodsPerNode != 0 {
		kubernetesSettings["max-pods"] = b.ng.MaxPodsPerNode
	}

	return nil
}

// validateBottlerocketSettings validates the supplied Kubernetes settings to ensure fields related to bootstrapping
// and fields available on the ManagedNodeGroup object are not set by the user.
func validateBottlerocketSettings(kubernetesSettings map[string]interface{}) error {
	clusterBootstrapKeys := []string{"cluster-certificate", "api-server", "cluster-name", "cluster-dns-ip"}
	for _, k := range clusterBootstrapKeys {
		if _, ok := kubernetesSettings[k]; ok {
			return errors.Errorf("cannot set settings.kubernetes.%s; EKS automatically injects cluster bootstrapping fields into user-data", k)
		}
	}

	apiFields := []string{"node-labels", "node-taints"}
	for _, k := range apiFields {
		if _, ok := kubernetesSettings[k]; ok {
			return errors.Errorf("cannot set settings.kubernetes.%s; labels and taints should be set on the managedNodeGroup object", k)
		}
	}

	return nil
}
