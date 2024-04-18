package accessentry

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	awsiam "github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/iam/types"
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
		logger.Warning("one or more iamidentitymapping could not be migrated to access entry, will not update authentication mode to %v", ekstypes.AuthenticationModeApi)
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
	acm, err := authconfigmap.NewFromClientSet(m.clientSet)
	if err != nil {
		return nil, err
	}

	cmEntries, err := acm.GetIdentities()
	if err != nil {
		return nil, err
	}

	for idx, cme := range cmEntries {
		lastIdx := strings.LastIndex(cme.ARN(), "/")
		cmeName := cme.ARN()[lastIdx+1:]
		var noSuchEntity *types.NoSuchEntityException

		switch cme.Type() {
		case iam.ResourceTypeRole:
			roleCme := iam.RoleIdentity{
				RoleARN: cme.ARN(),
				KubernetesIdentity: iam.KubernetesIdentity{
					KubernetesUsername: cme.Username(),
					KubernetesGroups:   cme.Groups(),
				},
			}

			if cmeName != "" {
				getRoleOutput, err := m.iamAPI.GetRole(ctx, &awsiam.GetRoleInput{RoleName: &cmeName})
				if err != nil {
					if errors.As(err, &noSuchEntity) {
						return nil, fmt.Errorf("role %s does not exists, either delete the iamidentitymapping using \"eksctl delete iamidentitymapping --cluster %s --arn %s\" or create the role in AWS", cmeName, m.clusterName, cme.ARN())
					}
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

			if cmeName != "" {
				getUserOutput, err := m.iamAPI.GetUser(ctx, &awsiam.GetUserInput{UserName: &cmeName})
				if err != nil {
					if errors.As(err, &noSuchEntity) {
						return nil, fmt.Errorf("user \"%s\" does not exists, either delete the iamidentitymapping using \"eksctl delete iamidentitymapping --cluster %s --arn %s\" or create the user in AWS", cmeName, m.clusterName, cme.ARN())
					}
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
	uniqueCmEntries := map[string]struct{}{}
	aeArns := map[string]struct{}{}

	// Create map for current access entry principal ARN
	for _, ae := range accessEntries {
		aeArns[ae.PrincipalARN] = struct{}{}
	}

	for _, cme := range cmEntries {
		if _, ok := uniqueCmEntries[cme.ARN()]; !ok { // Check if cmEntry is not duplicate
			uniqueCmEntries[cme.ARN()] = struct{}{} // Add ARN to cmEntries map

			if _, ok := aeArns[cme.ARN()]; !ok { // Check if the principal ARN is not present in existing access entries
				switch cme.Type() {
				case iam.ResourceTypeRole:
					if strings.Contains(cme.ARN(), ":role/aws-service-role/") { // Check if the principal ARN is service-linked-role
						logger.Warning("found service-linked role iamidentitymapping \"%s\", can not create access entry, skipping", cme.ARN())
						skipAPImode = true
					} else if cme.Username() == authconfigmap.RoleNodeGroupUsername {
						aeEntry := doBuildNodeRoleAccessEntry(cme)
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
					logger.Warning("found account iamidentitymapping \"%s\", can not create access entry", cme.Account())
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
	isLinux := true

	for _, group := range cme.Groups() {
		if group == "eks:kube-proxy-windows" {
			isLinux = false
		}
	}
	// For Linux Nodes
	if isLinux {
		return &api.AccessEntry{
			PrincipalARN: api.MustParseARN(cme.ARN()),
			Type:         "EC2_LINUX",
		}
	}
	// For windows Nodes
	return &api.AccessEntry{
		PrincipalARN: api.MustParseARN(cme.ARN()),
		Type:         "EC2_WINDOWS",
	}
}

func doBuildAccessEntry(cme iam.Identity) *api.AccessEntry {
	containsSys := false

	for _, group := range cme.Groups() {
		if strings.HasPrefix(group, "system:") {
			containsSys = true
			if group == "system:masters" { // Cluster Admin Role
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
		}
	}

	if containsSys { // Check if any GroupName start with "system:"" in name
		logger.Warning("at least one group name associated with %s starts with \"system:\", can not create access entry, skipping", cme.ARN())
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
