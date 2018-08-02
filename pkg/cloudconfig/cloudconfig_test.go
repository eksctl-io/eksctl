package cloudconfig_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/weaveworks/eksctl/pkg/cloudconfig"
)

var _ = Describe("cloudconfig", func() {
	var (
		err    error
		result string
		output *CloudConfig
	)

	input := New()
	input.AddPackages("curl", "jq")

	const (
		testScript1     = "curl --silen jsonip.com | jq ."
		testScript1Name = "foo.sh"
		testScript1Path = "/var/lib/cloud/scripts/per-instance/" + testScript1Name
	)

	input.AddShellCommand(testScript1)

	input.RunScript(testScript1Name, testScript1)

	It("encode cloud-config without errors, with a non-blank result", func() {
		result, err = input.Encode()
		Expect(err).NotTo(HaveOccurred())
		Expect(result).NotTo(Equal(""))
		Expect(result).NotTo(Equal("H4sIAAAAAAAA")) // empty base64+gzip string
	})

	It("decode config without an error", func() {
		output, err = DecodeCloudConfig(result)
		Expect(err).NotTo(HaveOccurred())
	})

	It("has the intended packages", func() {
		Expect(output.Packages).To(Equal(input.Packages))
	})

	It("has has correct shell command", func() {
		Expect(output.Commands[0].([]interface{})[0]).To(Equal(Shell))
		Expect(output.Commands[0].([]interface{})[1]).To(Equal("-c"))
		Expect(output.Commands[0].([]interface{})[2]).To(Equal(testScript1))
	})

	It("has has script file and calls it", func() {
		cmd := output.Commands[1].([]interface{})
		file := output.WriteFiles[0]
		Expect(len(cmd)).To(Equal(1))
		Expect(cmd[0]).To(Equal(testScript1Path))
		Expect(file.Path).To(Equal(testScript1Path))
		Expect(file.Content).To(Equal(testScript1))
		Expect(file.Permissions).To(Equal("0755"))
	})
})
