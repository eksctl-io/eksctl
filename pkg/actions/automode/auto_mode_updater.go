package automode

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/kris-nova/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
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

// A ClusterRoleManager manages the cluster role.
type ClusterRoleManager interface {
	UpdateRoleForAutoMode(ctx context.Context) error
	DeleteAutoModePolicies(ctx context.Context) error
}

// An Updater enables or disables Auto Mode.
type Updater struct {
	RoleManager        RoleManager
	ClusterRoleManager ClusterRoleManager
	PodsGetter         corev1client.PodsGetter
	EKSUpdater         EKSUpdater
	Drainer            NodeGroupDrainer
}

// Update updates the cluster to match the autoModeConfig settings supplied in clusterConfig.
func (u *Updater) Update(ctx context.Context, clusterConfig *api.ClusterConfig, currentCluster *ekstypes.Cluster) error {
	autoModeEnabled := func() bool {
		cc := currentCluster.ComputeConfig
		return cc != nil && *cc.Enabled
	}
	if clusterConfig.IsAutoModeEnabled() {
		amc := clusterConfig.AutoModeConfig
		if autoModeEnabled() {
			if !amc.NodeRoleARN.IsZero() && currentCluster.ComputeConfig.NodeRoleArn != nil &&
				*currentCluster.ComputeConfig.NodeRoleArn != amc.NodeRoleARN.String() {
				return errors.New("autoModeConfig.nodeRoleARN cannot be modified")
			}
			nodePoolsMatch := len(*amc.NodePools) == len(currentCluster.ComputeConfig.NodePools) && (len(*amc.NodePools) == 0 || slices.ContainsFunc(*amc.NodePools, func(np string) bool {
				return slices.Contains(currentCluster.ComputeConfig.NodePools, np)
			}))
			if nodePoolsMatch {
				logger.Info("Auto Mode is already enabled and up-to-date")
				return nil
			}
		} else {
			logger.Info("enabling Auto Mode")
		}
		if err := u.enableAutoMode(ctx, amc, currentCluster.ComputeConfig, clusterConfig.Metadata.Name); err != nil {
			return fmt.Errorf("enabling Auto Mode: %w", err)
		}
		if amc.HasNodePools() {
			logger.Info("cluster subnets will be used for nodes launched by Auto Mode; please create a new NodeClass " +
				"resource if you do not want to use cluster subnets")
		}
		logger.Info("Auto Mode enabled successfully")
		return nil
	}
	if !autoModeEnabled() {
		logger.Info("Auto Mode is already disabled")
		return nil
	}
	if err := u.disableAutoMode(ctx, clusterConfig.Metadata.Name); err != nil {
		return fmt.Errorf("disabling Auto Mode: %w", err)
	}
	logger.Info("Auto Mode disabled successfully")
	return nil
}

func (u *Updater) enableAutoMode(ctx context.Context, autoModeConfig *api.AutoModeConfig, currentClusterCompute *ekstypes.ComputeConfigResponse, clusterName string) error {
	if err := u.preflightCheck(ctx); err != nil {
		return err
	}
	if err := u.ClusterRoleManager.UpdateRoleForAutoMode(ctx); err != nil {
		return fmt.Errorf("updating cluster role to use Auto Mode: %w", err)
	}
	computeConfigReq := &ekstypes.ComputeConfigRequest{
		Enabled:   aws.Bool(true),
		NodePools: *autoModeConfig.NodePools,
	}
	if len(computeConfigReq.NodePools) > 0 {
		if currentClusterCompute != nil && currentClusterCompute.NodeRoleArn != nil {
			computeConfigReq.NodeRoleArn = currentClusterCompute.NodeRoleArn
		} else if autoModeConfig.NodeRoleARN.IsZero() {
			logger.Info("creating node role for Auto Mode")
			nodeRoleARN, err := u.RoleManager.CreateOrImport(ctx, clusterName)
			if err != nil {
				return fmt.Errorf("creating node role to use for Auto Mode nodes: %w", err)
			}
			computeConfigReq.NodeRoleArn = aws.String(nodeRoleARN)
		} else {
			computeConfigReq.NodeRoleArn = aws.String(autoModeConfig.NodeRoleARN.String())
		}
	}
	if err := u.updateClusterConfig(ctx, &eks.UpdateClusterConfigInput{
		Name:          aws.String(clusterName),
		ComputeConfig: computeConfigReq,
		KubernetesNetworkConfig: &ekstypes.KubernetesNetworkConfigRequest{
			ElasticLoadBalancing: &ekstypes.ElasticLoadBalancing{
				Enabled: aws.Bool(true),
			},
		},
		StorageConfig: &ekstypes.StorageConfigRequest{
			BlockStorage: &ekstypes.BlockStorage{
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
		"for a cluster using Auto Mode")
	return nil
}

func (u *Updater) disableAutoMode(ctx context.Context, clusterName string) error {
	logger.Info("disabling Auto Mode")
	if err := u.updateClusterConfig(ctx, &eks.UpdateClusterConfigInput{
		Name: aws.String(clusterName),
		ComputeConfig: &ekstypes.ComputeConfigRequest{
			Enabled: aws.Bool(false),
		},
		KubernetesNetworkConfig: &ekstypes.KubernetesNetworkConfigRequest{
			ElasticLoadBalancing: &ekstypes.ElasticLoadBalancing{
				Enabled: aws.Bool(false),
			},
		},
		StorageConfig: &ekstypes.StorageConfigRequest{
			BlockStorage: &ekstypes.BlockStorage{
				Enabled: aws.Bool(false),
			},
		},
	}); err != nil {
		return err
	}
	if err := u.RoleManager.DeleteIfRequired(ctx); err != nil {
		return fmt.Errorf("deleting IAM resources for Auto Mode: %w", err)
	}
	if err := u.ClusterRoleManager.DeleteAutoModePolicies(ctx); err != nil {
		return fmt.Errorf("deleting Auto Mode policies from cluster role: %w", err)
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
		podList, err := u.PodsGetter.Pods(ns).List(ctx, metav1.ListOptions{
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
