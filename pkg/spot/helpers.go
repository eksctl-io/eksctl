package spot

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/kris-nova/logger"
	"github.com/spotinst/spotinst-sdk-go/service/ocean"
	oceanaws "github.com/spotinst/spotinst-sdk-go/service/ocean/providers/aws"
	"github.com/spotinst/spotinst-sdk-go/spotinst"
	"github.com/spotinst/spotinst-sdk-go/spotinst/client"
	"github.com/spotinst/spotinst-sdk-go/spotinst/credentials"
	"github.com/spotinst/spotinst-sdk-go/spotinst/featureflag"
	"github.com/spotinst/spotinst-sdk-go/spotinst/log"
	"github.com/spotinst/spotinst-sdk-go/spotinst/session"
	"github.com/weaveworks/eksctl/pkg/addons"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/cfn/outputs"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	kubewrapper "github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/version"
)

// NewNodeGroup creates a new NodeGroup, and returns a pointer to it.
func NewNodeGroup() *api.NodeGroup {
	ng := api.NewNodeGroup()
	ng.SpotOcean = new(api.NodeGroupSpotOcean)
	return ng
}

// NewNodeGroupForOcean creates a new NodeGroup, and returns a pointer to it.
func NewNodeGroupForOcean() *api.NodeGroup {
	ng := NewNodeGroup()
	ng.Name = api.SpotOceanNodeGroupName
	return ng
}

// RunPreCreation executes post-creation actions.
func RunPostCreation(clusterConfig *api.ClusterConfig, clientSet kubernetes.Interface,
	rawClient *kubewrapper.RawClient, updateAuthConfigMap bool) error {
	logger.Debug("ocean: executing post-creation actions")

	for _, ng := range clusterConfig.NodeGroups {
		if ng.SpotOcean == nil {
			logger.Debug("ocean: skipping nodegroup %q", ng.Name)
			continue
		}
		logger.Debug("ocean: handling nodegroup %q", ng.Name)

		// Allow nodes that are launched by Ocean to join to the cluster.
		// We have to do it before all other nodegroups to prevent `WaitForNodes`
		// to wait forever.
		if updateAuthConfigMap {
			logger.Debug("ocean: updating auth configmap")
			if err := authconfigmap.AddNodeGroup(clientSet, ng); err != nil {
				return fmt.Errorf("ocean: error updaing auth configmap: %w", err)
			}
		}

		// Install the Ocean controller.
		if ng.Name == api.SpotOceanNodeGroupName {
			logger.Debug("ocean: installing addons")
			controller := addons.NewSpotOceanController(
				rawClient,
				clusterConfig,
				false,
				spotinst.StringValue(profileOrDefault(ng)))
			if err := controller.Deploy(); err != nil {
				return fmt.Errorf("ocean: error installing controller: %w", err)
			}
		}
	}

	return nil
}

// RunPreDeletion executes pre-deletion actions.
func RunPreDeletion(clusterProvider api.ClusterProvider,
	clusterConfig *api.ClusterConfig, nodeGroups []*api.NodeGroup,
	stacks []*cloudformation.Stack, shouldDelete DeleteFilter, roll bool,
	rollBatchSize int, plan bool) error {
	logger.Debug("ocean: executing pre-deletion actions")

	// Filter Ocean nodegroups that are marked for deletion.
	oceanNodeGroups := make([]*api.NodeGroup, 0, len(nodeGroups))
	for _, ng := range nodeGroups {
		if shouldDelete(ng.Name) && IsNodeGroupManagedByOcean(ng, stacks) {
			oceanNodeGroups = append(oceanNodeGroups, ng)
		}
	}

	LoadFeatureFlags()
	if !plan && len(oceanNodeGroups) > 0 {
		// Gracefully migrate all running workload.
		if roll {
			if err := initiateRoll(oceanNodeGroups, stacks, rollBatchSize); err != nil {
				return err
			}
		}

		// Update upstream credentials, if needed.
		if AllowCredentialsChanges.Enabled() {
			logger.Debug("ocean: updating credentials for existing nodegroups")

			for _, ng := range oceanNodeGroups {
				if err := UpdateCredentials(clusterProvider, clusterConfig,
					ng.Name, stacks); err != nil {
					return err
				}
			}
		}
	}

	_, shouldDeleteOcean, err := ShouldDeleteOceanNodeGroup(stacks, shouldDelete)
	if err != nil {
		return err
	}
	if shouldDeleteOcean {
		if !plan && AllowCredentialsChanges.Enabled() {
			if err = UpdateCredentials(clusterProvider,
				clusterConfig, api.SpotOceanNodeGroupName, stacks); err != nil {
				return err
			}
		}

		// Allow post-deletion actions to be performed on Ocean as well.
		clusterConfig.NodeGroups = append(clusterConfig.NodeGroups,
			NewNodeGroupForOcean())
	}

	return nil
}

