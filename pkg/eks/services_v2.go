package eks

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/weaveworks/eksctl/pkg/awsapi"
)

// ServicesV2 implements api.ServicesV2.
// The SDK clients are initialized lazily and guarded by a mutex.
type ServicesV2 struct {
	config aws.Config

	// mu guards initialization of SDK clients.
	// All service methods should ensure that their initialization is guarded by mu.
	mu             sync.Mutex
	sts            *sts.Client
	cloudformation *cloudformation.Client
}

// STSV2 implements the AWS STS service.
func (s *ServicesV2) STSV2() awsapi.STS {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.sts == nil {
		s.sts = sts.NewFromConfig(s.config, func(o *sts.Options) {
			// Disable retryer for STS
			// (see https://github.com/weaveworks/eksctl/issues/705)
			o.Retryer = aws.NopRetryer{}
		})
	}
	return s.sts
}

// CloudFormationV2 implements the AWS CloudFormation service.
func (s *ServicesV2) CloudFormationV2() awsapi.CloudFormation {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.cloudformation == nil {
		s.cloudformation = cloudformation.NewFromConfig(s.config, func(o *cloudformation.Options) {
			// Use adaptive mode for retrying CloudFormation requests to mimic
			// the logic used for AWS SDK v1.
			o.Retryer = retry.NewAdaptiveMode(func(o *retry.AdaptiveModeOptions) {
				o.StandardOptions = []func(*retry.StandardOptions){
					func(so *retry.StandardOptions) {
						so.MaxAttempts = maxRetries
					},
				}
			})
		})
	}
	return s.cloudformation
}
