package manager

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/autoscaling/autoscalingiface"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudformation/cloudformationiface"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/aws/aws-sdk-go/service/cloudtrail/cloudtrailiface"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/aws/aws-sdk-go/service/iam/iamiface"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/waiter"
	"github.com/weaveworks/eksctl/pkg/version"
)

const (
	resourcesRootPath = "Resources"
	outputsRootPath   = "Outputs"
	mappingsRootPath  = "Mappings"
	ourStackRegexFmt  = "^(eksctl|EKS)-%s-((cluster|nodegroup-.+|addon-.+|fargate|karpenter)|(VPC|ServiceRole|ControlPlane|DefaultNodeGroup))$"
	clusterStackRegex = "eksctl-.*-cluster"
)

var (
	stackCapabilitiesIAM      = aws.StringSlice([]string{cloudformation.CapabilityCapabilityIam})
	stackCapabilitiesNamedIAM = aws.StringSlice([]string{cloudformation.CapabilityCapabilityNamedIam})
)

// Stack represents the CloudFormation stack
type Stack = cloudformation.Stack

// StackInfo hold the stack along with template and resources
type StackInfo struct {
	Stack     *Stack
	Resources []*cloudformation.StackResource
	Template  *string
}

// TemplateData is a union (sum type) to describe template data.
type TemplateData interface {
	isTemplateData()
}

// TemplateBody allows to pass the full template.
type TemplateBody []byte

func (b TemplateBody) isTemplateData() {}

// TemplateURL allows to pass in a link to a template.
type TemplateURL string

func (u TemplateURL) isTemplateData() {}

// ChangeSet represents a CloudFormation ChangeSet
type ChangeSet = cloudformation.DescribeChangeSetOutput

// StackCollection stores the CloudFormation stack information
type StackCollection struct {
	cloudformationAPI cloudformationiface.CloudFormationAPI
	ec2API            ec2iface.EC2API
	eksAPI            eksiface.EKSAPI
	iamAPI            iamiface.IAMAPI
	cloudTrailAPI     cloudtrailiface.CloudTrailAPI
	asgAPI            autoscalingiface.AutoScalingAPI
	spec              *api.ClusterConfig
	disableRollback   bool
	roleARN           string
	region            string
	waitTimeout       time.Duration
	sharedTags        []*cloudformation.Tag
}

func newTag(key, value string) *cloudformation.Tag {
	return &cloudformation.Tag{Key: &key, Value: &value}
}

// NewStackCollection creates a stack manager for a single cluster
func NewStackCollection(provider api.ClusterProvider, spec *api.ClusterConfig) *StackCollection {
	tags := []*cloudformation.Tag{
		newTag(api.ClusterNameTag, spec.Metadata.Name),
		newTag(api.OldClusterNameTag, spec.Metadata.Name),
		newTag(api.EksctlVersionTag, version.GetVersion()),
	}
	for key, value := range spec.Metadata.Tags {
		tags = append(tags, newTag(key, value))
	}
	return &StackCollection{
		spec:              spec,
		sharedTags:        tags,
		cloudformationAPI: provider.CloudFormation(),
		ec2API:            provider.EC2(),
		eksAPI:            provider.EKS(),
		iamAPI:            provider.IAM(),
		cloudTrailAPI:     provider.CloudTrail(),
		asgAPI:            provider.ASG(),
		disableRollback:   provider.CloudFormationDisableRollback(),
		roleARN:           provider.CloudFormationRoleARN(),
		region:            provider.Region(),
		waitTimeout:       provider.WaitTimeout(),
	}
}

