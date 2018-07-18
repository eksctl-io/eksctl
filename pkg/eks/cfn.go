package eks

import (
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"

	"github.com/kubicorn/kubicorn/pkg/logger"
)

//go:generate go-bindata -pkg $GOPACKAGE -prefix assets/1.10.3/2018-07-18 -o cfn_templates.go assets/1.10.3/2018-07-18

type Stack = cloudformation.Stack

const (
	policyAmazonEKSWorkerNodePolicy           = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
	policyAmazonEKS_CNI_Policy                = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
	policyAmazonEC2ContainerRegistryPowerUser = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryPowerUser"
	policyAmazonEC2ContainerRegistryReadOnly  = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
)

var (
	defaultPolicyARNs = []string{
		policyAmazonEKSWorkerNodePolicy,
		policyAmazonEKS_CNI_Policy,
	}
)

func (c *ClusterProvider) CreateStack(name string, templateBody []byte, parameters map[string]string, withIAM bool, stack chan Stack, errs chan error) error {
	input := &cloudformation.CreateStackInput{}
	input.SetStackName(name)
	input.SetTags([]*cloudformation.Tag{
		&cloudformation.Tag{
			Key:   aws.String(ClusterNameTag),
			Value: aws.String(c.Spec.ClusterName),
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
	s, err := c.Provider.CloudFormation().CreateStack(input)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("creating CloudFormation stack %q", name))
	}
	logger.Debug("stack = %#v", s)

	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		timer := time.NewTimer(c.Spec.WaitTimeout)
		defer timer.Stop()

		defer close(errs)
		for {
			select {
			case <-timer.C:
				errs <- fmt.Errorf("timed out creating CloudFormation stack %q after %d", name, c.Spec.WaitTimeout)
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
					errs <- fmt.Errorf("unexpected status %q while creating CloudFormation stack %q", *s.StackStatus, name)
					// stack <- *s // this usually results in closed channel panic, but we don't need it really
					logger.Debug("stack = %#v", s)
					return
				}
			}
		}
	}()

	return nil
}

func (c *ClusterProvider) describeStack(name *string) (*Stack, error) {
	input := &cloudformation.DescribeStacksInput{
		StackName: name,
	}
	resp, err := c.Provider.CloudFormation().DescribeStacks(input)
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
	if err := c.Provider.CloudFormation().ListStacksPages(input, pager); err != nil {
		return nil, err
	}
	if subErr != nil {
		return nil, subErr
	}
	return stacks, nil
}

func (c *ClusterProvider) DeleteStack(name string) error {
	s, err := c.describeStack(&name)
	if err != nil {
		return errors.Wrapf(err, "not able to get stack %q for deletion", name)
	}

	for _, tag := range s.Tags {
		if *tag.Key == ClusterNameTag && *tag.Value == c.Spec.ClusterName {
			input := &cloudformation.DeleteStackInput{
				StackName: s.StackName,
			}

			if _, err := c.Provider.CloudFormation().DeleteStack(input); err != nil {
				return errors.Wrapf(err, "not able to delete stack %q", name)
			}
			return nil
		}
	}

	return fmt.Errorf("cannot delete stack %s as it doesn't bare our %q tag", *s.StackName,
		fmt.Sprintf("%s:%s", ClusterNameTag, c.Spec.ClusterName))
}

func (c *ClusterProvider) stackNameVPC() string {
	return "EKS-" + c.Spec.ClusterName + "-VPC"
}

func (c *ClusterProvider) stackParamsVPC() map[string]string {
	params := map[string]string{
		"AvailabilityZones": strings.Join(c.Status.availabilityZones, ","),
	}
	return params
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

		securityGroup := GetOutput(&s, "SecurityGroups")
		if securityGroup == nil {
			errs <- fmt.Errorf("SecurityGroups is nil")
			return
		}
		c.Spec.securityGroup = *securityGroup

		subnetsList := GetOutput(&s, "SubnetIds")
		if subnetsList == nil {
			errs <- fmt.Errorf("SubnetIds is nil")
			return
		}
		c.Spec.subnetsList = *subnetsList

		clusterVPC := GetOutput(&s, "VpcId")
		if clusterVPC == nil {
			errs <- fmt.Errorf("VpcId is nil")
			return
		}
		c.Spec.clusterVPC = *clusterVPC

		logger.Debug("clusterConfig = %#v", c.Spec)
		logger.Success("created VPC stack %q", name)

		errs <- nil
	}()
	return nil
}

func (c *ClusterProvider) DeleteStackVPC() error {
	return c.DeleteStack(c.stackNameVPC())
}

func (c *ClusterProvider) stackNameServiceRole() string {
	return "EKS-" + c.Spec.ClusterName + "-ServiceRole"
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
		c.Spec.clusterRoleARN = *clusterRoleARN

		logger.Debug("clusterConfig = %#v", c.Spec)
		logger.Success("created ServiceRole stack %q", name)

		errs <- nil
	}()
	return nil
}

func (c *ClusterProvider) DeleteStackServiceRole() error {
	return c.DeleteStack(c.stackNameServiceRole())
}

