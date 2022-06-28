package v1alpha5

import (
	"github.com/pkg/errors"
)

// IAMIdentityMapping contains IAM accounts, users, roles and services that will be added to the
// aws-auth configmap to enable access to the cluster
type IAMIdentityMapping struct {
	// +optional
	ARN             string   `json:"arn,omitempty"`
	Username        string   `json:"username,omitempty"`
	Groups          []string `json:"groups,omitempty"`
	Account         string   `json:"account,omitempty"`
	ServiceName     string   `json:"serviceName,omitempty"`
	Namespace       string   `json:"namespace,omitempty"`
	NoDuplicateARNs bool     `json:"noDuplicateARNs,omitempty"`
}

func (im *IAMIdentityMapping) hasARNOptions() bool {
	return im.hasARN() && im.hasUsername() && im.hasGroups()
}

func (im *IAMIdentityMapping) hasARN() bool {
	return im.ARN != ""
}

func (im *IAMIdentityMapping) hasUsername() bool {
	return im.Username != ""
}

func (im *IAMIdentityMapping) hasGroups() bool {
	return len(im.Groups) > 0
}

func (im *IAMIdentityMapping) Validate() error {

	if im.ServiceName != "" {
		if im.hasARNOptions() {
			return errors.New("cannot use arn, username, and groups with serviceName")
		}
		if im.Namespace == "" {
			return errors.New("namespace is required when using serviceName")
		}
	} else {
		if im.Namespace != "" {
			return errors.New("serviceName is required when using namespace")
		}
	}

	if im.Account != "" && (im.hasARN() || im.hasUsername() || im.hasGroups() || im.ServiceName != "" || im.Namespace != "") {
		return errors.New("account cannot be configured with any other options")
	}

	return nil
}
