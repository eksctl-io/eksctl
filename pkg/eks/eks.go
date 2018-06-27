package eks

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/eks"

	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/weaveworks/eksctl/pkg/utils"
)

func (c *ClusterProvider) CreateControlPlane() error {
	input := &eks.CreateClusterInput{
		Name:    &c.cfg.ClusterName,
		RoleArn: &c.cfg.clusterRoleARN,
		ResourcesVpcConfig: &eks.VpcConfigRequest{
			SubnetIds:        aws.StringSlice(strings.Split(c.cfg.subnetsList, ",")),
			SecurityGroupIds: aws.StringSlice([]string{c.cfg.securityGroup}),
		},
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
		Name: &c.cfg.ClusterName,
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
		Name: cluster.Name,
	}

	if _, err := c.svc.eks.DeleteCluster(input); err != nil {
		return errors.Wrap(err, "unable to delete cluster control plane")
	}
	return nil
}

func (c *ClusterProvider) createControlPlane() <-chan error {
	errs := make(chan error)

	go func() {
		defer close(errs)
		logger.Info("creating control plane %q", c.cfg.ClusterName)

		if err := c.CreateControlPlane(); err != nil {
			if utils.HasAwsErrorCode(err, eks.ErrCodeResourceInUseException) {
				logger.Info("using existing EKS Cluster stack %q", c.cfg.ClusterName)
			} else {
				errs <- err
				return
			}
		}

		ticker := utils.NewInstantTicker(20 * time.Second)
		defer ticker.Stop()

		timer := time.NewTimer(time.Duration(c.cfg.AWSOperationTimeoutSeconds) * time.Second)
		defer timer.Stop()

		for {
			select {
			case <-timer.C:
				errs <- fmt.Errorf("timed out creating control plane %q after %d seconds", c.cfg.ClusterName, c.cfg.AWSOperationTimeoutSeconds)
				return

			case <-ticker.C:
				cluster, err := c.DescribeControlPlane()
				if err != nil {
					logger.Warning("continue despite err=%q", err.Error())
					continue
				}
				logger.Debug("cluster = %#v", cluster)
				switch *cluster.Status {
				case eks.ClusterStatusCreating:
					continue
				case eks.ClusterStatusActive:
					logger.Debug("created control plane – processing outputs")

					if err := c.GetCredentials(*cluster); err != nil {
						errs <- err
					}

					logger.Debug("clusterConfig = %#v", c.cfg)
					logger.Success("created control plane %q", c.cfg.ClusterName)
					return
				default:
					errs <- fmt.Errorf("creating control plane: %s", *cluster.Status)
					return
				}
			}
		}
	}()

	return errs
}

func (c *ClusterProvider) GetCredentials(cluster eks.Cluster) error {
	c.cfg.MasterEndpoint = *cluster.Endpoint

	data, err := base64.StdEncoding.DecodeString(*cluster.CertificateAuthority.Data)
	if err != nil {
		return errors.Wrap(err, "decoding certificate authority data")
	}

	c.cfg.CertificateAuthorityData = data
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
	logger.Debug("clusters = %#v", output)
	for _, clusterName := range output.Clusters {
		if err := c.doListCluster(clusterName); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClusterProvider) doListCluster(clusterName *string) error {
	input := &eks.DescribeClusterInput{
		Name: clusterName,
	}
	output, err := c.svc.eks.DescribeCluster(input)
	if err != nil {
		return errors.Wrapf(err, "unable to describe control plane %q", *clusterName)
	}
	logger.Debug("cluster = %#v", output)
	if *output.Cluster.Status == eks.ClusterStatusActive {
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
