package fargate

import (
	"errors"
	"fmt"
	"strings"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/utils/names"
)

// Options groups the parameters required to interact with Fargate.
type Options struct {
	ProfileName string
}

// Validate validates this Options object's fields.
func (o *Options) Validate() error {
	if o.ProfileName == "" {
		return errors.New("invalid Fargate profile: empty name")
	}
	return nil
}

// CreateOptions groups the parameters required to create a Fargate profile.
type CreateOptions struct {
	Options
	ProfileSelectorNamespace string
	// +optional
	ProfileSelectorLabels map[string]string
}

// Validate validates this Options object's fields.
func (o *CreateOptions) Validate() error {
	if strings.HasPrefix(o.ProfileName, api.ReservedProfileNamePrefix) {
		return fmt.Errorf("invalid Fargate profile: name should NOT start with %q", api.ReservedProfileNamePrefix)
	}
	if o.ProfileSelectorNamespace == "" {
		return errors.New("invalid Fargate profile: empty selector namespace")
	}
	return nil
}

// ToFargateProfile creates a FargateProfile object from this Options object.
func (o CreateOptions) ToFargateProfile() *api.FargateProfile {
	return &api.FargateProfile{
		Name: names.ForFargateProfile(o.ProfileName),
		Selectors: []api.FargateProfileSelector{
			api.FargateProfileSelector{
				Namespace: o.ProfileSelectorNamespace,
				Labels:    o.ProfileSelectorLabels,
			},
		},
	}
}
