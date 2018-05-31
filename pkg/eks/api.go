package eks

import (
	"os"
	"sync"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/sts"

	"github.com/kubicorn/kubicorn/pkg/logger"
)

const ClusterNameTag = "eksctl.cluster.k8s.io/v1alpha1/cluster-name"

type ClusterProvider struct {
	cfg *ClusterConfig
	svc *providerServices
}

type providerServices struct {
	cfn *cloudformation.CloudFormation
	eks *eks.EKS
	ec2 *ec2.EC2
	sts *sts.STS
	arn string
}

// simple config, to be replaced with Cluster API
type ClusterConfig struct {
	Region      string
	ClusterName string
	NodeAMI     string
	NodeType    string
	Nodes       int
	MinNodes    int
	MaxNodes    int

	SSHPublicKeyPath string
	SSHPublicKey     []byte

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
	// we might want to use bits from kops, although right now it seems like too many thing we
	// don't want yet
	// https://github.com/kubernetes/kops/blob/master/upup/pkg/fi/cloudup/awsup/aws_cloud.go#L179
	config := aws.NewConfig()
	config = config.WithRegion(clusterConfig.Region)
	config = config.WithCredentialsChainVerboseErrors(true)

	s := session.Must(session.NewSession(config))

	cfn := &ClusterProvider{
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
		s := session.Must(session.NewSession(config.WithEndpoint(endpoint)))
		cfn.svc.cfn = cloudformation.New(s)
	}
	if endpoint, ok := os.LookupEnv("AWS_EKS_ENDPOINT"); ok {
		s := session.Must(session.NewSession(config.WithEndpoint(endpoint)))
		cfn.svc.eks = eks.New(s)
	}
	if endpoint, ok := os.LookupEnv("AWS_EC2_ENDPOINT"); ok {
		s := session.Must(session.NewSession(config.WithEndpoint(endpoint)))
		cfn.svc.ec2 = ec2.New(s)
	}
	if endpoint, ok := os.LookupEnv("AWS_STS_ENDPOINT"); ok {
		s := session.Must(session.NewSession(config.WithEndpoint(endpoint)))
		cfn.svc.sts = sts.New(s)
	}

	return cfn
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
		"createControlPlane":          func(errs chan error) error { return c.createControlPlane(errs) },
		"createStackDefaultNodeGroup": func(errs chan error) error { return c.createStackDefaultNodeGroup(errs) },
	}, taskErrs)
	close(taskErrs)
}
