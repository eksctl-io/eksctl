package kubeconf

import (
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/tools/clientcmd/api/latest"

	"github.com/kubicorn/kubicorn/pkg/logger"
)

// GetRecommendedPath returns the recommended kubeconf path based on the
// Kubernetes Go Client. If the KUBECONFIG environment variable is set it will
// use this path otherwise it will use the default path to the config file.
func GetRecommendedPath() string {
	kubeVar := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)

	if len(kubeVar) > 0 {
		return kubeVar
	}
	return clientcmd.RecommendedHomeFile
}

// WriteToFile will write Kubernetes client configuration to a file.
// If file doesn't exist it will be created. If the file exists then
// the configuration will be merged with the existing file.
func WriteToFile(path string, config *api.Config, setContext bool) error {
	exists, err := fileExists(path)
	if err != nil {
		return errors.Wrapf(err, "error trying to read config file %q", path)
	}

	if !exists {
		logger.Debug("Kube configuration file doesn't exist: %s", path)
		return writeConfToFile(path, config)
	}

	existing, err := readConfigurationFile(path)
	if err != nil {
		return errors.Wrapf(err, "unable to read existing kube configuration file %q", path)
	}

	logger.Debug("Merging kubeconfigurations files")
	merged, err := mergeConfigurations(existing, config)
	if err != nil {
		return errors.Wrapf(err, "unable to merge configuration with existing kube configuration file %q", path)
	}

	if setContext && len(config.CurrentContext) > 0 {
		logger.Debug("Setting current-context to %s", config.CurrentContext)
		merged.CurrentContext = config.CurrentContext
	}

	return writeConfToFile(path, merged)
}

func mergeConfigurations(existing *api.Config, tomerge *api.Config) (*api.Config, error) {
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

func writeConfToFile(path string, config *api.Config) error {
	if err := clientcmd.WriteToFile(*config, path); err != nil {
		return errors.Wrapf(err, "couldn't write client config file %q", path)
	}
	return nil
}

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func readConfigurationFile(path string) (*api.Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "error trying to read config file %q", path)
	}

	if len(data) == 0 {
		return api.NewConfig(), nil
	}

	config, _, err := latest.Codec.Decode(data, nil, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't decode kubeconf file %q", path)
	}

	return config.(*api.Config), nil
}
