package capability

import (
	"context"
	"crypto/sha1"
	"encoding/base32"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

type CreatorInterface interface {
	Create(ctx context.Context, capabilities []api.Capability) error
	CreateTasks(ctx context.Context, capabilities []api.Capability) *tasks.TaskTree
}

type StackCreator interface {
	CreateStack(ctx context.Context, stackName string, resourceSet builder.ResourceSetReader, tags, parameters map[string]string, errs chan error) error
}

type Creator struct {
	clusterName  string
	eksAPI       awsapi.EKS
	stackCreator StackCreator
	cmd          *cmdutils.Cmd
}

func NewCreator(clusterName string, stackCreator StackCreator, eksAPI awsapi.EKS, cmd *cmdutils.Cmd) *Creator {
	return &Creator{
		clusterName:  clusterName,
		stackCreator: stackCreator,
		eksAPI:       eksAPI,
		cmd:          cmd,
	}
}

func (c *Creator) Create(ctx context.Context, capabilities []api.Capability) error {
	taskTree := c.CreateTasks(ctx, capabilities)
	if errs := taskTree.DoAllSync(); len(errs) > 0 {
		var allErrs []string
		for _, err := range errs {
			allErrs = append(allErrs, err.Error())
		}
		return errors.New(strings.Join(allErrs, "\n"))
	}
	return nil
}

func (c *Creator) CreateTasks(ctx context.Context, capabilities []api.Capability) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{
		Parallel: true,
	}
	for _, cap := range capabilities {
		taskTree.Append(&tasks.GenericTask{
			Description: fmt.Sprintf("create and wait for capability %s", cap.Name),
			Doer: func() error {
				// Check cluster is ready before creating capability
				if err := c.ensureClusterReady(ctx); err != nil {
					return fmt.Errorf("cluster not ready for capability creation: %w", err)
				}
				if err := c.createIAMRoleStack(ctx, &cap); err != nil {
					return err
				}

				return c.createCapability(ctx, &cap)
			},
		})
	}
	return taskTree
}

func (c *Creator) ensureClusterReady(ctx context.Context) error {
	clusterProvider, err := c.cmd.NewProviderForExistingCluster(ctx)
	if err != nil {
		return fmt.Errorf("failed to create cluster provider: %w", err)
	}

	cluster := clusterProvider.GetClusterState()
	switch cluster.Status {
	case ekstypes.ClusterStatusCreating, ekstypes.ClusterStatusDeleting, ekstypes.ClusterStatusFailed:
		return fmt.Errorf("cluster is in %s state, cannot create capabilities", cluster.Status)
	default:
		return nil
	}
}

func (c *Creator) createCapability(ctx context.Context, capability *api.Capability) error {
	input := &awseks.CreateCapabilityInput{
		ClusterName:             aws.String(c.clusterName),
		CapabilityName:          aws.String(capability.Name),
		Type:                    ekstypes.CapabilityType(capability.Type),
		DeletePropagationPolicy: ekstypes.CapabilityDeletePropagationPolicyRetain,
		RoleArn:                 aws.String(capability.RoleARN),
	}

	// Override delete propagation policy if specified
	if capability.DeletePropagationPolicy != "" {
		input.DeletePropagationPolicy = ekstypes.CapabilityDeletePropagationPolicy(capability.DeletePropagationPolicy)
	}

	// Handle configuration if present
	if capability.Configuration != nil {
		config, err := c.convertConfiguration(capability.Configuration)
		if err != nil {
			return err
		}
		input.Configuration = config
	}

	// Add tags if present
	if len(capability.Tags) > 0 {
		input.Tags = capability.Tags
	}

	_, err := c.eksAPI.CreateCapability(ctx, input)
	if err != nil {
		return fmt.Errorf("creating capability %s: %w", capability.Name, err)
	}

	// Wait for capability to become active
	if err := c.waitForCapabilityActive(ctx, capability.Name); err != nil {
		return fmt.Errorf("waiting for capability %s to become active: %w", capability.Name, err)
	}
	logger.Success("capability %s is now active", capability.Name)

	// Associate access policies if specified
	if len(capability.AccessPolicies) > 0 {
		logger.Info("associating %d access policies for capability %s", len(capability.AccessPolicies), capability.Name)
		if err := c.waitForAndAssociateAccessPolicies(ctx, capability); err != nil {
			return fmt.Errorf("associating access policies for capability %s: %w", capability.Name, err)
		}
		logger.Info("successfully associated access policies for capability %s", capability.Name)
	}

	return nil
}

func (c *Creator) createIAMRoleStack(ctx context.Context, capability *api.Capability) error {
	rs := builder.NewIAMRoleResourceSetForCapability(capability)
	if err := rs.AddAllResources(); err != nil {
		return err
	}

	stackName := MakeIAMRoleStackName(c.clusterName, capability)
	tags := map[string]string{
		api.ClusterNameTag:    c.clusterName,
		api.CapabilityNameTag: capability.Name,
	}

	for k, v := range capability.Tags {
		tags[k] = v
	}

	stackCh := make(chan error)
	if err := c.stackCreator.CreateStack(ctx, stackName, rs, tags, nil, stackCh); err != nil {
		return fmt.Errorf("creating IAM role stack: %w", err)
	}

	select {
	case err := <-stackCh:
		if err != nil {
			return err
		}
	case <-ctx.Done():
		return fmt.Errorf("timed out waiting for IAM role creation: %w", ctx.Err())
	}

	return nil
}

