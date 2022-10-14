package kubeconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gofrs/flock"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/eksctl/pkg/utils/file"
)

const (
	// AWSIAMAuthenticator defines the name of the AWS IAM authenticator
	AWSIAMAuthenticator = "aws-iam-authenticator"
	// HeptioAuthenticatorAWS defines the old name of AWS IAM authenticator
	HeptioAuthenticatorAWS = "heptio-authenticator-aws"
	// AWSEKSAuthenticator defines the recently added `aws eks get-token` command
	AWSEKSAuthenticator = "aws"
	// AWSIAMAuthenticatorMinimumBetaVersion this is the minimum version at which aws-iam-authenticator uses v1beta1 as APIVersion
	AWSIAMAuthenticatorMinimumBetaVersion = "0.5.3"
	// AWSCLIv1MinimumBetaVersion this is the minimum version at which aws-cli v1 uses v1beta1 as APIVersion
	AWSCLIv1MinimumBetaVersion = "1.23.9"
	// AWSCLIv2MinimumBetaVersion this is the minimum version at which aws-cli v2 uses v1beta1 as APIVersion
	AWSCLIv2MinimumBetaVersion = "2.6.3"

	alphaAPIVersion = "client.authentication.k8s.io/v1alpha1"
	betaAPIVersion  = "client.authentication.k8s.io/v1beta1"
)

type ExecCommandFunc func(name string, arg ...string) *exec.Cmd

var execCommand = exec.Command

// DefaultPath defines the default path
func DefaultPath() string {
	if env := os.Getenv(clientcmd.RecommendedConfigPathEnvVar); len(env) > 0 {
		return env
	}
	return clientcmd.RecommendedHomeFile
}

// AuthenticatorCommands returns all of authenticator commands
func AuthenticatorCommands() []string {
	return []string{
		AWSIAMAuthenticator,
		HeptioAuthenticatorAWS,
		AWSEKSAuthenticator,
	}
}

// ConfigBuilder can create a client-go clientcmd Config
type ConfigBuilder struct {
	cluster     clientcmdapi.Cluster
	clusterName string
	username    string
}

// Build creates the Config with the ConfigBuilder settings
func (cb *ConfigBuilder) Build() *clientcmdapi.Config {
	contextName := fmt.Sprintf("%s@%s", cb.username, cb.clusterName)
	return &clientcmdapi.Config{
		Clusters: map[string]*clientcmdapi.Cluster{
			cb.clusterName: &cb.cluster,
		},
		Contexts: map[string]*clientcmdapi.Context{
			contextName: {
				Cluster:  cb.clusterName,
				AuthInfo: contextName,
			},
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			contextName: {},
		},
		CurrentContext: contextName,
	}
}

// NewBuilder returns a minimal ConfigBuilder
func NewBuilder(metadata *api.ClusterMeta, status *api.ClusterStatus, username string) *ConfigBuilder {
	cluster := clientcmdapi.Cluster{
		Server:                   status.Endpoint,
		CertificateAuthorityData: status.CertificateAuthorityData,
	}
	return &ConfigBuilder{
		cluster:     cluster,
		clusterName: metadata.String(),
		username:    username,
	}
}

// UseCertificateAuthorityFile sets the config to use CA from file for TLS
// communication instead of the CA retrieved from EKS
func (cb *ConfigBuilder) UseCertificateAuthorityFile(path string) *ConfigBuilder {
	cb.cluster.CertificateAuthority = path
	cb.cluster.CertificateAuthorityData = []byte{}
	return cb
}

// UseSystemCA sets the config to use the system CAs for TLS communication
// instead of the CA retrieved from EKS
func (cb *ConfigBuilder) UseSystemCA() *ConfigBuilder {
	cb.cluster.CertificateAuthority = ""
	cb.cluster.CertificateAuthorityData = []byte{}
	return cb
}

// ClusterInfo holds the cluster info.
type ClusterInfo interface {
	// ID returns the cluster ID.
	// This can either be the name of the cluster or a UUID.
	ID() string
	// Meta returns the cluster metadata.
	Meta() *api.ClusterMeta
	// GetStatus returns the cluster status.
	GetStatus() *api.ClusterStatus
}

