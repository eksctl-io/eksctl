package elb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kris-nova/logger"

	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayclient "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned"

	"github.com/weaveworks/eksctl/pkg/awsapi"
)

const (
	// AWS Load Balancer Controller names for Gateway API
	awsLBCALBController = "gateway.k8s.aws/alb"
	awsLBCNLBController = "gateway.k8s.aws/nlb"
)

// Gateway interface abstracts over different Gateway API versions
type Gateway interface {
	Delete(ctx context.Context, gwClient gatewayclient.Interface) error
	GetGatewayClassName() string
	GetMetadata() metav1.ObjectMeta
	GetLoadBalancerAddresses() []string
}

// v1Gateway wraps Gateway API v1 Gateway
type v1Gateway struct {
	gateway gatewayv1.Gateway
}

func (g *v1Gateway) Delete(ctx context.Context, gwClient gatewayclient.Interface) error {
	return gwClient.GatewayV1().Gateways(g.gateway.Namespace).
		Delete(ctx, g.gateway.Name, metav1.DeleteOptions{})
}

func (g *v1Gateway) GetGatewayClassName() string {
	return string(g.gateway.Spec.GatewayClassName)
}

func (g *v1Gateway) GetMetadata() metav1.ObjectMeta {
	return g.gateway.ObjectMeta
}

func (g *v1Gateway) GetLoadBalancerAddresses() []string {
	var addresses []string
	for _, addr := range g.gateway.Status.Addresses {
		if addr.Type != nil && *addr.Type == gatewayv1.HostnameAddressType {
			addresses = append(addresses, addr.Value)
		}
	}
	return addresses
}

// listGateway lists all Gateway resources across all namespaces using Gateway API v1
// Returns an empty list if Gateway API CRDs are not installed (no error)
func listGateway(ctx context.Context, restConfig *rest.Config) ([]Gateway, gatewayclient.Interface, error) {
	gwClient, err := gatewayclient.NewForConfig(restConfig)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create gateway client: %w", err)
	}

	logger.Debug("using v1 Gateway API")

	gateways, err := gwClient.GatewayV1().Gateways(metav1.NamespaceAll).List(ctx, metav1.ListOptions{})
	if err != nil {
		// Check if this is a "not found" error (CRDs not installed)
		if isGatewayAPINotFoundErr(err) {
			logger.Debug("Gateway API v1 CRDs not found, skipping Gateway cleanup")
			return []Gateway{}, gwClient, nil
		}
		return nil, nil, err
	}

	var gatewayList []Gateway
	for i := range gateways.Items {
		gatewayList = append(gatewayList, &v1Gateway{gateway: gateways.Items[i]})
	}
	return gatewayList, gwClient, nil
}

// isGatewayAPINotFoundErr checks if the error indicates Gateway API CRDs are not installed
func isGatewayAPINotFoundErr(err error) bool {
	if err == nil {
		return false
	}

	// Check for NotFound status error
	if k8serrors.IsNotFound(err) {
		return true
	}

	// Check for common error messages when CRDs are not installed
	errStr := err.Error()
	return strings.Contains(errStr, "not found") ||
		strings.Contains(errStr, "could not find the requested resource") ||
		strings.Contains(errStr, "no matches for kind") ||
		strings.Contains(errStr, "the server could not find the requested resource")
}

// getGatewayClass retrieves a GatewayClass by name and returns its controller name
// Returns empty string if the GatewayClass is not found (no error)
func getGatewayClass(ctx context.Context, gwClient gatewayclient.Interface, gatewayClassName string) (string, error) {
	if gatewayClassName == "" {
		return "", nil
	}

	gatewayClass, err := gwClient.GatewayV1().GatewayClasses().Get(ctx, gatewayClassName, metav1.GetOptions{})
	if err != nil {
		// If GatewayClass not found, return empty string (no error)
		if k8serrors.IsNotFound(err) {
			logger.Debug("GatewayClass %q not found", gatewayClassName)
			return "", nil
		}
		return "", fmt.Errorf("failed to get GatewayClass %q: %w", gatewayClassName, err)
	}
	return string(gatewayClass.Spec.ControllerName), nil
}

