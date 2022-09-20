package cluster

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	"github.com/pkg/errors"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/kris-nova/logger"

	"github.com/weaveworks/eksctl/pkg/actions/nodegroup"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/elb"
	"github.com/weaveworks/eksctl/pkg/fargate"
	"github.com/weaveworks/eksctl/pkg/kubernetes"
	ssh "github.com/weaveworks/eksctl/pkg/ssh/client"
	"github.com/weaveworks/eksctl/pkg/utils/kubeconfig"
)

type NodeGroupDrainer interface {
	Drain(ctx context.Context, input *nodegroup.DrainInput) error
}
type vpcCniDeleter func(clusterConfig *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface)

func deleteSharedResources(ctx context.Context, cfg *api.ClusterConfig, ctl *eks.ClusterProvider, stackManager manager.StackManager, clusterOperable bool, clientSet kubernetes.Interface) error {
	if clusterOperable && !cfg.IsControlPlaneOnOutposts() {
		if err := deleteFargateProfiles(ctx, cfg.Metadata, ctl, stackManager); err != nil {
			return err
		}
	}

	if hasDeprecatedStacks, err := deleteDeprecatedStacks(ctx, stackManager); hasDeprecatedStacks {
		if err != nil {
			return err
		}
		return nil
	}

	ssh.DeleteKeys(ctx, ctl.AWSProvider.EC2(), cfg.Metadata.Name)

	kubeconfig.MaybeDeleteConfig(cfg.Metadata)

	// only need to cleanup ELBs if the cluster has already been created.
	if clusterOperable {
		ctx, cleanup := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cleanup()

		cfg.Metadata.Version = *ctl.Status.ClusterInfo.Cluster.Version

		logger.Info("cleaning up AWS load balancers created by Kubernetes objects of Kind Service or Ingress")
		if err := elb.Cleanup(ctx, ctl.AWSProvider.EC2(), ctl.AWSProvider.ELB(), ctl.AWSProvider.ELBV2(), clientSet, cfg); err != nil {
			return err
		}
	}
	return nil
}

func handleErrors(errs []error, subject string) error {
	logger.Info("%d error(s) occurred while deleting %s", len(errs), subject)
	for _, err := range errs {
		logger.Critical("%s\n", err.Error())
	}
	return fmt.Errorf("failed to delete %s", subject)
}

