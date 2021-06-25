package cmdutils

import (
	"github.com/pkg/errors"
)

func NewUtilsKMSLoader(cmd *Cmd) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)
	l.flagsIncompatibleWithConfigFile.Insert(
		"key-arn",
		"cluster",
	)

	l.validateWithConfigFile = func() error {
		if cmd.NameArg != "" {
			return errors.Errorf("config file and key ARN %s", IncompatibleFlags)
		}
		if l.ClusterConfig.SecretsEncryption == nil || l.ClusterConfig.SecretsEncryption.KeyARN == "" {
			return errors.New("field secretsEncryption.keyARN is required")
		}
		return nil
	}

	l.validateWithoutConfigFile = func() error {
		if l.ClusterConfig.SecretsEncryption == nil || l.ClusterConfig.SecretsEncryption.KeyARN == "" {
			return errors.New("--key-arn is required when a config file is not specified")
		}
		return nil
	}
	return l
}
