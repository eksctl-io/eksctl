package eks

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/az"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/version"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/kris-nova/logger"
)

// ClusterProvider stores information about the cluster
type ClusterProvider struct {
	// core fields used for config and AWS APIs
	Provider api.ClusterProvider
	// informative fields, i.e. used as outputs
	Status *ProviderStatus
}

// ProviderServices stores the used APIs
type ProviderServices struct {
	spec *api.ProviderConfig
	cfn  cloudformationiface.CloudFormationAPI
	eks  eksiface.EKSAPI
	ec2  ec2iface.EC2API
	sts  stsiface.STSAPI
	iam  iamiface.IAMAPI
}

// CloudFormation returns a representation of the CloudFormation API
func (p ProviderServices) CloudFormation() cloudformationiface.CloudFormationAPI { return p.cfn }

// CloudFormationRoleARN returns, if any,  a service role used by CloudFormation to call AWS API on your behalf
func (p ProviderServices) CloudFormationRoleARN() string { return p.spec.CloudFormationRoleARN }

// EKS returns a representation of the EKS API
func (p ProviderServices) EKS() eksiface.EKSAPI { return p.eks }

// EC2 returns a representation of the EC2 API
func (p ProviderServices) EC2() ec2iface.EC2API { return p.ec2 }

// STS returns a representation of the STS API
func (p ProviderServices) STS() stsiface.STSAPI { return p.sts }

// IAM returns a representation of the IAM API
func (p ProviderServices) IAM() iamiface.IAMAPI { return p.iam }

// Region returns provider-level region setting
func (p ProviderServices) Region() string { return p.spec.Region }

// Profile returns provider-level profile name
func (p ProviderServices) Profile() string { return p.spec.Profile }

// WaitTimeout returns provider-level duration after which any wait operation has to timeout
func (p ProviderServices) WaitTimeout() time.Duration { return p.spec.WaitTimeout }

// ProviderStatus stores information about the used IAM role and the resulting session
type ProviderStatus struct {
	iamRoleARN   string
	sessionCreds *credentials.Credentials
}

// New creates a new setup of the used AWS APIs
func New(spec *api.ProviderConfig, clusterSpec *api.ClusterConfig) *ClusterProvider {
	provider := &ProviderServices{
		spec: spec,
	}
	c := &ClusterProvider{
		Provider: provider,
	}
	// Create a new session and save credentials for possible
	// later re-use if overriding sessions due to custom URL
	s := c.newSession(spec, "", nil)

	provider.cfn = cloudformation.New(s)
	provider.eks = awseks.New(s)
	provider.ec2 = ec2.New(s)
	provider.sts = sts.New(s)
	provider.iam = iam.New(s)

	c.Status = &ProviderStatus{
		sessionCreds: s.Config.Credentials,
	}

	// override sessions if any custom endpoints specified
	if endpoint, ok := os.LookupEnv("AWS_CLOUDFORMATION_ENDPOINT"); ok {
		logger.Debug("Setting CloudFormation endpoint to %s", endpoint)
		provider.cfn = cloudformation.New(c.newSession(spec, endpoint, c.Status.sessionCreds))
	}
	if endpoint, ok := os.LookupEnv("AWS_EKS_ENDPOINT"); ok {
		logger.Debug("Setting EKS endpoint to %s", endpoint)
		provider.eks = awseks.New(c.newSession(spec, endpoint, c.Status.sessionCreds))
	}
	if endpoint, ok := os.LookupEnv("AWS_EC2_ENDPOINT"); ok {
		logger.Debug("Setting EC2 endpoint to %s", endpoint)
		provider.ec2 = ec2.New(c.newSession(spec, endpoint, c.Status.sessionCreds))

	}
	if endpoint, ok := os.LookupEnv("AWS_STS_ENDPOINT"); ok {
		logger.Debug("Setting STS endpoint to %s", endpoint)
		provider.sts = sts.New(c.newSession(spec, endpoint, c.Status.sessionCreds))
	}
	if endpoint, ok := os.LookupEnv("AWS_IAM_ENDPOINT"); ok {
		logger.Debug("Setting IAM endpoint to %s", endpoint)
		provider.iam = iam.New(c.newSession(spec, endpoint, c.Status.sessionCreds))
	}

	if clusterSpec != nil {
		clusterSpec.Metadata.Region = c.Provider.Region()
	}

	return c
}

// LoadConfigFromFile populates cfg based on contents of configFile
func LoadConfigFromFile(configFile string, cfg *api.ClusterConfig) error {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return errors.Wrapf(err, "reading config file %q", configFile)
	}

	obj, err := runtime.Decode(scheme.Codecs.UniversalDeserializer(), data)
	if err != nil {
		return errors.Wrapf(err, "loading config file %q", configFile)
	}

	cfgLoaded, ok := obj.(*api.ClusterConfig)
	if !ok {
		return fmt.Errorf("decoded object of wrong type")
	}

	*cfg = *cfgLoaded // mutate the content, not the reference

	return nil
}

