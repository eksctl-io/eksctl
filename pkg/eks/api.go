package eks

import (
	"fmt"
	"os"
	"reflect"
	"strings"
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
	Interactive bool // for interactive use, i.e. eksctl

	Region      string `flag:"--regions"`
	ClusterName string `flag:"--cluster-name"`
	NodeOS      string `flag:"--node-os"`
	NodeAMI     string `flag:"--node-ami"`
	NodeType    string `flag:"--node-type"`
	Nodes       int    `flag:"--nodes"`
	MinNodes    int    `flag:"--nodes-min"`
	MaxNodes    int    `flag:"--nodes-max"`

	SSHPublicKeyPath string `flag:"--ssh-public-key"`
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

const (
	DEFAULT_NODE_COUNT = 2
	DEFAULT_NODE_TYPE  = "m5.large"

	REGION_US_WEST_2 = "us-west-2"
	REGION_US_EAST_2 = "us-east-2"

	DEFAULT_REGION = REGION_US_WEST_2

	NODE_OS_AMAZON_LINUX_2 = "Amazon Linux 2"

	DEFAULT_NODE_OS = NODE_OS_AMAZON_LINUX_2
)

var regionalAMIs = map[string]map[string]string{
	REGION_US_WEST_2: {
		NODE_OS_AMAZON_LINUX_2: "ami-73a6e20b",
	},
	REGION_US_EAST_2: {
		NODE_OS_AMAZON_LINUX_2: "ami-dea4d5a1",
	},
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

func (c *ClusterConfig) getFlagOrField(fieldName string) string {
	if !c.Interactive {
		return fieldName
	}
	field, ok := reflect.TypeOf(c).Elem().FieldByName(fieldName)
	if !ok {
		return fieldName
	}
	if flag := field.Tag.Get("flag"); flag != "" {
		return flag
	}
	return fieldName
}

func (c *ClusterProvider) CheckConfig() error {

	regionalAMIs, ok := regionalAMIs[c.cfg.Region]
	if !ok {
		supportedRegions := []string{}
		for r := range regionalAMIs {
			supportedRegions = append(supportedRegions, r)
		}
		return fmt.Errorf("unsupported %s %q, EKS is only availabe in the following regions: %s",
			c.cfg.getFlagOrField("Region"),
			c.cfg.Region,
			strings.Join(supportedRegions, ", "),
		)
	}

	c.cfg.NodeOS = DEFAULT_NODE_OS // will expose when more OSs are available

	if c.cfg.NodeAMI == "" {
		c.cfg.NodeAMI, ok = regionalAMIs[c.cfg.NodeOS]
		if !ok {
			return fmt.Errorf("unsuported %s %q, use %s to set custom AMI",
				c.cfg.getFlagOrField("NodeOS"),
				c.cfg.NodeOS,
				c.cfg.getFlagOrField("NodeAMI"),
			)
		}
	}

	if c.cfg.ClusterName == "" {
		return fmt.Errorf("%s must be set", c.cfg.getFlagOrField("ClusterName"))
	}
	return nil
}

func (c *ClusterProvider) CheckNodeCountConfig() error {
	// this is separate from CheckConfig as we would only call it on create
	if c.cfg.MinNodes < 0 || c.cfg.MaxNodes < 0 || c.cfg.Nodes < 0 {
		return fmt.Errorf("%s, %s or %s cannot be less than zero",
			c.cfg.getFlagOrField("MinNodes"),
			c.cfg.getFlagOrField("MaxNodes"),
			c.cfg.getFlagOrField("Nodes"),
		)
	}

	if c.cfg.MinNodes != 0 && c.cfg.MaxNodes != 0 && c.cfg.Nodes != 0 {
		return fmt.Errorf("%s, %s and %s cannot be specified all at the same time",
			c.cfg.getFlagOrField("MinNodes"),
			c.cfg.getFlagOrField("MaxNodes"),
			c.cfg.getFlagOrField("Nodes"),
		)
	}

	if c.cfg.MinNodes == 0 && c.cfg.MaxNodes == 0 {
		// defaults
		c.cfg.MinNodes = c.cfg.Nodes
		c.cfg.MaxNodes = c.cfg.Nodes
	} else {
		// ambiguities
		if c.cfg.MinNodes > c.cfg.MaxNodes {
			return fmt.Errorf("%s cannot be greater than %s",
				c.cfg.getFlagOrField("MinNodes"),
				c.cfg.getFlagOrField("MaxNodes"),
			)
		}
		if c.cfg.MinNodes > 0 && c.cfg.MaxNodes == 0 && c.cfg.Nodes > 0 {
			c.cfg.MaxNodes = c.cfg.Nodes
		}
	}

	return nil
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
