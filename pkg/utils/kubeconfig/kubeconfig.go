package kubeconfig

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/weaveworks/eksctl/pkg/utils"

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/kubicorn/kubicorn/pkg/logger"
)

// WriteKubeCfg will write Kubernetes client configuration to a file.
// If path isn't specified then the path will be determined by client-go.
// If file pointed to by the doesn't exist it will be created.
// If the file exists then the configuration will be merged with the existing file.
func WriteKubeCfg(path string, newConfig *api.Config, setContext bool) error {
	configAccess := getConfigAccess(path)

	config, err := configAccess.GetStartingConfig()

	logger.Debug("merging kubeconfig files")
	merged, err := merge(config, newConfig)
	if err != nil {
		return errors.Wrapf(err, "unable to merge configuration with existing kubeconfig file %q", path)
	}

	if setContext && len(newConfig.CurrentContext) > 0 {
		logger.Debug("setting current-context to %s", config.CurrentContext)
		merged.CurrentContext = newConfig.CurrentContext
	}

	if err := clientcmd.ModifyConfig(configAccess, *config, true); err != nil {
		return nil
	}

	return nil
}

func getConfigAccess(explicitPath string) clientcmd.ConfigAccess {
	pathOptions := clientcmd.NewDefaultPathOptions()
	if explicitPath != "" {
		pathOptions.LoadingRules.ExplicitPath = explicitPath
	}

	return interface{}(pathOptions).(clientcmd.ConfigAccess)
}

func merge(existing *api.Config, tomerge *api.Config) (*api.Config, error) {
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

func MaybeDeleteConfig(name string) {
	p := AutoPath(name)

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

	// Print message to manually remove from config file
	logger.Warning("as you are not using the auto-generated kubeconfig file you will need to remove the details of cluster %s manually", name)
}
