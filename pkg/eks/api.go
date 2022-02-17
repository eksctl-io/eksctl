package eks

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"

	"github.com/aws/aws-sdk-go/service/cloudwatchlogs/cloudwatchlogsiface"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
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
	ekscreds "github.com/weaveworks/eksctl/pkg/credentials"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/version"
)

// ClusterProvider stores information about the cluster
type ClusterProvider struct {
	// core fields used for config and AWS APIs
	Provider api.ClusterProvider
	// informative fields, i.e. used as outputs
	Status *ProviderStatus
}

//counterfeiter:generate -o fakes/fake_kube_provider.go . KubeProvider
// KubeProvider is an interface with helper funcs for k8s and EKS that are part of ClusterProvider
type KubeProvider interface {
	NewRawClient(spec *api.ClusterConfig) (*kubewrapper.RawClient, error)
	ServerVersion(rawClient *kubernetes.RawClient) (string, error)
	LoadClusterIntoSpecFromStack(spec *api.ClusterConfig, stackManager manager.StackManager) error
	SupportsManagedNodes(clusterConfig *api.ClusterConfig) (bool, error)
	ValidateClusterForCompatibility(cfg *api.ClusterConfig, stackManager manager.StackManager) error
	UpdateAuthConfigMap(nodeGroups []*api.NodeGroup, clientSet kubernetes.Interface) error
	WaitForNodes(clientSet kubernetes.Interface, ng KubeNodeGroup) error
}

// ProviderServices stores the used APIs
type ProviderServices struct {
	spec  *api.ProviderConfig
	cfn   cloudformationiface.CloudFormationAPI
	asg   autoscalingiface.AutoScalingAPI
	eks   eksiface.EKSAPI
	ec2   ec2iface.EC2API
	elb   elbiface.ELBAPI
	elbv2 elbv2iface.ELBV2API
	sts   stsiface.STSAPI
	ssm   ssmiface.SSMAPI
	iam   iamiface.IAMAPI

	cloudtrail     cloudtrailiface.CloudTrailAPI
	cloudwatchlogs cloudwatchlogsiface.CloudWatchLogsAPI

	session *session.Session
}

// CloudFormation returns a representation of the CloudFormation API
func (p ProviderServices) CloudFormation() cloudformationiface.CloudFormationAPI { return p.cfn }

// CloudFormationRoleARN returns, if any, a service role used by CloudFormation to call AWS API on your behalf
func (p ProviderServices) CloudFormationRoleARN() string { return p.spec.CloudFormationRoleARN }

// CloudFormationDisableRollback returns whether stacks should not rollback on failure
func (p ProviderServices) CloudFormationDisableRollback() bool {
	return p.spec.CloudFormationDisableRollback
}

// ASG returns a representation of the AutoScaling API
func (p ProviderServices) ASG() autoscalingiface.AutoScalingAPI { return p.asg }

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

// CloudWatchLogs returns a representation of the CloudWatchLogs API.
func (p ProviderServices) CloudWatchLogs() cloudwatchlogsiface.CloudWatchLogsAPI {
	return p.cloudwatchlogs
}

// Region returns provider-level region setting
func (p ProviderServices) Region() string { return p.spec.Region }

// Profile returns provider-level profile name
func (p ProviderServices) Profile() string { return p.spec.Profile }

// WaitTimeout returns provider-level duration after which any wait operation has to timeout
func (p ProviderServices) WaitTimeout() time.Duration { return p.spec.WaitTimeout }

func (p ProviderServices) ConfigProvider() client.ConfigProvider {
	return p.session
}

func (p ProviderServices) Session() *session.Session {
	return p.session
}

// ClusterInfo provides information about the cluster.
type ClusterInfo struct {
	Cluster *awseks.Cluster
}

// ProviderStatus stores information about the used IAM role and the resulting session
type ProviderStatus struct {
	iamRoleARN   string
	sessionCreds *credentials.Credentials
	ClusterInfo  *ClusterInfo
}

// New creates a new setup of the used AWS APIs
func New(spec *api.ProviderConfig, clusterSpec *api.ClusterConfig) (*ClusterProvider, error) {
	provider := &ProviderServices{
		spec: spec,
	}
	c := &ClusterProvider{
		Provider: provider,
	}
	// Create a new session and save credentials for possible
	// later re-use if overriding sessions due to custom URL
	s := c.newSession(spec)

	cache := os.Getenv(ekscreds.EksctlGlobalEnableCachingEnvName)
	if s.Config != nil && cache != "" {
		if cachedProvider, err := ekscreds.NewFileCacheProvider(spec.Profile, s.Config.Credentials, &ekscreds.RealClock{}); err == nil {
			s.Config.Credentials = credentials.NewCredentials(&cachedProvider)
		} else {
			logger.Warning("Failed to use cached provider: ", err)
		}
	}

	provider.session = s
	provider.asg = autoscaling.New(s)
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
	provider.cloudwatchlogs = cloudwatchlogs.New(s)

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

	return c, c.checkAuth()
}

