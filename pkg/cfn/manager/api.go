package manager

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/autoscaling"
	asTypes "github.com/aws/aws-sdk-go-v2/service/autoscaling/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/smithy-go"

	cttypes "github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/waiter"
	"github.com/weaveworks/eksctl/pkg/version"
)

const (
	resourcesRootPath            = "Resources"
	resourceTypeAutoScalingGroup = "auto-scaling-group"
	outputsRootPath              = "Outputs"
	mappingsRootPath             = "Mappings"
	ourStackRegexFmt             = "^(eksctl|EKS)-%s-((cluster|nodegroup-.+|addon-.+|fargate|karpenter)|(VPC|ServiceRole|ControlPlane|DefaultNodeGroup))$"
	clusterStackRegex            = "eksctl-.*-cluster"
)

var (
	stackCapabilitiesIAM      = []types.Capability{types.CapabilityCapabilityIam}
	stackCapabilitiesNamedIAM = []types.Capability{types.CapabilityCapabilityNamedIam}
)

// Stack represents the CloudFormation stack
type Stack = types.Stack

// StackInfo hold the stack along with template and resources
type StackInfo struct {
	Stack     *Stack
	Resources []types.StackResource
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
	cloudformationAPI awsapi.CloudFormation
	ec2API            awsapi.EC2
	eksAPI            eksiface.EKSAPI
	iamAPI            awsapi.IAM
	cloudTrailAPI     awsapi.CloudTrail
	asgAPI            awsapi.ASG

	spec            *api.ClusterConfig
	disableRollback bool
	roleARN         string
	region          string
	waitTimeout     time.Duration
	sharedTags      []types.Tag
}

func newTag(key, value string) types.Tag {
	return types.Tag{Key: &key, Value: &value}
}

