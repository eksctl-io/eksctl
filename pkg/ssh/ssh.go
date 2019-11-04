package ssh

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/weaveworks/eksctl/pkg/utils/file"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"

	"k8s.io/kops/pkg/pki"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

// LoadKey loads the SSH public key specified in NodeGroupSSH and returns it. The key should be specified
// in only one way: by name (for a key existing in EC2), by path (for a key in a local file)
// or by its contents (in the config-file). It also assumes that if ssh is enabled (SSH.Allow
// == true) then one key was specified
func LoadKey(sshConfig *api.NodeGroupSSH, clusterName, nodeGroupName string, ec2API ec2iface.EC2API) (string, error) {
	if sshConfig.Allow == nil || *sshConfig.Allow == false {
		return "", nil
	}

	switch {

	// Load Key by content
	case sshConfig.PublicKey != nil:
		keyName, err := LoadKeyByContent(sshConfig.PublicKey, clusterName, nodeGroupName, ec2API)
		if err != nil {
			return "", err
		}
		return keyName, nil

		// Use key by name in EC2
	case sshConfig.PublicKeyName != nil && *sshConfig.PublicKeyName != "":
		if err := CheckKeyExistsInEC2(*sshConfig.PublicKeyName, ec2API); err != nil {
			return "", err
		}
		logger.Info("using EC2 key pair %q", *sshConfig.PublicKeyName)
		return *sshConfig.PublicKeyName, nil

		// Local ssh key file
	case file.Exists(*sshConfig.PublicKeyPath):
		keyName, err := LoadKeyFromFile(*sshConfig.PublicKeyPath, clusterName, nodeGroupName, ec2API)
		if err != nil {
			return "", err
		}
		return keyName, nil

		// A keyPath, when specified as a flag, can mean a local key (checked above) or a key name in EC2
	default:
		err := CheckKeyExistsInEC2(*sshConfig.PublicKeyPath, ec2API)
		if err != nil {
			return "", err
		}
		logger.Info("using EC2 key pair %q", sshConfig.PublicKeyName)
		return *sshConfig.PublicKeyPath, nil
	}

}

// LoadKeyFromFile loads and imports a public SSH key from a file provided a path to that file.
// returns the name of the key
func LoadKeyFromFile(filePath, clusterName, ngName string, ec2API ec2iface.EC2API) (string, error) {
	if !file.Exists(filePath) {
		return "", fmt.Errorf("SSH public key file %q not found", filePath)
	}

	expandedPath := file.ExpandPath(filePath)
	fileContent, err := readFileContents(expandedPath)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("reading SSH public key file %q", filePath))
	}

	key := string(fileContent)

	fingerprint, err := pki.ComputeAWSKeyFingerprint(key)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("computing fingerprint for key %q", filePath))
	}
	keyName := getKeyName(clusterName, ngName, fingerprint)

	logger.Info("using SSH public key %q as %q ", expandedPath, keyName)

	// Import SSH key in EC2
	if err := importKey(keyName, fingerprint, &key, ec2API); err != nil {
		return "", err
	}
	return keyName, nil
}

// LoadKeyByContent loads and imports an SSH public key into EC2 if it doesn't exist
func LoadKeyByContent(key *string, clusterName, ngName string, ec2API ec2iface.EC2API) (string, error) {
	fingerprint, err := pki.ComputeAWSKeyFingerprint(*key)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("computing fingerprint for key %q", *key))
	}
	keyName := getKeyName(clusterName, ngName, fingerprint)

	logger.Info("using SSH public key %q ", *key)

	// Import SSH key in EC2
	if err := importKey(keyName, fingerprint, key, ec2API); err != nil {
		return "", err
	}
	return keyName, nil
}

