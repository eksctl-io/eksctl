package podidentityassociation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/iam"

	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

// AddonMigrator migrates EKS managed addons using IRSAv1 to EKS Pod Identity.
type AddonMigrator struct {
	ClusterName                   string
	AddonServiceAccountRoleMapper AddonServiceAccountRoleMapper
	IAMRoleGetter                 IAMRoleGetter
	StackDescriber                StackDescriber
	EKSAddonsAPI                  EKSAddonsAPI
	RoleMigrator                  RoleMigrator
}

type IAMRoleGetter interface {
	GetRole(ctx context.Context, params *iam.GetRoleInput, optFns ...func(*iam.Options)) (*iam.GetRoleOutput, error)
}

type EKSAddonsAPI interface {
	ListAddons(ctx context.Context, params *eks.ListAddonsInput, optFns ...func(*eks.Options)) (*eks.ListAddonsOutput, error)
	DescribeAddon(ctx context.Context, params *eks.DescribeAddonInput, optFns ...func(*eks.Options)) (*eks.DescribeAddonOutput, error)
	DescribeAddonConfiguration(ctx context.Context, params *eks.DescribeAddonConfigurationInput, optFns ...func(*eks.Options)) (*eks.DescribeAddonConfigurationOutput, error)
	UpdateAddon(ctx context.Context, params *eks.UpdateAddonInput, optFns ...func(*eks.Options)) (*eks.UpdateAddonOutput, error)
}

// A RoleMigrator updates an IAM role to use EKS Pod Identity.
type RoleMigrator interface {
	UpdateTrustPolicyForOwnedRoleTask(ctx context.Context, roleName, serviceAccountName string, stack IRSAv1StackSummary, removeOIDCProviderTrustRelationship bool) tasks.Task
	UpdateTrustPolicyForUnownedRoleTask(ctx context.Context, roleName string, removeOIDCProviderTrustRelationship bool) tasks.Task
}

// Migrate migrates all EKS addons to use EKS Pod Identity.
func (a *AddonMigrator) Migrate(ctx context.Context) (*tasks.TaskTree, error) {
	allTasks := &tasks.TaskTree{
		Parallel: true,
	}
	for serviceAccountRoleARN, addon := range a.AddonServiceAccountRoleMapper {
		taskTree, err := a.migrateAddon(ctx, addon, serviceAccountRoleARN)
		if err != nil {
			return nil, fmt.Errorf("migrating addon %s: %w", *addon.AddonName, err)
		}
		if taskTree != nil {
			allTasks.Append(taskTree)
		}
	}
	return allTasks, nil
}

func (a *AddonMigrator) migrateAddon(ctx context.Context, addon *ekstypes.Addon, serviceAccountRoleARN string) (*tasks.TaskTree, error) {
	if len(addon.PodIdentityAssociations) > 0 {
		logger.Info("addon %s is already using pod identity; skipping migration to pod identity", *addon.AddonName)
		return nil, nil
	}
	addonConfig, err := a.EKSAddonsAPI.DescribeAddonConfiguration(ctx, &eks.DescribeAddonConfigurationInput{
		AddonName:    addon.AddonName,
		AddonVersion: addon.AddonVersion,
	})
	if err != nil {
		return nil, fmt.Errorf("describing pod identity configuration for addon %s: %w", *addon.AddonName, err)
	}
	if len(addonConfig.PodIdentityConfiguration) == 0 {
		logger.Info("addon %s does not support pod identity; skipping migration to pod identity", *addon.AddonName)
		return nil, nil
	}

	logger.Info("will migrate addon %s with serviceAccountRoleARN %q to pod identity; OIDC provider trust relationship will also be removed", *addon.AddonName, *addon.ServiceAccountRoleArn)
	roleName, err := api.RoleNameFromARN(serviceAccountRoleARN)
	if err != nil {
		return nil, fmt.Errorf("parsing role ARN %s: %w", serviceAccountRoleARN, err)
	}
	serviceAccount, err := a.getRoleServiceAccount(ctx, roleName)
	if err != nil {
		return nil, fmt.Errorf("extracting service account from role %s: %w", *addon.ServiceAccountRoleArn, err)
	}
	if serviceAccount == "" {
		if len(addonConfig.PodIdentityConfiguration) != 1 {
			logger.Info("cannot choose a service account as addon %s supports pod identity for multiple service accounts; "+
				"skipping migration to pod identity", *addon.AddonName)
			return nil, nil
		}
		serviceAccount = *addonConfig.PodIdentityConfiguration[0].ServiceAccount
		logger.Info("could not find service account to use for addon %s from existing IAM role; defaulting to %s", *addon.AddonName, serviceAccount)
	}

	var addonTask tasks.TaskTree
	addonTask.IsSubTask = true
	stack, err := a.StackDescriber.DescribeStack(ctx, &manager.Stack{
		StackName: aws.String(manager.MakeAddonStackName(a.ClusterName, *addon.AddonName)),
	})
	if err != nil {
		if manager.IsStackDoesNotExistError(err) {
			return &tasks.TaskTree{
				IsSubTask: true,
				Tasks: []tasks.Task{
					a.RoleMigrator.UpdateTrustPolicyForUnownedRoleTask(ctx, roleName, true),
					a.updateAddonToUsePodIdentity(ctx, addon, serviceAccount),
				},
			}, nil
		}
		return nil, err
	}

	return &tasks.TaskTree{
		IsSubTask: true,
		Tasks: []tasks.Task{
			a.RoleMigrator.UpdateTrustPolicyForOwnedRoleTask(ctx, roleName, serviceAccount, toStackSummary(stack), true),
			a.updateAddonToUsePodIdentity(ctx, addon, serviceAccount),
		},
	}, nil
}

