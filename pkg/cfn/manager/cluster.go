package manager

import (
	"fmt"
	"strings"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha3"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/utils/ipnet"
)

func (c *StackCollection) makeClusterStackName() string {
	return "eksctl-" + c.spec.Metadata.Name + "-cluster"
}

// CreateCluster creates the cluster
func (c *StackCollection) CreateCluster(errs chan error, _ interface{}) error {
	name := c.makeClusterStackName()
	logger.Info("creating cluster stack %q", name)
	stack := builder.NewClusterResourceSet(c.provider, c.spec)
	if err := stack.AddAllResources(); err != nil {
		return err
	}

	// Unlike with `CreateNodeGroup`, all tags are already set for the cluster stack
	return c.CreateStack(name, stack, nil, nil, errs)
}

// DescribeClusterStack calls DescribeStacks and filters out cluster stack
func (c *StackCollection) DescribeClusterStack() (*Stack, error) {
	stacks, err := c.DescribeStacks()
	if err != nil {
		return nil, err
	}

	for _, s := range stacks {
		if *s.StackStatus == cfn.StackStatusDeleteComplete {
			continue
		}
		if getClusterName(s) != "" {
			return s, nil
		}
	}
	return nil, nil
}

// DeleteCluster deletes the cluster
func (c *StackCollection) DeleteCluster() error {
	_, err := c.DeleteStack(c.makeClusterStackName())
	return err
}

// WaitDeleteCluster waits till the cluster is deleted
func (c *StackCollection) WaitDeleteCluster() error {
	return c.BlockingWaitDeleteStack(c.makeClusterStackName())
}

// UpdateClusterForCompability will update cluster
// with new resources based on features that have
// a critical effect on forward-compatibility with
// respect to overal functionality and integrity
func (c *StackCollection) UpdateClusterForCompability() error {
	const resourceRootPath = "Resources"

	name := c.makeClusterStackName()

	currentStack, err := c.DescribeClusterStack()
	if err != nil {
		return err
	}

	// NOTE: currently we can only append new
	// resources to the stack, as there are a
	// few limitations;
	// we don't have a way of recompiling the
	// definition of the stack from it's current
	// template and we don't have all feature
	// indicators we would need (e.g. when
	// existing VPC is used);
	// to do that in a sensible manner we would
	// have to thoroughly insepect the current
	// template and see if e.g. VPC or SGs are
	// imported or managed by us;
	// addtionally, the EKS control plane itself
	// cannot yet be updated via CloudFormation

	missingSharedNodeSecurityGroup := true

	for _, x := range currentStack.Outputs {
		switch *x.OutputKey {
		case builder.CfnOutputClusterSharedNodeSecurityGroup:
			missingSharedNodeSecurityGroup = false
		}
	}

	// Get current stack
	currentTemplate, err := c.GetStackTemplate(name)
	if err != nil {
		return errors.Wrapf(err, "error getting stack template %s", name)
	}

	updateFeatureList := []string{}
	addResources := []string{}

	if missingSharedNodeSecurityGroup {
		updateFeatureList = append(updateFeatureList, "shared node security group")
		addResources = append(addResources,
			"ClusterSharedNodeSecurityGroup",
			"IngressInterNodeGroupSG",
		)
	}

	if len(addResources) == 0 {
		logger.Success("all resources in cluster stack %q are up-to-date", name)
		return nil
	}

	currentResources := gjson.Get(currentTemplate, resourceRootPath)
	if !currentResources.IsObject() {
		return fmt.Errorf("unexpected template format of the current stack ")
	}

	{
		// We need to use same subnet CIDRs in order to recompile the template
		// with all of default resources
		vpc := c.spec.VPC
		vpc.Subnets = map[api.SubnetTopology]map[string]api.Network{
			api.SubnetTopologyPublic:  map[string]api.Network{},
			api.SubnetTopologyPrivate: map[string]api.Network{},
		}
		currentResources.ForEach(func(resourceKey, resource gjson.Result) bool {
			if resource.Get("Type").Value() == "AWS::EC2::Subnet" {
				az := resource.Get("Properties.AvailabilityZone").String()
				cidr, _ := ipnet.ParseCIDR(resource.Get("Properties.CidrBlock").String())
				k := resourceKey.String()
				if strings.HasPrefix(k, "SubnetPrivate") {
					vpc.Subnets[api.SubnetTopologyPrivate][az] = api.Network{
						CIDR: cidr,
					}
				}
				if strings.HasPrefix(k, "SubnetPublic") {
					vpc.Subnets[api.SubnetTopologyPublic][az] = api.Network{
						CIDR: cidr,
					}
				}
			}
			return true
		})
	}

	logger.Info("creating cluster stack %q", name)
	newStack := builder.NewClusterResourceSet(c.provider, c.spec)
	if err := newStack.AddAllResources(); err != nil {
		return err
	}

	newTemplate, err := newStack.RenderJSON()
	if err != nil {
		return errors.Wrapf(err, "rendering template for %q stack", name)
	}
	logger.Debug("newTemplate = %s", newTemplate)

	newResources := gjson.Get(string(newTemplate), resourceRootPath)

	if !newResources.IsObject() {
		return fmt.Errorf("unexpected template format of the new version of the stack ")
	}

	logger.Debug("currentTemplate = %s", currentTemplate)

	for _, resourceKey := range addResources {
		var err error
		resource := newResources.Get(resourceKey)
		if !resource.Exists() {
			return fmt.Errorf("resource with key %q doesn't exist in the new version of the stack", resourceKey)
		}
		currentTemplate, err = sjson.Set(currentTemplate, resourceRootPath+"."+resourceKey, resource.Value())
		if err != nil {
			return errors.Wrapf(err, "unable to add resource with key %q to cluster stack", resourceKey)
		}
	}

	logger.Debug("currentTemplate = %s", currentTemplate)

	describeUpdate := fmt.Sprintf("updating stack to add new features: %s;",
		strings.Join(updateFeatureList, ", "))
	return c.UpdateStack(name, "update-cluster", describeUpdate, []byte(currentTemplate), nil)
}

func getClusterName(s *Stack) string {
	for _, tag := range s.Tags {
		if *tag.Key == api.ClusterNameTag {
			if strings.HasSuffix(*s.StackName, "-cluster") {
				return *tag.Value
			}
		}
	}

	if strings.HasPrefix(*s.StackName, "EKS-") && strings.HasSuffix(*s.StackName, "-ControlPlane") {
		return strings.TrimPrefix("EKS-", strings.TrimSuffix(*s.StackName, "-ControlPlane"))
	}
	return ""
}
