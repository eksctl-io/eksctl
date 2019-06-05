package cmdutils

import (
	"fmt"
	"os"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

// IncompatibleFlags is a common substring of an error message
const IncompatibleFlags = "cannot be used at the same time"

// AddPreRun chains cmd.PreRun handlers, as cobra only allows one, so we don't
// accidentially override one we registered earlier
func AddPreRun(cmd *cobra.Command, newFn func(cmd *cobra.Command, args []string)) {
	currentFn := cmd.PreRun
	cmd.PreRun = func(cmd *cobra.Command, args []string) {
		if currentFn != nil {
			currentFn(cmd, args)
		}
		newFn(cmd, args)
	}
}

// LogIntendedAction calls logger.Info with appropriate prefix
func LogIntendedAction(plan bool, msgFmt string, args ...interface{}) {
	prefix := "will "
	if plan {
		prefix = "(plan) would "
	}
	logger.Info(prefix+msgFmt, args...)
}

// LogCompletedAction calls logger.Success with appropriate prefix
func LogCompletedAction(plan bool, msgFmt string, args ...interface{}) {
	prefix := ""
	if plan {
		prefix = "(plan) would have "
	}
	logger.Success(prefix+msgFmt, args...)
}

// LogPlanModeWarning will log a message to inform user that they are in plan-mode
func LogPlanModeWarning(plan bool) {
	if plan {
		logger.Warning("no changes were applied, run again with '--approve' to apply the changes")
	}
}

// AddApproveFlag adds common `--approve` flag
func AddApproveFlag(plan *bool, cmd *cobra.Command, fs *pflag.FlagSet) {
	approve := fs.Bool("approve", !*plan, "Apply the changes")
	AddPreRun(cmd, func(cmd *cobra.Command, args []string) {
		if cmd.Flag("approve").Changed {
			*plan = !*approve
		}
	})
}

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
func AddCommonFlagsForAWS(group *NamedFlagSetGroup, p *api.ProviderConfig, cfnRole bool) {
	group.InFlagSet("AWS client", func(fs *pflag.FlagSet) {
		fs.StringVarP(&p.Profile, "profile", "p", "", "AWS credentials profile to use (overrides the AWS_PROFILE environment variable)")

		fs.DurationVar(&p.WaitTimeout, "aws-api-timeout", api.DefaultWaitTimeout, "")
		// TODO deprecate in 0.2.0
		if err := fs.MarkHidden("aws-api-timeout"); err != nil {
			logger.Debug("ignoring error %q", err.Error())
		}
		fs.DurationVar(&p.WaitTimeout, "timeout", api.DefaultWaitTimeout, "max wait time in any polling operations")
		if cfnRole {
			fs.StringVar(&p.CloudFormationRoleARN, "cfn-role-arn", "", "IAM role used by CloudFormation to call AWS API on your behalf")
		}
	})
}

// AddRegionFlag adds common --region flag
func AddRegionFlag(fs *pflag.FlagSet, p *api.ProviderConfig) {
	fs.StringVarP(&p.Region, "region", "r", "", "AWS region")
}

// AddVersionFlag adds common --version flag
func AddVersionFlag(fs *pflag.FlagSet, meta *api.ClusterMeta, extraUsageInfo string) {
	usage := fmt.Sprintf("Kubernetes version (valid options: %s)", strings.Join(api.SupportedVersions(), ", "))
	if extraUsageInfo != "" {
		usage = fmt.Sprintf("%s [%s]", usage, extraUsageInfo)
	}
	fs.StringVar(&meta.Version, "version", meta.Version, usage)
}

// AddWaitFlag adds common --wait flag
func AddWaitFlag(wait *bool, fs *pflag.FlagSet, description string) {
	fs.BoolVarP(wait, "wait", "w", *wait, fmt.Sprintf("wait for %s before exiting", description))
}

// AddUpdateAuthConfigMap adds common --update-auth-configmap flag
func AddUpdateAuthConfigMap(updateAuthConfigMap *bool, fs *pflag.FlagSet, description string) {
	fs.BoolVar(updateAuthConfigMap, "update-auth-configmap", true, description)
}

// AddCommonFlagsForKubeconfig adds common flags for controlling how output kubeconfig is written
func AddCommonFlagsForKubeconfig(fs *pflag.FlagSet, outputPath *string, setContext, autoPath *bool, exampleName string) {
	fs.StringVar(outputPath, "kubeconfig", kubeconfig.DefaultPath, "path to write kubeconfig (incompatible with --auto-kubeconfig)")
	fs.BoolVar(setContext, "set-kubeconfig-context", true, "if true then current-context will be set in kubeconfig; if a context is already set then it will be overwritten")
	fs.BoolVar(autoPath, "auto-kubeconfig", false, fmt.Sprintf("save kubeconfig file by cluster name, e.g. %q", kubeconfig.AutoPath(exampleName)))
}

// AddCommonFlagsForGetCmd adds common flafs for get commands
func AddCommonFlagsForGetCmd(fs *pflag.FlagSet, chunkSize *int, outputMode *string) {
	fs.IntVar(chunkSize, "chunk-size", 100, "return large lists in chunks rather than all at once, pass 0 to disable")
	fs.StringVarP(outputMode, "output", "o", "table", "specifies the output format (valid option: table, json, yaml)")
}

// ErrUnsupportedRegion is a common error message
func ErrUnsupportedRegion(p *api.ProviderConfig) error {
	return fmt.Errorf("--region=%s is not supported - use one of: %s", p.Region, strings.Join(api.SupportedRegions(), ", "))
}

// ErrNameFlagAndArg is a common error message
func ErrNameFlagAndArg(nameFlag, nameArg string) error {
	return ErrFlagAndArg("--name", nameFlag, nameArg)
}

// ErrFlagAndArg may be used to err for options that can be given
// as flags /and/ arg but only one is allowed to be used.
func ErrFlagAndArg(kind, flag, arg string) error {
	return fmt.Errorf("%s=%s and argument %s %s", kind, flag, arg, IncompatibleFlags)
}

// ErrMustBeSet is a common error message
func ErrMustBeSet(pathOrFlag string) error {
	return fmt.Errorf("%s must be set", pathOrFlag)
}

// ErrCannotUseWithConfigFile is a common error message
func ErrCannotUseWithConfigFile(what string) error {
	return fmt.Errorf("cannot use %s when --config-file/-f is set", what)
}
