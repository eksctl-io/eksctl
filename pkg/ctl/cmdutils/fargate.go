package cmdutils

import (
	"github.com/spf13/pflag"
	"github.com/weaveworks/eksctl/pkg/fargate"
)

const (
	name = "name" // Fargate profile name.
)

// AddCommonFlagsForFargate configures the flags required to interact with
// Fargate.
func AddCommonFlagsForFargate(fs *pflag.FlagSet, opts *fargate.Options) {
	fs.StringVar(&opts.ProfileName, name, "",
		"Fargate profile name")
}
