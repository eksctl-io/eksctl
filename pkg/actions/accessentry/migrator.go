package accessentry

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	awsiam "github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/kris-nova/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/eks/waiter"
	"github.com/weaveworks/eksctl/pkg/iam"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

type MigrationOptions struct {
	TargetAuthMode string
	Approve        bool
	Timeout        time.Duration
}

type Migrator struct {
	clusterName string
	eksAPI      awsapi.EKS
	iamAPI      awsapi.IAM
	clientSet   kubernetes.Interface
	aeCreator   Creator
	curAuthMode ekstypes.AuthenticationMode
	tgAuthMode  ekstypes.AuthenticationMode
}

func NewMigrator(
	clusterName string,
	eksAPI awsapi.EKS,
	iamAPI awsapi.IAM,
	clientSet kubernetes.Interface,
	aeCreator Creator,
	curAuthMode ekstypes.AuthenticationMode,
	tgAuthMode ekstypes.AuthenticationMode,
) *Migrator {
	return &Migrator{
		clusterName: clusterName,
		eksAPI:      eksAPI,
		iamAPI:      iamAPI,
		clientSet:   clientSet,
		aeCreator:   aeCreator,
		curAuthMode: curAuthMode,
		tgAuthMode:  tgAuthMode,
	}
}

func (m *Migrator) MigrateToAccessEntry(ctx context.Context, options MigrationOptions) error {
	taskTree := tasks.TaskTree{
		Parallel: false,
		PlanMode: !options.Approve,
	}

	if m.curAuthMode != m.tgAuthMode {
		taskTree.Append(&tasks.GenericTask{
			Description: fmt.Sprintf("update authentication mode to %v", ekstypes.AuthenticationModeApiAndConfigMap),
			Doer: func() error {
				if m.curAuthMode != ekstypes.AuthenticationModeApiAndConfigMap {
					logger.Info("target authentication mode %v is different than the current authentication mode %v, Updating the cluster authentication mode", m.tgAuthMode, m.curAuthMode)
					return m.doUpdateAuthenticationMode(ctx, ekstypes.AuthenticationModeApiAndConfigMap, options.Timeout)
				}
				m.curAuthMode = ekstypes.AuthenticationModeApiAndConfigMap
				return nil
			},
		})
	} else {
		logger.Info("target authentication mode %v is same as current authentication mode %v, not updating the cluster authentication mode", m.tgAuthMode, m.curAuthMode)
	}

	cmEntries, err := m.doGetIAMIdentityMappings(ctx)
	if err != nil {
		return err
	}

	curAccessEntries, err := m.doGetAccessEntries(ctx)
	if err != nil && m.curAuthMode != ekstypes.AuthenticationModeConfigMap {
		return err
	}

	newAccessEntries, skipAPImode, err := doFilterAccessEntries(cmEntries, curAccessEntries)
	if err != nil {
		return err
	}

	if newaelen := len(newAccessEntries); newaelen != 0 {
		logger.Info("%d new access entries will be created", newaelen)
		aeTasks := m.aeCreator.CreateTasks(ctx, newAccessEntries)
		aeTasks.IsSubTask = true
		taskTree.Append(aeTasks)
	}

	if !skipAPImode {
		if m.tgAuthMode == ekstypes.AuthenticationModeApi {
			taskTree.Append(&tasks.GenericTask{
				Description: fmt.Sprintf("update authentication mode to %v", ekstypes.AuthenticationModeApi),
				Doer: func() error {
					logger.Info("target authentication mode %v is different than the current authentication mode %v, updating the cluster authentication mode", m.tgAuthMode, m.curAuthMode)
					return m.doUpdateAuthenticationMode(ctx, m.tgAuthMode, options.Timeout)
				},
			})

			taskTree.Append(&tasks.GenericTask{
				Description: fmt.Sprintf("delete aws-auth configMap when authentication mode is %v", ekstypes.AuthenticationModeApi),
				Doer: func() error {
					return doDeleteAWSAuthConfigMap(ctx, m.clientSet, authconfigmap.ObjectNamespace, authconfigmap.ObjectName)
				},
			})
		}
	} else if m.tgAuthMode == ekstypes.AuthenticationModeApi {
		logger.Warning("one or more identitymapping could not be migrated to access entry, will not update authentication mode to %v", ekstypes.AuthenticationModeApi)
	}

	return runAllTasks(&taskTree)
}