// NewStackCollection creates a stack manager for a single cluster
func NewStackCollection(provider api.ClusterProvider, spec *api.ClusterConfig) StackManager {
	tags := []types.Tag{
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
func (c *StackCollection) DoCreateStackRequest(ctx context.Context, i *Stack, templateData TemplateData, tags, parameters map[string]string, withIAM bool, withNamedIAM bool) error {
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
		input.TemplateBody = aws.String(string(data))
	case TemplateURL:
		input.TemplateURL = aws.String(string(data))
	default:
		return fmt.Errorf("unknown template data type: %T", templateData)
	}

	if withIAM {
		input.Capabilities = stackCapabilitiesIAM
	}

	if withNamedIAM {
		input.Capabilities = stackCapabilitiesNamedIAM
	}

	if cfnRole := c.roleARN; cfnRole != "" {
		input.RoleARN = aws.String(cfnRole)
	}

	for k, v := range parameters {
		input.Parameters = append(input.Parameters, types.Parameter{
			ParameterKey:   aws.String(k),
			ParameterValue: aws.String(v),
		})
	}

	logger.Debug("CreateStackInput = %#v", input)
	s, err := c.cloudformationAPI.CreateStack(ctx, input)
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
func (c *StackCollection) CreateStack(ctx context.Context, stackName string, resourceSet builder.ResourceSetReader, tags, parameters map[string]string, errs chan error) error {
	stack, err := c.createStackRequest(ctx, stackName, resourceSet, tags, parameters)
	if err != nil {
		return err
	}

	go c.waitUntilStackIsCreated(ctx, stack, resourceSet, errs)
	return nil
}

// createClusterStack creates the cluster stack
func (c *StackCollection) createClusterStack(ctx context.Context, stackName string, resourceSet builder.ResourceSetReader, errCh chan error) error {
	// Unlike with `createNodeGroupTask`, all tags are already set for the cluster stack
	stack, err := c.createStackRequest(ctx, stackName, resourceSet, nil, nil)
	if err != nil {
		return err
	}
	go func() {
		defer close(errCh)
		troubleshoot := func() {
			stack, err := c.DescribeStack(ctx, stack)
			if err != nil {
				logger.Info("error describing stack to troubleshoot the cause of the failure; "+
					"check the CloudFormation console for further details", err)
				return
			}

			logger.Critical("unexpected status %q while waiting for CloudFormation stack %q", stack.StackStatus, *stack.StackName)
			c.troubleshootStackFailureCause(ctx, stack, string(types.StackStatusCreateComplete))
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

func (c *StackCollection) createStackRequest(ctx context.Context, stackName string, resourceSet builder.ResourceSetReader, tags, parameters map[string]string) (*Stack, error) {
	stack := &Stack{StackName: &stackName}
	templateBody, err := resourceSet.RenderJSON()
	if err != nil {
		return nil, errors.Wrapf(err, "rendering template for %q stack", *stack.StackName)
	}

	if err := c.DoCreateStackRequest(ctx, stack, TemplateBody(templateBody), tags, parameters, resourceSet.WithIAM(), resourceSet.WithNamedIAM()); err != nil {
		return nil, err
	}

	logger.Info("deploying stack %q", stackName)
	return stack, nil
}

func (c *StackCollection) PropagateManagedNodeGroupTagsToASG(ngName string, ngTags map[string]string, asgNames []string, errCh chan error) error {
	go func() {
		defer close(errCh)
		// build the input tags for all ASGs attached to the managed nodegroup
		asgTags := []asTypes.Tag{}

		for _, asgName := range asgNames {
			// skip directly if not tags are required to be created
			if len(ngTags) == 0 {
				continue
			}

			// check if the number of tags on the ASG would go over the defined limit
			if err := c.checkASGTagsNumber(ngName, asgName, ngTags); err != nil {
				errCh <- err
				return
			}
			// build the list of tags to attach to the ASG
			for ngTagKey, ngTagValue := range ngTags {
				asgTag := asTypes.Tag{
					ResourceId:        aws.String(asgName),
					ResourceType:      aws.String(resourceTypeAutoScalingGroup),
					Key:               aws.String(ngTagKey),
					Value:             aws.String(ngTagValue),
					PropagateAtLaunch: aws.Bool(false),
				}
				asgTags = append(asgTags, asgTag)
			}
		}

		// consider the maximum number of tags we can create at once...
		var chunkedASGTags [][]asTypes.Tag
		chunkSize := builder.MaximumCreatedTagNumberPerCall
		for start := 0; start < len(asgTags); start += chunkSize {
			end := start + chunkSize
			if end > len(asgTags) {
				end = len(asgTags)
			}
			chunkedASGTags = append(chunkedASGTags, asgTags[start:end])
		}
		// ...then create all of them in a loop
		for _, asgTags := range chunkedASGTags {
			input := &autoscaling.CreateOrUpdateTagsInput{Tags: asgTags}
			if _, err := c.asgAPI.CreateOrUpdateTags(context.Background(), input); err != nil {
				errCh <- errors.Wrapf(err, "creating or updating asg tags for managed nodegroup %q", ngName)
				return
			}
		}
		errCh <- nil
	}()
	return nil
}

// checkASGTagsNumber limit considering the new propagated tags
func (c *StackCollection) checkASGTagsNumber(ngName, asgName string, propagatedTags map[string]string) error {
	tagsFilter := &autoscaling.DescribeTagsInput{
		Filters: []asTypes.Filter{
			{
				Name:   aws.String(resourceTypeAutoScalingGroup),
				Values: []string{asgName},
			},
		},
	}
	output, err := c.asgAPI.DescribeTags(context.Background(), tagsFilter)
	if err != nil {
		return errors.Wrapf(err, "describing asg %q tags for managed nodegroup %q", asgName, ngName)
	}
	asgTags := output.Tags
	// intersection of key tags to consider the number of tags going
	// to be attached to the ASG
	uniqueTagKeyCount := len(asgTags) + len(propagatedTags)
	for ngTagKey := range propagatedTags {
		for _, asgTag := range asgTags {
			// decrease the unique tag key count if there is a match
			if aws.StringValue(asgTag.Key) == ngTagKey {
				uniqueTagKeyCount--
				break
			}
		}
	}
	if uniqueTagKeyCount > builder.MaximumTagNumber {
		return fmt.Errorf("number of tags is exceeding the maximum amount for asg %d, was: %d", builder.MaximumTagNumber, uniqueTagKeyCount)
	}
	return nil
}

// UpdateStack will update a CloudFormation stack by creating and executing a ChangeSet
func (c *StackCollection) UpdateStack(ctx context.Context, options UpdateStackOptions) error {
	logger.Info(options.Description)
	if options.Stack == nil {
		i := &Stack{StackName: &options.StackName}
		// Read existing tags
		s, err := c.DescribeStack(ctx, i)
		if err != nil {
			return err
		}
		options.Stack = s
	} else {
		options.StackName = *options.Stack.StackName
	}
	if err := c.doCreateChangeSetRequest(ctx,
		options.StackName,
		options.ChangeSetName,
		options.Description,
		options.TemplateData,
		options.Parameters,
		options.Stack.Capabilities,
		options.Stack.Tags,
	); err != nil {
		return err
	}
	if err := c.doWaitUntilChangeSetIsCreated(ctx, options.Stack, options.ChangeSetName); err != nil {
		if _, ok := err.(*noChangeError); ok {
			return nil
		}
		return err
	}
	changeSet, err := c.DescribeStackChangeSet(ctx, options.Stack, options.ChangeSetName)
	if err != nil {
		return err
	}
	logger.Debug("changes = %#v", changeSet.Changes)
	if err := c.doExecuteChangeSet(ctx, options.StackName, options.ChangeSetName); err != nil {
		logger.Warning("error executing Cloudformation changeSet %s in stack %s. Check the Cloudformation console for further details", options.ChangeSetName, options.StackName)
		return err
	}
	if options.Wait {
		return c.doWaitUntilStackIsUpdated(ctx, options.Stack)
	}
	return nil
}

// DescribeStack describes a cloudformation stack.
func (c *StackCollection) DescribeStack(ctx context.Context, i *Stack) (*Stack, error) {
	input := &cloudformation.DescribeStacksInput{
		StackName: i.StackName,
	}
	if api.IsSetAndNonEmptyString(i.StackId) {
		input.StackName = i.StackId
	}
	resp, err := c.cloudformationAPI.DescribeStacks(ctx, input)
	if err != nil {
		return nil, errors.Wrapf(err, "describing CloudFormation stack %q", *i.StackName)
	}
	if len(resp.Stacks) == 0 {
		return nil, fmt.Errorf("no CloudFormation stack found for %s", *i.StackName)
	}
	return &resp.Stacks[0], nil
}

func IsStackDoesNotExistError(err error) bool {
	awsError, ok := errors.Unwrap(errors.Unwrap(err)).(*smithy.OperationError)
	return ok && strings.Contains(awsError.Error(), "ValidationError")

}

// GetManagedNodeGroupTemplate returns the template for a ManagedNodeGroup resource
func (c *StackCollection) GetManagedNodeGroupTemplate(ctx context.Context, options GetNodegroupOption) (string, error) {
	nodeGroupType, err := c.GetNodeGroupStackType(ctx, options)
	if err != nil {
		return "", err
	}

	if nodeGroupType != api.NodeGroupTypeManaged {
		return "", fmt.Errorf("%q is not a managed nodegroup", options.NodeGroupName)
	}

	stackName := c.makeNodeGroupStackName(options.NodeGroupName)
	templateBody, err := c.GetStackTemplate(ctx, stackName)
	if err != nil {
		return "", err
	}

	return templateBody, nil
}

// UpdateNodeGroupStack updates the nodegroup stack with the specified template
func (c *StackCollection) UpdateNodeGroupStack(ctx context.Context, nodeGroupName, template string, wait bool) error {
	stackName := c.makeNodeGroupStackName(nodeGroupName)
	return c.UpdateStack(ctx, UpdateStackOptions{
		StackName:     stackName,
		ChangeSetName: c.MakeChangeSetName("update-nodegroup"),
		Description:   "updating nodegroup stack",
		TemplateData:  TemplateBody(template),
		Wait:          wait,
	})
}

// ListStacksMatching gets all of CloudFormation stacks with names matching nameRegex.
func (c *StackCollection) ListStacksMatching(ctx context.Context, nameRegex string, statusFilters ...types.StackStatus) ([]*Stack, error) {
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
		input.StackStatusFilter = statusFilters
	}
	stacks := []*Stack{}

	paginator := cloudformation.NewListStacksPaginator(c.cloudformationAPI, input)

	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, s := range out.StackSummaries {
			if re.MatchString(*s.StackName) {
				stack, subErr = c.DescribeStack(ctx, &Stack{StackName: s.StackName, StackId: s.StackId})
				if subErr != nil {
					// this shouldn't return the error, but just stop the pagination and return whatever it gathered so far.
					return stacks, nil
				}
				stacks = append(stacks, stack)
			}
		}
	}

	return stacks, nil
}

// ListClusterStackNames gets all stack names matching regex
func (c *StackCollection) ListClusterStackNames(ctx context.Context) ([]string, error) {
	var stacks []string
	re, err := regexp.Compile(clusterStackRegex)
	if err != nil {
		return nil, errors.Wrap(err, "cannot list stacks")
	}
	input := &cloudformation.ListStacksInput{
		StackStatusFilter: defaultStackStatusFilter(),
	}

	paginator := cloudformation.NewListStacksPaginator(c.cloudformationAPI, input)

	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, s := range out.StackSummaries {
			if re.MatchString(*s.StackName) {
				stacks = append(stacks, *s.StackName)
			}
		}
	}

	return stacks, nil
}

// ListStacks gets all of CloudFormation stacks
func (c *StackCollection) ListStacks(ctx context.Context, statusFilters ...types.StackStatus) ([]*Stack, error) {
	return c.ListStacksMatching(ctx, fmtStacksRegexForCluster(c.spec.Metadata.Name), statusFilters...)
}

// StackStatusIsNotTransitional will return true when stack status is non-transitional
func (*StackCollection) StackStatusIsNotTransitional(s *Stack) bool {
	for _, state := range nonTransitionalReadyStackStatuses() {
		if s.StackStatus == state {
			return true
		}
	}
	return false
}

func nonTransitionalReadyStackStatuses() []types.StackStatus {
	return []types.StackStatus{
		types.StackStatusCreateComplete,
		types.StackStatusUpdateComplete,
		types.StackStatusRollbackComplete,
		types.StackStatusUpdateRollbackComplete,
	}
}

// StackStatusIsNotReady will return true when stack statate is non-ready
func (*StackCollection) StackStatusIsNotReady(s *Stack) bool {
	for _, state := range nonReadyStackStatuses() {
		if s.StackStatus == state {
			return true
		}
	}
	return false
}

func nonReadyStackStatuses() []types.StackStatus {
	return []types.StackStatus{
		types.StackStatusCreateInProgress,
		types.StackStatusCreateFailed,
		types.StackStatusRollbackInProgress,
		types.StackStatusRollbackFailed,
		types.StackStatusDeleteInProgress,
		types.StackStatusDeleteFailed,
		types.StackStatusUpdateInProgress,
		types.StackStatusUpdateCompleteCleanupInProgress,
		types.StackStatusUpdateRollbackInProgress,
		types.StackStatusUpdateRollbackFailed,
		types.StackStatusUpdateRollbackCompleteCleanupInProgress,
		types.StackStatusReviewInProgress,
	}
}

func allNonDeletedStackStatuses() []types.StackStatus {
	return []types.StackStatus{
		types.StackStatusCreateInProgress,
		types.StackStatusCreateFailed,
		types.StackStatusCreateComplete,
		types.StackStatusRollbackInProgress,
		types.StackStatusRollbackFailed,
		types.StackStatusRollbackComplete,
		types.StackStatusDeleteInProgress,
		types.StackStatusDeleteFailed,
		types.StackStatusUpdateInProgress,
		types.StackStatusUpdateCompleteCleanupInProgress,
		types.StackStatusUpdateComplete,
		types.StackStatusUpdateRollbackInProgress,
		types.StackStatusUpdateRollbackFailed,
		types.StackStatusUpdateRollbackCompleteCleanupInProgress,
		types.StackStatusUpdateRollbackComplete,
		types.StackStatusReviewInProgress,
	}
}

func defaultStackStatusFilter() []types.StackStatus {
	return allNonDeletedStackStatuses()
}

// DeleteStackBySpec sends a request to delete the stack
func (c *StackCollection) DeleteStackBySpec(ctx context.Context, s *Stack) (*Stack, error) {
	if !matchesCluster(c.spec.Metadata.Name, s.Tags) {
		return nil, fmt.Errorf("cannot delete stack %q as it doesn't bear our %q, %q tags", *s.StackName,
			fmt.Sprintf("%s:%s", api.OldClusterNameTag, c.spec.Metadata.Name),
			fmt.Sprintf("%s:%s", api.ClusterNameTag, c.spec.Metadata.Name))
	}

	input := &cloudformation.DeleteStackInput{
		StackName: s.StackId,
	}

	if cfnRole := c.roleARN; cfnRole != "" {
		input.RoleARN = &cfnRole
	}

	if _, err := c.cloudformationAPI.DeleteStack(ctx, input); err != nil {
		return nil, errors.Wrapf(err, "not able to delete stack %q", *s.StackName)
	}
	logger.Info("will delete stack %q", *s.StackName)
	return s, nil
}

func matchesCluster(clusterName string, tags []types.Tag) bool {
	for _, tag := range tags {
		switch *tag.Key {
		case api.ClusterNameTag, api.OldClusterNameTag:
			if *tag.Value == clusterName {
				return true
			}
		}
	}
	return false
}

// DeleteStackBySpecSync sends a request to delete the stack, and waits until status is DELETE_COMPLETE;
// any errors will be written to errs channel, assume completion when nil is written, do not expect
// more then one error value on the channel, it's closed immediately after it is written to
func (c *StackCollection) DeleteStackBySpecSync(ctx context.Context, s *Stack, errs chan error) error {
	i, err := c.DeleteStackBySpec(ctx, s)
	if err != nil {
		return err
	}

	logger.Info("waiting for stack %q to get deleted", *i.StackName)

	go c.waitUntilStackIsDeleted(ctx, i, errs)

	return nil
}

// DeleteStackSync sends a request to delete the stack, and waits until status is DELETE_COMPLETE;
func (c *StackCollection) DeleteStackSync(ctx context.Context, s *Stack) error {
	i, err := c.DeleteStackBySpec(ctx, s)
	if err != nil {
		return err
	}

	logger.Info("waiting for stack %q to get deleted", *i.StackName)
	return c.doWaitUntilStackIsDeleted(ctx, s)
}

func fmtStacksRegexForCluster(name string) string {
	return fmt.Sprintf(ourStackRegexFmt, name)
}

// DescribeStacks describes the existing stacks
func (c *StackCollection) DescribeStacks(ctx context.Context) ([]*Stack, error) {
	stacks, err := c.ListStacks(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "describing CloudFormation stacks for %q", c.spec.Metadata.Name)
	}
	if len(stacks) == 0 {
		logger.Debug("No stacks found for %s", c.spec.Metadata.Name)
	}
	return stacks, nil
}

func (c *StackCollection) GetClusterStackIfExists(ctx context.Context) (*Stack, error) {
	clusterStackNames, err := c.ListClusterStackNames(ctx)
	if err != nil {
		return nil, err
	}
	return c.getClusterStackFromList(ctx, clusterStackNames, c.spec.Metadata.Name)
}

func (c *StackCollection) HasClusterStackFromList(ctx context.Context, clusterStackNames []string, clusterName string) (bool, error) {
	stack, err := c.getClusterStackFromList(ctx, clusterStackNames, clusterName)
	return stack != nil, err
}

func (c *StackCollection) getClusterStackFromList(ctx context.Context, clusterStackNames []string, clusterName string) (*Stack, error) {
	clusterStackName := c.MakeClusterStackName()
	if clusterName != "" {
		clusterStackName = c.MakeClusterStackNameFromName(clusterName)
	}

	for _, stack := range clusterStackNames {
		if stack == clusterStackName {
			stack, err := c.DescribeStack(ctx, &types.Stack{StackName: &clusterStackName})
			if err != nil {
				return nil, err
			}
			if matchesCluster(clusterName, stack.Tags) {
				return stack, nil

			}
		}
	}
	return nil, nil
}

// DescribeStackEvents describes the events that have occurred on the stack
func (c *StackCollection) DescribeStackEvents(ctx context.Context, i *Stack) ([]types.StackEvent, error) {
	input := &cloudformation.DescribeStackEventsInput{
		StackName: i.StackName,
	}
	if api.IsSetAndNonEmptyString(i.StackId) {
		input.StackName = i.StackId
	}

	stackEvents, err := c.cloudformationAPI.DescribeStackEvents(ctx, input)
	if err != nil {
		return nil, errors.Wrapf(err, "describing CloudFormation stack %q events", *i.StackName)
	}

	return stackEvents.StackEvents, nil
}

func (c *StackCollection) LookupCloudTrailEvents(ctx context.Context, i *Stack) ([]cttypes.Event, error) {
	input := &cloudtrail.LookupEventsInput{
		LookupAttributes: []cttypes.LookupAttribute{{
			AttributeKey:   cttypes.LookupAttributeKeyResourceName,
			AttributeValue: i.StackId,
		}},
	}

	var events []cttypes.Event
	paginator := cloudtrail.NewLookupEventsPaginator(c.cloudTrailAPI, input)
	for paginator.HasMorePages() {
		out, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, errors.Wrapf(err, "looking up CloudTrail events for stack %q", *i.StackName)
		}
		events = append(events, out.Events...)
	}

	return events, nil
}

