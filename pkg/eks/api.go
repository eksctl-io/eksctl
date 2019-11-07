package eks

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/aws/aws-sdk-go/service/cloudtrail/cloudtrailiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/yaml"

	"github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/az"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/utils"
	"github.com/weaveworks/eksctl/pkg/version"
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
	spec  *api.ProviderConfig
	cfn   cloudformationiface.CloudFormationAPI
	eks   eksiface.EKSAPI
	ec2   ec2iface.EC2API
	elb   elbiface.ELBAPI
	elbv2 elbv2iface.ELBV2API
	sts   stsiface.STSAPI
	ssm   ssmiface.SSMAPI
	iam   iamiface.IAMAPI

	cloudtrail cloudtrailiface.CloudTrailAPI
}

// CloudFormation returns a representation of the CloudFormation API
func (p ProviderServices) CloudFormation() cloudformationiface.CloudFormationAPI { return p.cfn }

// CloudFormationRoleARN returns, if any,  a service role used by CloudFormation to call AWS API on your behalf
func (p ProviderServices) CloudFormationRoleARN() string { return p.spec.CloudFormationRoleARN }

// EKS returns a representation of the EKS API
func (p ProviderServices) EKS() eksiface.EKSAPI { return p.eks }

// EC2 returns a representation of the EC2 API
func (p ProviderServices) EC2() ec2iface.EC2API { return p.ec2 }

// ELB returns a representation of the ELB API
func (p ProviderServices) ELB() elbiface.ELBAPI { return p.elb }

// ELBV2 returns a representation of the ELBV2 API
func (p ProviderServices) ELBV2() elbv2iface.ELBV2API { return p.elbv2 }

// STS returns a representation of the STS API
func (p ProviderServices) STS() stsiface.STSAPI { return p.sts }

// SSM returns a representation of the STS API
func (p ProviderServices) SSM() ssmiface.SSMAPI { return p.ssm }

// IAM returns a representation of the IAM API
func (p ProviderServices) IAM() iamiface.IAMAPI { return p.iam }

// CloudTrail returns a representation of the CloudTrail API
func (p ProviderServices) CloudTrail() cloudtrailiface.CloudTrailAPI { return p.cloudtrail }

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
	clusterInfo  *clusterInfo
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
	s := c.newSession(spec)

	provider.cfn = cloudformation.New(s)
	provider.eks = awseks.New(s)
	provider.ec2 = ec2.New(s)
	provider.elb = elb.New(s)
	provider.elbv2 = elbv2.New(s)
	provider.sts = sts.New(s,
		// STS retrier has to be disabled, as it's not very helpful
		// (see https://github.com/weaveworks/eksctl/issues/705)
		request.WithRetryer(s.Config.Copy(),
			&client.DefaultRetryer{
				NumMaxRetries: 1,
			},
		),
	)
	provider.ssm = ssm.New(s)
	provider.iam = iam.New(s)
	provider.cloudtrail = cloudtrail.New(s)

	c.Status = &ProviderStatus{
		sessionCreds: s.Config.Credentials,
	}

	// override sessions if any custom endpoints specified
	if endpoint, ok := os.LookupEnv("AWS_CLOUDFORMATION_ENDPOINT"); ok {
		logger.Debug("Setting CloudFormation endpoint to %s", endpoint)
		provider.cfn = cloudformation.New(s, s.Config.Copy().WithEndpoint(endpoint))
	}
	if endpoint, ok := os.LookupEnv("AWS_EKS_ENDPOINT"); ok {
		logger.Debug("Setting EKS endpoint to %s", endpoint)
		provider.eks = awseks.New(s, s.Config.Copy().WithEndpoint(endpoint))
	}
	if endpoint, ok := os.LookupEnv("AWS_EC2_ENDPOINT"); ok {
		logger.Debug("Setting EC2 endpoint to %s", endpoint)
		provider.ec2 = ec2.New(s, s.Config.Copy().WithEndpoint(endpoint))

	}
	if endpoint, ok := os.LookupEnv("AWS_ELB_ENDPOINT"); ok {
		logger.Debug("Setting ELB endpoint to %s", endpoint)
		provider.elb = elb.New(s, s.Config.Copy().WithEndpoint(endpoint))

	}
	if endpoint, ok := os.LookupEnv("AWS_ELBV2_ENDPOINT"); ok {
		logger.Debug("Setting ELBV2 endpoint to %s", endpoint)
		provider.elbv2 = elbv2.New(s, s.Config.Copy().WithEndpoint(endpoint))

	}
	if endpoint, ok := os.LookupEnv("AWS_STS_ENDPOINT"); ok {
		logger.Debug("Setting STS endpoint to %s", endpoint)
		provider.sts = sts.New(s, s.Config.Copy().WithEndpoint(endpoint))
	}
	if endpoint, ok := os.LookupEnv("AWS_IAM_ENDPOINT"); ok {
		logger.Debug("Setting IAM endpoint to %s", endpoint)
		provider.iam = iam.New(s, s.Config.Copy().WithEndpoint(endpoint))
	}
	if endpoint, ok := os.LookupEnv("AWS_CLOUDTRAIL_ENDPOINT"); ok {
		logger.Debug("Setting CloudTrail endpoint to %s", endpoint)
		provider.cloudtrail = cloudtrail.New(s, s.Config.Copy().WithEndpoint(endpoint))
	}

	if clusterSpec != nil {
		clusterSpec.Metadata.Region = c.Provider.Region()
	}

	return c
}