func (m *Migrator) doUpdateAuthenticationMode(ctx context.Context, authMode ekstypes.AuthenticationMode, timeout time.Duration) error {
	output, err := m.eksAPI.UpdateClusterConfig(ctx, &awseks.UpdateClusterConfigInput{
		Name: aws.String(m.clusterName),
		AccessConfig: &ekstypes.UpdateAccessConfigRequest{
			AuthenticationMode: authMode,
		},
	})
	if err != nil {
		return fmt.Errorf("failed to update cluster config: %v", err)
	}

	updateWaiter := waiter.NewUpdateWaiter(m.eksAPI, func(options *waiter.UpdateWaiterOptions) {
		options.RetryAttemptLogMessage = fmt.Sprintf("waiting for update %q in cluster %q to complete", *output.Update.Id, m.clusterName)
	})
	err = updateWaiter.Wait(ctx, &awseks.DescribeUpdateInput{
		Name:     &m.clusterName,
		UpdateId: output.Update.Id,
	}, timeout)

	switch e := err.(type) {
	case *waiter.UpdateFailedError:
		if e.Status == string(ekstypes.UpdateStatusCancelled) {
			return fmt.Errorf("request to update cluster authentication mode was cancelled: %s", e.UpdateError)
		}
		return fmt.Errorf("failed to update cluster authentication mode: %s", e.UpdateError)
	case nil:
		logger.Info("authentication mode was successfully updated to %s on cluster %s", authMode, m.clusterName)
		m.curAuthMode = authMode
		return nil
	default:
		return err
	}
}

func (m *Migrator) doGetAccessEntries(ctx context.Context) ([]Summary, error) {
	aeGetter := NewGetter(m.clusterName, m.eksAPI)
	return aeGetter.Get(ctx, api.ARN{})
}

func (m *Migrator) doGetIAMIdentityMappings(ctx context.Context) ([]iam.Identity, error) {

	nameRegex := regexp.MustCompile(`[^/]+$`)

	acm, err := authconfigmap.NewFromClientSet(m.clientSet)
	if err != nil {
		return nil, err
	}

	cmEntries, err := acm.GetIdentities()
	if err != nil {
		return nil, err
	}

	for idx, cme := range cmEntries {
		switch cme.Type() {
		case iam.ResourceTypeRole:
			roleCme := iam.RoleIdentity{
				RoleARN: cme.ARN(),
				KubernetesIdentity: iam.KubernetesIdentity{
					KubernetesUsername: cme.Username(),
					KubernetesGroups:   cme.Groups(),
				},
			}

			if match := nameRegex.FindStringSubmatch(roleCme.RoleARN); match != nil {
				getRoleOutput, err := m.iamAPI.GetRole(ctx, &awsiam.GetRoleInput{RoleName: &match[0]})
				if err != nil {
					return nil, err
				}

				roleCme.RoleARN = *getRoleOutput.Role.Arn
			}

			cmEntries[idx] = iam.Identity(roleCme)

		case iam.ResourceTypeUser:
			userCme := iam.UserIdentity{
				UserARN: cme.ARN(),
				KubernetesIdentity: iam.KubernetesIdentity{
					KubernetesUsername: cme.Username(),
					KubernetesGroups:   cme.Groups(),
				},
			}

			if match := nameRegex.FindStringSubmatch(userCme.UserARN); match != nil {
				getUserOutput, err := m.iamAPI.GetUser(ctx, &awsiam.GetUserInput{UserName: &match[0]})
				if err != nil {
					return nil, err
				}

				userCme.UserARN = *getUserOutput.User.Arn
			}
			cmEntries[idx] = iam.Identity(userCme)
		}
	}

	return cmEntries, nil
}

