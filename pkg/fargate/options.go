package fargate

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"time"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
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
		return fmt.Errorf("invalid Fargate profile: name should NOT start with \"%s\"", api.ReservedProfileNamePrefix)
	}
	if o.ProfileSelectorNamespace == "" {
		return errors.New("invalid Fargate profile: empty selector namespace")
	}
	return nil
}

// ToFargateProfile creates a FargateProfile object from this Options object.
func (o CreateOptions) ToFargateProfile() *api.FargateProfile {
	return &api.FargateProfile{
		Name: GetOrDefaultProfileName(o.ProfileName),
		Selectors: []api.FargateProfileSelector{
			api.FargateProfileSelector{
				Namespace: o.ProfileSelectorNamespace,
				Labels:    o.ProfileSelectorLabels,
			},
		},
	}
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// GetOrDefaultProfileName returns the provided name if non-empty, or else
// generates a random name matching: fp-[abcdef0123456789]{8}
func GetOrDefaultProfileName(name string) string {
	if name != "" {
		return name
	}
	length := 8                 // Length of automatically generated Fargate profile names.
	chars := "abcdef0123456789" // Characters used to generate Fargate profile names.
	randomName := make([]byte, length)
	for i := 0; i < length; i++ {
		randomName[i] = chars[r.Intn(len(chars))]
	}
	return fmt.Sprintf("fp-%s", string(randomName))
}
