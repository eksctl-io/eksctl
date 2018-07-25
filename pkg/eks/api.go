package eks

import (
	"fmt"
	"os"
	"sync"
	"time"

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

const (
	ClusterNameTag = "eksctl.cluster.k8s.io/v1alpha1/cluster-name"
	AWSDebugLevel  = 5
	//RequiredAvailabilityZones = 3
)

var DefaultWaitTimeout = 20 * time.Minute

type ClusterProvider struct {
	// core fields used for config and AWS APIs
	Spec     *ClusterConfig
	Provider Provider
	// informative fields, i.e. used as outputs
	Status *ProviderStatus
}

type Provider interface {
	CloudFormation() cloudformationiface.CloudFormationAPI
	EKS() eksiface.EKSAPI
	EC2() ec2iface.EC2API
	STS() stsiface.STSAPI
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

// simple config, to be replaced with Cluster API
type ClusterConfig struct {
	Region      string
	Profile     string
	ClusterName string
	NodeAMI     string
	NodeType    string
	Nodes       int
	MinNodes    int
	MaxNodes    int
	PolicyARNs  []string

	SSHPublicKeyPath string
	SSHPublicKey     []byte

	WaitTimeout time.Duration

	keyName        string
	clusterRoleARN string
	securityGroup  string
	subnetsList    string
	clusterVPC     string

	nodeInstanceRoleARN string

	MasterEndpoint           string
	CertificateAuthorityData []byte
	ClusterARN               string

	Addons ClusterAddons
}

type ClusterAddons struct {
	WithIAM AddonIAM
}

type AddonIAM struct {
	PolicyAmazonEC2ContainerRegistryPowerUser bool
}

func New(clusterConfig *ClusterConfig) *ClusterProvider {
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

func (c *ClusterProvider) SetAvailabilityZones(given []string) error {
	if len(given) == 0 {
		logger.Debug("determining availability zones")
		azSelector := az.NewSelectorWithDefaults(c.Provider.EC2())
		zones, err := azSelector.SelectZones(c.Spec.Region)
		if err != nil {
			return errors.Wrap(err, "getting availability zones")
		}

		logger.Info("setting availability zones to %v", zones)
		c.Status.availabilityZones = zones
		return nil
	}
	if len(given) < az.DefaultRequiredAvailabilityZones {
		return fmt.Errorf("only %d zones specified %v, %d are required (can be non-unque)", len(given), given, az.DefaultRequiredAvailabilityZones)
	}
	c.Status.availabilityZones = given
	return nil
}

func (c *ClusterProvider) runCreateTask(tasks map[string]func(chan error) error, taskErrs chan error) {
	wg := &sync.WaitGroup{}
	wg.Add(len(tasks))
	for taskName := range tasks {
		task := tasks[taskName]
		go func(tn string) {
			defer wg.Done()
			logger.Debug("task %q started", tn)
			errs := make(chan error)
			if err := task(errs); err != nil {
				taskErrs <- err
				return
			}
			if err := <-errs; err != nil {
				taskErrs <- err
				return
			}
			logger.Debug("task %q returned without errors", tn)
		}(taskName)
	}
	logger.Debug("waiting for %d tasks to complete", len(tasks))
	wg.Wait()
}

func (c *ClusterProvider) CreateCluster(taskErrs chan error) {
	c.runCreateTask(map[string]func(chan error) error{
		"createStackServiceRole": func(errs chan error) error { return c.createStackServiceRole(errs) },
		"createStackVPC":         func(errs chan error) error { return c.createStackVPC(errs) },
	}, taskErrs)
	c.runCreateTask(map[string]func(chan error) error{
		"createStackControlPlane": func(errs chan error) error { return c.createStackControlPlane(errs) },
	}, taskErrs)
	c.runCreateTask(map[string]func(chan error) error{
		"createStackDefaultNodeGroup": func(errs chan error) error { return c.createStackDefaultNodeGroup(errs) },
	}, taskErrs)
	close(taskErrs)
}

func newSession(clusterConfig *ClusterConfig, endpoint string, credentials *credentials.Credentials) *session.Session {
	// we might want to use bits from kops, although right now it seems like too many thing we
	// don't want yet
	// https://github.com/kubernetes/kops/blob/master/upup/pkg/fi/cloudup/awsup/aws_cloud.go#L179
	config := aws.NewConfig()
	config = config.WithRegion(clusterConfig.Region)
	config = config.WithCredentialsChainVerboseErrors(true)
	if logger.Level >= AWSDebugLevel {
		config = config.WithLogLevel(aws.LogDebug |
			aws.LogDebugWithHTTPBody |
			aws.LogDebugWithRequestRetries |
			aws.LogDebugWithRequestErrors |
			aws.LogDebugWithEventStreamBody)
		config = config.WithLogLevel(aws.LogDebugWithHTTPBody)
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
