package eks

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gofrs/flock"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"k8s.io/apimachinery/pkg/runtime"
	k8sclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/yaml"

	"github.com/weaveworks/eksctl/pkg/ami"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/az"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	ekscreds "github.com/weaveworks/eksctl/pkg/credentials"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
	"github.com/weaveworks/eksctl/pkg/utils/nodes"
	"github.com/weaveworks/eksctl/pkg/version"
)

// ClusterProvider stores information about the cluster
type ClusterProvider struct {
	// KubeProvider offers helper methods to handle Kubernetes operations
	KubeProvider

	// core fields used for config and AWS APIs
	AWSProvider api.ClusterProvider
	// informative fields, i.e. used as outputs
	Status *ProviderStatus
}

// KubernetesProvider provides helper methods to handle Kubernetes operations.
type KubernetesProvider struct {
	WaitTimeout time.Duration
	RoleARN     string
	Signer      api.STSPresigner
}

// KubeProvider is an interface with helper funcs for k8s and EKS that are part of ClusterProvider
//
//go:generate counterfeiter -o fakes/fake_kube_provider.go . KubeProvider
type KubeProvider interface {
	NewRawClient(clusterInfo kubeconfig.ClusterInfo) (*kubernetes.RawClient, error)
	NewStdClientSet(clusterInfo kubeconfig.ClusterInfo) (k8sclient.Interface, error)
	ServerVersion(rawClient *kubernetes.RawClient) (string, error)
	WaitForControlPlane(meta *api.ClusterMeta, clientSet *kubernetes.RawClient, waitTimeout time.Duration) error
}

// ProviderServices stores the used APIs
type ProviderServices struct {
	spec *api.ProviderConfig
	asg  awsapi.ASG

	cloudtrail     awsapi.CloudTrail
	cloudwatchlogs awsapi.CloudWatchLogs
	session        *session.Session

	*ServicesV2
}

// CloudFormationRoleARN returns, if any, a service role used by CloudFormation to call AWS API on your behalf
func (p ProviderServices) CloudFormationRoleARN() string { return p.spec.CloudFormationRoleARN }

// CloudFormationDisableRollback returns whether stacks should not rollback on failure
func (p ProviderServices) CloudFormationDisableRollback() bool {
	return p.spec.CloudFormationDisableRollback
}

// ASG returns a representation of the AutoScaling API
func (p ProviderServices) ASG() awsapi.ASG { return p.asg }

// CloudTrail returns a representation of the CloudTrail API
func (p ProviderServices) CloudTrail() awsapi.CloudTrail { return p.cloudtrail }

// CloudWatchLogs returns a representation of the CloudWatchLogs API.
func (p ProviderServices) CloudWatchLogs() awsapi.CloudWatchLogs {
	return p.cloudwatchlogs
}

// Region returns provider-level region setting
func (p ProviderServices) Region() string { return p.spec.Region }

// Profile returns the provider-level AWS profile.
func (p ProviderServices) Profile() api.Profile { return p.spec.Profile }

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
	Cluster *ekstypes.Cluster
}

// SessionProvider abstracts an aws credentials.Value provider.
type SessionProvider interface {
	Get() (credentials.Value, error)
}

// ProviderStatus stores information about the used IAM role and the resulting session
type ProviderStatus struct {
	IAMRoleARN   string
	ClusterInfo  *ClusterInfo
	SessionCreds SessionProvider
}

