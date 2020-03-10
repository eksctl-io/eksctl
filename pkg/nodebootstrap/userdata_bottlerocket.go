package nodebootstrap

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/pelletier/go-toml"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// NewUserDataForBottlerocket generates TOML userdata for bootstrapping a Bottlerocket node.
func NewUserDataForBottlerocket(spec *api.ClusterConfig, ng *api.NodeGroup) (string, error) {
	if ng.Bottlerocket.Settings == nil {
		ng.Bottlerocket.Settings = &api.InlineDocument{}
	}

	// Update settings based on NodeGroup configuration. Values set here are not
	// allowed to be set by the user - the values are owned by the NodeGroup and
	// expressly written into settings.
	if err := setDerivedBottlerocketSettings(ng); err != nil {
		return "", err
	}

	settings, err := toml.TreeFromMap(map[string]interface{}{
		"settings": *ng.Bottlerocket.Settings,
	})
	if err != nil {
		return "", errors.Wrap(err, "error loading user provided settings")
	}

	// All input settings key names need to be checked & protected against
	// TOML's dotted key semantics.
	protectTOMLKeys([]string{"settings"}, settings)

	// Generate TOML for launch in this NodeGroup.
	data, err := bottlerocketSettingsTOML(spec, ng, settings)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString([]byte(data)), nil
}

func setDerivedBottlerocketSettings(ng *api.NodeGroup) error {
	settings := *ng.Bottlerocket.Settings

	var kubernetesSettings map[string]interface{}

	if val, ok := settings["kubernetes"]; ok {
		kubernetesSettings, ok = val.(map[string]interface{})
		if !ok {
			return errors.Errorf("expected settings.kubernetes to be of type %T; got %T", kubernetesSettings, val)
		}
	} else {
		kubernetesSettings = make(map[string]interface{})
		settings["kubernetes"] = kubernetesSettings
	}

	if len(ng.Labels) != 0 {
		kubernetesSettings["node-labels"] = ng.Labels
	}
	if len(ng.Taints) != 0 {
		kubernetesSettings["node-taints"] = ng.Taints
	}
	if ng.MaxPodsPerNode != 0 {
		kubernetesSettings["max-pods"] = ng.MaxPodsPerNode
	}
	if ng.ClusterDNS != "" {
		kubernetesSettings["cluster-dns-ip"] = ng.ClusterDNS
	}
	return nil
}

// protectTOMLKeys processes a tree finding and replacing dotted keys
// with quoted keys to retain the configured settings. This prevents
// TOML parsers from deserializing keys into nested key-value pairs at
// each dot encountered - which is not uncommon in the context of
// Kubernetes' labels, annotations, and taints.
func protectTOMLKeys(path []string, tree *toml.Tree) {
	pathTree, ok := tree.GetPath(path).(*toml.Tree)
	if !ok {
		return
	}

	keys := pathTree.Keys()
	for _, key := range keys {
		val := pathTree.GetPath([]string{key})
		keyPath := append(path, key)

		if strings.Contains(key, ".") {
			quotedKey := fmt.Sprintf("%q", key)
			quotedPath := append(path, quotedKey)
			err := pathTree.DeletePath([]string{key})
			if err == nil {
				pathTree.SetPath([]string{quotedKey}, val)
				keyPath = quotedPath
			}
		}
		if _, ok := val.(*toml.Tree); ok {
			protectTOMLKeys(keyPath, tree)
		}
	}
}

// bottlerocketSettingsTOML generates TOML userdata for configuring
// settings on Bottlerocket nodes.
func bottlerocketSettingsTOML(spec *api.ClusterConfig, ng *api.NodeGroup, tree *toml.Tree) (string, error) {
	const insertWithoutComment = false // `false` indicates that the item should be inserted without commenting it out
	// Set, or override, cluster settings' keys to provide latest EKS cluster
	// data.
	tree.SetWithComment("settings.kubernetes.cluster-certificate", "Kubernetes Cluster CA Certificate",
		insertWithoutComment,
		base64.StdEncoding.EncodeToString([]byte(spec.Status.CertificateAuthorityData)))
	tree.SetWithComment("settings.kubernetes.api-server", "Kubernetes Control Plane API Endpoint",
		insertWithoutComment,
		spec.Status.Endpoint)
	tree.SetWithComment("settings.kubernetes.cluster-name", "Kubernetes Cluster Name",
		insertWithoutComment,
		spec.Metadata.Name)

	// Don't override user's explicit setting if they provided it in config.
	if !tree.Has("settings.host-containers.admin.enabled") {
		// Provide value only if given, with `enabled`
		// commented out otherwise.
		enabled := false
		isUnset := ng.Bottlerocket == nil || ng.Bottlerocket.EnableAdminContainer == nil
		if !isUnset {
			enabled = *ng.Bottlerocket.EnableAdminContainer
		}
		tree.SetWithComment("settings.host-containers.admin.enabled", "Bottlerocket Admin Container",
			isUnset, // comment out if not specified by config
			enabled)
	}

	userdata := tree.String()
	if userdata == "" {
		return "", errors.New("generated unexpected empty TOML user-data from input")
	}
	return userdata, nil
}
