package eks

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	middlewarev2 "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go/middleware"
	"github.com/gofrs/flock"
	"github.com/kris-nova/logger"
	"github.com/spf13/afero"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/credentials"
	"github.com/weaveworks/eksctl/pkg/version"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate -o fakes/fake_configuration_loader.go . AWSConfigurationLoader
type AWSConfigurationLoader interface {
	LoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) (cfg aws.Config, err error)
}

type ConfigurationLoader struct {
	AWSConfigurationLoader
}

func (cl ConfigurationLoader) LoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx, optFns...)
}

func newV2Config(pc *api.ProviderConfig, credentialsCacheFilePath string, configurationLoader AWSConfigurationLoader) (aws.Config, error) {
	var options []func(options *config.LoadOptions) error

	if pc.Region != "" {
		options = append(options, config.WithRegion(pc.Region))
	}
	clientLogMode := aws.ClientLogMode(1)

	if logger.Level >= api.AWSDebugLevel {
		clientLogMode = clientLogMode | aws.LogRequestWithBody | aws.LogRequestEventMessage | aws.LogResponseWithBody | aws.LogRetries
	}
	options = append(options, config.WithClientLogMode(clientLogMode))

	if endpointResolver := makeEndpointResolverFunc(); endpointResolver != nil {
		options = append(options, config.WithEndpointResolverWithOptions(endpointResolver))
	}

	if !pc.Profile.SourceIsEnvVar {
		options = append(options, config.WithSharedConfigProfile(pc.Profile.Name))
	}

	cfg, err := configurationLoader.LoadDefaultConfig(context.TODO(), append(options,
		config.WithRetryer(func() aws.Retryer {
			return NewRetryerV2()
		}),
		config.WithAssumeRoleCredentialOptions(func(o *stscreds.AssumeRoleOptions) {
			o.TokenProvider = stscreds.StdinTokenProvider
			o.Duration = 30 * time.Minute
		}),
		config.WithAPIOptions([]func(stack *middleware.Stack) error{
			middlewarev2.AddUserAgentKeyValue("eksctl", version.String()),
		}),
		// Some CloudFormation operations can take a long time to complete, and we
		// don't want any temporary credentials to expire before this occurs. So if
		// it's less than 30 minutes before the current credentials will expire, try
		// to renew them first.
		config.WithCredentialsCacheOptions(func(o *aws.CredentialsCacheOptions) {
			logger.Debug("Setting credentials expiry window to 30 minutes")
			o.ExpiryWindow = 30 * time.Minute
			o.ExpiryWindowJitterFrac = 0
		}),
	)...)

	if err != nil {
		return cfg, err
	}
	if credentialsCacheFilePath != "" {
		fileCache, err := credentials.NewFileCacheV2(cfg.Credentials, pc.Profile.Name, afero.NewOsFs(), func(path string) credentials.Flock {
			return flock.New(path)
		}, &credentials.RealClock{}, credentialsCacheFilePath)
		if err != nil {
			return cfg, fmt.Errorf("error creating credentials cache: %w", err)
		}
		cfg.Credentials = aws.NewCredentialsCache(fileCache)
	}
	return cfg, nil
}

func makeEndpointResolverFunc() aws.EndpointResolverWithOptionsFunc {
	serviceIDEnvMap := map[string]string{
		cloudformation.ServiceID:         "AWS_CLOUDFORMATION_ENDPOINT",
		eks.ServiceID:                    "AWS_EKS_ENDPOINT",
		ec2.ServiceID:                    "AWS_EC2_ENDPOINT",
		elasticloadbalancing.ServiceID:   "AWS_ELB_ENDPOINT",
		elasticloadbalancingv2.ServiceID: "AWS_ELBV2_ENDPOINT",
		sts.ServiceID:                    "AWS_STS_ENDPOINT",
		iam.ServiceID:                    "AWS_IAM_ENDPOINT",
		cloudtrail.ServiceID:             "AWS_CLOUDTRAIL_ENDPOINT",
	}

	hasCustomEndpoint := false
	for service, envName := range serviceIDEnvMap {
		if endpoint, ok := os.LookupEnv(envName); ok {
			logger.Debug(
				"Setting %s endpoint to %s", service, endpoint)
			hasCustomEndpoint = true
		}
	}

	if !hasCustomEndpoint {
		return nil
	}

	return func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if envName, ok := serviceIDEnvMap[service]; ok {
			if ok {
				if endpoint, ok := os.LookupEnv(envName); ok {
					return aws.Endpoint{
						URL:           endpoint,
						SigningRegion: region,
					}, nil
				}
			}
		}
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	}
}
