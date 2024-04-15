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

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/eks/waiter"
	"github.com/weaveworks/eksctl/pkg/iam"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	if m.curAuthMode != m.tgAuthMode {
		if m.curAuthMode != ekstypes.AuthenticationModeApiAndConfigMap {
			logger.Info("target authentication mode %v is different than the current authentication mode %v, Updating the cluster authentication mode", m.tgAuthMode, m.curAuthMode)
			err := m.doUpdateAuthenticationMode(ctx, ekstypes.AuthenticationModeApiAndConfigMap, options.Timeout)
			if err != nil {
				return err
			}
		}
		m.curAuthMode = ekstypes.AuthenticationModeApiAndConfigMap
	} else {
		logger.Info("target authentication mode %v is same as current authentication mode %v, not updating the cluster authentication mode", m.tgAuthMode, m.curAuthMode)
	}

	cmEntries, err := m.doGetIAMIdentityMappings()
	if err != nil {
		return err
	}

	curAccessEntries, err := m.doGetAccessEntries(ctx)
	if err != nil {
		return err
	}

	newAccessEntries, skipAPImode, err := doFilterAccessEntries(cmEntries, curAccessEntries)

	if err != nil {
		return err
	}

	newaelen := len(newAccessEntries)

	logger.Info("%d new access entries will be created", newaelen)

	if len(newAccessEntries) != 0 {
		err = m.aeCreator.Create(ctx, newAccessEntries)
		if err != nil {
			return err
		}
	}

	if !skipAPImode {
		if m.curAuthMode != m.tgAuthMode {
			logger.Info("target authentication mode %v is different than the current authentication mode %v, updating the cluster authentication mode", m.tgAuthMode, m.curAuthMode)
			err = m.doUpdateAuthenticationMode(ctx, m.tgAuthMode, options.Timeout)
			if err != nil {
				return err
			}

			err = m.doDeleteIAMIdentityMapping()
			if err != nil {
				return err
			}

			err = doDeleteAWSAuthConfigMap(m.clientSet, "kube-system", "aws-auth")
			if err != nil {
				return err
			}
		}
	}

	return nil

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

	aegetter := NewGetter(m.clusterName, m.eksAPI)
	accessEntries, err := aegetter.Get(ctx, api.ARN{})
	if err != nil {
		return nil, err
	}

	return accessEntries, nil
}

func (m *Migrator) doGetIAMIdentityMappings() ([]iam.Identity, error) {

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
				getRoleOutput, err := m.iamAPI.GetRole(context.Background(), &awsiam.GetRoleInput{RoleName: &match[0]})
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
				getUserOutput, err := m.iamAPI.GetUser(context.Background(), &awsiam.GetUserInput{UserName: &match[0]})
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
	toDoEntries := []api.AccessEntry{}
	uniqueCmEntries := map[string]bool{}

	aeArns := make(map[string]bool)

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
					}
				case iam.ResourceTypeUser:
					if aeEntry := doBuildAccessEntry(cme); aeEntry != nil {
						toDoEntries = append(toDoEntries, *aeEntry)
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

func (m Migrator) doDeleteIAMIdentityMapping() error {
	acm, err := authconfigmap.NewFromClientSet(m.clientSet)
	if err != nil {
		return err
	}

	cmEntries, err := acm.GetIdentities()
	if err != nil {
		return err
	}

	for _, cmEntry := range cmEntries {
		arn := cmEntry.ARN()
		if err := acm.RemoveIdentity(arn, true); err != nil {
			return err
		}
	}
	return acm.Save()
}

func doDeleteAWSAuthConfigMap(clientset kubernetes.Interface, namespace string, name string) error {
	logger.Info("Deleting %q ConfigMap as it is no longer needed in API mode", name)
	err := clientset.CoreV1().ConfigMaps(namespace).Delete(context.Background(), name, metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return nil
}
