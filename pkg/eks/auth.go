package eks

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/heptio/authenticator/pkg/token"
	"github.com/kubicorn/kubicorn/pkg/logger"

	clientset "k8s.io/client-go/kubernetes"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"k8s.io/kops/upup/pkg/fi/utils"
)

func (c *CloudFormation) LoadSSHPublicKey() error {
	c.cfg.SSHPublicKeyPath = utils.ExpandPath(c.cfg.SSHPublicKeyPath)
	sshPublicKey, err := ioutil.ReadFile(c.cfg.SSHPublicKeyPath)
	if err != nil {
		if os.IsNotExist(err) {
			// if file not found – try to use existing EC2 key pair
			logger.Info("SSH public key file %q does not exist; will assume existing EC2 key pair", c.cfg.SSHPublicKeyPath)
			input := &ec2.DescribeKeyPairsInput{
				KeyNames: aws.StringSlice([]string{c.cfg.SSHPublicKeyPath}),
			}
			output, err := c.ec2.DescribeKeyPairs(input)
			if err != nil {
				return errors.Wrap(err, "cannot find EC2 key pair")
			}
			if len(output.KeyPairs) != 1 {
				logger.Debug("output = %#v", output)
				return fmt.Errorf("coulnd't find existing EC2 key pair")
			}
			c.cfg.keyName = *output.KeyPairs[0].KeyName
			logger.Info("found EC2 key pair %q with finger print %q", c.cfg.keyName, *output.KeyPairs[0].KeyFingerprint)
		} else {
			return errors.Wrap(err, fmt.Sprintf("error reading SSH public key file %q", c.cfg.SSHPublicKeyPath))
		}
	} else {
		// on successfull read – import it
		c.cfg.SSHPublicKey = sshPublicKey
		c.cfg.keyName = "EKS-" + c.cfg.ClusterName
		input := &ec2.ImportKeyPairInput{
			KeyName:           &c.cfg.keyName,
			PublicKeyMaterial: c.cfg.SSHPublicKey,
		}
		logger.Info("importing SSH public key %q as %q", c.cfg.SSHPublicKeyPath, c.cfg.keyName)
		if _, err := c.ec2.ImportKeyPair(input); err != nil {
			return errors.Wrap(err, "importing SSH public key")
		}
	}
	return nil
}

func (c *CloudFormation) MaybeDeletePublicSSHKey() {
	input := &ec2.DeleteKeyPairInput{
		KeyName: aws.String("EKS-" + c.cfg.ClusterName),
	}
	c.ec2.DeleteKeyPair(input)
}

func (c *CloudFormation) getUsername() string {
	usernameParts := strings.Split(c.arn, "/")
	username := usernameParts[len(usernameParts)-1]
	return username
}

type ClientConfig struct {
	Client  *clientcmdapi.Config
	Cluster *Config
	roleARN string
}

// based on "k8s.io/kubernetes/cmd/kubeadm/app/util/kubeconfig"
// these are small, so we can copy these, and no need to deal with k/k as dependency
func (c *CloudFormation) NewClientConfig() (*ClientConfig, error) {
	clusterName := fmt.Sprintf("%s.%s.eksctl.io", c.cfg.ClusterName, c.cfg.Region)
	contextName := fmt.Sprintf("%s@%s", c.getUsername(), clusterName)

	clientConfig := &ClientConfig{
		Cluster: c.cfg,
		Client: &clientcmdapi.Config{
			Clusters: map[string]*clientcmdapi.Cluster{
				clusterName: {
					Server: c.cfg.MasterEndpoint,
					CertificateAuthorityData: c.cfg.CertificateAuthorityData,
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
		},
		roleARN: c.arn,
	}

	return clientConfig, nil
}

func (c *ClientConfig) WithExecHeptioAuthenticator() *ClientConfig {

	clientConfigCopy := *c
	x := clientConfigCopy.Client.AuthInfos[c.Client.CurrentContext]
	x.Exec = &clientcmdapi.ExecConfig{
		APIVersion: "client.authentication.k8s.io/v1alpha1",
		Command:    "heptio-authenticator-aws",
		Args:       []string{"token", "-i", c.Cluster.ClusterName},
		/*
			Args:       []string{"token", "-i", c.Cluster.ClusterName, "-r", c.roleARN},
		*/
	}
	return &clientConfigCopy
}

func (c *ClientConfig) WithEmbeddedToken() (*ClientConfig, error) {
	clientConfigCopy := *c

	gen, err := token.NewGenerator()
	if err != nil {
		return nil, errors.Wrap(err, "could not get token generator")
	}

	// could not get token: AccessDenied: User <ARN> is not authorized to perform: sts:AssumeRole on resource: <ARN>
	/*
		tok, err := gen.GetWithRole(c.Cluster.ClusterName, c.roleARN)
		if err != nil {
			return nil, errors.Wrap(err, "could not get token")
		}
	*/
	tok, err := gen.Get(c.Cluster.ClusterName)
	if err != nil {
		return nil, errors.Wrap(err, "could not get token")
	}

	x := c.Client.AuthInfos[c.Client.CurrentContext]
	x.Token = tok

	return &clientConfigCopy, nil
}

func (c *ClientConfig) WriteToFile(filename string) error {
	if err := clientcmd.WriteToFile(*c.Client, filename); err != nil {
		return errors.Wrapf(err, "couldn't write client config file %q", filename)
	}
	return nil
}

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
