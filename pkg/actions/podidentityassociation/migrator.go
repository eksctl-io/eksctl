package podidentityassociation

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	iamoidc "github.com/weaveworks/eksctl/pkg/iam/oidc"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

type AddonCreator interface {
	Create(ctx context.Context, addon *api.Addon, waitTimeout time.Duration) error
}

type PodIdentityMigrationOptions struct {
	RemoveOIDCProviderTrustRelationship bool
	// SkipAgentInstallation               bool
	Approve bool
	Timeout time.Duration
}

type Migrator struct {
	clusterName  string
	eksAPI       awsapi.EKS
	iamAPI       awsapi.IAM
	stackUpdater StackUpdater
	oidcManager  *iamoidc.OpenIDConnectManager
	clientSet    kubernetes.Interface
	addonCreator AddonCreator
}

func NewMigrator(
	clusterName string,
	eksAPI awsapi.EKS,
	iamAPI awsapi.IAM,
	stackUpdater StackUpdater,
	oidcManager *iamoidc.OpenIDConnectManager,
	clientSet kubernetes.Interface,
	addonCreator AddonCreator,
) *Migrator {
	return &Migrator{
		clusterName:  clusterName,
		eksAPI:       eksAPI,
		iamAPI:       iamAPI,
		stackUpdater: stackUpdater,
		oidcManager:  oidcManager,
		clientSet:    clientSet,
		addonCreator: addonCreator,
	}
}

func (m *Migrator) MigrateToPodIdentity(ctx context.Context, options PodIdentityMigrationOptions) error {
	taskTree := tasks.TaskTree{
		Parallel: false,
		PlanMode: !options.Approve,
	}
	defer cmdutils.LogPlanModeWarning(taskTree.PlanMode)

	// add task to install the pod identity agent addon
	isInstalled, err := isPodIdentityAgentInstalled(ctx, m.eksAPI, m.clusterName)
	if err != nil {
		return err
	}
	if !isInstalled {
		taskTree.Append(&tasks.GenericTask{
			Description: fmt.Sprintf("install %s addon", api.PodIdentityAgentAddon),
			Doer: func() error {
				return m.addonCreator.Create(ctx, &api.Addon{Name: api.PodIdentityAgentAddon}, options.Timeout)
			},
		})
	}

	// add tasks to update trust policies for IRSAv1 roles
	roleToStackName := map[string]string{}
	iamRoleStacks, err := m.stackUpdater.GetIAMServiceAccounts(ctx)
	if err != nil {
		return fmt.Errorf("getting IAM Role stacks: %w", err)
	}
	for _, s := range iamRoleStacks {
		roleToStackName[*s.Status.RoleARN] = *s.Status.StackName
	}

	serviceAccounts, err := m.clientSet.CoreV1().ServiceAccounts("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("listing k8s service accounts: %w", err)
	}

	updateTrustPolicyTasks := tasks.TaskTree{
		Parallel:  true,
		IsSubTask: true,
	}
	toBeCreated := []api.PodIdentityAssociation{}
	for _, sa := range serviceAccounts.Items {
		if roleARN, ok := sa.Annotations[api.AnnotationEKSRoleARN]; ok {
			// collect pod identity associations that need to be created
			toBeCreated = append(toBeCreated, api.PodIdentityAssociation{
				ServiceAccountName: sa.Name,
				Namespace:          sa.Namespace,
				RoleARN:            roleARN,
			})

			// infer role name to use in IAM API inputs
			roleName, err := getNameFromARN(roleARN)
			if err != nil {
				return err
			}

			// add tasks
			if stackName, ok := roleToStackName[roleARN]; ok {
				updateTrustPolicyTasks.Append(&updateTrustPolicyForOwnedRole{
					ctx:                                 ctx,
					info:                                fmt.Sprintf("update trust policy for role %s", roleName),
					roleName:                            roleName,
					stackName:                           stackName,
					removeOIDCProviderTrustRelationship: options.RemoveOIDCProviderTrustRelationship,
					iamAPI:                              m.iamAPI,
					stackUpdater:                        m.stackUpdater,
				})
				continue
			}
			updateTrustPolicyTasks.Append(&updateTrustPolicyForUnownedRole{
				ctx:                                 ctx,
				info:                                fmt.Sprintf("update trust policy for role %s", roleName),
				roleName:                            roleName,
				removeOIDCProviderTrustRelationship: options.RemoveOIDCProviderTrustRelationship,
				iamAPI:                              m.iamAPI,
			})
		}
	}
	taskTree.Append(&updateTrustPolicyTasks)

	// add tasks to create pod identity associations
	createAssociationsTasks := NewCreator(m.clusterName, nil, m.eksAPI).CreateTasks(ctx, toBeCreated)
	if createAssociationsTasks.Len() > 0 {
		createAssociationsTasks.IsSubTask = true
		taskTree.Append(createAssociationsTasks)
	}

	return runAllTasks(&taskTree)
}

func isPodIdentityAgentInstalled(ctx context.Context, eksAPI awsapi.EKS, clusterName string) (bool, error) {
	if _, err := eksAPI.DescribeAddon(ctx, &awseks.DescribeAddonInput{
		AddonName:   aws.String(api.PodIdentityAgentAddon),
		ClusterName: &clusterName,
	}); err != nil {
		var notFoundErr *ekstypes.ResourceNotFoundException
		if errors.As(err, &notFoundErr) {
			return false, nil
		}
		return false, fmt.Errorf("error calling `EKS::DescribeAddon::%s`: %v", api.PodIdentityAgentAddon, err)
	}
	return true, nil
}

func getNameFromARN(roleARN string) (string, error) {
	parts := strings.Split(roleARN, "/")
	if len(parts) != 2 {
		return "", fmt.Errorf("cannot parse role name from roleARN: %s", roleARN)
	}
	return parts[1], nil
}
