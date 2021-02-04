package fargate

import (
	"errors"
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