// NewForUser returns a Config suitable for a user by respecting
// provider settings
func NewForUser(cluster ClusterInfo, username string) *clientcmdapi.Config {
	configBuilder := NewBuilder(cluster.Meta(), cluster.GetStatus(), username)
	if os.Getenv("KUBECONFIG_USE_SYSTEM_CA") != "" {
		configBuilder.UseSystemCA()
	}
	return configBuilder.Build()
}

// NewForKubectl creates configuration for a user with kubectl by configuring
// a suitable authenticator and respecting provider settings
func NewForKubectl(cluster ClusterInfo, username, roleARN, profile string) *clientcmdapi.Config {
	config := NewForUser(cluster, username)
	authenticator, found := LookupAuthenticator()
	if !found {
		// fall back to aws-iam-authenticator
		authenticator = AWSIAMAuthenticator
	}
	AppendAuthenticator(config, cluster, authenticator, roleARN, profile)
	return config
}

// AppendAuthenticator appends the AWS IAM  authenticator, and
// if profile is non-empty string it sets AWS_PROFILE environment
// variable also
func AppendAuthenticator(config *clientcmdapi.Config, cluster ClusterInfo, authenticatorCMD, roleARN, profile string) {
	var (
		args        []string
		roleARNFlag string
	)

	execConfig := &clientcmdapi.ExecConfig{
		APIVersion: alphaAPIVersion,
		Command:    authenticatorCMD,
		Env: []clientcmdapi.ExecEnvVar{
			{
				Name:  "AWS_STS_REGIONAL_ENDPOINTS",
				Value: "regional",
			},
		},
		ProvideClusterInfo: false,
	}

	meta := cluster.Meta()

	switch authenticatorCMD {
	case AWSIAMAuthenticator:
		// if version is above or equal to v0.5.3 we change the APIVersion to v1beta1.
		if authenticatorIsBetaVersion, err := authenticatorIsAboveVersion(AWSIAMAuthenticatorMinimumBetaVersion); err != nil {
			logger.Warning("failed to determine authenticator version, leaving API version as default v1alpha1: %v", err)
		} else if authenticatorIsBetaVersion {
			execConfig.APIVersion = betaAPIVersion
		}
		args = []string{"token", "-i", cluster.ID()}
		roleARNFlag = "-r"
		if meta.Region != "" {
			execConfig.Env = append(execConfig.Env, clientcmdapi.ExecEnvVar{
				Name:  "AWS_DEFAULT_REGION",
				Value: meta.Region,
			})
		}
	case HeptioAuthenticatorAWS:
		args = []string{"token", "-i", cluster.ID()}
		roleARNFlag = "-r"
		if meta.Region != "" {
			execConfig.Env = append(execConfig.Env, clientcmdapi.ExecEnvVar{
				Name:  "AWS_DEFAULT_REGION",
				Value: meta.Region,
			})
		}
	case AWSEKSAuthenticator:
		// if [aws-cli v1/aws-cli v2] is above or equal to [v1.23.9/v2.6.3] respectively, we change the APIVersion to v1beta1.
		if awsCLIIsBetaVersion, err := awsCliIsAboveVersion(); err != nil {
			logger.Warning("failed to determine authenticator version, leaving API version as default v1alpha1: %v", err)
		} else if awsCLIIsBetaVersion {
			execConfig.APIVersion = betaAPIVersion
		}
		args = []string{"eks", "get-token", "--cluster-name", cluster.ID()}
		roleARNFlag = "--role-arn"
		if meta.Region != "" {
			args = append(args, "--region", meta.Region)
		}
	}
	// If the alpha API version is selected, check the kubectl version
	// If kubectl 1.24.0 or above is detected, override with the beta API version
	// kubectl 1.24.0 removes the alpha API version, so it will never work
	// Therefore as a best effort try the beta version even if it might not work
	if execConfig.APIVersion == alphaAPIVersion {
		if kubectlVersion := getKubectlVersion(); kubectlVersion != "" {
			// Silently ignore errors because kubectl is not required to run eksctl
			compareVersions, err := utils.CompareVersions(kubectlVersion, "1.24.0")
			if err == nil && compareVersions >= 0 {
				execConfig.APIVersion = betaAPIVersion
			}
		}
	}
	if roleARN != "" {
		args = append(args, roleARNFlag, roleARN)
	}

	execConfig.Args = args

	if profile != "" {
		execConfig.Env = append(execConfig.Env, clientcmdapi.ExecEnvVar{
			Name:  "AWS_PROFILE",
			Value: profile,
		})
	}

	config.AuthInfos[config.CurrentContext] = &clientcmdapi.AuthInfo{
		Exec: execConfig,
	}
}