// New creates a new setup of the used AWS APIs
func New(ctx context.Context, spec *api.ProviderConfig, clusterSpec *api.ClusterConfig) (*ClusterProvider, error) {
	provider := &ProviderServices{
		spec: spec,
	}
	c := &ClusterProvider{
		AWSProvider: provider,
	}
	// Create a new session and save credentials for possible
	// later re-use if overriding sessions due to custom URL
	s := c.newSession(spec)

	cacheCredentials := os.Getenv(ekscreds.EksctlGlobalEnableCachingEnvName) != ""
	var (
		credentialsCacheFilePath string
		err                      error
	)
	if cacheCredentials {
		if s.Config == nil {
			return nil, errors.New("expected Session.Config to be non-nil")
		}
		credentialsCacheFilePath, err = ekscreds.GetCacheFilePath()
		if err != nil {
			return nil, fmt.Errorf("error getting cache file path: %w", err)
		}
		if cachedProvider, err := ekscreds.NewFileCacheProvider(spec.Profile.Name, s.Config.Credentials, &ekscreds.RealClock{}, afero.NewOsFs(), func(path string) ekscreds.Flock {
			return flock.New(path)
		}, credentialsCacheFilePath); err == nil {
			s.Config.Credentials = credentials.NewCredentials(&cachedProvider)
		} else {
			logger.Warning("Failed to use cached provider: ", err)
		}
	}

	provider.session = s

	cfg, err := newV2Config(spec, c.AWSProvider.Region(), credentialsCacheFilePath)
	if err != nil {
		return nil, err
	}

	provider.ServicesV2 = &ServicesV2{
		config: cfg,
	}

	c.Status = &ProviderStatus{
		SessionCreds: s.Config.Credentials,
	}

	provider.asg = autoscaling.NewFromConfig(cfg)
	provider.cloudwatchlogs = cloudwatchlogs.NewFromConfig(cfg)
	provider.cloudtrail = cloudtrail.NewFromConfig(cfg)

	if endpoint, ok := os.LookupEnv("AWS_CLOUDTRAIL_ENDPOINT"); ok {
		logger.Debug("Setting CloudTrail endpoint to %s", endpoint)
		provider.cloudtrail = cloudtrail.NewFromConfig(cfg, func(o *cloudtrail.Options) {
			o.EndpointResolver = cloudtrail.EndpointResolverFromURL(endpoint)
		})
	}

	stsOutput, err := c.checkAuth(ctx)
	if err != nil {
		return nil, err
	}

	// c.Status.IAMRoleARN is later needed by the kubeProvider
	c.Status.IAMRoleARN = *stsOutput.Arn
	logger.Debug("role ARN for the current session is %q", c.Status.IAMRoleARN)

	if clusterSpec != nil {
		clusterSpec.Metadata.AccountID = *stsOutput.Account
		clusterSpec.Metadata.Region = c.AWSProvider.Region()
	}

	kubeProvider := &KubernetesProvider{
		WaitTimeout: spec.WaitTimeout,
		RoleARN:     c.Status.IAMRoleARN,
		Signer:      provider.STSPresigner(),
	}
	c.KubeProvider = kubeProvider

	return c, nil
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
		if c.AWSProvider.Region() == supportedRegion {
			return true
		}
	}
	return false
}

// GetCredentialsEnv returns the AWS credentials for env usage
func (c *ClusterProvider) GetCredentialsEnv() ([]string, error) {
	creds, err := c.Status.SessionCreds.Get()
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
func (c *ClusterProvider) checkAuth(ctx context.Context) (*sts.GetCallerIdentityOutput, error) {
	output, err := c.AWSProvider.STS().GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, errors.Wrap(err, "checking AWS STS access – cannot get role ARN for current session")
	}
	if output == nil || output.Arn == nil {
		return nil, fmt.Errorf("unexpected response from AWS STS")
	}
	return output, nil
}

// ResolveAMI ensures that the node AMI is set and is available
func ResolveAMI(ctx context.Context, provider api.ClusterProvider, version string, np api.NodePool) error {
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
	id, err := resolver.Resolve(ctx, provider.Region(), version, instanceType, ng.AMIFamily)
	if err != nil {
		return errors.Wrap(err, "unable to determine AMI to use")
	}
	if id == "" {
		return ami.NewErrFailedResolution(provider.Region(), version, instanceType, ng.AMIFamily)
	}
	ng.AMI = id
	return nil
}

// SetAvailabilityZones sets the given (or chooses) the availability zones
// Returns whether azs were set randomly or provided by a user.
// CheckInstanceAvailability is only run if azs were provided by the user. Random
// selection already performs this check and makes sure AZs support all given instances.
func SetAvailabilityZones(ctx context.Context, spec *api.ClusterConfig, given []string, ec2API awsapi.EC2, region string) (bool, error) {
	if count := len(given); count != 0 {
		if count < api.MinRequiredAvailabilityZones {
			return false, api.ErrTooFewAvailabilityZones(given)
		}
		spec.AvailabilityZones = given
		return true, nil
	}

	if count := len(spec.AvailabilityZones); count != 0 {
		if count < api.MinRequiredAvailabilityZones {
			return false, api.ErrTooFewAvailabilityZones(spec.AvailabilityZones)
		}
		return true, nil
	}

	logger.Debug("determining availability zones")
	zones, err := az.GetAvailabilityZones(ctx, ec2API, region, spec)
	if err != nil {
		return false, errors.Wrap(err, "getting availability zones")
	}

	logger.Info("setting availability zones to %v", zones)
	spec.AvailabilityZones = zones

	return false, nil
}

