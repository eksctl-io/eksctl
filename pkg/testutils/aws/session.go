package aws

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// NewSession creates a new session
func NewSession(region string) *session.Session {
	config := aws.NewConfig()
	config = config.WithRegion(region)
	opts := session.Options{
		Config: *config,
	}
	return session.Must(session.NewSessionWithOptions(opts))
}
