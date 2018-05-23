package cfn

import (
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

type CloudFormation struct {
	svc *cloudformation.CloudFormation
}

func New() *CloudFormation {
	return &CloudFormation{
		svc: cloudformation.New(session.Must(session.NewSession())),
	}
}

func (c *CloudFormation) CreateStack(name string, templateBody []byte, parameters map[string]string, withIAM bool, done chan struct{}, fail chan cloudformation.Stack) error {
	input := &cloudformation.CreateStackInput{}
	input.SetStackName(name)
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
		return err
	}

	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
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
					close(done)
					return
				case cloudformation.StackStatusCreateFailed:
					fail <- *s
					close(done)
					return
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
					fail <- *s
					close(done)
					return

				}
			}
		}
	}()

	return nil

}

func (c *CloudFormation) describeStack(name *string) (*cloudformation.Stack, error) {
	input := &cloudformation.DescribeStacksInput{
		StackName: name,
	}
	resp, err := c.svc.DescribeStacks(input)
	if err != nil {
		return nil, err
	}
	return resp.Stacks[0], nil
}
func (c *CloudFormation) ListReadyStacks(nameRegex string) ([]*cloudformation.Stack, error) {
	var (
		subErr error
		stack  *cloudformation.Stack
	)

	re := regexp.MustCompile(nameRegex)
	input := &cloudformation.ListStacksInput{
		StackStatusFilter: aws.StringSlice([]string{cloudformation.StackStatusCreateComplete}),
	}
	stacks := []*cloudformation.Stack{}

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

func StackParamsDefaultNodeGroup(clusterName, keyName, nodeAMI, nodeType, minNodes, maxNodes, securityGroup, subnetsList, clusterVPC string) map[string]string {
	return map[string]string{
		"clusterName":                      clusterName,
		"NodeGroupName":                    clusterName + "-DefaultNodeGroup",
		"keyName":                          keyName,
		"NodeImageId":                      nodeAMI,
		"NodeInstanceType":                 nodeType,
		"NodeAutoScalingGroupMinSize":      minNodes,
		"NodeAutoScalingGroupMaxSize":      maxNodes,
		"ClusterControlPlaneSecurityGroup": securityGroup,
		"Subnets":                          subnetsList,
		"VpcId":                            clusterVPC,
	}
}

func (c *CloudFormation) GetStack(name string) (*cloudformation.Stack, error) {
	return c.describeStack(&name)
}

func (c *CloudFormation) GetStackVPC(clusterName string) (*cloudformation.Stack, error) {
	return c.GetStack(strings.Join([]string{"^EKS", clusterName, "VPC$"}, "-"))
}

func (c *CloudFormation) GetStackServiceRole(clusterName string) (*cloudformation.Stack, error) {
	return c.GetStack(
		strings.Join([]string{"^EKS", clusterName, "ServiceRole$"}, "-"),
	)
}

func (c *CloudFormation) GetStackDefaultNodeGroup(clusterName string) (*cloudformation.Stack, error) {
	return c.GetStack(
		strings.Join([]string{"^EKS", clusterName, "DefaultNodeGroup$"}, "-"),
	)
}

func GetOutput(stack *cloudformation.Stack, key string) *string {
	for _, x := range stack.Outputs {
		if *x.OutputKey == key {
			return x.OutputValue
		}
	}
	return nil
}
