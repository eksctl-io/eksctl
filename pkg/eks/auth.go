package eks

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/ghodss/yaml"
	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/pkg/errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	_ "k8s.io/client-go/tools/clientcmd"
	_ "k8s.io/client-go/tools/clientcmd/api"

	//_ "k8s.io/kubernetes/cmd/kubeadm/app/util/kubeconfig"

	"k8s.io/kops/upup/pkg/fi/utils"

	_ "github.com/heptio/authenticator/pkg/token"
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
