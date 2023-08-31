package elb

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ELB Cleanup", func() {
	When("Getting ingress load balancers to clean up", func() {
		It("should get the right ELB hostname", func() {

			testELBDNSnames := [][]string{
				{"bf647c9e-default-appingres-350b-1622159649.eu-central-1.elb.amazonaws.com", ""},
				{"internal-k8s-default-testcase-93d8419948-541000771.us-west-2.elb.amazonaws.com", ""},
				{"abcdefgh-default-bugfixtest-001-1356525548.us-west-2.elb.amazonaws.com", ""},
				{"internal-abcdefgh-default-bugfixtest-002-541910110.us-west-2.elb.amazonaws.com", ""},
				{"abcdefghijklmnopqrstuvwxyz012345-2118752702.us-west-2.elb.amazonaws.com", ""},
				{"internal-abcdefghijklmnopqrstuvwxyz999999-67491582.us-west-2.elb.amazonaws.com", ""},
				{"k8s-default-testcase-98cdbf582b-1474733506.us-west-2.elb.amazonaws.com", ""},
				{"internal-k8s-default-testcase-fb10378931-824853021.us-west-2.elb.amazonaws.com", ""},
				{"abcdefghijklmnopqrstuvw000-1623371943.us-west-2.elb.amazonaws.com", ""},
				{"internal-abcdefghijklmnopqrstuvw001-774959707.us-west-2.elb.amazonaws.com", ""},
				{"myloadbalancer-1234567890.us-west-2.elb.amazonaws.com", ""},
				{"my-loadbalancer-1234567890.us-west-2.elb.amazonaws.com", ""},
			}
			expectELBName := []string{
				"bf647c9e-default-appingres-350b",
				"k8s-default-testcase-93d8419948",
				"abcdefgh-default-bugfixtest-001",
				"abcdefgh-default-bugfixtest-002",
				"abcdefghijklmnopqrstuvwxyz012345",
				"abcdefghijklmnopqrstuvwxyz999999",
				"k8s-default-testcase-98cdbf582b",
				"k8s-default-testcase-fb10378931",
				"abcdefghijklmnopqrstuvw000",
				"abcdefghijklmnopqrstuvw001",
				"myloadbalancer",
				"my-loadbalancer",
			}

			for i, elbDNSName := range testELBDNSnames {
				name := getIngressELBName(context.TODO(), "", elbDNSName)
				Expect(name).To(Equal(expectELBName[i]))
			}
		})
	})
})
