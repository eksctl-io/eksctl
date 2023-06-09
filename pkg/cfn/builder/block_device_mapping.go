package builder

import (
	gfnec2 "github.com/weaveworks/goformation/v4/cloudformation/ec2"
	gfnt "github.com/weaveworks/goformation/v4/cloudformation/types"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
)

const (
	rootDeviceForLinux   = "/dev/xvda"
	rootDeviceForWindows = "/dev/sda1"
)

func makeBlockDeviceMappings(ng *api.NodeGroupBase) []gfnec2.LaunchTemplate_BlockDeviceMapping {
	var (
		mappings       []gfnec2.LaunchTemplate_BlockDeviceMapping
		baseVolumeName string
	)

	if ng.VolumeName != nil {
		baseVolumeName = *ng.VolumeName
	} else if api.IsWindowsImage(ng.AMIFamily) {
		baseVolumeName = rootDeviceForWindows
	} else {
		baseVolumeName = rootDeviceForLinux
	}

	mapping := makeBlockDeviceMapping(&api.VolumeMapping{
		VolumeSize:       ng.VolumeSize,
		VolumeType:       ng.VolumeType,
		VolumeName:       &baseVolumeName,
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
				Encrypted:  gfnt.NewBoolean(*ng.VolumeEncrypted),
				KmsKeyId:   mapping.Ebs.KmsKeyId,
				VolumeType: mapping.Ebs.VolumeType,
			},
		})
	}

	return mappings
}

func makeBlockDeviceMapping(vm *api.VolumeMapping) *gfnec2.LaunchTemplate_BlockDeviceMapping {
	if vm.VolumeSize == nil || *vm.VolumeSize == 0 {
		return nil
	}

	if !api.IsSetAndNonEmptyString(vm.VolumeName) {
		return nil
	}

	mapping := gfnec2.LaunchTemplate_BlockDeviceMapping{
		Ebs: &gfnec2.LaunchTemplate_Ebs{
			VolumeSize: gfnt.NewInteger(*vm.VolumeSize),
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

	if api.IsSetAndNonEmptyString(vm.SnapshotID) {
		mapping.Ebs.SnapshotId = gfnt.NewString(*vm.SnapshotID)
	}

	mapping.DeviceName = gfnt.NewString(*vm.VolumeName)

	return &mapping
}
