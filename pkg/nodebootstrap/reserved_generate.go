//+build ignore

// Fetches AWS instance type data such as memory, cpus, and
// (ephemeral) storage. Writes it to a Go file to be included
// by eksctl.
//
// Since each region may support different instance types we
// need to query all regions. But the same instance type in
// different regions has consistent specs.
package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/dave/jennifer/jen"

	"github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

const outputFilename = "reserved_data.go"

func main() {
	infos := map[string]nodebootstrap.InstanceTypeInfo{}
	for _, region := range v1alpha5.SupportedRegions() {
		fmt.Printf("Fetching instance types in %q\n", region)
		client := ec2.New(newSession(region))
		if err := updateRegionInstanceTypes(client, infos); err != nil {
			checkError(region, err)
		}
	}
	if err := render(infos); err != nil {
		panic(err)
	}
	fmt.Println("Done.")
}

func checkError(region string, err error) {
	if err == nil {
		return
	}
	// Allow failing for Chinese regions
	if aerr, ok := err.(awserr.Error); ok {
		if aerr.Code() == "AuthFailure" && strings.HasPrefix(region, "cn-") {
			fmt.Println("  AuthFailure, skipping.")
			return
		}
	}
	fmt.Println("Make sure your environment provides credentials to run ec2:DescribeInstanceTypes in supported regions")
	panic(err)
}

func updateRegionInstanceTypes(client *ec2.EC2, infos map[string]nodebootstrap.InstanceTypeInfo) error {
	var token *string
	for {
		resp, err := client.DescribeInstanceTypes(&ec2.DescribeInstanceTypesInput{NextToken: token})
		if err != nil {
			return err
		}
		for _, t := range resp.InstanceTypes {
			infos[aws.StringValue(t.InstanceType)] = nodebootstrap.NewInstanceTypeInfo(t)
		}
		if resp.NextToken == nil {
			return nil
		}
		token = resp.NextToken
	}
}

func newSession(region string) *session.Session {
	config := aws.NewConfig().WithRegion(region)
	opts := session.Options{Config: *config}
	return session.Must(session.NewSessionWithOptions(opts))
}

func render(infos map[string]nodebootstrap.InstanceTypeInfo) error {
	f := NewFile("nodebootstrap")

	f.Commentf("This file was generated %s by reserved_generate.go; DO NOT EDIT.", time.Now().Format(time.RFC3339))
	f.Line()
	f.Comment("Data downloaded through the API.")

	f.Var().Id("instanceTypeInfos").Op("=").
		Map(String()).Id("InstanceTypeInfo").Values(DictFunc(func(d Dict) {
		for k, info := range infos {
			d[Lit(k)] = Values(Dict{
				Id("Storage"): Lit(info.Storage),
				Id("Memory"):  Lit(info.Memory),
				Id("CPU"):     Lit(info.CPU),
			})
		}
	}))

	return f.Save(outputFilename)
}
