package helm

import (
	"context"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/otiai10/copy"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/getter"
	kubefake "helm.sh/helm/v3/pkg/kube/fake"
	"helm.sh/helm/v3/pkg/registry"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"

	"github.com/weaveworks/eksctl/pkg/karpenter/providers"
	"github.com/weaveworks/eksctl/pkg/karpenter/providers/fakes"
)

var _ = Describe("HelmInstaller", func() {

	Context("InstallChart", func() {

		var (
			fakeGetter         *fakes.FakeURLGetter
			getters            getter.Providers
			tmp                string
			err                error
			installerUnderTest *Installer
			values             map[string]interface{}
			actionConfig       *action.Configuration
			fakeKubeClient     *fakes.PrintingKubeClient
			registryClient     *registry.Client
		)

		BeforeEach(func() {
			tmp, err = os.MkdirTemp("", "helm-testing")
			Expect(err).NotTo(HaveOccurred())
			fakeGetter = &fakes.FakeURLGetter{}
			provider := getter.Provider{
				Schemes: []string{"http", "https", registry.OCIScheme},
				New: func(options ...getter.Option) (getter.Getter, error) {
					return fakeGetter, nil
				},
			}
			getters = append(getters, provider)
			store := storage.Init(driver.NewMemory())
			fakeKubeClient = &fakes.PrintingKubeClient{Out: ioutil.Discard}
			actionConfig = &action.Configuration{
				Releases:     store,
				KubeClient:   fakeKubeClient,
				Capabilities: chartutil.DefaultCapabilities,
				Log:          func(format string, v ...interface{}) {},
			}
			installerUnderTest = &Installer{
				Getters: getters,
				Settings: &cli.EnvSettings{
					RepositoryCache:  tmp,
					RepositoryConfig: filepath.Join(tmp, "repositories.yaml"),
					Debug:            true,
				},
				ActionConfig: actionConfig,
			}
			values = map[string]interface{}{
				"some": "value",
			}
			registryClient, _ = registry.NewClient(
				registry.ClientOptEnableCache(true),
			)

		})

		AfterEach(func() {
			_ = os.RemoveAll(tmp)
		})

		It("can install a test chart", func() {
			// write out repo config
			Expect(os.WriteFile(filepath.Join(tmp, "repositories.yaml"), []byte(expectedRepositoryYaml), 0644)).To(Succeed())
			Expect(copy.Copy(filepath.Join("testdata", "karpenter-0.18.0.tgz"), filepath.Join(tmp, "karpenter-0.18.0.tgz"))).To(Succeed())
			Expect(copy.Copy(filepath.Join("testdata", "karpenter-index.yaml"), filepath.Join(tmp, "karpenter-index.yaml"))).To(Succeed())
			Expect(installerUnderTest.InstallChart(context.Background(), providers.InstallChartOpts{
				ChartName:       "oci://public.ecr.aws/karpenter/karpenter",
				CreateNamespace: true,
				Namespace:       "karpenter",
				ReleaseName:     "karpenter",
				Values:          values,
				Version:         "v0.18.0",
				RegistryClient:  registryClient,
			})).To(Succeed())
			Expect(fakeKubeClient.BuildCall).To(Equal(4))
			Expect(fakeKubeClient.CreateCall).To(Equal(3))
		})
		When("creating a namespace is disabled", func() {
			It("will not call build and create for that resource", func() {
				// write out repo config
				Expect(os.WriteFile(filepath.Join(tmp, "repositories.yaml"), []byte(expectedRepositoryYaml), 0644)).To(Succeed())
				Expect(copy.Copy(filepath.Join("testdata", "karpenter-0.18.0.tgz"), filepath.Join(tmp, "karpenter-0.18.0.tgz"))).To(Succeed())
				Expect(copy.Copy(filepath.Join("testdata", "karpenter-index.yaml"), filepath.Join(tmp, "karpenter-index.yaml"))).To(Succeed())
				Expect(installerUnderTest.InstallChart(context.Background(), providers.InstallChartOpts{
					ChartName:       "oci://public.ecr.aws/karpenter/karpenter",
					CreateNamespace: false,
					Namespace:       "karpenter",
					ReleaseName:     "karpenter",
					Values:          values,
					Version:         "v0.18.0",
					RegistryClient:  registryClient,
				})).To(Succeed())
				// Verifying the call number is easier than trying to mock RestClient calls through the fake kube client.
				// And the Printer does not work, because the Builder returns an empty list that the Creator gleefully
				// accepts and does nothing.
				Expect(fakeKubeClient.BuildCall).To(Equal(3))
				Expect(fakeKubeClient.CreateCall).To(Equal(2))
			})
		})
		When("locate chart is unable to find the requested chart", func() {
			It("errors", func() {
				err := installerUnderTest.InstallChart(context.Background(), providers.InstallChartOpts{
					ChartName:       "oci://public.ecr.aws/karpenter/karpenter",
					CreateNamespace: true,
					Namespace:       "karpenter",
					ReleaseName:     "karpenter",
					Values:          values,
					Version:         "0.8.0",
					RegistryClient:  registryClient,
				})
				Expect(err).To(MatchError(ContainSubstring("failed to locate chart: public.ecr.aws/karpenter/karpenter:0.8.0: not found")))
			})
		})
		When("the exact version is not specified", func() {
			It("errors", func() {
				Expect(os.WriteFile(filepath.Join(tmp, "repositories.yaml"), []byte(expectedRepositoryYaml), 0644)).To(Succeed())
				Expect(copy.Copy(filepath.Join("testdata", "karpenter-0.18.0.tgz"), filepath.Join(tmp, "karpenter-0.18.0.tgz"))).To(Succeed())
				Expect(copy.Copy(filepath.Join("testdata", "karpenter-index.yaml"), filepath.Join(tmp, "karpenter-index.yaml"))).To(Succeed())
				err := installerUnderTest.InstallChart(context.Background(), providers.InstallChartOpts{
					ChartName:       "oci://public.ecr.aws/karpenter/karpenter",
					CreateNamespace: true,
					Namespace:       "karpenter",
					ReleaseName:     "karpenter",
					Values:          values,
					Version:         "0.18.0",
					RegistryClient:  registryClient,
				})
				Expect(err).To(MatchError(ContainSubstring("failed to locate chart: public.ecr.aws/karpenter/karpenter:0.18.0: not found")))
			})
		})
		When("kube client fails to reach the cluster", func() {
			It("errors", func() {
				Expect(os.WriteFile(filepath.Join(tmp, "repositories.yaml"), []byte(expectedRepositoryYaml), 0644)).To(Succeed())
				Expect(copy.Copy(filepath.Join("testdata", "karpenter-0.18.0.tgz"), filepath.Join(tmp, "karpenter-0.18.0.tgz"))).To(Succeed())
				Expect(copy.Copy(filepath.Join("testdata", "karpenter-index.yaml"), filepath.Join(tmp, "karpenter-index.yaml"))).To(Succeed())
				fakeKube := &kubefake.FailingKubeClient{
					PrintingKubeClient: kubefake.PrintingKubeClient{},
					CreateError:        errors.New("nope"),
				}
				actionConfig.KubeClient = fakeKube
				installerUnderTest.ActionConfig = actionConfig
				err := installerUnderTest.InstallChart(context.Background(), providers.InstallChartOpts{
					ChartName:       "oci://public.ecr.aws/karpenter/karpenter",
					CreateNamespace: true,
					Namespace:       "karpenter",
					ReleaseName:     "karpenter",
					Values:          values,
					Version:         "v0.18.0",
					RegistryClient:  registryClient,
				})
				Expect(err).To(MatchError(ContainSubstring("failed to install chart: failed to install CRD crds/karpenter.k8s.aws_awsnodetemplates.yaml: nope")))
			})
		})
	})
})

var expectedRepositoryYaml = `apiVersion: ""
generated: "0001-01-01T00:00:00Z"
repositories:
- caFile: ""
  certFile: ""
  insecure_skip_tls_verify: false
  keyFile: ""
  name: karpenter
  pass_credentials_all: false
  password: ""
  url: oci://public.ecr.aws/karpenter/karpenter
  username: ""
`
