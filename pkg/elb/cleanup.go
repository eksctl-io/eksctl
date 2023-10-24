package elb

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/aws/smithy-go"

	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	elbtypes "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing/types"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
	elbv2types "github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2/types"

	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kris-nova/logger"
	corev1 "k8s.io/api/core/v1"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	cloudprovider "k8s.io/cloud-provider"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
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

const (
	serviceAnnotationLoadBalancerType = "service.beta.kubernetes.io/aws-load-balancer-type"
	tagNameKubernetesClusterPrefix    = "kubernetes.io/cluster/"
	elbv2ClusterTagKey                = "elbv2.k8s.aws/cluster"
)

// DescribeLoadBalancersAPI provides an interface to the AWS ELB service.
type DescribeLoadBalancersAPI interface {
	// DescribeLoadBalancers describes the specified load balancers or all load balancers.
	DescribeLoadBalancers(ctx context.Context, params *elasticloadbalancing.DescribeLoadBalancersInput, optFns ...func(*elasticloadbalancing.Options)) (*elasticloadbalancing.DescribeLoadBalancersOutput, error)
}

// DescribeLoadBalancersAPIV2 provides an interface to the AWS ELBv2 service.
type DescribeLoadBalancersAPIV2 interface {
	// DescribeLoadBalancers describes the specified load balancers or all load balancers for ELBv2.
	DescribeLoadBalancers(ctx context.Context, params *elasticloadbalancingv2.DescribeLoadBalancersInput, optFns ...func(*elasticloadbalancingv2.Options)) (*elasticloadbalancingv2.DescribeLoadBalancersOutput, error)
}

