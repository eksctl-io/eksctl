package elb

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elb/elbiface"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubernetes/pkg/cloudprovider"
	awsprovider "k8s.io/kubernetes/pkg/cloudprovider/providers/aws"

	"github.com/weaveworks/eksctl/pkg/eks"
)

const (
	elbClassic = iota
	elbV2      = iota
)

// CleanupLoadBalancers finds and deletes any dangling ELBs attached to a service
func CleanupLoadBalancers(elbapi elbiface.ELBAPI, elbv2api elbv2iface.ELBV2API, client *eks.Client) error {
	kubernetesCS, err := client.NewClientSet()
	if err != nil {
		return err
	}
	services, err := kubernetesCS.CoreV1().Services("").List(metav1.ListOptions{})
	if err != nil {
		return err
	}

	// Delete Services of type 'LoadBalancer'
	elbs := map[string]int{}
	var deletionErr error
	for _, s := range services.Items {
		if s.Spec.Type == corev1.ServiceTypeLoadBalancer {
			currentErr := kubernetesCS.CoreV1().Services(s.Namespace).Delete(s.Name, &metav1.DeleteOptions{})
			// Report any of the deletion errors if they happen, but go through the full list of services anyways
			if currentErr != nil {
				deletionErr = currentErr
			} else {
				elbs[cloudprovider.DefaultLoadBalancerName(&s)] = getELBType(&s)
			}
		}
	}

	// Wait for all the ELBs back the LoadBalancer services to disappear, for a maximum of 10 minutes
	waitDuration := 10 * time.Minute
	waitDeadline := time.Now().Add(waitDuration)
	ctx, cleanup := context.WithDeadline(context.Background(), waitDeadline)
	defer cleanup()
	for time.Now().Before(waitDeadline) && len(elbs) > 0 {
		for name, kind := range elbs {
			exists, err := elbExists(ctx, elbapi, elbv2api, name, kind)
			if err == nil && !exists {
				delete(elbs, name)
			}
			time.Sleep(100 * time.Millisecond)
		}
	}

	var finalErr error
	if deletionErr != nil {
		finalErr = fmt.Errorf("error deleting service with load balancer: %s", deletionErr)
	}
	if len(elbs) > 0 {
		errStr := fmt.Sprintf("deadline (%s) surpased waiting for load balancers to be deleted", waitDuration)
		if finalErr != nil {
			errStr = finalErr.Error() + "; " + errStr
		}
		finalErr = errors.New(errStr)
	}

	return nil
}

func getELBType(service *corev1.Service) int {
	// See https://github.com/kubernetes/kubernetes/blob/v1.12.6/pkg/cloudprovider/providers/aws/aws_loadbalancer.go#L51-L56
	if service.Annotations[awsprovider.ServiceAnnotationLoadBalancerType] == "nlb" {
		return elbV2
	}
	return elbClassic
}

func elbExists(ctx context.Context, elbAPI elbiface.ELBAPI, elbv2API elbv2iface.ELBV2API, name string, kind int) (bool, error) {
	if kind == elbV2 {
		return elbV2Exists(ctx, elbv2API, name)
	}
	return classicELBExists(ctx, elbAPI, name)
}

func classicELBExists(ctx context.Context, api elbiface.ELBAPI, name string) (bool, error) {
	request := &elb.DescribeLoadBalancersInput{
		LoadBalancerNames: []*string{&name},
	}

	_, err := api.DescribeLoadBalancersWithContext(ctx, request)
	if err != nil {
		if awsError, ok := err.(awserr.Error); ok {
			if awsError.Code() == "LoadBalancerNotFound" {
				return false, nil
			}
		}
		return false, err
	}

	return true, nil
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
