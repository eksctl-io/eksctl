package eks

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/kubernetes-sigs/aws-iam-authenticator/pkg/token"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/kops/pkg/pki"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

func (c *ClusterProvider) getKeyPairName(clusterName string, ng *api.NodeGroup, fingerprint *string) string {
	keyNameParts := []string{"eksctl", clusterName}
	if ng != nil {
		keyNameParts = append(keyNameParts, fmt.Sprintf("nodegroup-%s", ng.Name))
	}
	if fingerprint != nil {
		keyNameParts = append(keyNameParts, *fingerprint)
	}
	return strings.Join(keyNameParts, "-")
}

func (c *ClusterProvider) getKeyPair(name string) (*ec2.KeyPairInfo, error) {
	input := &ec2.DescribeKeyPairsInput{
		KeyNames: aws.StringSlice([]string{name}),
	}
	output, err := c.Provider.EC2().DescribeKeyPairs(input)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot find EC2 key pair %q", name)
	}
	if len(output.KeyPairs) != 1 {
		logger.Debug("output = %#v", output)
		return nil, fmt.Errorf("unexpected number of key pairs found (expected: 1, got: %d)", len(output.KeyPairs))
	}
	return output.KeyPairs[0], nil
}

func (c *ClusterProvider) tryExistingSSHPublicKeyFromPath(ng *api.NodeGroup) error {
	logger.Info("SSH public key file %q does not exist; will assume existing EC2 key pair", ng.SSHPublicKeyPath)
	existing, err := c.getKeyPair(ng.SSHPublicKeyPath)
	if err != nil {
		return err
	}
	ng.SSHPublicKeyName = *existing.KeyName
	logger.Info("found EC2 key pair %q", ng.SSHPublicKeyName)
	return nil
}

func (c *ClusterProvider) importSSHPublicKeyIfNeeded(clusterName string, ng *api.NodeGroup) error {
	fingerprint, err := pki.ComputeAWSKeyFingerprint(string(ng.SSHPublicKey))
	if err != nil {
		return err
	}
	ng.SSHPublicKeyName = c.getKeyPairName(clusterName, ng, &fingerprint)
	existing, err := c.getKeyPair(ng.SSHPublicKeyName)
	if err != nil {
		if strings.HasPrefix(err.Error(), "cannot find EC2 key pair") {
			input := &ec2.ImportKeyPairInput{
				KeyName:           &ng.SSHPublicKeyName,
				PublicKeyMaterial: ng.SSHPublicKey,
			}
			logger.Info("importing SSH public key %q as %q", ng.SSHPublicKeyPath, ng.SSHPublicKeyName)
			if _, err = c.Provider.EC2().ImportKeyPair(input); err != nil {
				return errors.Wrap(err, "importing SSH public key")
			}
			return nil
		}
		return errors.Wrap(err, "checking existing key pair")
	}
	if *existing.KeyFingerprint != fingerprint {
		return fmt.Errorf("SSH public key %s already exists, but fingerprints don't match (exected: %q, got: %q)", ng.SSHPublicKeyName, fingerprint, *existing.KeyFingerprint)
	}
	logger.Debug("SSH public key %s already exists", ng.SSHPublicKeyName)
	return nil
}

// LoadSSHPublicKey loads the given SSH public key
func (c *ClusterProvider) LoadSSHPublicKey(clusterName string, ng *api.NodeGroup) error {
	if !ng.AllowSSH {
		// TODO: https://github.com/weaveworks/eksctl/issues/144
		return nil
	}
	ng.SSHPublicKeyPath = utils.ExpandPath(ng.SSHPublicKeyPath)
	sshPublicKey, err := ioutil.ReadFile(ng.SSHPublicKeyPath)
	if err != nil {
		if os.IsNotExist(err) {
			// if file not found – try to use existing EC2 key pair
			return c.tryExistingSSHPublicKeyFromPath(ng)
		}
		return errors.Wrap(err, fmt.Sprintf("reading SSH public key file %q", ng.SSHPublicKeyPath))
	}
	// on successful read – import it
	ng.SSHPublicKey = sshPublicKey
	if err := c.importSSHPublicKeyIfNeeded(clusterName, ng); err != nil {
		return err
	}
	return nil
}

