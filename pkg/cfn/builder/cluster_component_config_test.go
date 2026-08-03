package builder

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go-v2/aws"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	gfneks "github.com/weaveworks/eksctl/pkg/goformation/cloudformation/eks"
	gfnt "github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"
)

// The template-level tests in cluster_test.go cover the fully-populated and entirely-absent
// cases by asserting on the rendered ControlPlane resource. These call the converters
// directly to cover the nil and partially-populated cases in between, mirroring the table in
// pkg/ctl/utils/update_control_plane_component_config_request_test.go.
var _ = Describe("control plane component config template properties", func() {

	Describe("makeKubeAPIServerConfig", func() {
		type entry struct {
			config   *api.KubeAPIServerConfig
			expected *gfneks.Cluster_KubeApiServerConfig
		}

		DescribeTable("converts the config", func(e entry) {
			Expect(makeKubeAPIServerConfig(e.config)).To(Equal(e.expected))
		},
			Entry("nil config", entry{
				config:   nil,
				expected: nil,
			}),
			// An empty config is reachable via the update command, whose loader only requires
			// that one of the three components is present.
			Entry("empty config is omitted", entry{
				config:   &api.KubeAPIServerConfig{},
				expected: nil,
			}),
			Entry("empty service node port range is omitted", entry{
				config:   &api.KubeAPIServerConfig{ServiceNodePortRange: &api.ServiceNodePortRange{}},
				expected: nil,
			}),
			Entry("eventTTL only", entry{
				config: &api.KubeAPIServerConfig{EventTTL: aws.String("30m")},
				expected: &gfneks.Cluster_KubeApiServerConfig{
					EventTtl: gfnt.NewString("30m"),
				},
			}),
			// A partially specified range is emitted so that EKS can reject it with a clear
			// error, rather than silently dropping what the user wrote.
			Entry("min port only is emitted as a partial range", entry{
				config: &api.KubeAPIServerConfig{
					ServiceNodePortRange: &api.ServiceNodePortRange{MinPort: aws.Int(30000)},
				},
				expected: &gfneks.Cluster_KubeApiServerConfig{
					ServiceNodePortRange: &gfneks.Cluster_ServiceNodePortRange{
						MinPort: gfnt.NewInteger(30000),
					},
				},
			}),
			Entry("max port only is emitted as a partial range", entry{
				config: &api.KubeAPIServerConfig{
					ServiceNodePortRange: &api.ServiceNodePortRange{MaxPort: aws.Int(32767)},
				},
				expected: &gfneks.Cluster_KubeApiServerConfig{
					ServiceNodePortRange: &gfneks.Cluster_ServiceNodePortRange{
						MaxPort: gfnt.NewInteger(32767),
					},
				},
			}),
			Entry("eventTTL is emitted alongside a partial port range", entry{
				config: &api.KubeAPIServerConfig{
					EventTTL:             aws.String("30m"),
					ServiceNodePortRange: &api.ServiceNodePortRange{MinPort: aws.Int(30000)},
				},
				expected: &gfneks.Cluster_KubeApiServerConfig{
					EventTtl: gfnt.NewString("30m"),
					ServiceNodePortRange: &gfneks.Cluster_ServiceNodePortRange{
						MinPort: gfnt.NewInteger(30000),
					},
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
				expected: &gfneks.Cluster_KubeApiServerConfig{
					EventTtl: gfnt.NewString("30m"),
					ServiceNodePortRange: &gfneks.Cluster_ServiceNodePortRange{
						MinPort: gfnt.NewInteger(30000),
						MaxPort: gfnt.NewInteger(32767),
					},
				},
			}),
		)
	})

	Describe("makeKubeSchedulerConfig", func() {
		type entry struct {
			config   *api.KubeSchedulerConfig
			expected *gfneks.Cluster_KubeSchedulerConfig
		}

		DescribeTable("converts the config", func(e entry) {
			Expect(makeKubeSchedulerConfig(e.config)).To(Equal(e.expected))
		},
			Entry("nil config", entry{
				config:   nil,
				expected: nil,
			}),
			Entry("empty config is omitted", entry{
				config:   &api.KubeSchedulerConfig{},
				expected: nil,
			}),
			Entry("empty node resources fit is omitted", entry{
				config:   &api.KubeSchedulerConfig{NodeResourcesFit: &api.NodeResourcesFitConfig{}},
				expected: nil,
			}),
			Entry("empty scoring strategy is omitted", entry{
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
				expected: &gfneks.Cluster_KubeSchedulerConfig{
					NodeResourcesFit: &gfneks.Cluster_NodeResourcesFitConfig{
						ScoringStrategy: &gfneks.Cluster_ScoringStrategy{
							Type: gfnt.NewString("MostAllocated"),
						},
					},
				},
			}),
			// Resources alone is enough to emit the strategy, so a weight list with no type
			// is not dropped.
			Entry("resource weights only", entry{
				config: &api.KubeSchedulerConfig{
					NodeResourcesFit: &api.NodeResourcesFitConfig{
						ScoringStrategy: &api.ScoringStrategy{
							Resources: []api.ResourceWeight{
								{Name: aws.String("cpu"), Weight: aws.Int(3)},
							},
						},
					},
				},
				expected: &gfneks.Cluster_KubeSchedulerConfig{
					NodeResourcesFit: &gfneks.Cluster_NodeResourcesFitConfig{
						ScoringStrategy: &gfneks.Cluster_ScoringStrategy{
							Resources: []gfneks.Cluster_ResourceWeight{
								{Name: gfnt.NewString("cpu"), Weight: gfnt.NewInteger(3)},
							},
						},
					},
				},
			}),
			// Nothing rejects a resource entry with only one of the two fields set, so each
			// is emitted on its own rather than dereferenced unconditionally.
			Entry("resource weight with no weight", entry{
				config: &api.KubeSchedulerConfig{
					NodeResourcesFit: &api.NodeResourcesFitConfig{
						ScoringStrategy: &api.ScoringStrategy{
							Resources: []api.ResourceWeight{{Name: aws.String("cpu")}},
						},
					},
				},
				expected: &gfneks.Cluster_KubeSchedulerConfig{
					NodeResourcesFit: &gfneks.Cluster_NodeResourcesFitConfig{
						ScoringStrategy: &gfneks.Cluster_ScoringStrategy{
							Resources: []gfneks.Cluster_ResourceWeight{
								{Name: gfnt.NewString("cpu")},
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
				expected: &gfneks.Cluster_KubeSchedulerConfig{
					NodeResourcesFit: &gfneks.Cluster_NodeResourcesFitConfig{
						ScoringStrategy: &gfneks.Cluster_ScoringStrategy{
							Resources: []gfneks.Cluster_ResourceWeight{
								{Weight: gfnt.NewInteger(5)},
							},
						},
					},
				},
			}),
			// An entry with neither field still counts as set, so the strategy is emitted and
			// EKS reports the problem rather than eksctl silently discarding the list.
			Entry("empty resource weight entry is still emitted", entry{
				config: &api.KubeSchedulerConfig{
					NodeResourcesFit: &api.NodeResourcesFitConfig{
						ScoringStrategy: &api.ScoringStrategy{
							Resources: []api.ResourceWeight{{}},
						},
					},
				},
				expected: &gfneks.Cluster_KubeSchedulerConfig{
					NodeResourcesFit: &gfneks.Cluster_NodeResourcesFitConfig{
						ScoringStrategy: &gfneks.Cluster_ScoringStrategy{
							Resources: []gfneks.Cluster_ResourceWeight{{}},
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
				expected: &gfneks.Cluster_KubeSchedulerConfig{
					NodeResourcesFit: &gfneks.Cluster_NodeResourcesFitConfig{
						ScoringStrategy: &gfneks.Cluster_ScoringStrategy{
							Type: gfnt.NewString("LeastAllocated"),
							Resources: []gfneks.Cluster_ResourceWeight{
								{Name: gfnt.NewString("cpu"), Weight: gfnt.NewInteger(3)},
								{Name: gfnt.NewString("memory"), Weight: gfnt.NewInteger(1)},
							},
						},
					},
				},
			}),
		)
	})

	Describe("makeKubeControllerManagerConfig", func() {
		type entry struct {
			config   *api.KubeControllerManagerConfig
			expected *gfneks.Cluster_KubeControllerManagerConfig
		}

		DescribeTable("converts the config", func(e entry) {
			Expect(makeKubeControllerManagerConfig(e.config)).To(Equal(e.expected))
		},
			Entry("nil config", entry{
				config:   nil,
				expected: nil,
			}),
			Entry("empty config is omitted", entry{
				config:   &api.KubeControllerManagerConfig{},
				expected: nil,
			}),
			Entry("empty horizontal pod autoscaler config is omitted", entry{
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
				expected: &gfneks.Cluster_KubeControllerManagerConfig{
					HorizontalPodAutoscalerControllerConfig: &gfneks.Cluster_HorizontalPodAutoscalerControllerConfig{
						HorizontalPodAutoscalerSyncPeriod: gfnt.NewString("15s"),
					},
				},
			}),
		)
	})
})
