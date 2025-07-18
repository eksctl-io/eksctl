package v1alpha5

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	. "github.com/onsi/ginkgo/v2"

	. "github.com/onsi/gomega"
)

type nodeGroupCase struct {
	ng     *ManagedNodeGroup
	errMsg string
}

var _ = Describe("Managed Nodegroup Validation", func() {
	DescribeTable("Supported and unsupported field combinations", func(n *nodeGroupCase) {
		err := SetManagedNodeGroupDefaults(n.ng, &ClusterMeta{Name: "managed-cluster"}, false)
		Expect(err).NotTo(HaveOccurred())
		err = ValidateManagedNodeGroup(0, n.ng)
		if n.errMsg == "" {
			Expect(err).NotTo(HaveOccurred())
			return
		}
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring(n.errMsg))

	},
		Entry("Supported Windows as AMI family", &nodeGroupCase{
			ng: &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMIFamily: "WindowsServer2022FullContainer",
				},
			},
		}),
		Entry("Unsupported OverrideBootstrapCommand for Windows AMI", &nodeGroupCase{
			ng: &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMI:                      "",
					OverrideBootstrapCommand: aws.String(`bootstrap.sh`),
					AMIFamily:                "WindowsServer2019CoreContainer",
				},
			},
			errMsg: "overrideBootstrapCommand is not supported for WindowsServer2019CoreContainer nodegroups",
		}),
		Entry("Supported AMI family", &nodeGroupCase{
			ng: &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMIFamily: "AmazonLinux2",
				},
			},
		}),
		Entry("Custom AMI without AMI Family", &nodeGroupCase{
			ng: &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMI: "ami-custom",
				},
			},
			errMsg: "when using a custom AMI, amiFamily needs to be explicitly set via config file or via --node-ami-family flag",
		}),
		Entry("Custom AMI without overrideBootstrapCommand", &nodeGroupCase{
			ng: &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMI:       "ami-custom",
					AMIFamily: DefaultNodeImageFamily,
				},
			},
			errMsg: fmt.Sprintf("overrideBootstrapCommand is required when using a custom AMI based on %s", DefaultNodeImageFamily),
		}),
		Entry("Custom AMI with Windows AMI family without overrideBootstrapCommand", &nodeGroupCase{
			ng: &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMI:       "ami-custom",
					AMIFamily: "WindowsServer2019FullContainer",
				},
			},
			errMsg: "cannot set amiFamily to WindowsServer2019FullContainer when using a custom AMI",
		}),
		Entry("Custom AMI with Bottlerocket AMI family without overrideBootstrapCommand", &nodeGroupCase{
			ng: &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMI:       "ami-custom",
					AMIFamily: "Bottlerocket",
				},
			},
		}),
		Entry("Custom AMI with overrideBootstrapCommand", &nodeGroupCase{
			ng: &ManagedNodeGroup{
				NodeGroupBase: &NodeGroupBase{
					AMI:                      "ami-custom",
					AMIFamily:                DefaultNodeImageFamily,
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
		err := SetManagedNodeGroupDefaults(mng, &ClusterMeta{Name: "managed-cluster"}, false)
		Expect(err).NotTo(HaveOccurred())
		err = ValidateManagedNodeGroup(0, mng)
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
		err := SetManagedNodeGroupDefaults(mng, &ClusterMeta{Name: "managed-cluster"}, false)
		Expect(err).NotTo(HaveOccurred())
		err = ValidateManagedNodeGroup(0, mng)
		if e.valid {
			Expect(err).NotTo(HaveOccurred())
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
		Entry("returns an error if both maxUnavailable and maxUnavailablePercentage are not set", updateConfigEntry{
			valid: false,
		}),
	)
})
