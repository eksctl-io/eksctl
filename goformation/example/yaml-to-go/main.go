package main

import (
	"log"

	"github.com/awslabs/goformation/v4"
)

func main() {

	// Open a template from file (can be JSON or YAML)
	template, err := goformation.Open("template.yaml")
	if err != nil {
		log.Fatalf("There was an error processing the template: %s", err)
	}

	// You can extract all resources of a certain type
	// Each AWS CloudFormation resource is a strongly typed struct
	topics := template.GetAllSNSTopicResources()
	for name, topic := range topics {

		// E.g. Found a AWS::SNS::Topic with Logical ID ExampleTopic and TopicName 'example'
		log.Printf("Found a %s with Logical ID %s and TopicName %s\n", topic.AWSCloudFormationType(), name, topic.TopicName)

	}

	// You can also search for specific resources by their logicalId
	search := "ExampleTopic"
	topic, err := template.GetSNSTopicWithName(search)
	if err != nil {
		log.Fatalf("SNS topic with logical ID %s not found", search)
	}

	// E.g. Found a AWS::Serverless::Function named GetHelloWorld (runtime: nodejs6.10)
	log.Printf("Found a %s with Logical ID %s and TopicName %s\n", topic.AWSCloudFormationType(), search, topic.TopicName)

}
