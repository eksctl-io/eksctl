package spot

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudformation"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/kris-nova/logger"
	"github.com/spotinst/spotinst-sdk-go/service/ocean"
	oceanaws "github.com/spotinst/spotinst-sdk-go/service/ocean/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/spotinst-sdk-go/spotinst/client"
	"github.com/spotinst/spotinst-sdk-go/spotinst/featureflag"
	"github.com/spotinst/spotinst-sdk-go/spotinst/log"
	"github.com/spotinst/spotinst-sdk-go/spotinst/session"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/version"
	gfn "github.com/weaveworks/goformation/v4"
	gfncfn "github.com/weaveworks/goformation/v4/cloudformation"
)

// NewOceanVirtualNodeGroup returns a new NodeGroup object for the Ocean Virtual NodeGroup.
func NewOceanVirtualNodeGroup() *api.NodeGroup {
	ng := api.NewNodeGroup()
	ng.SpotOcean = new(api.SpotOceanVirtualNodeGroup)
	return ng
}

// NewOceanClusterNodeGroup returns a new NodeGroup object for the Ocean Cluster.
func NewOceanClusterNodeGroup(clusterSpec *api.ClusterConfig) *api.NodeGroup {
	ng := NewOceanVirtualNodeGroup()
	ng.Name = api.SpotOceanClusterNodeGroupName
	ng.PrivateNetworking = shouldUsePrivateNetworking(clusterSpec)
	// TODO(liran): Support is not available at the Ocean Virtual Nodegroup level.
	ng.EBSOptimized = shouldUseEBSOptimization(clusterSpec)
	api.SetNodeGroupDefaults(ng, clusterSpec.Metadata, false)
	return ng
}

// RunPreDelete executes pre-delete actions.
func RunPreDelete(ctx context.Context, provider api.ClusterProvider,
	clusterSpec *api.ClusterConfig, nodeGroups []*api.NodeGroup, stacks []*types.Stack,
	shouldDelete FilterFunc, roll bool, rollBatchSize int, plan bool) error {
	logger.Debug("ocean: executing pre-delete actions")

	// Only delete groups that are managed by Ocean and need to be deleted.
	oceanNodeGroups := make([]*api.NodeGroup, 0, len(nodeGroups))
	for _, ng := range nodeGroups {
		if shouldDelete(ng.Name) && IsNodeGroupManagedByOcean(ng, stacks) {
			oceanNodeGroups = append(oceanNodeGroups, ng)
		}
	}

	LoadFeatureFlags()
	if !plan {
		// Gracefully migrate all running workload.
		if len(oceanNodeGroups) > 0 && roll {
			if err := rollingUpdate(ctx, oceanNodeGroups, stacks, rollBatchSize); err != nil {
				return err
			}
		}

		// Update upstream credentials, if needed.
		if AllowCredentialsChanges.Enabled() {
			logger.Debug("ocean: updating credentials for existing nodegroups")
			for _, ng := range oceanNodeGroups {
				if err := UpdateCredentials(ctx, provider, ng, stacks); err != nil {
					return err
				}
			}
		}
	}

	stack, err := ShouldDeleteOceanCluster(stacks, shouldDelete)
	if err != nil {
		return err
	}
	if stack != nil {
		if !plan && AllowCredentialsChanges.Enabled() {
			ng := NewOceanClusterNodeGroup(clusterSpec)
			if err = UpdateCredentials(ctx, provider, ng, stacks); err != nil {
				return err
			}
		}

		// Allow post-delete actions on Ocean Cluster's nodegroup.
		clusterSpec.NodeGroups = append(clusterSpec.NodeGroups,
			NewOceanClusterNodeGroup(clusterSpec))
	}

	logger.Debug("ocean: successfully executed pre-delete actions")
	return nil
}

