package eks

import (
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/ec2"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
)

// ServicesV2 implements api.ServicesV2.
// The SDK clients are initialized lazily and guarded by a mutex.
type ServicesV2 struct {
	config aws.Config

	// mu guards initialization of SDK clients.
	// All service methods should ensure that their initialization is guarded by mu.
	mu                     sync.Mutex
	sts                    *sts.Client
	stsPresigned           *sts.PresignClient
	cloudformation         *cloudformation.Client
	elasticloadbalancing   *elasticloadbalancing.Client
	elasticloadbalancingV2 *elasticloadbalancingv2.Client
	ssm                    *ssm.Client
	iam                    *iam.Client
	ec2                    *ec2.Client
	eks                    *eks.Client
}

// STS implements the AWS STS service.
func (s *ServicesV2) STS() awsapi.STS {
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

// STSPresigner provides a signed STS client for calls to Kubernetes.
func (s *ServicesV2) STSPresigner() api.STSPresigner {
	// set up sts client.
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.stsPresigned == nil {
		client := sts.NewFromConfig(s.config, func(o *sts.Options) {
			// Disable retryer for STS
			// (see https://github.com/weaveworks/eksctl/issues/705)
			o.Retryer = aws.NopRetryer{}
		})
		s.stsPresigned = sts.NewPresignClient(client)
	}
	return s.stsPresigned
}

// CloudFormation implements the AWS CloudFormation service.
func (s *ServicesV2) CloudFormation() awsapi.CloudFormation {
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

// ELB implements the AWS ELB service.
func (s *ServicesV2) ELB() awsapi.ELB {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.elasticloadbalancing == nil {
		s.elasticloadbalancing = elasticloadbalancing.NewFromConfig(s.config)
	}
	return s.elasticloadbalancing
}

// ELBV2 implements the ELBV2 service.
func (s *ServicesV2) ELBV2() awsapi.ELBV2 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.elasticloadbalancingV2 == nil {
		s.elasticloadbalancingV2 = elasticloadbalancingv2.NewFromConfig(s.config)
	}
	return s.elasticloadbalancingV2
}

// SSM implements the AWS SSM service.
func (s *ServicesV2) SSM() awsapi.SSM {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ssm == nil {
		s.ssm = ssm.NewFromConfig(s.config)
	}
	return s.ssm
}

// IAM implements the AWS IAM service.
func (s *ServicesV2) IAM() awsapi.IAM {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.iam == nil {
		s.iam = iam.NewFromConfig(s.config)
	}
	return s.iam
}

// EC2 implements the AWS EC2 service.
func (s *ServicesV2) EC2() awsapi.EC2 {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.ec2 == nil {
		s.ec2 = ec2.NewFromConfig(s.config)
	}
	return s.ec2
}

func (s *ServicesV2) EKS() awsapi.EKS {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.eks == nil {
		s.eks = eks.NewFromConfig(s.config)
	}
	return s.eks
}
