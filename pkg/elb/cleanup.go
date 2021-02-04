package elb

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	"github.com/weaveworks/logger"
	corev1 "k8s.io/api/core/v1"
	networkingv1beta1 "k8s.io/api/networking/v1beta1"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	cloudprovider "k8s.io/cloud-provider"
	awsprovider "k8s.io/legacy-cloud-providers/aws"

	"github.com/aws/aws-sdk-go/aws/request"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

type loadBalancerKind int

const (
	classic loadBalancerKind = iota
	network
	application
)

type loadBalancer struct {
	name                  string
	kind                  loadBalancerKind
	ownedSecurityGroupIDs map[string]struct{}
}

// Cleanup finds and deletes any dangling ELBs associated to a Kubernetes Service
func Cleanup(ctx context.Context, ec2API ec2iface.EC2API, elbAPI elbiface.ELBAPI, elbv2API elbv2iface.ELBV2API,
	kubernetesCS kubernetes.Interface, clusterConfig *api.ClusterConfig) error {

	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("no context deadline set in call to elb.Cleanup()")
	}
	services, err := kubernetesCS.CoreV1().Services(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		errStr := fmt.Sprintf("cannot list Kubernetes Services: %s", err)
		if k8serrors.IsForbidden(err) {
			errStr = fmt.Sprintf("%s (deleting a cluster requires permission to list Kubernetes Services)", errStr)
		}
		return errors.New(errStr)
	}

	ingresses, err := kubernetesCS.NetworkingV1beta1().Ingresses(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		errStr := fmt.Sprintf("cannot list Kubernetes Ingresses: %s", err)
		if k8serrors.IsForbidden(err) {
			errStr = fmt.Sprintf("%s (deleting a cluster requires permission to list Kubernetes Ingresses)", errStr)
		}
		return errors.New(errStr)
	}

	// Delete Services of type 'LoadBalancer' and Ingresses kubernetes.io/ingress.class: alb
	// collecting their ELBs, NLBs and ALBs to wait for them to be deleted later on.
	awsLoadBalancers := map[string]loadBalancer{}
	// For k8s kind Service
	for _, s := range services.Items {
		lb, err := getServiceLoadBalancer(ctx, ec2API, elbAPI, clusterConfig.Metadata.Name, &s)
		if err != nil {
			return fmt.Errorf("cannot obtain information for ELB %s from LoadBalancer Service %s/%s: %s",
				cloudprovider.DefaultLoadBalancerName(&s), s.Namespace, s.Name, err)
		}
		if lb == nil {
			continue
		}
		logger.Debug(
			"tracking deletion of load balancer %s of kind %d with security groups %v",
			lb.name, lb.kind, convertStringSetToSlice(lb.ownedSecurityGroupIDs),
		)
		awsLoadBalancers[lb.name] = *lb
		logger.Debug("deleting 'type: LoadBalancer' Service %s/%s", s.Namespace, s.Name)
		if err := kubernetesCS.CoreV1().Services(s.Namespace).Delete(context.TODO(), s.Name, metav1.DeleteOptions{}); err != nil {
			errStr := fmt.Sprintf("cannot delete Kubernetes Service %s/%s: %s", s.Namespace, s.Name, err)
			if k8serrors.IsForbidden(err) {
				errStr = fmt.Sprintf("%s (deleting a cluster requires permission to delete Kubernetes Services)", errStr)
			}
			return errors.New(errStr)
		}
	}
	// For k8s Kind Ingress
	for _, i := range ingresses.Items {
		lb := getIngressLoadBalancer(ctx, i)
		if lb == nil {
			continue
		}
		logger.Debug(
			"tracking deletion of load balancer %s of kind %d with security groups %v",
			lb.name, lb.kind, convertStringSetToSlice(lb.ownedSecurityGroupIDs),
		)
		awsLoadBalancers[lb.name] = *lb
		logger.Debug("deleting 'kubernetes.io/ingress.class: alb' Ingress %s/%s", i.Namespace, i.Name)
		if err := kubernetesCS.NetworkingV1beta1().Ingresses(i.Namespace).Delete(context.TODO(), i.Name, metav1.DeleteOptions{}); err != nil {
			errStr := fmt.Sprintf("cannot delete Kubernetes Ingress %s/%s: %s", i.Namespace, i.Name, err)
			if k8serrors.IsForbidden(err) {
				errStr = fmt.Sprintf("%s (deleting a cluster requires permission to delete Kubernetes Ingress)", errStr)
			}
			return errors.New(errStr)
		}
	}

	// Wait for all the load balancers backing the LoadBalancer services to disappear
	pollInterval := 2 * time.Second
	for ; time.Now().Before(deadline) && len(awsLoadBalancers) > 0; time.Sleep(pollInterval) {
		for name, lb := range awsLoadBalancers {
			exists, err := loadBalancerExists(ctx, ec2API, elbAPI, elbv2API, lb)
			if err != nil {
				logger.Warning("error when checking existence of load balancer %s: %s", lb.name, err)
			}
			if exists {
				continue
			}
			logger.Debug("load balancer %s and its security groups were deleted by the cloud provider", name)
			// The load balancer and its security groups have been deleted
			delete(awsLoadBalancers, name)
		}
	}

	if numLB := len(awsLoadBalancers); numLB > 0 {
		lbs := make([]string, 0, numLB)
		for name := range awsLoadBalancers {
			lbs = append(lbs, name)
		}
		return fmt.Errorf("deadline surpassed waiting for AWS load balancers to be deleted: %s", strings.Join(lbs, ","))
	}
	logger.Debug("deleting load balancer Security Group orphans")
	// Orphan security-group deletion is needed due to https://github.com/kubernetes/kubernetes/issues/79994`
	// and because we could have started the service deletion when a service didn't finish its creation
	if err := deleteOrphanLoadBalancerSecurityGroups(ctx, ec2API, elbAPI, clusterConfig); err != nil {
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

func getIngressLoadBalancer(ctx context.Context, ingress networkingv1beta1.Ingress) (lb *loadBalancer) {
	ingressCls := "kubernetes.io/ingress.class"
	if ingress.Annotations[ingressCls] != "alb" {
		logger.Debug("%s is not ALB Ingress, it is '%s': '%s', skip", ingress.Name, ingressCls, ingress.Annotations[ingressCls])
		return nil
	}

	// Check if field status.loadBalancer.ingress[].hostname is set, value corresponds with name for AWS ALB
	// if does not pass ALB hadn't been provisioned so nothing to return.
	if len(ingress.Status.LoadBalancer.Ingress) == 0 ||
		len(ingress.Status.LoadBalancer.Ingress[0].Hostname) == 0 {
		logger.Debug("%s is ALB Ingress, but probably not provisioned, skip", ingress.Name)
		return nil
	}
	// Expected e.g. bf647c9e-default-appingres-350b-1622159649.eu-central-1.elb.amazonaws.com where AWS ALB name is
	// bf647c9e-default-appingres-350b (cannot be longer than 32 characters).
	hostNameParts := strings.Split(ingress.Status.LoadBalancer.Ingress[0].Hostname, ".")
	if len(hostNameParts[0]) == 0 {
		logger.Debug("%s is ALB Ingress, but probably not provisioned or something other unexpected, skip", ingress.Name)
		return nil
	}
	name := strings.TrimPrefix(hostNameParts[0], "internal-") // Trim 'internal-' prefix for ALB DNS name which is not a part of name.
	if len(name) > 31 {
		name = name[:31]
	}
	return &loadBalancer{
		name:                  name,
		kind:                  application,
		ownedSecurityGroupIDs: map[string]struct{}{},
	}
}

func convertStringSetToSlice(set map[string]struct{}) []string {
	result := make([]string, 0, len(set))
	for value := range set {
		result = append(result, value)
	}
	return result
}

// Load balancers provisioned by the AWS cloud-provider integration are named k8s-elb-$loadBalancerName
var sgNameRegex = regexp.MustCompile(`^k8s-elb-([^-]{1-32})$`)

func deleteOrphanLoadBalancerSecurityGroups(ctx context.Context, ec2API ec2iface.EC2API, elbAPI elbiface.ELBAPI, clusterConfig *api.ClusterConfig) error {
	clusterName := clusterConfig.Metadata.Name
	describeRequest := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("tag-key"),
				Values: []*string{aws.String(awsprovider.TagNameKubernetesClusterPrefix + clusterName)},
			},
		},
	}
	var failedSGs []*ec2.SecurityGroup
	for {
		result, err := ec2API.DescribeSecurityGroupsWithContext(ctx, describeRequest)
		if err != nil {
			return fmt.Errorf("cannot describe security groups: %s", err)
		}
		for _, sg := range result.SecurityGroups {
			if !sgNameRegex.MatchString(*sg.GroupName) {
				logger.Debug("ignoring non-matching security group %q", *sg.GroupName)
				continue
			}
			if err := deleteSecurityGroup(ctx, ec2API, sg); err != nil {
				if awsError, ok := err.(awserr.Error); ok && awsError.Code() == "DependencyViolation" {
					logger.Debug("failed to delete security group, possibly because its load balancer is still being deleted")
					failedSGs = append(failedSGs, sg)
				} else {
					return errors.Wrapf(err, "cannot delete security group %q", *sg.GroupName)
				}
			}
		}
		if result.NextToken == nil {
			break
		}
		describeRequest.NextToken = result.NextToken
	}

	if len(failedSGs) > 0 {
		logger.Debug("retrying deletion of %d security groups", len(failedSGs))
		return deleteFailedSecurityGroups(ctx, ec2API, elbAPI, failedSGs)
	}

	return nil
}

