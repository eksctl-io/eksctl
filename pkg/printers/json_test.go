package printers_test

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"time"

	. "github.com/weaveworks/eksctl/pkg/printers"

	"github.com/aws/aws-sdk-go/aws"
	awseks "github.com/aws/aws-sdk-go/service/eks"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("JSON Printer", func() {

	Describe("When creating New JSON printer", func() {
		var (
			printer OutputPrinter
		)

		BeforeEach(func() {
			printer = NewJSONPrinter()
		})

		It("should not be nil", func() {
			Expect(printer).ShouldNot(BeNil())
		})

		It("should be the correct type", func() {
			_ = printer.(*JSONPrinter)
		})

		Context("given a cluster struct and calling PrintObjWithKind", func() {
			var (
				cluster     *awseks.Cluster
				err         error
				actualBytes bytes.Buffer
			)

			BeforeEach(func() {
				created := &time.Time{}
				cluster = &awseks.Cluster{
					Name:      aws.String("test-cluster"),
					Status:    aws.String(awseks.ClusterStatusActive),
					Arn:       aws.String("arn-12345678"),
					CreatedAt: created,
					ResourcesVpcConfig: &awseks.VpcConfigResponse{
						VpcId:     aws.String("vpc-1234"),
						SubnetIds: []*string{aws.String("sub1"), aws.String("sub2")},
					},
				}
			})

			JustBeforeEach(func() {
				w := bufio.NewWriter(&actualBytes)
				err = printer.PrintObjWithKind("clusters", []*awseks.Cluster{cluster}, w)
				w.Flush()
			})

			AfterEach(func() {
				actualBytes.Reset()
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("the output should equal the golden file jsontest_single.golden", func() {
				g, err := ioutil.ReadFile("testdata/jsontest_single.golden")
				if err != nil {
					GinkgoT().Fatalf("failed reading .golden: %s", err)
				}

				Expect(actualBytes.Bytes()).Should(MatchJSON(g))
			})
		})

		Context("given 2 cluster structs and calling PrintObjWithKind", func() {
			var (
				clusters    []*awseks.Cluster
				err         error
				actualBytes bytes.Buffer
			)

			BeforeEach(func() {
				created := &time.Time{}
				clusters = []*awseks.Cluster{
					{
						Name:      aws.String("test-cluster-1"),
						Status:    aws.String(awseks.ClusterStatusActive),
						Arn:       aws.String("arn-12345678"),
						CreatedAt: created,
						ResourcesVpcConfig: &awseks.VpcConfigResponse{
							VpcId:     aws.String("vpc-1234"),
							SubnetIds: []*string{aws.String("sub1"), aws.String("sub2")},
						},
					},
					{
						Name:      aws.String("test-cluster-2"),
						Status:    aws.String(awseks.ClusterStatusActive),
						Arn:       aws.String("arn-87654321"),
						CreatedAt: created,
						ResourcesVpcConfig: &awseks.VpcConfigResponse{
							VpcId:     aws.String("vpc-1234"),
							SubnetIds: []*string{aws.String("sub1"), aws.String("sub2")},
						},
					},
				}
			})

			JustBeforeEach(func() {
				w := bufio.NewWriter(&actualBytes)
				err = printer.PrintObjWithKind("clusters", clusters, w)
				w.Flush()
			})

			AfterEach(func() {
				actualBytes.Reset()
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("the output should equal the golden file jsontest_2clusters.golden", func() {
				g, err := ioutil.ReadFile("testdata/jsontest_2clusters.golden")
				if err != nil {
					GinkgoT().Fatalf("failed reading .golden: %s", err)
				}

				Expect(actualBytes.Bytes()).Should(MatchJSON(g))
			})
		})
	})
})
