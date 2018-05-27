package cfn

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"

	"github.com/kubicorn/kubicorn/pkg/logger"

	"github.com/pkg/errors"
)

//go:generate go-bindata -pkg $GOPACKAGE -prefix ../../vendor/1.10.0/2018-05-09 ../../vendor/1.10.0/2018-05-09

const ClusterNameTag = "eksctl.cluster.k8s.io/v1alpha1/cluster-name"

type CloudFormation struct {
	svc *cloudformation.CloudFormation
	cfg *Config
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

	SSHPublicKey []byte

	keyName        string
	clusterRoleARN string
	securityGroup  string
	subnetsList    string
	clusterVPC     string
}

type Stack = cloudformation.Stack

func New(clusterConfig *Config) *CloudFormation {
	// we might want to use bits from kops, although right now it seems like too many thing we
	// don't want yet
	// https://github.com/kubernetes/kops/blob/master/upup/pkg/fi/cloudup/awsup/aws_cloud.go#L179
	config := aws.NewConfig().WithRegion(clusterConfig.Region)
	config = config.WithCredentialsChainVerboseErrors(true)

	return &CloudFormation{
		svc: cloudformation.New(session.Must(session.NewSession(config))),
		cfg: clusterConfig,
	}
}

func (c *CloudFormation) CheckAuth() error {
	input := &cloudformation.ListStacksInput{}
	if _, err := c.svc.ListStacks(input); err != nil {
		return errors.Wrap(err, "checking AWS CloudFormation access")
	}
	return nil
}

func (c *CloudFormation) CreateStack(name string, templateBody []byte, parameters map[string]string, withIAM bool, done chan error, stack chan Stack) error {
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

	_, err := c.svc.CreateStack(input)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("creating CloudFormation stack %q", name))
	}

	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		defer close(done)
		for {
			select {
			case <-ticker.C:
				s, err := c.describeStack(&name)
				if err != nil {
					log.Print(err)
					continue
				}
				switch *s.StackStatus {
				case cloudformation.StackStatusCreateInProgress:
					continue
				case cloudformation.StackStatusCreateComplete:
					done <- nil
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
					done <- fmt.Errorf("creating CloudFormation stack %q: %s", name, *s.StackStatus)
					stack <- *s
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

func (c *CloudFormation) stackParamsVPC() map[string]string {
	return map[string]string{
		"ClusterName": c.cfg.ClusterName,
	}
}

func (c *CloudFormation) CreateStackVPC(done chan error) error {
	name := "EKS-" + c.cfg.ClusterName + "-VPC"
	logger.Info("creating VPC stack %q", name)
	templateBody, err := amazonEksVpcSampleYamlBytes()
	if err != nil {
		return errors.Wrap(err, "decompressing bundled template for VPC stack")
	}

	stackChan := make(chan Stack)
	doneSubChan := make(chan error)

	if err := c.CreateStack(name, templateBody, c.stackParamsVPC(), false, doneSubChan, stackChan); err != nil {
		return err
	}

	go func() {
		defer close(done)
		defer close(stackChan)

		if err := <-doneSubChan; err != nil {
			done <- err
			return
		}
		done <- nil

		s := <-stackChan

		securityGroup := GetOutput(&s, "SecurityGroup")
		if securityGroup == nil {
			done <- fmt.Errorf("SecurityGroup is nil")
			return
		}
		c.cfg.securityGroup = *securityGroup

		subnetsList := GetOutput(&s, "SubnetsList")
		if subnetsList == nil {
			done <- fmt.Errorf("SubnetsList is nil")
			return
		}
		c.cfg.subnetsList = *subnetsList

		clusterVPC := GetOutput(&s, "ClusterVPC")
		if clusterVPC == nil {
			done <- fmt.Errorf("ClusterVPC is nil")
			return
		}
		c.cfg.clusterVPC = *clusterVPC

		logger.Debug("clusterConfig = %#v", c.cfg)
		logger.Success("created VPC stack %q", name)
	}()
	return nil
}

func (c *CloudFormation) DeleteStackVPC() error {
	s, err := c.GetStackVPC()
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
func (c *CloudFormation) CreateStackServiceRole(done chan error) error {
	name := "EKS-" + c.cfg.ClusterName + "-ServiceRole"
	logger.Info("creating ServiceRole stack %q", name)
	templateBody, err := amazonEksServiceRoleYamlBytes()
	if err != nil {
		return errors.Wrap(err, "decompressing bundled template for ServiceRole stack")
	}

	stackChan := make(chan Stack)
	doneSubChan := make(chan error)

	if err := c.CreateStack(name, templateBody, nil, true, doneSubChan, stackChan); err != nil {
		return err
	}

	go func() {
		defer close(done)
		defer close(stackChan)

		if err := <-doneSubChan; err != nil {
			done <- err
			return
		}
		done <- nil

		s := <-stackChan
		clusterRoleARN := GetOutput(&s, "RoleArn")
		if clusterRoleARN == nil {
			done <- fmt.Errorf("RoleArn is nil")
			return
		}
		c.cfg.clusterRoleARN = *clusterRoleARN

		logger.Debug("clusterConfig = %#v", c.cfg)
		logger.Success("created ServiceRole stack %q", name)
	}()
	return nil
}

func (c *CloudFormation) DeleteStackServiceRole() error {
	s, err := c.GetStackServiceRole()
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
		"NodeGroupName":                    c.cfg.ClusterName + "-DefaultNodeGroup",
		"KeyName":                          c.cfg.keyName,
		"NodeImageId":                      c.cfg.NodeAMI,
		"NodeInstanceType":                 c.cfg.NodeType,
		"NodeAutoScalingGroupMinSize":      string(c.cfg.MinNodes),
		"NodeAutoScalingGroupMaxSize":      string(c.cfg.MaxNodes),
		"ClusterControlPlaneSecurityGroup": c.cfg.securityGroup,
		"Subnets":                          c.cfg.subnetsList,
		"VpcId":                            c.cfg.clusterVPC,
	}
}

func (c *CloudFormation) CreateStackDefaultNodeGroup() error {
	_, err := amazonEksNodegroupYamlBytes()
	if err != nil {
		return errors.Wrap(err, "decompressing bundled template for DefaultNodeGroup stack")
	}
	return nil
}

func (c *CloudFormation) GetStack(name string) (*Stack, error) {
	return c.describeStack(&name)
}

func (c *CloudFormation) GetStackVPC() (*Stack, error) {
	return c.GetStack("EKS-" + c.cfg.ClusterName + "-VPC")
}

func (c *CloudFormation) GetStackServiceRole() (*Stack, error) {
	return c.GetStack("EKS-" + c.cfg.ClusterName + "-ServiceRole")
}

func (c *CloudFormation) GetStackDefaultNodeGroup() (*Stack, error) {
	return c.GetStack("EKS-" + c.cfg.ClusterName + "-DefaultNodeGroup")
}

func GetOutput(stack *Stack, key string) *string {
	for _, x := range stack.Outputs {
		if *x.OutputKey == key {
			return x.OutputValue
		}
	}
	return nil
}