// ErrSpotMultipleDefaultLaunchSpecs represents an error in detecting the default
// nodegroup since more than one has been configured as such.
var ErrSpotMultipleDefaultLaunchSpecs = errors.New("ocean: unable to detect " +
	"default ocean launch spec: multiple nodegroups configured with " +
	"`spot.metadata.defaultLaunchSpec: true`")

// ShouldCreateOceanNodeGroup checks whether the nodegroup of the Ocean cluster
// should be created and, if so, returns its NodeGroup configuration.
func ShouldCreateOceanNodeGroup(nodeGroups []*api.NodeGroup,
	stacks []*cloudformation.Stack) (*api.NodeGroup, bool, error) {
	logger.Debug("ocean: checking whether cluster should be created")

	var (
		oceanNodeGroups           = make([]*api.NodeGroup, 0, len(nodeGroups))
		oceanNodeGroup            *api.NodeGroup
		desired, minimum, maximum *int
	)

	// If there are no Ocean nodegroups, let's bail early.
	for _, ng := range nodeGroups {
		if ng.SpotOcean != nil {
			oceanNodeGroups = append(oceanNodeGroups, ng)
		}
	}
	if len(oceanNodeGroups) == 0 {
		logger.Debug("ocean: no nodegroups found")
		return nil, false, nil
	}

	// If there is already an existing cluster, we're done.
	clusterID := getOceanClusterIDFromStacks(stacks)
	if clusterID != "" {
		logger.Debug("ocean: cluster already exists")
		return nil, false, nil
	}

	// Find the default nodegroup and calculate the capacity.
	for _, ng := range oceanNodeGroups {
		// Is this the default nodegroup?
		if ng.SpotOcean.Metadata != nil &&
			spotinst.BoolValue(ng.SpotOcean.Metadata.DefaultLaunchSpec) {
			if oceanNodeGroup != nil {
				logger.Debug("ocean: multiple default nodegroups (%q and %q)",
					ng.Name, oceanNodeGroup.Name)
				return nil, false, ErrSpotMultipleDefaultLaunchSpecs
			}
			oceanNodeGroup = ng.DeepCopy()
		}

		// Sum up the capacity from all nodegroups.
		if ng.DesiredCapacity != nil {
			if desired == nil {
				desired = new(int)
			}
			*desired += spotinst.IntValue(ng.DesiredCapacity)
		}
		if ng.MinSize != nil {
			if minimum == nil {
				minimum = new(int)
			}
			*minimum += spotinst.IntValue(ng.MinSize)
		}
		if ng.MaxSize != nil {
			if maximum == nil {
				maximum = new(int)
			}
			*maximum += spotinst.IntValue(ng.MaxSize)
		}
	}

	// No default nodegroup. Take the first one.
	if oceanNodeGroup == nil {
		oceanNodeGroup = oceanNodeGroups[0].DeepCopy()
	}
	logger.Debug("ocean: using default nodegroup %q", oceanNodeGroup.Name)

	// Set the capacity.
	oceanNodeGroup.DesiredCapacity = desired
	oceanNodeGroup.MinSize = minimum
	oceanNodeGroup.MaxSize = maximum

	// Default of one node to run cluster-controller/metrics-server.
	if desired == nil && minimum == nil {
		oceanNodeGroup.DesiredCapacity = spotinst.Int(1)
		oceanNodeGroup.MinSize = spotinst.Int(1)
	}

	// Configure the nodegroup name.
	oceanNodeGroup.Name = api.SpotOceanNodeGroupName
	oceanNodeGroup.Labels[api.NodeGroupNameLabel] = api.SpotOceanNodeGroupName

	logger.Debug("ocean: cluster should be created")
	return oceanNodeGroup, true, nil
}