// ShouldCreateOceanCluster checks whether the Ocean Cluster should be created
// and returns its configuration.
func ShouldCreateOceanCluster(clusterSpec *api.ClusterConfig, stacks []*types.Stack) *api.NodeGroup {
	logger.Debug("ocean: checking whether cluster should be created")

	// If there's already an existing Ocean Cluster, we're done.
	clusterID := getOceanClusterIDFromStacks(stacks)
	if clusterID != "" {
		logger.Debug("ocean: cluster already exists (%s)", clusterID)
		return nil
	}

	var ng *api.NodeGroup
	if clusterSpec.SpotOcean != nil { // explicit
		ng = NewOceanClusterNodeGroup(clusterSpec)
	} else { // implicit
		for _, n := range clusterSpec.NodeGroups {
			if n.SpotOcean != nil {
				ng = NewOceanClusterNodeGroup(clusterSpec)
				break
			}
		}
	}
	if ng != nil {
		return ng
	}

	logger.Debug("ocean: no nodegroups found")
	return nil
}

// shouldUsePrivateNetworking returns true whether the Ocean Cluster should make
// use of private networking.
func shouldUsePrivateNetworking(clusterSpec *api.ClusterConfig) (private bool) {
	if clusterSpec.PrivateCluster != nil {
		private = clusterSpec.PrivateCluster.Enabled
	}
	if !private && clusterSpec.VPC != nil && clusterSpec.VPC.Subnets != nil {
		private = len(clusterSpec.VPC.Subnets.Private) > 0
	}
	if !private && len(clusterSpec.NodeGroups) > 0 {
		for _, ng := range clusterSpec.NodeGroups {
			if ng.SpotOcean != nil && ng.PrivateNetworking {
				private = true
				break
			}
		}
	}
	return
}

// shouldUseEBSOptimization returns true whether the Ocean Cluster should make
// use of EBS optimization.
func shouldUseEBSOptimization(clusterSpec *api.ClusterConfig) *bool {
	optimized := false
	if len(clusterSpec.NodeGroups) > 0 {
		for _, ng := range clusterSpec.NodeGroups {
			if ng.SpotOcean != nil && spotinst.BoolValue(ng.EBSOptimized) {
				optimized = true
				break
			}
		}
	}
	return spotinst.Bool(optimized)
}

// ShouldDeleteOceanCluster checks whether the Ocean Cluster should be deleted
// and returns its CloudFormation stack.
func ShouldDeleteOceanCluster(stacks []*types.Stack,
	shouldDelete FilterFunc) (*types.Stack, error) {
	logger.Debug("ocean: checking whether cluster should be deleted")

	// If there is no nodegroup for the Ocean Cluster, let's bail early.
	var oceanNodeGroupStack *types.Stack
	for _, s := range stacks {
		if getNodeGroupNameFromStack(s) == api.SpotOceanClusterNodeGroupName {
			oceanNodeGroupStack = s
			break
		}
	}
	if oceanNodeGroupStack == nil {
		logger.Debug("ocean: cluster does not exist")
		return nil, nil
	}

	// Check if there is at least one nodegroup that is not marked for deletion.
	if !shouldDelete(api.SpotOceanClusterNodeGroupName) {
		for _, s := range stacks {
			ngName := getNodeGroupNameFromStack(s)
			ng := &api.NodeGroup{NodeGroupBase: &api.NodeGroupBase{Name: ngName}}
			if !shouldDelete(ngName) && ngName != api.SpotOceanClusterNodeGroupName &&
				isStackStatusNotTransitional(s) && IsNodeGroupManagedByOcean(ng, stacks) {
				logger.Debug("ocean: at least one nodegroup remains "+
					"active (%s), skipping ocean cluster deletion", ngName)
				return nil, nil
			}
		}
	}

	logger.Debug("ocean: cluster should be deleted")
	return oceanNodeGroupStack, nil // all nodegroups are marked for deletion
}

// getOceanClusterIDFromStacks returns the Ocean Cluster identifier.
func getOceanClusterIDFromStacks(stacks []*types.Stack) (clusterID string) {
	collectors := map[string]outputs.Collector{
		outputs.NodeGroupSpotOceanClusterID: func(s string) error {
			clusterID = s
			return nil
		},
	}
	if len(stacks) > 0 {
		for _, s := range stacks {
			if getNodeGroupNameFromStack(s) != api.SpotOceanClusterNodeGroupName ||
				!isStackStatusNotTransitional(s) {
				continue
			}
			if err := outputs.Collect(*s, collectors, nil); err != nil {
				continue
			}
			if clusterID != "" {
				break
			}
		}
	}
	return
}

