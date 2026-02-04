package elb

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ELB Cleanup", func() {
	When("Getting ingress load balancers to clean up", func() {
		It("should get the right ELB hostname", func() {

			testCases := []struct {
				hostname string
				expected string
			}{
				{
					hostname: "bf647c9e-default-appingres-350b-1622159649.eu-central-1.elb.amazonaws.com",
					expected: "bf647c9e-default-appingres-350b",
				},
				{
					hostname: "internal-k8s-default-testcase-93d8419948-541000771.us-west-2.elb.amazonaws.com",
					expected: "k8s-default-testcase-93d8419948",
				},
				{
					hostname: "abcdefgh-default-bugfixtest-001-1356525548.us-west-2.elb.amazonaws.com",
					expected: "abcdefgh-default-bugfixtest-001",
				},
				{
					hostname: "internal-abcdefgh-default-bugfixtest-002-541910110.us-west-2.elb.amazonaws.com",
					expected: "abcdefgh-default-bugfixtest-002",
				},
				{
					hostname: "abcdefghijklmnopqrstuvwxyz012345-2118752702.us-west-2.elb.amazonaws.com",
					expected: "abcdefghijklmnopqrstuvwxyz012345",
				},
				{
					hostname: "internal-abcdefghijklmnopqrstuvwxyz999999-67491582.us-west-2.elb.amazonaws.com",
					expected: "abcdefghijklmnopqrstuvwxyz999999",
				},
				{
					hostname: "k8s-default-testcase-98cdbf582b-1474733506.us-west-2.elb.amazonaws.com",
					expected: "k8s-default-testcase-98cdbf582b",
				},
				{
					hostname: "internal-k8s-default-testcase-fb10378931-824853021.us-west-2.elb.amazonaws.com",
					expected: "k8s-default-testcase-fb10378931",
				},
				{
					hostname: "abcdefghijklmnopqrstuvw000-1623371943.us-west-2.elb.amazonaws.com",
					expected: "abcdefghijklmnopqrstuvw000",
				},
				{
					hostname: "internal-abcdefghijklmnopqrstuvw001-774959707.us-west-2.elb.amazonaws.com",
					expected: "abcdefghijklmnopqrstuvw001",
				},
				{
					hostname: "myloadbalancer-1234567890.us-west-2.elb.amazonaws.com",
					expected: "myloadbalancer",
				},
				{
					hostname: "my-loadbalancer-1234567890.us-west-2.elb.amazonaws.com",
					expected: "my-loadbalancer",
				},
			}

			for _, tc := range testCases {
				name, err := getIngressELBName([]string{tc.hostname})
				Expect(err).NotTo(HaveOccurred())
				Expect(name).To(Equal(tc.expected))
			}
		})
	})

	When("Getting ingress load balancers but cannot get the hostname", func() {
		It("should have error", func() {

			testCases := []struct {
				hostname string
			}{
				{
					hostname: "",
				},
				{
					hostname: ".us-east-1.elb.amazonaws.com",
				},
			}

			for _, tc := range testCases {
				_, err := getIngressELBName([]string{tc.hostname})
				Expect(err).To(HaveOccurred(), "Expected an error for hostname: %s", tc.hostname)
			}
		})
	})

	When("Getting ingress load balancers but parsed ELB name exceeds 32 characters", func() {
		It("should have error", func() {

			testCases := []struct {
				hostname string
			}{
				{
					hostname: "this-is-not-an-expected-elb-resource-name.us-east-1.elb.amazonaws.com",
				},
				{
					hostname: "this-is-not-an-expected-elb-resource-name-1234567890.us-east-1.elb.amazonaws.com",
				},
				{
					hostname: "internal-this-is-not-an-expected-elb-resource-name-1234567890.us-east-1.elb.amazonaws.com",
				},
			}

			for _, tc := range testCases {
				_, err := getIngressELBName([]string{tc.hostname})
				Expect(err).To(HaveOccurred(), "Expected an error for hostname: %s", tc.hostname)
			}
		})
	})

	When("Verifying security group cleanup works for Gateway API load balancers", func() {
		It("should handle Gateway API load balancers the same as Ingress ALBs", func() {
			// Gateway API ALBs use the same naming pattern as Ingress ALBs

			testCases := []struct {
				description string
				hostname    string
				expected    string
			}{
				{
					description: "Gateway LB with standard format",
					hostname:    "k8s-default-mygateway-abc123-1234567890.us-west-2.elb.amazonaws.com",
					expected:    "k8s-default-mygateway-abc123",
				},
				{
					description: "Internal Gateway LB",
					hostname:    "internal-k8s-default-gateway-xyz789-987654321.us-west-2.elb.amazonaws.com",
					expected:    "k8s-default-gateway-xyz789",
				},
				{
					description: "Gateway LB with namespace and gateway name",
					hostname:    "k8s-prod-api-gw-hash1234-1111111111.eu-central-1.elb.amazonaws.com",
					expected:    "k8s-prod-api-gw-hash1234",
				},
			}

			for _, tc := range testCases {
				// Test that getGatewayLBName produces the expected name
				name, err := getGatewayLBName([]string{tc.hostname})
				Expect(err).NotTo(HaveOccurred(), "Failed for: %s", tc.description)
				Expect(name).To(Equal(tc.expected), "Failed for: %s", tc.description)

				// Verify the name is also compatible with getIngressELBName
				// This demonstrates that Gateway and Ingress ALBs use the same naming pattern
				ingressName, err := getIngressELBName([]string{tc.hostname})
				Expect(err).NotTo(HaveOccurred(), "Failed for: %s", tc.description)
				Expect(ingressName).To(Equal(name), "Gateway and Ingress name parsing should produce identical results for: %s", tc.description)
			}
		})
	})
})
