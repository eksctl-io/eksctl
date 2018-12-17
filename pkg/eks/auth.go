package eks

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/kris-nova/logger"
	"github.com/kubernetes-sigs/aws-iam-authenticator/pkg/token"
	"github.com/pkg/errors"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/kops/pkg/pki"
)

func (c *ClusterProvider) getKeyPairName(clusterName string, ng *api.NodeGroup, fingerprint *string) string {
	keyNameParts := []string{"eksctl", clusterName}
	if ng != nil {
		keyNameParts = append(keyNameParts, fmt.Sprintf("ng%d", ng.ID))
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
	Cluster     *api.ClusterConfig
	ClusterName string
	ContextName string
	roleARN     string
	sts         stsiface.STSAPI
	profile     string
}

// NewClientConfig creates a new client config
// based on "k8s.io/kubernetes/cmd/kubeadm/app/util/kubeconfig"
// these are small, so we can copy these, and no need to deal with k/k as dependency
func (c *ClusterProvider) NewClientConfig(spec *api.ClusterConfig) (*ClientConfig, error) {
	client, clusterName, contextName := kubeconfig.New(spec, c.getUsername(), "")
	clientConfig := &ClientConfig{
		Cluster:     spec,
		Client:      client,
		ClusterName: clusterName,
		ContextName: contextName,
		roleARN:     c.Status.iamRoleARN,
		sts:         c.Provider.STS(),
		profile:     c.Provider.Profile(),
	}

	return clientConfig, nil
}

// WithExecAuthenticator creates a copy of ClientConfig with authenticator exec plugin
// it ensures that AWS_PROFILE environment variable gets added to config also
func (c *ClientConfig) WithExecAuthenticator() *ClientConfig {
	clientConfigCopy := *c

	kubeconfig.AppendAuthenticator(clientConfigCopy.Client, c.Cluster, utils.DetectAuthenticator())

	if len(c.profile) > 0 {
		clientConfigCopy.Client.AuthInfos[c.ContextName].Exec.Env = []clientcmdapi.ExecEnvVar{
			clientcmdapi.ExecEnvVar{
				Name:  "AWS_PROFILE",
				Value: c.profile,
			},
		}
	}

	return &clientConfigCopy
}

// WithEmbeddedToken embeds the STS token
func (c *ClientConfig) WithEmbeddedToken() (*ClientConfig, error) {
	clientConfigCopy := *c

	gen, err := token.NewGenerator(true)
	if err != nil {
		return nil, errors.Wrap(err, "could not get token generator")
	}

	tok, err := gen.GetWithSTS(c.Cluster.Metadata.Name, c.sts.(*sts.STS))
	if err != nil {
		return nil, errors.Wrap(err, "could not get token")
	}

	x := c.Client.AuthInfos[c.ContextName]
	x.Token = tok

	return &clientConfigCopy, nil
}

// NewClientSet creates a new API client
func (c *ClientConfig) NewClientSet() (*clientset.Clientset, error) {
	clientConfig, err := clientcmd.NewDefaultClientConfig(*c.Client, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create API client configuration from client config")
	}

	client, err := clientset.NewForConfig(clientConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create API client")
	}
	return client, nil
}

// NewClientSetWithEmbeddedToken creates a new API client with an embedded STS token
func (c *ClientConfig) NewClientSetWithEmbeddedToken() (*clientset.Clientset, error) {
	clientConfig, err := c.WithEmbeddedToken()
	if err != nil {
		return nil, errors.Wrap(err, "creating Kubernetes client config with embedded token")
	}
	clientSet, err := clientConfig.NewClientSet()
	if err != nil {
		return nil, errors.Wrap(err, "creating Kubernetes client")
	}
	return clientSet, nil
}
