package podidentityassociation

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	awsiam "github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

type createIAMRoleTask struct {
	ctx                    context.Context
	info                   string
	clusterName            string
	podIdentityAssociation *api.PodIdentityAssociation
	stackCreator           StackCreator
}

func (t *createIAMRoleTask) Describe() string {
	return t.info
}

func (t *createIAMRoleTask) Do(errorCh chan error) error {
	rs := builder.NewIAMRoleResourceSetForPodIdentity(t.podIdentityAssociation)
	if err := rs.AddAllResources(); err != nil {
		return err
	}
	if err := t.stackCreator.CreateStack(t.ctx,
		MakeStackName(
			t.clusterName,
			t.podIdentityAssociation.Namespace,
			t.podIdentityAssociation.ServiceAccountName),
		rs, nil, nil, errorCh); err != nil {
		return fmt.Errorf("creating IAM role for pod identity association for service account %s in namespace %s: %w",
			t.podIdentityAssociation.ServiceAccountName, t.podIdentityAssociation.Namespace, err)
	}
	return nil
}

type createPodIdentityAssociationTask struct {
	ctx                    context.Context
	info                   string
	clusterName            string
	podIdentityAssociation *api.PodIdentityAssociation
	eksAPI                 awsapi.EKS
}

func (t *createPodIdentityAssociationTask) Describe() string {
	return t.info
}

func (t *createPodIdentityAssociationTask) Do(errorCh chan error) error {
	defer close(errorCh)

	if _, err := t.eksAPI.CreatePodIdentityAssociation(t.ctx, &awseks.CreatePodIdentityAssociationInput{
		ClusterName:    &t.clusterName,
		Namespace:      &t.podIdentityAssociation.Namespace,
		RoleArn:        &t.podIdentityAssociation.RoleARN,
		ServiceAccount: &t.podIdentityAssociation.ServiceAccountName,
		Tags:           t.podIdentityAssociation.Tags,
	}); err != nil {
		return fmt.Errorf(
			"creating pod identity association for service account `%s` in namespace `%s`: %w",
			t.podIdentityAssociation.ServiceAccountName, t.podIdentityAssociation.Namespace, err)
	}
	logger.Info(fmt.Sprintf("created pod identity association for service account `%s` in namespace `%s`",
		t.podIdentityAssociation.ServiceAccountName, t.podIdentityAssociation.Namespace))
	return nil
}

type updateTrustPolicyForOwnedRole struct {
	ctx                                 context.Context
	info                                string
	roleName                            string
	stackName                           string
	removeOIDCProviderTrustRelationship bool
	iamAPI                              awsapi.IAM
	stackUpdater                        StackUpdater
}

func (t *updateTrustPolicyForOwnedRole) Describe() string {
	return t.info
}

func (t *updateTrustPolicyForOwnedRole) Do(errorCh chan error) error {
	defer close(errorCh)

	trustStatements, err := updateTrustStatements(t.removeOIDCProviderTrustRelationship, func() (*awsiam.GetRoleOutput, error) {
		return t.iamAPI.GetRole(t.ctx, &awsiam.GetRoleInput{RoleName: &t.roleName})
	})
	if err != nil {
		return fmt.Errorf("updating trust statements for role %s: %w", t.roleName, err)
	}

	rs := builder.NewIAMRoleResourceSetForPodIdentityWithTrustStatements(&api.PodIdentityAssociation{}, trustStatements)
	if err := rs.AddAllResources(); err != nil {
		return fmt.Errorf("adding resources to CloudFormation template: %w", err)
	}
	template, err := rs.RenderJSON()
	if err != nil {
		return fmt.Errorf("generating CloudFormation template: %w", err)
	}

	if err := t.stackUpdater.MustUpdateStack(t.ctx, manager.UpdateStackOptions{
		StackName:     t.stackName,
		ChangeSetName: fmt.Sprintf("eksctl-%s-update-%d", t.roleName, time.Now().Unix()),
		Description:   fmt.Sprintf("updating IAM resources stack %q for role %q", t.stackName, t.roleName),
		TemplateData:  manager.TemplateBody(template),
		Wait:          true,
	}); err != nil {
		if _, ok := err.(*manager.NoChangeError); ok {
			logger.Info("IAM resources for role %q are already up-to-date", t.roleName)
			return nil
		}
		return fmt.Errorf("updating IAM resources for role %q: %w", t.roleName, err)
	}
	logger.Info("updated IAM resources stack %q for role %q", t.stackName, t.roleName)

	return nil
}

