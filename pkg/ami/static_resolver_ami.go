package ami

// StaticImages is a map that holds the list of amis to be used by
// for static ami resolution
var StaticImages = map[string]map[int]map[string]string{
	ImageFamilyAmazonLinux2: {
		ImageClassGeneral: {
			"eu-west-1": "ami-0c7a4976cb6fafd3a",
			"us-east-1": "ami-0440e4f6b9713faf6",
			"us-west-2": "ami-0a54c984b9f908c81",
		},
		ImageClassGPU: {
			"eu-west-1": "ami-0706dc8a5eed2eed9",
			"us-east-1": "ami-058bfb8c236caae89",
			"us-west-2": "ami-0731694d53ef9604b",
		},
	},
}
