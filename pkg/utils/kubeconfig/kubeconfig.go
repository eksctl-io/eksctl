package kubeconfig

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/weaveworks/eksctl/pkg/utils"

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/tools/clientcmd/api/latest"

	"github.com/kubicorn/kubicorn/pkg/logger"
)

var DefaultPath = clientcmd.RecommendedHomeFile

// WriteToFile will write Kubernetes client configuration to a file.
// If file doesn't exist it will be created. If the file exists then
// the configuration will be merged with the existing file.
func WriteToFile(path string, config *api.Config, setContext bool) error {
	exists, err := utils.FileExists(path)
	if err != nil {
		return errors.Wrapf(err, "error trying to read config file %q", path)
	}

	if !exists {
		logger.Debug("kubeconfig file doesn't exist: %s", path)
		return write(path, config)
	}

	existing, err := read(path)
	if err != nil {
		return errors.Wrapf(err, "unable to read existing kubeconfig file %q", path)
	}

	logger.Debug("merging kubeconfig files")
	merged, err := merge(existing, config)
	if err != nil {
		return errors.Wrapf(err, "unable to merge configuration with existing kubeconfig file %q", path)
	}

	if setContext && len(config.CurrentContext) > 0 {
		logger.Debug("setting current-context to %s", config.CurrentContext)
		merged.CurrentContext = config.CurrentContext
	}

	return write(path, merged)
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

func write(path string, config *api.Config) error {
	if err := clientcmd.WriteToFile(*config, path); err != nil {
		return errors.Wrapf(err, "couldn't write client config file %q", path)
	}
	return nil
}

func read(path string) (*api.Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "error trying to read config file %q", path)
	}

	if len(data) == 0 {
		return api.NewConfig(), nil
	}

	config, _, err := latest.Codec.Decode(data, nil, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't decode kubeconfig file %q", path)
	}

	return config.(*api.Config), nil
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