// getOceanVirtualNodeGroupIDFromStacks returns the Ocean Virtual Node Group ID.
func getOceanVirtualNodeGroupIDFromStacks(
	stacks []*types.Stack, ngName string) (vngID string) {
	collectors := map[string]outputs.Collector{
		outputs.NodeGroupSpotOceanLaunchSpecID: func(s string) error {
			vngID = s
			return nil
		},
	}
	for _, s := range stacks {
		if getNodeGroupNameFromStack(s) != ngName || !isStackStatusNotTransitional(s) {
			continue
		}
		if err := outputs.Collect(*s, collectors, nil); err != nil {
			continue
		}
		if vngID != "" {
			break
		}
	}
	return
}

// getNodeGroupNameFromStack returns the name of the nodegroup.
func getNodeGroupNameFromStack(stack *types.Stack) string {
	for _, tag := range stack.Tags {
		switch *tag.Key {
		case api.NodeGroupNameTag:
			return *tag.Value
		}
	}
	return ""
}

// GetStackByNodeGroupName returns the nodegroup by name.
func GetStackByNodeGroupName(name string, stacks []*types.Stack) *types.Stack {
	for _, stack := range stacks {
		if getNodeGroupNameFromStack(stack) == name {
			return stack
		}
	}
	return nil
}

// isStackStatusNotTransitional returns true when nodegroup status is non-transitional.
func isStackStatusNotTransitional(stack *types.Stack) bool {
	states := map[types.StackStatus]struct{}{
		types.StackStatusCreateComplete:         {},
		types.StackStatusUpdateComplete:         {},
		types.StackStatusRollbackComplete:       {},
		types.StackStatusUpdateRollbackComplete: {},
	}
	_, ok := states[stack.StackStatus]
	return ok
}

// IsNodeGroupManagedByOcean returns a boolean indicating whether the nodegroup is managed by Ocean.
func IsNodeGroupManagedByOcean(nodeGroup *api.NodeGroup, stacks []*types.Stack) bool {
	if nodeGroup.SpotOcean != nil { // fast path when using a config file
		return true
	}
	for _, stack := range stacks { // slow path when using a flag
		if nodeGroup.Name != getNodeGroupNameFromStack(stack) {
			continue
		}
		for _, tag := range stack.Tags {
			if spotinst.StringValue(tag.Key) == api.SpotOceanResourceTypeTag {
				return true
			}
		}
	}
	return false
}

const (
	// CredentialsTokenParameterKey specifies the name of the key associated
	// with the parameter which holds the user token.
	CredentialsTokenParameterKey = "SpotToken"
	// CredentialsAccountParameterKey specifies the name of the key associated
	// with the parameter which holds the user account.
	CredentialsAccountParameterKey = "SpotAccount"
)

// UpdateCredentials loads the user credentials from its local environment and
// updates the upstream credentials, stored in AWS CloudFormation, by updating
// the stack parameters. Users should set the `AllowCredentialsChanges` feature
// flag to avoid unnecessary calls caused by updating the AWS CloudFormation
// stack parameters.
func UpdateCredentials(ctx context.Context, provider api.ClusterProvider,
	ng *api.NodeGroup, stacks []*types.Stack) error {
	logger.Debug("ocean: updating credentials for nodegroup %q", ng.Name)

	// Find the stack by the name of the nodegroup.
	stack := GetStackByNodeGroupName(ng.Name, stacks)
	if stack == nil {
		logger.Debug("ocean: couldn't find stack for nodegroup %q", ng.Name)
		return nil
	}

	// Load user credentials.
	token, account, err := LoadCredentials()
	if err != nil {
		return err
	}

	// Update upstream credentials.
	if err := updateStackCredentials(ctx, provider, stack, token, account); err != nil {
		return err
	}

	logger.Debug("ocean: successfully updated upstream credentials for nodegroup %q", ng.Name)
	return nil
}

