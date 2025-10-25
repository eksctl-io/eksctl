package v1alpha5

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	. "github.com/onsi/ginkgo/v2"
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

	Describe("NodeGroupNodeRepairConfig", func() {
		var (
			nodeRepairConfig *NodeGroupNodeRepairConfig
		)

		BeforeEach(func() {
			nodeRepairConfig = &NodeGroupNodeRepairConfig{}
		})

		Describe("JSON marshaling and unmarshaling", func() {
			When("all fields are set", func() {
				It("should marshal and unmarshal correctly", func() {
					nodeRepairConfig.Enabled = aws.Bool(true)
					nodeRepairConfig.MaxUnhealthyNodeThresholdPercentage = aws.Int(20)
					nodeRepairConfig.MaxUnhealthyNodeThresholdCount = aws.Int(5)
					nodeRepairConfig.MaxParallelNodesRepairedPercentage = aws.Int(15)
					nodeRepairConfig.MaxParallelNodesRepairedCount = aws.Int(2)
					nodeRepairConfig.NodeRepairConfigOverrides = []NodeRepairConfigOverride{
						{
							NodeMonitoringCondition: "AcceleratedInstanceNotReady",
							NodeUnhealthyReason:     "NvidiaXID13Error",
							MinRepairWaitTimeMins:   10,
							RepairAction:            "Terminate",
						},
						{
							NodeMonitoringCondition: "NetworkNotReady",
							NodeUnhealthyReason:     "InterfaceNotUp",
							MinRepairWaitTimeMins:   20,
							RepairAction:            "Restart",
						},
					}

					// Test JSON marshaling
					jsonData, err := json.Marshal(nodeRepairConfig)
					Expect(err).NotTo(HaveOccurred())
					Expect(string(jsonData)).To(ContainSubstring(`"enabled":true`))
					Expect(string(jsonData)).To(ContainSubstring(`"maxUnhealthyNodeThresholdPercentage":20`))
					Expect(string(jsonData)).To(ContainSubstring(`"maxUnhealthyNodeThresholdCount":5`))
					Expect(string(jsonData)).To(ContainSubstring(`"maxParallelNodesRepairedPercentage":15`))
					Expect(string(jsonData)).To(ContainSubstring(`"maxParallelNodesRepairedCount":2`))
					Expect(string(jsonData)).To(ContainSubstring(`"nodeRepairConfigOverrides"`))
					Expect(string(jsonData)).To(ContainSubstring(`"AcceleratedInstanceNotReady"`))
					Expect(string(jsonData)).To(ContainSubstring(`"NvidiaXID13Error"`))

					// Test JSON unmarshaling
					var unmarshaled NodeGroupNodeRepairConfig
					err = json.Unmarshal(jsonData, &unmarshaled)
					Expect(err).NotTo(HaveOccurred())
					Expect(*unmarshaled.Enabled).To(BeTrue())
					Expect(*unmarshaled.MaxUnhealthyNodeThresholdPercentage).To(Equal(20))
					Expect(*unmarshaled.MaxUnhealthyNodeThresholdCount).To(Equal(5))
					Expect(*unmarshaled.MaxParallelNodesRepairedPercentage).To(Equal(15))
					Expect(*unmarshaled.MaxParallelNodesRepairedCount).To(Equal(2))
					Expect(len(unmarshaled.NodeRepairConfigOverrides)).To(Equal(2))
					Expect(unmarshaled.NodeRepairConfigOverrides[0].NodeMonitoringCondition).To(Equal("AcceleratedInstanceNotReady"))
					Expect(unmarshaled.NodeRepairConfigOverrides[0].NodeUnhealthyReason).To(Equal("NvidiaXID13Error"))
					Expect(unmarshaled.NodeRepairConfigOverrides[0].MinRepairWaitTimeMins).To(Equal(10))
					Expect(unmarshaled.NodeRepairConfigOverrides[0].RepairAction).To(Equal("Terminate"))
				})
			})

			When("only enabled field is set", func() {
				It("should marshal and unmarshal correctly with minimal config", func() {
					nodeRepairConfig.Enabled = aws.Bool(true)

					// Test JSON marshaling
					jsonData, err := json.Marshal(nodeRepairConfig)
					Expect(err).NotTo(HaveOccurred())
					Expect(string(jsonData)).To(ContainSubstring(`"enabled":true`))
					Expect(string(jsonData)).NotTo(ContainSubstring(`"maxUnhealthyNodeThresholdPercentage"`))

					// Test JSON unmarshaling
					var unmarshaled NodeGroupNodeRepairConfig
					err = json.Unmarshal(jsonData, &unmarshaled)
					Expect(err).NotTo(HaveOccurred())
					Expect(*unmarshaled.Enabled).To(BeTrue())
					Expect(unmarshaled.MaxUnhealthyNodeThresholdPercentage).To(BeNil())
					Expect(unmarshaled.MaxUnhealthyNodeThresholdCount).To(BeNil())
					Expect(unmarshaled.MaxParallelNodesRepairedPercentage).To(BeNil())
					Expect(unmarshaled.MaxParallelNodesRepairedCount).To(BeNil())
					Expect(len(unmarshaled.NodeRepairConfigOverrides)).To(Equal(0))
				})
			})

			When("enabled is false", func() {
				It("should handle disabled state correctly", func() {
					nodeRepairConfig.Enabled = aws.Bool(false)

					jsonData, err := json.Marshal(nodeRepairConfig)
					Expect(err).NotTo(HaveOccurred())
					Expect(string(jsonData)).To(ContainSubstring(`"enabled":false`))

					var unmarshaled NodeGroupNodeRepairConfig
					err = json.Unmarshal(jsonData, &unmarshaled)
					Expect(err).NotTo(HaveOccurred())
					Expect(*unmarshaled.Enabled).To(BeFalse())
				})
			})
		})

		Describe("NodeRepairConfigOverride", func() {
			var override NodeRepairConfigOverride

			BeforeEach(func() {
				override = NodeRepairConfigOverride{
					NodeMonitoringCondition: "NetworkNotReady",
					NodeUnhealthyReason:     "InterfaceNotUp",
					MinRepairWaitTimeMins:   15,
					RepairAction:            "Restart",
				}
			})

			It("should have all required fields", func() {
				Expect(override.NodeMonitoringCondition).To(Equal("NetworkNotReady"))
				Expect(override.NodeUnhealthyReason).To(Equal("InterfaceNotUp"))
				Expect(override.MinRepairWaitTimeMins).To(Equal(15))
				Expect(override.RepairAction).To(Equal("Restart"))
			})

			It("should marshal to JSON correctly", func() {
				jsonData, err := json.Marshal(override)
				Expect(err).NotTo(HaveOccurred())
				Expect(string(jsonData)).To(ContainSubstring(`"nodeMonitoringCondition":"NetworkNotReady"`))
				Expect(string(jsonData)).To(ContainSubstring(`"nodeUnhealthyReason":"InterfaceNotUp"`))
				Expect(string(jsonData)).To(ContainSubstring(`"minRepairWaitTimeMins":15`))
				Expect(string(jsonData)).To(ContainSubstring(`"repairAction":"Restart"`))
			})

			It("should unmarshal from JSON correctly", func() {
				jsonStr := `{
					"nodeMonitoringCondition": "AcceleratedInstanceNotReady",
					"nodeUnhealthyReason": "NvidiaXID13Error",
					"minRepairWaitTimeMins": 25,
					"repairAction": "Terminate"
				}`

				var unmarshaled NodeRepairConfigOverride
				err := json.Unmarshal([]byte(jsonStr), &unmarshaled)
				Expect(err).NotTo(HaveOccurred())
				Expect(unmarshaled.NodeMonitoringCondition).To(Equal("AcceleratedInstanceNotReady"))
				Expect(unmarshaled.NodeUnhealthyReason).To(Equal("NvidiaXID13Error"))
				Expect(unmarshaled.MinRepairWaitTimeMins).To(Equal(25))
				Expect(unmarshaled.RepairAction).To(Equal("Terminate"))
			})
		})

		Describe("Pointer field handling", func() {
			It("should distinguish between nil and zero values", func() {
				// Test nil values
				config1 := &NodeGroupNodeRepairConfig{}
				Expect(config1.Enabled).To(BeNil())
				Expect(config1.MaxUnhealthyNodeThresholdPercentage).To(BeNil())

				// Test zero values
				config2 := &NodeGroupNodeRepairConfig{
					Enabled:                                 aws.Bool(false),
					MaxUnhealthyNodeThresholdPercentage:     aws.Int(0),
					MaxUnhealthyNodeThresholdCount:          aws.Int(0),
					MaxParallelNodesRepairedPercentage:      aws.Int(0),
					MaxParallelNodesRepairedCount:           aws.Int(0),
				}
				Expect(config2.Enabled).NotTo(BeNil())
				Expect(*config2.Enabled).To(BeFalse())
				Expect(config2.MaxUnhealthyNodeThresholdPercentage).NotTo(BeNil())
				Expect(*config2.MaxUnhealthyNodeThresholdPercentage).To(Equal(0))
			})
		})
	})

})
