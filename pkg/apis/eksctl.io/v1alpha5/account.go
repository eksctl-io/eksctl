package v1alpha5

import (
	"os"

	"github.com/pkg/errors"
)

const (
	IsoEKSAccountIDEnv = "ISO_EKS_ACCOUNT_ID"

	// eksResourceAccountStandard defines the AWS EKS account ID that provides node resources in default regions
	// for standard AWS partition
	eksResourceAccountStandard = "602401143452"

	// eksResourceAccountAPEast1 defines the AWS EKS account ID that provides node resources in ap-east-1 region
	eksResourceAccountAPEast1 = "800184023465"

	// eksResourceAccountMESouth1 defines the AWS EKS account ID that provides node resources in me-south-1 region
	eksResourceAccountMESouth1 = "558608220178"

	// eksResourceAccountCNNorthWest1 defines the AWS EKS account ID that provides node resources in cn-northwest-1 region
	eksResourceAccountCNNorthWest1 = "961992271922"

	// eksResourceAccountCNNorth1 defines the AWS EKS account ID that provides node resources in cn-north-1
	eksResourceAccountCNNorth1 = "918309763551"

	// eksResourceAccountAFSouth1 defines the AWS EKS account ID that provides node resources in af-south-1
	eksResourceAccountAFSouth1 = "877085696533"

	// eksResourceAccountEUSouth1 defines the AWS EKS account ID that provides node resources in eu-south-1
	eksResourceAccountEUSouth1 = "590381155156"

	// eksResourceAccountUSGovWest1 defines the AWS EKS account ID that provides node resources in us-gov-west-1
	eksResourceAccountUSGovWest1 = "013241004608"

	// eksResourceAccountUSGovEast1 defines the AWS EKS account ID that provides node resources in us-gov-east-1
	eksResourceAccountUSGovEast1 = "151742754352"
)

// EKSResourceAccountID provides worker node resources(ami/ecr image) in different aws account
// for different aws partitions & opt-in regions.
func EKSResourceAccountID(region string) (string, error) {
	switch region {
	case RegionUSIsoEast1:
		return lookupIsoAccountID(region)
	case RegionUSIsobEast1:
		return lookupIsoAccountID(region)
	default:
		return publicResourceAccountID(region), nil
	}
}

// ValidateEKSAccountID check that the ISO_EKS_ACCOUNT_ID env var has been set
// when an isolated region is requested
func ValidateEKSAccountID(region string) error {
	if region == RegionUSIsoEast1 || region == RegionUSIsobEast1 {
		if _, err := lookupIsoAccountID(region); err != nil {
			return err
		}
	}

	return nil
}

func lookupIsoAccountID(region string) (string, error) {
	id := os.Getenv(IsoEKSAccountIDEnv)
	if id == "" {
		return "", errors.Errorf("%s not set, required for use of region: %s", IsoEKSAccountIDEnv, region)
	}

	return id, nil
}

func publicResourceAccountID(region string) string {
	switch region {
	case RegionAPEast1:
		return eksResourceAccountAPEast1
	case RegionMESouth1:
		return eksResourceAccountMESouth1
	case RegionCNNorthwest1:
		return eksResourceAccountCNNorthWest1
	case RegionCNNorth1:
		return eksResourceAccountCNNorth1
	case RegionUSGovWest1:
		return eksResourceAccountUSGovWest1
	case RegionUSGovEast1:
		return eksResourceAccountUSGovEast1
	case RegionAFSouth1:
		return eksResourceAccountAFSouth1
	case RegionEUSouth1:
		return eksResourceAccountEUSouth1
	default:
		return eksResourceAccountStandard
	}
}
