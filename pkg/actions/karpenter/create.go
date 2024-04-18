package karpenter

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/kris-nova/logger"

	"github.com/aws/aws-sdk-go-v2/aws"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/iam"
	"github.com/weaveworks/eksctl/pkg/karpenter"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
)

// Create creates a Karpenter installer task and waits for it to finish.
func (i *Installer) Create(ctx context.Context) error {
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
	instanceProfileName := fmt.Sprintf("eksctl-%s-%s", builder.KarpenterNodeInstanceProfile, i.Config.Metadata.Name)
	if i.Config.Karpenter.DefaultInstanceProfile != nil {
		instanceProfileName = aws.ToString(i.Config.Karpenter.DefaultInstanceProfile)
	}

	// Create IAM roles
	taskTree := newTasksToInstallKarpenterIAMRoles(ctx, i.Config, i.StackManager, i.CTL.AWSProvider.EC2(), instanceProfileName)
	if err := doTasks(taskTree); err != nil {
		return err
	}

	// Set up service account
	// Because we prefix with eksctl and to avoid having to get the name again,
	// we always pass in the name and overwrite with the service account label.
	roleName := fmt.Sprintf("eksctl-%s-iamservice-role", i.Config.Metadata.Name)
	roleARN := fmt.Sprintf("arn:%s:iam::%s:role/%s", parsedARN.Partition, parsedARN.AccountID, roleName)
	policyArn := fmt.Sprintf("arn:%s:iam::%s:policy/eksctl-%s-%s", parsedARN.Partition, parsedARN.AccountID, builder.KarpenterManagedPolicy, i.Config.Metadata.Name)
	iamServiceAccount := &api.ClusterIAMServiceAccount{
		ClusterIAMMeta: api.ClusterIAMMeta{
			Name:      karpenter.DefaultServiceAccountName,
			Namespace: karpenter.DefaultNamespace,
		},
		AttachPolicyARNs: []string{policyArn},
		RoleName:         roleName,
	}
	if api.IsEnabled(i.Config.Karpenter.CreateServiceAccount) {
		// Create the service account role only.
		iamServiceAccount.RoleOnly = api.Enabled()
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
	identityArn := fmt.Sprintf("arn:%s:iam::%s:role/eksctl-%s-%s", parsedARN.Partition, parsedARN.AccountID, builder.KarpenterNodeRoleName, i.Config.Metadata.Name)
	id, err := iam.NewIdentity(identityArn, authconfigmap.RoleNodeGroupUsername, authconfigmap.RoleNodeGroupGroups)
	if err != nil {
		return fmt.Errorf("failed to create new identity: %w", err)
	}
	if err := acm.AddIdentity(id); err != nil {
		return fmt.Errorf("failed to add new identity: %w", err)
	}
	if err := acm.Save(); err != nil {
		return fmt.Errorf("failed to save the identity config: %w", err)
	}

	// Install Karpenter
	return i.KarpenterInstaller.Install(context.Background(), roleARN, instanceProfileName)
}
