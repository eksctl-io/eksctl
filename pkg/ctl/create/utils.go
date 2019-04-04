package create

import (
	"fmt"
	"github.com/weaveworks/eksctl/pkg/ssh"
	"strings"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func checkSubnetsGivenAsFlags() bool {
	return len(*subnets[api.SubnetTopologyPrivate])+len(*subnets[api.SubnetTopologyPublic]) != 0
}

func checkVersion(ctl *eks.ClusterProvider, meta *api.ClusterMeta) error {
	switch meta.Version {
	case "auto":
		break
	case "":
		meta.Version = "auto"
	case "latest":
		meta.Version = api.LatestVersion
		logger.Info("will use version latest version (%s) for new nodegroup(s)", meta.Version)
	default:
		validVersion := false
		for _, v := range api.SupportedVersions() {
			if meta.Version == v {
				validVersion = true
			}
		}
		if !validVersion {
			return fmt.Errorf("invalid version %s, supported values: auto, latest, %s", meta.Version, strings.Join(api.SupportedVersions(), ", "))
		}
	}

	if v := ctl.ControlPlaneVersion(); v == "" {
		return fmt.Errorf("unable to get control plane version")
	} else if meta.Version == "auto" {
		meta.Version = v
		logger.Info("will use version %s for new nodegroup(s) based on control plane version", meta.Version)
	} else if meta.Version != v {
		hint := "--version=auto"
		if clusterConfigFile != "" {
			hint = "metadata.version: auto"
		}
		logger.Warning("will use version %s for new nodegroup(s), while control plane version is %s; to automatically inherit the version use %q", meta.Version, v, hint)
	}

	return nil
}

// loadSSHKey loads the ssh public key specified in the NodeGroup. The key should be specified
// in only one way: by name (for a key existing in EC2), by path (for a key in a local file)
// or by its contents (in the config-file).
func loadSSHKey(ng *api.NodeGroup, clusterName string, provider api.ClusterProvider) error {
	sshConfig := ng.SSH
	if sshConfig.Allow == nil || *sshConfig.Allow == false {
		return nil
	}

	switch {

	// Load Key by content
	case sshConfig.PublicKey != nil:
		ssh.LoadSSHKeyByContent(sshConfig.PublicKey, clusterName, provider, ng)

	// Use key by name in EC2
	case sshConfig.PublicKeyName != nil && *sshConfig.PublicKeyName != "":
		if err := ssh.CheckKeyExistsInEc2(*sshConfig.PublicKeyName, provider); err != nil {
			return err
		}
		logger.Info("using EC2 key pair %q", *sshConfig.PublicKeyName)

	// Local ssh key file
	case ssh.FileExists(*sshConfig.PublicKeyPath):
		ssh.LoadSSHKeyFromFile(*sshConfig.PublicKeyPath, clusterName, provider, ng)

	// A keyPath, when specified as a flag, can mean a local key or a key name in EC2
	default:
		err := ssh.CheckKeyExistsInEc2(*sshConfig.PublicKeyPath, provider)
		if err != nil {
			ng.SSH.PublicKeyName = sshConfig.PublicKeyPath
			ng.SSH.PublicKeyPath = nil
			return err
		}
		logger.Info("using EC2 key pair %q", *ng.SSH.PublicKeyName)
	}

	return nil
}