// LoadConfigFromFile loads ClusterConfig from configFile
func LoadConfigFromFile(configFile string) (*api.ClusterConfig, error) {
	data, err := readConfig(configFile)
	if err != nil {
		return nil, errors.Wrapf(err, "reading config file %q", configFile)
	}

	// strict mode is not available in runtime.Decode, so we use the parser
	// directly; we don't store the resulting object, this is just the means
	// of detecting any unknown keys
	// NOTE: we must use sigs.k8s.io/yaml, as it behaves differently from
	// github.com/ghodss/yaml, which didn't handle nested structs well
	if err := yaml.UnmarshalStrict(data, &api.ClusterConfig{}); err != nil {
		return nil, errors.Wrapf(err, "loading config file %q", configFile)
	}

	obj, err := runtime.Decode(scheme.Codecs.UniversalDeserializer(), data)
	if err != nil {
		return nil, errors.Wrapf(err, "loading config file %q", configFile)
	}

	cfg, ok := obj.(*api.ClusterConfig)
	if !ok {
		return nil, fmt.Errorf("expected to decode object of type %T; got %T", &api.ClusterConfig{}, cfg)
	}
	return cfg, nil
}

func readConfig(configFile string) ([]byte, error) {
	if configFile == "-" {
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(configFile)
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

	input := &sts.GetCallerIdentityInput{}
	output, err := c.Provider.STS().GetCallerIdentity(input)
	if err != nil {
		return errors.Wrap(err, "checking AWS STS access – cannot get role ARN for current session")
	}
	if output == nil || output.Arn == nil {
		return fmt.Errorf("unexpected response from AWS STS")
	}
	c.Status.iamRoleARN = *output.Arn
	logger.Debug("role ARN for the current session is %q", c.Status.iamRoleARN)
	return nil
}

// EnsureAMI ensures that the node AMI is set and is available
func (c *ClusterProvider) EnsureAMI(version string, ng *api.NodeGroup) error {
	if api.IsAMI(ng.AMI) {
		return ami.Use(c.Provider.EC2(), ng)
	}

	var resolver ami.Resolver
	switch ng.AMI {
	case api.NodeImageResolverAuto:
		resolver = ami.NewAutoResolver(c.Provider.EC2())
	case api.NodeImageResolverAutoSSM:
		resolver = ami.NewSSMResolver(c.Provider.SSM())
	default:
		resolver = ami.NewDefaultResolver()
	}

	instanceType := selectInstanceType(ng)
	id, err := resolver.Resolve(c.Provider.Region(), version, instanceType, ng.AMIFamily)
	if err != nil {
		return errors.Wrap(err, "unable to determine AMI to use")
	}
	if id == "" {
		return ami.NewErrFailedResolution(c.Provider.Region(), version, instanceType, ng.AMIFamily)
	}
	ng.AMI = id

	// Check the AMI is available and populate RootDevice information
	return ami.Use(c.Provider.EC2(), ng)

}

// selectInstanceType determines which instanceType is relevant for selecting an AMI
// If the nodegroup has mixed instances it will prefer a GPU instance type over a general class one
// This is to make sure that the AMI that is selected later is valid for all the types
func selectInstanceType(ng *api.NodeGroup) string {
	if api.HasMixedInstances(ng) {
		for _, instanceType := range ng.InstancesDistribution.InstanceTypes {
			if utils.IsGPUInstanceType(instanceType) {
				return instanceType
			}
		}
		return ng.InstancesDistribution.InstanceTypes[0]
	}
	return ng.InstanceType
}

func errTooFewAvailabilityZones(azs []string) error {
	return fmt.Errorf("only %d zones specified %v, %d are required (can be non-unique)", len(azs), azs, az.MinRequiredAvailabilityZones)
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

func (c *ClusterProvider) newSession(spec *api.ProviderConfig) *session.Session {
	// we might want to use bits from kops, although right now it seems like too many thing we
	// don't want yet
	// https://github.com/kubernetes/kops/blob/master/upup/pkg/fi/cloudup/awsup/aws_cloud.go#L179
	config := aws.NewConfig()

	if c.Provider.Region() != "" {
		config = config.WithRegion(c.Provider.Region())
	}

	config = config.WithCredentialsChainVerboseErrors(true)
	config = request.WithRetryer(config, newLoggingRetryer())
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

	s := session.Must(session.NewSessionWithOptions(opts))

	s.Handlers.Build.PushFrontNamed(request.NamedHandler{
		Name: "eksctlUserAgent",
		Fn: request.MakeAddToUserAgentHandler(
			"eksctl", version.String()),
	})

	if spec.Region == "" {
		if api.IsSetAndNonEmptyString(s.Config.Region) {
			// set cluster config region, based on session config
			spec.Region = *s.Config.Region
		} else {
			// if session config doesn't have region set, make recursive call forcing default region
			logger.Debug("no region specified in flags or config, setting to %s", api.DefaultRegion)
			spec.Region = api.DefaultRegion
			return c.newSession(spec)
		}
	}

	return s
}

// NewStackManager returns a new stack manager
func (c *ClusterProvider) NewStackManager(spec *api.ClusterConfig) *manager.StackCollection {
	return manager.NewStackCollection(c.Provider, spec)
}