// ShouldDeleteOceanNodeGroup checks whether the nodegroup of the Ocean cluster
// should be deleted and, if so, returns its Cloud Formation stack.
func ShouldDeleteOceanNodeGroup(stacks []*cloudformation.Stack,
	shouldDelete func(string) bool) (*cloudformation.Stack, bool, error) {

	logger.Debug("ocean: checking whether cluster should be deleted")
	var oceanNodeGroupStack *cloudformation.Stack

	// If there is no nodegroup for the Ocean cluster, let's bail early.
	for _, s := range stacks {
		if getNodeGroupNameFromStack(s) == api.SpotOceanNodeGroupName {
			oceanNodeGroupStack = s
			break
		}
	}
	if oceanNodeGroupStack == nil {
		logger.Debug("ocean: cluster does not exist; nothing to delete")
		return nil, false, nil
	}

	// Do not delete if there is at least one nodegroup that is not marked for deletion.
	if !shouldDelete(api.SpotOceanNodeGroupName) {
		for _, s := range stacks {
			ngName := getNodeGroupNameFromStack(s)
			ng := &api.NodeGroup{NodeGroupBase: &api.NodeGroupBase{Name: ngName}}

			if !shouldDelete(ngName) && ngName != api.SpotOceanNodeGroupName &&
				isStackStatusNotTransitional(s) && IsNodeGroupManagedByOcean(ng, stacks) {
				logger.Debug("ocean: at least one nodegroup remains "+
					"active (%s); skipping ocean cluster deletion", ngName)
				return nil, false, nil
			}
		}
	}

	logger.Debug("ocean: cluster should be deleted")
	return oceanNodeGroupStack, true, nil // all nodegroups are marked for deletion
}

