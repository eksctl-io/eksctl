package v1alpha5

import (
	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Types", func() {
	var (
		cfg *ClusterConfig
	)
	Describe("HasNodes", func() {
		BeforeEach(func() {
			cfg = NewClusterConfig()
		})

		When("no nodegroups or managed nodegroups exist", func() {
			It("returns false", func() {
				Expect(cfg.HasNodes()).To(BeFalse())
			})
		})

		When("a nodegroup exists but is scaled to 0", func() {
			It("returns false", func() {
				cfg.NodeGroups = []*NodeGroup{
					{
						NodeGroupBase: &NodeGroupBase{
							ScalingConfig: &ScalingConfig{
								DesiredCapacity: aws.Int(0),
							},
						},
					},
				}
				Expect(cfg.HasNodes()).To(BeFalse())
			})
		})

		When("a managed nodegroup exists but is scaled to 0", func() {
			It("returns false", func() {
				cfg.ManagedNodeGroups = []*ManagedNodeGroup{
					{
						NodeGroupBase: &NodeGroupBase{
							ScalingConfig: &ScalingConfig{
								DesiredCapacity: aws.Int(0),
							},
						},
					},
				}
				Expect(cfg.HasNodes()).To(BeFalse())
			})
		})

		When("a nodegroup exists and is scaled above 0", func() {
			It("returns true", func() {
				cfg.NodeGroups = []*NodeGroup{
					{
						NodeGroupBase: &NodeGroupBase{
							ScalingConfig: &ScalingConfig{
								DesiredCapacity: aws.Int(1),
							},
						},
					},
				}
				Expect(cfg.HasNodes()).To(BeTrue())
			})
		})

		When("a managed nodegroup exists and is scaled above 0", func() {
			It("returns true", func() {
				cfg.ManagedNodeGroups = []*ManagedNodeGroup{
					{
						NodeGroupBase: &NodeGroupBase{
							ScalingConfig: &ScalingConfig{
								DesiredCapacity: aws.Int(1),
							},
						},
					},
				}
				Expect(cfg.HasNodes()).To(BeTrue())
			})
		})
	})

})
