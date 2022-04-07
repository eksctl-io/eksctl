package builder

import (
	"context"
	"math"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/pkg/errors"
	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	"github.com/weaveworks/eksctl/pkg/awsapi"
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
	ctx context.Context,
	launchTemplateData *gfnec2.LaunchTemplate_LaunchTemplateData,
	instanceTypes []string,
	efaEnabled bool,
	securityGroups []*gfnt.Value,
	ec2API awsapi.EC2,
) error {
	firstNI := defaultNetworkInterface(securityGroups, 0, 0)
	if efaEnabled {
		var instanceTypeList []ec2types.InstanceType
		for _, it := range instanceTypes {
			instanceTypeList = append(instanceTypeList, ec2types.InstanceType(it))
		}
		input := &ec2.DescribeInstanceTypesInput{
			InstanceTypes: instanceTypeList,
		}

		info, err := ec2API.DescribeInstanceTypes(ctx, input)
		if err != nil {
			return errors.Wrapf(err, "couldn't retrieve instance type description for %v", instanceTypes)
		}

		var numEFAs = math.MaxFloat64
		for _, it := range info.InstanceTypes {
			networkInfo := it.NetworkInfo
			numEFAs = math.Min(float64(aws.ToInt32(networkInfo.MaximumNetworkCards)), numEFAs)
			if !aws.ToBool(networkInfo.EfaSupported) {
				return errors.Errorf("instance type %s does not support EFA", it.InstanceType)
			}
		}

		firstNI.InterfaceType = gfnt.NewString("efa")
		nis := []gfnec2.LaunchTemplate_NetworkInterface{firstNI}
		// Only one card can be on deviceIndex=0
		// Additional cards are on deviceIndex=1
		// Due to ASG incompatibilities, we create each network card
		// with its own device
		for i := 1; i < int(numEFAs); i++ {
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
