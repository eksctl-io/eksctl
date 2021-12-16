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
	clientSetGetter := &kubernetes.CallbackClientSet{
		Callback: func() (kubernetes.Interface, error) {
			return i.ClientSet, nil
		},
	}
	var roleARN string
	iamServiceAccount := &api.ClusterIAMServiceAccount{
		ClusterIAMMeta: api.ClusterIAMMeta{
			Name:      karpenter.DefaultServiceAccountName,
			Namespace: karpenter.DefaultNamespace,
		},
	}
	if api.IsDisabled(i.Config.Karpenter.CreateServiceAccount) {
		policyArn := fmt.Sprintf("arn:aws:iam::%s:policy/%s-%s", parsedARN.AccountID, builder.KarpenterManagedPolicy, i.Config.Metadata.Name)
		iamServiceAccount.AttachPolicyARNs = []string{policyArn}
	} else {
		// Create the service account role only.
		roleName := fmt.Sprintf("eksctl-%s-iamservice-role", i.Config.Metadata.Name)
		roleARN = fmt.Sprintf("arn:aws:iam::%s:role/%s", parsedARN.AccountID, roleName)
		iamServiceAccount.RoleOnly = api.Enabled()
		iamServiceAccount.RoleName = roleName
		iamServiceAccount.AttachPolicy = makePolicyDocument()
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
	identityArn := fmt.Sprintf("arn:aws:iam::%s:role/%s-%s", parsedARN.AccountID, builder.KarpenterNodeRoleName, i.Config.Metadata.Name)
	id, err := iam.NewIdentity(identityArn, authconfigmap.RoleNodeGroupUsername, authconfigmap.RoleNodeGroupGroups)
	if err != nil {
		return fmt.Errorf("failed to create new identity: %w", err)
	}
	if err := acm.AddIdentity(id); err != nil {
		return fmt.Errorf("failed to add new identity: %w", err)
	}
	taskTree := NewTasksToInstallKarpenter(i.Config, i.StackManager, i.CTL.Provider.EC2(), i.KarpenterInstaller, roleARN)
	return doTasks(taskTree)
}

// makePolicyDocument is used in case we don't attach an ARN to the iam service role.
// The arn that would have been attached (KarpenterControllerPolicy) contains these permissions.
// To keep this updated, check required permissions under https://karpenter.sh/docs/getting-started/cloudformation.yaml
// Alternatively, parse out the permissions from that CF template.
func makePolicyDocument() map[string]interface{} {
	return map[string]interface{}{
		"Version": "2012-10-17",
		"Statement": []map[string]interface{}{
			{
				"Effect": "Allow",
				"Action": []string{
					"ec2:CreateLaunchTemplate",
					"ec2:CreateFleet",
					"ec2:RunInstances",
					"ec2:CreateTags",
					"iam:PassRole",
					"ec2:TerminateInstances",
					"ec2:DescribeLaunchTemplates",
					"ec2:DescribeInstances",
					"ec2:DescribeSecurityGroups",
					"ec2:DescribeSubnets",
					"ec2:DescribeInstanceTypes",
					"ec2:DescribeInstanceTypeOfferings",
					"ec2:DescribeAvailabilityZones",
					"ssm:GetParameter",
				},
				"Resource": "*",
			},
		},
	}
}
