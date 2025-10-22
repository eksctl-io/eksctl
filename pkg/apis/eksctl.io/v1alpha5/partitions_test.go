package v1alpha5

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestV1SDKDNSPrefixReturnsForAllRegions(t *testing.T) {
	for _, region := range SupportedRegions() {
		_, err := Partitions.V1SDKDNSPrefixForRegion(region)
		require.NoError(t, err)
	}
}

func TestV1SDKDNSPrefix(t *testing.T) {
	tests := []struct {
		region   string
		expected string
	}{
		{RegionUSWest1, "amazonaws.com"},
		{RegionUSGovEast1, "amazonaws.com"},
		{RegionCNNorthwest1, "amazonaws.com.cn"},
		{RegionUSISOEast1, "c2s.ic.gov"},
		{RegionUSISOBEast1, "sc2s.sgov.gov"},
		{RegionEUISOEWest1, "cloud.adc-e.uk"},
		{RegionUSISOFSouth1, "csp.hci.ic.gov"},
	}
	for _, test := range tests {
		dns, err := Partitions.V1SDKDNSPrefixForRegion(test.region)
		require.NoError(t, err)
		require.Equal(t, test.expected, dns)
	}
}

// Commented out so we don't depend on v1 sdk.
//func TestAgainstV1SDK(t *testing.T) {
//	var partitionNames []string
//	for _, p := range Partitions {
//		partitionNames = append(partitionNames, p.Name())
//	}
//	var expected []string
//	for _, p := range endpoints.DefaultPartitions() {
//		expected = append(expected, p.ID())
//	}
//	require.ElementsMatch(t, expected, partitionNames)
//
//	for _, r := range SupportedRegions() {
//		var expected string
//		for _, p := range endpoints.DefaultPartitions() {
//			_, ok := p.Regions()[r]
//			if ok {
//				expected = p.DNSSuffix()
//				break
//			}
//		}
//		if expected == "" { // SDK version may be missing some endpoints.
//			continue
//		}
//
//		actual, err := Partitions.V1SDKDNSPrefixForRegion(r)
//		require.NoError(t, err)
//		require.Equalf(t, expected, actual, "wrong dns for region %s", r)
//	}
//}
