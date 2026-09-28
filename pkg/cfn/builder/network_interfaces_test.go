package builder

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	gfnec2 "github.com/weaveworks/eksctl/pkg/goformation/cloudformation/ec2"
	gfnt "github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

func TestBuildNetworkInterfaces(t *testing.T) {
	tests := []struct {
		name                      string
		instanceTypes             []string
		efaEnabled                bool
		securityGroups            []*gfnt.Value
		connectionTracking        *api.ConnectionTracking
		mockInstanceTypes         []ec2types.InstanceTypeInfo
		expectedNetworkInterfaces int
		expectedInterfaceType     string
		expectedError             string
	}{
		{
			name:          "non-EFA nodegroup",
			instanceTypes: []string{"m5.large"},
			efaEnabled:    false,
			securityGroups: []*gfnt.Value{
				gfnt.NewString("sg-12345"),
				gfnt.NewString("sg-67890"),
			},
			expectedNetworkInterfaces: 1,
			expectedInterfaceType:     "",
		},
		{
			name:          "EFA nodegroup with single network card",
			instanceTypes: []string{"c5n.large"},
			efaEnabled:    true,
			securityGroups: []*gfnt.Value{
				gfnt.NewString("sg-12345"),
				gfnt.NewString("sg-67890"),
			},
			mockInstanceTypes: []ec2types.InstanceTypeInfo{
				{
					InstanceType: ec2types.InstanceTypeC5nLarge,
					NetworkInfo: &ec2types.NetworkInfo{
						MaximumNetworkCards: aws.Int32(1),
						EfaSupported:        aws.Bool(true),
					},
				},
			},
			expectedNetworkInterfaces: 1,
			expectedInterfaceType:     "efa",
		},
		{
			name:          "EFA nodegroup with multiple network cards",
			instanceTypes: []string{"c5n.18xlarge"},
			efaEnabled:    true,
			securityGroups: []*gfnt.Value{
				gfnt.NewString("sg-12345"),
				gfnt.NewString("sg-67890"),
			},
			mockInstanceTypes: []ec2types.InstanceTypeInfo{
				{
					InstanceType: ec2types.InstanceTypeC5n18xlarge,
					NetworkInfo: &ec2types.NetworkInfo{
						MaximumNetworkCards: aws.Int32(4),
						EfaSupported:        aws.Bool(true),
					},
				},
			},
			expectedNetworkInterfaces: 4,
			expectedInterfaceType:     "efa",
		},
		{
			name:          "EFA nodegroup with mixed instance types",
			instanceTypes: []string{"c5n.large", "c5n.xlarge"},
			efaEnabled:    true,
			securityGroups: []*gfnt.Value{
				gfnt.NewString("sg-12345"),
				gfnt.NewString("sg-67890"),
			},
			mockInstanceTypes: []ec2types.InstanceTypeInfo{
				{
					InstanceType: ec2types.InstanceTypeC5nLarge,
					NetworkInfo: &ec2types.NetworkInfo{
						MaximumNetworkCards: aws.Int32(1),
						EfaSupported:        aws.Bool(true),
					},
				},
				{
					InstanceType: ec2types.InstanceTypeC5nXlarge,
					NetworkInfo: &ec2types.NetworkInfo{
						MaximumNetworkCards: aws.Int32(2),
						EfaSupported:        aws.Bool(true),
					},
				},
			},
			expectedNetworkInterfaces: 1, // Should use minimum across instance types
			expectedInterfaceType:     "efa",
		},
		{
			name:          "EFA nodegroup with non-EFA instance type",
			instanceTypes: []string{"m5.large"},
			efaEnabled:    true,
			securityGroups: []*gfnt.Value{
				gfnt.NewString("sg-12345"),
			},
			mockInstanceTypes: []ec2types.InstanceTypeInfo{
				{
					InstanceType: ec2types.InstanceTypeM5Large,
					NetworkInfo: &ec2types.NetworkInfo{
						MaximumNetworkCards: aws.Int32(2),
						EfaSupported:        aws.Bool(false),
					},
				},
			},
			expectedError: "instance type m5.large does not support EFA",
		},
		{
			name:          "EFA nodegroup with default security groups (1.33+ scenario)",
			instanceTypes: []string{"c5n.18xlarge"},
			efaEnabled:    true,
			securityGroups: []*gfnt.Value{
				gfnt.NewString("sg-default-shared"),
				gfnt.NewString("sg-nodegroup-local"),
			},
			mockInstanceTypes: []ec2types.InstanceTypeInfo{
				{
					InstanceType: ec2types.InstanceTypeC5n18xlarge,
					NetworkInfo: &ec2types.NetworkInfo{
						MaximumNetworkCards: aws.Int32(4),
						EfaSupported:        aws.Bool(true),
					},
				},
			},
			expectedNetworkInterfaces: 4,
			expectedInterfaceType:     "efa",
		},
		{
			name:          "EFA nodegroup with custom EFA security groups (1.32 scenario)",
			instanceTypes: []string{"c5n.18xlarge"},
			efaEnabled:    true,
			securityGroups: []*gfnt.Value{
				gfnt.NewString("sg-default-shared"),
				gfnt.NewString("sg-nodegroup-local"),
				gfnt.NewString("sg-custom-efa"),
			},
			mockInstanceTypes: []ec2types.InstanceTypeInfo{
				{
					InstanceType: ec2types.InstanceTypeC5n18xlarge,
					NetworkInfo: &ec2types.NetworkInfo{
						MaximumNetworkCards: aws.Int32(4),
						EfaSupported:        aws.Bool(true),
					},
				},
			},
			expectedNetworkInterfaces: 4,
			expectedInterfaceType:     "efa",
		},
		{
			name:          "non-EFA nodegroup with connection tracking",
			instanceTypes: []string{"c8in.8xlarge"},
			efaEnabled:    false,
			securityGroups: []*gfnt.Value{
				gfnt.NewString("sg-12345"),
			},
			connectionTracking: &api.ConnectionTracking{
				TCPEstablishedTimeout: aws.Int(432000),
				UDPStreamTimeout:      aws.Int(180),
				UDPTimeout:            aws.Int(60),
			},
			mockInstanceTypes: []ec2types.InstanceTypeInfo{
				{
					InstanceType: ec2types.InstanceTypeC8in8xlarge,
					Hypervisor:   ec2types.InstanceTypeHypervisorNitro,
				},
			},
			expectedNetworkInterfaces: 1,
			expectedInterfaceType:     "",
		},
		{
			name:          "non-EFA nodegroup with a single connection tracking timeout",
			instanceTypes: []string{"c8in.8xlarge"},
			efaEnabled:    false,
			securityGroups: []*gfnt.Value{
				gfnt.NewString("sg-12345"),
			},
			connectionTracking: &api.ConnectionTracking{
				TCPEstablishedTimeout: aws.Int(3600),
			},
			mockInstanceTypes: []ec2types.InstanceTypeInfo{
				{
					InstanceType: ec2types.InstanceTypeC8in8xlarge,
					Hypervisor:   ec2types.InstanceTypeHypervisorNitro,
				},
			},
			expectedNetworkInterfaces: 1,
			expectedInterfaceType:     "",
		},
		{
			name:          "EFA nodegroup with connection tracking",
			instanceTypes: []string{"c5n.18xlarge"},
			efaEnabled:    true,
			securityGroups: []*gfnt.Value{
				gfnt.NewString("sg-12345"),
			},
			connectionTracking: &api.ConnectionTracking{
				TCPEstablishedTimeout: aws.Int(432000),
				UDPStreamTimeout:      aws.Int(180),
				UDPTimeout:            aws.Int(60),
			},
			mockInstanceTypes: []ec2types.InstanceTypeInfo{
				{
					InstanceType: ec2types.InstanceTypeC5n18xlarge,
					Hypervisor:   ec2types.InstanceTypeHypervisorNitro,
					NetworkInfo: &ec2types.NetworkInfo{
						MaximumNetworkCards: aws.Int32(4),
						EfaSupported:        aws.Bool(true),
					},
				},
			},
			expectedNetworkInterfaces: 4,
			expectedInterfaceType:     "efa",
		},
		{
			// Bare metal instance types report no hypervisor at all, but they are delivered on
			// the Nitro system, so they must not be rejected.
			name:          "nodegroup with connection tracking on a bare metal instance type",
			instanceTypes: []string{"c5n.metal"},
			efaEnabled:    false,
			securityGroups: []*gfnt.Value{
				gfnt.NewString("sg-12345"),
			},
			connectionTracking: &api.ConnectionTracking{
				TCPEstablishedTimeout: aws.Int(432000),
			},
			mockInstanceTypes: []ec2types.InstanceTypeInfo{
				{
					InstanceType: ec2types.InstanceTypeC5nMetal,
					BareMetal:    aws.Bool(true),
				},
			},
			expectedNetworkInterfaces: 1,
			expectedInterfaceType:     "",
		},
		{
			name:          "nodegroup with connection tracking on a non-Nitro instance type",
			instanceTypes: []string{"m4.large"},
			efaEnabled:    false,
			securityGroups: []*gfnt.Value{
				gfnt.NewString("sg-12345"),
			},
			connectionTracking: &api.ConnectionTracking{
				TCPEstablishedTimeout: aws.Int(432000),
			},
			mockInstanceTypes: []ec2types.InstanceTypeInfo{
				{
					InstanceType: ec2types.InstanceTypeM4Large,
					Hypervisor:   ec2types.InstanceTypeHypervisorXen,
					BareMetal:    aws.Bool(false),
				},
			},
			expectedError: "connectionTracking is only supported on Nitro-based instance types; m4.large cannot be used with it",
		},
		{
			name:          "nodegroup with connection tracking on a mixed instance list containing non-Nitro types",
			instanceTypes: []string{"m5.large", "m4.large", "t2.medium"},
			efaEnabled:    false,
			securityGroups: []*gfnt.Value{
				gfnt.NewString("sg-12345"),
			},
			connectionTracking: &api.ConnectionTracking{
				TCPEstablishedTimeout: aws.Int(432000),
			},
			mockInstanceTypes: []ec2types.InstanceTypeInfo{
				{
					InstanceType: ec2types.InstanceTypeM5Large,
					Hypervisor:   ec2types.InstanceTypeHypervisorNitro,
				},
				{
					InstanceType: ec2types.InstanceTypeM4Large,
					Hypervisor:   ec2types.InstanceTypeHypervisorXen,
				},
				{
					InstanceType: ec2types.InstanceTypeT2Medium,
					Hypervisor:   ec2types.InstanceTypeHypervisorXen,
				},
			},
			expectedError: "connectionTracking is only supported on Nitro-based instance types; m4.large, t2.medium cannot be used with it",
		},
		{
			// EFA enabled and connection tracking set share a single DescribeInstanceTypes call,
			// and the support check runs before the EFA interface count.
			name:          "EFA nodegroup with connection tracking on a non-Nitro instance type",
			instanceTypes: []string{"m4.large"},
			efaEnabled:    true,
			securityGroups: []*gfnt.Value{
				gfnt.NewString("sg-12345"),
			},
			connectionTracking: &api.ConnectionTracking{
				TCPEstablishedTimeout: aws.Int(432000),
			},
			mockInstanceTypes: []ec2types.InstanceTypeInfo{
				{
					InstanceType: ec2types.InstanceTypeM4Large,
					Hypervisor:   ec2types.InstanceTypeHypervisorXen,
					NetworkInfo: &ec2types.NetworkInfo{
						MaximumNetworkCards: aws.Int32(1),
						EfaSupported:        aws.Bool(false),
					},
				},
			},
			expectedError: "connectionTracking is only supported on Nitro-based instance types; m4.large cannot be used with it",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProvider := mockprovider.NewMockProvider()

			// The instance type description is fetched once for EFA, for connection tracking,
			// or for both together.
			if tt.efaEnabled || tt.connectionTracking != nil {
				mockProvider.MockEC2().On("DescribeInstanceTypes",
					context.Background(),
					&ec2.DescribeInstanceTypesInput{
						InstanceTypes: func() []ec2types.InstanceType {
							var types []ec2types.InstanceType
							for _, it := range tt.instanceTypes {
								types = append(types, ec2types.InstanceType(it))
							}
							return types
						}(),
					}).Return(&ec2.DescribeInstanceTypesOutput{
					InstanceTypes: tt.mockInstanceTypes,
				}, nil).Once()
			}

			launchTemplateData := &gfnec2.LaunchTemplate_LaunchTemplateData{}

			err := buildNetworkInterfaces(
				context.Background(),
				launchTemplateData,
				tt.instanceTypes,
				tt.efaEnabled,
				tt.securityGroups,
				tt.connectionTracking,
				mockProvider.EC2(),
			)

			mockProvider.MockEC2().AssertExpectations(t)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, launchTemplateData.NetworkInterfaces)
			assert.Len(t, launchTemplateData.NetworkInterfaces, tt.expectedNetworkInterfaces)

			// Verify first network interface
			firstNI := launchTemplateData.NetworkInterfaces[0]
			assert.Equal(t, gfnt.Integer(0), firstNI.DeviceIndex.Raw())
			assert.Equal(t, gfnt.Integer(0), firstNI.NetworkCardIndex.Raw())
			assert.Nil(t, firstNI.AssociatePublicIpAddress)

			// Verify security groups are correctly assigned
			require.NotNil(t, firstNI.Groups)
			groupsSlice := firstNI.Groups.Raw().(gfnt.Slice)
			assert.Len(t, groupsSlice, len(tt.securityGroups))
			for i, sg := range tt.securityGroups {
				assert.Equal(t, sg.Raw(), groupsSlice[i].Raw())
			}

			// Verify connection tracking is applied to every network interface
			for _, ni := range launchTemplateData.NetworkInterfaces {
				assertConnectionTracking(t, tt.connectionTracking, ni.ConnectionTrackingSpecification)
			}

			if tt.efaEnabled && tt.expectedError == "" {
				// Verify EFA interface type
				require.NotNil(t, firstNI.InterfaceType)
				assert.Equal(t, gfnt.String(tt.expectedInterfaceType), firstNI.InterfaceType.Raw())

				// Verify additional network interfaces for multi-card instances
				for i := 1; i < tt.expectedNetworkInterfaces; i++ {
					ni := launchTemplateData.NetworkInterfaces[i]
					assert.Equal(t, gfnt.Integer(1), ni.DeviceIndex.Raw())
					assert.Equal(t, gfnt.Integer(i), ni.NetworkCardIndex.Raw())
					assert.Nil(t, ni.AssociatePublicIpAddress)
					require.NotNil(t, ni.InterfaceType)
					assert.Equal(t, gfnt.String(tt.expectedInterfaceType), ni.InterfaceType.Raw())

					// Verify security groups are correctly assigned to additional interfaces
					require.NotNil(t, ni.Groups)
					groupsSlice := ni.Groups.Raw().(gfnt.Slice)
					assert.Len(t, groupsSlice, len(tt.securityGroups))
					for j, sg := range tt.securityGroups {
						assert.Equal(t, sg.Raw(), groupsSlice[j].Raw())
					}
				}
			} else if !tt.efaEnabled {
				// Verify non-EFA interface
				assert.Nil(t, firstNI.InterfaceType)
			}
		})
	}
}