// isAWSLoadBalancerController checks if the controller name matches AWS LBC patterns
func isAWSLoadBalancerController(controllerName string) bool {
	return controllerName == awsLBCALBController || controllerName == awsLBCNLBController
}

// getGatewayLBName parses the load balancer name from Gateway DNS addresses
// Gateway load balancers follow a similar naming pattern to Ingress ALBs:
// k8s-<namespace>-<gateway>-<hash>.<region>.elb.amazonaws.com
// The load balancer name is extracted by removing the region and domain suffix, and the hash suffix
func getGatewayLBName(addresses []string) (string, error) {
	if len(addresses) == 0 {
		return "", fmt.Errorf("no addresses provided")
	}

	// Expected format: k8s-namespace-gateway-hash.region.elb.amazonaws.com
	// or internal-k8s-namespace-gateway-hash.region.elb.amazonaws.com for internal load balancers
	hostNameParts := strings.Split(addresses[0], ".")
	if len(hostNameParts) == 0 || len(hostNameParts[0]) == 0 {
		return "", fmt.Errorf("cannot get the hostname: %v", hostNameParts)
	}

	name := strings.TrimPrefix(hostNameParts[0], "internal-")

	idIdx := strings.LastIndex(name, "-")
	if idIdx != -1 {
		name = name[:idIdx]
	}

	// AWS load balancer names cannot exceed 32 characters
	if len(name) > 32 {
		return "", fmt.Errorf("parsed name exceeds maximum of 32 characters: %s", name)
	}

	return name, nil
}

// getGatewayLoadBalancer extracts load balancer information from a Gateway resource
// Returns nil if the Gateway is not managed by AWS LBC or has not been provisioned
func getGatewayLoadBalancer(ctx context.Context, ec2API awsapi.EC2, elbAPI DescribeLoadBalancersAPI,
	elbv2API DescribeLoadBalancersAPIV2, gwClient gatewayclient.Interface, clusterName string,
	gateway Gateway) (*loadBalancer, error) {

	metadata := gateway.GetMetadata()
	gatewayClassName := gateway.GetGatewayClassName()

	// Get the GatewayClass to check if it's managed by AWS LBC
	controllerName, err := getGatewayClass(ctx, gwClient, gatewayClassName)
	if err != nil {
		return nil, fmt.Errorf("cannot get GatewayClass %q: %w", gatewayClassName, err)
	}

	// Skip Gateways not managed by AWS Load Balancer Controller
	if !isAWSLoadBalancerController(controllerName) {
		logger.Debug("Gateway %s/%s uses controller %q, not AWS LBC, skip",
			metadata.Namespace, metadata.Name, controllerName)
		return nil, nil
	}

	// Check if the Gateway has been provisioned (status.addresses populated)
	addresses := gateway.GetLoadBalancerAddresses()
	if len(addresses) == 0 {
		logger.Debug("Gateway %s/%s is managed by AWS LBC, but not provisioned yet, skip",
			metadata.Namespace, metadata.Name)
		return nil, nil
	}

	// Parse the load balancer name from the DNS address
	name, err := getGatewayLBName(addresses)
	if err != nil {
		logger.Debug("Gateway %s/%s is managed by AWS LBC, but cannot parse load balancer name, skip: %s",
			metadata.Namespace, metadata.Name, err)
		return nil, nil
	}

	logger.Debug("Gateway load balancer resource name: %s", name)

	// Determine load balancer kind based on controller name
	var kind loadBalancerKind
	switch controllerName {
	case awsLBCALBController:
		kind = application
	case awsLBCNLBController:
		kind = network
	default:
		// This should not happen due to isAWSLoadBalancerController check above
		return nil, fmt.Errorf("unexpected AWS LBC controller name: %s", controllerName)
	}

	// Retrieve security groups owned by the load balancer
	ctx, cleanup := context.WithTimeout(ctx, 30*time.Second)
	defer cleanup()
	securityGroupIDs, err := getSecurityGroupsOwnedByLoadBalancer(ctx, ec2API, elbAPI, elbv2API, clusterName, name, kind)
	if err != nil {
		return nil, fmt.Errorf("cannot obtain security groups for Gateway load balancer %s: %w", name, err)
	}

	return &loadBalancer{
		name:                  name,
		kind:                  kind,
		ownedSecurityGroupIDs: securityGroupIDs,
	}, nil
}
