// Package ipnet wraps net.IPNet to get CIDR serialization.
package ipnet

// This was code was copied from [1] and cobined with version in master [2].
// [1]: https://github.com/wking/openshift-installer/commit/06739c665e9d3bc9b813ad873bed862a5bde24f1
// [2]: https://github.com/openshift/installer/tree/5e7b36d6351c9cc773f1dadc64abf9d7041151b1/pkg/ipnet
// TODO: this is not ideal, we should move this out or do something else about it.

import (
	"encoding/json"
	"net"
	"reflect"

	"github.com/pkg/errors"
)

var nullString = "null"
var nullBytes = []byte(nullString)
var emptyIPNet = net.IPNet{}

// IPNet wraps net.IPNet to get CIDR serialization.
type IPNet struct {
	net.IPNet
}

// String returns a CIDR serialization of the subnet, or an empty
// string if the subnet is nil.
func (ipnet *IPNet) String() string {
	if ipnet == nil {
		return ""
	}
	return ipnet.IPNet.String()
}

// DeepCopyInto copies the receiver into out.  out must be non-nil.
func (ipnet *IPNet) DeepCopyInto(out *IPNet) {
	if ipnet == nil {
		*out = *new(IPNet)
	} else {
		*out = *ipnet
	}
}

// DeepCopy copies the receiver, creating a new IPNet.
func (ipnet *IPNet) DeepCopy() *IPNet {
	if ipnet == nil {
		return nil
	}
	out := new(IPNet)
	ipnet.DeepCopyInto(out)
	return out
}

// MarshalJSON interface for an IPNet
func (ipnet IPNet) MarshalJSON() (data []byte, err error) {
	if reflect.DeepEqual(ipnet.IPNet, emptyIPNet) {
		return nullBytes, nil
	}

	return json.Marshal(ipnet.String())
}

// UnmarshalJSON interface for an IPNet
func (ipnet *IPNet) UnmarshalJSON(b []byte) (err error) {
	if string(b) == nullString {
		ipnet.IP = net.IP{}
		ipnet.Mask = net.IPMask{}
		return nil
	}

	var cidr string
	err = json.Unmarshal(b, &cidr)
	if err != nil {
		return errors.Wrap(err, "failed to Unmarshal string")
	}

	ip, net, err := net.ParseCIDR(cidr)
	if err != nil {
		return errors.Wrap(err, "failed to Parse cidr string to net.IPNet")
	}

	// This check is needed in order to work around a strange quirk in the Go
	// standard library. All of the addresses returned by net.ParseCIDR() are
	// 16-byte addresses. This does _not_ imply that they are IPv6 addresses,
	// which is what some libraries (e.g. github.com/apparentlymart/go-cidr)
	// assume. By forcing the address to be the expected length, we can work
	// around these bugs.
	if ip.To4() != nil {
		ipnet.IP = ip.To4()
	} else {
		ipnet.IP = ip
	}
	ipnet.Mask = net.Mask

	return nil
}

// ParseCIDR parses a CIDR from its string representation.
func ParseCIDR(s string) (*IPNet, error) {
	_, cidr, err := net.ParseCIDR(s)
	if err != nil {
		return nil, err
	}
	return &IPNet{IPNet: *cidr}, nil
}

// MustParseCIDR parses a CIDR from its string representation. If the parse fails,
// the function will panic.
func MustParseCIDR(s string) *IPNet {
	cidr, err := ParseCIDR(s)
	if err != nil {
		panic(err)
	}
	return cidr
}