func TestBuildNetworkInterfaces_EC2APIError(t *testing.T) {
	mockProvider := mockprovider.NewMockProvider()

	instanceTypes := []string{"c5n.large"}
	mockProvider.MockEC2().On("DescribeInstanceTypes",
		context.Background(),
		&ec2.DescribeInstanceTypesInput{
			InstanceTypes: []ec2types.InstanceType{ec2types.InstanceTypeC5nLarge},
		}).Return(nil, assert.AnError)

	launchTemplateData := &gfnec2.LaunchTemplate_LaunchTemplateData{}
	securityGroups := []*gfnt.Value{gfnt.NewString("sg-12345")}

	err := buildNetworkInterfaces(
		context.Background(),
		launchTemplateData,
		instanceTypes,
		true, // EFA enabled
		securityGroups,
		nil,
		mockProvider.EC2(),
	)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "couldn't retrieve instance type description")
}

func TestDefaultNetworkInterface(t *testing.T) {
	securityGroups := []*gfnt.Value{
		gfnt.NewString("sg-12345"),
		gfnt.NewString("sg-67890"),
	}

	ni := defaultNetworkInterface(securityGroups, 1, 2, nil)

	assert.Nil(t, ni.AssociatePublicIpAddress)
	assert.Equal(t, gfnt.Integer(1), ni.DeviceIndex.Raw())
	assert.Equal(t, gfnt.Integer(2), ni.NetworkCardIndex.Raw())
	assert.Nil(t, ni.ConnectionTrackingSpecification)
	require.NotNil(t, ni.Groups)
	groupsSlice := ni.Groups.Raw().(gfnt.Slice)
	assert.Len(t, groupsSlice, 2)
	assert.Equal(t, gfnt.String("sg-12345"), groupsSlice[0].Raw())
	assert.Equal(t, gfnt.String("sg-67890"), groupsSlice[1].Raw())
}

