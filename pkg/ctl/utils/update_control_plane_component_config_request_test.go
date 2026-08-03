package utils

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go-v2/aws"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

var _ = Describe("control plane component config update requests", func() {

	Describe("makeKubeAPIServerConfigRequest", func() {
		type entry struct {
			config   *api.KubeAPIServerConfig
			expected *ekstypes.KubeApiServerConfigRequest
		}

		DescribeTable("converts the config", func(e entry) {
			Expect(makeKubeAPIServerConfigRequest(e.config)).To(Equal(e.expected))
		},
			Entry("nil config", entry{
				config:   nil,
				expected: nil,
			}),
			// An empty config is reachable: the loader only requires that one of the three
			// components is present, so `kubeAPIServerConfig: {}` passes validation.
			Entry("empty config is not sent", entry{
				config:   &api.KubeAPIServerConfig{},
				expected: nil,
			}),
			Entry("empty service node port range is not sent", entry{
				config:   &api.KubeAPIServerConfig{ServiceNodePortRange: &api.ServiceNodePortRange{}},
				expected: nil,
			}),
			Entry("eventTTL only", entry{
				config: &api.KubeAPIServerConfig{EventTTL: aws.String("30m")},
				expected: &ekstypes.KubeApiServerConfigRequest{
					EventTtl: aws.String("30m"),
				},
			}),
			// A partially specified range is passed through so that EKS can reject it with a
			// clear error, rather than silently dropping what the user wrote. The unset end
			// is a zero int32, which the SDK serializer omits from the request body.
			Entry("min port only is sent as a partial range", entry{
				config: &api.KubeAPIServerConfig{
					ServiceNodePortRange: &api.ServiceNodePortRange{MinPort: aws.Int(30000)},
				},
				expected: &ekstypes.KubeApiServerConfigRequest{
					ServiceNodePortRange: &ekstypes.ServiceNodePortRange{MinPort: 30000},
				},
			}),
			Entry("max port only is sent as a partial range", entry{
				config: &api.KubeAPIServerConfig{
					ServiceNodePortRange: &api.ServiceNodePortRange{MaxPort: aws.Int(32767)},
				},
				expected: &ekstypes.KubeApiServerConfigRequest{
					ServiceNodePortRange: &ekstypes.ServiceNodePortRange{MaxPort: 32767},
				},
			}),
			Entry("eventTTL is sent alongside a partial port range", entry{
				config: &api.KubeAPIServerConfig{
					EventTTL:             aws.String("30m"),
					ServiceNodePortRange: &api.ServiceNodePortRange{MinPort: aws.Int(30000)},
				},
				expected: &ekstypes.KubeApiServerConfigRequest{
					EventTtl:             aws.String("30m"),
					ServiceNodePortRange: &ekstypes.ServiceNodePortRange{MinPort: 30000},
				},
			}),
			Entry("all values set", entry{
				config: &api.KubeAPIServerConfig{
					EventTTL: aws.String("30m"),
					ServiceNodePortRange: &api.ServiceNodePortRange{
						MinPort: aws.Int(30000),
						MaxPort: aws.Int(32767),
					},
				},
				expected: &ekstypes.KubeApiServerConfigRequest{
					EventTtl: aws.String("30m"),
					ServiceNodePortRange: &ekstypes.ServiceNodePortRange{
						MinPort: 30000,
						MaxPort: 32767,
					},
				},
			}),
		)
	})

	Describe("makeKubeSchedulerConfigRequest", func() {
		type entry struct {
			config   *api.KubeSchedulerConfig
			expected *ekstypes.KubeSchedulerConfigRequest
		}

		DescribeTable("converts the config", func(e entry) {
			Expect(makeKubeSchedulerConfigRequest(e.config)).To(Equal(e.expected))
		},
			Entry("nil config", entry{
				config:   nil,
				expected: nil,
			}),
			Entry("empty config is not sent", entry{
				config:   &api.KubeSchedulerConfig{},
				expected: nil,
			}),
			Entry("empty node resources fit is not sent", entry{
				config:   &api.KubeSchedulerConfig{NodeResourcesFit: &api.NodeResourcesFitConfig{}},
				expected: nil,
			}),
			Entry("empty scoring strategy is not sent", entry{
				config: &api.KubeSchedulerConfig{
					NodeResourcesFit: &api.NodeResourcesFitConfig{
						ScoringStrategy: &api.ScoringStrategy{},
					},
				},
				expected: nil,
			}),
			Entry("scoring strategy type only", entry{
				config: &api.KubeSchedulerConfig{
					NodeResourcesFit: &api.NodeResourcesFitConfig{
						ScoringStrategy: &api.ScoringStrategy{Type: aws.String("MostAllocated")},
					},
				},
				expected: &ekstypes.KubeSchedulerConfigRequest{
					NodeResourcesFit: &ekstypes.NodeResourcesFitConfig{
						ScoringStrategy: &ekstypes.ScoringStrategy{
							Type: ekstypes.ScoringStrategyType("MostAllocated"),
						},
					},
				},
			}),
			// Nothing rejects a resource entry with only one of the two fields set, so each
			// is sent on its own rather than dereferenced unconditionally.
			Entry("resource weight with no weight", entry{
				config: &api.KubeSchedulerConfig{
					NodeResourcesFit: &api.NodeResourcesFitConfig{
						ScoringStrategy: &api.ScoringStrategy{
							Resources: []api.ResourceWeight{{Name: aws.String("cpu")}},
						},
					},
				},
				expected: &ekstypes.KubeSchedulerConfigRequest{
					NodeResourcesFit: &ekstypes.NodeResourcesFitConfig{
						ScoringStrategy: &ekstypes.ScoringStrategy{
							Resources: []ekstypes.ResourceWeight{
								{Name: aws.String("cpu")},
							},
						},
					},
				},
			}),
			Entry("resource weight with no name", entry{
				config: &api.KubeSchedulerConfig{
					NodeResourcesFit: &api.NodeResourcesFitConfig{
						ScoringStrategy: &api.ScoringStrategy{
							Resources: []api.ResourceWeight{{Weight: aws.Int(5)}},
						},
					},
				},
				expected: &ekstypes.KubeSchedulerConfigRequest{
					NodeResourcesFit: &ekstypes.NodeResourcesFitConfig{
						ScoringStrategy: &ekstypes.ScoringStrategy{
							Resources: []ekstypes.ResourceWeight{
								{Weight: aws.Int32(5)},
							},
						},
					},
				},
			}),
			// An entry with neither field still counts as set, so the strategy is sent and EKS
			// reports the problem rather than eksctl silently discarding the list.
			Entry("empty resource weight entry is still sent", entry{
				config: &api.KubeSchedulerConfig{
					NodeResourcesFit: &api.NodeResourcesFitConfig{
						ScoringStrategy: &api.ScoringStrategy{
							Resources: []api.ResourceWeight{{}},
						},
					},
				},
				expected: &ekstypes.KubeSchedulerConfigRequest{
					NodeResourcesFit: &ekstypes.NodeResourcesFitConfig{
						ScoringStrategy: &ekstypes.ScoringStrategy{
							Resources: []ekstypes.ResourceWeight{{}},
						},
					},
				},
			}),
			Entry("type and resource weights", entry{
				config: &api.KubeSchedulerConfig{
					NodeResourcesFit: &api.NodeResourcesFitConfig{
						ScoringStrategy: &api.ScoringStrategy{
							Type: aws.String("LeastAllocated"),
							Resources: []api.ResourceWeight{
								{Name: aws.String("cpu"), Weight: aws.Int(3)},
								{Name: aws.String("memory"), Weight: aws.Int(1)},
							},
						},
					},
				},
				expected: &ekstypes.KubeSchedulerConfigRequest{
					NodeResourcesFit: &ekstypes.NodeResourcesFitConfig{
						ScoringStrategy: &ekstypes.ScoringStrategy{
							Type: ekstypes.ScoringStrategyType("LeastAllocated"),
							Resources: []ekstypes.ResourceWeight{
								{Name: aws.String("cpu"), Weight: aws.Int32(3)},
								{Name: aws.String("memory"), Weight: aws.Int32(1)},
							},
						},
					},
				},
			}),
		)
	})

	Describe("makeKubeControllerManagerConfigRequest", func() {
		type entry struct {
			config   *api.KubeControllerManagerConfig
			expected *ekstypes.KubeControllerManagerConfigRequest
		}

		DescribeTable("converts the config", func(e entry) {
			Expect(makeKubeControllerManagerConfigRequest(e.config)).To(Equal(e.expected))
		},
			Entry("nil config", entry{
				config:   nil,
				expected: nil,
			}),
			Entry("empty config is not sent", entry{
				config:   &api.KubeControllerManagerConfig{},
				expected: nil,
			}),
			Entry("empty horizontal pod autoscaler config is not sent", entry{
				config: &api.KubeControllerManagerConfig{
					HorizontalPodAutoscalerControllerConfig: &api.HorizontalPodAutoscalerControllerConfig{},
				},
				expected: nil,
			}),
			Entry("sync period set", entry{
				config: &api.KubeControllerManagerConfig{
					HorizontalPodAutoscalerControllerConfig: &api.HorizontalPodAutoscalerControllerConfig{
						HorizontalPodAutoscalerSyncPeriod: aws.String("15s"),
					},
				},
				expected: &ekstypes.KubeControllerManagerConfigRequest{
					HorizontalPodAutoscalerControllerConfig: &ekstypes.HorizontalPodAutoscalerControllerConfigRequest{
						HorizontalPodAutoscalerSyncPeriod: aws.String("15s"),
					},
				},
			}),
		)
	})
})
