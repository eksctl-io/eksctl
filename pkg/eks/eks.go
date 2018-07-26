package eks

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/weaveworks/eksctl/pkg/printers"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"

	"github.com/aws/aws-sdk-go/aws"
	awseks "github.com/aws/aws-sdk-go/service/eks"

	"github.com/kubicorn/kubicorn/pkg/logger"
)

func (c *ClusterProvider) CreateControlPlane() error {
	input := &awseks.CreateClusterInput{
		Name:    &c.Spec.ClusterName,
		RoleArn: &c.Spec.clusterRoleARN,
		ResourcesVpcConfig: &awseks.VpcConfigRequest{
			SubnetIds:        aws.StringSlice(strings.Split(c.Spec.subnetsList, ",")),
			SecurityGroupIds: aws.StringSlice([]string{c.Spec.securityGroup}),
		},
	}
	output, err := c.Provider.EKS().CreateCluster(input)
	if err != nil {
		return errors.Wrap(err, "unable to create cluster control plane")
	}
	logger.Debug("output = %#v", output)
	return nil
}

func (c *ClusterProvider) DescribeControlPlane() (*awseks.Cluster, error) {
	input := &awseks.DescribeClusterInput{
		Name: &c.Spec.ClusterName,
	}
	output, err := c.Provider.EKS().DescribeCluster(input)
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

	input := &awseks.DeleteClusterInput{
		Name: cluster.Name,
	}

	if _, err := c.Provider.EKS().DeleteCluster(input); err != nil {
		return errors.Wrap(err, "unable to delete cluster control plane")
	}
	return nil
}

func (c *ClusterProvider) createControlPlane(errs chan error) error {
	logger.Info("creating control plane %q", c.Spec.ClusterName)

	clusterChan := make(chan awseks.Cluster)
	taskErrs := make(chan error)

	if err := c.CreateControlPlane(); err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(20 * time.Second)
		defer ticker.Stop()

		timer := time.NewTimer(c.Spec.WaitTimeout)
		defer timer.Stop()

		defer close(taskErrs)
		defer close(clusterChan)

		for {
			select {
			case <-timer.C:
				taskErrs <- fmt.Errorf("timed out creating control plane %q after %s", c.Spec.ClusterName, c.Spec.WaitTimeout)
				return

			case <-ticker.C:
				cluster, err := c.DescribeControlPlane()
				if err != nil {
					logger.Warning("continue despite err=%q", err.Error())
					continue
				}
				logger.Debug("cluster = %#v", cluster)
				switch *cluster.Status {
				case awseks.ClusterStatusCreating:
					continue
				case awseks.ClusterStatusActive:
					taskErrs <- nil
					clusterChan <- *cluster
					return
				default:
					taskErrs <- fmt.Errorf("unexpected status %q while creating control plane %q", *cluster.Status, c.Spec.ClusterName)
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

		if err := c.GetCredentials(cluster); err != nil {
			errs <- err
		}

		logger.Debug("clusterConfig = %#v", c.Spec)
		logger.Success("created control plane %q", c.Spec.ClusterName)

		errs <- nil
	}()

	return nil
}

func (c *ClusterProvider) GetCredentials(cluster awseks.Cluster) error {
	c.Spec.MasterEndpoint = *cluster.Endpoint

	data, err := base64.StdEncoding.DecodeString(*cluster.CertificateAuthority.Data)
	if err != nil {
		return errors.Wrap(err, "decoding certificate authority data")
	}

	c.Spec.CertificateAuthorityData = data
	return nil
}

// ListClusters display details of all the EKS cluster in your account
func (c *ClusterProvider) ListClusters(chunkSize int, printer printers.OutputPrinter) error {
	if c.Spec.ClusterName != "" {
		return c.doListCluster(&c.Spec.ClusterName, printer)
	}

	// TODO: https://github.com/weaveworks/eksctl/issues/27
	input := &awseks.ListClustersInput{}
	output, err := c.Provider.EKS().ListClusters(input)
	if err != nil {
		return errors.Wrap(err, "listing control planes")
	}
	logger.Debug("clusters = %#v", output)
	for _, clusterName := range output.Clusters {
		if err := c.doListCluster(clusterName, printer); err != nil {
			return err
		}
	}
	return nil
}

func (c *ClusterProvider) doListCluster(clusterName *string, printer printers.OutputPrinter) error {
	input := &awseks.DescribeClusterInput{
		Name: clusterName,
	}
	output, err := c.Provider.EKS().DescribeCluster(input)
	if err != nil {
		return errors.Wrapf(err, "unable to describe control plane %q", *clusterName)
	}
	logger.Debug("cluster = %#v", output)
	if *output.Cluster.Status == awseks.ClusterStatusActive {

		printer.PrintObj(output, os.Stdout)

		if logger.Level >= 4 {
			stacks, err := c.ListReadyStacks(fmt.Sprintf("^EKS-%s-.*$", *clusterName))
			if err != nil {
				return errors.Wrapf(err, "listing CloudFormation stack for %q", *clusterName)
			}
			for _, s := range stacks {
				logger.Debug("stack = %#v", *s)
			}
		}
	}
	return nil
}

func (c *ClusterProvider) ListAllTaggedResources() error {
	// TODO: https://github.com/weaveworks/eksctl/issues/26
	return nil
}

func (c *ClusterConfig) WaitForControlPlane(clientSet *kubernetes.Clientset) error {
	if _, err := clientSet.ServerVersion(); err == nil {
		return nil
	}

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	timer := time.NewTimer(c.WaitTimeout)
	defer timer.Stop()

	for {
		select {
		case <-ticker.C:
			_, err := clientSet.ServerVersion()
			if err == nil {
				return nil
			}
			logger.Debug("control plane not ready yet – %s", err.Error())
		case <-timer.C:
			return fmt.Errorf("timed out waiting for control plane %q after %s", c.ClusterName, c.WaitTimeout)
		}
	}
}