func (a *AddonMigrator) updateAddonToUsePodIdentity(ctx context.Context, addon *ekstypes.Addon, serviceAccount string) tasks.Task {
	return &tasks.GenericTask{
		Description: fmt.Sprintf("migrate addon %s to pod identity", *addon.AddonName),
		Doer: func() error {
			logger.Info("creating a pod identity for addon %s with service account %s", *addon.AddonName, serviceAccount)
			if _, err := a.EKSAddonsAPI.UpdateAddon(ctx, &eks.UpdateAddonInput{
				AddonName:           addon.AddonName,
				AddonVersion:        addon.AddonVersion,
				ClusterName:         addon.ClusterName,
				ConfigurationValues: addon.ConfigurationValues,
				PodIdentityAssociations: []ekstypes.AddonPodIdentityAssociations{
					{
						RoleArn:        addon.ServiceAccountRoleArn,
						ServiceAccount: aws.String(serviceAccount),
					},
				},
			}); err != nil {
				return fmt.Errorf("updating addon %s to use pod identity for service account %s: %w", *addon.AddonName, serviceAccount, err)
			}
			return nil
		},
	}
}

func (a *AddonMigrator) getRoleServiceAccount(ctx context.Context, roleName string) (string, error) {
	role, err := a.IAMRoleGetter.GetRole(ctx, &iam.GetRoleInput{
		RoleName: aws.String(roleName),
	})
	if err != nil {
		return "", err
	}
	assumeRolePolicyDoc, err := url.PathUnescape(*role.Role.AssumeRolePolicyDocument)

	if err != nil {
		return "", fmt.Errorf("unquoting assume role policy document: %w", err)
	}
	var policyDoc api.IAMPolicyDocument
	if err := json.Unmarshal([]byte(assumeRolePolicyDoc), &policyDoc); err != nil {
		return "", fmt.Errorf("parsing assume role policy document: %w", err)
	}
	for _, stmt := range policyDoc.Statements {
		if len(stmt.Condition) == 0 {
			logger.Info("no IAM statements for IRSA found; skipping migration to pod identity")
			return "", nil
		}
		conditions := map[string]map[string]string{}
		if err := json.Unmarshal(stmt.Condition, &conditions); err != nil {
			return "", fmt.Errorf("unmarshaling IAM statement condition: %w", err)
		}
		strEquals, ok := conditions["StringEquals"]
		if !ok {
			continue
		}
		for k, v := range strEquals {
			if strings.HasSuffix(k, ":sub") && strings.HasPrefix(v, "system:serviceaccount:") {
				parts := strings.Split(v, ":")
				if len(parts) != 4 {
					return "", fmt.Errorf("unexpected format %q for service account subject", v)
				}
				return parts[len(parts)-1], nil
			}
		}
	}
	return "", nil
}

func toStackSummary(stack *cfntypes.Stack) IRSAv1StackSummary {
	tags := map[string]string{}
	for _, tag := range stack.Tags {
		tags[*tag.Key] = *tag.Value
	}
	var capabilities []string
	for _, capability := range stack.Capabilities {
		capabilities = append(capabilities, string(capability))
	}
	return IRSAv1StackSummary{
		Name:         *stack.StackName,
		Tags:         tags,
		Capabilities: capabilities,
	}
}
