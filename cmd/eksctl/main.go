package main

import (
	"fmt"

	"github.com/weaveworks/eksctl/pkg/cfn"
)

const DEFAULT_EKS_REGION = "us-west-2"

func main() {

	cfn := cfn.New(DEFAULT_EKS_REGION)

	if err := cfn.CheckAuth(); err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

	// {
	// 	s, _ := cfn.CreateStackVPC("cluster-1")
	// 	fmt.Println(s)
	// }
}
