package elb

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/kris-nova/logger"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/cloudprovider"
	awsprovider "k8s.io/kubernetes/pkg/cloudprovider/providers/aws"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

const (
	loadBalancerKindClassic = iota
	loadBalancerKindNetwork = iota
)

type loadBalancer struct {
	name                  string
	kind                  int
	ownedSecurityGroupIDs map[string]struct{}
}

// Cleanup finds and deletes any dangling ELBs associated to a Kubernetes Service
func Cleanup(ctx context.Context, ec2API ec2iface.EC2API, elbAPI elbiface.ELBAPI, elbv2API elbv2iface.ELBV2API,
	kubernetesCS kubernetes.Interface, clusterConfig *api.ClusterConfig) error {

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("no context deadline set in call to elb.Cleanup()")
	}
	services, err := kubernetesCS.CoreV1().Services(metav1.NamespaceAll).List(metav1.ListOptions{})
	if err != nil {
		errStr := fmt.Sprintf("cannot list Kubernetes Services: %s", err)
		if k8serrors.IsForbidden(err) {
			errStr = fmt.Sprintf("%s (deleting a cluster requires permission to list Kubernetes services)", errStr)
		}
		return errors.New(errStr)
	}

	// Delete Services of type 'LoadBalancer', collecting their ELBs to wait for them to be deleted later on
	loadBalancers := map[string]loadBalancer{}
	for _, s := range services.Items {
		lb, err := getServiceLoadBalancer(ctx, ec2API, elbAPI, clusterConfig.Metadata.Name, &s)
		if err != nil {
			return fmt.Errorf("cannot obtain information for ELB %s from LoadBalancer service %s/%s: %s",
				cloudprovider.DefaultLoadBalancerName(&s), s.Namespace, s.Name, err)
		}
		if lb == nil {
			continue
		}
		logger.Debug("tracking deletion of Load Balancer %s of kind %d with security groups %v",
			lb.name, lb.kind, convertStringSetToSlice(lb.ownedSecurityGroupIDs))
		loadBalancers[lb.name] = *lb
		logger.Debug("deleting 'type: LoadBalancer' service %s/%s", s.Namespace, s.Name)
		if err := kubernetesCS.CoreV1().Services(s.Namespace).Delete(s.Name, &metav1.DeleteOptions{}); err != nil {
			errStr := fmt.Sprintf("cannot delete Kubernetes Service %s/%s: %s", s.Namespace, s.Name, err)
			if k8serrors.IsForbidden(err) {
				errStr = fmt.Sprintf("%s (deleting a cluster requires permission to delete Kubernetes services)", errStr)
			}
			return errors.New(errStr)
		}
	}

	// Wait for all the load balancers backing the LoadBalancer services to disappear
	pollInterval := 2 * time.Second
	for ; time.Now().Before(deadline) && len(loadBalancers) > 0; time.Sleep(pollInterval) {
		for name, lb := range loadBalancers {
			exists, err := loadBalancerExists(ctx, ec2API, elbAPI, elbv2API, lb)
			if err != nil {
				logger.Warning("error when checking existence of load balancer %s: %s", lb.name, err)
			}
			if exists {
				continue
			}
			logger.Debug("load balancer %s and its security groups were deleted by the cloud provider", name)
			// The load balancer and its security groups have been deleted
			delete(loadBalancers, name)
		}
	}

	if len(loadBalancers) > 0 {
		return fmt.Errorf("deadline surpased waiting for load balancers to be deleted")
	}
	logger.Debug("deleting Load Balancer Security Group orphans")
	// Orphan security-group deletion is needed due to https://github.com/kubernetes/kubernetes/issues/79994
	// and because we could have started the service deletion when a service didn't finish its creation
	if err := deleteOrphanLoadBalancerSecurityGroups(ctx, ec2API, clusterConfig); err != nil {
		return fmt.Errorf("cannot delete orphan ELB Security Groups: %s", err)
	}
	return nil
}

