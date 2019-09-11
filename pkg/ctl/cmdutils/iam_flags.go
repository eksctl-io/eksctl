package cmdutils

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/iam"
)

// AddIAMServiceAccountFilterFlags add common `--include` and `--exclude` flags for filtering iamserviceaccounts
func AddIAMServiceAccountFilterFlags(fs *pflag.FlagSet, includeGlobs, excludeGlobs *[]string) {
	fs.StringSliceVar(includeGlobs, "include", nil,
		"iamserviceaccounts to include (list of globs), e.g.: 'default/s3-reader,*/dynamo-*'")

	fs.StringSliceVar(excludeGlobs, "exclude", nil,
		"iamserviceaccounts to exclude (list of globs), e.g.: 'default/s3-reader,*/dynamo-*'")
}

// AddIAMIdentityMappingARNFlags adds --arn and deprecated --role flags
func AddIAMIdentityMappingARNFlags(fs *pflag.FlagSet, cmd *Cmd, arn iam.ARN) {
	fs.Var(&arn, "arn", "ARN of the IAM role or user to create")
	// Add deprecated --role
	var role string
	fs.StringVar(&role, "role", "", "")
	_ = fs.MarkDeprecated("role", "see --arn")
	AddPreRunE(cmd.CobraCommand, func(cobraCmd *cobra.Command, args []string) error {
		var err error
		if role != "" {
			arn, err = iam.Parse(role)
		}
		return err
	})
}