// DoCreateStackRequest requests the creation of a CloudFormation stack
func (c *StackCollection) DoCreateStackRequest(i *Stack, templateData TemplateData, tags, parameters map[string]string, withIAM bool, withNamedIAM bool) error {
	input := &cloudformation.CreateStackInput{
		StackName:       i.StackName,
		DisableRollback: aws.Bool(c.disableRollback),
	}
	input.Tags = append(input.Tags, c.sharedTags...)
	for k, v := range tags {
		input.Tags = append(input.Tags, newTag(k, v))
	}

	switch data := templateData.(type) {
	case TemplateBody:
		input.SetTemplateBody(string(data))
	case TemplateURL:
		input.SetTemplateURL(string(data))
	default:
		return fmt.Errorf("unknown template data type: %T", templateData)
	}

	if withIAM {
		input.SetCapabilities(stackCapabilitiesIAM)
	}

	if withNamedIAM {
		input.SetCapabilities(stackCapabilitiesNamedIAM)
	}

	if cfnRole := c.roleARN; cfnRole != "" {
		input = input.SetRoleARN(cfnRole)
	}

	for k, v := range parameters {
		p := &cloudformation.Parameter{
			ParameterKey:   aws.String(k),
			ParameterValue: aws.String(v),
		}
		input.Parameters = append(input.Parameters, p)
	}

	logger.Debug("CreateStackInput = %#v", input)
	s, err := c.cloudformationAPI.CreateStack(input)
	if err != nil {
		return errors.Wrapf(err, "creating CloudFormation stack %q", *i.StackName)
	}
	i.StackId = s.StackId
	return nil
}

// CreateStack with given name, stack builder instance and parameters;
// any errors will be written to errs channel, when nil is written,
// assume completion, do not expect more then one error value on the
// channel, it's closed immediately after it is written to
func (c *StackCollection) CreateStack(stackName string, resourceSet builder.ResourceSet, tags, parameters map[string]string, errs chan error) error {
	stack, err := c.createStackRequest(stackName, resourceSet, tags, parameters)
	if err != nil {
		return err
	}

	go c.waitUntilStackIsCreated(stack, resourceSet, errs)
	return nil
}

// createClusterStack creates the cluster stack
func (c *StackCollection) createClusterStack(stackName string, resourceSet builder.ResourceSet, errCh chan error) error {
	// Unlike with `createNodeGroupTask`, all tags are already set for the cluster stack
	stack, err := c.createStackRequest(stackName, resourceSet, nil, nil)
	if err != nil {
		return err
	}
	go func() {
		defer close(errCh)
		troubleshoot := func() {
			stack, err := c.DescribeStack(stack)
			if err != nil {
				logger.Info("error describing stack to troubleshoot the cause of the failure; "+
					"check the CloudFormation console for further details", err)
				return
			}

			logger.Critical("unexpected status %q while waiting for CloudFormation stack %q", *stack.StackStatus, *stack.StackName)
			c.troubleshootStackFailureCause(stack, cloudformation.StackStatusCreateComplete)
		}

		ctx, cancelFunc := context.WithTimeout(context.Background(), c.waitTimeout)
		defer cancelFunc()

		stack, err := waiter.WaitForStack(ctx, c.cloudformationAPI, *stack.StackId, *stack.StackName, func(attempts int) time.Duration {
			// Wait 30s for the first two requests, and 1m for subsequent requests.
			if attempts <= 2 {
				return 30 * time.Second
			}
			return 1 * time.Minute
		})

		if err != nil {
			troubleshoot()
			errCh <- err
			return
		}

		if err := resourceSet.GetAllOutputs(*stack); err != nil {
			errCh <- errors.Wrapf(err, "getting stack %q outputs", *stack.StackName)
			return
		}

		errCh <- nil
	}()

	return nil
}

func (c *StackCollection) createStackRequest(stackName string, resourceSet builder.ResourceSet, tags, parameters map[string]string) (*Stack, error) {
	stack := &Stack{StackName: &stackName}
	templateBody, err := resourceSet.RenderJSON()
	if err != nil {
		return nil, errors.Wrapf(err, "rendering template for %q stack", *stack.StackName)
	}

	if err := c.DoCreateStackRequest(stack, TemplateBody(templateBody), tags, parameters, resourceSet.WithIAM(), resourceSet.WithNamedIAM()); err != nil {
		return nil, err
	}

	logger.Info("deploying stack %q", stackName)
	return stack, nil
}

