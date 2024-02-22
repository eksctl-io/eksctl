package cmdutils

import "k8s.io/apimachinery/pkg/util/sets"

type ValidateCmdFunc = func(cmd *Cmd) error

type ConfigLoaderBuilder struct {
	FlagsIncompatibleWithConfigFile    sets.Set[string]
	FlagsIncompatibleWithoutConfigFile sets.Set[string]
	validateWithConfigFile             []ValidateCmdFunc
	validateWithoutConfigFile          []ValidateCmdFunc
	validate                           []ValidateCmdFunc
}

func (b *ConfigLoaderBuilder) ValidateWithoutConfigFile(f ValidateCmdFunc) {
	b.validateWithoutConfigFile = append(b.validateWithoutConfigFile, f)
}

func (b *ConfigLoaderBuilder) ValidateWithConfigFile(f ValidateCmdFunc) {
	b.validateWithConfigFile = append(b.validateWithConfigFile, f)
}
func (b *ConfigLoaderBuilder) Validate(f ValidateCmdFunc) {
	b.validate = append(b.validate, f)
}

func (b *ConfigLoaderBuilder) Build(cmd *Cmd) ClusterConfigLoader {
	return &commonClusterConfigLoader{
		Cmd:                                cmd,
		flagsIncompatibleWithConfigFile:    b.FlagsIncompatibleWithConfigFile,
		flagsIncompatibleWithoutConfigFile: b.FlagsIncompatibleWithoutConfigFile,
		validateWithConfigFile: func() error {
			for _, f := range append(b.validateWithConfigFile, b.validate...) {
				if err := f(cmd); err != nil {
					return err
				}
			}
			return nil
		},
		validateWithoutConfigFile: func() error {
			for _, f := range append(b.validateWithoutConfigFile, b.validate...) {
				if err := f(cmd); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func NewConfigLoaderBuilder() ConfigLoaderBuilder {
	return ConfigLoaderBuilder{
		FlagsIncompatibleWithConfigFile:    sets.New[string](defaultFlagsIncompatibleWithConfigFile[:]...),
		FlagsIncompatibleWithoutConfigFile: sets.New[string](defaultFlagsIncompatibleWithoutConfigFile[:]...),
		validateWithoutConfigFile: []func(cmd *Cmd) error{
			validateMetadataWithoutConfigFile,
		},
		validateWithConfigFile: []func(cmd *Cmd) error{},
	}
}