// CheckInstanceAvailability verifies that if any instances are provided in any node groups
// that those instances are available in the selected AZs.
func CheckInstanceAvailability(ctx context.Context, spec *api.ClusterConfig, ec2API awsapi.EC2) error {
	logger.Debug("determining instance availability in zones")

	// This map will use either globally configured AZs or, if set, the AZ defined by the nodegroup.
	// map["ng-1"]["c2.large"]=[]string{"us-west-1a", "us-west-1b"}
	instanceMap := make(map[string]map[string][]string)
	uniqueInstances := sets.NewString()

	pool := nodes.ToNodePools(spec)
	for _, ng := range pool {
		if _, ok := instanceMap[ng.BaseNodeGroup().Name]; !ok {
			instanceMap[ng.BaseNodeGroup().Name] = make(map[string][]string)
		}
		for _, instanceType := range ng.InstanceTypeList() {
			if instanceType == "mixed" {
				continue
			}
			uniqueInstances.Insert(instanceType)
			if len(ng.BaseNodeGroup().AvailabilityZones) > 0 {
				instanceMap[ng.BaseNodeGroup().Name][instanceType] = ng.BaseNodeGroup().AvailabilityZones
			} else {
				instanceMap[ng.BaseNodeGroup().Name][instanceType] = spec.AvailabilityZones
			}
		}
	}

	// Do an early exit if we don't have anything.
	if uniqueInstances.Len() == 0 {
		// nothing to do
		return nil
	}

	var instanceTypeOfferings []ec2types.InstanceTypeOffering

	p := ec2.NewDescribeInstanceTypeOfferingsPaginator(ec2API, &ec2.DescribeInstanceTypeOfferingsInput{
		Filters: []ec2types.Filter{
			{
				Name:   awsv2.String("instance-type"),
				Values: uniqueInstances.List(),
			},
		},
		LocationType: ec2types.LocationTypeAvailabilityZone,
		MaxResults:   awsv2.Int32(100),
	})
	for p.HasMorePages() {
		output, err := p.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("unable to list offerings for instance types %w", err)
		}
		instanceTypeOfferings = append(instanceTypeOfferings, output.InstanceTypeOfferings...)
	}
	// construct a map so instance types can easily be checked
	// map["c2.large"]["us-east-1a"]=struct{}{}
	offers := make(map[string]map[string]struct{})
	for _, offer := range instanceTypeOfferings {
		if _, ok := offers[string(offer.InstanceType)]; !ok {
			offers[string(offer.InstanceType)] = map[string]struct{}{
				awsv2.ToString(offer.Location): {},
			}
		} else {
			offers[string(offer.InstanceType)][awsv2.ToString(offer.Location)] = struct{}{}
		}
	}
	// check if the instance type is available in at least one of the offered zones
	// per nodegroup.
	for k, v := range instanceMap {
		var (
			notAvailableIn []string
			available      bool
		)
		for instance, azs := range v {
			if zones, ok := offers[instance]; ok {
				for _, az := range azs {
					if _, ok := zones[az]; ok {
						available = true
						break
					} else {
						notAvailableIn = append(notAvailableIn, az)
					}
				}
			}
			if !available {
				return fmt.Errorf("none of the provided AZs %q support instance type %s in nodegroup %s", strings.Join(notAvailableIn, ","), instance, k)
			}
		}
	}

	return nil
}

// ValidateLocalZones validates that the specified local zones exist.
func ValidateLocalZones(ctx context.Context, ec2API awsapi.EC2, localZones []string, region string) error {
	output, err := ec2API.DescribeAvailabilityZones(ctx, &ec2.DescribeAvailabilityZonesInput{
		ZoneNames: localZones,
		Filters: []ec2types.Filter{
			{
				Name:   awsv2.String("region-name"),
				Values: []string{region},
			},
			{
				Name:   awsv2.String("state"),
				Values: []string{string(ec2types.AvailabilityZoneStateAvailable)},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("error describing availability zones: %w", err)
	}
	if len(output.AvailabilityZones) != len(localZones) {
		return fmt.Errorf("failed to find all local zones; expected to find %d available local zones but found only %d", len(localZones), len(output.AvailabilityZones))
	}
	for _, z := range output.AvailabilityZones {
		if *z.ZoneType != "local-zone" {
			return fmt.Errorf("non local-zone %q specified in localZones", *z.ZoneName)
		}
	}
	return nil
}

func (c *ClusterProvider) newSession(spec *api.ProviderConfig) *session.Session {
	// we might want to use bits from kops, although right now it seems like too many things we
	// don't want yet
	// https://github.com/kubernetes/kops/blob/master/upup/pkg/fi/cloudup/awsup/aws_cloud.go#L179
	config := aws.NewConfig().WithCredentialsChainVerboseErrors(true)

	if c.AWSProvider.Region() != "" {
		config = config.WithRegion(c.AWSProvider.Region()).WithSTSRegionalEndpoint(endpoints.RegionalSTSEndpoint)
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
		Profile:                 spec.Profile.Name,
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
	return manager.NewStackCollection(c.AWSProvider, spec)
}

// LoadClusterIntoSpecFromStack uses stack information to load the cluster
// configuration into the spec
// At the moment VPC and KubernetesNetworkConfig are respected
func (c *ClusterProvider) LoadClusterIntoSpecFromStack(ctx context.Context, spec *api.ClusterConfig, stack *manager.Stack) error {
	if err := c.LoadClusterVPC(ctx, spec, stack); err != nil {
		return err
	}
	return c.RefreshClusterStatus(ctx, spec)
}
