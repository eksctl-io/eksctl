package v1alpha5_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("Endpoint Service", func() {
	type endpointServiceEntry struct {
		controlPlaneOnOutposts   bool
		expectedEndpointServices []string
	}

	DescribeTable("Required endpoint services", func(e endpointServiceEntry) {
		endpointServices := api.RequiredEndpointServices(e.controlPlaneOnOutposts)
		endpointServiceNames := make([]string, 0, len(endpointServices))
		for _, es := range endpointServices {
			endpointServiceNames = append(endpointServiceNames, es.Name)
		}
		Expect(endpointServiceNames).To(Equal(e.expectedEndpointServices))
	},
		Entry("standard cluster", endpointServiceEntry{
			controlPlaneOnOutposts: false,
			expectedEndpointServices: []string{
				"ec2",
				"ecr.api",
				"ecr.dkr",
				"s3",
				"sts",
			},
		}),
		Entry("Outposts", endpointServiceEntry{
			controlPlaneOnOutposts: true,
			expectedEndpointServices: []string{
				"ec2",
				"ecr.api",
				"ecr.dkr",
				"s3",
				"sts",
				"ssm",
				"ssmmessages",
				"ec2messages",
				"secretsmanager",
			},
		}),
	)

	type optionalEndpointEntry struct {
		endpointServiceNames     []string
		cloudWatchLoggingEnabled bool

		expectedEndpointServiceNames []string
		expectedErr                  string
	}

	DescribeTable("Map optional endpoint services", func(e optionalEndpointEntry) {
		endpointServices, err := api.MapOptionalEndpointServices(e.endpointServiceNames, e.cloudWatchLoggingEnabled)
		if e.expectedErr != "" {
			Expect(err).To(MatchError(e.expectedErr))
			return
		}
		endpointServiceNames := make([]string, 0, len(endpointServices))
		for _, es := range endpointServices {
			endpointServiceNames = append(endpointServiceNames, es.Name)
		}
		Expect(endpointServiceNames).To(Equal(e.expectedEndpointServiceNames))
	},
		Entry("valid endpoint services", optionalEndpointEntry{
			endpointServiceNames: []string{
				"cloudformation",
				"autoscaling",
				"logs",
			},
			expectedEndpointServiceNames: []string{
				"cloudformation",
				"autoscaling",
				"logs",
			},
		}),
		Entry("invalid endpoint services", optionalEndpointEntry{
			endpointServiceNames: []string{
				"cloudformation",
				"glue",
				"logs",
			},
			expectedErr: `invalid optional endpoint service: "glue"`,
		}),
		Entry("CloudWatch logging enabled", optionalEndpointEntry{
			endpointServiceNames: []string{
				"cloudformation",
				"autoscaling",
			},
			cloudWatchLoggingEnabled: true,
			expectedEndpointServiceNames: []string{
				"cloudformation",
				"autoscaling",
				"logs",
			},
		}),
		Entry("CloudWatch logging enabled in both", optionalEndpointEntry{
			endpointServiceNames: []string{
				"cloudformation",
				"autoscaling",
				"logs",
			},
			cloudWatchLoggingEnabled: true,
			expectedEndpointServiceNames: []string{
				"cloudformation",
				"autoscaling",
				"logs",
			},
		}),
	)

})