// IsSupportedRegion check if given region is supported
func (c *ClusterProvider) IsSupportedRegion() bool {
	for _, supportedRegion := range api.SupportedRegions() {
		if c.Provider.Region() == supportedRegion {
			return true
		}
	}
	return false
}

// GetCredentialsEnv returns the AWS credentials for env usage
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

// CheckAuth checks the AWS authentication
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

// EnsureAMI ensures that the node AMI is set and is available
func (c *ClusterProvider) EnsureAMI(version string, ng *api.NodeGroup) error {
	if ng.AMI == ami.ResolverAuto {
		ami.DefaultResolvers = []ami.Resolver{ami.NewAutoResolver(c.Provider.EC2())}
	}
	if ng.AMI == ami.ResolverStatic || ng.AMI == ami.ResolverAuto {
		id, err := ami.Resolve(c.Provider.Region(), version, ng.InstanceType, ng.AMIFamily)
		if err != nil {
			return errors.Wrap(err, "Unable to determine AMI to use")
		}
		if id == "" {
			return ami.NewErrFailedResolution(c.Provider.Region(), version, ng.InstanceType, ng.AMIFamily)
		}
		ng.AMI = id
	}

	// Check the AMI is available
	available, err := ami.IsAvailable(c.Provider.EC2(), ng.AMI)
	if err != nil {
		return errors.Wrapf(err, "%s is not available", ng.AMI)
	}

	if !available {
		return ami.NewErrNotFound(ng.AMI)
	}

	return nil
}

// SetNodeLabels initialises and validate node labels based on cluster and nodegroup names
func (c *ClusterProvider) SetNodeLabels(ng *api.NodeGroup, meta *api.ClusterMeta) error {
	if ng.Labels == nil {
		ng.Labels = make(map[string]string)
	}

	ng.Labels[api.ClusterNameLabel] = meta.Name
	ng.Labels[api.NodeGroupNameLabel] = ng.Name

	return api.ValidateNodeGroupLabels(ng)
}

func errTooFewAvailabilityZones(azs []string) error {
	return fmt.Errorf("only %d zones specified %v, %d are required (can be non-unque)", len(azs), azs, az.MinRequiredAvailabilityZones)
}

// SetAvailabilityZones sets the given (or chooses) the availability zones
func (c *ClusterProvider) SetAvailabilityZones(spec *api.ClusterConfig, given []string) error {
	if count := len(given); count != 0 {
		if count < az.MinRequiredAvailabilityZones {
			return errTooFewAvailabilityZones(given)
		}
		spec.AvailabilityZones = given
		return nil
	}

	if count := len(spec.AvailabilityZones); count != 0 {
		if count < az.MinRequiredAvailabilityZones {
			return errTooFewAvailabilityZones(spec.AvailabilityZones)
		}
		return nil
	}

	logger.Debug("determining availability zones")
	azSelector := az.NewSelectorWithDefaults(c.Provider.EC2())
	if c.Provider.Region() == api.RegionUSEast1 {
		azSelector = az.NewSelectorWithMinRequired(c.Provider.EC2())
	}
	zones, err := azSelector.SelectZones(c.Provider.Region())
	if err != nil {
		return errors.Wrap(err, "getting availability zones")
	}

	logger.Info("setting availability zones to %v", zones)
	spec.AvailabilityZones = zones

	return nil
}

func (c *ClusterProvider) newSession(spec *api.ProviderConfig, endpoint string, credentials *credentials.Credentials) *session.Session {
	// we might want to use bits from kops, although right now it seems like too many thing we
	// don't want yet
	// https://github.com/kubernetes/kops/blob/master/upup/pkg/fi/cloudup/awsup/aws_cloud.go#L179
	config := aws.NewConfig()

	if c.Provider.Region() != "" {
		config = config.WithRegion(c.Provider.Region())
	}

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
		Profile:                 spec.Profile,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	}

	stscreds.DefaultDuration = 30 * time.Minute

	if len(endpoint) > 0 {
		opts.Config.Endpoint = &endpoint
	}

	if credentials != nil {
		opts.Config.Credentials = credentials
	}

	s := session.Must(session.NewSessionWithOptions(opts))

	s.Handlers.Build.PushFrontNamed(request.NamedHandler{
		Name: "eksctlUserAgent",
		Fn: request.MakeAddToUserAgentHandler(
			"eksctl", version.String()),
	})

	if spec.Region == "" {
		if *s.Config.Region != "" {
			// set cluster config region, based on session config
			spec.Region = *s.Config.Region
		} else {
			// if session config doesn't have region set, make recursive call forcing default region
			logger.Debug("no region specified in flags or config, setting to %s", api.DefaultRegion)
			spec.Region = api.DefaultRegion
			return c.newSession(spec, endpoint, credentials)
		}
	}

	return s
}

// NewStackManager returns a new stack manager
func (c *ClusterProvider) NewStackManager(spec *api.ClusterConfig) *manager.StackCollection {
	return manager.NewStackCollection(c.Provider, spec)
}
