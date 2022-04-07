package client

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/smithy-go"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"golang.org/x/crypto/ssh"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"k8s.io/kops/pkg/pki"

	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/utils/file"
)

// LoadKeyFromFile loads and imports a public SSH key from a file provided a path to that file.
// returns the name of the key
func LoadKeyFromFile(ctx context.Context, filePath, clusterName, ngName string, ec2API awsapi.EC2) (string, error) {
	if !file.Exists(filePath) {
		return "", fmt.Errorf("SSH public key file %q not found", filePath)
	}

	expandedPath := file.ExpandPath(filePath)
	fileContent, err := readFileContents(expandedPath)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("reading SSH public key file %q", filePath))
	}

	fingerprint, err := fingerprint(filePath, fileContent)
	if err != nil {
		return "", err
	}

	key := string(fileContent)
	keyName := getKeyName(clusterName, ngName, fingerprint)
	logger.Info("using SSH public key %q as %q ", expandedPath, keyName)

	// Import SSH key in EC2
	if err := importKey(ctx, keyName, fingerprint, &key, ec2API); err != nil {
		return "", err
	}
	return keyName, nil
}

func fingerprint(filePath string, key []byte) (string, error) {
	pk, _, _, _, err := ssh.ParseAuthorizedKey(key)
	if err != nil {
		return "", fmt.Errorf("parsing key %q: %w", filePath, err)
	}

	if pk.Type() == "ssh-ed25519" {
		return strings.TrimPrefix(ssh.FingerprintSHA256(pk), "SHA256:"), nil
	}

	fingerprint, err := pki.ComputeAWSKeyFingerprint(string(key))
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("computing fingerprint for key %q", filePath))
	}

	return fingerprint, nil
}

// LoadKeyByContent loads and imports an SSH public key into EC2 if it doesn't exist
func LoadKeyByContent(ctx context.Context, key *string, clusterName, ngName string, ec2API awsapi.EC2) (string, error) {
	fingerprint, err := pki.ComputeAWSKeyFingerprint(*key)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("computing fingerprint for key %q", *key))
	}
	keyName := getKeyName(clusterName, ngName, fingerprint)

	logger.Info("using SSH public key %q ", *key)

	// Import SSH key in EC2
	if err := importKey(ctx, keyName, fingerprint, key, ec2API); err != nil {
		return "", err
	}
	return keyName, nil
}

// DeleteKeys will delete the public SSH key, if it exists
func DeleteKeys(ctx context.Context, ec2API awsapi.EC2, clusterName string) {
	existing, err := ec2API.DescribeKeyPairs(ctx, &ec2.DescribeKeyPairsInput{})
	if err != nil {
		logger.Debug("cannot describe keys: %v", err)
		return
	}
	var matching []*string
	prefix := getKeyName(clusterName, "", "")
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
		if _, err := ec2API.DeleteKeyPair(ctx, input); err != nil {
			logger.Debug("key pair couldn't be deleted: %v", err)
		}
	}
}

// CheckKeyExistsInEC2 returns whether a public ssh key already exists in EC2 or error if it couldn't be checked
func CheckKeyExistsInEC2(ctx context.Context, ec2API awsapi.EC2, sshKeyName string) error {
	existing, err := findKeyInEC2(ctx, ec2API, sshKeyName)
	if err != nil {
		return errors.Wrap(err, "checking existing key pair")
	}
	if existing == nil {
		return fmt.Errorf("cannot find EC2 key pair %q", sshKeyName)

	}
	return nil
}

func importKey(ctx context.Context, keyName, fingerprint string, keyContent *string, ec2API awsapi.EC2) error {
	if existing, err := findKeyInEC2(ctx, ec2API, keyName); err != nil {
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

	if _, err := ec2API.ImportKeyPair(ctx, input); err != nil {
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

func findKeyInEC2(ctx context.Context, ec2API awsapi.EC2, name string) (*ec2types.KeyPairInfo, error) {
	input := &ec2.DescribeKeyPairsInput{
		KeyNames: []string{name},
	}
	output, err := ec2API.DescribeKeyPairs(ctx, input)

	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) && ae.ErrorCode() == "InvalidKeyPair.NotFound" {
			return nil, nil
		}
		return nil, errors.Wrapf(err, fmt.Sprintf("searching for SSH public key %q in EC2", name))
	}

	if len(output.KeyPairs) != 1 {
		logger.Debug("output = %#v", output)
		return nil, fmt.Errorf("unexpected number of key pairs found (expected: 1, got: %d)", len(output.KeyPairs))
	}
	return &output.KeyPairs[0], nil
}

func readFileContents(filePath string) ([]byte, error) {
	fileContents, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("reading SSH public key file %q", filePath))
	}
	return fileContents, nil
}
