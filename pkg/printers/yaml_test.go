package printers_test

import (
	"bufio"
	"bytes"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	. "github.com/weaveworks/eksctl/pkg/printers"
)

var _ = Describe("YAML Printer", func() {

	Describe("When creating New YAML printer", func() {
		var (
			printer OutputPrinter
		)

		BeforeEach(func() {
			printer = NewYAMLPrinter()
		})

		It("should not be nil", func() {
			Expect(printer).ShouldNot(BeNil())
		})

		It("should be the correct type", func() {
			_ = printer.(*YAMLPrinter)
		})

		Context("given a cluster struct and calling PrintObj", func() {
			var (
				cluster     ekstypes.Cluster
				err         error
				actualBytes bytes.Buffer
			)

			BeforeEach(func() {
				created := &time.Time{}
				cluster = ekstypes.Cluster{
					Name:      aws.String("test-cluster"),
					Status:    ekstypes.ClusterStatusActive,
					Arn:       aws.String("arn-12345678"),
					CreatedAt: created,
					ResourcesVpcConfig: &ekstypes.VpcConfigResponse{
						VpcId:     aws.String("vpc-1234"),
						SubnetIds: []string{"sub1", "sub2"},
					},
				}
			})

			JustBeforeEach(func() {
				w := bufio.NewWriter(&actualBytes)
				err = printer.PrintObjWithKind("clusters", []ekstypes.Cluster{cluster}, w)
				w.Flush()
			})

			AfterEach(func() {
				actualBytes.Reset()
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("the output should equal the golden file yamltest_single.golden", func() {
				g, err := os.ReadFile("testdata/yamltest_single.golden")
				if err != nil {
					GinkgoT().Fatalf("failed reading .golden: %s", err)
				}

				Expect(actualBytes.Bytes()).To(MatchYAML(g))
			})
		})

		Context("given 2 cluster structs and calling PrintObj", func() {
			var (
				clusters    []ekstypes.Cluster
				err         error
				actualBytes bytes.Buffer
			)

			BeforeEach(func() {
				created := &time.Time{}
				clusters = []ekstypes.Cluster{
					{
						Name:      aws.String("test-cluster-1"),
						Status:    ekstypes.ClusterStatusActive,
						Arn:       aws.String("arn-12345678"),
						CreatedAt: created,
						ResourcesVpcConfig: &ekstypes.VpcConfigResponse{
							VpcId:     aws.String("vpc-1234"),
							SubnetIds: []string{"sub1", "sub2"},
						},
					},
					{
						Name:      aws.String("test-cluster-2"),
						Status:    ekstypes.ClusterStatusActive,
						Arn:       aws.String("arn-87654321"),
						CreatedAt: created,
						ResourcesVpcConfig: &ekstypes.VpcConfigResponse{
							VpcId:     aws.String("vpc-1234"),
							SubnetIds: []string{"sub1", "sub2"},
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

			It("the output should equal the golden file yamltest_2clusters.golden", func() {
				g, err := os.ReadFile("testdata/yamltest_2clusters.golden")
				if err != nil {
					GinkgoT().Fatalf("failed reading .golden: %s", err)
				}

				Expect(actualBytes.Bytes()).To(MatchYAML(g))
			})
		})

		Context("given an empty cluster list and calling PrintObj", func() {
			var (
				clusters    []ekstypes.Cluster
				err         error
				actualBytes bytes.Buffer
			)

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

			It("the output should be an empty array", func() {
				g := "[]"

				Expect(actualBytes.Bytes()).To(MatchYAML(g))
			})
		})
	})
})
