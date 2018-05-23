// +build example

package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudformation"
)

func main() {

	clientSession := session.Must(session.NewSession())

	cfn := cloudformation.New(clientSession)

	listStacksInput := &cloudformation.ListStacksInput{}

	i := 0
	pager := func(p *cloudformation.ListStacksOutput, last bool) (shouldContinue bool) {
		fmt.Println("Page,", i)
		i++

		for _, s := range p.StackSummaries {
			fmt.Println("stack:", *s)
		}
		return true
	}
	err := cfn.ListStacksPages(listStacksInput, pager)
	if err != nil {
		fmt.Println("failed to list stacks", err)
		return
	}
}
