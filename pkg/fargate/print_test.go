package fargate_test

import (
	"bytes"

	. "github.com/onsi/ginkgo/v2"
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
			err := fargate.PrintProfiles(profiles, out, "foo")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("unknown output printer type: expected {\"yaml\",\"json\",\"table\"} but got \"foo\""))
		})
	})
})

const expectedTable = `NAME	SELECTOR_NAMESPACE	SELECTOR_LABELS		POD_EXECUTION_ROLE_ARN		SUBNETS				TAGS			STATUS
fp-prod	prod			env=prod		arn:aws:iam::123:role/root	subnet-prod,subnet-d34dc0w	<none>			ACTIVE
fp-test	default			<none>			arn:aws:iam::123:role/root	<none>				app=my-app,env=test	ACTIVE
fp-test	kube-system		app=my-app,env=test	arn:aws:iam::123:role/root	<none>				app=my-app,env=test	ACTIVE
`

const expectedYAML = `- name: fp-test
  podExecutionRoleARN: arn:aws:iam::123:role/root
  selectors:
  - labels:
      app: my-app
      env: test
    namespace: kube-system
  - namespace: default
  status: ACTIVE
  tags:
    app: my-app
    env: test
- name: fp-prod
  podExecutionRoleARN: arn:aws:iam::123:role/root
  selectors:
  - labels:
      env: prod
    namespace: prod
  status: ACTIVE
  subnets:
  - subnet-prod
  - subnet-d34dc0w
`

const expectedJSON = `[
    {
        "name": "fp-test",
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
        ],
        "tags": {
            "app": "my-app",
            "env": "test"
        },
        "status": "ACTIVE"
    },
    {
        "name": "fp-prod",
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
            "subnet-prod",
            "subnet-d34dc0w"
        ],
        "status": "ACTIVE"
    }
]`

func sampleProfiles() []*api.FargateProfile {
	return []*api.FargateProfile{
		{
			Name:                "fp-test",
			PodExecutionRoleARN: "arn:aws:iam::123:role/root",
			Selectors: []api.FargateProfileSelector{
				{
					Namespace: "kube-system",
					Labels: map[string]string{
						"app": "my-app",
						"env": "test",
					},
				},
				{
					Namespace: "default",
				},
			},
			Tags: map[string]string{
				"app": "my-app",
				"env": "test",
			},
			Status: "ACTIVE",
		},
		{
			Name:    "fp-prod",
			Subnets: []string{"subnet-prod", "subnet-d34dc0w"},
			Selectors: []api.FargateProfileSelector{
				{
					Namespace: "prod",
					Labels: map[string]string{
						"env": "prod",
					},
				},
			},
			PodExecutionRoleARN: "arn:aws:iam::123:role/root",
			Status:              "ACTIVE",
		},
	}
}
