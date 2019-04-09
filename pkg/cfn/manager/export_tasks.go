package manager

import (
	gfn "github.com/awslabs/goformation/cloudformation"
	"github.com/kris-nova/logger"
	"k8s.io/apimachinery/pkg/util/sets"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
)

// ExportClusterWithNodeGroups prepares all templates for export
func (c *StackCollection) ExportClusterWithNodeGroups(onlySubset sets.String) (map[string]gfn.Template, []error) {
	name := c.makeClusterStackName()
	logger.Debug("preparing cluster stack %q", name)

	stack := builder.NewClusterResourceSet(c.provider, c.spec)

	templates := map[string]gfn.Template{}
	templates[name] = stack.Template()

	errs := []error{}
	if err := stack.AddAllResources(); err != nil {
		errs = append(errs, err)
	}

	for i := range c.spec.NodeGroups {
		ng := c.spec.NodeGroups[i]
		if onlySubset != nil && !onlySubset.Has(ng.Name) {
			continue
		}
		name := c.makeNodeGroupStackName(ng.Name)
		logger.Info("exporting nodegroup stack %q", name)

		stack := builder.NewNodeGroupResourceSet(c.provider, c.spec, c.makeClusterStackName(), ng)
		templates[name] = stack.Template()
		if err := stack.AddAllResources(); err != nil {
			errs = append(errs, err)
		}
		if ng.Tags == nil {
			ng.Tags = make(map[string]string)
		}
		ng.Tags[api.NodeGroupNameTag] = ng.Name
	}
	if len(c.spec.NodeGroups) == 0 {
		logger.Warning("no node groups to export")
	}

	return templates, errs
}