func (c *StackCollection) doCreateChangeSetRequest(ctx context.Context, stackName, changeSetName, description string, templateData TemplateData,
	parameters map[string]string, capabilities []types.Capability, tags []types.Tag) error {
	input := &cloudformation.CreateChangeSetInput{
		StackName:     &stackName,
		ChangeSetName: &changeSetName,
		Description:   &description,
		Tags:          append(tags, c.sharedTags...),
	}

	input.ChangeSetType = types.ChangeSetTypeUpdate

	switch data := templateData.(type) {
	case TemplateBody:
		input.TemplateBody = aws.String(string(data))
	case TemplateURL:
		input.TemplateURL = aws.String(string(data))
	default:
		return fmt.Errorf("unknown template data type: %T", templateData)
	}

	input.Capabilities = capabilities
	if cfnRole := c.roleARN; cfnRole != "" {
		input.RoleARN = &cfnRole
	}

	for k, v := range parameters {
		p := types.Parameter{
			ParameterKey:   aws.String(k),
			ParameterValue: aws.String(v),
		}
		input.Parameters = append(input.Parameters, p)
	}

	logger.Debug("creating changeSet, input = %#v", input)
	s, err := c.cloudformationAPI.CreateChangeSet(ctx, input)
	if err != nil {
		return errors.Wrapf(err, "creating ChangeSet %q for stack %q", changeSetName, stackName)
	}
	logger.Debug("changeSet = %#v", s)
	return nil
}

func (c *StackCollection) doExecuteChangeSet(ctx context.Context, stackName string, changeSetName string) error {
	input := &cloudformation.ExecuteChangeSetInput{
		ChangeSetName: &changeSetName,
		StackName:     &stackName,
	}

	logger.Debug("executing changeSet, input = %#v", input)

	if _, err := c.cloudformationAPI.ExecuteChangeSet(ctx, input); err != nil {
		return errors.Wrapf(err, "executing CloudFormation ChangeSet %q for stack %q", changeSetName, stackName)
	}
	return nil
}

// DescribeStackChangeSet describes a ChangeSet by name
func (c *StackCollection) DescribeStackChangeSet(ctx context.Context, i *Stack, changeSetName string) (*ChangeSet, error) {
	input := &cloudformation.DescribeChangeSetInput{
		StackName:     i.StackName,
		ChangeSetName: &changeSetName,
	}
	if api.IsSetAndNonEmptyString(i.StackId) {
		input.StackName = i.StackId
	}
	resp, err := c.cloudformationAPI.DescribeChangeSet(ctx, input)
	if err != nil {
		return nil, errors.Wrapf(err, "describing CloudFormation ChangeSet %s for stack %s", changeSetName, *i.StackName)
	}
	return resp, nil
}