func deleteFargateProfiles(ctx context.Context, clusterMeta *api.ClusterMeta, ctl *eks.ClusterProvider, stackManager manager.StackManager) error {
	manager := fargate.NewFromProvider(
		clusterMeta.Name,
		ctl.AWSProvider,
		stackManager,
	)
	profileNames, err := manager.ListProfiles(ctx)
	if err != nil {
		if fargate.IsUnauthorizedError(err) {
			logger.Debug("Fargate: unauthorized error: %v", err)
			logger.Info("either account is not authorized to use Fargate or region %s is not supported. Ignoring error",
				clusterMeta.Region)
			return nil
		}
		return err
	}

	// Linearise the deleting of Fargate profiles by passing as the API
	// otherwise errors out with:
	//   ResourceInUseException: Cannot delete Fargate Profile ${name2} because
	//   cluster ${clusterName} currently has Fargate profile ${name1} in
	//   status DELETING

	for _, profileName := range profileNames {
		logger.Info("deleting Fargate profile %q", profileName)
		// All Fargate profiles must be completely deleted by waiting for the deletion to complete, before deleting
		// the cluster itself, otherwise it can result in this error:
		//   Cannot delete because cluster <cluster> currently has Fargate profile <profile> in status DELETING
		if err := manager.DeleteProfile(ctx, profileName, true); err != nil {
			return err
		}
		logger.Info("deleted Fargate profile %q", profileName)
	}
	logger.Info("deleted %v Fargate profile(s)", len(profileNames))

	stack, err := stackManager.GetFargateStack(ctx)
	if err != nil {
		return err
	}

	if stack != nil {
		_, err := stackManager.DeleteStackBySpec(ctx, stack)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteDeprecatedStacks(ctx context.Context, stackManager manager.StackManager) (bool, error) {
	tasks, err := stackManager.DeleteTasksForDeprecatedStacks(ctx)
	if err != nil {
		return true, err
	}
	if count := tasks.Len(); count > 0 {
		logger.Info(tasks.Describe())
		if errs := tasks.DoAllSync(); len(errs) > 0 {
			return true, handleErrors(errs, "deprecated stacks")
		}
		logger.Success("deleted all %s deprecated stacks", count)
		return true, nil
	}
	return false, nil
}

func checkForUndeletedStacks(ctx context.Context, stackManager manager.StackManager) error {
	stacks, err := stackManager.DescribeStacks(ctx)
	if err != nil {
		return err
	}

	var undeletedStacks []string

	for _, stack := range stacks {
		if stack.StackStatus == types.StackStatusDeleteInProgress {
			continue
		}

		undeletedStacks = append(undeletedStacks, *stack.StackName)
	}

	if len(undeletedStacks) > 0 {
		logger.Warning("found the following undeleted stacks: %s", strings.Join(undeletedStacks, ","))
		return errors.New("failed to delete all resources")
	}

	return nil
}

func drainAllNodeGroups(ctx context.Context, cfg *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface, allStacks []manager.NodeGroupStack,
	disableEviction bool, parallel int, nodeGroupDrainer NodeGroupDrainer, vpcCniDeleter vpcCniDeleter, podEvictionWaitPeriod time.Duration) error {
	if len(allStacks) == 0 {
		return nil
	}

	cfg.NodeGroups = []*api.NodeGroup{}
	for _, s := range allStacks {
		if s.Type == api.NodeGroupTypeUnmanaged {
			cmdutils.PopulateUnmanagedNodegroup(s.NodeGroupName, cfg)
		}
	}

	logger.Info("will drain %d unmanaged nodegroup(s) in cluster %q", len(cfg.NodeGroups), cfg.Metadata.Name)

	drainInput := &nodegroup.DrainInput{
		NodeGroups:            cmdutils.ToKubeNodeGroups(cfg),
		MaxGracePeriod:        ctl.AWSProvider.WaitTimeout(),
		DisableEviction:       disableEviction,
		PodEvictionWaitPeriod: podEvictionWaitPeriod,
		Parallel:              parallel,
	}
	if err := nodeGroupDrainer.Drain(ctx, drainInput); err != nil {
		return err
	}

	vpcCniDeleter(cfg, ctl, clientSet)
	return nil
}

// Attempts to delete the vpc-cni, and fails silently if an error occurs. This is an attempt
// to prevent a race condition in the vpc-cni #1849
func attemptVpcCniDeletion(ctx context.Context, clusterConfig *api.ClusterConfig, ctl *eks.ClusterProvider, clientSet kubernetes.Interface) {
	if !clusterConfig.IsControlPlaneOnOutposts() {
		const vpcCNI = "vpc-cni"
		logger.Debug("deleting EKS addon %q if it exists", vpcCNI)
		_, err := ctl.AWSProvider.EKS().DeleteAddon(ctx, &awseks.DeleteAddonInput{
			ClusterName: aws.String(clusterConfig.Metadata.Name),
			AddonName:   aws.String(vpcCNI),
		})

		if err != nil {
			var notFoundErr *ekstypes.ResourceNotFoundException
			if errors.As(err, &notFoundErr) {
				logger.Debug("EKS addon %q does not exist", vpcCNI)
			} else {
				logger.Debug("failed to delete addon %q: %v", vpcCNI, err)
			}
		}
	}

	logger.Debug("deleting kube-system/aws-node DaemonSet")
	err := clientSet.AppsV1().DaemonSets("kube-system").Delete(ctx, "aws-node", metav1.DeleteOptions{})
	if err != nil {
		logger.Debug("failed to delete kube-system/aws-node DaemonSet: %w", err)
	}
}
