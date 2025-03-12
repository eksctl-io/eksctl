package v1alpha5

import "slices"

var KnownAddons = map[string]struct {
	IsDefault             bool
	CreateBeforeNodeGroup bool
	IsDefaultAutoMode     bool
	ExcludedRegions       []string
}{
	VPCCNIAddon: {
		IsDefault:             true,
		CreateBeforeNodeGroup: true,
	},
	KubeProxyAddon: {
		IsDefault:             true,
		CreateBeforeNodeGroup: true,
	},
	CoreDNSAddon: {
		IsDefault:             true,
		CreateBeforeNodeGroup: true,
	},
	PodIdentityAgentAddon: {
		CreateBeforeNodeGroup: true,
	},
	AWSEBSCSIDriverAddon: {},
	AWSEFSCSIDriverAddon: {},
	MetricsServerAddon: {
		IsDefault:             true,
		CreateBeforeNodeGroup: true,
		IsDefaultAutoMode:     true,
		ExcludedRegions: []string{
			RegionCNNorthwest1,
			RegionCNNorth1,
			RegionUSISOEast1,
			RegionUSISOWest1,
			RegionUSISOBEast1,
			RegionUSGovWest1,
			RegionUSGovEast1,
			RegionUSISOFEast1,
			RegionUSISOFSouth1,
			RegionEUISOEWest1,
		},
	},
}

// HasDefaultNonAutoAddon reports whether addons contains at least one non-auto mode default addon
func HasDefaultNonAutoAddon(addons []*Addon) bool {
	for _, addon := range addons {
		addonConfig, ok := KnownAddons[addon.Name]
		if ok && addonConfig.IsDefault && !addonConfig.IsDefaultAutoMode {
			return true
		}
	}
	return false
}

// HasAllDefaultAddons reports whether addonNames contains all default addons.
func HasAllDefaultAddons(addonNames []string) bool {
	for addonName, addon := range KnownAddons {
		if addon.IsDefault && !slices.Contains(addonNames, addonName) {
			return false
		}
	}
	return true
}