// Cleanup finds and deletes any dangling ELBs associated to a Kubernetes Service
func Cleanup(ctx context.Context, ec2API awsapi.EC2, elbAPI DescribeLoadBalancersAPI, elbv2API DescribeLoadBalancersAPIV2,
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

	ingresses, err := listIngress(kubernetesCS, clusterConfig)
	if err != nil {
		errStr := fmt.Sprintf("cannot list Kubernetes Ingresses: %s", err)
		if k8serrors.IsForbidden(err) {
			errStr = fmt.Sprintf("%s (deleting a cluster requires permission to list Kubernetes Ingresses)", errStr)
		}
		return errors.New(errStr)
	}

	// Delete Services of type 'LoadBalancer' and Ingresses with IngressClass of alb
	// collecting their ELBs, NLBs and ALBs to wait for them to be deleted later on.
	awsLoadBalancers := map[string]loadBalancer{}
	// For k8s kind Service
	for _, s := range services.Items {
		lb, err := getServiceLoadBalancer(ctx, ec2API, elbAPI, elbv2API, clusterConfig.Metadata.Name, &s)
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
	for _, i := range ingresses {
		ingressMetadata := i.GetMetadata()

		lb, err := getIngressLoadBalancer(ctx, ec2API, elbAPI, elbv2API, clusterConfig.Metadata.Name, i)
		if err != nil {
			return fmt.Errorf("cannot obtain information for ALB from Ingress %s/%s: %w",
				ingressMetadata.Namespace, ingressMetadata.Name, err)
		}
		if lb == nil {
			continue
		}
		logger.Debug(
			"tracking deletion of load balancer %s of kind %d with security groups %v",
			lb.name, lb.kind, convertStringSetToSlice(lb.ownedSecurityGroupIDs),
		)
		awsLoadBalancers[lb.name] = *lb
		logger.Debug("deleting ALB Ingress %s/%s", ingressMetadata.Namespace, ingressMetadata.Name)
		if err := i.Delete(kubernetesCS); err != nil {
			errStr := fmt.Sprintf("cannot delete Kubernetes Ingress %s/%s: %s", ingressMetadata.Namespace, ingressMetadata.Name, err)
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

func getServiceLoadBalancer(ctx context.Context, ec2API awsapi.EC2, elbAPI DescribeLoadBalancersAPI, elbv2API DescribeLoadBalancersAPIV2,
	clusterName string, service *corev1.Service) (*loadBalancer, error) {
	if service.Spec.Type != corev1.ServiceTypeLoadBalancer {
		return nil, nil
	}
	name := cloudprovider.DefaultLoadBalancerName(service)
	kind := getLoadBalancerKind(service)
	ctx, cleanup := context.WithTimeout(ctx, 30*time.Second)
	securityGroupIDs, err := getSecurityGroupsOwnedByLoadBalancer(ctx, ec2API, elbAPI, elbv2API, clusterName, name, kind)
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

func getIngressELBName(hosts []string) (string, error) {
	// Expected e.g. bf647c9e-default-appingres-350b-1622159649.eu-central-1.elb.amazonaws.com where AWS ALB name is
	// bf647c9e-default-appingres-350b (cannot be longer than 32 characters).
	hostNameParts := strings.Split(hosts[0], ".")
	if len(hostNameParts[0]) == 0 {
		return "", fmt.Errorf("cannot get the hostname: %v", hostNameParts)
	}
	name := strings.TrimPrefix(hostNameParts[0], "internal-") // Trim 'internal-' prefix for ALB DNS name which is not a part of name.
	idIdx := strings.LastIndex(name, "-")
	if (idIdx) != -1 {
		name = name[:idIdx] // Remove the ELB ID and last hyphen at the end of the hostname (ELB name cannot end with a hyphen)
	}
	if len(name) > 32 {
		return "", fmt.Errorf("parsed name exceeds maximum of 32 characters: %s", name)
	}
	return name, nil
}

func getIngressLoadBalancer(ctx context.Context, ec2API awsapi.EC2, elbAPI DescribeLoadBalancersAPI, elbv2API DescribeLoadBalancersAPIV2,
	clusterName string, ingress Ingress) (*loadBalancer, error) {
	metadata := ingress.GetMetadata()

	if ingress.GetIngressClass() != "alb" {
		logger.Debug("%s is not ALB Ingress, Ingress Class is '%s', skip", metadata.Name, ingress.GetIngressClass())
		return nil, nil
	}

	// Check if field status.loadBalancer.ingress[].hostname is set, value corresponds with name for AWS ALB
	// if does not pass ALB hadn't been provisioned so nothing to return.
	hosts := ingress.GetLoadBalancersHosts()
	if len(hosts) == 0 {
		logger.Debug("%s is ALB Ingress, but probably not provisioned, skip", metadata.Name)
		return nil, nil
	}
	name, err := getIngressELBName(hosts)
	if err != nil {
		logger.Debug("%s is ALB Ingress, but probably not provisioned or something other unexpected when getting ALB resource name, skip: %s", metadata.Name, err)
		return nil, nil
	}
	logger.Debug("ALB resource name: %s", name)
	ctx, cleanup := context.WithTimeout(ctx, 30*time.Second)
	defer cleanup()
	securityGroupIDs, err := getSecurityGroupsOwnedByLoadBalancer(ctx, ec2API, elbAPI, elbv2API, clusterName, name, application)
	if err != nil {
		return nil, fmt.Errorf("cannot obtain security groups for ALB %s: %w", name, err)
	}
	return &loadBalancer{
		name:                  name,
		kind:                  application,
		ownedSecurityGroupIDs: securityGroupIDs,
	}, nil
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

func deleteOrphanLoadBalancerSecurityGroups(ctx context.Context, ec2API awsapi.EC2, elbAPI DescribeLoadBalancersAPI, clusterConfig *api.ClusterConfig) error {
	clusterName := clusterConfig.Metadata.Name
	describeRequest := &ec2.DescribeSecurityGroupsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("tag-key"),
				Values: []string{tagNameKubernetesClusterPrefix + clusterName},
			},
		},
	}
	var failedSGs []ec2types.SecurityGroup
	paginator := ec2.NewDescribeSecurityGroupsPaginator(ec2API, describeRequest)
	for paginator.HasMorePages() {
		result, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("cannot describe security groups: %s", err)
		}
		for _, sg := range result.SecurityGroups {
			if !sgNameRegex.MatchString(*sg.GroupName) {
				logger.Debug("ignoring non-matching security group %q", *sg.GroupName)
				continue
			}
			if err := deleteSecurityGroup(ctx, ec2API, sg); err != nil {
				var ae smithy.APIError
				if errors.As(err, &ae) && ae.ErrorCode() == "DependencyViolation" {
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

func deleteFailedSecurityGroups(ctx context.Context, ec2API awsapi.EC2, elbAPI DescribeLoadBalancersAPI, securityGroups []ec2types.SecurityGroup) error {
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

func ensureLoadBalancerDeleted(ctx context.Context, elbAPI DescribeLoadBalancersAPI, sg ec2types.SecurityGroup) error {
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
		input := &elasticloadbalancing.DescribeLoadBalancersInput{
			LoadBalancerNames: []string{loadBalancerName},
		}
		if _, err := elbAPI.DescribeLoadBalancers(timeoutCtx, input); err != nil {
			switch {
			case isELBNotFoundErr(err):
				return nil
			case err == context.DeadlineExceeded:
				return errors.Wrap(err, "timed out looking up load balancer")
			case retry.IsErrorRetryables(retry.DefaultRetryables).IsErrorRetryable(err).Bool():
				// This is not required when a retryer is configured. It exists to maintain the existing behaviour
				// of manually retrying requests.
				logger.Debug("retrying request after %v", lbRetryAfter)
			default:
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

func deleteSecurityGroup(ctx context.Context, ec2API awsapi.EC2, sg ec2types.SecurityGroup) error {
	logger.Debug("deleting orphan Load Balancer security group %s with description %q",
		aws.ToString(sg.GroupId), aws.ToString(sg.Description))
	input := &ec2.DeleteSecurityGroupInput{
		GroupId: sg.GroupId,
	}
	_, err := ec2API.DeleteSecurityGroup(ctx, input)
	return err
}

func describeSecurityGroupsByID(ctx context.Context, ec2API awsapi.EC2, groupIDs []string) (*ec2.DescribeSecurityGroupsOutput, error) {
	return ec2API.DescribeSecurityGroups(ctx, &ec2.DescribeSecurityGroupsInput{
		Filters: []ec2types.Filter{
			{
				Name:   aws.String("group-id"),
				Values: groupIDs,
			},
		},
	})
}

func tagsIncludeClusterName(tags []ec2types.Tag, clusterName string) bool {
	k8sClusterTagKey := tagNameKubernetesClusterPrefix + clusterName
	for _, tag := range tags {
		switch aws.ToString(tag.Key) {
		case k8sClusterTagKey, elbv2ClusterTagKey:
			return true
		}
	}
	return false
}

func getSecurityGroupsOwnedByLoadBalancer(ctx context.Context, ec2API awsapi.EC2, elbAPI DescribeLoadBalancersAPI,
	elbv2API DescribeLoadBalancersAPIV2, clusterName string, loadBalancerName string, loadBalancerKind loadBalancerKind) (map[string]struct{}, error) {

	var groupIDs []string

	switch loadBalancerKind {
	case network:
		// V2 ELBs just use the Security Group of the EC2 instances
		return map[string]struct{}{}, nil
	case application:
		alb, err := describeApplicationLoadBalancer(ctx, elbv2API, loadBalancerName)

		if err != nil {
			return nil, fmt.Errorf("cannot describe ELB: %w", err)
		}
		if alb == nil {
			// The load balancer wasn't found
			return map[string]struct{}{}, nil
		}
		groupIDs = alb.SecurityGroups

	case classic:
		clb, err := describeClassicLoadBalancer(ctx, elbAPI, loadBalancerName)

		if err != nil {
			return nil, fmt.Errorf("cannot describe ELB: %w", err)
		}
		if clb == nil {
			// The load balancer wasn't found
			return map[string]struct{}{}, nil
		}
		groupIDs = clb.SecurityGroups
	}

	sgResponse, err := describeSecurityGroupsByID(ctx, ec2API, groupIDs)

	if err != nil {
		return nil, fmt.Errorf("error obtaining security groups for ELB: %s", err)
	}

	result := map[string]struct{}{}

	for _, sg := range sgResponse.SecurityGroups {
		sgID := aws.ToString(sg.GroupId)

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
	if service.Annotations[serviceAnnotationLoadBalancerType] == "nlb" {
		return network
	}
	return classic
}

func loadBalancerExists(ctx context.Context, ec2API awsapi.EC2, elbAPI DescribeLoadBalancersAPI, elbv2API DescribeLoadBalancersAPIV2, lb loadBalancer) (bool, error) {
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

func elbExists(ctx context.Context, elbAPI DescribeLoadBalancersAPI, elbv2API DescribeLoadBalancersAPIV2,
	name string, kind loadBalancerKind) (bool, error) {
	if kind == network || kind == application {
		return elbV2Exists(ctx, elbv2API, name)
	}
	desc, err := describeClassicLoadBalancer(ctx, elbAPI, name)
	return desc != nil, err
}

func describeApplicationLoadBalancer(ctx context.Context, elbv2API DescribeLoadBalancersAPIV2,
	name string) (*elbv2types.LoadBalancer, error) {

	response, err := elbv2API.DescribeLoadBalancers(ctx, &elasticloadbalancingv2.DescribeLoadBalancersInput{
		Names: []string{name},
	})
	if err != nil {
		if isELBv2NotFoundErr(err) {
			return nil, nil
		}
		return nil, err
	}

	var ret elbv2types.LoadBalancer
	switch {
	case len(response.LoadBalancers) > 1:
		logger.Warning("found multiple load balancers with name: %s", name)
		fallthrough
	case len(response.LoadBalancers) > 0:
		ret = response.LoadBalancers[0]
	}
	return &ret, nil
}

func describeClassicLoadBalancer(ctx context.Context, elbAPI DescribeLoadBalancersAPI,
	name string) (*elbtypes.LoadBalancerDescription, error) {

	response, err := elbAPI.DescribeLoadBalancers(ctx, &elasticloadbalancing.DescribeLoadBalancersInput{
		LoadBalancerNames: []string{name},
	})
	if err != nil {
		if isELBNotFoundErr(err) {
			return nil, nil
		}
		return nil, err
	}

	var ret elbtypes.LoadBalancerDescription
	switch {
	case len(response.LoadBalancerDescriptions) > 1:
		logger.Warning("found multiple load balancers with name: %s", name)
		fallthrough
	case len(response.LoadBalancerDescriptions) > 0:
		ret = response.LoadBalancerDescriptions[0]
	}
	return &ret, nil
}

func elbV2Exists(ctx context.Context, api DescribeLoadBalancersAPIV2, name string) (bool, error) {
	_, err := api.DescribeLoadBalancers(ctx, &elasticloadbalancingv2.DescribeLoadBalancersInput{
		Names: []string{name},
	})
	if err != nil {
		var notFoundErr *elbv2types.LoadBalancerNotFoundException
		if errors.As(err, &notFoundErr) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func isELBNotFoundErr(err error) bool {
	var notFoundErr *elbtypes.AccessPointNotFoundException
	return errors.As(err, &notFoundErr)
}

func isELBv2NotFoundErr(err error) bool {
	var notFoundErr *elbv2types.LoadBalancerNotFoundException
	return errors.As(err, &notFoundErr)
}
