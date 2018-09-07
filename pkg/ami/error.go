package ami

import (
	"fmt"
)

// ErrFailedAMIResolution is an error type that represents
// failure to resolve a region/instance type to an AMI
type ErrFailedAMIResolution struct {
	region       string
	instanceType string
}

// NewErrFailedAMIResolution creates a new instance of ErrFailedAMIResolution for a
// give region and instance type
func NewErrFailedAMIResolution(region string, instanceType string) *ErrFailedAMIResolution {
	return &ErrFailedAMIResolution{
		region:       region,
		instanceType: instanceType,
	}
}

// Error return the error message
func (e *ErrFailedAMIResolution) Error() string {
	return fmt.Sprintf("Unable to determine AMI for region %s and instance type %s", e.region, e.instanceType)
}

// ErrAmiNotFound is an error type that represents
// failure to find a given ami
type ErrAmiNotFound struct {
	ami string
}

// NewErrAmiNotFound creates a new instance of ErrAmiNotFound for a
// given ami
func NewErrAmiNotFound(ami string) *ErrAmiNotFound {
	return &ErrAmiNotFound{
		ami: ami,
	}
}

// Error return the error message
func (e *ErrAmiNotFound) Error() string {
	return fmt.Sprintf("Unable to find AMI %s", e.ami)
}