func deleteFailedSecurityGroups(ctx aws.Context, ec2API ec2iface.EC2API, elbAPI elbiface.ELBAPI, securityGroups []*ec2.SecurityGroup) error {
	for _, sg := range securityGroups {
		// wait for the security group's load balancer to complete deletion
		if err := ensureLoadBalancerDeleted(ctx, elbAPI, sg); err != nil {
			return err
		}
		if err := deleteSecurityGroup(ctx, ec2API, sg); err != nil {
			return errors.Wrapf(err, "failed to delete security group (name: %q, id: %q)", *sg.GroupName, *sg.GroupId)
		}
	}
	return nil
}

func ensureLoadBalancerDeleted(ctx context.Context, elbAPI elbiface.ELBAPI, sg *ec2.SecurityGroup) error {
	// extract load balancer name from the SG ID
	match := sgNameRegex.FindStringSubmatch(*sg.GroupName)
	if len(match) != 2 {
		return fmt.Errorf("unexpected security group name format: %q", *sg.GroupName)
	}
	loadBalancerName := match[1]

	var (
		lbDeleteTimeout = 30 * time.Second
		lbRetryAfter    = 3 * time.Second
	)

	timeoutCtx, cancelFunc := context.WithTimeout(ctx, lbDeleteTimeout)
	defer cancelFunc()

	for {
		input := &elb.DescribeLoadBalancersInput{
			LoadBalancerNames: []*string{aws.String(loadBalancerName)},
		}

		if _, err := elbAPI.DescribeLoadBalancersWithContext(timeoutCtx, input); err != nil {
			if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == elb.ErrCodeAccessPointNotFoundException {
				return nil
			} else if err == context.DeadlineExceeded {
				return errors.Wrap(err, "timed out looking up load balancer")
			} else if request.IsErrorRetryable(err) {
				logger.Debug("retrying request after %v", lbRetryAfter)
			} else {
				return err
			}
		}

		timer := time.NewTimer(lbRetryAfter)
		select {
		case <-timeoutCtx.Done():
			timer.Stop()
			return errors.Wrap(timeoutCtx.Err(), "timed out waiting for load balancer's deletion")
		case <-timer.C:
		}

	}
}

func deleteSecurityGroup(ctx context.Context, ec2API ec2iface.EC2API, sg *ec2.SecurityGroup) error {
	logger.Debug("deleting orphan Load Balancer security group %s with description %q",
		aws.StringValue(sg.GroupId), aws.StringValue(sg.Description))
	input := &ec2.DeleteSecurityGroupInput{
		GroupId: sg.GroupId,
	}
	_, err := ec2API.DeleteSecurityGroupWithContext(ctx, input)
	return err
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
	clusterName string, loadBalancerName string, loadBalancerKind loadBalancerKind) (map[string]struct{}, error) {
	if loadBalancerKind == network {
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

func getLoadBalancerKind(service *corev1.Service) loadBalancerKind {
	// See https://github.com/kubernetes/legacy-cloud-providers/blob/master/aws/aws_loadbalancer.go#L65-L70
	if service.Annotations[awsprovider.ServiceAnnotationLoadBalancerType] == "nlb" {
		return network
	}
	return classic
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
	name string, kind loadBalancerKind) (bool, error) {
	if kind == network || kind == application {
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
