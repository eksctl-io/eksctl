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

	"k8s.io/kops/pkg/pki"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/utils"
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
	logger.Info("SSH public key file %q does not exist; will assume existing EC2 key pair", ng.SSH.PublicKeyPath)
	existing, err := c.getKeyPair(ng.SSH.PublicKeyPath)
	if err != nil {
		return err
	}
	ng.SSH.PublicKeyName = *existing.KeyName
	logger.Info("found EC2 key pair %q", ng.SSH.PublicKeyName)
	return nil
}

func (c *ClusterProvider) importSSHPublicKeyIfNeeded(clusterName string, ng *api.NodeGroup) error {
	fingerprint, err := pki.ComputeAWSKeyFingerprint(string(ng.SSH.PublicKey))
	if err != nil {
		return err
	}
	ng.SSH.PublicKeyName = c.getKeyPairName(clusterName, ng, &fingerprint)
	existing, err := c.getKeyPair(ng.SSH.PublicKeyName)
	if err != nil {
		if strings.HasPrefix(err.Error(), "cannot find EC2 key pair") {
			input := &ec2.ImportKeyPairInput{
				KeyName:           &ng.SSH.PublicKeyName,
				PublicKeyMaterial: ng.SSH.PublicKey,
			}
			logger.Info("importing SSH public key %q as %q", ng.SSH.PublicKeyPath, ng.SSH.PublicKeyName)
			if _, err = c.Provider.EC2().ImportKeyPair(input); err != nil {
				return errors.Wrap(err, "importing SSH public key")
			}
			return nil
		}
		return errors.Wrap(err, "checking existing key pair")
	}
	if *existing.KeyFingerprint != fingerprint {
		return fmt.Errorf("SSH public key %s already exists, but fingerprints don't match (exected: %q, got: %q)", ng.SSH.PublicKeyName, fingerprint, *existing.KeyFingerprint)
	}
	logger.Debug("SSH public key %s already exists", ng.SSH.PublicKeyName)
	return nil
}

// LoadSSHPublicKey loads the given SSH public key
func (c *ClusterProvider) LoadSSHPublicKey(clusterName string, ng *api.NodeGroup) error {
	if !ng.SSH.Allow {
		// TODO: https://github.com/weaveworks/eksctl/issues/144
		return nil
	}
	ng.SSH.PublicKeyPath = utils.ExpandPath(ng.SSH.PublicKeyPath)
	sshPublicKey, err := ioutil.ReadFile(ng.SSH.PublicKeyPath)
	if err != nil {
		if os.IsNotExist(err) {
			// if file not found – try to use existing EC2 key pair
			return c.tryExistingSSHPublicKeyFromPath(ng)
		}
		return errors.Wrap(err, fmt.Sprintf("reading SSH public key file %q", ng.SSH.PublicKeyPath))
	}
	// on successful read – import it
	ng.SSH.PublicKey = sshPublicKey
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