// updateStackCredentials updates the credentials stored in stack parameters.
func updateStackCredentials(ctx context.Context, provider api.ClusterProvider,
	stack *types.Stack, token, account string) error {

	var (
		cfnAPI  = provider.CloudFormation()
		cfnWait = true
	)

	template, err := updateTemplateCredentials(ctx, provider, stack, token, account)
	if err != nil {
		return err
	}

	input := &cloudformation.UpdateStackInput{
		StackName:    stack.StackName,
		Capabilities: []types.Capability{types.CapabilityCapabilityIam},
		TemplateBody: template,
		Parameters: []types.Parameter{
			{
				ParameterKey:   spotinst.String(CredentialsTokenParameterKey),
				ParameterValue: spotinst.String(token),
			},
			{
				ParameterKey:   spotinst.String(CredentialsAccountParameterKey),
				ParameterValue: spotinst.String(account),
			},
			{
				ParameterKey:   spotinst.String(FeatureFlagsParameterKey),
				ParameterValue: spotinst.String(LoadFeatureFlags()),
			},
		},
	}

	// isIgnorableError ignores errors that may occur while updating a stack.
	isIgnorableError := func(err string) bool {
		errs := []string{
			"no updates are to be performed",
		}
		for _, e := range errs {
			if strings.Contains(strings.ToLower(err), e) {
				return true
			}
		}
		return false
	}

	logger.Debug("ocean: updating stack %q", spotinst.StringValue(stack.StackName))
	if _, err = cfnAPI.UpdateStack(ctx, input); err != nil {
		if !isIgnorableError(err.Error()) {
			return fmt.Errorf("ocean: upstream error: %w", err)
		}
		cfnWait = false
		logger.Debug("ocean: local and upstream credentials are the "+
			"same; no updates needed for stack %q", stack.StackName)
	}

	// Wait until stack status is UPDATE_COMPLETE.
	if cfnWait {
		logger.Debug("ocean: waiting for stack update to complete")
		waiter := cloudformation.NewStackUpdateCompleteWaiter(provider.CloudFormation())
		params := &cloudformation.DescribeStacksInput{
			StackName: stack.StackName,
		}
		if err := waiter.Wait(ctx, params, provider.WaitTimeout()); err != nil {
			return fmt.Errorf("ocean: error waiting for stack update: %w", err)
		}
	}

	logger.Debug("ocean: successfully updated stack %q", spotinst.StringValue(stack.StackName))
	return nil
}

// updateTemplateCredentials updates the credentials stored in template parameters.
func updateTemplateCredentials(ctx context.Context, provider api.ClusterProvider,
	stack *types.Stack, token, account string) (*string, error) {

	input := &cloudformation.GetTemplateInput{
		StackName: stack.StackName,
	}

	output, err := provider.CloudFormation().GetTemplate(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("ocean: error getting template: %w", err)
	}
	if output.TemplateBody == nil {
		return nil, fmt.Errorf("ocean: empty template for stack %q", spotinst.StringValue(stack.StackName))
	}

	template, err := gfn.ParseJSON([]byte(spotinst.StringValue(output.TemplateBody)))
	if err != nil {
		return nil, fmt.Errorf("ocean: unexpected error parsing template: %w", err)
	}

	if _, ok := template.Parameters[CredentialsTokenParameterKey]; ok {
		template.Parameters[CredentialsTokenParameterKey] = gfncfn.Parameter{
			Type:    "String",
			Default: token,
		}
	}

	if _, ok := template.Parameters[CredentialsAccountParameterKey]; ok {
		template.Parameters[CredentialsAccountParameterKey] = gfncfn.Parameter{
			Type:    "String",
			Default: account,
		}
	}

	if _, ok := template.Parameters[FeatureFlagsParameterKey]; ok {
		template.Parameters[FeatureFlagsParameterKey] = gfncfn.Parameter{
			Type:    "String",
			Default: LoadFeatureFlags(),
		}
	}

	b, err := template.JSON()
	if err != nil {
		return nil, fmt.Errorf("ocean: unexpected error marshaling template: %w", err)
	}

	return spotinst.String(string(b)), nil
}

