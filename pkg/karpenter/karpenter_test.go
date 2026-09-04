package karpenter

import (
	"context"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"helm.sh/helm/v3/pkg/registry"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/karpenter/providers"
	"github.com/weaveworks/eksctl/pkg/karpenter/providers/fakes"
)

var _ = Describe("Install", func() {

	Context("Install", func() {

		var (
			fakeHelmInstaller  *fakes.FakeHelmInstaller
			installerUnderTest *Installer
			cfg                *api.ClusterConfig
		)

		BeforeEach(func() {
			cfg = api.NewClusterConfig()
			cfg.Metadata.Name = "test-cluster"
			cfg.Karpenter = &api.Karpenter{
				Version:                "0.15.3",
				CreateServiceAccount:   api.Disabled(),
				DefaultInstanceProfile: nil,
				// The queue name is only sent when eksctl provisioned the queue,
				// so the specs that assert on it enable it explicitly.
				WithSpotInterruptionQueue: api.Enabled(),
			}
			cfg.Status = &api.ClusterStatus{
				Endpoint: "https://endpoint.com",
			}
			fakeHelmInstaller = &fakes.FakeHelmInstaller{}
			installerUnderTest = &Installer{
				Options: Options{
					HelmInstaller: fakeHelmInstaller,
					Namespace:     "karpenter",
					ClusterConfig: cfg,
				},
			}
		})

		It("installs karpenter into an existing cluster", func() {
			Expect(installerUnderTest.Install(context.Background(), "role-arn", "role/profile")).To(Succeed())
			_, args := fakeHelmInstaller.InstallChartArgsForCall(0)
			args.RegistryClient = &registry.Client{}
			values := map[string]interface{}{
				clusterName:     cfg.Metadata.Name,
				clusterEndpoint: cfg.Status.Endpoint,
				serviceAccount: map[string]interface{}{
					create: api.IsEnabled(cfg.Karpenter.CreateServiceAccount),
					serviceAccountAnnotation: map[string]interface{}{
						api.AnnotationEKSRoleARN: "role-arn",
					},
					serviceAccountName: DefaultServiceAccountName,
				},
				aws: map[string]interface{}{
					defaultInstanceProfile: "role/profile",
				},
				settings: map[string]interface{}{
					aws: map[string]interface{}{
						defaultInstanceProfile: "role/profile",
						clusterName:            cfg.Metadata.Name,
						clusterEndpoint:        cfg.Status.Endpoint,
						interruptionQueueName:  cfg.Metadata.Name,
					},
				},
			}
			Expect(args).To(Equal(providers.InstallChartOpts{
				ChartName:       "oci://public.ecr.aws/karpenter/karpenter",
				CreateNamespace: true,
				Namespace:       "karpenter",
				ReleaseName:     "karpenter",
				Values:          values,
				Version:         "0.15.3",
				RegistryClient:  &registry.Client{},
			}))
		})

		It("installs karpenter with expanded settings.aws values for version greater or equal to v0.33.0", func() {
			installerUnderTest.ClusterConfig.Karpenter.Version = "0.33.0"
			Expect(installerUnderTest.Install(context.Background(), "dummy", "dummy")).To(Succeed())
			_, opts := fakeHelmInstaller.InstallChartArgsForCall(0)
			values := map[string]interface{}{
				settings: map[string]interface{}{
					defaultInstanceProfile: "dummy",
					clusterName:            cfg.Metadata.Name,
					clusterEndpoint:        cfg.Status.Endpoint,
					// The flattened Karpenter chart names this value
					// "interruptionQueue", not "interruptionQueueName" --
					// see charts/karpenter/values.yaml from v0.32.0 onwards.
					// Asserted as a literal rather than via a constant so the
					// test pins the key the chart actually reads.
					"interruptionQueue": cfg.Metadata.Name,
				},
			}
			Expect(opts.Values[settings]).To(Equal(values[settings]))
			// The legacy specs assert the whole values map; do the same here so
			// the top-level keys are guarded on the flattened path too.
			Expect(opts.Values[aws]).To(Equal(map[string]interface{}{defaultInstanceProfile: "dummy"}))
			Expect(opts.Values[clusterName]).To(Equal(cfg.Metadata.Name))
			Expect(opts.Values[clusterEndpoint]).To(Equal(cfg.Status.Endpoint))
		})

		When("withSpotInterruptionQueue is disabled", func() {

			BeforeEach(func() {
				cfg.Karpenter.WithSpotInterruptionQueue = api.Disabled()
			})

			// pkg/cfn/builder only creates the SQS queue, and only grants the
			// controller role sqs:ReceiveMessage on it, when the queue is
			// enabled. Advertising a queue name in either chart layout would
			// point Karpenter at a queue that does not exist and that it
			// cannot poll.
			It("omits the queue name from the legacy settings.aws values", func() {
				Expect(installerUnderTest.Install(context.Background(), "dummy", "dummy")).To(Succeed())
				_, opts := fakeHelmInstaller.InstallChartArgsForCall(0)
				Expect(opts.Values[settings]).To(Equal(map[string]interface{}{
					aws: map[string]interface{}{
						defaultInstanceProfile: "dummy",
						clusterName:            cfg.Metadata.Name,
						clusterEndpoint:        cfg.Status.Endpoint,
					},
				}))
			})

			It("omits the queue name from the flattened settings values", func() {
				installerUnderTest.ClusterConfig.Karpenter.Version = "0.33.0"
				Expect(installerUnderTest.Install(context.Background(), "dummy", "dummy")).To(Succeed())
				_, opts := fakeHelmInstaller.InstallChartArgsForCall(0)
				Expect(opts.Values[settings]).To(Equal(map[string]interface{}{
					defaultInstanceProfile: "dummy",
					clusterName:            cfg.Metadata.Name,
					clusterEndpoint:        cfg.Status.Endpoint,
				}))
			})
		})

		When("install chart fails", func() {

			BeforeEach(func() {
				fakeHelmInstaller.AddRepoReturns(nil)
				fakeHelmInstaller.InstallChartReturns(errors.New("nope"))
			})

			It("errors", func() {
				Expect(installerUnderTest.Install(context.Background(), "", "role/profile")).
					To(MatchError(ContainSubstring("failed to install Karpenter chart: nope")))
			})
		})

		When("service account is defined", func() {
			It("add role to the values for the helm chart", func() {
				Expect(installerUnderTest.Install(context.Background(), "role/account", "role/profile")).To(Succeed())
				_, opts := fakeHelmInstaller.InstallChartArgsForCall(0)
				values := map[string]interface{}{
					clusterName:     cfg.Metadata.Name,
					clusterEndpoint: cfg.Status.Endpoint,
					serviceAccount: map[string]interface{}{
						create: api.IsEnabled(cfg.Karpenter.CreateServiceAccount),
						serviceAccountAnnotation: map[string]interface{}{
							api.AnnotationEKSRoleARN: "role/account",
						},
						serviceAccountName: DefaultServiceAccountName,
					},
					aws: map[string]interface{}{
						defaultInstanceProfile: "role/profile",
					},
					settings: map[string]interface{}{
						aws: map[string]interface{}{
							defaultInstanceProfile: "role/profile",
							clusterName:            cfg.Metadata.Name,
							clusterEndpoint:        cfg.Status.Endpoint,
							interruptionQueueName:  cfg.Metadata.Name,
						},
					},
				}
				Expect(opts.Values).To(Equal(values))
			})
		})
	})
})
