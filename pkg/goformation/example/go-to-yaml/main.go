package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/sns"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"
)

func main() {

	// Create a new CloudFormation template
	template := cloudformation.NewTemplate()

	// Create an Amazon SNS topic, with a unique name based off the current timestamp
	template.Resources["MyTopic"] = &sns.Topic{
		TopicName: types.NewString("my-topic-" + strconv.FormatInt(time.Now().Unix(), 10)),
	}

	// Create a subscription, connected to our topic, that forwards notifications to an email address
	template.Resources["MyTopicSubscription"] = &sns.Subscription{
		TopicArn: types.MakeRef("MyTopic"),
		Protocol: types.NewString("email"),
		Endpoint: types.NewString("some.email@example.com"),
	}

	// Let's see the JSON
	j, err := template.JSON()
	if err != nil {
		fmt.Printf("Failed to generate JSON: %s\n", err)
	} else {
		fmt.Printf("%s\n", string(j))
	}

	y, err := template.YAML()
	if err != nil {
		fmt.Printf("Failed to generate YAML: %s\n", err)
	} else {
		fmt.Printf("%s\n", string(y))
	}

}
