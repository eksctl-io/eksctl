package eks

import (
	"fmt"
	"regexp"
	"time"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/weaveworks/eksctl/pkg/utils"
)

//go:generate go-bindata -pkg $GOPACKAGE -prefix assets/1.10.3/2018-06-05 -o cfn_templates.go assets/1.10.3/2018-06-05

type Stack = cloudformation.Stack

func (c *ClusterProvider) CreateStack(name string, templateBody []byte, parameters map[string]string, withIAM bool, stack chan Stack, errs chan error) error {
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

	logger.Debug("input = %#v", input)
	s, err := c.svc.cfn.CreateStack(input)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("creating CloudFormation stack %q", name))
	}
	logger.Debug("stack = %#v", s)

	go c.waitForStack(name, stack, errs)

	return nil
}

func (c *ClusterProvider) waitForStack(name string, stack chan Stack, errs chan error) {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	timer := time.NewTimer(time.Duration(c.cfg.AWSOperationTimeoutSeconds) * time.Second)
	defer timer.Stop()

	defer close(errs)
	for {
		select {
		case <-timer.C:
			errs <- fmt.Errorf("creating CloudFormation stack %q timed out after %d seconds", name, c.cfg.AWSOperationTimeoutSeconds)
			logger.Debug("stack = %q", name)
			return

		case <-ticker.C:
			s, err := c.describeStack(&name)
			if err != nil {
				logger.Warning("continue despite err=%q", err.Error())
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
				fallthrough // TODO: https://github.com/weaveworks/eksctl/issues/24
			default:
				errs <- fmt.Errorf("creating CloudFormation stack %q: %s", name, *s.StackStatus)
				// stack <- *s // this usually results in closed channel panic, but we don't need it really
				logger.Debug("stack = %#v", s)
				return
			}
		}
	}
}

func (c *ClusterProvider) describeStack(name *string) (*Stack, error) {
	input := &cloudformation.DescribeStacksInput{
		StackName: name,
	}
	resp, err := c.svc.cfn.DescribeStacks(input)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("describing CloudFormation stack %q", *name))
	}
	return resp.Stacks[0], nil
}

func (c *ClusterProvider) ListReadyStacks(nameRegex string) ([]*Stack, error) {
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
	if err := c.svc.cfn.ListStacksPages(input, pager); err != nil {
		return nil, err
	}
	if subErr != nil {
		return nil, subErr
	}
	return stacks, nil
}

func (c *ClusterProvider) stackNameVPC() string {
	return "EKS-" + c.cfg.ClusterName + "-VPC"
}

func (c *ClusterProvider) createStackVPC(errs chan error) error {
	name := c.stackNameVPC()
	logger.Info("creating VPC stack %q", name)
	templateBody, err := amazonEksVpcSampleYamlBytes()
	if err != nil {
		return errors.Wrap(err, "decompressing bundled template for VPC stack")
	}

	stackChan := make(chan Stack)
	taskErrs := make(chan error)

	if err := c.CreateStack(name, templateBody, nil, false, stackChan, taskErrs); err != nil {
		if utils.HasAwsErrorCode(err, cloudformation.ErrCodeAlreadyExistsException) {
			logger.Info("using existing VPC stack %q", name)
			go c.waitForStack(name, stackChan, taskErrs)
		} else {
			return err
		}
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

		securityGroup := GetOutput(&s, "SecurityGroups")
		if securityGroup == nil {
			errs <- fmt.Errorf("SecurityGroups is nil")
			return
		}
		c.cfg.securityGroup = *securityGroup

		subnetsList := GetOutput(&s, "SubnetIds")
		if subnetsList == nil {
			errs <- fmt.Errorf("SubnetIds is nil")
			return
		}
		c.cfg.subnetsList = *subnetsList

		clusterVPC := GetOutput(&s, "VpcId")
		if clusterVPC == nil {
			errs <- fmt.Errorf("VpcId is nil")
			return
		}
		c.cfg.clusterVPC = *clusterVPC

		logger.Debug("clusterConfig = %#v", c.cfg)
		logger.Success("created VPC stack %q", name)

		errs <- nil
	}()
	return nil
}

