package kubeconfig

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/weaveworks/eksctl/pkg/eks/api"
	"github.com/weaveworks/eksctl/pkg/utils"

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/kubicorn/kubicorn/pkg/logger"
)

var DefaultPath = clientcmd.RecommendedHomeFile

const (
	HeptioAuthenticatorAWS = "heptio-authenticator-aws"
	AWSIAMAuthenticator    = "aws-iam-authenticator"
)

// New creates Kubernetes client configuration for a given username
// if certificateAuthorityPath is no empty, it is used instead of
// embedded certificate-authority-data
func New(spec *api.ClusterConfig, username, certificateAuthorityPath string) (*clientcmdapi.Config, string, string) {
	clusterName := getCompleteClusterName(spec)
	contextName := fmt.Sprintf("%s@%s", username, clusterName)

	c := &clientcmdapi.Config{
		Clusters: map[string]*clientcmdapi.Cluster{
			clusterName: {
				Server: spec.Endpoint,
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			contextName: {
				Cluster:  clusterName,
				AuthInfo: contextName,
			},
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			contextName: &clientcmdapi.AuthInfo{},
		},
		CurrentContext: contextName,
	}

	if certificateAuthorityPath == "" {
		c.Clusters[clusterName].CertificateAuthorityData = spec.CertificateAuthorityData
	} else {
		c.Clusters[clusterName].CertificateAuthority = certificateAuthorityPath
	}

	return c, clusterName, contextName
}

func AppendAuthenticator(c *clientcmdapi.Config, spec *api.ClusterConfig, command string) {
	c.AuthInfos[c.CurrentContext].Exec = &clientcmdapi.ExecConfig{
		APIVersion: "client.authentication.k8s.io/v1alpha1",
		Command:    command,
		Args:       []string{"token", "-i", spec.ClusterName},
		/*
			Args:       []string{"token", "-i", c.Cluster.ClusterName, "-r", c.roleARN},
		*/
	}
}

// Write will write Kubernetes client configuration to a file.
// If path isn't specified then the path will be determined by client-go.
// If file pointed to by path doesn't exist it will be created.
// If the file already exists then the configuration will be merged with the existing file.
func Write(path string, newConfig *clientcmdapi.Config, setContext bool) (string, error) {
	configAccess := getConfigAccess(path)

	config, err := configAccess.GetStartingConfig()

	logger.Debug("merging kubeconfig files")
	merged, err := merge(config, newConfig)
	if err != nil {
		return "", errors.Wrapf(err, "unable to merge configuration with existing kubeconfig file %q", path)
	}

	if setContext && newConfig.CurrentContext != "" {
		logger.Debug("setting current-context to %s", newConfig.CurrentContext)
		merged.CurrentContext = newConfig.CurrentContext
	}

	if err := clientcmd.ModifyConfig(configAccess, *merged, true); err != nil {
		return "", nil
	}

	return configAccess.GetDefaultFilename(), nil
}

func getCompleteClusterName(spec *api.ClusterConfig) string {
	return fmt.Sprintf("%s.%s.eksctl.io", spec.ClusterName, spec.Region)
}

func getConfigAccess(explicitPath string) clientcmd.ConfigAccess {
	pathOptions := clientcmd.NewDefaultPathOptions()
	if explicitPath != "" && explicitPath != DefaultPath {
		pathOptions.LoadingRules.ExplicitPath = explicitPath
	}

	return interface{}(pathOptions).(clientcmd.ConfigAccess)
}
func merge(existing *clientcmdapi.Config, tomerge *clientcmdapi.Config) (*clientcmdapi.Config, error) {
	for k, v := range tomerge.Clusters {
		existing.Clusters[k] = v
	}
	for k, v := range tomerge.AuthInfos {
		existing.AuthInfos[k] = v
	}
	for k, v := range tomerge.Contexts {
		existing.Contexts[k] = v
	}

	return existing, nil
}

func AutoPath(name string) string {
	return path.Join(clientcmd.RecommendedConfigDir, "eksctl", "clusters", name)
}

func isValidConfig(p, name string) error {
	clientConfig, err := clientcmd.LoadFromFile(p)
	if err != nil {
		return errors.Wrapf(err, "unable to load config %q", p)
	}

	if err := clientcmd.ConfirmUsable(*clientConfig, ""); err != nil {
		return errors.Wrapf(err, "unable to parse config %q", p)
	}

	// we want to make sure we only delete config files that haven't be modified by the user
	// checking context name is a good start, we might want ot do deeper checks later, e.g. checksum,
	// as we don't want to delete any files by accident that didn't belong to us
	ctxFmtErr := fmt.Errorf("unable to verify ownership of config %q, unexpected contex name %q", p, clientConfig.CurrentContext)

	ctx := strings.Split(clientConfig.CurrentContext, "@")
	if len(ctx) != 2 {
		return ctxFmtErr
	}
	if strings.HasPrefix(ctx[1], name+".") && strings.HasSuffix(ctx[1], ".eksctl.io") {
		return nil
	}
	return ctxFmtErr
}

func tryDeleteConfig(p, name string) {
	if err := isValidConfig(p, name); err != nil {
		logger.Debug("ignoring error while checking config file %q: %s", p, err.Error())
		return
	}
	if err := os.Remove(p); err != nil {
		logger.Debug("ignoring error while removing config file %q: %s", p, err.Error())
	}
}

func MaybeDeleteConfig(ctl *api.ClusterConfig) {
	p := AutoPath(ctl.ClusterName)

	autoConfExists, err := utils.FileExists(p)
	if err != nil {
		logger.Debug("error checking if auto-generated kubeconfig file exists %q: %s", p, err.Error())
		return
	}
	if autoConfExists {
		if err := os.Remove(p); err != nil {
			logger.Debug("ignoring error while removing auto-generated config file %q: %s", p, err.Error())
		}
		return
	}

	configAccess := getConfigAccess(DefaultPath)
	config, _ := configAccess.GetStartingConfig()
	if !deleteClusterInfo(config, ctl) {
		return
	}

	if err := clientcmd.ModifyConfig(configAccess, *config, true); err != nil {
		logger.Debug("ignoring error while failing to update config file %q: %s", DefaultPath, err)
	} else {
		logger.Success("kubeconfig has been updated")
	}
}

// deleteClusterInfo removes a cluster's information from the kubeconfig if the cluster name
// provided by ctl matches a eksctl-created cluster in the kubeconfig
// returns 'true' if the existing config has changes and 'false' otherwise
func deleteClusterInfo(existing *clientcmdapi.Config, ctl *api.ClusterConfig) bool {
	isChanged := false
	clusterName := getCompleteClusterName(ctl)

	if existing.Clusters[clusterName] != nil {
		delete(existing.Clusters, clusterName)
		logger.Debug("removed cluster %q from kubeconfig", clusterName)
		isChanged = true
	}

	for username, context := range existing.Contexts {
		if context.Cluster == clusterName {
			delete(existing.Contexts, username)
			logger.Debug("removed context for %q from kubeconfig", username)
			isChanged = true
			if existing.AuthInfos[username] != nil {
				delete(existing.AuthInfos, username)
				logger.Debug("removed auth info for %q from kubeconfig", username)
			}
			break
		}
	}

	return isChanged
}
