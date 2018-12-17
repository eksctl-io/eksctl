package cmdutils

import (
	"fmt"
	"os"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

// IncompatibleFlags is a common substring of an error message
const IncompatibleFlags = "cannot be used at the same time"

// GetNameArg tests to ensure there is only 1 name argument
func GetNameArg(args []string) string {
	if len(args) > 1 {
		logger.Critical("only one argument is allowed to be used as a name")
		os.Exit(1)
	}
	if len(args) == 1 {
		return (strings.TrimSpace(args[0]))
	}
	return ""
}

// AddCommonFlagsForAWS adds common flags for api.ProviderConfig
func AddCommonFlagsForAWS(group *NamedFlagSetGroup, p *api.ProviderConfig) {
	group.InFlagSet("AWS client", func(fs *pflag.FlagSet) {
		fs.StringVarP(&p.Profile, "profile", "p", "", "AWS credentials profile to use (overrides the AWS_PROFILE environment variable)")

		fs.DurationVar(&p.WaitTimeout, "aws-api-timeout", api.DefaultWaitTimeout, "")
		// TODO deprecate in 0.2.0
		if err := fs.MarkHidden("aws-api-timeout"); err != nil {
			logger.Debug("ignoring error %q", err.Error())
		}
		fs.DurationVar(&p.WaitTimeout, "timeout", api.DefaultWaitTimeout, "max wait time in any polling operations")
	})
}

// AddRegionFlag adds common --region flag
func AddRegionFlag(fs *pflag.FlagSet, p *api.ProviderConfig) {
	fs.StringVarP(&p.Region, "region", "r", "", "AWS region")
}

// AddCommonFlagsForKubeconfig adds common flags for controlling how output kubeconfig is written
func AddCommonFlagsForKubeconfig(fs *pflag.FlagSet, outputPath *string, setContext, autoPath *bool, exampleName string) {
	fs.StringVar(outputPath, "kubeconfig", kubeconfig.DefaultPath, "path to write kubeconfig (incompatible with --auto-kubeconfig)")
	fs.BoolVar(setContext, "set-kubeconfig-context", true, "if true then current-context will be set in kubeconfig; if a context is already set then it will be overwritten")
	fs.BoolVar(autoPath, "auto-kubeconfig", false, fmt.Sprintf("save kubeconfig file by cluster name, e.g. %q", kubeconfig.AutoPath(exampleName)))
}

// ErrUnsupportedRegion is a common error message
func ErrUnsupportedRegion(p *api.ProviderConfig) error {
	return fmt.Errorf("--region=%s is not supported - use one of: %s", p.Region, strings.Join(api.SupportedRegions(), ", "))
}

// ErrNameFlagAndArg is a common error message
func ErrNameFlagAndArg(nameFlag, nameArg string) error {
	return fmt.Errorf("--name=%s and argument %s %s", nameFlag, nameArg, IncompatibleFlags)
}
