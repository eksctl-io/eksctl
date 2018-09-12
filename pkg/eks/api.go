package eks

import (
	"fmt"
	"os"
	"time"

	"github.com/weaveworks/eksctl/pkg/ami"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks/api"

	"github.com/weaveworks/eksctl/pkg/az"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	"github.com/kubicorn/kubicorn/pkg/logger"
)

type ClusterProvider struct {
	// core fields used for config and AWS APIs
	Spec     *api.ClusterConfig
	Provider api.ClusterProvider
	// informative fields, i.e. used as outputs
	Status *ProviderStatus
}

type ProviderServices struct {
	cfn cloudformationiface.CloudFormationAPI
	eks eksiface.EKSAPI
	ec2 ec2iface.EC2API
	sts stsiface.STSAPI
}

func (p ProviderServices) CloudFormation() cloudformationiface.CloudFormationAPI { return p.cfn }

func (p ProviderServices) EKS() eksiface.EKSAPI { return p.eks }
func (p ProviderServices) EC2() ec2iface.EC2API { return p.ec2 }
func (p ProviderServices) STS() stsiface.STSAPI { return p.sts }

type ProviderStatus struct {
	iamRoleARN        string
	sessionCreds      *credentials.Credentials
	availabilityZones []string
}

func New(clusterConfig *api.ClusterConfig) *ClusterProvider {
	// Create a new session and save credentials for possible
	// later re-use if overriding sessions due to custom URL
	s := newSession(clusterConfig, "", nil)

	provider := &ProviderServices{
		cfn: cloudformation.New(s),
		eks: awseks.New(s),
		ec2: ec2.New(s),
		sts: sts.New(s),
	}

	status := &ProviderStatus{
		sessionCreds: s.Config.Credentials,
	}

	// override sessions if any custom endpoints specified
	if endpoint, ok := os.LookupEnv("AWS_CLOUDFORMATION_ENDPOINT"); ok {
		logger.Debug("Setting CloudFormation endpoint to %s", endpoint)
		provider.cfn = cloudformation.New(newSession(clusterConfig, endpoint, status.sessionCreds))
	}
	if endpoint, ok := os.LookupEnv("AWS_EKS_ENDPOINT"); ok {
		logger.Debug("Setting EKS endpoint to %s", endpoint)
		provider.eks = awseks.New(newSession(clusterConfig, endpoint, status.sessionCreds))
	}
	if endpoint, ok := os.LookupEnv("AWS_EC2_ENDPOINT"); ok {
		logger.Debug("Setting EC2 endpoint to %s", endpoint)
		provider.ec2 = ec2.New(newSession(clusterConfig, endpoint, status.sessionCreds))

	}
	if endpoint, ok := os.LookupEnv("AWS_STS_ENDPOINT"); ok {
		logger.Debug("Setting STS endpoint to %s", endpoint)
		provider.sts = sts.New(newSession(clusterConfig, endpoint, status.sessionCreds))
	}

	return &ClusterProvider{
		Spec:     clusterConfig,
		Provider: provider,
		Status:   status,
	}
}

func (c *ClusterProvider) GetCredentialsEnv() ([]string, error) {
	creds, err := c.Status.sessionCreds.Get()
	if err != nil {
		return nil, errors.Wrap(err, "getting effective credentials")
	}
	return []string{
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", creds.AccessKeyID),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", creds.SecretAccessKey),
		fmt.Sprintf("AWS_SESSION_TOKEN=%s", creds.SessionToken),
	}, nil
}

func (c *ClusterProvider) CheckAuth() error {
	{
		input := &sts.GetCallerIdentityInput{}
		output, err := c.Provider.STS().GetCallerIdentity(input)
		if err != nil {
			return errors.Wrap(err, "checking AWS STS access – cannot get role ARN for current session")
		}
		c.Status.iamRoleARN = *output.Arn
		logger.Debug("role ARN for the current session is %q", c.Status.iamRoleARN)
	}
	{
		input := &cloudformation.ListStacksInput{}
		if _, err := c.Provider.CloudFormation().ListStacks(input); err != nil {
			return errors.Wrap(err, "checking AWS CloudFormation access – cannot list stacks")
		}
	}
	return nil
}

