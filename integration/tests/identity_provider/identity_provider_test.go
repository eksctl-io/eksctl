//go:build integration
// +build integration

//revive:disable
package identity_provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	// Register the OIDC provider
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"github.com/sethvargo/go-password/password"

	. "github.com/weaveworks/eksctl/integration/matchers"
	. "github.com/weaveworks/eksctl/integration/runner"

	"github.com/weaveworks/eksctl/integration/tests"
	"github.com/weaveworks/eksctl/integration/utilities/kube"
	"github.com/weaveworks/eksctl/pkg/testutils"
)

const (
	oidcGroupName = "oidc-reader"
)

var (
	params                  *tests.Params
	oidcConfig              *OIDCConfig
	cleanupCognitoResources func() error
)

type OIDCConfig struct {
	clientID     string
	idToken      string
	refreshToken string
	idpIssuerURL string
}

func init() {
	// Call testing.Init() prior to tests.NewParams(), as otherwise -test.* will not be recognised. See also: https://golang.org/doc/go1.13#testing
	testing.Init()
	params = tests.NewParams("identity-provider")
}

func TestIdentityProvider(t *testing.T) {
	testutils.RegisterAndRun(t)
}

var _ = BeforeSuite(func() {
	if !params.SkipCreate {
		cmd := params.EksctlCreateCmd.WithArgs(
			"cluster",
			"--verbose", "4",
			"--name", params.ClusterName,
			"--kubeconfig", params.KubeconfigPath,
			"--nodes", "1",
		)
		Expect(cmd).To(RunSuccessfully())
	}

	var err error
	fmt.Fprintf(GinkgoWriter, "creating Cognito OIDC provider\n")
	oidcConfig, err = setupCognitoProvider(params.ClusterName, params.Region)
	Expect(err).NotTo(HaveOccurred())
	fmt.Fprintf(GinkgoWriter, "created Cognito provider; client ID: %s\n", oidcConfig.clientID)
})

var _ = Describe("(Integration) [Identity Provider]", func() {

	It("should associate, get and disassociate identity provider", func() {
		By("associating a new identity provider")
		identityProviderClusterConfig := makeIdentityProviderClusterConfig(oidcConfig, params.ClusterName, params.Region)

		cmd := params.EksctlCmd.WithArgs(
			"associate",
			"identityprovider",
			"--config-file", "-",
			"--verbose", "4",
			"--wait",
		).
			WithStdin(strings.NewReader(identityProviderClusterConfig)).
			WithoutArg("--region", params.Region).
			WithTimeout(1 * time.Hour)

		Expect(cmd).To(RunSuccessfully())

		By("getting the identity provider")
		cmd = params.EksctlGetCmd.
			WithArgs(
				"identityprovider",
				"--cluster", params.ClusterName,
				"-o", "yaml",
			)
		Expect(cmd).To(RunSuccessfullyWithOutputStringLines(
			ContainElement(ContainSubstring(fmt.Sprintf("ClientID: %s", oidcConfig.clientID)))),
		)

		By("creating RBAC resources")
		test, err := kube.NewTest(params.KubeconfigPath)
		Expect(err).NotTo(HaveOccurred())
		defer test.Close()

		test.CreateClusterRoleFromFile("testdata/cluster-role.yaml")
		test.CreateClusterRoleBindingFromFile("testdata/cluster-role-binding.yaml")

		By("creating an OIDC Clientset")
		clientset, err := createOIDCClientset(eks.New(NewSession(params.Region)), oidcConfig, params.ClusterName)
		Expect(err).NotTo(HaveOccurred())

		By("reading Kubernetes resources")
		list, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		Expect(err).NotTo(HaveOccurred())
		Expect(list.Items).To(HaveLen(1))

		secrets, err := clientset.CoreV1().Secrets(metav1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
		Expect(err).NotTo(HaveOccurred())
		Expect(secrets.Items).NotTo(BeEmpty())

		By("ensuring the client does not have write access")
		_, err = clientset.CoreV1().ConfigMaps(metav1.NamespaceDefault).Create(context.TODO(), &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name: "testdata",
			},
			Data: map[string]string{
				"key": "value",
			},
		}, metav1.CreateOptions{})

		Expect(err).To(HaveOccurred())

		By("disassociating the identity provider")
		cmd = params.EksctlCmd.
			WithArgs(
				"disassociate",
				"identityprovider",
				"--config-file", "-",
				"--wait",
			).
			WithStdin(strings.NewReader(identityProviderClusterConfig)).
			WithoutArg("--region", params.Region).
			WithTimeout(1 * time.Hour)

		Expect(cmd).To(RunSuccessfully())

	})

})

