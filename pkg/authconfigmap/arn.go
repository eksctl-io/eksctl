package authconfigmap

import "github.com/aws/aws-sdk-go/aws/arn"

// Implement the pflag.Value interface for arn.ARN
type ARN arn.ARN

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

func (a *ARN) String() string {
	tmp := arn.ARN(*a)
	return tmp.String()
}

func (a *ARN) Type() string {
	return "aws arn"
}