// UpdateStack will update a CloudFormation stack by creating and executing a ChangeSet
func (c *StackCollection) UpdateStack(options UpdateStackOptions) error {
	logger.Info(options.Description)
	i := &Stack{StackName: &options.StackName}
	// Read existing tags
	s, err := c.DescribeStack(i)
	if err != nil {
		return err
	}
	if err := c.doCreateChangeSetRequest(options.StackName, options.ChangeSetName, options.Description, options.TemplateData, options.Parameters, s.Capabilities, s.Tags); err != nil {
		return err
	}
	if err := c.doWaitUntilChangeSetIsCreated(i, options.ChangeSetName); err != nil {
		if _, ok := err.(*noChangeError); ok {
			return nil
		}
		return err
	}
	changeSet, err := c.DescribeStackChangeSet(i, options.ChangeSetName)
	if err != nil {
		return err
	}
	logger.Debug("changes = %#v", changeSet.Changes)
	if err := c.doExecuteChangeSet(options.StackName, options.ChangeSetName); err != nil {
		logger.Warning("error executing Cloudformation changeSet %s in stack %s. Check the Cloudformation console for further details", options.ChangeSetName, options.StackName)
		return err
	}
	if options.Wait {
		return c.doWaitUntilStackIsUpdated(i)
	}
	return nil
}

// DescribeStack describes a cloudformation stack.
func (c *StackCollection) DescribeStack(i *Stack) (*Stack, error) {
	input := &cloudformation.DescribeStacksInput{
		StackName: i.StackName,
	}
	if api.IsSetAndNonEmptyString(i.StackId) {
		input.StackName = i.StackId
	}
	resp, err := c.cloudformationAPI.DescribeStacks(input)
	if err != nil {
		return nil, errors.Wrapf(err, "describing CloudFormation stack %q", *i.StackName)
	}
	return resp.Stacks[0], nil
}

// GetManagedNodeGroupTemplate returns the template for a ManagedNodeGroup resource
func (c *StackCollection) GetManagedNodeGroupTemplate(nodeGroupName string) (string, error) {
	nodeGroupType, err := c.GetNodeGroupStackType(nodeGroupName)
	if err != nil {
		return "", err
	}

	if nodeGroupType != api.NodeGroupTypeManaged {
		return "", fmt.Errorf("%q is not a managed nodegroup", nodeGroupName)
	}

	stackName := c.makeNodeGroupStackName(nodeGroupName)
	templateBody, err := c.GetStackTemplate(stackName)
	if err != nil {
		return "", err
	}

	return templateBody, nil
}

// UpdateNodeGroupStack updates the nodegroup stack with the specified template
func (c *StackCollection) UpdateNodeGroupStack(nodeGroupName, template string, wait bool) error {
	stackName := c.makeNodeGroupStackName(nodeGroupName)
	return c.UpdateStack(UpdateStackOptions{
		StackName:     stackName,
		ChangeSetName: c.MakeChangeSetName("update-nodegroup"),
		Description:   "updating nodegroup stack",
		TemplateData:  TemplateBody(template),
		Wait:          wait,
	})
}

// ListStacksMatching gets all of CloudFormation stacks with names matching nameRegex.
func (c *StackCollection) ListStacksMatching(nameRegex string, statusFilters ...string) ([]*Stack, error) {
	var (
		subErr error
		stack  *Stack
	)

	re, err := regexp.Compile(nameRegex)
	if err != nil {
		return nil, errors.Wrap(err, "cannot list stacks")
	}
	input := &cloudformation.ListStacksInput{
		StackStatusFilter: defaultStackStatusFilter(),
	}
	if len(statusFilters) > 0 {
		input.StackStatusFilter = aws.StringSlice(statusFilters)
	}
	stacks := []*Stack{}

	pager := func(p *cloudformation.ListStacksOutput, _ bool) bool {
		for _, s := range p.StackSummaries {
			if re.MatchString(*s.StackName) {
				stack, subErr = c.DescribeStack(&Stack{StackName: s.StackName, StackId: s.StackId})
				if subErr != nil {
					return false
				}
				stacks = append(stacks, stack)
			}
		}
		return true
	}
	if err := c.cloudformationAPI.ListStacksPages(input, pager); err != nil {
		return nil, err
	}
	if subErr != nil {
		return nil, subErr
	}
	return stacks, nil
}