func (c *Creator) convertConfiguration(config *api.CapabilityConfiguration) (*ekstypes.CapabilityConfigurationRequest, error) {
	if config.ArgoCD == nil {
		return nil, nil
	}

	if config.ArgoCD.AWSIDC == nil {
		return nil, fmt.Errorf("awsIdc configuration is required for ArgoCD capability")
	}

	req := &ekstypes.CapabilityConfigurationRequest{
		ArgoCd: &ekstypes.ArgoCdConfigRequest{},
	}

	if config.ArgoCD.Namespace != "" {
		req.ArgoCd.Namespace = aws.String(config.ArgoCD.Namespace)
	}

	if config.ArgoCD.NetworkAccess != nil && len(config.ArgoCD.NetworkAccess.VPCEIDs) > 0 {
		req.ArgoCd.NetworkAccess = &ekstypes.ArgoCdNetworkAccessConfigRequest{
			VpceIds: config.ArgoCD.NetworkAccess.VPCEIDs,
		}
	}

	if len(config.ArgoCD.RBACRoleMappings) > 0 {
		for _, mapping := range config.ArgoCD.RBACRoleMappings {
			var identities []ekstypes.SsoIdentity
			for _, identity := range mapping.Identities {
				identities = append(identities, ekstypes.SsoIdentity{
					Id:   aws.String(identity.ID),
					Type: ekstypes.SsoIdentityType(identity.Type),
				})
			}
			req.ArgoCd.RbacRoleMappings = append(req.ArgoCd.RbacRoleMappings, ekstypes.ArgoCdRoleMapping{
				Role:       ekstypes.ArgoCdRole(mapping.Role),
				Identities: identities,
			})
		}
	}

	req.ArgoCd.AwsIdc = &ekstypes.ArgoCdAwsIdcConfigRequest{
		IdcInstanceArn: aws.String(config.ArgoCD.AWSIDC.IDCInstanceARN),
	}
	if config.ArgoCD.AWSIDC.IDCRegion != "" {
		req.ArgoCd.AwsIdc.IdcRegion = aws.String(config.ArgoCD.AWSIDC.IDCRegion)
	}

	return req, nil
}

func (c *Creator) waitForCapabilityActive(ctx context.Context, capabilityName string) error {
	timeout := 15 * time.Minute
	interval := 15 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	logger.Info("waiting for capability %s to become active", capabilityName)

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for capability %s to become active", capabilityName)
		case <-ticker.C:
			resp, err := c.eksAPI.DescribeCapability(ctx, &awseks.DescribeCapabilityInput{
				ClusterName:    aws.String(c.clusterName),
				CapabilityName: aws.String(capabilityName),
			})
			if err != nil {
				continue
			}

			if resp.Capability != nil && resp.Capability.Status == ekstypes.CapabilityStatusActive {
				return nil
			}
		}
	}
}

func (c *Creator) waitForAndAssociateAccessPolicies(ctx context.Context, capability *api.Capability) error {
	// Wait for EKS to create the AccessEntry
	accessEntry, err := c.waitForAccessEntryCreation(ctx, capability.RoleARN)
	if err != nil {
		return err
	}

	// Associate AccessPolicies to the AccessEntry
	return c.associateAccessPolicies(ctx, capability, accessEntry)
}

func (c *Creator) waitForAccessEntryCreation(ctx context.Context, capabilityRoleARN string) (*ekstypes.AccessEntry, error) {
	timeout := 5 * time.Minute
	interval := 10 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout waiting for access entry creation for capability role %s", capabilityRoleARN)
		case <-ticker.C:
			resp, err := c.eksAPI.ListAccessEntries(ctx, &awseks.ListAccessEntriesInput{
				ClusterName: aws.String(c.clusterName),
			})
			if err != nil {
				continue
			}

			for _, entryArn := range resp.AccessEntries {
				entry, err := c.eksAPI.DescribeAccessEntry(ctx, &awseks.DescribeAccessEntryInput{
					ClusterName:  aws.String(c.clusterName),
					PrincipalArn: aws.String(entryArn),
				})
				if err != nil {
					continue
				}

				if entry.AccessEntry != nil && *entry.AccessEntry.PrincipalArn == capabilityRoleARN {
					return entry.AccessEntry, nil
				}
			}
		}
	}
}

func (c *Creator) associateAccessPolicies(ctx context.Context, capability *api.Capability, accessEntry *ekstypes.AccessEntry) error {
	for _, policy := range capability.AccessPolicies {
		logger.Info("associating access policy %s for capability %s", policy.PolicyARN, capability.Name)
		_, err := c.eksAPI.AssociateAccessPolicy(ctx, &awseks.AssociateAccessPolicyInput{
			ClusterName:  aws.String(c.clusterName),
			PrincipalArn: accessEntry.PrincipalArn,
			PolicyArn:    aws.String(policy.PolicyARN.String()),
			AccessScope: &ekstypes.AccessScope{
				Type:       ekstypes.AccessScopeType(policy.AccessScope.Type),
				Namespaces: policy.AccessScope.Namespaces,
			},
		})
		if err != nil {
			return fmt.Errorf("associating access policy %s: %w", policy.PolicyARN, err)
		}
	}
	return nil
}

func MakeIAMRoleStackName(clusterName string, capability *api.Capability) string {
	s := sha1.Sum([]byte(capability.Name))
	return fmt.Sprintf("eksctl-%s-capability-role-%s", clusterName, base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(s[:]))
}
