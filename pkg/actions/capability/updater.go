package capability

import (
	"context"
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
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

type UpdaterInterface interface {
	Update(ctx context.Context, capabilities []api.Capability) error
	UpdateTasks(ctx context.Context, capabilities []api.Capability) *tasks.TaskTree
}

type Updater struct {
	clusterName string
	eksAPI      awsapi.EKS
}

func NewUpdater(clusterName string, eksAPI awsapi.EKS) *Updater {
	return &Updater{
		clusterName: clusterName,
		eksAPI:      eksAPI,
	}
}

func (u *Updater) Update(ctx context.Context, capabilities []api.Capability) error {
	taskTree := u.UpdateTasks(ctx, capabilities)
	if errs := taskTree.DoAllSync(); len(errs) > 0 {
		var allErrs []string
		for _, err := range errs {
			allErrs = append(allErrs, err.Error())
		}
		return errors.New(strings.Join(allErrs, "\n"))
	}
	return nil
}

func (u *Updater) UpdateTasks(ctx context.Context, capabilities []api.Capability) *tasks.TaskTree {
	taskTree := &tasks.TaskTree{
		Parallel: true,
	}
	for _, cap := range capabilities {
		taskTree.Append(&tasks.GenericTask{
			Description: fmt.Sprintf("update capability %s", cap.Name),
			Doer: func() error {
				return u.updateCapability(ctx, &cap)
			},
		})
	}
	return taskTree
}

func (u *Updater) updateCapability(ctx context.Context, capability *api.Capability) error {
	input := &awseks.UpdateCapabilityInput{
		ClusterName:    aws.String(u.clusterName),
		CapabilityName: aws.String(capability.Name),
	}

	// Handle role ARN if present
	if capability.RoleARN != "" {
		input.RoleArn = aws.String(capability.RoleARN)
	}

	// Handle delete propagation policy if present
	if capability.DeletePropagationPolicy != "" {
		input.DeletePropagationPolicy = ekstypes.CapabilityDeletePropagationPolicy(capability.DeletePropagationPolicy)
	}

	// Handle configuration if present
	if capability.Configuration != nil {
		config, err := u.convertConfiguration(capability.Configuration)
		if err != nil {
			return err
		}
		input.Configuration = config
	}

	logger.Info("updating capability %s", capability.Name)
	_, err := u.eksAPI.UpdateCapability(ctx, input)
	if err != nil {
		return fmt.Errorf("updating capability %s: %w", capability.Name, err)
	}

	// Wait for capability to become active
	if err := u.waitForCapabilityActive(ctx, capability.Name); err != nil {
		return fmt.Errorf("waiting for capability %s to become active: %w", capability.Name, err)
	}
	logger.Success("capability %s is now active", capability.Name)

	return nil
}

func (u *Updater) convertConfiguration(config *api.CapabilityConfiguration) (*ekstypes.UpdateCapabilityConfiguration, error) {
	if config.ArgoCD == nil {
		return nil, nil
	}

	req := &ekstypes.UpdateCapabilityConfiguration{
		ArgoCd: &ekstypes.UpdateArgoCdConfig{},
	}

	if config.ArgoCD.NetworkAccess != nil && len(config.ArgoCD.NetworkAccess.VPCEIDs) > 0 {
		req.ArgoCd.NetworkAccess = &ekstypes.ArgoCdNetworkAccessConfigRequest{
			VpceIds: config.ArgoCD.NetworkAccess.VPCEIDs,
		}
	}

	if len(config.ArgoCD.RBACRoleMappings) > 0 {
		var roleMappings []ekstypes.ArgoCdRoleMapping
		for _, mapping := range config.ArgoCD.RBACRoleMappings {
			var identities []ekstypes.SsoIdentity
			for _, identity := range mapping.Identities {
				identities = append(identities, ekstypes.SsoIdentity{
					Id:   aws.String(identity.ID),
					Type: ekstypes.SsoIdentityType(identity.Type),
				})
			}
			roleMappings = append(roleMappings, ekstypes.ArgoCdRoleMapping{
				Role:       ekstypes.ArgoCdRole(mapping.Role),
				Identities: identities,
			})
		}
		req.ArgoCd.RbacRoleMappings = &ekstypes.UpdateRoleMappings{
			AddOrUpdateRoleMappings: roleMappings,
		}
	}

	return req, nil
}

func (u *Updater) waitForCapabilityActive(ctx context.Context, capabilityName string) error {
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
			resp, err := u.eksAPI.DescribeCapability(ctx, &awseks.DescribeCapabilityInput{
				ClusterName:    aws.String(u.clusterName),
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
