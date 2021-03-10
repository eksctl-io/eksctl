package cmdutils

import (
	"fmt"

	"github.com/spf13/pflag"
)

// AddIAMServiceAccountFilterFlags add common `--include` and `--exclude` flags for filtering iamserviceaccounts
func AddIAMServiceAccountFilterFlags(fs *pflag.FlagSet, includeGlobs, excludeGlobs *[]string) {
	fs.StringSliceVar(includeGlobs, "include", nil,
		"iamserviceaccounts to include (list of globs), e.g.: 'default/s3-reader,*/dynamo-*'")

	fs.StringSliceVar(excludeGlobs, "exclude", nil,
		"iamserviceaccounts to exclude (list of globs), e.g.: 'default/s3-reader,*/dynamo-*'")
}

// AddIAMIdentityMappingARNFlags adds --arn and deprecated --role flags
func AddIAMIdentityMappingARNFlags(fs *pflag.FlagSet, cmd *Cmd, arn *string, verb string) {
	fs.StringVar(arn, "arn", "", fmt.Sprintf("ARN of the IAM role or user to %s", verb))
	fs.StringVar(arn, "role", "", "")
	_ = fs.MarkDeprecated("role", "use --arn")
}
