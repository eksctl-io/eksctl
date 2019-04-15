package manager

import (
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/aws/aws-sdk-go/service/cloudtrail"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

const (
	resourcesRootPath = "Resources"
	outputsRootPath   = "Outputs"
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

// ChangeSet represents a CloudFormation ChangeSet
type ChangeSet = cloudformation.DescribeChangeSetOutput

// StackCollection stores the CloudFormation stack information
type StackCollection struct {
	provider   api.ClusterProvider
	spec       *api.ClusterConfig
	sharedTags []*cloudformation.Tag
}

func newTag(key, value string) *cloudformation.Tag {
	return &cloudformation.Tag{Key: &key, Value: &value}
}

// NewStackCollection create a stack manager for a single cluster
func NewStackCollection(provider api.ClusterProvider, spec *api.ClusterConfig) *StackCollection {
	tags := []*cloudformation.Tag{
		newTag(api.ClusterNameTag, spec.Metadata.Name),
	}
	for key, value := range spec.Metadata.Tags {
		tags = append(tags, newTag(key, value))
	}
	return &StackCollection{
		provider:   provider,
		spec:       spec,
		sharedTags: tags,
	}
}

// DoCreateStackRequest requests the creation of a CloudFormation stack
func (c *StackCollection) DoCreateStackRequest(i *Stack, templateBody []byte, tags, parameters map[string]string, withIAM bool, withNamedIAM bool) error {
	input := &cloudformation.CreateStackInput{
		StackName: i.StackName,
	}

	for _, t := range c.sharedTags {
		input.Tags = append(input.Tags, t)
	}
	for k, v := range tags {
		input.Tags = append(input.Tags, newTag(k, v))
	}

	input.SetTemplateBody(string(templateBody))

	if withIAM {
		input.SetCapabilities(stackCapabilitiesIAM)
	}

	if withNamedIAM {
		input.SetCapabilities(stackCapabilitiesNamedIAM)
	}

	if cfnRole := c.provider.CloudFormationRoleARN(); cfnRole != "" {
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
	s, err := c.provider.CloudFormation().CreateStack(input)
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
func (c *StackCollection) CreateStack(name string, stack builder.ResourceSet, tags, parameters map[string]string, errs chan error) error {
	i := &Stack{StackName: &name}
	templateBody, err := stack.RenderJSON()
	if err != nil {
		return errors.Wrapf(err, "rendering template for %q stack", *i.StackName)
	}

	if err := c.DoCreateStackRequest(i, templateBody, tags, parameters, stack.WithIAM(), stack.WithNamedIAM()); err != nil {
		return err
	}

	logger.Info("deploying stack %q", name)

	go c.waitUntilStackIsCreated(i, stack, errs)

	return nil
}

// UpdateStack will update a CloudFormation stack by creating and executing a ChangeSet
func (c *StackCollection) UpdateStack(stackName string, changeSetName string, description string, template []byte, parameters map[string]string) error {
	logger.Info(description)
	i := &Stack{StackName: &stackName}
	if err := c.doCreateChangeSetRequest(i, changeSetName, description, template, parameters, true); err != nil {
		return err
	}
	if err := c.doWaitUntilChangeSetIsCreated(i, changeSetName); err != nil {
		return err
	}
	changeSet, err := c.DescribeStackChangeSet(i, changeSetName)
	if err != nil {
		return err
	}
	logger.Debug("changes = %#v", changeSet.Changes)
	if err := c.doExecuteChangeSet(stackName, changeSetName); err != nil {
		logger.Warning("error executing Cloudformation changeSet %s in stack %s. Check the Cloudformation console for further details", changeSetName, stackName)
		return err
	}
	return c.doWaitUntilStackIsUpdated(i)
}

// DescribeStack describes a cloudformation stack.
func (c *StackCollection) DescribeStack(i *Stack) (*Stack, error) {
	input := &cloudformation.DescribeStacksInput{
		StackName: i.StackName,
	}
	if api.IsSetAndNonEmptyString(i.StackId) {
		input.StackName = i.StackId
	}
	resp, err := c.provider.CloudFormation().DescribeStacks(input)
	if err != nil {
		return nil, errors.Wrapf(err, "describing CloudFormation stack %q", *i.StackName)
	}
	return resp.Stacks[0], nil
}

// ListStacks gets all of CloudFormation stacks
func (c *StackCollection) ListStacks(nameRegex string, statusFilters ...string) ([]*Stack, error) {
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
	if err := c.provider.CloudFormation().ListStacksPages(input, pager); err != nil {
		return nil, err
	}
	if subErr != nil {
		return nil, subErr
	}
	return stacks, nil
}

// StackStatusIsNotTransitional will return true when stack statate is non-transitional
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
		cloudformation.StackStatusUpdateCompleteCleanupInProgress,
		cloudformation.StackStatusRollbackComplete,
		cloudformation.StackStatusUpdateRollbackCompleteCleanupInProgress,
		cloudformation.StackStatusUpdateRollbackComplete,
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
	i := &Stack{StackName: &name}
	s, err := c.DescribeStack(i)
	if err != nil {
		err = errors.Wrapf(err, "not able to get stack %q for deletion", name)
		stacks, newErr := c.ListStacks(fmt.Sprintf("^%s$", name), cloudformation.StackStatusDeleteComplete)
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
	i.StackId = s.StackId
	return c.DeleteStackBySpec(i)
}

// DeleteStackByNameSync sends a request to delete the stack, and waits until status is DELETE_COMPLETE;
// any errors will be written to errs channel, assume completion when nil is written, do not expect
// more then one error value on the channel, it's closed immediately after it is written to
func (c *StackCollection) DeleteStackByNameSync(name string, errs chan error) error {
	i, err := c.DeleteStackByName(name)
	if err != nil {
		return err
	}

	logger.Info("waiting for stack %q to get deleted", *i.StackName)

	go c.waitUntilStackIsDeleted(i, errs)

	return nil
}

// DeleteStackBySpec sends a request to delete the stack
func (c *StackCollection) DeleteStackBySpec(s *Stack) (*Stack, error) {
	for _, tag := range s.Tags {
		if matchesClusterName(*tag.Key, *tag.Value, c.spec.Metadata.Name) {
			input := &cloudformation.DeleteStackInput{
				StackName: s.StackId,
			}

			if cfnRole := c.provider.CloudFormationRoleARN(); cfnRole != "" {
				input = input.SetRoleARN(cfnRole)
			}

			if _, err := c.provider.CloudFormation().DeleteStack(input); err != nil {
				return nil, errors.Wrapf(err, "not able to delete stack %q", *s.StackName)
			}
			logger.Info("will delete stack %q", *s.StackName)
			return s, nil
		}
	}

	return nil, fmt.Errorf("cannot delete stack %q as it doesn't bare our %q, %q tags", *s.StackName,
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
	const ourStackRegexFmt = "^(eksctl|EKS)-%s-((cluster|nodegroup-.+)|(VPC|ServiceRole|ControlPlane|DefaultNodeGroup))$"
	return fmt.Sprintf(ourStackRegexFmt, name)
}

// DescribeStacks describes the existing stacks
func (c *StackCollection) DescribeStacks() ([]*Stack, error) {
	stacks, err := c.ListStacks(fmtStacksRegexForCluster(c.spec.Metadata.Name))
	if err != nil {
		return nil, errors.Wrapf(err, "describing CloudFormation stacks for %q", c.spec.Metadata.Name)
	}
	if len(stacks) == 0 {
		return nil, fmt.Errorf("no eksctl-managed CloudFormation stacks found for %q", c.spec.Metadata.Name)
	}
	return stacks, nil
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
	if err := c.provider.CloudFormation().DescribeStackEventsPages(input, pager); err != nil {
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
	if err := c.provider.CloudTrail().LookupEventsPages(input, pager); err != nil {
		return nil, errors.Wrapf(err, "looking up CloduTrail events for stack %q", *i.StackName)
	}

	return events, nil
}

func (c *StackCollection) doCreateChangeSetRequest(i *Stack, changeSetName string, description string, templateBody []byte,
	parameters map[string]string, withIAM bool) error {
	input := &cloudformation.CreateChangeSetInput{
		StackName:     i.StackName,
		ChangeSetName: &changeSetName,
		Description:   &description,
	}

	input.SetChangeSetType(cloudformation.ChangeSetTypeUpdate)

	input.SetTemplateBody(string(templateBody))

	if withIAM {
		input.SetCapabilities(stackCapabilitiesIAM)
	}

	if cfnRole := c.provider.CloudFormationRoleARN(); cfnRole != "" {
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
	s, err := c.provider.CloudFormation().CreateChangeSet(input)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("creating ChangeSet %q for stack %q", changeSetName, *i.StackName))
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

	if _, err := c.provider.CloudFormation().ExecuteChangeSet(input); err != nil {
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
	resp, err := c.provider.CloudFormation().DescribeChangeSet(input)
	if err != nil {
		return nil, errors.Wrapf(err, "describing CloudFormation ChangeSet %s for stack %s", changeSetName, *i.StackName)
	}
	return resp, nil
}
