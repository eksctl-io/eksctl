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

			err := NewMetadataLoader(nil, cfg, "", "foo-1", nil).Load()
			Expect(err).ToNot(HaveOccurred())
			Expect(cfg.Metadata.Name).To(Equal("foo-1"))

			err = NewMetadataLoader(nil, cfg, "", "foo-2", nil).Load()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("--name=foo-1 and argument foo-2 cannot be used at the same time"))

			cmd := newCmd()

			err = NewMetadataLoader(nil, cfg, examplesDir+"01-simple-cluster.yaml", "foo-3", cmd).Load()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(ErrCannotUseWithConfigFile(`name argument "foo-3"`).Error()))

			fs := cmd.Flags()

			fs.StringVar(&cfg.Metadata.Name, "name", "", "")
			cmd.Flag("name").Changed = true

			Expect(cmd.Flag("name").Changed).To(BeTrue())
			err = NewMetadataLoader(nil, cfg, examplesDir+"01-simple-cluster.yaml", "foo-3", cmd).Load()
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal(ErrCannotUseWithConfigFile("--name").Error()))
		})

		It("load all of example file", func() {
			examples, err := filepath.Glob(examplesDir + "*.yaml")
			Expect(err).ToNot(HaveOccurred())

			Expect(examples).To(HaveLen(7))
			for _, example := range examples {
				cfg := api.NewClusterConfig()

				p := &api.ProviderConfig{}
				err := NewMetadataLoader(p, cfg, example, "", newCmd()).Load()

				Expect(err).ToNot(HaveOccurred())
				Expect(cfg.Metadata.Name).ToNot(BeEmpty())
				Expect(cfg.Metadata.Region).ToNot(BeEmpty())
				Expect(cfg.Metadata.Region).To(Equal(p.Region))
				Expect(cfg.Metadata.Version).To(BeEmpty())
			}
		})
	})
})
