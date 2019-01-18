package eks

import (
	"encoding/base64"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/vpc"

	awseks "github.com/aws/aws-sdk-go/service/eks"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha3"
	"github.com/weaveworks/eksctl/pkg/printers"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
)

// DescribeControlPlane describes the cluster control plane
func (c *ClusterProvider) DescribeControlPlane(cl *api.ClusterMeta) (*awseks.Cluster, error) {
	input := &awseks.DescribeClusterInput{
		Name: &cl.Name,
	}
	output, err := c.Provider.EKS().DescribeCluster(input)
	if err != nil {
		return nil, errors.Wrap(err, "unable to describe cluster control plane")
	}
	return output.Cluster, nil
}

// DescribeControlPlaneMustBeActive describes the cluster control plane and checks if status is active
func (c *ClusterProvider) DescribeControlPlaneMustBeActive(cl *api.ClusterMeta) (*awseks.Cluster, error) {
	cluster, err := c.DescribeControlPlane(cl)
	if err != nil {
		return nil, errors.Wrap(err, "unable to describe cluster control plane")
	}
	if *cluster.Status != awseks.ClusterStatusActive {
		return nil, fmt.Errorf("status of cluster %q is %q, has to be %q", *cluster.Name, *cluster.Status, awseks.ClusterStatusActive)
	}

	return cluster, nil
}