// ListStackNamesMatching gets all stack names matching regex
func (c *StackCollection) ListClusterStackNames() ([]string, error) {
	stacks := []string{}
	re, err := regexp.Compile(clusterStackRegex)
	if err != nil {
		return nil, errors.Wrap(err, "cannot list stacks")
	}
	input := &cloudformation.ListStacksInput{
		StackStatusFilter: defaultStackStatusFilter(),
	}

	pager := func(p *cloudformation.ListStacksOutput, _ bool) bool {
		for _, s := range p.StackSummaries {
			if re.MatchString(*s.StackName) {
				stacks = append(stacks, *s.StackName)
			}
		}
		return true
	}
	if err := c.cloudformationAPI.ListStacksPages(input, pager); err != nil {
		return nil, err
	}

	return stacks, nil
}

// ListStacks gets all of CloudFormation stacks
func (c *StackCollection) ListStacks(statusFilters ...string) ([]*Stack, error) {
	return c.ListStacksMatching(fmtStacksRegexForCluster(c.spec.Metadata.Name), statusFilters...)
}

// StackStatusIsNotTransitional will return true when stack status is non-transitional
func (*StackCollection) StackStatusIsNotTransitional(s *Stack) bool {
	for _, state := range nonTransitionalReadyStackStatuses() {
		if *s.StackStatus == state {
			return true
		}
	}
	return false
}

func nonTransitionalReadyStackStatuses() []string {
	return []string{
		cloudformation.StackStatusCreateComplete,
		cloudformation.StackStatusUpdateComplete,
		cloudformation.StackStatusRollbackComplete,
		cloudformation.StackStatusUpdateRollbackComplete,
	}
}

// StackStatusIsNotReady will return true when stack statate is non-ready
func (*StackCollection) StackStatusIsNotReady(s *Stack) bool {
	for _, state := range nonReadyStackStatuses() {
		if *s.StackStatus == state {
			return true
		}
	}
	return false
}

func nonReadyStackStatuses() []string {
	return []string{
		cloudformation.StackStatusCreateInProgress,
		cloudformation.StackStatusCreateFailed,
		cloudformation.StackStatusRollbackInProgress,
		cloudformation.StackStatusRollbackFailed,
		cloudformation.StackStatusDeleteInProgress,
		cloudformation.StackStatusDeleteFailed,
		cloudformation.StackStatusUpdateInProgress,
		cloudformation.StackStatusUpdateCompleteCleanupInProgress,
		cloudformation.StackStatusUpdateRollbackInProgress,
		cloudformation.StackStatusUpdateRollbackFailed,
		cloudformation.StackStatusUpdateRollbackCompleteCleanupInProgress,
		cloudformation.StackStatusReviewInProgress,
	}
}

func allNonDeletedStackStatuses() []string {
	return []string{
		cloudformation.StackStatusCreateInProgress,
		cloudformation.StackStatusCreateFailed,
		cloudformation.StackStatusCreateComplete,
		cloudformation.StackStatusRollbackInProgress,
		cloudformation.StackStatusRollbackFailed,
		cloudformation.StackStatusRollbackComplete,
		cloudformation.StackStatusDeleteInProgress,
		cloudformation.StackStatusDeleteFailed,
		cloudformation.StackStatusUpdateInProgress,
		cloudformation.StackStatusUpdateCompleteCleanupInProgress,
		cloudformation.StackStatusUpdateComplete,
		cloudformation.StackStatusUpdateRollbackInProgress,
		cloudformation.StackStatusUpdateRollbackFailed,
		cloudformation.StackStatusUpdateRollbackCompleteCleanupInProgress,
		cloudformation.StackStatusUpdateRollbackComplete,
		cloudformation.StackStatusReviewInProgress,
	}
}

func defaultStackStatusFilter() []*string {
	return aws.StringSlice(allNonDeletedStackStatuses())
}

// DeleteStackByName sends a request to delete the stack
func (c *StackCollection) DeleteStackByName(name string) (*Stack, error) {
	s, err := c.DescribeStack(&Stack{StackName: &name})
	if err != nil {
		err = errors.Wrapf(err, "not able to get stack %q for deletion", name)
		stacks, newErr := c.ListStacksMatching(fmt.Sprintf("^%s$", name), cloudformation.StackStatusDeleteComplete)
		if newErr != nil {
			logger.Critical("not able double-check if stack was already deleted: %s", newErr.Error())
		}
		if count := len(stacks); count > 0 {
			logger.Debug("%d deleted stacks found {%v}", count, stacks)
			logger.Info("stack %q was already deleted", name)
			return nil, nil
		}
		return nil, err
	}
	return c.DeleteStackBySpec(s)
}

