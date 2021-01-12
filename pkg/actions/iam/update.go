package iam

import (
	"fmt"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
)

func (a *Manager) UpdateIAMServiceAccount(iamServiceAccounts *api.ClusterIAMServiceAccount, plan bool) error {
	stackName := makeIAMServiceAccountStackName(a.clusterName, iamServiceAccounts.Namespace, iamServiceAccounts.Name)
	stacks, err := a.stackManager.ListStacksMatching(stackName)
	if err != nil {
		return err
	}

	if len(stacks) == 0 {
		return fmt.Errorf("IAMServiceAccount %s/%s does not exist", iamServiceAccounts.Namespace, iamServiceAccounts.Name)
	}

	rs := builder.NewIAMServiceAccountResourceSet(iamServiceAccounts, a.oidcManager)
	err = rs.AddAllResources()
	if err != nil {
		return err
	}

	template, err := rs.RenderJSON()
	if err != nil {
		return err
	}

	var templateBody manager.TemplateBody = template
	err = a.stackManager.UpdateStack(stackName, "updating-policy", "updating policies", templateBody, nil)
	if err != nil {
		return err
	}
	return nil
}

func makeIAMServiceAccountStackName(clusterName, namespace, name string) string {
	return fmt.Sprintf("eksctl-%s-addon-iamserviceaccount-%s-%s", clusterName, namespace, name)
}
