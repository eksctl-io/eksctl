package eks

import (
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/kubicorn/kubicorn/pkg/logger"

	"github.com/pkg/errors"
)

//go:generate go-bindata -pkg $GOPACKAGE -prefix ../../vendor/1.10.0/2018-05-09 -o cfn_templates.go ../../vendor/1.10.0/2018-05-09

const ClusterNameTag = "eksctl.cluster.k8s.io/v1alpha1/cluster-name"

type CloudFormation struct {
	cfg *Config
	svc *cloudformation.CloudFormation
	ec2 *ec2.EC2
}

// simple config, to be replaced with Cluster API
type Config struct {
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
}
type Stack = cloudformation.Stack

func New(clusterConfig *Config) *CloudFormation {
	// we might want to use bits from kops, although right now it seems like too many thing we
	// don't want yet
	// https://github.com/kubernetes/kops/blob/master/upup/pkg/fi/cloudup/awsup/aws_cloud.go#L179
	config := aws.NewConfig().WithRegion(clusterConfig.Region)
	config = config.WithCredentialsChainVerboseErrors(true)
	session := session.Must(session.NewSession(config))

	return &CloudFormation{
		cfg: clusterConfig,
		svc: cloudformation.New(session),
		ec2: ec2.New(session),
	}
}

func (c *CloudFormation) CheckAuth() error {
	input := &cloudformation.ListStacksInput{}
	if _, err := c.svc.ListStacks(input); err != nil {
		return errors.Wrap(err, "checking AWS CloudFormation access")
	}
	return nil
}

func (c *CloudFormation) CreateStack(name string, templateBody []byte, parameters map[string]string, withIAM bool, stack chan Stack, errs chan error) error {
	input := &cloudformation.CreateStackInput{}
	input.SetStackName(name)
	input.SetTags([]*cloudformation.Tag{
		&cloudformation.Tag{
			Key:   aws.String(ClusterNameTag),
			Value: aws.String(c.cfg.ClusterName),
		},
	})
	input.SetTemplateBody(string(templateBody))
	if withIAM {
		input.SetCapabilities(aws.StringSlice([]string{cloudformation.CapabilityCapabilityIam}))
	}
	for k, v := range parameters {
		p := &cloudformation.Parameter{
			ParameterKey:   aws.String(k),
			ParameterValue: aws.String(v),
		}
		input.Parameters = append(input.Parameters, p)
	}

	// TODO(p0): looks like we can block on this forever, if parameters are invalid
	logger.Debug("input = %#v", input)
	s, err := c.svc.CreateStack(input)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("creating CloudFormation stack %q", name))
	}
	logger.Debug("stack = %#v", s)

	go func() {
		// TODO(eksctld): should probably use SNS notifications instead of polling
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		defer close(errs)
		for {
			select {
			case <-ticker.C:
				s, err := c.describeStack(&name)
				if err != nil {
					log.Print(err)
					continue
				}
				logger.Debug("stack = %#v", s)
				switch *s.StackStatus {
				case cloudformation.StackStatusCreateInProgress:
					continue
				case cloudformation.StackStatusCreateComplete:
					errs <- nil
					stack <- *s
					return
				case cloudformation.StackStatusCreateFailed:
					fallthrough
					// TODO: technically, any of these may occur, but we may want to ignore some of these
					// case cloudformation.StackStatusRollbackInProgress:
					// case cloudformation.StackStatusRollbackFailed:
					// case cloudformation.StackStatusRollbackComplete:
					// case cloudformation.StackStatusDeleteInProgress:
					// case cloudformation.StackStatusDeleteFailed:
					// case cloudformation.StackStatusDeleteComplete:
					// case cloudformation.StackStatusUpdateInProgress:
					// case cloudformation.StackStatusUpdateCompleteCleanupInProgress:
					// case cloudformation.StackStatusUpdateComplete:
					// case cloudformation.StackStatusUpdateRollbackInProgress:
					// case cloudformation.StackStatusUpdateRollbackFailed:
					// case cloudformation.StackStatusUpdateRollbackCompleteCleanupInProgress:
					// case cloudformation.StackStatusUpdateRollbackComplete:
					// case cloudformation.StackStatusReviewInProgress:
				default:
					errs <- fmt.Errorf("creating CloudFormation stack %q: %s", name, *s.StackStatus)
					// stack <- *s // this usually results in closed channel panic, but we don't need it really
					logger.Debug("stack = %#v", s)
					return
				}
			}
		}
	}()

	return nil

}

func (c *CloudFormation) describeStack(name *string) (*Stack, error) {
	input := &cloudformation.DescribeStacksInput{
		StackName: name,
	}
	resp, err := c.svc.DescribeStacks(input)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("describing CloudFormation stack %q", *name))
	}
	return resp.Stacks[0], nil
}

