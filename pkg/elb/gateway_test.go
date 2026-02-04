package elb

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	gatewayfake "sigs.k8s.io/gateway-api/pkg/client/clientset/versioned/fake"
)

var _ = Describe("Gateway API", func() {
	Describe("Gateway Interface", func() {
		Describe("v1Gateway", func() {
			It("should return correct gateway class name", func() {
				className := gatewayv1.ObjectName("test-class")
				gateway := &v1Gateway{
					gateway: gatewayv1.Gateway{
						Spec: gatewayv1.GatewaySpec{
							GatewayClassName: className,
						},
					},
				}
				Expect(gateway.GetGatewayClassName()).To(Equal("test-class"))
			})

			It("should return correct metadata", func() {
				gateway := &v1Gateway{
					gateway: gatewayv1.Gateway{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "test-gateway",
							Namespace: "default",
						},
					},
				}
				meta := gateway.GetMetadata()
				Expect(meta.Name).To(Equal("test-gateway"))
				Expect(meta.Namespace).To(Equal("default"))
			})

			It("should return load balancer addresses for hostname type", func() {
				addrType := gatewayv1.HostnameAddressType
				gateway := &v1Gateway{
					gateway: gatewayv1.Gateway{
						Status: gatewayv1.GatewayStatus{
							Addresses: []gatewayv1.GatewayStatusAddress{
								{
									Type:  &addrType,
									Value: "test.elb.amazonaws.com",
								},
							},
						},
					},
				}
				addresses := gateway.GetLoadBalancerAddresses()
				Expect(addresses).To(HaveLen(1))
				Expect(addresses[0]).To(Equal("test.elb.amazonaws.com"))
			})

			It("should filter out non-hostname address types", func() {
				hostnameType := gatewayv1.HostnameAddressType
				ipType := gatewayv1.IPAddressType
				gateway := &v1Gateway{
					gateway: gatewayv1.Gateway{
						Status: gatewayv1.GatewayStatus{
							Addresses: []gatewayv1.GatewayStatusAddress{
								{
									Type:  &hostnameType,
									Value: "test.elb.amazonaws.com",
								},
								{
									Type:  &ipType,
									Value: "192.168.1.1",
								},
							},
						},
					},
				}
				addresses := gateway.GetLoadBalancerAddresses()
				Expect(addresses).To(HaveLen(1))
				Expect(addresses[0]).To(Equal("test.elb.amazonaws.com"))
			})

			It("should return empty slice when no hostname addresses exist", func() {
				gateway := &v1Gateway{
					gateway: gatewayv1.Gateway{
						Status: gatewayv1.GatewayStatus{
							Addresses: []gatewayv1.GatewayStatusAddress{},
						},
					},
				}
				addresses := gateway.GetLoadBalancerAddresses()
				Expect(addresses).To(BeEmpty())
			})
		})
	})

	Describe("listGateway", func() {
		Context("graceful handling of missing CRDs", func() {
			It("should detect CRD not found errors", func() {
				testCases := []struct {
					errorMsg string
					expected bool
				}{
					{"not found", true},
					{"could not find the requested resource", true},
					{"no matches for kind", true},
					{"the server could not find the requested resource", true},
					{"some other error", false},
					{"", false},
				}

				for _, tc := range testCases {
					var err error
					if tc.errorMsg != "" {
						err = fmt.Errorf("%s", tc.errorMsg)
					}
					result := isGatewayAPINotFoundErr(err)
					Expect(result).To(Equal(tc.expected), "Error message: %s", tc.errorMsg)
				}
			})

			It("should return false for nil error", func() {
				Expect(isGatewayAPINotFoundErr(nil)).To(BeFalse())
			})
		})
	})

	Describe("getGatewayClass", func() {
		var (
			ctx      context.Context
			gwClient *gatewayfake.Clientset
		)

		BeforeEach(func() {
			ctx = context.Background()
		})

		It("should return AWS LBC ALB controller name", func() {
			gatewayClass := &gatewayv1.GatewayClass{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-alb",
				},
				Spec: gatewayv1.GatewayClassSpec{
					ControllerName: awsLBCALBController,
				},
			}
			gwClient = gatewayfake.NewSimpleClientset(gatewayClass)

			controllerName, err := getGatewayClass(ctx, gwClient, "aws-alb")
			Expect(err).NotTo(HaveOccurred())
			Expect(controllerName).To(Equal(awsLBCALBController))
		})

		It("should return AWS LBC NLB controller name", func() {
			gatewayClass := &gatewayv1.GatewayClass{
				ObjectMeta: metav1.ObjectMeta{
					Name: "aws-nlb",
				},
				Spec: gatewayv1.GatewayClassSpec{
					ControllerName: awsLBCNLBController,
				},
			}
			gwClient = gatewayfake.NewSimpleClientset(gatewayClass)

			controllerName, err := getGatewayClass(ctx, gwClient, "aws-nlb")
			Expect(err).NotTo(HaveOccurred())
			Expect(controllerName).To(Equal(awsLBCNLBController))
		})

		It("should return non-AWS controller name", func() {
			gatewayClass := &gatewayv1.GatewayClass{
				ObjectMeta: metav1.ObjectMeta{
					Name: "istio",
				},
				Spec: gatewayv1.GatewayClassSpec{
					ControllerName: "istio.io/gateway-controller",
				},
			}
			gwClient = gatewayfake.NewSimpleClientset(gatewayClass)

			controllerName, err := getGatewayClass(ctx, gwClient, "istio")
			Expect(err).NotTo(HaveOccurred())
			Expect(controllerName).To(Equal("istio.io/gateway-controller"))
		})

		It("should return empty string when GatewayClass is not found", func() {
			gwClient = gatewayfake.NewSimpleClientset()

			controllerName, err := getGatewayClass(ctx, gwClient, "non-existent")
			Expect(err).NotTo(HaveOccurred())
			Expect(controllerName).To(BeEmpty())
		})

		It("should return empty string when gatewayClassName is empty", func() {
			gwClient = gatewayfake.NewSimpleClientset()

			controllerName, err := getGatewayClass(ctx, gwClient, "")
			Expect(err).NotTo(HaveOccurred())
			Expect(controllerName).To(BeEmpty())
		})
	})

	Describe("isAWSLoadBalancerController", func() {
		It("should return true for AWS LBC ALB controller", func() {
			Expect(isAWSLoadBalancerController(awsLBCALBController)).To(BeTrue())
		})

		It("should return true for AWS LBC NLB controller", func() {
			Expect(isAWSLoadBalancerController(awsLBCNLBController)).To(BeTrue())
		})

		It("should return false for non-AWS controllers", func() {
			testCases := []string{
				"istio.io/gateway-controller",
				"nginx.org/gateway-controller",
				"traefik.io/gateway-controller",
				"kong.io/gateway-controller",
				"",
				"gateway.k8s.aws/other",
				"gateway.k8s.aws",
			}

			for _, controllerName := range testCases {
				Expect(isAWSLoadBalancerController(controllerName)).To(BeFalse(),
					"Controller name: %s should not be identified as AWS LBC", controllerName)
			}
		})
	})

	Describe("getGatewayLBName", func() {
		Context("valid Gateway DNS names", func() {
			It("should parse external ALB DNS name", func() {
				addresses := []string{"k8s-default-testgw-abc123.us-west-2.elb.amazonaws.com"}
				name, err := getGatewayLBName(addresses)
				Expect(err).NotTo(HaveOccurred())
				Expect(name).To(Equal("k8s-default-testgw"))
			})

			It("should parse internal ALB DNS name", func() {
				addresses := []string{"internal-k8s-default-testgw-xyz789.eu-central-1.elb.amazonaws.com"}
				name, err := getGatewayLBName(addresses)
				Expect(err).NotTo(HaveOccurred())
				Expect(name).To(Equal("k8s-default-testgw"))
			})

			It("should parse NLB DNS name", func() {
				addresses := []string{"k8s-kube-system-gateway-def456.us-east-1.elb.amazonaws.com"}
				name, err := getGatewayLBName(addresses)
				Expect(err).NotTo(HaveOccurred())
				Expect(name).To(Equal("k8s-kube-system-gateway"))
			})

			It("should parse internal NLB DNS name", func() {
				addresses := []string{"internal-k8s-prod-api-gw-hash123.ap-southeast-1.elb.amazonaws.com"}
				name, err := getGatewayLBName(addresses)
				Expect(err).NotTo(HaveOccurred())
				Expect(name).To(Equal("k8s-prod-api-gw"))
			})
		})

		Context("error cases", func() {
			It("should return error for empty addresses", func() {
				_, err := getGatewayLBName([]string{})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("no addresses provided"))
			})

			It("should return error for invalid format", func() {
				_, err := getGatewayLBName([]string{""})
				Expect(err).To(HaveOccurred())
			})

			It("should return error for name exceeding 32 characters", func() {
				// Create a very long hostname that would result in >32 char name
				longName := "k8s-verylongnamespace-verylonggatewayname-hash123.us-west-2.elb.amazonaws.com"
				_, err := getGatewayLBName([]string{longName})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("exceeds maximum of 32 characters"))
			})
		})
	})
})
