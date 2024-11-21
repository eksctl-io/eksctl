package autonomousmode

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/kris-nova/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/autonomousmode"
	"github.com/weaveworks/eksctl/pkg/eks/waiter"
)

const updateTimeout = 30 * time.Minute

// A NodeGroupDrainer drains nodegroups.
type NodeGroupDrainer interface {
	// Drain drains nodegroups.
	Drain(ctx context.Context) error
}

// EKSUpdater updates an EKS cluster.
type EKSUpdater interface {
	// UpdateClusterConfig updates the cluster config.
	UpdateClusterConfig(ctx context.Context, params *eks.UpdateClusterConfigInput, optFns ...func(*eks.Options)) (*eks.UpdateClusterConfigOutput, error)
	// UpdateDescriber describes an update.
	waiter.UpdateDescriber
}

// A RoleManager creates or deletes IAM roles.
type RoleManager interface {
	CreateOrImport(ctx context.Context, clusterName string) (string, error)
	DeleteIfRequired(ctx context.Context) error
}

// An Updater enables or disables Autonomous Mode.
type Updater struct {
	RoleManager     RoleManager
	CoreV1Interface corev1client.CoreV1Interface
	EKSUpdater      EKSUpdater
	Drainer         NodeGroupDrainer
	RBACApplier     *autonomousmode.RBACApplier
}

// Update updates the cluster to match the autonomousModeConfig settings supplied in clusterConfig.
func (u *Updater) Update(ctx context.Context, clusterConfig *api.ClusterConfig, currentCluster *ekstypes.Cluster) error {
	autonomousModeEnabled := func() bool {
		cc := currentCluster.ComputeConfig
		return cc != nil && *cc.Enabled
	}
	if clusterConfig.IsAutonomousModeEnabled() {
		if autonomousModeEnabled() {
			amc := clusterConfig.AutonomousModeConfig
			if !amc.NodeRoleARN.IsZero() && currentCluster.ComputeConfig.NodeRoleArn != nil &&
				*currentCluster.ComputeConfig.NodeRoleArn != amc.NodeRoleARN.String() {
				return errors.New("autonomousModeConfig.nodeRoleARN cannot be modified")
			}
			nodePoolsMatch := len(*amc.NodePools) == len(currentCluster.ComputeConfig.NodePools) && (len(*amc.NodePools) == 0 || slices.ContainsFunc(*amc.NodePools, func(np string) bool {
				return slices.Contains(currentCluster.ComputeConfig.NodePools, np)
			}))
			if nodePoolsMatch {
				logger.Info("Autonomous Mode is already enabled and up-to-date")
				return nil
			}
		} else {
			logger.Info("enabling Autonomous Mode")
		}
		if err := u.enableAutonomousMode(ctx, clusterConfig.AutonomousModeConfig, currentCluster.ComputeConfig, clusterConfig.Metadata.Name); err != nil {
			return fmt.Errorf("enabling Autonomous Mode: %w", err)
		}
		if clusterConfig.AutonomousModeConfig.HasNodePools() {
			logger.Info("cluster subnets will be used for nodes launched by Autonomous Mode; please create a new NodeClass " +
				"resource if you do not want to use cluster subnets")
		}
		logger.Info("applying node RBAC resources for Autonomous Mode")
		if err := u.RBACApplier.ApplyRBACResources(); err != nil {
			return err
		}
		logger.Info("Autonomous Mode enabled successfully")
		return nil
	}
	if !autonomousModeEnabled() {
		logger.Info("Autonomous Mode is already disabled")
		return nil
	}
	if err := u.disableAutonomousMode(ctx, clusterConfig.Metadata.Name); err != nil {
		return fmt.Errorf("disabling Autonomous Mode: %w", err)
	}
	if err := u.RBACApplier.DeleteRBACResources(); err != nil {
		return err
	}
	logger.Info("Autonomous Mode disabled successfully")
	return nil
}

