package podidentityassociation

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
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
	clientSet    kubernetes.Interface
	addonCreator AddonCreator
}

func NewMigrator(
	clusterName string,
	eksAPI awsapi.EKS,
	iamAPI awsapi.IAM,
	stackUpdater StackUpdater,
	clientSet kubernetes.Interface,
	addonCreator AddonCreator,
) *Migrator {
	return &Migrator{
		clusterName:  clusterName,
		eksAPI:       eksAPI,
		iamAPI:       iamAPI,
		stackUpdater: stackUpdater,
		clientSet:    clientSet,
		addonCreator: addonCreator,
	}
}

func (m *Migrator) MigrateToPodIdentity(ctx context.Context, options PodIdentityMigrationOptions) error {
	taskTree := tasks.TaskTree{
		Parallel: false,
		PlanMode: !options.Approve,
	}

	// add task to install the pod identity agent addon
	isInstalled, err := IsPodIdentityAgentInstalled(ctx, m.eksAPI, m.clusterName)
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

	/*
	 add tasks to:
	 update trust policies for IRSAv1 roles
	 AND
	 remove IRSAv1 annotation from service accounts
	*/
	resolver := IRSAv1StackNameResolver{}
	if err := resolver.Populate(func() ([]*api.ClusterIAMServiceAccount, error) {
		return m.stackUpdater.GetIAMServiceAccounts(ctx)
	}); err != nil {
		return err
	}

	serviceAccounts, err := m.clientSet.CoreV1().ServiceAccounts("").List(ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("listing k8s service accounts: %w", err)
	}

	updateTrustPolicyTasks := tasks.TaskTree{
		Parallel:  true,
		IsSubTask: true,
	}
	removeIRSAv1AnnotationTasks := tasks.TaskTree{
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

			// add updateTrustPolicyTasks
			if stackSummary, hasStack := resolver.GetStack(roleARN); hasStack {
				updateTrustPolicyTasks.Append(&updateTrustPolicyForOwnedRole{
					ctx:                                 ctx,
					info:                                fmt.Sprintf("update trust policy for owned role %q", roleName),
					roleName:                            roleName,
					stack:                               stackSummary,
					removeOIDCProviderTrustRelationship: options.RemoveOIDCProviderTrustRelationship,
					iamAPI:                              m.iamAPI,
					stackUpdater:                        m.stackUpdater,
				})

			} else {
				updateTrustPolicyTasks.Append(&updateTrustPolicyForUnownedRole{
					ctx:                                 ctx,
					info:                                fmt.Sprintf("update trust policy for unowned role %q", roleName),
					roleName:                            roleName,
					removeOIDCProviderTrustRelationship: options.RemoveOIDCProviderTrustRelationship,
					iamAPI:                              m.iamAPI,
				})
			}

			// add removeIRSAv1AnnotationTasks
			if !options.RemoveOIDCProviderTrustRelationship {
				continue
			}

			saNameString := sa.Namespace + "/" + sa.Name
			saCopy := &corev1.ServiceAccount{
				ObjectMeta: sa.ObjectMeta,
			}
			removeIRSAv1AnnotationTasks.Append(&tasks.GenericTask{
				Description: fmt.Sprintf("remove iamserviceaccount EKS role annotation for %q", saNameString),
				Doer: func() error {
					delete(saCopy.Annotations, api.AnnotationEKSRoleARN)
					_, err := m.clientSet.CoreV1().ServiceAccounts(saCopy.Namespace).Update(ctx, saCopy, metav1.UpdateOptions{})
					if err != nil {
						return fmt.Errorf("updating serviceaccount %q: %w", saNameString, err)
					}
					logger.Info("removed iamserviceaccount annotation with key %q for %q", api.AnnotationEKSRoleARN, saNameString)
					return nil
				},
			})
		}
	}
	if updateTrustPolicyTasks.Len() == 0 {
		logger.Info("no iamserviceacconts found, there is no need to migrate to pod identity")
		return nil
	}
	taskTree.Append(&updateTrustPolicyTasks)
	if removeIRSAv1AnnotationTasks.Len() > 0 {
		taskTree.Append(&removeIRSAv1AnnotationTasks)
	}

	// add tasks to create pod identity associations
	createAssociationsTasks := NewCreator(m.clusterName, nil, m.eksAPI).CreateTasks(ctx, toBeCreated)
	if createAssociationsTasks.Len() > 0 {
		createAssociationsTasks.IsSubTask = true
		taskTree.Append(createAssociationsTasks)
	}

	// add suggestive logs
	cmdutils.LogIntendedAction(taskTree.PlanMode, "migrate %d iamserviceaccount(s) to pod identity association(s) by executing the following tasks", len(toBeCreated))
	defer cmdutils.LogPlanModeWarning(taskTree.PlanMode)

	return runAllTasks(&taskTree)
}

func IsPodIdentityAgentInstalled(ctx context.Context, eksAPI awsapi.EKS, clusterName string) (bool, error) {
	if _, err := eksAPI.DescribeAddon(ctx, &awseks.DescribeAddonInput{
		AddonName:   aws.String(api.PodIdentityAgentAddon),
		ClusterName: &clusterName,
	}); err != nil {
		var notFoundErr *ekstypes.ResourceNotFoundException
		if errors.As(err, &notFoundErr) {
			return false, nil
		}
		return false, fmt.Errorf("calling %q: %w", fmt.Sprintf("EKS::DescribeAddon::%s", api.PodIdentityAgentAddon), err)
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

type IRSAv1StackNameResolver map[string]IRSAv1StackSummary

type IRSAv1StackSummary struct {
	Name         string
	Tags         map[string]string
	Capabilities []string
}

func (r *IRSAv1StackNameResolver) Populate(
	getIAMServiceAccounts func() ([]*api.ClusterIAMServiceAccount, error),
) error {
	serviceAccounts, err := getIAMServiceAccounts()
	if err != nil {
		return fmt.Errorf("getting iamserviceaccount role stacks: %w", err)
	}
	for _, sa := range serviceAccounts {
		(*r)[*sa.Status.RoleARN] = IRSAv1StackSummary{
			Name:         *sa.Status.StackName,
			Tags:         sa.Status.Tags,
			Capabilities: sa.Status.Capabilities,
		}
	}
	return nil
}

func (r *IRSAv1StackNameResolver) GetStack(roleARN string) (IRSAv1StackSummary, bool) {
	if stack, ok := (*r)[roleARN]; ok {
		return stack, true
	}
	return IRSAv1StackSummary{}, false
}
