package cloudformation

import (
	"encoding/base64"
	"encoding/json"
	"strings"
)

// Note: Intrinsic objects are Base64 encoded, to prevent escaping (backslash) issues
// with nested intrinsic functions.

// Ref creates a CloudFormation Reference to another resource in the template
func Ref(logicalName string) string {
	i := `{ "Ref": "` + logicalName + `" }`
	return base64.StdEncoding.EncodeToString([]byte(i))
}

// GetAtt returns the value of an attribute from a resource in the template.
func GetAtt(logicalName string, attribute string) string {
	i := `{ "Fn::GetAtt": [ "` + logicalName + `", "` + attribute + `" ] }`
	return base64.StdEncoding.EncodeToString([]byte(i))
}

// ImportValue returns the value of an output exported by another stack. You typically use this function to create cross-stack references. In the following example template snippets, Stack A exports VPC security group values and Stack B imports them.
func ImportValue(name string) string {
	i := `{ "Fn::ImportValue": "` + name + `" }`
	return base64.StdEncoding.EncodeToString([]byte(i))
}

// Base64 returns the Base64 representation of the input string. This function is typically used to pass encoded data to Amazon EC2 instances by way of the UserData property
func Base64(input string) string {
	i := `{ "Fn::Base64": "` + input + `" }`
	return base64.StdEncoding.EncodeToString([]byte(i))
}

// CIDR returns an array of CIDR address blocks. The number of CIDR blocks returned is dependent on the count parameter.
func CIDR(ipBlock, count, cidrBits string) string {
	i := `{ "Fn::Cidr" : [ "` + ipBlock + `", "` + count + `", "` + cidrBits + `" ] }`
	return base64.StdEncoding.EncodeToString([]byte(i))
}

// FindInMap returns the value corresponding to keys in a two-level map that is declared in the Mappings section.
func FindInMap(mapName, topLevelKey, secondLevelKey string) string {
	i := `{ "Fn::FindInMap" : [ "` + mapName + `", "` + topLevelKey + `", "` + secondLevelKey + `" ] }`
	return base64.StdEncoding.EncodeToString([]byte(i))
}

// GetAZs returns an array that lists Availability Zones for a specified region. Because customers have access to different Availability Zones, the intrinsic function Fn::GetAZs enables template authors to write templates that adapt to the calling user's access. That way you don't have to hard-code a full list of Availability Zones for a specified region.
func GetAZs(region string) string {
	i := `{ "Fn::GetAZs": "` + region + `" }`
	return base64.StdEncoding.EncodeToString([]byte(i))
}

// Join appends a set of values into a single value, separated by the specified delimiter. If a delimiter is the empty string, the set of values are concatenated with no delimiter.
func Join(delimiter string, values []string) string {
	i := `{ "Fn::Join": [ "` + delimiter + `", [ "` + strings.Trim(strings.Join(values, `", "`), `, "`) + `" ] ] }`
	return base64.StdEncoding.EncodeToString([]byte(i))
}

// Select returns a single object from a list of objects by index.
func Select(index string, list []string) string {
	i := `{ "Fn::Select": [ "` + index + `", [ "` + strings.Trim(strings.Join(list, `", "`), `, "`) + `" ] ] }`
	return base64.StdEncoding.EncodeToString([]byte(i))
}

// Split splits a string into a list of string values so that you can select an element from the resulting string list, use the Fn::Split intrinsic function. Specify the location of splits with a delimiter, such as , (a comma). After you split a string, use the Fn::Select function to pick a specific element.
func Split(delimiter, source string) string {
	i := `{ "Fn::Split" : [ "` + delimiter + `", "` + source + `" ] }`
	return base64.StdEncoding.EncodeToString([]byte(i))
}

// Sub substitutes variables in an input string with values that you specify. In your templates, you can use this function to construct commands or outputs that include values that aren't available until you create or update a stack.
func Sub(value string) string {
	i := `{ "Fn::Sub" : "` + value + `" }`
	return base64.StdEncoding.EncodeToString([]byte(i))
}

// processIntrinsics is a post processor that hydrates all intrinsic functions in the template
func processIntrinsics(input interface{}) (interface{}, error) {

	// Marshal to JSON and back to convert from a typed template object to simple primitives
	b, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	// Recurse through the object tree, replacing any Goformation references
	return replaceIntrinsicsRecursive(m), nil

}

// replaceReferencesRecursive recurses through an object, and replaces any strings that
// contain '%%Ref:(.*)%%' with a proper AWS CloudFormation reference object
func replaceIntrinsicsRecursive(input interface{}) interface{} {

	switch value := input.(type) {

	case map[string]interface{}:
		result := map[string]interface{}{}
		for k, v := range value {
			result[k] = replaceIntrinsicsRecursive(v)
		}
		return result

	case []interface{}:
		result := []interface{}{}
		for _, v := range value {
			result = append(result, replaceIntrinsicsRecursive(v))
		}
		return result

	case string:

		// Check if the string can be unmarshalled into an intrinsic object
		var decoded []byte
		decoded, err := base64.StdEncoding.DecodeString(value)
		if err != nil {
			// The string value is not base64 encoded, so it's not an intrinsic so just pass it back
			return value
		}

		var intrinsic map[string]interface{}
		if err := json.Unmarshal([]byte(decoded), &intrinsic); err != nil {
			// The string value is not JSON, so it's not an intrinsic so just pass it back
			return value
		}

		// An intrinsic should be an object, with a single key containing a valid intrinsic name
		if len(intrinsic) != 1 {
			return value
		}

		supported := []string{
			"Ref",
			"Fn::Base64",
			"Fn::Cidr",
			"Fn::FindInMap",
			"Fn::GetAtt",
			"Fn::GetAZs",
			"Fn::ImportValue",
			"Fn::Join",
			"Fn::Select",
			"Fn::Split",
			"Fn::Sub",
			"Fn::Transform",
		}

		for name, v := range intrinsic {
			for _, i := range supported {
				if name == i {
					return map[string]interface{}{
						name: replaceIntrinsicsRecursive(v),
					}
				}
			}
		}

		return value

	default:
		return value
	}

}
