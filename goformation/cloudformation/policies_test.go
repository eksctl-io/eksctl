package cloudformation_test

import (
	"goformation/v4/cloudformation/types"

	"github.com/sanathkr/yaml"

	"goformation/v4/cloudformation"
	"goformation/v4/cloudformation/autoscaling"
	"goformation/v4/cloudformation/policies"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Goformation", func() {

	Context("with a template that contains an UpdatePolicy", func() {

		tests := []struct {
			Name     string
			Input    *policies.UpdatePolicy
			Expected interface{}
		}{
			{
				Name: "AutoScalingReplacingUpdate",
				Input: &policies.UpdatePolicy{
					AutoScalingReplacingUpdate: &policies.AutoScalingReplacingUpdate{
						WillReplace: true,
					},
				},
				Expected: map[string]interface{}{
					"AutoScalingReplacingUpdate": map[string]interface{}{
						"WillReplace": true,
					},
				},
			},
			{
				Name: "AutoScalingReplacingUpdate",
				Input: &policies.UpdatePolicy{
					AutoScalingRollingUpdate: &policies.AutoScalingRollingUpdate{
						MaxBatchSize:                  types.NewInteger(10),
						MinInstancesInService:         types.NewInteger(11),
						MinSuccessfulInstancesPercent: float64(12),
						PauseTime:                     "test-pause-time",
						SuspendProcesses:              []string{"test-suspend1", "test-suspend2"},
						WaitOnResourceSignals:         true,
					},
				},
				Expected: map[string]interface{}{
					"AutoScalingRollingUpdate": map[string]interface{}{
						"MaxBatchSize":                  float64(10),
						"MinInstancesInService":         float64(11),
						"MinSuccessfulInstancesPercent": float64(12),
						"PauseTime":                     "test-pause-time",
						"SuspendProcesses":              []interface{}{"test-suspend1", "test-suspend2"},
						"WaitOnResourceSignals":         true,
					},
				},
			},
			{
				Name: "AutoScalingScheduledAction",
				Input: &policies.UpdatePolicy{
					AutoScalingScheduledAction: &policies.AutoScalingScheduledAction{
						IgnoreUnmodifiedGroupSizeProperties: true,
					},
				},
				Expected: map[string]interface{}{
					"AutoScalingScheduledAction": map[string]interface{}{
						"IgnoreUnmodifiedGroupSizeProperties": true,
					},
				},
			},
			{
				Name: "CodeDeployLambdaAliasUpdate",
				Input: &policies.UpdatePolicy{
					CodeDeployLambdaAliasUpdate: &policies.CodeDeployLambdaAliasUpdate{
						ApplicationName:        "test-application-name",
						DeploymentGroupName:    "test-deployment-group-name",
						AfterAllowTrafficHook:  "test-after-allow-traffic-hook",
						BeforeAllowTrafficHook: "test-before-allow-traffic-hook",
					},
				},
				Expected: map[string]interface{}{
					"CodeDeployLambdaAliasUpdate": map[string]interface{}{
						"ApplicationName":        "test-application-name",
						"DeploymentGroupName":    "test-deployment-group-name",
						"AfterAllowTrafficHook":  "test-after-allow-traffic-hook",
						"BeforeAllowTrafficHook": "test-before-allow-traffic-hook",
					},
				},
			},
		}

		for _, test := range tests {
			test := test

			It("should have the correct values for "+test.Name, func() {

				asg := autoscaling.AutoScalingGroup{}
				asg.AWSCloudFormationUpdatePolicy = test.Input

				template := &cloudformation.Template{
					Resources: cloudformation.Resources{"AutoScalingGroup": &asg},
				}

				data, err := template.JSON()
				Expect(err).To(BeNil())

				var result map[string]interface{}
				if err := yaml.Unmarshal(data, &result); err != nil {
					Fail(err.Error())
				}

				resources, ok := result["Resources"].(map[string]interface{})
				Expect(ok).To(BeTrue())

				bucket, ok := resources["AutoScalingGroup"].(map[string]interface{})
				Expect(ok).To(BeTrue())

				policy, ok := bucket["UpdatePolicy"].(map[string]interface{})
				Expect(ok).To(BeTrue())
				Expect(policy).To(BeEquivalentTo(test.Expected))

			})

		}

	})

	Context("with a template that contains an CreationPolicy", func() {

		tests := []struct {
			Name     string
			Input    *policies.CreationPolicy
			Expected interface{}
		}{
			{
				Name: "CreationPolicy",
				Input: &policies.CreationPolicy{
					AutoScalingCreationPolicy: &policies.AutoScalingCreationPolicy{
						MinSuccessfulInstancesPercent: 10,
					},
					ResourceSignal: &policies.ResourceSignal{
						Count:   11,
						Timeout: "test-timeout",
					},
				},
				Expected: map[string]interface{}{
					"AutoScalingCreationPolicy": map[string]interface{}{
						"MinSuccessfulInstancesPercent": float64(10),
					},
					"ResourceSignal": map[string]interface{}{
						"Count":   float64(11),
						"Timeout": "test-timeout",
					},
				},
			},
		}

		for _, test := range tests {
			test := test

			It("should have the correct values for "+test.Name, func() {

				asg := autoscaling.AutoScalingGroup{}
				asg.AWSCloudFormationCreationPolicy = test.Input

				template := &cloudformation.Template{
					Resources: cloudformation.Resources{"AutoScalingGroup": &asg},
				}

				data, err := template.JSON()
				Expect(err).To(BeNil())

				var result map[string]interface{}
				if err := yaml.Unmarshal(data, &result); err != nil {
					Fail(err.Error())
				}

				resources, ok := result["Resources"].(map[string]interface{})
				Expect(ok).To(BeTrue())

				bucket, ok := resources["AutoScalingGroup"].(map[string]interface{})
				Expect(ok).To(BeTrue())

				policy, ok := bucket["CreationPolicy"].(map[string]interface{})
				Expect(ok).To(BeTrue())
				Expect(policy).To(BeEquivalentTo(test.Expected))

			})

		}

	})

	Context("with a template that contains an DeletionPolicy", func() {

		tests := []struct {
			Name     string
			Input    policies.DeletionPolicy
			Expected interface{}
		}{
			{
				Name:     "DeletionPolicy-Delete",
				Input:    "Delete",
				Expected: "Delete",
			},
			{
				Name:     "DeletionPolicy-Retain",
				Input:    "Retain",
				Expected: "Retain",
			},
			{
				Name:     "DeletionPolicy-Snapshot",
				Input:    "Snapshot",
				Expected: "Snapshot",
			},
		}

		for _, test := range tests {
			test := test

			It("should have the correct values for "+test.Name, func() {

				asg := autoscaling.AutoScalingGroup{}
				asg.AWSCloudFormationDeletionPolicy = test.Input

				template := &cloudformation.Template{
					Resources: cloudformation.Resources{"AutoScalingGroup": &asg},
				}

				data, err := template.JSON()
				Expect(err).To(BeNil())

				var result map[string]interface{}
				if err := yaml.Unmarshal(data, &result); err != nil {
					Fail(err.Error())
				}

				resources, ok := result["Resources"].(map[string]interface{})
				Expect(ok).To(BeTrue())

				bucket, ok := resources["AutoScalingGroup"].(map[string]interface{})
				Expect(ok).To(BeTrue())

				policy, ok := bucket["DeletionPolicy"].(string)
				Expect(ok).To(BeTrue())
				Expect(policy).To(Equal(test.Expected))

			})

		}

	})

})
