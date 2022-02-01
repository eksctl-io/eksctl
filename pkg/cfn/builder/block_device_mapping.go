package builder

import (
	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

const rootDevice = "/dev/xvda"

func makeBlockDeviceMappings(ng *api.NodeGroupBase) []gfnec2.LaunchTemplate_BlockDeviceMapping {
	var mappings []gfnec2.LaunchTemplate_BlockDeviceMapping
	mapping := makeBlockDeviceMapping(api.VolumeMapping{
		VolumeSize:       ng.VolumeSize,
		VolumeType:       ng.VolumeType,
		VolumeName:       ng.VolumeName,
		VolumeEncrypted:  ng.VolumeEncrypted,
		VolumeKmsKeyID:   ng.VolumeKmsKeyID,
		VolumeIOPS:       ng.VolumeIOPS,
		VolumeThroughput: ng.VolumeThroughput,
	})

	if mapping != nil {
		mappings = append(mappings, *mapping)
	}

	for _, volume := range ng.AdditionalVolumes {
		if dm := makeBlockDeviceMapping(volume); dm != nil {
			mappings = append(mappings, *dm)
		}
	}

	if api.IsEnabled(ng.VolumeEncrypted) && ng.AdditionalEncryptedVolume != "" {
		mappings = append(mappings, gfnec2.LaunchTemplate_BlockDeviceMapping{
			DeviceName: gfnt.NewString(ng.AdditionalEncryptedVolume),
			Ebs: &gfnec2.LaunchTemplate_Ebs{
				Encrypted: gfnt.NewBoolean(*ng.VolumeEncrypted),
				KmsKeyId:  mapping.Ebs.KmsKeyId,
			},
		})
	}

	return mappings
}

func makeBlockDeviceMapping(vm api.VolumeMapping) *gfnec2.LaunchTemplate_BlockDeviceMapping {
	volumeSize := vm.VolumeSize
	if volumeSize == nil || *volumeSize == 0 {
		return nil
	}

	mapping := gfnec2.LaunchTemplate_BlockDeviceMapping{
		Ebs: &gfnec2.LaunchTemplate_Ebs{
			VolumeSize: gfnt.NewInteger(*volumeSize),
			VolumeType: gfnt.NewString(*vm.VolumeType),
		},
	}
	if vm.VolumeEncrypted != nil {
		mapping.Ebs.Encrypted = gfnt.NewBoolean(*vm.VolumeEncrypted)
	}
	if api.IsSetAndNonEmptyString(vm.VolumeKmsKeyID) {
		mapping.Ebs.KmsKeyId = gfnt.NewString(*vm.VolumeKmsKeyID)
	}

	if (*vm.VolumeType == api.NodeVolumeTypeIO1 || *vm.VolumeType == api.NodeVolumeTypeGP3) && vm.VolumeIOPS != nil {
		mapping.Ebs.Iops = gfnt.NewInteger(*vm.VolumeIOPS)
	}

	if *vm.VolumeType == api.NodeVolumeTypeGP3 && vm.VolumeThroughput != nil {
		mapping.Ebs.Throughput = gfnt.NewInteger(*vm.VolumeThroughput)
	}

	if vm.VolumeName != nil {
		mapping.DeviceName = gfnt.NewString(*vm.VolumeName)
	} else {
		mapping.DeviceName = gfnt.NewString(rootDevice)
	}

	return &mapping
}
