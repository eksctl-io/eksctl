package builder

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/pkg/errors"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"
)

func defaultNetworkInterface(securityGroups []*gfnt.Value, device, card int) gfnec2.LaunchTemplate_NetworkInterface {
	return gfnec2.LaunchTemplate_NetworkInterface{
		// Explicitly un-setting this so that it doesn't get defaulted to true
		AssociatePublicIpAddress: nil,
		DeviceIndex:              gfnt.NewInteger(device),
		Groups:                   gfnt.NewSlice(securityGroups...),
		NetworkCardIndex:         gfnt.NewInteger(card),
	}
}

func buildNetworkInterfaces(
	launchTemplateData *gfnec2.LaunchTemplate_LaunchTemplateData,
	spec *api.NodeGroupBase,
	securityGroups []*gfnt.Value,
	ec2api ec2iface.EC2API,
) error {
	firstNI := defaultNetworkInterface(securityGroups, 0, 0)
	if api.IsEnabled(spec.EFAEnabled) {
		input := ec2.DescribeInstanceTypesInput{
			InstanceTypes: aws.StringSlice([]string{spec.InstanceType}),
		}
		info, err := ec2api.DescribeInstanceTypes(&input)
		if err != nil {
			return errors.Wrapf(err, "couldn't retrieve instance type description for %s", spec.InstanceType)
		}
		numEFAs := int(aws.Int64Value(info.InstanceTypes[0].NetworkInfo.MaximumNetworkCards))
		firstNI.InterfaceType = gfnt.NewString("efa")
		nis := []gfnec2.LaunchTemplate_NetworkInterface{firstNI}
		// Only one card can be on deviceIndex=0
		// Additional cards are on deviceIndex=1
		// Due to ASG incompatibilities, we create each network card
		// with its own device
		for i := 1; i < numEFAs; i++ {
			ni := defaultNetworkInterface(securityGroups, i, i)
			ni.InterfaceType = gfnt.NewString("efa")
			nis = append(nis, ni)
		}
		launchTemplateData.NetworkInterfaces = nis
	} else {
		launchTemplateData.NetworkInterfaces = []gfnec2.LaunchTemplate_NetworkInterface{firstNI}
	}
	return nil
}