var _ = AfterSuite(func() {
	if !params.SkipCreate && !params.SkipDelete {
		params.DeleteClusters()
	}
	if cleanupCognitoResources != nil {
		Expect(cleanupCognitoResources()).To(Succeed())
	}
})

func createOIDCClientset(eksAPI eksiface.EKSAPI, o *OIDCConfig, clusterName string) (kubernetes.Interface, error) {
	contextName := fmt.Sprintf("%s@%s", "test", clusterName)

	cluster, err := eksAPI.DescribeCluster(&eks.DescribeClusterInput{
		Name: aws.String(clusterName),
	})
	if err != nil {
		return nil, fmt.Errorf("describing cluster: %w", err)
	}

	certData, err := base64.StdEncoding.DecodeString(*cluster.Cluster.CertificateAuthority.Data)
	if err != nil {
		return nil, fmt.Errorf("unexpected error decoding certificate authority data: %w", err)
	}

	config := clientcmdapi.Config{
		Clusters: map[string]*clientcmdapi.Cluster{
			clusterName: {
				Server:                   *cluster.Cluster.Endpoint,
				CertificateAuthorityData: certData,
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			contextName: {
				Cluster:  clusterName,
				AuthInfo: contextName,
			},
		},
		AuthInfos: map[string]*clientcmdapi.AuthInfo{
			contextName: {
				AuthProvider: &clientcmdapi.AuthProviderConfig{
					Name: "oidc",
					Config: map[string]string{
						"client-id":      o.clientID,
						"id-token":       o.idToken,
						"refresh-token":  o.refreshToken,
						"idp-issuer-url": o.idpIssuerURL,
					},
				},
			},
		},
		CurrentContext: contextName,
	}

	clientConfig, err := clientcmd.NewDefaultClientConfig(config, &clientcmd.ConfigOverrides{}).ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("creating default client config: %w", err)
	}

	return kubernetes.NewForConfig(clientConfig)
}

type userPoolClient struct {
	userPoolID *string
	clientID   *string
}

func setupCognitoProvider(clusterName, region string) (*OIDCConfig, error) {
	c := cognitoidentityprovider.New(NewSession(region))

	userPoolClient, err := createCognitoUserPoolClient(c, clusterName)
	if err != nil {
		return nil, err
	}

	userPass, err := password.Generate(10, 2, 3, false, false)
	if err != nil {
		return nil, fmt.Errorf("generating password: %w", err)
	}

	clientUsername := aws.String(fmt.Sprintf("%s@weave.works", clusterName))
	if err := createCognitoUserGroup(c, userPoolClient.userPoolID, clientUsername, userPass); err != nil {
		return nil, err
	}

	auth, err := c.AdminInitiateAuth(&cognitoidentityprovider.AdminInitiateAuthInput{
		AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeAdminUserPasswordAuth),
		AuthParameters: map[string]*string{
			"USERNAME": clientUsername,
			"PASSWORD": aws.String(userPass),
		},
		ClientId:   userPoolClient.clientID,
		UserPoolId: userPoolClient.userPoolID,
	})

	if err != nil {
		return nil, fmt.Errorf("initiating auth: %w", err)
	}

	return &OIDCConfig{
		clientID:     *userPoolClient.clientID,
		idToken:      *auth.AuthenticationResult.IdToken,
		refreshToken: *auth.AuthenticationResult.RefreshToken,
		idpIssuerURL: fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", region, *userPoolClient.userPoolID),
	}, nil
}