// DeleteStackByNameSync sends a request to delete the stack, and waits until status is DELETE_COMPLETE;
// any errors will be written to errs channel, assume completion when nil is written, do not expect
// more then one error value on the channel, it's closed immediately after it is written to
func (c *StackCollection) DeleteStackByNameSync(name string) error {
	stack, err := c.DeleteStackByName(name)
	if err != nil {
		return err
	}

	logger.Info("waiting for stack %q to get deleted", *stack.StackName)

	return c.doWaitUntilStackIsDeleted(stack)
}

// DeleteStackBySpec sends a request to delete the stack
func (c *StackCollection) DeleteStackBySpec(s *Stack) (*Stack, error) {
	for _, tag := range s.Tags {
		if matchesClusterName(*tag.Key, *tag.Value, c.spec.Metadata.Name) {
			input := &cloudformation.DeleteStackInput{
				StackName: s.StackId,
			}

			if cfnRole := c.roleARN; cfnRole != "" {
				input = input.SetRoleARN(cfnRole)
			}

			if _, err := c.cloudformationAPI.DeleteStack(input); err != nil {
				return nil, errors.Wrapf(err, "not able to delete stack %q", *s.StackName)
			}
			logger.Info("will delete stack %q", *s.StackName)
			return s, nil
		}
	}

	return nil, fmt.Errorf("cannot delete stack %q as it doesn't bear our %q, %q tags", *s.StackName,
		fmt.Sprintf("%s:%s", api.OldClusterNameTag, c.spec.Metadata.Name),
		fmt.Sprintf("%s:%s", api.ClusterNameTag, c.spec.Metadata.Name))
}

func matchesClusterName(key, value, name string) bool {
	if key == api.ClusterNameTag && value == name {
		return true
	}

	if key == api.OldClusterNameTag && value == name {
		return true
	}
	return false
}

// DeleteStackBySpecSync sends a request to delete the stack, and waits until status is DELETE_COMPLETE;
// any errors will be written to errs channel, assume completion when nil is written, do not expect
// more then one error value on the channel, it's closed immediately after it is written to
func (c *StackCollection) DeleteStackBySpecSync(s *Stack, errs chan error) error {
	i, err := c.DeleteStackBySpec(s)
	if err != nil {
		return err
	}

	logger.Info("waiting for stack %q to get deleted", *i.StackName)

	go c.waitUntilStackIsDeleted(i, errs)

	return nil
}

func fmtStacksRegexForCluster(name string) string {
	return fmt.Sprintf(ourStackRegexFmt, name)
}

// DescribeStacks describes the existing stacks
func (c *StackCollection) DescribeStacks() ([]*Stack, error) {
	stacks, err := c.ListStacks()
	if err != nil {
		return nil, errors.Wrapf(err, "describing CloudFormation stacks for %q", c.spec.Metadata.Name)
	}
	if len(stacks) == 0 {
		logger.Debug("No stacks found for %s", c.spec.Metadata.Name)
	}
	return stacks, nil
}

func (c *StackCollection) GetClusterStackIfExists() (*Stack, error) {
	clusterStackNames, err := c.ListClusterStackNames()
	if err != nil {
		return nil, err
	}
	return c.getClusterStackUsingCachedList(clusterStackNames)
}

func (c *StackCollection) HasClusterStackUsingCachedList(clusterStackNames []string) (bool, error) {
	stack, err := c.getClusterStackUsingCachedList(clusterStackNames)
	return stack != nil, err
}

func (c *StackCollection) getClusterStackUsingCachedList(clusterStackNames []string) (*Stack, error) {
	clusterStackName := c.MakeClusterStackName()
	for _, stack := range clusterStackNames {
		if stack == clusterStackName {
			stack, err := c.DescribeStack(&cloudformation.Stack{StackName: &clusterStackName})
			if err != nil {
				return nil, err
			}
			for _, tag := range stack.Tags {
				if matchesClusterName(*tag.Key, *tag.Value, c.spec.Metadata.Name) {
					return stack, nil
				}
			}
		}
	}
	return nil, nil
}