// DeprecatedDeleteControlPlane deletes the control plane
func (c *ClusterProvider) DeprecatedDeleteControlPlane(cl *api.ClusterMeta) error {
	cluster, err := c.DescribeControlPlane(cl)
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

// GetCredentials retrieves the certificate authority data
func (c *ClusterProvider) GetCredentials(spec *api.ClusterConfig) error {
	// Check the cluster exists and is active
	cluster, err := c.DescribeControlPlaneMustBeActive(spec.Metadata)
	if err != nil {
		return err
	}
	logger.Debug("cluster = %#v", cluster)

	data, err := base64.StdEncoding.DecodeString(*cluster.CertificateAuthority.Data)
	if err != nil {
		return errors.Wrap(err, "decoding certificate authority data")
	}

	if spec.Status == nil {
		spec.Status = &api.ClusterStatus{}
	}

	spec.Status.Endpoint = *cluster.Endpoint
	spec.Status.CertificateAuthorityData = data
	spec.Status.ARN = *cluster.Arn

	return nil
}

// GetClusterVPC retrieves the VPC configuration
func (c *ClusterProvider) GetClusterVPC(spec *api.ClusterConfig, ignoreMissingKeys ...string) error {
	cluster, err := c.NewStackManager(spec).DescribeClusterStack()
	if err != nil {
		return err
	}

	outputs := map[string]string{}
	for _, x := range cluster.Outputs {
		outputs[*x.OutputKey] = *x.OutputValue
	}

	if spec.VPC == nil {
		spec.VPC = &api.ClusterVPC{}
	}

	requiredKeyErrFmt := "cluster stack has no output key %q"
	isKeyRequired := func(k string) bool {
		for _, key := range ignoreMissingKeys {
			if key == k {
				return false
			}
		}
		return true
	}
	if vpc, ok := outputs[builder.CfnOutputClusterVPC]; ok {
		spec.VPC.ID = vpc
	} else {
		return fmt.Errorf(requiredKeyErrFmt, builder.CfnOutputClusterVPC)
	}

	if securityGroup, ok := outputs[builder.CfnOutputClusterSecurityGroup]; ok {
		spec.VPC.SecurityGroup = securityGroup
	} else if isKeyRequired(builder.CfnOutputClusterSharedNodeSecurityGroup) {
		return fmt.Errorf(requiredKeyErrFmt, builder.CfnOutputClusterSecurityGroup)
	}

	if sharedNodeSecurityGroup, ok := outputs[builder.CfnOutputClusterSharedNodeSecurityGroup]; ok {
		spec.VPC.SharedNodeSecurityGroup = sharedNodeSecurityGroup
	} else if isKeyRequired(builder.CfnOutputClusterSharedNodeSecurityGroup) {
		return fmt.Errorf(requiredKeyErrFmt, builder.CfnOutputClusterSharedNodeSecurityGroup)
	}

	for _, topology := range api.SubnetTopologies() {
		// either of subnet topologies are optional
		if subnets, ok := outputs[builder.CfnOutputClusterSubnets+string(topology)]; ok {
			subnets := strings.Split(subnets, ",")
			if err := vpc.UseSubnetsFromList(c.Provider, spec, topology, subnets); err != nil {
				return err
			}
		}
	}

	return nil
}

// ListClusters display details of all the EKS cluster in your account
func (c *ClusterProvider) ListClusters(clusterName string, chunkSize int, output string, eachRegion bool) error {
	// NOTE: this needs to be reworked in the future so that the functionality
	// is combined. This require the ability to return details of all clusters
	// in a single call.
	printer, err := printers.NewPrinter(output)
	if err != nil {
		return err
	}

	if clusterName != "" {
		if output == "table" {
			addSummaryTableColumns(printer.(*printers.TablePrinter))
		}
		return c.doGetCluster(clusterName, printer)
	}

	if output == "table" {
		addListTableColumns(printer.(*printers.TablePrinter))
	}
	allClusters := []*api.ClusterMeta{}
	if err := c.doListClusters(int64(chunkSize), printer, &allClusters, eachRegion); err != nil {
		return err
	}
	return printer.PrintObjWithKind("clusters", allClusters, os.Stdout)
}

func (c *ClusterProvider) getClustersRequest(chunkSize int64, nextToken string) ([]*string, *string, error) {
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

func (c *ClusterProvider) doListClusters(chunkSize int64, printer printers.OutputPrinter, allClusters *[]*api.ClusterMeta, eachRegion bool) error {
	if eachRegion {
		// reset region and re-create the client, then make a recursive call
		for _, region := range api.SupportedRegions() {
			spec := &api.ProviderConfig{
				Region:      region,
				Profile:     c.Provider.Profile(),
				WaitTimeout: c.Provider.WaitTimeout(),
			}
			if err := New(spec, nil).doListClusters(chunkSize, printer, allClusters, false); err != nil {
				return err
			}
		}
		return nil
	}

	token := ""
	for {
		clusters, nextToken, err := c.getClustersRequest(chunkSize, token)
		if err != nil {
			return err
		}

		for _, clusterName := range clusters {
			*allClusters = append(*allClusters, &api.ClusterMeta{
				Name:   *clusterName,
				Region: c.Provider.Region(),
			})
		}

		if nextToken != nil && *nextToken != "" {
			token = *nextToken
		} else {
			break
		}
	}

	return nil
}

func (c *ClusterProvider) doGetCluster(clusterName string, printer printers.OutputPrinter) error {
	input := &awseks.DescribeClusterInput{
		Name: &clusterName,
	}
	output, err := c.Provider.EKS().DescribeCluster(input)
	if err != nil {
		return errors.Wrapf(err, "unable to describe control plane %q", clusterName)
	}
	logger.Debug("cluster = %#v", output)

	clusters := []*awseks.Cluster{output.Cluster} // TODO: in the future this will have multiple clusters
	if err := printer.PrintObjWithKind("clusters", clusters, os.Stdout); err != nil {
		return err
	}

	if *output.Cluster.Status == awseks.ClusterStatusActive {

		if logger.Level >= 4 {
			spec := &api.ClusterConfig{Metadata: &api.ClusterMeta{Name: clusterName}}
			stacks, err := c.NewStackManager(spec).ListReadyStacks(fmt.Sprintf("^(eksclt|EKS)-%s-.*$", clusterName))
			if err != nil {
				return errors.Wrapf(err, "listing CloudFormation stack for %q", clusterName)
			}
			for _, s := range stacks {
				logger.Debug("stack = %#v", *s)
			}
		}
	}
	return nil
}

// WaitForControlPlane waits till the control plane is ready
func (c *ClusterProvider) WaitForControlPlane(id *api.ClusterMeta, clientSet *kubernetes.Clientset) error {
	if _, err := clientSet.ServerVersion(); err == nil {
		return nil
	}

	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	timer := time.NewTimer(c.Provider.WaitTimeout())
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
			return fmt.Errorf("timed out waiting for control plane %q after %s", id.Name, c.Provider.WaitTimeout())
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
	printer.AddColumn("NAME", func(c *api.ClusterMeta) string {
		return c.Name
	})
	printer.AddColumn("REGION", func(c *api.ClusterMeta) string {
		return c.Region
	})
}