// LoadCredentials loads and returns the user credentials.
func LoadCredentials() (string, string, error) {
	logger.Debug("ocean: loading credentials")

	config := spotinst.DefaultConfig()
	c, err := config.Credentials.Get()
	if err != nil {
		return "", "", fmt.Errorf("ocean: error loading credentials: %w", err)
	}

	logger.Debug("ocean: loaded credentials for account: %s", c.Account)
	return c.Token, c.Account, nil
}

const (
	// Default ARN of the AWS Lambda function that should handle AWS CloudFormation requests.
	defaultServiceToken = "arn:aws:lambda:${AWS::Region}:178579023202:function:spotinst-cloudformation"
	// Name of the environment variable to read when loading a custom service token.
	envServiceToken = "SPOTINST_SERVICE_TOKEN"
)

// LoadServiceToken loads and returns the service token that should be use by
// AWS CloudFormation.
func LoadServiceToken() string {
	logger.Debug("ocean: loading service token")

	token := os.Getenv(envServiceToken)
	if token == "" {
		token = defaultServiceToken
	}

	logger.Debug("ocean: loaded service token: %s", token)
	return token
}

// AllowCredentialsChanges is a feature flag that controls whether eksctl should
// allow credentials changes.  When true, eksctl reloads the user credentials
// and attempts to update the relevant AWS CloudFormation stacks.
var AllowCredentialsChanges = featureflag.New("AllowCredentialsChanges", false)

// FeatureFlagsParameterKey specifies the name of the key associated with the
// parameter which holds all feature flags.
const FeatureFlagsParameterKey = "SpotFeatureFlags"

// LoadFeatureFlags reads the local feature flags from an environment variable
// and returns the upstream feature flags that should be configured for the
// resource handler.
func LoadFeatureFlags() (ff string) {
	// Local feature flags.
	{
		// Read feature flags from the environment.
		featureflag.Set(os.Getenv(featureflag.EnvVar))
		logger.Debug("ocean: loaded feature flags: %s", featureflag.All())
	}

	// Upstream feature flags.
	{
		// Avoid `Parameters: [SpotFeatureFlags] must have values` errors.
		ff = "None"

		// Credentials changes.
		if AllowCredentialsChanges.Enabled() {
			// When the user allows credentials changes, we have to configure the
			// opposite feature flag for the resource handler to avoid unnecessary
			// calls caused by updating the AWS CloudFormation stack parameters.
			ff = "IgnoreCredentialsChanges=true"
		}

		logger.Debug("ocean: configuring resource handler's feature flags: %s", ff)
	}

	return ff
}

// FilterFunc is a function that takes a nodegroup name and returns whether it
// should be filtered.
type FilterFunc func(ngName string) bool

// NewAlwaysFilter returns a FilterFunc that always returns true.
func NewAlwaysFilter() FilterFunc { return func(string) bool { return true } }

// NewContainsFilter returns a FilterFunc that return true if a nodegroup with
// the same name exists in the list of nodegroups.
func NewContainsFilter(nodeGroups []*api.NodeGroup) FilterFunc {
	return func(ngName string) bool {
		for _, ng := range nodeGroups {
			if ng.Name == ngName {
				return true
			}
		}
		return false
	}
}

