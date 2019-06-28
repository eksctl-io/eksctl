package authconfigmap

import "github.com/aws/aws-sdk-go/aws/arn"

// ARN implements the pflag.Value interface for aws-sdk-go/aws/arn.ARN
type ARN arn.ARN

// Parse wraps the aws-sdk-go/aws/arn.Parse function and instead returns a
// authconfigmap.ARN
func Parse(s string) (ARN, error) {
	a, err := arn.Parse(s)
	return ARN(a), err
}

// Set parses the given string into an arn.ARN and sets the receiver pointer to the
// populated struct
func (a *ARN) Set(s string) error {
	arn, err := arn.Parse(s)
	if err != nil {
		return err
	}
	*a = ARN(arn)
	return nil
}

// String returns the canonical representation of the ARN
func (a *ARN) String() string {
	tmp := arn.ARN(*a)
	return tmp.String()
}

// Type returns a string representation that describes the type
func (a *ARN) Type() string {
	return "aws arn"
}