func doFilterAccessEntries(cmEntries []iam.Identity, accessEntries []Summary) ([]api.AccessEntry, bool, error) {

	skipAPImode := false
	var toDoEntries []api.AccessEntry
	uniqueCmEntries := map[string]bool{}
	aeArns := map[string]bool{}

	// Create ARN Map for current access entries
	for _, ae := range accessEntries {
		aeArns[ae.PrincipalARN] = true
	}

	for _, cme := range cmEntries {
		if !uniqueCmEntries[cme.ARN()] { // Check if cmEntry is not duplicate
			if !aeArns[cme.ARN()] { // Check if the ARN is not in existing access entries
				switch cme.Type() {
				case iam.ResourceTypeRole:
					if aeEntry := doBuildNodeRoleAccessEntry(cme); aeEntry != nil {
						toDoEntries = append(toDoEntries, *aeEntry)
					} else if aeEntry := doBuildAccessEntry(cme); aeEntry != nil {
						toDoEntries = append(toDoEntries, *aeEntry)
					} else {
						skipAPImode = true
					}
				case iam.ResourceTypeUser:
					if aeEntry := doBuildAccessEntry(cme); aeEntry != nil {
						toDoEntries = append(toDoEntries, *aeEntry)
					} else {
						skipAPImode = true
					}
				case iam.ResourceTypeAccount:
					logger.Warning("found account mapping %s, can not create access entry for account mapping, skipping", cme.Account())
					skipAPImode = true
				}
			} else {
				logger.Warning("%s already exists in access entry, skipping", cme.ARN())
			}
		}
	}

	return toDoEntries, skipAPImode, nil
}

func doBuildNodeRoleAccessEntry(cme iam.Identity) *api.AccessEntry {

	groupsStr := strings.Join(cme.Groups(), ",")

	if strings.Contains(groupsStr, "system:nodes") && !strings.Contains(groupsStr, "eks:kube-proxy-windows") { // For Windows Nodes
		return &api.AccessEntry{
			PrincipalARN: api.MustParseARN(cme.ARN()),
			Type:         "EC2_LINUX",
		}
	}

	if strings.Contains(groupsStr, "system:nodes") && strings.Contains(groupsStr, "eks:kube-proxy-windows") { // For Linux Nodes
		return &api.AccessEntry{
			PrincipalARN: api.MustParseARN(cme.ARN()),
			Type:         "EC2_WINDOWS",
		}
	}

	return nil
}

func doBuildAccessEntry(cme iam.Identity) *api.AccessEntry {

	groupsStr := strings.Join(cme.Groups(), ",")

	if strings.Contains(groupsStr, "system:masters") { // Admin Role
		return &api.AccessEntry{
			PrincipalARN: api.MustParseARN(cme.ARN()),
			Type:         "STANDARD",
			AccessPolicies: []api.AccessPolicy{
				{
					PolicyARN: api.MustParseARN("arn:aws:eks::aws:cluster-access-policy/AmazonEKSClusterAdminPolicy"),
					AccessScope: api.AccessScope{
						Type: ekstypes.AccessScopeTypeCluster,
					},
				},
			},
			KubernetesUsername: cme.Username(),
		}
	}

	if strings.Contains(groupsStr, "system") { // Admin Role
		logger.Warning("at least one group name associated with %s starts with \"system\", can not create access entry with such group name, skipping", cme.ARN())
		return nil
	}

	return &api.AccessEntry{
		PrincipalARN:       api.MustParseARN(cme.ARN()),
		Type:               "STANDARD",
		KubernetesGroups:   cme.Groups(),
		KubernetesUsername: cme.Username(),
	}

}

func doDeleteAWSAuthConfigMap(ctx context.Context, clientset kubernetes.Interface, namespace, name string) error {
	logger.Info("deleting %q ConfigMap as it is no longer needed in API mode", name)
	return clientset.CoreV1().ConfigMaps(namespace).Delete(ctx, name, metav1.DeleteOptions{})

}