func (c *CloudFormation) ListReadyStacks(nameRegex string) ([]*Stack, error) {
	var (
		subErr error
		stack  *Stack
	)

	re := regexp.MustCompile(nameRegex)
	input := &cloudformation.ListStacksInput{
		StackStatusFilter: aws.StringSlice([]string{cloudformation.StackStatusCreateComplete}),
	}
	stacks := []*Stack{}

	pager := func(p *cloudformation.ListStacksOutput, last bool) (shouldContinue bool) {
		for _, s := range p.StackSummaries {
			if re.MatchString(*s.StackName) {
				stack, subErr = c.describeStack(s.StackName)
				if subErr != nil {
					return false
				}
				stacks = append(stacks, stack)
			}
		}
		return true
	}
	if err := c.svc.ListStacksPages(input, pager); err != nil {
		return nil, err
	}
	if subErr != nil {
		return nil, subErr
	}
	return stacks, nil
}

func (c *CloudFormation) CreateStacks(tasks map[string]func(chan error) error, taskErrs chan error) {
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
	close(taskErrs)
}

func (c *CloudFormation) CreateCoreStacks(taskErrs chan error) {
	c.CreateStacks(map[string]func(chan error) error{
		"createStackServiceRole": func(errs chan error) error { return c.createStackServiceRole(errs) },
		"createStackVPC":         func(errs chan error) error { return c.createStackVPC(errs) },
	}, taskErrs)
}

func (c *CloudFormation) CreateNodeGroupStack(taskErrs chan error) {
	c.CreateStacks(map[string]func(chan error) error{
		"createStackDefaultNodeGroup": func(errs chan error) error { return c.createStackDefaultNodeGroup(errs) },
	}, taskErrs)
}

func (c *CloudFormation) stackNameVPC() string {
	return "EKS-" + c.cfg.ClusterName + "-VPC"
}

func (c *CloudFormation) stackParamsVPC() map[string]string {
	return map[string]string{
		"ClusterName": c.cfg.ClusterName,
	}
}

func (c *CloudFormation) createStackVPC(errs chan error) error {
	name := c.stackNameVPC()
	logger.Info("creating VPC stack %q", name)
	templateBody, err := amazonEksVpcSampleYamlBytes()
	if err != nil {
		return errors.Wrap(err, "decompressing bundled template for VPC stack")
	}

	stackChan := make(chan Stack)
	taskErrs := make(chan error)

	if err := c.CreateStack(name, templateBody, c.stackParamsVPC(), false, stackChan, taskErrs); err != nil {
		return err
	}

	go func() {
		defer close(errs)
		defer close(stackChan)

		if err := <-taskErrs; err != nil {
			errs <- err
			return
		}

		s := <-stackChan

		logger.Debug("created VPC stack %q – processing outputs", name)

		securityGroup := GetOutput(&s, "SecurityGroup")
		if securityGroup == nil {
			errs <- fmt.Errorf("SecurityGroup is nil")
			return
		}
		c.cfg.securityGroup = *securityGroup

		subnetsList := GetOutput(&s, "SubnetsList")
		if subnetsList == nil {
			errs <- fmt.Errorf("SubnetsList is nil")
			return
		}
		c.cfg.subnetsList = *subnetsList

		clusterVPC := GetOutput(&s, "ClusterVPC")
		if clusterVPC == nil {
			errs <- fmt.Errorf("ClusterVPC is nil")
			return
		}
		c.cfg.clusterVPC = *clusterVPC

		logger.Debug("clusterConfig = %#v", c.cfg)
		logger.Success("created VPC stack %q", name)

		errs <- nil
	}()
	return nil
}

func (c *CloudFormation) DeleteStackVPC() error {
	name := c.stackNameVPC()
	s, err := c.describeStack(&name)
	if err != nil {
		return errors.Wrap(err, "not able to get VPC stack for deletion")
	}

	input := &cloudformation.DeleteStackInput{
		StackName: s.StackName,
	}

	if _, err := c.svc.DeleteStack(input); err != nil {
		return errors.Wrap(err, "not able to delete VPC stack")
	}
	return nil
}

func (c *CloudFormation) stackNameServiceRole() string {
	return "EKS-" + c.cfg.ClusterName + "-ServiceRole"
}

