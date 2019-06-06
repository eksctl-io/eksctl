package cmdutils

import (
	"github.com/spf13/cobra"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type CommonParams struct {
	Command *cobra.Command

	Plan, Wait bool

	NameArg string

	ClusterConfigFile string

	ProviderConfig *api.ProviderConfig
	ClusterConfig  *api.ClusterConfig
}

func NewCommonParams(spec *api.ClusterConfig) *CommonParams {
	return &CommonParams{
		ProviderConfig: &api.ProviderConfig{},
		ClusterConfig:  spec,
	}
}
