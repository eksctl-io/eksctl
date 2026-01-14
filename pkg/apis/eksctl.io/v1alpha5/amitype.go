package v1alpha5

import (
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	instanceutils "github.com/weaveworks/eksctl/pkg/utils/instance"
)

// GetAMIType returns the most appropriate amiType for the amiFamily and
// instanceType provided. Parameter `strict` controls whether or not fallbacks
// should be applied when searching for specialized amiTypes (eg. accerated
// instance types). If `strict` is false a fallback may be applied, otherwise a
// valid value is not guaranteed to be returned (empty string).
func GetAMIType(amiFamily, instanceType string, strict bool) ekstypes.AMITypes {
	amiTypeMapping := map[string]struct {
		X86x64      ekstypes.AMITypes
		X86Nvidia   ekstypes.AMITypes
		X86Neuron   ekstypes.AMITypes
		ARM         ekstypes.AMITypes
		ARM64Nvidia ekstypes.AMITypes
		ARM64Neuron ekstypes.AMITypes
	}{
		NodeImageFamilyAmazonLinux2023: {
			X86x64:      ekstypes.AMITypesAl2023X8664Standard,
			X86Nvidia:   ekstypes.AMITypesAl2023X8664Nvidia,
			X86Neuron:   ekstypes.AMITypesAl2023X8664Neuron,
			ARM:         ekstypes.AMITypesAl2023Arm64Standard,
			ARM64Nvidia: ekstypes.AMITypesAl2023Arm64Nvidia,
		},
		NodeImageFamilyAmazonLinux2: {
			X86x64:    ekstypes.AMITypesAl2X8664,
			X86Nvidia: ekstypes.AMITypesAl2X8664Gpu,
			X86Neuron: ekstypes.AMITypesAl2X8664Gpu,
			ARM:       ekstypes.AMITypesAl2Arm64,
		},
		NodeImageFamilyBottlerocket: {
			X86x64:      ekstypes.AMITypesBottlerocketX8664,
			X86Nvidia:   ekstypes.AMITypesBottlerocketX8664Nvidia,
			X86Neuron:   ekstypes.AMITypesBottlerocketX8664,
			ARM:         ekstypes.AMITypesBottlerocketArm64,
			ARM64Nvidia: ekstypes.AMITypesBottlerocketArm64Nvidia,
			ARM64Neuron: ekstypes.AMITypesBottlerocketArm64,
		},
		NodeImageFamilyWindowsServer2019FullContainer: {
			X86x64:    ekstypes.AMITypesWindowsFull2019X8664,
			X86Nvidia: ekstypes.AMITypesWindowsFull2019X8664,
		},
		NodeImageFamilyWindowsServer2019CoreContainer: {
			X86x64:    ekstypes.AMITypesWindowsCore2019X8664,
			X86Nvidia: ekstypes.AMITypesWindowsCore2019X8664,
		},
		NodeImageFamilyWindowsServer2022FullContainer: {
			X86x64:    ekstypes.AMITypesWindowsFull2022X8664,
			X86Nvidia: ekstypes.AMITypesWindowsFull2022X8664,
		},
		NodeImageFamilyWindowsServer2022CoreContainer: {
			X86x64:    ekstypes.AMITypesWindowsCore2022X8664,
			X86Nvidia: ekstypes.AMITypesWindowsCore2022X8664,
		},
		NodeImageFamilyWindowsServer2025FullContainer: {
			X86x64:    ekstypes.AMITypes("WINDOWS_FULL_2025_x86_64"),
			X86Nvidia: ekstypes.AMITypes("WINDOWS_FULL_2025_x86_64"),
		},
		NodeImageFamilyWindowsServer2025CoreContainer: {
			X86x64:    ekstypes.AMITypes("WINDOWS_CORE_2025_x86_64"),
			X86Nvidia: ekstypes.AMITypes("WINDOWS_CORE_2025_x86_64"),
		},
	}

	amiType, ok := amiTypeMapping[amiFamily]
	if !ok {
		return ekstypes.AMITypesCustom
	}

	// this helper short circuits the check for missing entries for amiTypes in
	// ami families based on the value of `strict`.
	isValid := func(amiType ekstypes.AMITypes) bool {
		return strict || amiType != ""
	}

	if instanceutils.IsARMInstanceType(instanceType) {
		switch {
		case instanceutils.IsNvidiaInstanceType(instanceType) && isValid(amiType.ARM64Nvidia):
			return amiType.ARM64Nvidia
		case instanceutils.IsNeuronInstanceType(instanceType) && isValid(amiType.ARM64Neuron):
			return amiType.ARM64Neuron
		default:
			return amiType.ARM
		}
	} else {
		switch {
		case instanceutils.IsNvidiaInstanceType(instanceType) && isValid(amiType.X86Nvidia):
			return amiType.X86Nvidia
		case instanceutils.IsNeuronInstanceType(instanceType) && isValid(amiType.X86Neuron):
			return amiType.X86Neuron
		default:
			return amiType.X86x64
		}
	}
}
