// Package authconfigmap allows manipulation of the EKS auth ConfigMap (aws-auth),
// which maps IAM entities to Kubernetes groups.
//
// See for more information:
// - https://docs.aws.amazon.com/eks/latest/userguide/add-user-role.html
// - https://github.com/kubernetes-sigs/aws-iam-authenticator/blob/master/README.md#full-configuration-format
package authconfigmap

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awsutil"
	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/typed/core/v1"
	"sigs.k8s.io/yaml"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/iam"
)

const (
	// ObjectName is the Kubernetes resource name of the auth ConfigMap
	ObjectName = "aws-auth"
	// ObjectNamespace is the namespace the object can be found
	ObjectNamespace = metav1.NamespaceSystem

	rolesData = "mapRoles"
	roleKey   = "rolearn"

	usersData = "mapUsers"
	userKey   = "userarn"

	accountsData = "mapAccounts"

	// GroupMasters is the admin group which is also automatically
	// granted to the IAM role that creates the cluster.
	GroupMasters = "system:masters"

	// RoleNodeGroupUsername is the default username for a nodegroup
	// role mapping.
	RoleNodeGroupUsername = "system:node:{{EC2PrivateDNSName}}"
)

// RoleNodeGroupGroups are the groups to allow roles to interact
// with the cluster, required for the instance role ARNs of nodegroups.
var RoleNodeGroupGroups = []string{"system:bootstrappers", "system:nodes"}

// MapIdentity represents an IAM identity with an ARN.
type MapIdentity struct {
	iam.Identity `json:",inline"`
	ARN          string `json:"-"` // This field is (un)marshaled manually
}

// MapIdentities is a list of IAM identities with a role or user ARN.
type MapIdentities []MapIdentity

// UnmarshalJSON for MapIdentity makes sure there is either a "rolearn" or "userarn" key
// and places the data of that key into the "ARN" field. All other fields are handled by
// the default implementation of unmarshal.
func (m *MapIdentity) UnmarshalJSON(data []byte) error {
	// Handle everything as usual
	type alias MapIdentity
	var a alias
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	// We want to unmarshal "(rolearn|userarn)" into the "ARN" field and then unmarshal
	// the rest as usual
	outerKeys := map[string]*json.RawMessage{}
	if err := json.Unmarshal(data, &outerKeys); err != nil {
		return err
	}

	arn, ok := outerKeys[roleKey]
	if !ok {
		arn, ok = outerKeys[userKey]
		if !ok {
			return errors.New("missing arn")
		}
	}

	if err := json.Unmarshal(*arn, &a.ARN); err != nil {
		return err
	}

	*m = MapIdentity(a)
	return nil
}

// MarshalJSON for MapIdentity marshals the ARN field into either "rolearn" or "userarn"
// depending on what is appropriate and returns and error if it cannot determine that.
// All other fields are marshaled with the default marshaler
func (m *MapIdentity) MarshalJSON() ([]byte, error) {
	// Marshal everything as one would normally do
	type alias MapIdentity
	a := (*alias)(m)
	data, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	// Partially unmarshal everything into a map in order to more easily append the (role|user)arn
	partial := map[string]*json.RawMessage{}
	if err := json.Unmarshal(data, &partial); err != nil {
		return nil, err
	}

	arn, err := json.Marshal(m.ARN)
	if err != nil {
		return nil, err
	}

	switch {
	case m.role():
		partial[roleKey] = (*json.RawMessage)(&arn)
	case m.user():
		partial[userKey] = (*json.RawMessage)(&arn)
	default:
		return nil, errors.Errorf("cannot determine if %q refers to a user or role", arn)
	}

	return json.Marshal(partial)
}

// resource returns the resource portion of the ARN as described by
// https://docs.aws.amazon.com/general/latest/gr/aws-arns-and-namespaces.html
//
// arn:partition:service:region:account-id:resource
// arn:partition:service:region:account-id:resourcetype/resource
// arn:partition:service:region:account-id:resourcetype/resource/qualifier
// arn:partition:service:region:account-id:resourcetype/resource:qualifier
// arn:partition:service:region:account-id:resourcetype:resource
// arn:partition:service:region:account-id:resourcetype:resource:qualifier
func (m MapIdentity) resource() string {
	portions := strings.Split(m.ARN, ":")
	if len(portions) < 6 {
		// malformed arn
		return ""
	}

	if idx := strings.Index(portions[5], "/"); idx >= 0 {
		return portions[5][:idx] // Only return up until the forward slash (if one is present)
	}

	return portions[5]
}

