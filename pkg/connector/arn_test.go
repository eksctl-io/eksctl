package connector_test

import (
	"fmt"
	"testing"

	"github.com/weaveworks/eksctl/pkg/connector"
)

var arnTests = []struct {
	arn      string // input arn
	expected string // canonacalized arn
	err      error  // expected error value
}{
	{"NOT AN ARN", "", fmt.Errorf("Not an arn")},
	{"arn:aws:iam::123456789012:user/Alice", "arn:aws:iam::123456789012:user/Alice", nil},
	{"arn:aws:iam::123456789012:role/Users", "arn:aws:iam::123456789012:role/Users", nil},
	{"arn:aws:sts::123456789012:assumed-role/Admin/Session", "arn:aws:iam::123456789012:role/Admin", nil},
	{"arn:aws:sts::123456789012:federated-user/Bob", "arn:aws:sts::123456789012:federated-user/Bob", nil},
	{"arn:aws:iam::123456789012:root", "arn:aws:iam::123456789012:root", nil},
	{"arn:aws:sts::123456789012:assumed-role/Org/Team/Admin/Session", "arn:aws:iam::123456789012:role/Org/Team/Admin", nil},
	{"arn:aws-iso:iam::123456789012:user/Chris", "arn:aws-iso:iam::123456789012:user/Chris", nil},
	{"arn:aws-iso-b:iam::123456789012:user/Chris", "arn:aws-iso-b:iam::123456789012:user/Chris", nil},
}

func TestUserARN(t *testing.T) {
	for _, tc := range arnTests {
		actual, err := connector.Canonicalize(tc.arn)
		if err != nil && tc.err == nil || err == nil && tc.err != nil {
			t.Errorf("Canoncialize(%s) expected err: %v, actual err: %v", tc.arn, tc.err, err)
			continue
		}
		if actual != tc.expected {
			t.Errorf("Canonicalize(%s) expected: %s, actual: %s", tc.arn, tc.expected, actual)
		}
	}
}
