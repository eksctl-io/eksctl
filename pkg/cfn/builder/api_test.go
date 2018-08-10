package builder_test

import (
	"encoding/base64"
	"encoding/json"
	"path/filepath"
	"strings"

	cfn "github.com/aws/aws-sdk-go/service/cloudformation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cloudconfig"
	"github.com/weaveworks/eksctl/pkg/eks/api"
	"github.com/weaveworks/eksctl/pkg/nodebootstrap"
)

type Template struct {
	Resources map[string]struct {
		Properties struct {
			Tags []struct {
				Key   interface{}
				Value interface{}

				PropagateAtLaunch string
			}
			UserData string
		}
	}
}

func newStackWithOutputs(outputs map[string]string) cfn.Stack {
	s := cfn.Stack{}
	for k, v := range outputs {
		func(k, v string) {
			s.Outputs = append(s.Outputs,
				&cfn.Output{
					OutputKey:   &k,
					OutputValue: &v,
				})
		}(k, v)
	}
	return s
}

var _ = Describe("CloudFormation template builder API", func() {

	testAZs := []string{"us-west-2b", "us-west-2a", "us-west-2c"}

	Describe("GetAllOutputsFromClusterStack", func() {
		caCert := "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUN5RENDQWJDZ0F3SUJBZ0lCQURBTkJna3Foa2lHOXcwQkFRc0ZBREFWTVJNd0VRWURWUVFERXdwcmRXSmwKY201bGRHVnpNQjRYRFRFNE1EWXdOekExTlRBMU5Wb1hEVEk0TURZd05EQTFOVEExTlZvd0ZURVRNQkVHQTFVRQpBeE1LYTNWaVpYSnVaWFJsY3pDQ0FTSXdEUVlKS29aSWh2Y05BUUVCQlFBRGdnRVBBRENDQVFvQ2dnRUJBTWJoCnpvZElYR0drckNSZE1jUmVEN0YvMnB1NFZweTdvd3FEVDgrdk9zeGs2bXFMNWxQd3ZicFhmYkE3R0xzMDVHa0wKaDdqL0ZjcU91cnMwUFZSK3N5REtuQXltdDFORWxGNllGQktSV1dUQ1hNd2lwN1pweW9XMXdoYTlJYUlPUGxCTQpPTEVlckRabFVrVDFVV0dWeVdsMmxPeFgxa2JhV2gvakptWWdkeW5jMXhZZ3kxa2JybmVMSkkwLzVUVTRCajJxClB1emtrYW5Xd3lKbGdXQzhBSXlpWW82WFh2UVZmRzYrM3RISE5XM1F1b3ZoRng2MTFOYnl6RUI3QTdtZGNiNmgKR0ZpWjdOeThHZnFzdjJJSmI2Nk9FVzBSdW9oY1k3UDZPdnZmYnlKREhaU2hqTStRWFkxQXN5b3g4Ri9UelhHSgpQUWpoWUZWWEVhZU1wQmJqNmNFQ0F3RUFBYU1qTUNFd0RnWURWUjBQQVFIL0JBUURBZ0trTUE4R0ExVWRFd0VCCi93UUZNQU1CQWY4d0RRWUpLb1pJaHZjTkFRRUxCUUFEZ2dFQkFCa2hKRVd4MHk1LzlMSklWdXJ1c1hZbjN6Z2EKRkZ6V0JsQU44WTlqUHB3S2t0Vy9JNFYyUGg3bWY2Z3ZwZ3Jhc2t1Slk1aHZPcDdBQmcxSTFhaHUxNUFpMUI0ZApuMllRaDlOaHdXM2pKMmhuRXk0VElpb0gza2JFdHRnUVB2bWhUQzNEYUJreEpkbmZJSEJCV1RFTTU1czRwRmxUClpzQVJ3aDc1Q3hYbjdScVU0akpKcWNPaTRjeU5qeFVpRDBqR1FaTmNiZWEyMkRCeTJXaEEzUWZnbGNScGtDVGUKRDVPS3NOWlF4MW9MZFAwci9TSmtPT1NPeUdnbVJURTIrODQxN21PRW02Z3RPMCszdWJkbXQ0aENsWEtFTTZYdwpuQWNlK0JxVUNYblVIN2ZNS3p2TDE5UExvMm5KbFU1TnlCbU1nL1pNVHVlUy80eFZmKy94WnpsQ0Q1WT0KLS0tLS1FTkQgQ0VSVElGSUNBVEUtLS0tLQo="
		caCertData, err := base64.StdEncoding.DecodeString(caCert)
		It("should not error", func() { Expect(err).NotTo(HaveOccurred()) })

		expected := &api.ClusterConfig{
			SecurityGroup:            "sg-0b44c48bcba5b7362",
			Subnets:                  []string{"subnet-0f98135715dfcf55f", "subnet-0ade11bad78dced9e", "subnet-0e2e63ff1712bf6ef"},
			VPC:                      "vpc-0e265ad953062b94b",
			Endpoint:                 "https://DE37D8AFB23F7275D2361AD6B2599143.yl4.us-west-2.eks.amazonaws.com",
			CertificateAuthorityData: caCertData,
			ARN:                      "arn:aws:eks:us-west-2:376248598259:cluster/ferocious-mushroom-1532594698",
			NodeInstanceRoleARN:      "",
			AvailabilityZones:        testAZs,
		}

		initial := &api.ClusterConfig{
			AvailabilityZones: testAZs,
		}

		rs := NewClusterResourceSet(initial)
		rs.AddAllResources()

		sampleStack := newStackWithOutputs(map[string]string{
			"SecurityGroup":            "sg-0b44c48bcba5b7362",
			"Subnets":                  "subnet-0f98135715dfcf55f,subnet-0ade11bad78dced9e,subnet-0e2e63ff1712bf6ef",
			"VPC":                      "vpc-0e265ad953062b94b",
			"Endpoint":                 "https://DE37D8AFB23F7275D2361AD6B2599143.yl4.us-west-2.eks.amazonaws.com",
			"CertificateAuthorityData": caCert,
			"ARN":                      "arn:aws:eks:us-west-2:376248598259:cluster/ferocious-mushroom-1532594698",
			"ClusterStackName":         "",
		})

		It("should not error", func() {
			err := rs.GetAllOutputs(sampleStack)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should be equal", func() {
			Expect(initial).To(Equal(expected))
		})
	})

	Describe("AutoNameTag", func() {
		rs := NewNodeGroupResourceSet(&api.ClusterConfig{
			AvailabilityZones: testAZs,
			NodeType:          "t2.medium",
			Region:            "us-west-2",
		}, "eksctl-test-123-cluster", 0)
		rs.AddAllResources()

		template, err := rs.RenderJSON()
		It("should serialise JSON without errors", func() {
			Expect(err).NotTo(HaveOccurred())
		})
		obj := Template{}
		It("should parse JSON withon errors", func() {
			err := json.Unmarshal(template, &obj)
			Expect(err).NotTo(HaveOccurred())
		})

		It("SG should have correct tags", func() {
			Expect(len(obj.Resources)).ToNot(Equal(0))
			Expect(len(obj.Resources["SG"].Properties.Tags)).To(Equal(2))
			Expect(obj.Resources["SG"].Properties.Tags[0].Key).To(Equal(map[string]interface{}{
				"Fn::Sub": "kubernetes.io/cluster/${ClusterName}",
			}))
			Expect(obj.Resources["SG"].Properties.Tags[0].Value).To(Equal("owned"))
			Expect(obj.Resources["SG"].Properties.Tags[1].Key).To(Equal("Name"))
			Expect(obj.Resources["SG"].Properties.Tags[1].Value).To(Equal(map[string]interface{}{
				"Fn::Sub": "${AWS::StackName}/SG",
			}))
		})
	})

	Describe("NodeGroupTags", func() {
		rs := NewNodeGroupResourceSet(&api.ClusterConfig{
			AvailabilityZones: testAZs,
			NodeType:          "t2.medium",
			Region:            "us-west-2",
		}, "eksctl-test-123-cluster", 0)
		rs.AddAllResources()

		template, err := rs.RenderJSON()
		It("should serialise JSON without errors", func() {
			Expect(err).NotTo(HaveOccurred())
		})
		obj := Template{}
		It("should parse JSON withon errors", func() {
			err := json.Unmarshal(template, &obj)
			Expect(err).NotTo(HaveOccurred())
		})

		It("NodeGroup should have correct tags", func() {
			Expect(len(obj.Resources)).ToNot(Equal(0))
			Expect(len(obj.Resources["NodeGroup"].Properties.Tags)).To(Equal(2))
			Expect(obj.Resources["NodeGroup"].Properties.Tags[0].Key).To(Equal("Name"))
			Expect(obj.Resources["NodeGroup"].Properties.Tags[0].Value).To(Equal(map[string]interface{}{
				"Fn::Sub": "${ClusterName}-${NodeGroupID}-Node",
			}))
			Expect(obj.Resources["NodeGroup"].Properties.Tags[0].PropagateAtLaunch).To(Equal("true"))
			Expect(obj.Resources["NodeGroup"].Properties.Tags[1].Key).To(Equal(map[string]interface{}{
				"Fn::Sub": "kubernetes.io/cluster/${ClusterName}",
			}))
			Expect(obj.Resources["NodeGroup"].Properties.Tags[1].Value).To(Equal("owned"))
			Expect(obj.Resources["NodeGroup"].Properties.Tags[1].PropagateAtLaunch).To(Equal("true"))
		})
	})

	Describe("UserData", func() {
		rs := NewNodeGroupResourceSet(&api.ClusterConfig{
			AvailabilityZones: testAZs,
			NodeType:          "m5.large",
			Region:            "us-west-2",
		}, "eksctl-test-123-cluster", 0)
		rs.AddAllResources()

		var (
			err error
			c   *cloudconfig.CloudConfig
		)

		template, err := rs.RenderJSON()
		It("should serialise JSON without errors", func() {
			Expect(err).NotTo(HaveOccurred())
		})

		It("should parse JSON withon errors and extract valid cloud-config using our implementation", func() {
			obj := Template{}
			err = json.Unmarshal(template, &obj)
			Expect(err).NotTo(HaveOccurred())
			Expect(len(obj.Resources)).ToNot(Equal(0))

			userData := obj.Resources["NodeLaunchConfig"].Properties.UserData
			Expect(userData).ToNot(BeEmpty())

			c, err = cloudconfig.DecodeCloudConfig(userData)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should have packages, scripts and commands in cloud-config", func() {
			Expect(c).ToNot(BeNil())

			Expect(c.Packages).To(Equal([]string{"jq"}))

			cmds := len(c.Commands)
			Expect(c.Commands[0]).To(Equal([]interface{}{"pip", "install", "--upgrade", "awscli"}))
			Expect(c.Commands[cmds-2]).To(Equal([]interface{}{"systemctl", "daemon-reload"}))
			Expect(c.Commands[cmds-1]).To(Equal([]interface{}{"systemctl", "restart", "kubelet"}))

			getFile := func(p string) *cloudconfig.File {
				for _, f := range c.WriteFiles {
					if f.Path == p {
						return &f
					}
				}
				return nil
			}

			checkAsset := func(name, expectedContent string) {
				assetContent, err := nodebootstrap.Asset(name)
				Expect(err).ToNot(HaveOccurred())
				Expect(string(assetContent)).ToNot(BeEmpty())
				Expect(expectedContent).To(Equal(string(assetContent)))
			}

			authenticator := getFile("/etc/eksctl/authenticator.sh")
			Expect(authenticator).ToNot(BeNil())
			Expect(authenticator.Permissions).To(Equal("0755"))
			checkAsset("authenticator.sh", authenticator.Content)

			kubeconfig := getFile("/etc/eksctl/kubeconfig.yaml")
			Expect(kubeconfig).ToNot(BeNil())
			Expect(kubeconfig.Permissions).To(Equal("0644"))
			checkAsset("kubeconfig.yaml", kubeconfig.Content)

			kubeletEnv := getFile("/etc/eksctl/kubelet.env")
			Expect(kubeletEnv).ToNot(BeNil())
			Expect(kubeletEnv.Permissions).To(Equal("0644"))
			Expect(strings.Split(kubeletEnv.Content, "\n")).To(Equal([]string{
				"MAX_PODS=29",
				"CLUSTER_DNS=10.100.0.10",
			}))

			kubeletDropInUnit := getFile("/etc/systemd/system/kubelet.service.d/10-eksclt.al2.conf")
			Expect(kubeletDropInUnit).ToNot(BeNil())
			Expect(kubeletDropInUnit.Permissions).To(Equal("0644"))
			checkAsset("10-eksclt.al2.conf", kubeletDropInUnit.Content)

			checkScript := func(p string, assetContent bool) {
				script := getFile(p)
				Expect(script).ToNot(BeNil())
				Expect(script.Permissions).To(Equal("0755"))
				scriptRuns := false
				for _, s := range c.Commands {
					if s.([]interface{})[0] == script.Path {
						scriptRuns = true
					}
				}
				Expect(scriptRuns).To(BeTrue())
				if assetContent {
					checkAsset(filepath.Base(p), script.Content)
				}
			}

			checkScript("/var/lib/cloud/scripts/per-instance/get_metadata.sh", true)
			checkScript("/var/lib/cloud/scripts/per-instance/get_credentials.sh", true)
		})
	})
})