func (c *ClusterProvider) stackNameControlPlane() string {
	return "EKS-" + c.cfg.ClusterName + "-ControlPlane"
}

func (c *ClusterProvider) createStackControlPlane(errs chan error) error {
	stackName := c.stackNameControlPlane()
	stackChan := make(chan Stack)
	taskErrs := make(chan error)

	params := make(map[string]string)
	params["ClusterName"] = c.cfg.ClusterName
	params["Subnets"] = c.cfg.subnetsList
	params["ControlPlaneSecurityGroups"] = c.cfg.securityGroup
	params["KubernetesVersion"] = "1.10"
	params["ServiceRoleARN"] = c.cfg.clusterRoleARN

	templateBody, err := amazonEksClusterYamlBytes()
	if err != nil {
		return errors.Wrap(err, "decompressing bundled template for Control Plane stack")
	}

	if err := c.CreateStack(stackName, templateBody, params, true, stackChan, taskErrs); err != nil {
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

		logger.Debug("created ControlPlane stack %q – processing outputs", stackName)

		clusterCA := GetOutput(&s, "EKSClusterCA")
		if clusterCA == nil {
			errs <- fmt.Errorf("cluster CA is nil")
			return
		}

		data, err := base64.StdEncoding.DecodeString(*clusterCA)
		if err != nil {
			errs <- errors.Wrap(err, "decoding certificate authority data")
			return
		}

		c.cfg.CertificateAuthorityData = data

		clusterEndpoint := GetOutput(&s, "EKSClusterEndpoint")
		if clusterEndpoint == nil {
			errs <- fmt.Errorf("cluster endpoint is nil")
			return
		}
		c.cfg.MasterEndpoint = *clusterEndpoint

		clusterARN := GetOutput(&s, "EKSClusterARN")
		if clusterARN == nil {
			errs <- fmt.Errorf("cluster ARN is nil")
			return
		}
		c.cfg.ClusterARN = *clusterARN

		logger.Debug("clusterConfig = %#v", c.cfg)
		logger.Success("created Control Plane stack %q", stackName)

		errs <- nil
	}()
	return nil

}

func (c *ClusterProvider) stackNameDefaultNodeGroup() string {
	return "EKS-" + c.Spec.ClusterName + "-DefaultNodeGroup"
}

func (c *ClusterProvider) stackParamsDefaultNodeGroup() map[string]string {
	regionalAMIs := map[string]string{
		// TODO: https://github.com/weaveworks/eksctl/issues/49
		// currently source of truth for these is here:
		// https://docs.aws.amazon.com/eks/latest/userguide/launch-workers.html
		"us-west-2": "ami-73a6e20b",
		"us-east-1": "ami-dea4d5a1",
	}

	if c.Spec.NodeAMI == "" {
		c.Spec.NodeAMI = regionalAMIs[c.Spec.Region]
	}

	if c.Spec.MinNodes == 0 && c.Spec.MaxNodes == 0 {
		c.Spec.MinNodes = c.Spec.Nodes
		c.Spec.MaxNodes = c.Spec.Nodes
	}

	if len(c.Spec.PolicyARNs) == 0 {
		c.Spec.PolicyARNs = defaultPolicyARNs
	}
	if c.Spec.Addons.WithIAM.PolicyAmazonEC2ContainerRegistryPowerUser {
		c.Spec.PolicyARNs = append(c.Spec.PolicyARNs, policyAmazonEC2ContainerRegistryPowerUser)
	} else {
		c.Spec.PolicyARNs = append(c.Spec.PolicyARNs, policyAmazonEC2ContainerRegistryReadOnly)
	}

	params := map[string]string{
		"ClusterName":                      c.Spec.ClusterName,
		"NodeGroupName":                    "default",
		"KeyName":                          c.Spec.keyName,
		"NodeImageId":                      c.Spec.NodeAMI,
		"NodeInstanceType":                 c.Spec.NodeType,
		"NodeAutoScalingGroupMinSize":      fmt.Sprintf("%d", c.Spec.MinNodes),
		"NodeAutoScalingGroupMaxSize":      fmt.Sprintf("%d", c.Spec.MaxNodes),
		"ClusterControlPlaneSecurityGroup": c.Spec.securityGroup,
		"Subnets":                          c.Spec.subnetsList,
		"VpcId":                            c.Spec.clusterVPC,
		"ManagedPolicyArns":                strings.Join(c.Spec.PolicyARNs, ","),
	}

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
			errs <- fmt.Errorf("NodeInstanceRole is nil")
			return
		}
		c.Spec.nodeInstanceRoleARN = *nodeInstanceRoleARN

		logger.Debug("clusterConfig = %#v", c.Spec)
		logger.Success("created DefaultNodeGroup stack %q", name)

		errs <- nil
	}()
	return nil
}

func (c *ClusterProvider) DeleteStackDefaultNodeGroup() error {
	return c.DeleteStack(c.stackNameDefaultNodeGroup())
}

func GetOutput(stack *Stack, key string) *string {
	for _, x := range stack.Outputs {
		if *x.OutputKey == key {
			return x.OutputValue
		}
	}
	return nil
}