type updateTrustPolicyForUnownedRole struct {
	ctx                                 context.Context
	info                                string
	roleName                            string
	iamAPI                              awsapi.IAM
	removeOIDCProviderTrustRelationship bool
}

func (t *updateTrustPolicyForUnownedRole) Describe() string {
	return t.info
}

func (t *updateTrustPolicyForUnownedRole) Do(errorCh chan error) error {
	defer close(errorCh)

	trustStatements, err := updateTrustStatements(t.removeOIDCProviderTrustRelationship, func() (*awsiam.GetRoleOutput, error) {
		return t.iamAPI.GetRole(t.ctx, &awsiam.GetRoleInput{RoleName: &t.roleName})
	})
	if err != nil {
		return fmt.Errorf("updating trust statements for role %s: %w", t.roleName, err)
	}

	documentString, err := json.Marshal(api.IAMPolicyDocument{
		Version:    "2012-10-17",
		Statements: trustStatements,
	})
	if err != nil {
		return fmt.Errorf("marshalling trust policy document: %w", err)
	}

	if _, err := t.iamAPI.UpdateAssumeRolePolicy(t.ctx, &awsiam.UpdateAssumeRolePolicyInput{
		RoleName:       &t.roleName,
		PolicyDocument: aws.String(string(documentString)),
	}); err != nil {
		return fmt.Errorf("updating trust policy for role %s: %w", t.roleName, err)
	}
	logger.Info(fmt.Sprintf("updated trust policy for role %s", t.roleName))

	return nil
}

func updateTrustStatements(
	removeOIDCProviderTrustRelationship bool,
	getRole func() (*awsiam.GetRoleOutput, error),
) ([]api.IAMStatement, error) {
	var trustStatements []api.IAMStatement
	var trustPolicy api.IAMPolicyDocument

	output, err := getRole()
	if err != nil {
		return trustStatements, err
	}
	documentJsonString, err := url.PathUnescape(*output.Role.AssumeRolePolicyDocument)
	if err != nil {
		return trustStatements, err
	}
	if err := json.Unmarshal([]byte(documentJsonString), &trustPolicy); err != nil {
		return trustStatements, err
	}

	shouldRemoveStatement := func(s api.IAMStatement) bool {
		value, ok := s.Principal["Federated"]
		if ok && len(value) == 1 &&
			strings.Contains(value[0], "oidc-provider") &&
			removeOIDCProviderTrustRelationship {
			return true
		}
		return false
	}

	// remove OIDC provider trust relationship if instructed so
	for _, s := range trustPolicy.Statements {
		if shouldRemoveStatement(s) {
			continue
		}
		trustStatements = append(trustStatements, s)
	}

	// add trust relationship with new EKS Service Principal
	trustStatements = append(trustStatements, api.EKSServicePrincipalTrustStatement)

	return trustStatements, nil
}

func makeStackNamePrefix(clusterName string) string {
	return fmt.Sprintf("eksctl-%s-podidentityrole-ns-", clusterName)
}

// MakeStackName creates a stack name for the specified access entry.
func MakeStackName(clusterName, namespace, serviceAccountName string) string {
	return fmt.Sprintf("%s%s-sa-%s", makeStackNamePrefix(clusterName), namespace, serviceAccountName)
}

func runAllTasks(taskTree *tasks.TaskTree) error {
	logger.Info(taskTree.Describe())
	if errs := taskTree.DoAllSync(); len(errs) > 0 {
		var allErrs []string
		for _, err := range errs {
			allErrs = append(allErrs, err.Error())
		}
		return fmt.Errorf(strings.Join(allErrs, "\n"))
	}
	completedAction := func() string {
		if taskTree.PlanMode {
			return "skipped"
		}
		return "completed successfully"
	}
	logger.Info("all tasks were %s", completedAction())
	return nil
}
