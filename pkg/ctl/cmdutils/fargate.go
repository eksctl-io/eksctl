package cmdutils

import (
	"fmt"

	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/fargate"
	"k8s.io/apimachinery/pkg/util/sets"
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
		if err := validateCluster(cmd); err != nil {
			return err
		}
		if err := validateNameFlagAndArgCreate(cmd, options); err != nil {
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

func validateNameFlagAndArgCreate(cmd *Cmd, options *fargate.CreateOptions) error {
	if options.ProfileName != "" && cmd.NameArg != "" {
		return ErrFlagAndArg(fmt.Sprintf("--%s", fargateProfileName), options.ProfileName, cmd.NameArg)
	}
	if options.ProfileName == "" && cmd.NameArg != "" {
		options.ProfileName = cmd.NameArg
	}
	return nil
}

// NewGetFargateProfileLoader will load config or use flags for
// 'eksctl get fargateprofile'
func NewGetFargateProfileLoader(cmd *Cmd, options *fargate.Options) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)
	// We optionally want to be able to filter profiles by name:
	l.flagsIncompatibleWithConfigFile = flagsIncompatibleWithConfigFileExcept(fargateProfileName)
	l.flagsIncompatibleWithoutConfigFile.Insert(fargateProfileFlagsIncompatibleWithoutConfigFile...)
	l.validateWithoutConfigFile = func() error {
		return validate(cmd, options)
	}
	l.validateWithConfigFile = func() error {
		return validate(cmd, options)
	}
	return l
}

func flagsIncompatibleWithConfigFileExcept(items ...string) sets.String {
	set := sets.NewString(fargateProfileFlagsIncompatibleWithConfigFile...)
	set = set.Union(defaultFlagsIncompatibleWithConfigFile)
	set.Delete(items...)
	return set
}

func validate(cmd *Cmd, options *fargate.Options) error {
	if err := validateCluster(cmd); err != nil {
		return err
	}
	return validateNameFlagAndArg(cmd, options)
}

func validateCluster(cmd *Cmd) error {
	if cmd.ClusterConfig.Metadata.Name == "" {
		return ErrMustBeSet(ClusterNameFlag(cmd))
	}
	return nil
}

func validateNameFlagAndArg(cmd *Cmd, options *fargate.Options) error {
	if options.ProfileName != "" && cmd.NameArg != "" {
		return ErrFlagAndArg(fmt.Sprintf("--%s", fargateProfileName), options.ProfileName, cmd.NameArg)
	}
	if options.ProfileName == "" && cmd.NameArg != "" {
		options.ProfileName = cmd.NameArg
	}
	return nil
}

// NewDeleteFargateProfileLoader will load config or use flags for
// 'eksctl delete fargateprofile'
func NewDeleteFargateProfileLoader(cmd *Cmd, options *fargate.Options) ClusterConfigLoader {
	l := newCommonClusterConfigLoader(cmd)
	// We want to be able to pass the name of the profile to delete, even if we
	// use a ClusterConfig file to set metadata (cluster name, region, etc.):
	l.flagsIncompatibleWithConfigFile = flagsIncompatibleWithConfigFileExcept(fargateProfileName)
	l.flagsIncompatibleWithoutConfigFile.Insert(fargateProfileFlagsIncompatibleWithoutConfigFile...)
	l.validateWithoutConfigFile = func() error {
		if err := validate(cmd, options); err != nil {
			return err
		}
		return options.Validate()
	}
	l.validateWithConfigFile = func() error {
		if err := validate(cmd, options); err != nil {
			return err
		}
		return options.Validate()
	}
	return l
}