// getOceanClusterIDFromStacks returns the Ocean Cluster identifier.
func getOceanClusterIDFromStacks(stacks []*cloudformation.Stack) string {
	var clusterID string

	collectors := map[string]outputs.Collector{
		outputs.NodeGroupSpotOceanClusterID: func(s string) error {
			clusterID = s
			return nil
		},
	}

	for _, s := range stacks {
		if getNodeGroupNameFromStack(s) != api.SpotOceanNodeGroupName ||
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

	return clusterID
}

// getOceanLaunchSpecIDFromStacks returns the Ocean Launch Spec identifier.
func getOceanLaunchSpecIDFromStacks(stacks []*cloudformation.Stack, ngName string) string {
	var specID string

	collectors := map[string]outputs.Collector{
		outputs.NodeGroupSpotOceanLaunchSpecID: func(s string) error {
			specID = s
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
		if specID != "" {
			break
		}
	}

	return specID
}

// getNodeGroupNameFromStack returns the name of the nodegroup.
func getNodeGroupNameFromStack(stack *cloudformation.Stack) string {
	for _, tag := range stack.Tags {
		switch *tag.Key {
		case api.NodeGroupNameTag:
			return *tag.Value
		}
	}
	return ""
}

// getStackByNodeGroupName returns the nodegroup by name.
func getStackByNodeGroupName(name string, stacks []*cloudformation.Stack) *cloudformation.Stack {
	for _, stack := range stacks {
		if getNodeGroupNameFromStack(stack) == name {
			return stack
		}
	}
	return nil
}

// isStackStatusNotTransitional returns true when nodegroup status is non-transitional.
func isStackStatusNotTransitional(stack *cloudformation.Stack) bool {
	states := map[string]struct{}{
		cloudformation.StackStatusCreateComplete:         {},
		cloudformation.StackStatusUpdateComplete:         {},
		cloudformation.StackStatusRollbackComplete:       {},
		cloudformation.StackStatusUpdateRollbackComplete: {},
	}
	_, ok := states[*stack.StackStatus]
	return ok
}

// IsNodeGroupManagedByOcean returns a boolean indicating whether the nodegroup is managed by Ocean.
func IsNodeGroupManagedByOcean(nodeGroup *api.NodeGroup, stacks []*cloudformation.Stack) bool {
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
	// Name of the key associated with the parameter that holds the user token.
	CredentialsTokenParameterKey = "SpotToken"
	// Name of the key associated with the parameter that holds the user account.
	CredentialsAccountParameterKey = "SpotAccount"
)

// UpdateCredentials loads the user credentials from its local environment and
// updates the upstream credentials, stored in AWS CloudFormation, by updating
// the stack parameters.  Users should set the `AllowCredentialsChanges` feature
// flag to avoid unnecessary calls caused by updating the AWS CloudFormation
// stack parameters.
func UpdateCredentials(clusterProvider api.ClusterProvider, clusterConfig *api.ClusterConfig,
	ngName string, stacks []*cloudformation.Stack) error {
	logger.Debug("ocean: updating credentials for nodegroup %q", ngName)

	// Find the stack by the name of the nodegroup.
	stack := getStackByNodeGroupName(ngName, stacks)
	if stack == nil {
		logger.Debug("ocean: couldn't find stack for nodegroup %q", ngName)
		return nil
	}

	// Set the credentials profile, if any.
	var profile *string
	for _, ng := range clusterConfig.NodeGroups {
		if ng.Name == getNodeGroupNameFromStack(stack) && ng.SpotOcean != nil {
			profile = profileOrDefault(ng)
		}
	}

	// Load user credentials.
	token, account, err := LoadCredentials(profile)
	if err != nil {
		return err
	}

	// Update upstream credentials and reuse the existing template.
	if err := updateUpstreamCredentials(clusterProvider, stack, token, account); err != nil {
		return err
	}

	logger.Debug("ocean: successfully updated credentials for nodegroup %q", ngName)
	return nil
}

// updateUpstreamCredentials updates the upstream credentials, stored in AWS
// CloudFormation, by updating the stack parameters.
func updateUpstreamCredentials(clusterProvider api.ClusterProvider,
	stack *cloudformation.Stack, token, account string) error {

	var (
		cfnAPI  = clusterProvider.CloudFormation()
		cfnWait = true
	)

	// Set parameters.
	input := &cloudformation.UpdateStackInput{
		StackName:           stack.StackName,
		Capabilities:        spotinst.StringSlice([]string{cloudformation.CapabilityCapabilityIam}),
		UsePreviousTemplate: spotinst.Bool(true),
		Parameters: []*cloudformation.Parameter{
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
				ParameterValue: spotinst.String(convertFeatureFlags()),
			},
		},
	}

	// Update stack parameters.
	logger.Debug("ocean: updating stack %q", spotinst.StringValue(stack.StackName))
	_, err := cfnAPI.UpdateStack(input)
	if err != nil {
		awsErr, ok := err.(awserr.Error)
		if !ok {
			return err
		}
		if !ignoreUpdateStackError(awsErr.Message()) {
			return err
		}
		cfnWait = false
		logger.Debug("ocean: new and old credentials are same; no updates "+
			"are to be performed for stack %q", stack.StackName)
	}

	// Wait until stack status is UPDATE_COMPLETE.
	if cfnWait {
		logger.Debug("ocean: waiting for stack update to complete")
		input := &cloudformation.DescribeStacksInput{
			StackName: stack.StackName,
		}
		if err := cfnAPI.WaitUntilStackUpdateComplete(input); err != nil {
			return err
		}
	}

	return nil
}

// ignoreUpdateStackError ignores errors that may occur while updating a stack.
func ignoreUpdateStackError(errMsg string) bool {
	errMsgs := []string{
		"no updates are to be performed",
	}
	for _, msg := range errMsgs {
		if strings.Contains(strings.ToLower(errMsg), msg) {
			return true
		}
	}
	return false
}

// LoadCredentials loads and returns the user credentials.
func LoadCredentials(profile *string) (string, string, error) {
	logger.Debug("ocean: loading credentials")

	providers := []credentials.Provider{
		&credentials.EnvProvider{},
		&credentials.FileProvider{Profile: spotinst.StringValue(profile)},
	}

	config := spotinst.DefaultConfig()
	config.WithCredentials(credentials.NewChainCredentials(providers...))

	c, err := config.Credentials.Get()
	if err != nil {
		return "", "", err
	}

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
func LoadServiceToken() (string, error) {
	logger.Debug("ocean: loading service token")

	v := os.Getenv(envServiceToken)
	if v == "" {
		v = defaultServiceToken
	}

	logger.Debug("ocean: will use service token %q", v)
	return v, nil
}

// AllowCredentialsChanges is a feature flag that controls whether eksctl should
// allow credentials changes.  When true, eksctl reloads the user credentials
// and attempts to update the relevant AWS CloudFormation stacks.
var AllowCredentialsChanges = featureflag.New("AllowCredentialsChanges", false)

// Name of the key associated with the parameter that holds all feature flags.
const FeatureFlagsParameterKey = "SpotFeatureFlags"

// LoadFeatureFlags reads the feature flags from an environment variable.
func LoadFeatureFlags() string {
	featureflag.Set(os.Getenv(featureflag.EnvVar))
	logger.Debug("ocean: will use feature flags %q", featureflag.All())

	// Convert to upstream feature flags, if needed.
	return convertFeatureFlags()
}

// convertFeatureFlags returns the upstream feature flags that should be
// configured for the resource handler.
func convertFeatureFlags() string {
	ff := "None" // avoids `Parameters: [SpotFeatureFlags] must have values` errors

	if AllowCredentialsChanges.Enabled() {
		// When the user allows credentials changes, we have to configure the
		// opposite feature flag for the resource handler to avoid unnecessary
		// calls caused by updating the AWS CloudFormation stack parameters.
		ff = "IgnoreCredentialsChanges=true"
	}

	logger.Debug("ocean: will set feature flags %q", ff)
	return ff
}

// profileOrDefault returns the credentials profile that should be loaded.
func profileOrDefault(nodeGroup *api.NodeGroup) *string {
	profile := credentials.DefaultProfile()

	if nodeGroup.SpotOcean != nil && nodeGroup.SpotOcean.Metadata != nil {
		profile = spotinst.StringValue(nodeGroup.SpotOcean.Metadata.Profile)
	}

	return &profile
}

// DeleteFilter represents the type definition for a delete filter.
type DeleteFilter func(ngName string) bool

// NewDeleteAllFilter returns a DeleteFilter that always returns true.
func NewDeleteAllFilter() DeleteFilter {
	return func(_ string) bool {
		return true
	}
}

// NewDeleteIncludedFilter returns a DeleteFilter that returns true whether the
// nodegroup is included.
func NewDeleteIncludedFilter(nodeGroups []*api.NodeGroup) DeleteFilter {
	return func(ngName string) bool {
		for _, ng := range nodeGroups {
			if ng.Name == ngName {
				return true
			}
		}
		return false
	}
}

// initiateRoll initiates a cluster roll to gracefully migrate all running workload.
func initiateRoll(nodeGroups []*api.NodeGroup, stacks []*cloudformation.Stack, batchSize int) error {
	logger.Debug("ocean: initiating a cluster roll")

	// Resolve Cluster ID.
	clusterID := getOceanClusterIDFromStacks(stacks)
	if clusterID == "" {
		return fmt.Errorf("ocean: couldn't find cluster")
	}

	// Resolve Launch Spec IDs.
	var specIDs []string
	for _, ng := range nodeGroups {
		specID := getOceanLaunchSpecIDFromStacks(stacks, ng.Name)
		if specID == "" {
			continue
		}
		specIDs = append(specIDs, specID)
	}
	if len(specIDs) == 0 {
		return fmt.Errorf("ocean: couldn't find launch specs")
	}

	// Roll parameters.
	input := &oceanaws.CreateRollInput{
		Roll: &oceanaws.RollSpec{
			LaunchSpecIDs:                specIDs,
			ClusterID:                    spotinst.String(clusterID),
			Comment:                      spotinst.String(fmt.Sprintf("created by eksctl")),
			DisableLaunchSpecAutoScaling: spotinst.Bool(true),
		},
	}
	if batchSize > 0 {
		input.Roll.BatchSizePercentage = spotinst.Int(batchSize)
	}

	ctx := context.Background()
	svc := newService()

	logger.Debug("ocean: rolling launch specs %q", strings.Join(specIDs, "; "))
	output, err := svc.CreateRoll(ctx, input)
	if err != nil {
		spotErrs, ok := err.(client.Errors)
		if !ok {
			return err
		}
		for _, spotErr := range spotErrs {
			if !ignoreCreateRollError(spotErr.Message) {
				return err
			}
		}
		logger.Debug("ocean: launch specs has no active instances; nothing to roll")
		return nil
	}

	// Wait for cluster roll to complete.
	return waitUntilRollComplete(ctx, svc, clusterID, spotinst.StringValue(output.Roll.ID))
}

// ignoreCreateRollError ignores errors that may occur while initiating a cluster roll.
func ignoreCreateRollError(errMsg string) bool {
	errMsgs := []string{
		"cluster has no active instances",
	}
	for _, msg := range errMsgs {
		if strings.Contains(strings.ToLower(errMsg), msg) {
			return true
		}
	}
	return false
}

// waitUntilRollComplete waits for the cluster roll to complete.
func waitUntilRollComplete(ctx context.Context, svc oceanaws.Service, clusterID, rollID string) error {
	condFn := func() (bool, error) {
		input := &oceanaws.ReadRollInput{
			ClusterID: spotinst.String(clusterID),
			RollID:    spotinst.String(rollID),
		}

		output, err := svc.ReadRoll(ctx, input)
		if err != nil {
			return true, err
		}

		return checkRollCompletionState(output.Roll.Status)
	}

	maxAttempts := 120
	delay := 30 * time.Second

	for attempt := 1; ; attempt++ {
		logger.Debug("ocean: waiting for cluster roll to complete (attempt: %d)", attempt)

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
		if err := aws.SleepWithContext(ctx, delay); err != nil {
			return fmt.Errorf("ocean: waiter context canceled: %v", err)
		}
	}

	return nil
}

// checkRollCompletionState returns true if a completion state is reached.
func checkRollCompletionState(status *string) (bool, error) {
	states := map[string]request.WaiterState{
		"COMPLETED": request.SuccessWaiterState,
		"STOPPED":   request.SuccessWaiterState,
		"FAILED":    request.FailureWaiterState,
	}
	state, completed := states[strings.ToUpper(spotinst.StringValue(status))]

	if completed {
		switch state {
		case request.SuccessWaiterState:
			// waiter completed
			return true, nil
		case request.FailureWaiterState:
			// waiter failure state triggered
			return true, fmt.Errorf("ocean: failed waiting for successful state")
		}
	}

	return false, nil
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