// AWSAuthenticatorVersionFormat is the format in which aws-iam-authenticator displays version information:
// {"Version":"0.5.5","Commit":"85e50980d9d916ae95882176c18f14ae145f916f"}
type AWSAuthenticatorVersionFormat struct {
	Version string `json:"Version"`
}

func awsCliIsAboveVersion() (bool, error) {
	awsCliVersion := getAWSCLIVersion()
	compareVersions, err := utils.CompareVersions(awsCliVersion, "2.0.0")
	if err != nil {
		return false, fmt.Errorf("failed to parse versions: %w", err)
	}
	// AWS CLI provides beta in two separate major versions. One for v1 and one for v2.
	// Being above a single version doesn't necessarily mean that we are in the clear.
	// Thus, first check which major version we are dealing with, then check if beta is
	// supported for that major version.
	if compareVersions < 0 {
		compareVersions, err = utils.CompareVersions(awsCliVersion, AWSCLIv1MinimumBetaVersion)
	} else {
		compareVersions, err = utils.CompareVersions(awsCliVersion, AWSCLIv2MinimumBetaVersion)
	}
	if err != nil {
		return false, fmt.Errorf("failed to parse versions: %w", err)
	}
	return compareVersions >= 0, nil
}

func getAWSCLIVersion() string {
	cmd := execCommand("aws", "--version")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	r := regexp.MustCompile(`aws-cli/([\d.]*)`)
	matches := r.FindStringSubmatch(string(output))
	if len(matches) != 2 {
		return ""
	}
	return matches[1]
}

func authenticatorIsAboveVersion(version string) (bool, error) {
	authenticatorVersion, err := getAWSIAMAuthenticatorVersion()
	if err != nil {
		return false, fmt.Errorf("failed to retrieve authenticator version: %w", err)
	}
	compareVersions, err := utils.CompareVersions(authenticatorVersion, version)
	if err != nil {
		return false, fmt.Errorf("failed to parse versions: %w", err)
	}
	return compareVersions >= 0, nil
}

func getAWSIAMAuthenticatorVersion() (string, error) {
	cmd := execCommand(AWSIAMAuthenticator, "version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to run aws-iam-authenticator version command: %w", err)
	}
	var parsedVersion AWSAuthenticatorVersionFormat
	if err := json.Unmarshal(output, &parsedVersion); err != nil {
		return "", fmt.Errorf("failed to parse version information: %w", err)
	}
	return parsedVersion.Version, nil
}

/* KubectlVersionFormat is the format used by kubectl version --format=json, example output:
{
  "clientVersion": {
    "major": "1",
    "minor": "23",
    "gitVersion": "v1.23.6",
    "gitCommit": "ad3338546da947756e8a88aa6822e9c11e7eac22",
    "gitTreeState": "archive",
    "buildDate": "2022-04-29T06:39:16Z",
    "goVersion": "go1.18.1",
    "compiler": "gc",
    "platform": "linux/amd64"
  }
} */
type KubectlVersionData struct {
	Version string `json:"gitVersion"`
}

type KubectlVersionFormat struct {
	ClientVersion KubectlVersionData `json:"clientVersion"`
}

func getKubectlVersion() string {
	cmd := execCommand("kubectl", "version", "--client", "--output=json")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	var parsedVersion KubectlVersionFormat
	if err := json.Unmarshal(output, &parsedVersion); err != nil {
		return ""
	}
	return strings.TrimLeft(parsedVersion.ClientVersion.Version, "v")
}

func lockFileName(filePath string) string {
	return filePath + ".eksctl.lock"
}

