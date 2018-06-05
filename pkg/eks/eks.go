package eks

import (
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"

	"github.com/kubicorn/kubicorn/pkg/logger"
)

func (c *ClusterProvider) CreateControlPlane() error {
	input := &eks.CreateClusterInput{
		ClusterName:    &c.cfg.ClusterName,
		RoleArn:        &c.cfg.clusterRoleARN,
		Subnets:        aws.StringSlice(strings.Split(c.cfg.subnetsList, ",")),
		SecurityGroups: aws.StringSlice([]string{c.cfg.securityGroup}),
	}
	output, err := c.svc.eks.CreateCluster(input)
	if err != nil {
		return errors.Wrap(err, "unable to create cluster control plane")
	}
	logger.Debug("output = %#v", output)
	return nil
}

func (c *ClusterProvider) DescribeControlPlane() (*eks.Cluster, error) {
	input := &eks.DescribeClusterInput{
		ClusterName: &c.cfg.ClusterName,
	}
	output, err := c.svc.eks.DescribeCluster(input)
	if err != nil {
		return nil, errors.Wrap(err, "unable to describe cluster control plane")
	}
	return output.Cluster, nil
}

func (c *ClusterProvider) DeleteControlPlane() error {
	cluster, err := c.DescribeControlPlane()
	if err != nil {
		return errors.Wrap(err, "not able to get control plane for deletion")
	}

	input := &eks.DeleteClusterInput{
		ClusterName: cluster.ClusterName,
	}

	if _, err := c.svc.eks.DeleteCluster(input); err != nil {
		return errors.Wrap(err, "unable to delete cluster control plane")
	}
	return nil
}

func (c *ClusterProvider) createControlPlane(errs chan error) error {
	logger.Info("creating control plane %q", c.cfg.ClusterName)

	clusterChan := make(chan eks.Cluster)
	taskErrs := make(chan error)

	if err := c.CreateControlPlane(); err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()
		defer close(taskErrs)
		defer close(clusterChan)
		for {
			select {
			// TODO: https://github.com/weaveworks/eksctl/issues/23
			case <-ticker.C:
				cluster, err := c.DescribeControlPlane()
				if err != nil {
					logger.Warning("continue despite err=%q", err.Error())
					continue
				}
				logger.Debug("cluster = %#v", cluster)
				switch *cluster.Status {
				case eks.ClusterStatusProvisioning:
					continue
				case eks.ClusterStatusActive:
					taskErrs <- nil
					clusterChan <- *cluster
					return
				default:
					taskErrs <- fmt.Errorf("creating control plane: %s", *cluster.Status)
					return
				}
			}
		}
	}()

	go func() {
		defer close(errs)
		if err := <-taskErrs; err != nil {
			errs <- err
			return
		}

		cluster := <-clusterChan

		logger.Debug("created control plane – processing outputs")

		c.cfg.MasterEndpoint = *cluster.MasterEndpoint
		c.cfg.CertificateAuthorityData = []byte(*cluster.CertificateAuthority.Data)

		logger.Debug("clusterConfig = %#v", c.cfg)
		logger.Success("created control plane %q", c.cfg.ClusterName)

		errs <- nil
	}()

	return nil
}

func (c *ClusterProvider) ListClusters() error {
	if c.cfg.ClusterName != "" {
		return c.doListCluster(&c.cfg.ClusterName)
	}

	// TODO: https://github.com/weaveworks/eksctl/issues/27
	input := &eks.ListClustersInput{}
	output, err := c.svc.eks.ListClusters(input)
	if err != nil {
		return errors.Wrap(err, "listing control planes")
	}
	for _, clusterName := range output.Clusters {
		if err := c.doListCluster(clusterName); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClusterProvider) doListCluster(clusterName *string) error {
	input := &eks.DescribeClusterInput{
		ClusterName: clusterName,
	}
	output, err := c.svc.eks.DescribeCluster(input)
	if err != nil {
		return errors.Wrapf(err, "unable to describe control plane %q", *clusterName)
	}
	if *output.Cluster.Status == "Ready" {
		logger.Info("cluster = %#v", *output.Cluster)
		stacks, err := c.ListReadyStacks(fmt.Sprintf("^EKS-%s-.*$", *clusterName))
		if err != nil {
			return errors.Wrapf(err, "listing CloudFormation stack for %q", *clusterName)
		}
		for _, s := range stacks {
			logger.Info("stack = %#v", *s)
		}
	}
	return nil
}

func (c *ClusterProvider) ListAllTaggedResources() error {
	// TODO: https://github.com/weaveworks/eksctl/issues/26
	return nil
}
