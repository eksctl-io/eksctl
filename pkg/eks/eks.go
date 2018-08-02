package eks

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/weaveworks/eksctl/pkg/printers"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/util/sets"
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
func (c *ClusterProvider) ListClusters(chunkSize int, output string) error {
	// NOTE: this needs to be reworked in the future so that the functionality
	// is combined. This require the ability to return details of all clusters
	// in a single call.
	printer, err := printers.NewPrinter(output)
	if err != nil {
		return err
	}

	if c.Spec.ClusterName != "" {
		if output == "table" {
			addSummaryTableColumns(printer.(*printers.TablePrinter))
		}
		return c.doGetCluster(&c.Spec.ClusterName, printer)
	}

	if output == "table" {
		addListTableColumns(printer.(*printers.TablePrinter))
	}
	return c.doListClusters(int64(chunkSize), printer)
}

func (c *ClusterProvider) doListClusters(chunkSize int64, printer printers.OutputPrinter) error {
	allClusterNames := []*string{}

	getFunc := func(chunkSize int64, nextToken string) ([]*string, *string, error) {
		input := &awseks.ListClustersInput{MaxResults: &chunkSize}
		if nextToken != "" {
			input = input.SetNextToken(nextToken)
		}
		output, err := c.Provider.EKS().ListClusters(input)
		if err != nil {
			return nil, nil, errors.Wrap(err, "listing control planes")
		}
		return output.Clusters, output.NextToken, nil
	}

	token := ""
	for {
		clusters, nextToken, err := getFunc(chunkSize, token)
		if err != nil {
			return err
		}
		allClusterNames = append(allClusterNames, clusters...)
		if nextToken != nil && *nextToken != "" {
			token = *nextToken
		} else {
			break
		}
	}

	printer.PrintObj(allClusterNames, os.Stdout)

	return nil
}

func (c *ClusterProvider) doGetCluster(clusterName *string, printer printers.OutputPrinter) error {
	input := &awseks.DescribeClusterInput{
		Name: clusterName,
	}
	output, err := c.Provider.EKS().DescribeCluster(input)
	if err != nil {
		return errors.Wrapf(err, "unable to describe control plane %q", *clusterName)
	}
	logger.Debug("cluster = %#v", output)
	if *output.Cluster.Status == awseks.ClusterStatusActive {
		clusters := []*awseks.Cluster{output.Cluster} // TODO: in the future this will have multiple clusters
		printer.PrintObj(clusters, os.Stdout)

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

func addSummaryTableColumns(printer *printers.TablePrinter) {
	printer.AddColumn("NAME", func(c *awseks.Cluster) string {
		return *c.Name
	})
	printer.AddColumn("VERSION", func(c *awseks.Cluster) string {
		return *c.Version
	})
	printer.AddColumn("STATUS", func(c *awseks.Cluster) string {
		return *c.Status
	})
	printer.AddColumn("CREATED", func(c *awseks.Cluster) string {
		return c.CreatedAt.Format(time.RFC3339)
	})
	printer.AddColumn("VPC", func(c *awseks.Cluster) string {
		return *c.ResourcesVpcConfig.VpcId
	})
	printer.AddColumn("SUBNETS", func(c *awseks.Cluster) string {
		subnets := sets.NewString()
		for _, subnetid := range c.ResourcesVpcConfig.SubnetIds {
			if *subnetid != "" {
				subnets.Insert(*subnetid)
			}
		}
		return strings.Join(subnets.List(), ",")
	})
	printer.AddColumn("SECURITYGROUPS", func(c *awseks.Cluster) string {
		groups := sets.NewString()
		for _, sg := range c.ResourcesVpcConfig.SecurityGroupIds {
			if *sg != "" {
				groups.Insert(*sg)
			}
		}
		return strings.Join(groups.List(), ",")
	})
}

func addListTableColumns(printer *printers.TablePrinter) {
	printer.AddColumn("NAME", func(c *string) string {
		return *c
	})
}
