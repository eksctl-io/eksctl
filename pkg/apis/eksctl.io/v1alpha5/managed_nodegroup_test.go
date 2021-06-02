package v1alpha5

import (
	"github.com/aws/aws-sdk-go/aws"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

type nodeGroupCase struct {
	ng     *ManagedNodeGroup
	errMsg string
}

var _ = Describe("Managed Nodegroup Validation", func() {
	DescribeTable("Supported and unsupported field combinations", func(n *nodeGroupCase) {
		SetManagedNodeGroupDefaults(n.ng, &ClusterMeta{Name: "managed-cluster"})
		err := ValidateManagedNodeGroup(n.ng, 0)
		if n.errMsg == "" {
			Expect(err).ToNot(HaveOccurred())
			return
		}
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(n.errMsg))

	},
		Entry("Unsupported AMI family", &nodeGroupCase{
			ng: &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMI:                      "",
					OverrideBootstrapCommand: nil,
					PreBootstrapCommands:     nil,
					AMIFamily:                "WindowsServer2019FullContainer",
				},
			},
			errMsg: `"WindowsServer2019FullContainer" is not supported for managed nodegroups`,
		}),
		Entry("Supported AMI family", &nodeGroupCase{
			ng: &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMIFamily: "AmazonLinux2",
				},
			},
		}),
		Entry("Custom AMI without overrideBootstrapCommand", &nodeGroupCase{
			ng: &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMI: "ami-custom",
				},
			},
			errMsg: "overrideBootstrapCommand is required when using a custom AMI",
		}),
		Entry("Custom AMI with overrideBootstrapCommand", &nodeGroupCase{
			ng: &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMI:                      "ami-custom",
					OverrideBootstrapCommand: aws.String(`bootstrap.sh`),
				},
			},
		}),
		Entry("launchTemplate with no ID", &nodeGroupCase{
			ng: &ManagedNodeGroup{
				NodeGroupBase:  &NodeGroupBase{},
				LaunchTemplate: &LaunchTemplate{},
			},
			errMsg: "launchTemplate.id is required",
		}),
		Entry("launchTemplate with ID", &nodeGroupCase{
			ng: &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{},
				LaunchTemplate: &LaunchTemplate{
					ID: "lt-1234",
				},
			},
		}),
		Entry("launchTemplate with invalid version", &nodeGroupCase{
			ng: &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{},
				LaunchTemplate: &LaunchTemplate{
					ID:      "lt-custom",
					Version: aws.String("0"),
				},
			},
			errMsg: "launchTemplate.version must be >= 1",
		}),
		Entry("launchTemplate with valid version", &nodeGroupCase{
			ng: &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{},
				LaunchTemplate: &LaunchTemplate{
					ID:      "lt-custom",
					Version: aws.String("3"),
				},
			},
		}),
		Entry("launchTemplate with instanceTypes", &nodeGroupCase{
			ng: &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{},
				InstanceTypes: []string{"c3.large", "c4.large"},
				LaunchTemplate: &LaunchTemplate{
					ID: "lt-custom",
				},
			},
		}),
		Entry("instanceSelector and instanceTypes", &nodeGroupCase{
			ng: &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{
					InstanceSelector: &InstanceSelector{
						VCPUs:  2,
						Memory: "4",
					},
				},
				InstanceTypes: []string{"c3.large", "c4.large"},
			},
		}),
		Entry("instanceSelector and instanceType", &nodeGroupCase{
			ng: &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{
					InstanceSelector: &InstanceSelector{
						VCPUs:  2,
						Memory: "4",
					},
					InstanceType: "c4.large",
				},
			},
			errMsg: "cannot set instanceType when instanceSelector is specified",
		}),
	)

	DescribeTable("User-supplied launch template with unsupported fields", func(ngBase *NodeGroupBase) {
		mng := &ManagedNodeGroup{
			NodeGroupBase: ngBase,
			LaunchTemplate: &LaunchTemplate{
				ID: "lt-custom",
			},
		}
		SetManagedNodeGroupDefaults(mng, &ClusterMeta{Name: "managed-cluster"})
		err := ValidateManagedNodeGroup(mng, 0)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("cannot set instanceType, ami, ssh.allow, ssh.enableSSM, ssh.sourceSecurityGroupIds, securityGroups, " +
			"volumeSize, instanceName, instancePrefix, maxPodsPerNode, disableIMDSv1, disablePodIMDS, preBootstrapCommands, overrideBootstrapCommand, placement in managedNodeGroup when a launch template is supplied"))
	},
		Entry("instanceType", &NodeGroupBase{
			InstanceType: "m5.xlarge",
		}),
		Entry("AMI", &NodeGroupBase{
			AMI: "ami-custom",
		}),
		Entry("SSH", &NodeGroupBase{
			SSH: &NodeGroupSSH{
				Allow: Enabled(),
			},
		}),
		Entry("enableSSM", &NodeGroupBase{
			SSH: &NodeGroupSSH{
				EnableSSM: Enabled(),
			},
		}),
		Entry("volumeSize", &NodeGroupBase{
			VolumeSize: aws.Int(100),
		}),
		Entry("preBootstrapCommands", &NodeGroupBase{
			PreBootstrapCommands: []string{"echo test"},
		}),
		Entry("overrideBootstrapCommand", &NodeGroupBase{
			OverrideBootstrapCommand: aws.String("bootstrap.sh"),
		}),
		Entry("securityGroups.attachIDs", &NodeGroupBase{
			SecurityGroups: &NodeGroupSGs{
				AttachIDs: []string{"sg-custom"},
			},
		}),
	)

	type updateConfigEntry struct {
		unavailable           *int
		unavailablePercentage *int
		maxSize               *int
		valid                 bool
	}

	DescribeTable("UpdateConfig", func(e updateConfigEntry) {
		mng := &ManagedNodeGroup{
			NodeGroupBase: &NodeGroupBase{
				AMIFamily: "AmazonLinux2",
				ScalingConfig: &ScalingConfig{
					MaxSize: e.maxSize,
				},
			},
			UpdateConfig: &NodeGroupUpdateConfig{
				MaxUnavailable:           e.unavailable,
				MaxUnavailablePercentage: e.unavailablePercentage,
			},
		}
		SetManagedNodeGroupDefaults(mng, &ClusterMeta{Name: "managed-cluster"})
		err := ValidateManagedNodeGroup(mng, 0)
		if e.valid {
			Expect(err).ToNot(HaveOccurred())
		} else {
			Expect(err).To(HaveOccurred())
		}
	},
		Entry("max unavailable set", updateConfigEntry{
			unavailable: aws.Int(1),
			valid:       true,
		}),
		Entry("max unavailable specified in percentage", updateConfigEntry{
			unavailablePercentage: aws.Int(1),
			valid:                 true,
		}),
		Entry("returns an error if both are set", updateConfigEntry{
			unavailable:           aws.Int(1),
			unavailablePercentage: aws.Int(1),
			valid:                 false,
		}),
		Entry("returns an error if max unavailable is greater than maxSize", updateConfigEntry{
			unavailable: aws.Int(100),
			maxSize:     aws.Int(5),
			valid:       false,
		}),
		Entry("returns an error if max unavailable percentage is greater than maxSize", updateConfigEntry{
			unavailablePercentage: aws.Int(100),
			maxSize:               aws.Int(5),
			valid:                 false,
		}),
	)
})
