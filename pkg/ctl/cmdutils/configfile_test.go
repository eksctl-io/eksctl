package cmdutils_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/spf13/cobra"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	. "github.com/weaveworks/eksctl/pkg/ctl/cmdutils"
	// "github.com/weaveworks/eksctl/pkg/printers"
)

var _ = Describe("cmdutils configfile", func() {

	newCmd := func() *cobra.Command {
		return &cobra.Command{
			Use: "test",
			Run: func(_ *cobra.Command, _ []string) {},
		}
	}

	const examplesDir = "../../../examples/"

	Context("load configfiles", func() {

		It("should handle name argument", func() {
			cfg := api.NewClusterConfig()

			{
				rc := &ResourceCmd{
					ClusterConfig: cfg,
					NameArg:       "foo-1",
				}

				err := NewMetadataLoader(rc).Load()
				Expect(err).ToNot(HaveOccurred())
				Expect(cfg.Metadata.Name).To(Equal("foo-1"))
			}

			{
				rc := &ResourceCmd{
					ClusterConfig: cfg,
					NameArg:       "foo-2",
				}

				err := NewMetadataLoader(rc).Load()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("--name=foo-1 and argument foo-2 cannot be used at the same time"))
			}

			{
				rc := &ResourceCmd{
					ClusterConfig:     cfg,
					NameArg:           "foo-3",
					Command:           newCmd(),
					ClusterConfigFile: examplesDir + "01-simple-cluster.yaml",
				}

				err := NewMetadataLoader(rc).Load()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(ErrCannotUseWithConfigFile(`name argument "foo-3"`).Error()))

				fs := rc.Command.Flags()

				fs.StringVar(&cfg.Metadata.Name, "name", "", "")
				rc.Command.Flag("name").Changed = true

				Expect(rc.Command.Flag("name").Changed).To(BeTrue())

				err = NewMetadataLoader(rc).Load()

				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(ErrCannotUseWithConfigFile("--name").Error()))
			}
		})

		It("load all of example file", func() {
			examples, err := filepath.Glob(examplesDir + "*.yaml")
			Expect(err).ToNot(HaveOccurred())

			Expect(examples).To(HaveLen(10))
			for _, example := range examples {
				cfg := api.NewClusterConfig()

				rc := &ResourceCmd{
					Command:           newCmd(),
					ClusterConfigFile: example,
					ClusterConfig:     cfg,
					ProviderConfig:    &api.ProviderConfig{},
				}

				err := NewMetadataLoader(rc).Load()

				Expect(err).ToNot(HaveOccurred())
				Expect(cfg.Metadata.Name).ToNot(BeEmpty())
				Expect(cfg.Metadata.Region).ToNot(BeEmpty())
				Expect(cfg.Metadata.Region).To(Equal(rc.ProviderConfig.Region))
				Expect(cfg.Metadata.Version).To(BeEmpty())
			}
		})
	})
})