// MaybeDeletePublicSSHKey will delete the public SSH key, if it exists
func (c *ClusterProvider) MaybeDeletePublicSSHKey(clusterName string) {
	existing, err := c.Provider.EC2().DescribeKeyPairs(&ec2.DescribeKeyPairsInput{})
	if err != nil {
		logger.Debug("cannot describe keys: %v", err)
		return
	}
	matching := []*string{}
	prefix := c.getKeyPairName(clusterName, nil, nil)
	logger.Debug("existing = %#v", existing)
	for _, e := range existing.KeyPairs {
		if strings.HasPrefix(*e.KeyName, prefix) {
			nameParts := strings.Split(*e.KeyName, "-")
			logger.Debug("existing key %q matches prefix", *e.KeyName)
			if nameParts[len(nameParts)-1] == *e.KeyFingerprint {
				logger.Debug("existing key %q matches fingerprint", *e.KeyName)
				matching = append(matching, e.KeyName)
			}
		}
	}
	for i := range matching {
		input := &ec2.DeleteKeyPairInput{
			KeyName: matching[i],
		}
		logger.Debug("deleting key %q", *matching[i])
		if _, err := c.Provider.EC2().DeleteKeyPair(input); err != nil {
			logger.Debug("key pair couldn't be deleted: %v", err)
		}
	}
}

func (c *ClusterProvider) getUsername() string {
	usernameParts := strings.Split(c.Status.iamRoleARN, "/")
	if len(usernameParts) > 1 {
		return usernameParts[len(usernameParts)-1]
	}
	return "iam-root-account"
}

// ClientConfig stores information about the client config
type ClientConfig struct {
	Client      *clientcmdapi.Config
	ContextName string
}

// NewClientConfig creates a new client config, if withEmbeddedToken is true
// it will embed the STS token, otherwise it will use authenticator exec plugin
// and ensures that AWS_PROFILE environment variable gets set also
func (c *ClusterProvider) NewClientConfig(spec *api.ClusterConfig, withEmbeddedToken bool) (*ClientConfig, error) {
	client, _, contextName := kubeconfig.New(spec, c.getUsername(), "")

	config := &ClientConfig{
		Client:      client,
		ContextName: contextName,
	}

	if withEmbeddedToken {
		if err := config.useEmbeddedToken(spec, c.Provider.STS().(*sts.STS)); err != nil {
			return nil, err
		}
	} else {
		kubeconfig.AppendAuthenticator(config.Client, spec, utils.DetectAuthenticator(), c.Provider.Profile())
	}

	return config, nil
}

func (c *ClientConfig) useEmbeddedToken(spec *api.ClusterConfig, sts *sts.STS) error {
	gen, err := token.NewGenerator(true)
	if err != nil {
		return errors.Wrap(err, "could not get token generator")
	}

	tok, err := gen.GetWithSTS(spec.Metadata.Name, sts)
	if err != nil {
		return errors.Wrap(err, "could not get token")
	}

	c.Client.AuthInfos[c.ContextName].Token = tok
	return nil
}

// NewClientSet creates a new API client
func (c *ClientConfig) NewClientSet() (*kubernetes.Clientset, error) {
	clientConfig, err := clientcmd.NewDefaultClientConfig(*c.Client, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create API client configuration from client config")
	}

	client, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create API client")
	}
	return client, nil
}

// NewStdClientSet creates a new API client in one go with an embedded STS token, this is most commonly used option
func (c *ClusterProvider) NewStdClientSet(spec *api.ClusterConfig) (*kubernetes.Clientset, error) {
	clientConfig, err := c.NewClientConfig(spec, true)
	if err != nil {
		return nil, errors.Wrap(err, "creating Kubernetes client config with embedded token")
	}

	clientSet, err := clientConfig.NewClientSet()
	if err != nil {
		return nil, errors.Wrap(err, "creating Kubernetes client")
	}
	return clientSet, nil
}
