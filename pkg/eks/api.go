package eks

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"

	"github.com/kubicorn/kubicorn/pkg/logger"
)

const (
	ClusterNameTag = "eksctl.cluster.k8s.io/v1alpha1/cluster-name"
	AWSDebugLevel  = 5
)

type ClusterProvider struct {
	cfg *ClusterConfig
	svc *providerServices
}

type providerServices struct {
	cfn cloudformationiface.CloudFormationAPI
	eks eksiface.EKSAPI
	ec2 ec2iface.EC2API
	sts stsiface.STSAPI
	arn string
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

	SSHPublicKeyPath string
	SSHPublicKey     []byte

	AWSOperationTimeout time.Duration

	keyName        string
	clusterRoleARN string
	securityGroup  string
	subnetsList    string
	clusterVPC     string

	nodeInstanceRoleARN string

	MasterEndpoint           string
	CertificateAuthorityData []byte
}

func New(clusterConfig *ClusterConfig) *ClusterProvider {

	// Create a new session and save credentials for possible
	// later re-use if overriding sessions due to custom URL
	s := newSession(clusterConfig, "", nil)
	creds := s.Config.Credentials

	c := &ClusterProvider{
		cfg: clusterConfig,
		svc: &providerServices{
			cfn: cloudformation.New(s),
			eks: eks.New(s),
			ec2: ec2.New(s),
			sts: sts.New(s),
		},
	}

	// override sessions if any custom endpoints specified
	if endpoint, ok := os.LookupEnv("AWS_CLOUDFORMATION_ENDPOINT"); ok {
		logger.Debug("Setting CloudFormation endpoint to %s", endpoint)
		s := newSession(clusterConfig, endpoint, creds)
		c.svc.cfn = cloudformation.New(s)
	}
	if endpoint, ok := os.LookupEnv("AWS_EKS_ENDPOINT"); ok {
		logger.Debug("Setting EKS endpoint to %s", endpoint)
		s := newSession(clusterConfig, endpoint, creds)
		c.svc.eks = eks.New(s)
	}
	if endpoint, ok := os.LookupEnv("AWS_EC2_ENDPOINT"); ok {
		logger.Debug("Setting EC2 endpoint to %s", endpoint)
		s := newSession(clusterConfig, endpoint, creds)
		c.svc.ec2 = ec2.New(s)
	}
	if endpoint, ok := os.LookupEnv("AWS_STS_ENDPOINT"); ok {
		logger.Debug("Setting STS endpoint to %s", endpoint)
		s := newSession(clusterConfig, endpoint, creds)
		c.svc.sts = sts.New(s)
	}

	return c
}

func (c *ClusterProvider) CheckAuth() error {
	{
		input := &sts.GetCallerIdentityInput{}
		output, err := c.svc.sts.GetCallerIdentity(input)
		if err != nil {
			return errors.Wrap(err, "checking AWS STS access – cannot get role ARN for current session")
		}
		c.svc.arn = *output.Arn
		logger.Debug("role ARN for the current session is %q", c.svc.arn)
	}
	{
		input := &cloudformation.ListStacksInput{}
		if _, err := c.svc.cfn.ListStacks(input); err != nil {
			return errors.Wrap(err, "checking AWS CloudFormation access – cannot list stacks")
		}
	}
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
		"createControlPlane": func(errs chan error) error { return c.createControlPlane(errs) },
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

	if len(endpoint) > 0 {
		opts.Config.Endpoint = &endpoint
	}

	if credentials != nil {
		opts.Config.Credentials = credentials
	}

	return session.Must(session.NewSessionWithOptions(opts))
}