func (u *Updater) enableAutonomousMode(ctx context.Context, autonomousModeConfig *api.AutonomousModeConfig, currentClusterCompute *ekstypes.ComputeConfigResponse, clusterName string) error {
	if err := u.preflightCheck(ctx); err != nil {
		return err
	}
	computeConfigReq := &ekstypes.ComputeConfigRequest{
		Enabled:   aws.Bool(true),
		NodePools: *autonomousModeConfig.NodePools,
	}
	if len(computeConfigReq.NodePools) > 0 {
		if currentClusterCompute != nil && currentClusterCompute.NodeRoleArn != nil {
			computeConfigReq.NodeRoleArn = currentClusterCompute.NodeRoleArn
		} else if autonomousModeConfig.NodeRoleARN.IsZero() {
			logger.Info("creating node role for Autonomous Mode")
			nodeRoleARN, err := u.RoleManager.CreateOrImport(ctx, clusterName)
			if err != nil {
				return fmt.Errorf("creating node role to use for Autonomous Mode nodes: %w", err)
			}
			computeConfigReq.NodeRoleArn = aws.String(nodeRoleARN)
		} else {
			computeConfigReq.NodeRoleArn = aws.String(autonomousModeConfig.NodeRoleARN.String())
		}
	}
	if err := u.updateClusterConfig(ctx, &eks.UpdateClusterConfigInput{
		Name:          aws.String(clusterName),
		ComputeConfig: computeConfigReq,
		KubernetesNetworkConfig: &ekstypes.KubernetesNetworkConfigRequest{
			ElasticLoadBalancing: &ekstypes.ElasticLoadBalancingRequest{
				Enabled: aws.Bool(true),
			},
		},
		StorageConfig: &ekstypes.StorageConfigRequest{
			BlockStorage: &ekstypes.BlockStorageRequest{
				Enabled: aws.Bool(true),
			},
		},
	}); err != nil {
		return err
	}
	if computeConfigReq.NodeRoleArn == nil && currentClusterCompute != nil && currentClusterCompute.NodeRoleArn != nil {
		if err := u.RoleManager.DeleteIfRequired(ctx); err != nil {
			return err
		}
	}
	if u.Drainer != nil {
		if err := u.Drainer.Drain(ctx); err != nil {
			return fmt.Errorf("draining nodegroups: %w", err)
		}
	}
	logger.Info("core networking addons can now be deleted using `eksctl delete addon` as they are not required " +
		"for a cluster using Autonomous Mode")
	return nil
}

func (u *Updater) disableAutonomousMode(ctx context.Context, clusterName string) error {
	logger.Info("disabling Autonomous Mode")
	if err := u.updateClusterConfig(ctx, &eks.UpdateClusterConfigInput{
		Name: aws.String(clusterName),
		ComputeConfig: &ekstypes.ComputeConfigRequest{
			Enabled: aws.Bool(false),
		},
		KubernetesNetworkConfig: &ekstypes.KubernetesNetworkConfigRequest{
			ElasticLoadBalancing: &ekstypes.ElasticLoadBalancingRequest{
				Enabled: aws.Bool(false),
			},
		},
		StorageConfig: &ekstypes.StorageConfigRequest{
			BlockStorage: &ekstypes.BlockStorageRequest{
				Enabled: aws.Bool(false),
			},
		},
	}); err != nil {
		return err
	}
	if err := u.RoleManager.DeleteIfRequired(ctx); err != nil {
		return fmt.Errorf("deleting IAM resources for Autonomous Mode: %w", err)
	}
	return nil
}

func (u *Updater) updateClusterConfig(ctx context.Context, input *eks.UpdateClusterConfigInput) error {
	logger.Info("updating compute config")
	update, err := u.EKSUpdater.UpdateClusterConfig(ctx, input)
	if err != nil {
		return err
	}
	updateWaiter := waiter.NewUpdateWaiter(u.EKSUpdater, func(options *waiter.UpdateWaiterOptions) {
		options.RetryAttemptLogMessage = fmt.Sprintf("waiting for update %q to complete", *update.Update.Id)
	})
	if err := updateWaiter.Wait(ctx, &eks.DescribeUpdateInput{
		Name:     input.Name,
		UpdateId: update.Update.Id,
	}, updateTimeout); err != nil {
		return fmt.Errorf("waiting for cluster update to complete: %w", err)
	}
	return nil
}

func (u *Updater) preflightCheck(ctx context.Context) error {
	if u.Drainer == nil {
		return nil
	}
	knownKarpenterNamespaces := []string{metav1.NamespaceSystem, "karpenter"}
	for _, ns := range knownKarpenterNamespaces {
		podList, err := u.CoreV1Interface.Pods(ns).List(ctx, metav1.ListOptions{
			LabelSelector: "app.kubernetes.io/instance=karpenter",
		})
		if err != nil {
			logger.Warning("error checking for Karpenter pods: %v", err)
			continue
		}
		if len(podList.Items) > 0 {
			return fmt.Errorf("found Karpenter pods in namespace %s; "+
				"either delete Karpenter or scale it down to zero and rerun the command", ns)
		}
	}
	return nil
}
