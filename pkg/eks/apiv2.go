package eks

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	middlewarev2 "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
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

	if !pc.Profile.SourceIsEnvVar {
		options = append(options, config.WithSharedConfigProfile(pc.Profile.Name))
	}

	cfg, err := configurationLoader.LoadDefaultConfig(context.TODO(), append(options,
		config.WithRetryer(func() aws.Retryer {
			return NewRetryerV2()
		}),
		config.WithAssumeRoleCredentialOptions(func(o *stscreds.AssumeRoleOptions) {
			o.TokenProvider = stscreds.StdinTokenProvider
			o.Duration = 60 * time.Minute
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
