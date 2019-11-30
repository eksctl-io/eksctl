package cmdutils

import (
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/fargate"
)

const (
	fargateProfileName              = "name"      // Fargate profile name.
	fargateProfileSelectorNamespace = "namespace" // Fargate profile selector's namespace.
	fargateProfileSelectorLabels    = "labels"    // Fargate profile selector's labels.
)

// AddFlagsForFargate configures the flags required to interact with Fargate.
func AddFlagsForFargate(fs *pflag.FlagSet, options *fargate.Options) {
	addFargateProfileName(fs, &options.ProfileName)
}

// AddFlagsForFargateProfileCreation configures the flags required to
// create a Fargate profile.
func AddFlagsForFargateProfileCreation(fs *pflag.FlagSet, options *fargate.CreateOptions) {
	addFargateProfileName(fs, &options.ProfileName)

	fs.StringVar(&options.ProfileSelectorNamespace, fargateProfileSelectorNamespace, "",
		"Kubernetes namespace of the workloads to schedule on Fargate")

	fs.StringToStringVarP(&options.ProfileSelectorLabels, fargateProfileSelectorLabels, "l", nil,
		"Kubernetes selector labels of the workloads to schedule on Fargate, e.g. k1=v1,k2=v2")
}

func addFargateProfileName(fs *pflag.FlagSet, profileName *string) {
	fs.StringVar(profileName, fargateProfileName, "",
		"Fargate profile's name")
}

// Flags which should NOT be provided when a ClusterConfig file is also
// provided, in order to prevent duplicated, conflicting input.
var fargateProfileFlagsIncompatibleWithConfigFile = []string{
	fargateProfileName,
	fargateProfileSelectorNamespace,
	fargateProfileSelectorLabels,
}

// Flags which also require a ClusterConfig file to be provided.
var fargateProfileFlagsIncompatibleWithoutConfigFile = []string{}

// NewCreateFargateProfileLoader will load config or use flags for
// 'eksctl create fargateprofile'
func NewCreateFargateProfileLoader(cmd *Cmd, options *fargate.CreateOptions) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)
	l.flagsIncompatibleWithConfigFile.Insert(fargateProfileFlagsIncompatibleWithConfigFile...)
	l.flagsIncompatibleWithoutConfigFile.Insert(fargateProfileFlagsIncompatibleWithoutConfigFile...)
	l.validateWithConfigFile = func() error {
		return validateFargateProfiles(l)
	}
	l.validateWithoutConfigFile = func() error {
		if err := l.validateMetadataWithoutConfigFile(); err != nil {
			return err
		}
		if err := options.Validate(); err != nil {
			return err
		}
		cmd.ClusterConfig.FargateProfiles = []*api.FargateProfile{
			options.ToFargateProfile(),
		}
		return validateFargateProfiles(l)
	}
	return l
}

func validateFargateProfiles(l *commonClusterConfigLoader) error {
	for _, profile := range l.ClusterConfig.FargateProfiles {
		if err := profile.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// NewDeleteFargateProfileLoader will load config or use flags for
// 'eksctl delete fargateprofile'
func NewDeleteFargateProfileLoader(cmd *Cmd, options *fargate.Options) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)
	l.flagsIncompatibleWithConfigFile.Insert(fargateProfileFlagsIncompatibleWithConfigFile...)
	l.flagsIncompatibleWithoutConfigFile.Insert(fargateProfileFlagsIncompatibleWithoutConfigFile...)
	l.validateWithoutConfigFile = func() error {
		if err := l.validateMetadataWithoutConfigFile(); err != nil {
			return err
		}
		return options.Validate()
	}
	return l
}