func createCognitoUserPoolClient(c cognitoidentityprovideriface.CognitoIdentityProviderAPI, clusterName string) (*userPoolClient, error) {
	pool, err := c.CreateUserPool(&cognitoidentityprovider.CreateUserPoolInput{
		Policies: &cognitoidentityprovider.UserPoolPolicyType{
			PasswordPolicy: &cognitoidentityprovider.PasswordPolicyType{
				MinimumLength:    aws.Int64(10),
				RequireLowercase: aws.Bool(false),
				RequireNumbers:   aws.Bool(true),
				RequireSymbols:   aws.Bool(true),
				RequireUppercase: aws.Bool(false),
			},
		},
		PoolName:           aws.String(clusterName),
		UsernameAttributes: aws.StringSlice([]string{"email"}),
	})

	if err != nil {
		return nil, fmt.Errorf("creating user pool: %w", err)
	}

	cleanupCognitoResources = func() error {
		_, err := c.DeleteUserPool(&cognitoidentityprovider.DeleteUserPoolInput{
			UserPoolId: pool.UserPool.Id,
		})
		return err
	}

	client, err := c.CreateUserPoolClient(&cognitoidentityprovider.CreateUserPoolClientInput{
		ClientName: aws.String("eks-client"),
		ExplicitAuthFlows: aws.StringSlice([]string{
			cognitoidentityprovider.ExplicitAuthFlowsTypeAllowAdminUserPasswordAuth,
			cognitoidentityprovider.ExplicitAuthFlowsTypeAllowUserPasswordAuth,
			cognitoidentityprovider.ExplicitAuthFlowsTypeAllowRefreshTokenAuth,
		}),
		AllowedOAuthFlows: aws.StringSlice([]string{
			cognitoidentityprovider.OAuthFlowTypeImplicit,
		}),
		AllowedOAuthScopes: aws.StringSlice([]string{
			"profile",
			"phone",
			"email",
			"openid",
			"aws.cognito.signin.user.admin",
		}),
		UserPoolId:                 pool.UserPool.Id,
		GenerateSecret:             aws.Bool(false),
		SupportedIdentityProviders: aws.StringSlice([]string{"COGNITO"}),
		// TODO this is likely not required, check if this can be removed.
		CallbackURLs: aws.StringSlice([]string{"https://example.com"}),
	})

	if err != nil {
		return nil, fmt.Errorf("creating user pool client: %w", err)
	}
	return &userPoolClient{
		userPoolID: pool.UserPool.Id,
		clientID:   client.UserPoolClient.ClientId,
	}, nil
}

func createCognitoUserGroup(c cognitoidentityprovideriface.CognitoIdentityProviderAPI, userPoolID, clientUsername *string, userPass string) error {
	_, err := c.AdminCreateUser(&cognitoidentityprovider.AdminCreateUserInput{
		UserPoolId: userPoolID,
		Username:   clientUsername,
	})

	if err != nil {
		return fmt.Errorf("creating user: %w", err)
	}

	_, err = c.AdminSetUserPassword(&cognitoidentityprovider.AdminSetUserPasswordInput{
		UserPoolId: userPoolID,
		Username:   clientUsername,
		Password:   aws.String(userPass),
		Permanent:  aws.Bool(true),
	})

	if err != nil {
		return fmt.Errorf("setting user password: %w", err)
	}

	groupName := aws.String(oidcGroupName)

	_, err = c.CreateGroup(&cognitoidentityprovider.CreateGroupInput{
		UserPoolId: userPoolID,
		GroupName:  groupName,
	})

	if err != nil {
		return fmt.Errorf("creating group: %w", err)
	}

	_, err = c.AdminAddUserToGroup(&cognitoidentityprovider.AdminAddUserToGroupInput{
		GroupName:  groupName,
		UserPoolId: userPoolID,
		Username:   clientUsername,
	})

	if err != nil {
		return fmt.Errorf("adding user to group: %w", err)
	}
	return nil
}

func makeIdentityProviderClusterConfig(o *OIDCConfig, clusterName, region string) string {
	return fmt.Sprintf(`apiVersion: eksctl.io/v1alpha5
kind: ClusterConfig

metadata:
  name: %s
  region: %s

identityProviders:
  - name: cognito
    issuerURL: %s
    clientID: %s
    usernameClaim: email
    groupsClaim: "cognito:groups"
    groupsPrefix: "gid:"
    type: oidc
`, clusterName, region, o.idpIssuerURL, o.clientID)
}
