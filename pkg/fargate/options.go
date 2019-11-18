package fargate

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
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
