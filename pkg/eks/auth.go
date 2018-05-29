package eks

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/ghodss/yaml"
	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"

	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"k8s.io/kops/upup/pkg/fi/utils"

	"github.com/heptio/authenticator/pkg/token"
)

func (c *Config) nodeAuthConfigMap() (*corev1.ConfigMap, error) {

	/*
		apiVersion: v1
		kind: ConfigMap
		metadata:
		  name: aws-auth
		  namespace: default
		data:
		  mapRoles: |
		    - rolearn: "${nodeInstanceRoleARN}"
		      username: system:node:{{EC2PrivateDNSName}}
		      groups:
		        - system:bootstrappers
		        - system:nodes
		        - system:node-proxier
	*/

	mapRoles := make([]map[string]interface{}, 1)
	mapRoles[0] = make(map[string]interface{})

	mapRoles[0]["rolearn"] = c.nodeInstanceRoleARN
	mapRoles[0]["username"] = "system:node:{{EC2PrivateDNSName}}"
	mapRoles[0]["groups"] = []string{
		"system:bootstrappers",
		"system:nodes",
		"system:nodes",
	}

	mapRolesBytes, err := yaml.Marshal(mapRoles)
	if err != nil {
		return nil, err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws-auth",
			Namespace: "default",
		},
		BinaryData: map[string][]byte{
			"mapRoles": mapRolesBytes,
		},
	}

	return cm, nil
}

// def generate_sts_token(name):
//     sts = setupSTSBoto()
//     prefix = "k8s-aws-v1."

//     signedURL = sts.generate_presigned_url(ClientMethod='get_caller_identity',  Params={}, ExpiresIn=60)
//     encodedURL = base64.b64encode(signedURL)

//     return prefix+encodedURL```

// the issue with boto is it doesn't allow you to generate a signed url with additional headers like the golang package so I have to rewrite this to manually sign the url

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

func (c *CloudFormation) getUsername() (string, error) {
	input := &sts.GetCallerIdentityInput{}
	output, err := c.sts.GetCallerIdentity(input)
	if err != nil {
		return "", errors.Wrap(err, "cannot get username for current session")
	}
	arn := strings.Split(*output.Arn, "/")
	username := arn[len(arn)-1]
	return username, nil
}

type ClientConfig struct {
	Client  *clientcmdapi.Config
	Cluster *Config
}

// based on "k8s.io/kubernetes/cmd/kubeadm/app/util/kubeconfig"
// these are small, so we can copy these, and no need to deal with k/k as dependency
func (c *CloudFormation) NewClientConfig() (*ClientConfig, error) {
	username, err := c.getUsername()
	if err != nil {
		return nil, err
	}

	clusterName := fmt.Sprintf("%s.%s.eksctl.k8s.io", c.cfg.ClusterName, c.cfg.Region)
	contextName := fmt.Sprintf("%s@%s", username, clusterName)

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
		// TODO(p2): should we add -r <role-arn> here?
	}
	return &clientConfigCopy
}

func (c *ClientConfig) WithEmbeddedToken() (*ClientConfig, error) {
	clientConfigCopy := *c

	gen, err := token.NewGenerator()
	if err != nil {
		return nil, errors.Wrap(err, "could not get token")
	}

	// otherwise sign the token with immediately available credentials
	tok, err := gen.Get(c.Cluster.ClusterName)
	if err != nil {
		return nil, errors.Wrap(err, "could not get token")
	}

	//c.Client.AuthInfos[c.Client.CurrentContext].

	logger.Debug("token = %#v", tok)
	logger.Debug("token = %s", gen.FormatJSON(tok))

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
