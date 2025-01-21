package main

import (
	"errors"
	"strings"
)

// filename takes a resource or property name (e.g. AWS::CloudFront::Distribution.Restrictions)
// and returns an appropriate filename for the generated struct (e.g. aws-cloudfront-distribution_restrictions.go)
func filename(input string) string {

	// Convert to lowercase
	output := strings.ToLower(input)

	// Replace :: with -
	output = strings.Replace(output, "::", "-", -1)

	// Replace . with _
	output = strings.Replace(output, ".", "_", -1)

	// Suffix with .go
	output += ".go"

	return output

}

// structName takes a resource or property name (e.g. AWS::CloudFront::Distribution.Restrictions)
// and returns an appropriate struct name for the generated struct (e.g. AWSCloudfrontDistributionRestrictions)
func structName(input string) (string, error) {

	// Remove ::
	//output := strings.Replace(input, "::", "", -1)

	if input == "Tag" {
		return "Tag", nil
	}

	parts := strings.Split(input, "::")
	if len(parts) < 2 {
		return "", errors.New("invalid CloudFormation resource type: " + input)
	}

	// Remove .
	output := strings.Replace(parts[2], ".", "_", -1)

	return output, nil

}

// packageName generates a go package name based on the AWS CloudFormation resource type
// For example, AWS::S3::Bucket would generate a package named 's3'
func packageName(input string, lowercase bool) (string, error) {

	if input == "Tag" {
		return "cloudformation", nil
	}

	parts := strings.Split(input, "::")
	if len(parts) < 2 {
		return "", errors.New("invalid CloudFormation resource type: " + input)
	}

	if lowercase {
		return strings.ToLower(parts[1]), nil
	}

	return parts[1], nil

}
