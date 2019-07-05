package kubeconfig

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/weaveworks/eksctl/pkg/utils/file"

	"os/exec"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// DefaultPath defines the default path
var DefaultPath = clientcmd.RecommendedHomeFile

const (
	// AWSIAMAuthenticator defines the name of the AWS IAM authenticator
	AWSIAMAuthenticator = "aws-iam-authenticator"
	// HeptioAuthenticatorAWS defines the old name of AWS IAM authenticator
	HeptioAuthenticatorAWS = "heptio-authenticator-aws"
	// AWSEKSAuthenticator defines the recently added `aws eks get-token` command
	AWSEKSAuthenticator = "aws"
)

// AuthenticatorCommands returns all of authenticator commands
func AuthenticatorCommands() []string {
	return []string{
		AWSIAMAuthenticator,
		HeptioAuthenticatorAWS,
		AWSEKSAuthenticator,
	}
}

// New creates Kubernetes client configuration for a given username
// if certificateAuthorityPath is not empty, it is used instead of
// embedded certificate-authority-data
func New(spec *api.ClusterConfig, username, certificateAuthorityPath string) (*clientcmdapi.Config, string, string) {
	clusterName := spec.Metadata.String()
	contextName := fmt.Sprintf("%s@%s", username, clusterName)

	c := &clientcmdapi.Config{
		Clusters: map[string]*clientcmdapi.Cluster{
			clusterName: {
				Server: spec.Status.Endpoint,
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
		c.Clusters[clusterName].CertificateAuthorityData = spec.Status.CertificateAuthorityData
	} else {
		c.Clusters[clusterName].CertificateAuthority = certificateAuthorityPath
	}

	return c, clusterName, contextName
}

// NewForKubectl creates configuration for kubectl using a suitable authenticator
func NewForKubectl(spec *api.ClusterConfig, username, roleARN, profile string) *clientcmdapi.Config {
	config, _, _ := New(spec, username, "")
	authenticator, found := LookupAuthenticator()
	if !found {
		// fall back to aws-iam-authenticator
		authenticator = AWSIAMAuthenticator
	}
	AppendAuthenticator(config, spec, authenticator, roleARN, profile)
	return config
}

// AppendAuthenticator appends the AWS IAM  authenticator, and
// if profile is non-empty string it sets AWS_PROFILE environment
// variable also
func AppendAuthenticator(config *clientcmdapi.Config, spec *api.ClusterConfig, authenticatorCMD, roleARN, profile string) {
	var (
		args        []string
		roleARNFlag string
	)

	switch authenticatorCMD {
	case AWSIAMAuthenticator, HeptioAuthenticatorAWS:
		args = []string{"token", "-i", spec.Metadata.Name}
		roleARNFlag = "-r"
	case AWSEKSAuthenticator:
		args = []string{"eks", "get-token", "--cluster-name", spec.Metadata.Name}
		roleARNFlag = "--role-arn"
		if spec.Metadata.Region != "" {
			args = append(args, "--region", spec.Metadata.Region)
		}
	}
	if roleARN != "" {
		args = append(args, roleARNFlag, roleARN)
	}

	execConfig := &clientcmdapi.ExecConfig{
		APIVersion: "client.authentication.k8s.io/v1alpha1",
		Command:    authenticatorCMD,
		Args:       args,
	}

	if profile != "" {
		execConfig.Env = []clientcmdapi.ExecEnvVar{
			{
				Name:  "AWS_PROFILE",
				Value: profile,
			},
		}
	}

	config.AuthInfos[config.CurrentContext] = &clientcmdapi.AuthInfo{
		Exec: execConfig,
	}
}

// Write will write Kubernetes client configuration to a file.
// If path isn't specified then the path will be determined by client-go.
// If file pointed to by path doesn't exist it will be created.
// If the file already exists then the configuration will be merged with the existing file.
func Write(path string, newConfig clientcmdapi.Config, setContext bool) (string, error) {
	configAccess := getConfigAccess(path)

	config, err := configAccess.GetStartingConfig()
	if err != nil {
		return "", errors.Wrapf(err, "enable to read existing kubeconfig file %q", path)
	}

	logger.Debug("merging kubeconfig files")
	merged := merge(config, &newConfig)

	if setContext && newConfig.CurrentContext != "" {
		logger.Debug("setting current-context to %s", newConfig.CurrentContext)
		merged.CurrentContext = newConfig.CurrentContext
	}

	if err := clientcmd.ModifyConfig(configAccess, *merged, true); err != nil {
		return "", nil
	}

	return configAccess.GetDefaultFilename(), nil
}

func getConfigAccess(explicitPath string) clientcmd.ConfigAccess {
	pathOptions := clientcmd.NewDefaultPathOptions()
	if explicitPath != "" && explicitPath != DefaultPath {
		pathOptions.LoadingRules.ExplicitPath = explicitPath
	}

	return interface{}(pathOptions).(clientcmd.ConfigAccess)
}
func merge(existing *clientcmdapi.Config, tomerge *clientcmdapi.Config) *clientcmdapi.Config {
	for k, v := range tomerge.Clusters {
		existing.Clusters[k] = v
	}
	for k, v := range tomerge.AuthInfos {
		existing.AuthInfos[k] = v
	}
	for k, v := range tomerge.Contexts {
		existing.Contexts[k] = v
	}

	return existing
}

// AutoPath returns the path for the auto-generated kubeconfig
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

// MaybeDeleteConfig will delete the auto-generated kubeconfig, if it exists
func MaybeDeleteConfig(cl *api.ClusterMeta) {
	p := AutoPath(cl.Name)

	if file.Exists(p) {
		if err := isValidConfig(p, cl.Name); err != nil {
			logger.Debug(err.Error())
			return
		}
		if err := os.Remove(p); err != nil {
			logger.Debug("ignoring error while removing auto-generated config file %q: %s", p, err.Error())
		}
		return
	}

	configAccess := getConfigAccess(DefaultPath)
	config, err := configAccess.GetStartingConfig()
	if err != nil {
		logger.Debug("error reading kubeconfig file %q: %s", DefaultPath, err.Error())
		return
	}

	if !deleteClusterInfo(config, cl) {
		return
	}

	if err := clientcmd.ModifyConfig(configAccess, *config, true); err != nil {
		logger.Debug("ignoring error while failing to update config file %q: %s", DefaultPath, err.Error())
	} else {
		logger.Success("kubeconfig has been updated")
	}
}

// deleteClusterInfo removes a cluster's information from the kubeconfig if the cluster name
// provided by ctl matches a eksctl-created cluster in the kubeconfig
// returns 'true' if the existing config has changes and 'false' otherwise
func deleteClusterInfo(existing *clientcmdapi.Config, cl *api.ClusterMeta) bool {
	isChanged := false
	clusterName := cl.String()

	if _, ok := existing.Clusters[clusterName]; ok {
		delete(existing.Clusters, clusterName)
		logger.Debug("removed cluster %q from kubeconfig", clusterName)
		isChanged = true
	}

	var currentContextName string
	for name, context := range existing.Contexts {
		if context.Cluster == clusterName {
			delete(existing.Contexts, name)
			logger.Debug("removed context for %q from kubeconfig", name)
			isChanged = true
			if _, ok := existing.AuthInfos[name]; ok {
				delete(existing.AuthInfos, name)
				logger.Debug("removed user for %q from kubeconfig", name)
			}
			currentContextName = name
			break
		}
	}

	if existing.CurrentContext == currentContextName {
		existing.CurrentContext = ""
		logger.Debug("reset current-context in kubeconfig", currentContextName)
		isChanged = true
	}

	if parts := strings.Split(existing.CurrentContext, "@"); len(parts) == 2 {
		if strings.HasSuffix(parts[1], "eksctl.io") {
			if _, ok := existing.Contexts[existing.CurrentContext]; !ok {
				existing.CurrentContext = ""
				logger.Debug("reset stale current-context in kubeconfig", currentContextName)
				isChanged = true
			}
		}
	}

	return isChanged
}

// LookupAuthenticator looks up an available authenticator
func LookupAuthenticator() (string, bool) {
	for _, cmd := range AuthenticatorCommands() {
		_, err := exec.LookPath(cmd)
		if err == nil {
			return cmd, true
		}
	}
	return "", false
}