// ensureDirectory should probably be handled in flock
func ensureDirectory(filePath string) error {
	dir := filepath.Dir(filePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func lockConfigFile(filePath string) (*flock.Flock, error) {
	lockFileName := lockFileName(filePath)
	// Make sure the directory exists, otherwise flock fails
	if err := ensureDirectory(lockFileName); err != nil {
		return nil, err
	}
	flock := flock.New(lockFileName)
	err := flock.Lock()
	if err != nil {
		return nil, errors.Wrap(err, "failed to obtain exclusive lock on kubeconfig lockfile")
	}

	return flock, nil
}

func unlockConfigFile(fl *flock.Flock) error {
	err := fl.Unlock()
	if err != nil {
		return errors.Wrap(err, "failed to release exclusive lock on kubeconfig lockfile")
	}

	return nil
}

// Write will write Kubernetes client configuration to a file.
// If path isn't specified then the path will be determined by client-go.
// If file pointed to by path doesn't exist it will be created.
// If the file already exists then the configuration will be merged with the existing file.
func Write(path string, newConfig clientcmdapi.Config, setContext bool) (string, error) {
	configAccess := getConfigAccess(path)
	configFileName := configAccess.GetDefaultFilename()
	fl, err := lockConfigFile(configFileName)
	if err != nil {
		return "", err
	}

	defer func() {
		if err := unlockConfigFile(fl); err != nil {
			logger.Critical(err.Error())
		}
	}()

	config, err := configAccess.GetStartingConfig()
	if err != nil {
		return "", errors.Wrapf(err, "unable to read existing kubeconfig file %q", path)
	}

	logger.Debug("merging kubeconfig files")
	merged := merge(config, &newConfig)

	if setContext && newConfig.CurrentContext != "" {
		logger.Debug("setting current-context to %s", newConfig.CurrentContext)
		merged.CurrentContext = newConfig.CurrentContext
	}

	if err := clientcmd.ModifyConfig(configAccess, *merged, true); err != nil {
		return "", errors.Wrapf(err, "unable to modify kubeconfig %s", path)
	}

	return configFileName, nil
}

func getConfigAccess(explicitPath string) clientcmd.ConfigAccess {
	pathOptions := clientcmd.NewDefaultPathOptions()
	if explicitPath != "" && explicitPath != DefaultPath() {
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
func MaybeDeleteConfig(meta *api.ClusterMeta) {
	p := AutoPath(meta.Name)

	if file.Exists(p) {
		fl, err := lockConfigFile(p)
		if err != nil {
			logger.Critical(err.Error())
			return
		}

		defer func() {
			if err := unlockConfigFile(fl); err != nil {
				logger.Critical(err.Error())
			}
		}()

		if err := isValidConfig(p, meta.Name); err != nil {
			logger.Debug(err.Error())
			return
		}
		if err := os.Remove(p); err != nil {
			logger.Debug("ignoring error while removing auto-generated config file %q: %s", p, err.Error())
		}
		return
	}

	configAccess := getConfigAccess(DefaultPath())
	defaultFilename := configAccess.GetDefaultFilename()
	fl, err := lockConfigFile(defaultFilename)
	if err != nil {
		logger.Critical(err.Error())
		return
	}

	defer func() {
		if err := unlockConfigFile(fl); err != nil {
			logger.Critical(err.Error())
		}
	}()

	config, err := configAccess.GetStartingConfig()
	if err != nil {
		logger.Debug("error reading kubeconfig file %q: %s", DefaultPath(), err.Error())
		return
	}

	if !deleteClusterInfo(config, meta) {
		return
	}

	if err := clientcmd.ModifyConfig(configAccess, *config, true); err != nil {
		logger.Debug("ignoring error while failing to update config file %q: %s", DefaultPath(), err.Error())
	} else {
		logger.Success("kubeconfig has been updated")
	}
}

// deleteClusterInfo removes a cluster's information from the kubeconfig if the cluster name
// provided by ctl matches a eksctl-created cluster in the kubeconfig
// returns 'true' if the existing config has changes and 'false' otherwise
func deleteClusterInfo(existing *clientcmdapi.Config, meta *api.ClusterMeta) bool {
	isChanged := false
	clusterName := meta.String()

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
		logger.Debug("reset current-context in kubeconfig to %q", currentContextName)
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
