package podidentityassociation

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	cfntypes "github.com/aws/aws-sdk-go-v2/service/cloudformation/types"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	awsiam "github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/builder"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

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
			"creating pod identity association for service account %q in namespace %q: %w",
			t.podIdentityAssociation.ServiceAccountName, t.podIdentityAssociation.Namespace, err)
	}
	logger.Info(fmt.Sprintf("created pod identity association for service account %q in namespace %q",
		t.podIdentityAssociation.ServiceAccountName, t.podIdentityAssociation.Namespace))
	return nil
}

type trustPolicyUpdater struct {
	iamAPI       awsapi.IAM
	stackUpdater StackUpdater
}

func (t *trustPolicyUpdater) UpdateTrustPolicyForOwnedRoleTask(ctx context.Context, roleName, serviceAccountName string, stack IRSAv1StackSummary, removeOIDCProviderTrustRelationship bool) tasks.Task {
	return &tasks.GenericTask{
		Description: fmt.Sprintf("update trust policy for owned role %q", roleName),
		Doer: func() error {
			trustStatements, err := updateTrustStatements(removeOIDCProviderTrustRelationship, func() (*awsiam.GetRoleOutput, error) {
				return t.iamAPI.GetRole(ctx, &awsiam.GetRoleInput{RoleName: &roleName})
			})
			if err != nil {
				return fmt.Errorf("updating trust statements for role %s: %w", roleName, err)
			}

			// build template for updating trust policy
			rs := builder.NewIAMRoleResourceSetForPodIdentityWithTrustStatements(&api.PodIdentityAssociation{}, trustStatements)
			if err := rs.AddAllResources(); err != nil {
				return fmt.Errorf("adding resources to CloudFormation template: %w", err)
			}
			template, err := rs.RenderJSON()
			if err != nil {
				return fmt.Errorf("generating CloudFormation template: %w", err)
			}

			// update stack tags to reflect migration to IRSAv2
			cfnTags := []cfntypes.Tag{}
			for key, value := range stack.Tags {
				if key == api.IAMServiceAccountNameTag && removeOIDCProviderTrustRelationship {
					continue
				}
				cfnTags = append(cfnTags, cfntypes.Tag{
					Key:   &key,
					Value: &value,
				})
			}

			getIAMServiceAccountName := func() string {
				if serviceAccountName != "" {
					return serviceAccountName
				}
				return strings.Replace(strings.Split(stack.Name, "-iamserviceaccount-")[1], "-", "/", 1)
			}
			cfnTags = append(cfnTags, cfntypes.Tag{
				Key:   aws.String(api.PodIdentityAssociationNameTag),
				Value: aws.String(getIAMServiceAccountName()),
			})

			// propagate capabilities
			cfnCapabilities := []cfntypes.Capability{}
			for _, c := range stack.Capabilities {
				cfnCapabilities = append(cfnCapabilities, cfntypes.Capability(c))
			}

			if err := t.stackUpdater.MustUpdateStack(ctx, manager.UpdateStackOptions{
				Stack: &cfntypes.Stack{
					StackName:    &stack.Name,
					Tags:         cfnTags,
					Capabilities: cfnCapabilities,
				},
				ChangeSetName: fmt.Sprintf("eksctl-%s-update-%d", roleName, time.Now().Unix()),
				Description:   fmt.Sprintf("updating IAM resources stack %q for role %q", stack.Name, roleName),
				TemplateData:  manager.TemplateBody(template),
				Wait:          true,
			}); err != nil {
				var noChangeErr *manager.NoChangeError
				if errors.As(err, &noChangeErr) {
					logger.Info("IAM resources for role %q are already up-to-date", roleName)
					return nil
				}
				return fmt.Errorf("updating IAM resources for role %q: %w", roleName, err)
			}
			logger.Info("updated IAM resources stack %q for role %q", stack.Name, roleName)

			return nil
		},
	}
}

func (t *trustPolicyUpdater) UpdateTrustPolicyForUnownedRoleTask(ctx context.Context, roleName string, removeOIDCProviderTrustRelationship bool) tasks.Task {
	return &tasks.GenericTask{
		Description: fmt.Sprintf("update trust policy for unowned role %q", roleName),
		Doer: func() error {
			trustStatements, err := updateTrustStatements(removeOIDCProviderTrustRelationship, func() (*awsiam.GetRoleOutput, error) {
				return t.iamAPI.GetRole(ctx, &awsiam.GetRoleInput{RoleName: &roleName})
			})
			if err != nil {
				return fmt.Errorf("updating trust statements for role %s: %w", roleName, err)
			}

			documentString, err := json.Marshal(api.IAMPolicyDocument{
				Version:    "2012-10-17",
				Statements: trustStatements,
			})
			if err != nil {
				return fmt.Errorf("marshalling trust policy document: %w", err)
			}

			if _, err := t.iamAPI.UpdateAssumeRolePolicy(ctx, &awsiam.UpdateAssumeRolePolicyInput{
				RoleName:       &roleName,
				PolicyDocument: aws.String(string(documentString)),
			}); err != nil {
				return fmt.Errorf("updating trust policy for role %s: %w", roleName, err)
			}
			logger.Info(fmt.Sprintf("updated trust policy for role %s", roleName))
			return nil
		},
	}
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
	documentJSONString, err := url.PathUnescape(*output.Role.AssumeRolePolicyDocument)
	if err != nil {
		return trustStatements, err
	}
	if err := json.Unmarshal([]byte(documentJSONString), &trustPolicy); err != nil {
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