func (m MapIdentity) role() bool {
	return m.resource() == "role"
}

func (m MapIdentity) user() bool {
	return m.resource() == "user"
}

// Get returns all matching role mappings. Note that at this moment
// aws-iam-authenticator only considers the last one!
func (rs MapIdentities) Get(arn string) MapIdentities {
	var m MapIdentities
	for _, r := range rs {
		if r.ARN == arn {
			m = append(m, r)
		}
	}
	return m
}

// AuthConfigMap allows modifying the auth ConfigMap.
type AuthConfigMap struct {
	client v1.ConfigMapInterface
	cm     *corev1.ConfigMap
}

// New creates an AuthConfigMap instance that manipulates
// a ConfigMap. If it is nil, one is created.
func New(client v1.ConfigMapInterface, cm *corev1.ConfigMap) *AuthConfigMap {
	if cm == nil {
		cm = &corev1.ConfigMap{
			ObjectMeta: ObjectMeta(),
			Data:       map[string]string{},
		}
	}
	if cm.Data == nil {
		cm.ObjectMeta = ObjectMeta()
		cm.Data = map[string]string{}
	}
	return &AuthConfigMap{client: client, cm: cm}
}

// NewFromClientSet fetches the auth ConfigMap.
func NewFromClientSet(clientSet kubernetes.Interface) (*AuthConfigMap, error) {
	client := clientSet.CoreV1().ConfigMaps(ObjectNamespace)

	cm, err := client.Get(ObjectName, metav1.GetOptions{})
	// It is fine for the configmap not to exist. Any other error is fatal.
	if err != nil && !kerr.IsNotFound(err) {
		return nil, errors.Wrapf(err, "getting auth ConfigMap")
	}
	logger.Debug("aws-auth = %s", awsutil.Prettify(cm))
	return New(client, cm), nil
}

// AddAccount appends an IAM account to the `mapAccounts` entry
// in the Configmap. It also deduplicates.
func (a *AuthConfigMap) AddAccount(account string) error {
	accounts, err := a.accounts()
	if err != nil {
		return err
	}
	// Distinct and sorted account numbers
	accounts = append(accounts, account)
	accounts = sets.NewString(accounts...).List()
	logger.Info("adding account %q to auth ConfigMap", account)
	return a.setAccounts(accounts)
}

// RemoveAccount removes the given IAM account entry in mapAccounts.
func (a *AuthConfigMap) RemoveAccount(account string) error {
	accounts, err := a.accounts()
	if err != nil {
		return err
	}

	var newaccounts []string
	found := false
	for _, acc := range accounts {
		if acc == account {
			found = true
			continue
		}
		newaccounts = append(newaccounts, acc)
	}
	if !found {
		return fmt.Errorf("account %q not found in auth ConfigMap", account)
	}
	logger.Info("removing account %q from auth ConfigMap", account)
	return a.setAccounts(newaccounts)
}

func (a *AuthConfigMap) accounts() ([]string, error) {
	var accounts []string
	if err := yaml.Unmarshal([]byte(a.cm.Data[accountsData]), &accounts); err != nil {
		return nil, errors.Wrap(err, "unmarshalling mapAccounts")
	}
	return accounts, nil
}

func (a *AuthConfigMap) setAccounts(accounts []string) error {
	bs, err := yaml.Marshal(accounts)
	if err != nil {
		return errors.Wrap(err, "marshalling mapAccounts")
	}
	a.cm.Data[accountsData] = string(bs)
	return nil
}

// AddIdentity maps an IAM role or user ARN to a k8s group dynamically. It modifies the
// role or user with given groups. If you are calling
// this as part of node creation you should use DefaultNodeGroups.
func (a *AuthConfigMap) AddIdentity(arn string, username string, groups []string) error {
	identities, err := a.Identities()
	if err != nil {
		return err
	}
	identities = append(identities, MapIdentity{
		ARN: arn,
		Identity: iam.Identity{
			Username: username,
			Groups:   groups,
		},
	})
	logger.Info("adding identity %q to auth ConfigMap", arn)
	return a.setIdentities(identities)
}