func TestMakeConnectionTrackingSpecification(t *testing.T) {
	t.Run("nil connection tracking", func(t *testing.T) {
		assert.Nil(t, makeConnectionTrackingSpecification(nil))
	})

	t.Run("all timeouts set", func(t *testing.T) {
		spec := makeConnectionTrackingSpecification(&api.ConnectionTracking{
			TCPEstablishedTimeout: aws.Int(432000),
			UDPStreamTimeout:      aws.Int(180),
			UDPTimeout:            aws.Int(60),
		})
		require.NotNil(t, spec)
		assert.Equal(t, gfnt.Integer(432000), spec.TcpEstablishedTimeout.Raw())
		assert.Equal(t, gfnt.Integer(180), spec.UdpStreamTimeout.Raw())
		assert.Equal(t, gfnt.Integer(60), spec.UdpTimeout.Raw())
	})

	t.Run("unset timeouts are omitted", func(t *testing.T) {
		spec := makeConnectionTrackingSpecification(&api.ConnectionTracking{
			UDPTimeout: aws.Int(30),
		})
		require.NotNil(t, spec)
		assert.Nil(t, spec.TcpEstablishedTimeout)
		assert.Nil(t, spec.UdpStreamTimeout)
		assert.Equal(t, gfnt.Integer(30), spec.UdpTimeout.Raw())
	})
}

func assertConnectionTracking(t *testing.T, expected *api.ConnectionTracking, actual *gfnec2.LaunchTemplate_ConnectionTrackingSpecification) {
	t.Helper()
	if expected == nil {
		assert.Nil(t, actual)
		return
	}
	require.NotNil(t, actual)
	assertTimeout := func(expected *int, actual *gfnt.Value) {
		if expected == nil {
			assert.Nil(t, actual)
			return
		}
		require.NotNil(t, actual)
		assert.Equal(t, gfnt.Integer(*expected), actual.Raw())
	}
	assertTimeout(expected.TCPEstablishedTimeout, actual.TcpEstablishedTimeout)
	assertTimeout(expected.UDPStreamTimeout, actual.UdpStreamTimeout)
	assertTimeout(expected.UDPTimeout, actual.UdpTimeout)
}
