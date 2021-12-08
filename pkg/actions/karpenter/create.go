package karpenter

import (
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
	roleArn := fmt.Sprintf("arn:aws:iam::%s:policy/%s-%s", parsedARN.AccountID, builder.KarpenterManagedPolicy, i.Config.Metadata.Name)
	clientSetGetter := &kubernetes.CallbackClientSet{
		Callback: func() (kubernetes.Interface, error) {
			return i.ClientSet, nil
		},
	}
	karpenterServiceAccountTaskTree := i.StackManager.NewTasksToCreateIAMServiceAccounts([]*api.ClusterIAMServiceAccount{
		{
			ClusterIAMMeta: api.ClusterIAMMeta{
				Name:      karpenter.DefaultServiceAccountName,
				Namespace: karpenter.DefaultNamespace,
			},
			AttachRoleARN: roleArn,
		},
	}, nil, clientSetGetter)
	logger.Info(karpenterServiceAccountTaskTree.Describe())
	if err := doTasks(karpenterServiceAccountTaskTree); err != nil {
		return fmt.Errorf("failed to create service account: %w", err)
	}
	// create identity mapping for EC2 nodes to be able to join the cluster.
	acm, err := authconfigmap.NewFromClientSet(i.ClientSet)
	if err != nil {
		return fmt.Errorf("failed to create client for auth config: %w", err)
	}
	identityArn := fmt.Sprintf("arn:aws:iam::%s:role/%s-%s", parsedARN.AccountID, builder.KarpenterNodeRoleName, i.Config.Metadata.Name)
	id, err := iam.NewIdentity(identityArn, authconfigmap.RoleNodeGroupUsername, authconfigmap.RoleNodeGroupGroups)
	if err != nil {
		return fmt.Errorf("failed to create new identity: %w", err)
	}
	if err := acm.AddIdentity(id); err != nil {
		return fmt.Errorf("failed to add new identity: %w", err)
	}
	taskTree := NewTasksToInstallKarpenter(i.Config, i.StackManager, i.CTL.Provider.EC2(), i.KarpenterInstaller)
	return doTasks(taskTree)
}