// EnsureAMI ensures that the node AMI is set and isavailable
func (c *ClusterProvider) EnsureAMI() error {
	// TODO: https://github.com/weaveworks/eksctl/issues/28
	// - imporve validation of parameter set overall, probably in another package
	if c.Spec.NodeAMI == ami.ResolverAuto {
		ami.DefaultResolvers = []ami.Resolver{ami.NewAutoResolver(c.Provider.EC2())}
	}
	if c.Spec.NodeAMI == ami.ResolverStatic || c.Spec.NodeAMI == ami.ResolverAuto {
		id, err := ami.Resolve(c.Spec.Region, c.Spec.NodeType)
		if err != nil {
			return errors.Wrap(err, "Unable to determine AMI to use")
		}
		if id == "" {
			return ami.NewErrFailedResolution(c.Spec.Region, c.Spec.NodeType)
		}
		c.Spec.NodeAMI = id
	}

	// Check the AMI is available
	available, err := ami.IsAvailable(c.Provider.EC2(), c.Spec.NodeAMI)
	if err != nil {
		return errors.Wrapf(err, "%s is not available", c.Spec.NodeAMI)
	}

	if !available {
		return ami.NewErrNotFound(c.Spec.NodeAMI)
	}

	logger.Info("using %q for nodes", c.Spec.NodeAMI)

	return nil
}

func (c *ClusterProvider) SetAvailabilityZones(given []string) error {
	if len(given) == 0 {
		logger.Debug("determining availability zones")
		azSelector := az.NewSelectorWithDefaults(c.Provider.EC2())
		zones, err := azSelector.SelectZones(c.Spec.Region)
		if err != nil {
			return errors.Wrap(err, "getting availability zones")
		}

		logger.Info("setting availability zones to %v", zones)
		c.Spec.AvailabilityZones = zones
		return nil
	}
	if len(given) < az.DefaultRequiredAvailabilityZones {
		return fmt.Errorf("only %d zones specified %v, %d are required (can be non-unque)", len(given), given, az.DefaultRequiredAvailabilityZones)
	}
	c.Spec.AvailabilityZones = given
	return nil
}

func newSession(clusterConfig *api.ClusterConfig, endpoint string, credentials *credentials.Credentials) *session.Session {
	// we might want to use bits from kops, although right now it seems like too many thing we
	// don't want yet
	// https://github.com/kubernetes/kops/blob/master/upup/pkg/fi/cloudup/awsup/aws_cloud.go#L179
	config := aws.NewConfig()
	config = config.WithRegion(clusterConfig.Region)
	config = config.WithCredentialsChainVerboseErrors(true)
	if logger.Level >= api.AWSDebugLevel {
		config = config.WithLogLevel(aws.LogDebug |
			aws.LogDebugWithHTTPBody |
			aws.LogDebugWithRequestRetries |
			aws.LogDebugWithRequestErrors |
			aws.LogDebugWithEventStreamBody)
		config = config.WithLogger(aws.LoggerFunc(func(args ...interface{}) {
			logger.Debug(fmt.Sprintln(args...))
		}))
	}

	// Create the options for the session
	opts := session.Options{
		Config:                  *config,
		SharedConfigState:       session.SharedConfigEnable,
		Profile:                 clusterConfig.Profile,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	}

	stscreds.DefaultDuration = 30 * time.Minute

	if len(endpoint) > 0 {
		opts.Config.Endpoint = &endpoint
	}

	if credentials != nil {
		opts.Config.Credentials = credentials
	}

	return session.Must(session.NewSessionWithOptions(opts))
}

func (c *ClusterProvider) NewStackManager() *manager.StackCollection {
	return manager.NewStackCollection(c.Provider, c.Spec)
}