func getServiceLoadBalancer(ctx context.Context, ec2API ec2iface.EC2API, elbAPI elbiface.ELBAPI, clusterName string,
	service *corev1.Service) (*loadBalancer, error) {
	if service.Spec.Type != corev1.ServiceTypeLoadBalancer {
		return nil, nil
	}
	name := cloudprovider.DefaultLoadBalancerName(service)
	kind := getLoadBalancerKind(service)
	ctx, cleanup := context.WithTimeout(ctx, 30*time.Second)
	securityGroupIDs, err := getSecurityGroupsOwnedByLoadBalancer(ctx, ec2API, elbAPI, clusterName, name, kind)
	cleanup()
	if err != nil {
		return nil, fmt.Errorf("cannot obtain security groups for ELB %s: %s", name, err)
	}
	lb := loadBalancer{
		name:                  name,
		kind:                  kind,
		ownedSecurityGroupIDs: securityGroupIDs,
	}
	return &lb, nil
}

func convertStringSetToSlice(set map[string]struct{}) []string {
	result := make([]string, 0, len(set))
	for value := range set {
		result = append(result, value)
	}
	return result
}

func deleteOrphanLoadBalancerSecurityGroups(ctx context.Context, ec2API ec2iface.EC2API, clusterConfig *api.ClusterConfig) error {
	clusterName := clusterConfig.Metadata.Name
	describeRequest := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag-key"),
				Values: []*string{aws.String(awsprovider.TagNameKubernetesClusterPrefix + clusterName)},
			},
		},
	}
	for {
		result, err := ec2API.DescribeSecurityGroupsWithContext(ctx, describeRequest)
		if err != nil {
			return fmt.Errorf("cannot describe security groups: %s", err)
		}
		for _, sg := range result.SecurityGroups {
			if !strings.HasPrefix(*sg.GroupName, "k8s-elb-") {
				continue
			}
			logger.Debug("deleting orphan Load Balancer security group %s with description %q",
				aws.StringValue(sg.GroupId), aws.StringValue(sg.Description))
			input := &ec2.DeleteSecurityGroupInput{
				GroupId: sg.GroupId,
			}
			if _, err := ec2API.DeleteSecurityGroupWithContext(ctx, input); err != nil {
				return fmt.Errorf("cannot delete security group %s: %s", *sg.GroupName, err)
			}
		}
		if result.NextToken == nil {
			break
		}
		describeRequest.NextToken = result.NextToken
	}
	return nil
}

func describeSecurityGroupsByID(ctx context.Context, ec2API ec2iface.EC2API, groupIDs []string) (*ec2.DescribeSecurityGroupsOutput, error) {
	filter := &ec2.Filter{
		Name: aws.String("group-id"),
	}
	for _, groupID := range groupIDs {
		filter.Values = append(filter.Values, aws.String(groupID))
	}
	describeRequest := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{filter},
	}
	return ec2API.DescribeSecurityGroupsWithContext(ctx, describeRequest)
}

func tagsIncludeClusterName(tags []*ec2.Tag, clusterName string) bool {
	clusterTagKey := awsprovider.TagNameKubernetesClusterPrefix + clusterName
	for _, tag := range tags {
		if aws.StringValue(tag.Key) == clusterTagKey {
			return true
		}
	}
	return false
}