// ParseConfig parses data into a ClusterConfig
func ParseConfig(data []byte) (*api.ClusterConfig, error) {
	// strict mode is not available in runtime.Decode, so we use the parser
	// directly; we don't store the resulting object, this is just the means
	// of detecting any unknown keys
	// NOTE: we must use sigs.k8s.io/yaml, as it behaves differently from
	// github.com/ghodss/yaml, which didn't handle nested structs well
	if err := yaml.UnmarshalStrict(data, &api.ClusterConfig{}); err != nil {
		return nil, err
	}

	obj, err := runtime.Decode(scheme.Codecs.UniversalDeserializer(), data)
	if err != nil {
		return nil, err
	}

	cfg, ok := obj.(*api.ClusterConfig)
	if !ok {
		return nil, fmt.Errorf("expected to decode object of type %T; got %T", &api.ClusterConfig{}, cfg)
	}
	return cfg, nil
}

// LoadConfigFromFile loads ClusterConfig from configFile
func LoadConfigFromFile(configFile string) (*api.ClusterConfig, error) {
	data, err := readConfig(configFile)
	if err != nil {
		return nil, errors.Wrapf(err, "reading config file %q", configFile)
	}
	clusterConfig, err := ParseConfig(data)
	if err != nil {
		return nil, errors.Wrapf(err, "loading config file %q", configFile)
	}
	return clusterConfig, nil

}

func readConfig(configFile string) ([]byte, error) {
	if configFile == "-" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(configFile)
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

// checkAuth checks the AWS authentication
func (c *ClusterProvider) checkAuth() error {

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

// ResolveAMI ensures that the node AMI is set and is available
func ResolveAMI(provider api.ClusterProvider, version string, np api.NodePool) error {
	var resolver ami.Resolver
	ng := np.BaseNodeGroup()
	switch ng.AMI {
	case api.NodeImageResolverAuto:
		resolver = ami.NewAutoResolver(provider.EC2())
	case api.NodeImageResolverAutoSSM:
		resolver = ami.NewSSMResolver(provider.SSM())
	case "":
		resolver = ami.NewMultiResolver(
			ami.NewSSMResolver(provider.SSM()),
			ami.NewAutoResolver(provider.EC2()),
		)
	default:
		return errors.Errorf("invalid AMI value: %q", ng.AMI)
	}

	instanceType := api.SelectInstanceType(np)
	id, err := resolver.Resolve(provider.Region(), version, instanceType, ng.AMIFamily)
	if err != nil {
		return errors.Wrap(err, "unable to determine AMI to use")
	}
	if id == "" {
		return ami.NewErrFailedResolution(provider.Region(), version, instanceType, ng.AMIFamily)
	}
	ng.AMI = id
	return nil
}

func errTooFewAvailabilityZones(azs []string) error {
	return fmt.Errorf("only %d zones specified %v, %d are required (can be non-unique)", len(azs), azs, api.MinRequiredAvailabilityZones)
}

// SetAvailabilityZones sets the given (or chooses) the availability zones
func (c *ClusterProvider) SetAvailabilityZones(spec *api.ClusterConfig, given []string) error {
	if count := len(given); count != 0 {
		if count < api.MinRequiredAvailabilityZones {
			return errTooFewAvailabilityZones(given)
		}
		spec.AvailabilityZones = given
		return nil
	}

	if count := len(spec.AvailabilityZones); count != 0 {
		if count < api.MinRequiredAvailabilityZones {
			return errTooFewAvailabilityZones(spec.AvailabilityZones)
		}
		return nil
	}

	logger.Debug("determining availability zones")
	zones, err := az.GetAvailabilityZones(c.Provider.EC2(), c.Provider.Region())
	if err != nil {
		return errors.Wrap(err, "getting availability zones")
	}

	logger.Info("setting availability zones to %v", zones)
	spec.AvailabilityZones = zones

	return nil
}

func (c *ClusterProvider) newSession(spec *api.ProviderConfig) *session.Session {
	// we might want to use bits from kops, although right now it seems like too many things we
	// don't want yet
	// https://github.com/kubernetes/kops/blob/master/upup/pkg/fi/cloudup/awsup/aws_cloud.go#L179
	config := aws.NewConfig().WithCredentialsChainVerboseErrors(true)

	if c.Provider.Region() != "" {
		config = config.WithRegion(c.Provider.Region()).WithSTSRegionalEndpoint(endpoints.RegionalSTSEndpoint)
	}

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
func (c *ClusterProvider) NewStackManager(spec *api.ClusterConfig) manager.StackManager {
	return manager.NewStackCollection(c.Provider, spec)
}