// DeleteKeys will delete the public SSH key, if it exists
func DeleteKeys(clusterName string, ec2API ec2iface.EC2API) {
	existing, err := ec2API.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{})
	if err != nil {
		logger.Debug("cannot describe keys: %v", err)
		return
	}
	var matching []*string
	prefix := getKeyName(clusterName, "", "")
	logger.Debug("existing = %#v", existing)
	for _, e := range existing.KeyPairs {
		if !strings.HasPrefix(*e.KeyName, prefix) {
			continue
		}
		nameParts := strings.Split(*e.KeyName, "-")
		logger.Debug("existing key %q matches prefix", *e.KeyName)
		if nameParts[len(nameParts)-1] == *e.KeyFingerprint {
			logger.Debug("existing key %q matches fingerprint", *e.KeyName)
			matching = append(matching, e.KeyName)
		}
	}
	for i := range matching {
		input := &ec2.DeleteKeyPairInput{
			KeyName: matching[i],
		}
		logger.Debug("deleting key %q", *matching[i])
		if _, err := ec2API.DeleteKeyPair(input); err != nil {
			logger.Debug("key pair couldn't be deleted: %v", err)
		}
	}
}

// CheckKeyExistsInEC2 returns whether a public ssh key already exists in EC2 or error if it couldn't be checked
func CheckKeyExistsInEC2(sshKeyName string, ec2API ec2iface.EC2API) error {
	existing, err := findKeyInEc2(sshKeyName, ec2API)
	if err != nil {
		return errors.Wrap(err, "checking existing key pair")
	}

	if existing == nil {
		return fmt.Errorf("cannot find EC2 key pair %q", sshKeyName)
	}

	return nil
}

func importKey(keyName, fingerprint string, keyContent *string, ec2API ec2iface.EC2API) error {
	if existing, err := findKeyInEc2(keyName, ec2API); err != nil {
		return err
	} else if existing != nil {
		if *existing.KeyFingerprint != fingerprint {
			return fmt.Errorf("SSH public key %s already exists, but fingerprints don't match (exected: %q, got: %q)", keyName, fingerprint, *existing.KeyFingerprint)
		}

		logger.Debug("SSH public key %s already exists", keyName)
		return nil
	}

	// Import it
	input := &ec2.ImportKeyPairInput{
		KeyName:           &keyName,
		PublicKeyMaterial: []byte(*keyContent),
	}
	logger.Debug("importing SSH public key %q", keyName)

	if _, err := ec2API.ImportKeyPair(input); err != nil {
		return errors.Wrap(err, "importing SSH public key")
	}
	return nil
}

// getKeyName generates the name of an SSH key based on the cluster name, nodegroup name and fingerprint
// in the form "eksctl-<clusterName>-nodegroup-<nodeGroupName>-<fingerprint>"
func getKeyName(clusterName, nodeGroupName, fingerprint string) string {
	keyNameParts := []string{"eksctl", clusterName}
	if nodeGroupName != "" {
		keyNameParts = append(keyNameParts, fmt.Sprintf("nodegroup-%s", nodeGroupName))
	}

	if fingerprint != "" {
		keyNameParts = append(keyNameParts, fingerprint)
	}

	return strings.Join(keyNameParts, "-")
}

func findKeyInEc2(name string, ec2API ec2iface.EC2API) (*ec2.KeyPairInfo, error) {
	input := &ec2.DescribeKeyPairsInput{
		KeyNames: aws.StringSlice([]string{name}),
	}
	output, err := ec2API.DescribeKeyPairs(input)

	if err != nil {
		awsError := err.(awserr.Error)
		if awsError.Code() == "InvalidKeyPair.NotFound" {
			return nil, nil
		}
		return nil, errors.Wrapf(err, fmt.Sprintf("searching for SSH public key %q in EC2", name))
	}

	if len(output.KeyPairs) != 1 {
		logger.Debug("output = %#v", output)
		return nil, fmt.Errorf("unexpected number of key pairs found (expected: 1, got: %d)", len(output.KeyPairs))
	}
	return output.KeyPairs[0], nil
}

func readFileContents(filePath string) ([]byte, error) {
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("reading SSH public key file %q", filePath))
	}
	return fileContents, nil
}
