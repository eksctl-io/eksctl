package karpenter

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/iam"
	"github.com/weaveworks/eksctl/pkg/karpenter"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

// Create creates a Karpenter installer task and waits for it to finish.
func (i *Installer) Create() error {
	// create the needed service account before Karpenter, otherwise, Karpenter will fail to be created.
	parsedARN, err := arn.Parse(i.Config.Status.ARN)
	if err != nil {
		return fmt.Errorf("unexpected or invalid ARN: %q, %w", i.Config.Status.ARN, err)
	}
	clientSetGetter := &kubernetes.CallbackClientSet{
		Callback: func() (kubernetes.Interface, error) {
			return i.ClientSet, nil
		},
	}
	// Create IAM roles
	taskTree := newTasksToInstallKarpenterIAMRoles(i.Config, i.StackManager, i.CTL.Provider.EC2())
	if err := doTasks(taskTree); err != nil {
		return err
	}

	// Set up service account
	var roleARN string
	policyArn := fmt.Sprintf("arn:aws:iam::%s:policy/eksctl-%s-%s", parsedARN.AccountID, builder.KarpenterManagedPolicy, i.Config.Metadata.Name)
	iamServiceAccount := &api.ClusterIAMServiceAccount{
		ClusterIAMMeta: api.ClusterIAMMeta{
			Name:      karpenter.DefaultServiceAccountName,
			Namespace: karpenter.DefaultNamespace,
		},
		AttachPolicyARNs: []string{policyArn},
	}
	if api.IsEnabled(i.Config.Karpenter.CreateServiceAccount) {
		// Create the service account role only.
		roleName := fmt.Sprintf("eksctl-%s-iamservice-role", i.Config.Metadata.Name)
		roleARN = fmt.Sprintf("arn:aws:iam::%s:role/%s", parsedARN.AccountID, roleName)
		iamServiceAccount.RoleOnly = api.Enabled()
		iamServiceAccount.RoleName = roleName
	}
	karpenterServiceAccountTaskTree := i.StackManager.NewTasksToCreateIAMServiceAccounts([]*api.ClusterIAMServiceAccount{iamServiceAccount}, i.OIDC, clientSetGetter)
	logger.Info(karpenterServiceAccountTaskTree.Describe())
	if err := doTasks(karpenterServiceAccountTaskTree); err != nil {
		return fmt.Errorf("failed to create/attach service account: %w", err)
	}

	// create identity mapping for EC2 nodes to be able to join the cluster.
	acm, err := authconfigmap.NewFromClientSet(i.ClientSet)
	if err != nil {
		return fmt.Errorf("failed to create client for auth config: %w", err)
	}
	identityArn := fmt.Sprintf("arn:aws:iam::%s:role/eksctl-%s-%s", parsedARN.AccountID, builder.KarpenterNodeRoleName, i.Config.Metadata.Name)
	id, err := iam.NewIdentity(identityArn, authconfigmap.RoleNodeGroupUsername, authconfigmap.RoleNodeGroupGroups)
	if err != nil {
		return fmt.Errorf("failed to create new identity: %w", err)
	}
	if err := acm.AddIdentity(id); err != nil {
		return fmt.Errorf("failed to add new identity: %w", err)
	}

	// Install Karpenter
	return i.KarpenterInstaller.Install(context.Background(), roleARN)
}