func (c *CloudFormation) createStackServiceRole(errs chan error) error {
	name := c.stackNameServiceRole()
	logger.Info("creating ServiceRole stack %q", name)
	templateBody, err := amazonEksServiceRoleYamlBytes()
	if err != nil {
		return errors.Wrap(err, "decompressing bundled template for ServiceRole stack")
	}

	stackChan := make(chan Stack)
	taskErrs := make(chan error)

	if err := c.CreateStack(name, templateBody, nil, true, stackChan, taskErrs); err != nil {
		return err
	}

	go func() {
		defer close(errs)
		defer close(stackChan)

		if err := <-taskErrs; err != nil {
			errs <- err
			return
		}

		s := <-stackChan

		logger.Debug("created ServiceRole stack %q – processing outputs", name)

		clusterRoleARN := GetOutput(&s, "RoleArn")
		if clusterRoleARN == nil {
			errs <- fmt.Errorf("RoleArn is nil")
			return
		}
		c.cfg.clusterRoleARN = *clusterRoleARN

		logger.Debug("clusterConfig = %#v", c.cfg)
		logger.Success("created ServiceRole stack %q", name)

		errs <- nil
	}()
	return nil
}

func (c *CloudFormation) DeleteStackServiceRole() error {
	name := c.stackNameServiceRole()
	s, err := c.describeStack(&name)
	if err != nil {
		return errors.Wrap(err, "not able to get ServiceRole stack for deletion")
	}

	input := &cloudformation.DeleteStackInput{
		StackName: s.StackName,
	}

	if _, err := c.svc.DeleteStack(input); err != nil {
		return errors.Wrap(err, "not able to delete ServiceRole stack")
	}
	return nil
}

func (c *CloudFormation) stackNameDefaultNodeGroup() string {
	return "EKS-" + c.cfg.ClusterName + "-DefaultNodeGroup"
}
func (c *CloudFormation) stackParamsDefaultNodeGroup() map[string]string {
	regionalAMIs := map[string]string{
		"us-west-2": "ami-993141e1",
	}

	if c.cfg.NodeAMI == "" {
		c.cfg.NodeAMI = regionalAMIs[c.cfg.Region]
	}

	if c.cfg.MinNodes == 0 && c.cfg.MaxNodes == 0 {
		c.cfg.MinNodes = c.cfg.Nodes
		c.cfg.MaxNodes = c.cfg.Nodes
	}

	return map[string]string{
		"ClusterName":                      c.cfg.ClusterName,
		"NodeGroupName":                    "default",
		"KeyName":                          c.cfg.keyName,
		"NodeImageId":                      c.cfg.NodeAMI,
		"NodeInstanceType":                 c.cfg.NodeType,
		"NodeAutoScalingGroupMinSize":      fmt.Sprintf("%d", c.cfg.MinNodes),
		"NodeAutoScalingGroupMaxSize":      fmt.Sprintf("%d", c.cfg.MaxNodes),
		"ClusterControlPlaneSecurityGroup": c.cfg.securityGroup,
		"Subnets":                          c.cfg.subnetsList,
		"VpcId":                            c.cfg.clusterVPC,
	}
}

func (c *CloudFormation) createStackDefaultNodeGroup(errs chan error) error {
	name := c.stackNameDefaultNodeGroup()
	logger.Info("creating DefaultNodeGroup stack %q", name)
	templateBody, err := amazonEksNodegroupYamlBytes()
	if err != nil {
		return errors.Wrap(err, "decompressing bundled template for DefaultNodeGroup stack")
	}

	stackChan := make(chan Stack)
	taskErrs := make(chan error)

	if err := c.CreateStack(name, templateBody, c.stackParamsDefaultNodeGroup(), true, stackChan, taskErrs); err != nil {
		return err
	}

	go func() {
		defer close(errs)
		defer close(stackChan)

		if err := <-taskErrs; err != nil {
			errs <- err
			return
		}

		s := <-stackChan

		logger.Debug("created DefaultNodeGroup stack %q – processing outputs", name)

		nodeInstanceRoleARN := GetOutput(&s, "NodeInstanceRole")
		if nodeInstanceRoleARN == nil {
			// TODO(p2): confirm if this can actually block if key was wrong and find out why
			errs <- fmt.Errorf("NodeInstanceRole is nil")
			return
		}
		c.cfg.nodeInstanceRoleARN = *nodeInstanceRoleARN

		logger.Debug("clusterConfig = %#v", c.cfg)
		logger.Success("created DefaultNodeGroup stack %q", name)

		errs <- nil
	}()
	return nil
}

func (c *CloudFormation) DeleteStackDefaultNodeGroup() error {
	name := c.stackNameDefaultNodeGroup()
	s, err := c.describeStack(&name)
	if err != nil {
		return errors.Wrap(err, "not able to get DefaultNodeGroup stack for deletion")
	}

	input := &cloudformation.DeleteStackInput{
		StackName: s.StackName,
	}

	if _, err := c.svc.DeleteStack(input); err != nil {
		return errors.Wrap(err, "not able to delete DefaultNodeGroup stack")
	}
	return nil
}

func GetOutput(stack *Stack, key string) *string {
	for _, x := range stack.Outputs {
		if *x.OutputKey == key {
			return x.OutputValue
		}
	}
	return nil
}
