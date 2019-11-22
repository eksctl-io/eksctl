package fargate_test

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/fargate"
	"github.com/weaveworks/eksctl/pkg/printers"
)

var _ = Describe("fargate", func() {
	Describe("PrintProfiles", func() {
		It("formats profiles & prints them as a table", func() {
			profiles := sampleProfiles()
			out := bytes.NewBufferString("")
			err := fargate.PrintProfiles(profiles, out, printers.TableType)
			Expect(err).To(Not(HaveOccurred()))
			Expect(out.String()).To(Equal(expectedTable))
		})

		It("formats profiles & prints them as a YAML object", func() {
			profiles := sampleProfiles()
			out := bytes.NewBufferString("")
			err := fargate.PrintProfiles(profiles, out, printers.YAMLType)
			Expect(err).To(Not(HaveOccurred()))
			Expect(out.String()).To(Equal(expectedYAML))
		})

		It("formats profiles & prints them as a JSON object", func() {
			profiles := sampleProfiles()
			out := bytes.NewBufferString("")
			err := fargate.PrintProfiles(profiles, out, printers.JSONType)
			Expect(err).To(Not(HaveOccurred()))
			Expect(out.String()).To(Equal(expectedJSON))
		})

		It("returns an error for unsupported printer type", func() {
			profiles := sampleProfiles()
			out := bytes.NewBufferString("")
			err := fargate.PrintProfiles(profiles, out, printers.Type("foo"))
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown output printer type: expected {\"yaml\",\"json\",\"table\"} but got \"foo\""))
		})
	})
})

const expectedTable = `NAME	POD_EXECUTION_ROLE_ARN		SUBNETS	SELECTOR_NAMESPACE	SELECTOR_LABELS
default	arn:aws:iam::123:role/root	<none>	default			<none>
default	arn:aws:iam::123:role/root	<none>	kube-system		app=my-app,env=test
prod	arn:aws:iam::123:role/root	prod	prod			env=prod
`

const expectedYAML = `- name: default
  podExecutionRoleARN: arn:aws:iam::123:role/root
  selectors:
  - labels:
      app: my-app
      env: test
    namespace: kube-system
  - namespace: default
- name: prod
  podExecutionRoleARN: arn:aws:iam::123:role/root
  selectors:
  - labels:
      env: prod
    namespace: prod
  subnets:
  - prod
`

const expectedJSON = `[
    {
        "name": "default",
        "podExecutionRoleARN": "arn:aws:iam::123:role/root",
        "selectors": [
            {
                "namespace": "kube-system",
                "labels": {
                    "app": "my-app",
                    "env": "test"
                }
            },
            {
                "namespace": "default"
            }
        ]
    },
    {
        "name": "prod",
        "podExecutionRoleARN": "arn:aws:iam::123:role/root",
        "selectors": [
            {
                "namespace": "prod",
                "labels": {
                    "env": "prod"
                }
            }
        ],
        "subnets": [
            "prod"
        ]
    }
]`

func sampleProfiles() []*api.FargateProfile {
	return []*api.FargateProfile{
		&api.FargateProfile{
			Name:                "default",
			PodExecutionRoleARN: "arn:aws:iam::123:role/root",
			Selectors: []api.FargateProfileSelector{
				api.FargateProfileSelector{
					Namespace: "kube-system",
					Labels: map[string]string{
						"app": "my-app",
						"env": "test",
					},
				},
				api.FargateProfileSelector{
					Namespace: "default",
				},
			},
		},
		&api.FargateProfile{
			Name:    "prod",
			Subnets: []string{"prod"},
			Selectors: []api.FargateProfileSelector{
				api.FargateProfileSelector{
					Namespace: "prod",
					Labels: map[string]string{
						"env": "prod",
					},
				},
			},
			PodExecutionRoleARN: "arn:aws:iam::123:role/root",
		},
	}
}