// DescribeStackEvents describes the events that have occurred on the stack
func (c *StackCollection) DescribeStackEvents(i *Stack) ([]*cloudformation.StackEvent, error) {
	input := &cloudformation.DescribeStackEventsInput{
		StackName: i.StackName,
	}
	if api.IsSetAndNonEmptyString(i.StackId) {
		input.StackName = i.StackId
	}

	events := []*cloudformation.StackEvent{}

	pager := func(p *cloudformation.DescribeStackEventsOutput, _ bool) bool {
		events = append(events, p.StackEvents...)
		return true
	}
	if err := c.cloudformationAPI.DescribeStackEventsPages(input, pager); err != nil {
		return nil, errors.Wrapf(err, "describing CloudFormation stack %q events", *i.StackName)
	}

	return events, nil
}

// LookupCloudTrailEvents looks up stack events in CloudTrail
func (c *StackCollection) LookupCloudTrailEvents(i *Stack) ([]*cloudtrail.Event, error) {
	input := &cloudtrail.LookupEventsInput{
		LookupAttributes: []*cloudtrail.LookupAttribute{{
			AttributeKey:   aws.String(cloudtrail.LookupAttributeKeyResourceName),
			AttributeValue: i.StackId,
		}},
	}

	events := []*cloudtrail.Event{}

	pager := func(p *cloudtrail.LookupEventsOutput, _ bool) bool {
		events = append(events, p.Events...)
		return true
	}
	if err := c.cloudTrailAPI.LookupEventsPages(input, pager); err != nil {
		return nil, errors.Wrapf(err, "looking up CloudTrail events for stack %q", *i.StackName)
	}

	return events, nil
}

func (c *StackCollection) doCreateChangeSetRequest(stackName, changeSetName, description string, templateData TemplateData,
	parameters map[string]string, capabilities []*string, tags []*cloudformation.Tag) error {
	input := &cloudformation.CreateChangeSetInput{
		StackName:     &stackName,
		ChangeSetName: &changeSetName,
		Description:   &description,
		Tags:          append(tags, c.sharedTags...),
	}

	input.SetChangeSetType(cloudformation.ChangeSetTypeUpdate)

	switch data := templateData.(type) {
	case TemplateBody:
		input.SetTemplateBody(string(data))
	case TemplateURL:
		input.SetTemplateURL(string(data))
	default:
		return fmt.Errorf("unknown template data type: %T", templateData)
	}

	input.SetCapabilities(capabilities)
	if cfnRole := c.roleARN; cfnRole != "" {
		input.SetRoleARN(cfnRole)
	}

	for k, v := range parameters {
		p := &cloudformation.Parameter{
			ParameterKey:   aws.String(k),
			ParameterValue: aws.String(v),
		}
		input.Parameters = append(input.Parameters, p)
	}

	logger.Debug("creating changeSet, input = %#v", input)
	s, err := c.cloudformationAPI.CreateChangeSet(input)
	if err != nil {
		return errors.Wrapf(err, "creating ChangeSet %q for stack %q", changeSetName, stackName)
	}
	logger.Debug("changeSet = %#v", s)
	return nil
}

func (c *StackCollection) doExecuteChangeSet(stackName string, changeSetName string) error {
	input := &cloudformation.ExecuteChangeSetInput{
		ChangeSetName: &changeSetName,
		StackName:     &stackName,
	}

	logger.Debug("executing changeSet, input = %#v", input)

	if _, err := c.cloudformationAPI.ExecuteChangeSet(input); err != nil {
		return errors.Wrapf(err, "executing CloudFormation ChangeSet %q for stack %q", changeSetName, stackName)
	}
	return nil
}

// DescribeStackChangeSet describes a ChangeSet by name
func (c *StackCollection) DescribeStackChangeSet(i *Stack, changeSetName string) (*ChangeSet, error) {
	input := &cloudformation.DescribeChangeSetInput{
		StackName:     i.StackName,
		ChangeSetName: &changeSetName,
	}
	if api.IsSetAndNonEmptyString(i.StackId) {
		input.StackName = i.StackId
	}
	resp, err := c.cloudformationAPI.DescribeChangeSet(input)
	if err != nil {
		return nil, errors.Wrapf(err, "describing CloudFormation ChangeSet %s for stack %s", changeSetName, *i.StackName)
	}
	return resp, nil
}