func getSecurityGroupsOwnedByLoadBalancer(ctx context.Context, ec2API ec2iface.EC2API, elbAPI elbiface.ELBAPI,
	clusterName string, loadBalancerName string, loadBalancerKind int) (map[string]struct{}, error) {
	if loadBalancerKind == loadBalancerKindNetwork {
		// V2 ELBs just use the Security Group of the EC2 instances
		return map[string]struct{}{}, nil
	}

	lb, err := describeClassicLoadBalancer(ctx, elbAPI, loadBalancerName)
	if err != nil {
		return nil, fmt.Errorf("cannot describe ELB: %s", err)
	}
	if lb == nil {
		// The load balancer wasn't found
		return map[string]struct{}{}, nil
	}

	sgResponse, err := describeSecurityGroupsByID(ctx, ec2API, aws.StringValueSlice(lb.SecurityGroups))

	if err != nil {
		return nil, fmt.Errorf("error obtaining security groups for ELB: %s", err)
	}

	result := map[string]struct{}{}

	for _, sg := range sgResponse.SecurityGroups {
		sgID := aws.StringValue(sg.GroupId)

		// FIXME(fons): AWS' CloudConfig accepts a global ELB security group, which shouldn't be deleted.
		//              However, there doesn't seem to be a way to access the CloudConfiguration through the API Server.
		//              Regardless, EKS doesn't expose the CloudConfiguration at the time of writing (which in
		//              turn doesn't allow setting the ELB security group).
		// if sgID == cfg.Global.ElbSecurityGroup {
		//	//We don't want to delete a security group that was defined in the Cloud Configuration.
		//	continue
		// }

		if sgID == "" {
			continue
		}

		// Only delete security groups created by the Kubernetes Cloud provider
		if !tagsIncludeClusterName(sg.Tags, clusterName) {
			continue
		}

		result[sgID] = struct{}{}
	}

	return result, nil
}

func getLoadBalancerKind(service *corev1.Service) int {
	// See https://github.com/kubernetes/kubernetes/blob/v1.12.6/pkg/cloudprovider/providers/aws/aws_loadbalancer.go#L51-L56
	if service.Annotations[awsprovider.ServiceAnnotationLoadBalancerType] == "nlb" {
		return loadBalancerKindNetwork
	}
	return loadBalancerKindClassic
}

func loadBalancerExists(ctx context.Context, ec2API ec2iface.EC2API, elbAPI elbiface.ELBAPI, elbv2API elbv2iface.ELBV2API, lb loadBalancer) (bool, error) {
	exists, err := elbExists(ctx, elbAPI, elbv2API, lb.name, lb.kind)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil
	}

	// Check whether all the security groups owned by the load balancer have also been deleted
	// (they are deleted after the load balancer)
	sgResponse, err := describeSecurityGroupsByID(ctx, ec2API, convertStringSetToSlice(lb.ownedSecurityGroupIDs))
	if err != nil {
		return false, err
	}
	return len(sgResponse.SecurityGroups) != 0, nil
}

func elbExists(ctx context.Context, elbAPI elbiface.ELBAPI, elbv2API elbv2iface.ELBV2API,
	name string, kind int) (bool, error) {
	if kind == loadBalancerKindNetwork {
		return elbV2Exists(ctx, elbv2API, name)
	}
	desc, err := describeClassicLoadBalancer(ctx, elbAPI, name)
	return desc != nil, err
}

func describeClassicLoadBalancer(ctx context.Context, elbAPI elbiface.ELBAPI,
	name string) (*elb.LoadBalancerDescription, error) {

	request := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{&name},
	}

	response, err := elbAPI.DescribeLoadBalancersWithContext(ctx, request)
	if err != nil {
		if awsError, ok := err.(awserr.Error); ok {
			if awsError.Code() == elb.ErrCodeAccessPointNotFoundException {
				return nil, nil
			}
		}
		return nil, err
	}

	var ret *elb.LoadBalancerDescription
	switch {
	case len(response.LoadBalancerDescriptions) > 1:
		logger.Warning("found multiple load balancers with name: %s", name)
		fallthrough
	case len(response.LoadBalancerDescriptions) > 0:
		ret = response.LoadBalancerDescriptions[0]
	}
	return ret, nil
}

func elbV2Exists(ctx context.Context, api elbv2iface.ELBV2API, name string) (bool, error) {
	request := &elbv2.DescribeLoadBalancersInput{
		Names: []*string{aws.String(name)},
	}

	_, err := api.DescribeLoadBalancersWithContext(ctx, request)
	if err != nil {
		if awsError, ok := err.(awserr.Error); ok {
			if awsError.Code() == elbv2.ErrCodeLoadBalancerNotFoundException {
				return false, nil
			}
		}
		return false, err
	}

	return true, nil
}