func (c *ClusterProvider) DeleteStackVPC() error {
	name := c.stackNameVPC()
	s, err := c.describeStack(&name)
	if err != nil {
		return errors.Wrap(err, "not able to get VPC stack for deletion")
	}

	input := &cloudformation.DeleteStackInput{
		StackName: s.StackName,
	}

	if _, err := c.svc.cfn.DeleteStack(input); err != nil {
		return errors.Wrap(err, "not able to delete VPC stack")
	}
	return nil
}

func (c *ClusterProvider) stackNameServiceRole() string {
	return "EKS-" + c.cfg.ClusterName + "-ServiceRole"
}

func (c *ClusterProvider) createStackServiceRole(errs chan error) error {
	name := c.stackNameServiceRole()
	logger.Info("creating ServiceRole stack %q", name)
	templateBody, err := amazonEksServiceRoleYamlBytes()
	if err != nil {
		return errors.Wrap(err, "decompressing bundled template for ServiceRole stack")
	}

	stackChan := make(chan Stack)
	taskErrs := make(chan error)

	if err := c.CreateStack(name, templateBody, nil, true, stackChan, taskErrs); err != nil {
		if utils.HasAwsErrorCode(err, cloudformation.ErrCodeAlreadyExistsException) {
			logger.Info("using existing ServiceRole stack %q", name)
			go c.waitForStack(name, stackChan, taskErrs)
		} else {
			return err
		}
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

func (c *ClusterProvider) DeleteStackServiceRole() error {
	name := c.stackNameServiceRole()
	s, err := c.describeStack(&name)
	if err != nil {
		return errors.Wrap(err, "not able to get ServiceRole stack for deletion")
	}

	input := &cloudformation.DeleteStackInput{
		StackName: s.StackName,
	}

	if _, err := c.svc.cfn.DeleteStack(input); err != nil {
		return errors.Wrap(err, "not able to delete ServiceRole stack")
	}
	return nil
}

func (c *ClusterProvider) stackNameDefaultNodeGroup() string {
	return "EKS-" + c.cfg.ClusterName + "-DefaultNodeGroup"
}

func (c *ClusterProvider) stackParamsDefaultNodeGroup() map[string]string {
	regionalAMIs := map[string]string{
		// TODO: https://github.com/weaveworks/eksctl/issues/49
		// currently source of truth for these is here:
		// https://docs.aws.amazon.com/eks/latest/userguide/launch-workers.html
		"us-west-2": "ami-73a6e20b",
		"us-east-1": "ami-dea4d5a1",
	}

	if c.cfg.NodeAMI == "" {
		c.cfg.NodeAMI = regionalAMIs[c.cfg.Region]
	}

	if c.cfg.MinNodes == 0 && c.cfg.MaxNodes == 0 {
		c.cfg.MinNodes = c.cfg.Nodes
		c.cfg.MaxNodes = c.cfg.Nodes
	}

	params := map[string]string{
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

	logger.Info("Node stack params: %#v", params)
	return params
}

func (c *ClusterProvider) createStackDefaultNodeGroup(errs chan error) error {
	name := c.stackNameDefaultNodeGroup()
	logger.Info("creating DefaultNodeGroup stack %q", name)
	templateBody, err := amazonEksNodegroupYamlBytes()
	if err != nil {
		return errors.Wrap(err, "decompressing bundled template for DefaultNodeGroup stack")
	}

	stackChan := make(chan Stack)
	taskErrs := make(chan error)

	if err := c.CreateStack(name, templateBody, c.stackParamsDefaultNodeGroup(), true, stackChan, taskErrs); err != nil {
		if utils.HasAwsErrorCode(err, cloudformation.ErrCodeAlreadyExistsException) {
			logger.Info("using existing DefaultNodeGroup stack %q", name)
			go c.waitForStack(name, stackChan, taskErrs)
		} else {
			return err
		}
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

func (c *ClusterProvider) DeleteStackDefaultNodeGroup() error {
	name := c.stackNameDefaultNodeGroup()
	s, err := c.describeStack(&name)
	if err != nil {
		return errors.Wrap(err, "not able to get DefaultNodeGroup stack for deletion")
	}

	input := &cloudformation.DeleteStackInput{
		StackName: s.StackName,
	}

	if _, err := c.svc.cfn.DeleteStack(input); err != nil {
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
