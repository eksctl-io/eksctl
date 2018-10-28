package ami

import (
	"fmt"
)

// ErrFailedResolution is an error type that represents
// failure to resolve a region/instance type/image family to an AMI
type ErrFailedResolution struct {
	region       string
	instanceType string
	imageFamily  string
}

// NewErrFailedResolution creates a new instance of ErrFailedResolution for a
// give region, instance type and image family
func NewErrFailedResolution(region string, instanceType string, imageFamily string) *ErrFailedResolution {
	return &ErrFailedResolution{
		region:       region,
		instanceType: instanceType,
		imageFamily:  imageFamily,
	}
}

// Error return the error message
func (e *ErrFailedResolution) Error() string {
	return fmt.Sprintf("Unable to determine AMI for region %s, instance type %s and image family %s", e.region, e.instanceType, e.imageFamily)
}

// ErrNotFound is an error type that represents
// failure to find a given ami
type ErrNotFound struct {
	ami string
}

// NewErrNotFound creates a new instance of ErrNotFound for a
// given ami
func NewErrNotFound(ami string) *ErrNotFound {
	return &ErrNotFound{
		ami: ami,
	}
}

// Error return the error message
func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("Unable to find AMI %s", e.ami)
}