// rollingUpdate gracefully migrates running workload in a rolling update fashion.
func rollingUpdate(ctx context.Context, nodeGroups []*api.NodeGroup,
	stacks []*types.Stack, batchSize int) error {
	logger.Debug("ocean: initiating a rolling update")

	// Resolve Cluster ID.
	clusterID := getOceanClusterIDFromStacks(stacks)
	if clusterID == "" {
		return fmt.Errorf("ocean: couldn't find cluster")
	}

	// Resolve Virtual Node Group IDs.
	var vngIDs []string
	for _, ng := range nodeGroups {
		vngID := getOceanVirtualNodeGroupIDFromStacks(stacks, ng.Name)
		if vngID == "" {
			continue
		}
		vngIDs = append(vngIDs, vngID)
	}
	if len(vngIDs) == 0 {
		return fmt.Errorf("ocean: couldn't find virtual node groups")
	}

	// Roll parameters.
	input := &oceanaws.CreateRollInput{
		Roll: &oceanaws.RollSpec{
			LaunchSpecIDs:                vngIDs,
			ClusterID:                    spotinst.String(clusterID),
			Comment:                      spotinst.String("created by @weaveworks/eksctl"),
			DisableLaunchSpecAutoScaling: spotinst.Bool(true),
		},
	}
	if batchSize > 0 {
		input.Roll.BatchSizePercentage = spotinst.Int(batchSize)
	}

	// isIgnorableError ignores errors that may occur while initiating a rolling update.
	isIgnorableError := func(err string) bool {
		errs := []string{
			"cluster has no active instances",
		}
		for _, e := range errs {
			if strings.Contains(strings.ToLower(err), e) {
				return true
			}
		}
		return false
	}

	logger.Debug("ocean: rolling virtual node groups: %s", strings.Join(vngIDs, "; "))
	svc := newService()
	output, err := svc.CreateRoll(ctx, input)
	if err != nil {
		spotErrs, ok := err.(client.Errors)
		if !ok {
			return fmt.Errorf("ocean: unexpected error: %w", err)
		}
		for _, spotErr := range spotErrs {
			if !isIgnorableError(spotErr.Message) {
				return fmt.Errorf("ocean: upstream error: %w", err)
			}
		}
		logger.Debug("ocean: no running instances, skipping rolling update")
		return nil
	}

	// Wait for the rolling update to complete.
	return waitUntilRollingUpdateComplete(ctx, svc, clusterID,
		spotinst.StringValue(output.Roll.ID))
}

// waitUntilRollingUpdateComplete waits for the rolling update to complete.
func waitUntilRollingUpdateComplete(
	ctx context.Context, svc oceanaws.Service, clusterID, rollID string) error {

	condFn := func() (bool, error) {
		input := &oceanaws.ReadRollInput{
			ClusterID: spotinst.String(clusterID),
			RollID:    spotinst.String(rollID),
		}
		output, err := svc.ReadRoll(ctx, input)
		if err != nil {
			return true, err
		}
		return checkRollingUpdateCompletionState(spotinst.StringValue(output.Roll.Status))
	}

	maxAttempts := 120 // 1 hour
	delay := 30 * time.Second

	for attempt := 1; ; attempt++ {
		logger.Debug("ocean: waiting for rolling update to complete (attempt: %d)", attempt)

		// Execute the condition function.
		done, err := condFn()
		if err != nil {
			return err
		}
		if done {
			break
		}

		// Fail if the maximum number of attempts is reached.
		if attempt == maxAttempts {
			return fmt.Errorf("ocean: exceeded wait attempts")
		}

		// Delay to wait before inspecting the resource again.
		if err := sleepWithContext(ctx, delay); err != nil {
			return fmt.Errorf("ocean: waiter context canceled: %w", err)
		}
	}

	logger.Debug("ocean: waiting for nodes to be drained")
	if err := sleepWithContext(ctx, 5*time.Minute); err != nil {
		return fmt.Errorf("ocean: waiter context canceled: %w", err)
	}

	logger.Debug("ocean: rolling update has been completed successfully")
	return nil
}

// checkRollingUpdateCompletionState returns true if a completion state is reached.
func checkRollingUpdateCompletionState(status string) (bool, error) {
	switch strings.ToUpper(status) {
	case "COMPLETED", "STOPPED":
		return true, nil
	case "FAILED":
		return true, fmt.Errorf("ocean: failed waiting for successful state")
	default:
		return false, nil
	}
}

// newService returns a new Ocean service.
func newService() oceanaws.Service {
	cfg := spotinst.DefaultConfig()
	cfg.WithLogger(newServiceLogger())
	cfg.WithUserAgent("weaveworks-eksctl/" + version.GetVersion())
	return ocean.New(session.New(cfg)).CloudProviderAWS()
}

// newServiceLogger returns a logger adapter.
func newServiceLogger() log.Logger {
	return log.LoggerFunc(func(format string, args ...interface{}) {
		logger.Debug(format+"\n", args...)
	})
}

// sleepWithContext will wait for the timer duration to expire, or the context
// is canceled. Which ever happens first. If the context is canceled the
// Context's error will be returned.
func sleepWithContext(ctx context.Context, dur time.Duration) error {
	t := time.NewTimer(dur)
	defer t.Stop()

	select {
	case <-t.C:
		break
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}
