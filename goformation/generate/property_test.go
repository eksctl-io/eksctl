package main_test

import (
	"encoding/json"

	"goformation/v4/cloudformation/serverless"
	"goformation/v4/cloudformation/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Goformation Code Generator", func() {

	Context("with a single generated custom property", func() {

		Context("specified as a Go struct", func() {

			property := &serverless.Function_S3Location{
				Bucket:  types.NewString("test-bucket"),
				Key:     types.NewString("test-key"),
				Version: types.NewInteger(123),
			}
			expected := []byte(`{"Bucket":"test-bucket","Key":"test-key","Version":123}`)

			result, err := json.Marshal(property)
			It("should marshal to JSON successfully", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(Equal(expected))
			})

		})
		Context("specified as JSON", func() {

			property := []byte(`{"Bucket":"test-bucket","Key":"test-key","Version":123}`)
			expected := serverless.Function_S3Location{
				Bucket:  types.NewString("test-bucket"),
				Key:     types.NewString("test-key"),
				Version: types.NewInteger(123),
			}

			result := serverless.Function_S3Location{}
			err := json.Unmarshal(property, &result)
			It("should unmarshal to a Go struct successfully", func() {
				Expect(err).ToNot(HaveOccurred())
				Expect(result).To(Equal(expected))
			})

		})

	})

	Context("with a polymorphic property", func() {

		Context("with multiple types", func() {

			Context("properly marshals and unmarshals values", func() {

				property := []byte(`{"Properties":{"BatchSize":10,"StartingPosition":"LATEST","Stream":"arn"},"Type":"Kinesis"}`)

				result := &serverless.Function_EventSource{}
				err := json.Unmarshal(property, result)
				output, err2 := json.Marshal(result)

				It("should marshal and unmarhal to same value", func() {
					Expect(err).To(BeNil())
					Expect(err2).To(BeNil())
					Expect(output).To(Equal(property))
				})

			})

			Context("properly marshals and unmarshals polymorphic values", func() {

				property := []byte(`{"Properties":{"Bucket":"asd","Events":"LATEST"},"Type":"S3"}`)

				result := &serverless.Function_EventSource{}
				err := json.Unmarshal(property, result)
				output, err2 := json.Marshal(result)

				It("should marshal and unmarhal to same value", func() {
					Expect(err).To(BeNil())
					Expect(err2).To(BeNil())
					Expect(output).To(Equal(property))
				})

			})

			Context("properly Marshals best value", func() {
				expected := []byte(`{"BatchSize":10,"Stream":"arn"}`)

				result := &serverless.Function_Properties{
					SQSEvent:     &serverless.Function_SQSEvent{BatchSize: types.NewInteger(10)},
					KinesisEvent: &serverless.Function_KinesisEvent{BatchSize: types.NewInteger(10), Stream: types.NewString("arn")},
				}

				output, err := result.MarshalJSON()

				It("should marshal and unmarhal to same value", func() {
					Expect(err).To(BeNil())
					Expect(output).To(Equal(expected))
				})

			})

		})

		Context("with a primitive value", func() {

			Context("specified as a Go struct", func() {

				value := types.NewString("test-primitive-value")
				property := &serverless.Function_CodeUri{
					String: &value,
				}

				expected := []byte(`"test-primitive-value"`)
				result, err := json.Marshal(property)
				It("should marshal to JSON successfully", func() {
					Expect(result).To(Equal(expected))
					Expect(err).To(BeNil())
				})

			})

			Context("specified as JSON", func() {

				property := []byte(`"test-primitive-value"`)
				value := types.NewString("test-primitive-value")
				expected := &serverless.Function_CodeUri{
					String: &value,
				}

				result := &serverless.Function_CodeUri{}
				err := json.Unmarshal(property, result)
				It("should unmarshal to a Go struct successfully", func() {
					Expect(result).To(Equal(expected))
					Expect(err).To(BeNil())
				})

			})

		})

		Context("with a custom type value", func() {

			Context("specified as a Go struct", func() {

				property := &serverless.Function_CodeUri{
					S3Location: &serverless.Function_S3Location{
						Bucket:  types.NewString("test-bucket"),
						Key:     types.NewString("test-key"),
						Version: types.NewInteger(123),
					},
				}

				expected := []byte(`{"Bucket":"test-bucket","Key":"test-key","Version":123}`)

				result, err := json.Marshal(property)
				It("should marshal to JSON successfully", func() {
					Expect(result).To(Equal(expected))
					Expect(err).To(BeNil())
				})

			})

			Context("specified as JSON", func() {

				property := []byte(`{"Bucket":"test-bucket","Key":"test-key","Version":123}`)

				expected := &serverless.Function_CodeUri{
					S3Location: &serverless.Function_S3Location{
						Bucket:  types.NewString("test-bucket"),
						Key:     types.NewString("test-key"),
						Version: types.NewInteger(123),
					},
				}

				result := &serverless.Function_CodeUri{}
				err := json.Unmarshal(property, result)
				It("should unmarshal to a Go struct successfully", func() {
					Expect(result).To(Equal(expected))
					Expect(err).To(BeNil())
				})

			})

		})

	})

})
