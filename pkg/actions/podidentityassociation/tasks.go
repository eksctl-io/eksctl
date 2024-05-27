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
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	awsiam "github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/kris-nova/logger"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/awsapi"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	cft "github.com/weaveworks/eksctl/pkg/cfn/template"
	"github.com/weaveworks/eksctl/pkg/utils/tasks"
)

const (
	resourceTypeIAMRole = "AWS::IAM::Role"
)

type createPodIdentityAssociationTask struct {
	ctx                        context.Context
	info                       string
	clusterName                string
	podIdentityAssociation     *api.PodIdentityAssociation
	eksAPI                     awsapi.EKS
	ignorePodIdentityExistsErr bool
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
		if t.ignorePodIdentityExistsErr {
			var inUseErr *ekstypes.ResourceInUseException
			if errors.As(err, &inUseErr) {
				logger.Info("pod identity association %s already exists", t.podIdentityAssociation.NameString())
				return nil
			}
		}
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

			currentTemplate, err := t.stackUpdater.GetStackTemplate(ctx, stack.Name)
			if err != nil {
				return fmt.Errorf("fetching current template for stack %q", stack.Name)
			}

			cfnTemplate := cft.NewTemplate()
			if err := cfnTemplate.LoadJSON([]byte(currentTemplate)); err != nil {
				return fmt.Errorf("unmarshalling current template for stack %q", stack.Name)
			}

			for i, r := range cfnTemplate.Resources {
				if r.Type != resourceTypeIAMRole {
					continue
				}
				role, err := r.ToIAMRole()
				if err != nil {
					return fmt.Errorf("fetching properties for role %s: %w", roleName, err)
				}
				role.AssumeRolePolicyDocument["Statement"] = trustStatements
				r.Properties = role
				cfnTemplate.Resources[i] = r
			}

			updatedTemplate, err := cfnTemplate.RenderJSON()
			if err != nil {
				return fmt.Errorf("marshalling updated template for stack %q", stack.Name)
			}
			logger.Debug("updated template for role %s: %v", string(updatedTemplate))

			// update stack tags to reflect migration to IRSAv2
			cfnTags := []cfntypes.Tag{}
			for key, value := range stack.Tags {
				if key == api.IAMServiceAccountNameTag && removeOIDCProviderTrustRelationship {
					continue
				}
				cfnTags = append(cfnTags, cfntypes.Tag{
					Key:   aws.String(key),
					Value: aws.String(value),
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
				TemplateData:  manager.TemplateBody(updatedTemplate),
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
