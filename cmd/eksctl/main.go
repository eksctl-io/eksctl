package main

import (
	"fmt"

	"github.com/weaveworks/eksctl/pkg/cfn"
)

func main() {

	cfn := cfn.New()

	{
		s, _ := cfn.SubStackVPC("cluster-1")
		fmt.Println(s)
	}
}