// RemoveIdentity removes an identity. If `all` is false it will only
// remove the first it encounters and return an error if it cannot
// find it.
// If `all` is true it will remove all of them and not return an
// error if it cannot be found.
func (a *AuthConfigMap) RemoveIdentity(arn string, all bool) error {
	identities, err := a.Identities()
	if err != nil {
		return err
	}

	newidentities := MapIdentities{}
	for i, identity := range identities {
		if identity.ARN == arn {
			logger.Info("removing identity %q from auth ConfigMap (username = %q, groups = %q)", arn, identity.Username, identity.Groups)
			if !all {
				identities = append(identities[:i], identities[i+1:]...)
				return a.setIdentities(identities)
			}
		} else if all {
			newidentities = append(newidentities, identity)
		}
	}
	if !all {
		return fmt.Errorf("instance identity ARN %q not found in auth ConfigMap", arn)
	}
	return a.setIdentities(newidentities)
}

// Identities returns a list of iam users and roles that are currently in the (cached) configmap.
func (a *AuthConfigMap) Identities() (MapIdentities, error) {
	var roles, users MapIdentities
	if err := yaml.Unmarshal([]byte(a.cm.Data[rolesData]), &roles); err != nil {
		return nil, errors.Wrapf(err, "unmarshalling %q", rolesData)
	}

	if err := yaml.Unmarshal([]byte(a.cm.Data[usersData]), &users); err != nil {
		return nil, errors.Wrapf(err, "unmarshalling %q", usersData)
	}
	return append(roles, users...), nil
}

func (a *AuthConfigMap) setIdentities(identities MapIdentities) error {
	// Determine which are users and which are roles
	var users, roles MapIdentities
	for _, identity := range identities {
		switch {
		case identity.role():
			roles = append(roles, identity)
		case identity.user():
			users = append(users, identity)
		default:
			return errors.Errorf("cannot determine if %q refers to a user or role during setIdentities preprocessing", identity)
		}
	}

	// Update the corresponding keys
	_roles, err := yaml.Marshal(roles)
	if err != nil {
		return errors.Wrapf(err, "marshalling %q", rolesData)
	}
	a.cm.Data[rolesData] = string(_roles)

	_users, err := yaml.Marshal(users)
	if err != nil {
		return errors.Wrapf(err, "marshalling %q", usersData)
	}
	a.cm.Data[usersData] = string(_users)

	return nil
}

// Save persists the ConfigMap to the cluster. It determines
// whether to create or update by looking at the ConfigMap's UID.
func (a *AuthConfigMap) Save() (err error) {
	if a.cm.UID == "" {
		a.cm, err = a.client.Create(a.cm)
		return err
	}

	a.cm, err = a.client.Update(a.cm)
	return err
}

// ObjectMeta constructs metadata for the ConfigMap.
func ObjectMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      ObjectName,
		Namespace: ObjectNamespace,
	}
}

// AddNodeGroup creates or adds a nodegroup IAM role in the auth
// ConfigMap for the given nodegroup.
func AddNodeGroup(clientSet kubernetes.Interface, ng *api.NodeGroup) error {
	acm, err := NewFromClientSet(clientSet)
	if err != nil {
		return err
	}
	if err := acm.AddIdentity(ng.IAM.InstanceRoleARN, RoleNodeGroupUsername, RoleNodeGroupGroups); err != nil {
		return errors.Wrap(err, "adding nodegroup to auth ConfigMap")
	}
	if err := acm.Save(); err != nil {
		return errors.Wrap(err, "saving auth ConfigMap")
	}
	logger.Debug("saved auth ConfigMap for %q", ng.Name)
	return nil
}

// RemoveNodeGroup removes a nodegroup from the ConfigMap and
// does a client update.
func RemoveNodeGroup(clientSet kubernetes.Interface, ng *api.NodeGroup) error {
	arn := ng.IAM.InstanceRoleARN
	if arn == "" {
		return errors.New("nodegroup instance role ARN is not set")
	}
	acm, err := NewFromClientSet(clientSet)
	if err != nil {
		return err
	}
	if err := acm.RemoveIdentity(arn, false); err != nil {
		return errors.Wrap(err, "removing nodegroup from auth ConfigMap")
	}
	if err := acm.Save(); err != nil {
		return errors.Wrap(err, "updating auth ConfigMap after removing role")
	}
	logger.Debug("updated auth ConfigMap for %s", ng.Name)
	return nil
}
