// Package authconfigmap allows manipulation of the EKS configmap,
// which maps IAM roles to Kubernetes groups.
// See https://docs.aws.amazon.com/eks/latest/userguide/add-user-role.html
// for more information.
package authconfigmap

import (
	"fmt"
	"sort"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/typed/core/v1"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"
)

type mapRole map[string]interface{}
type mapRoles []mapRole

const (
	objectName      = "aws-auth"
	objectNamespace = metav1.NamespaceSystem

	rolesData    = "mapRoles"
	accountsData = "mapAccounts"

	// GroupMasters is the admin group which is also automatically
	// granted to the IAM role that creates the cluster.
	GroupMasters = "system:masters"
)

// DefaultNodeGroups are the groups to allow roles to interact
// with the cluster, required for the instance role ARNs of node groups.
var DefaultNodeGroups = []string{"system:bootstrappers", "system:nodes"}

// AuthConfigMap allows modifying the auth configmap.
type AuthConfigMap struct {
	cm *corev1.ConfigMap
}

// New creates an AuthConfigMap instance that manipulates
// a config map. If it is nil, one is created.
func New(cm *corev1.ConfigMap) *AuthConfigMap {
	a := &AuthConfigMap{cm: cm}
	if a.cm == nil {
		a.cm = &corev1.ConfigMap{
			ObjectMeta: ObjectMeta(),
			Data:       map[string]string{},
		}
	}
	return a
}

// AddAccount appends an IAM account to the `mapAccounts` entry
// in the configmap. It also deduplicates.
func (a *AuthConfigMap) AddAccount(account string) error {
	accounts, err := a.accounts()
	if err != nil {
		return err
	}
	distinct := map[string]struct{}{account: {}}
	for _, acc := range accounts {
		distinct[acc] = struct{}{}
	}
	accounts = accounts[:0]
	for acc := range distinct {
		accounts = append(accounts, acc)
	}
	// List order matters in yamls, maintain deterministic output
	sort.Strings(accounts)
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
		return fmt.Errorf("account %q not found in config map", account)
	}
	logger.Info("removing account %s from config map", account)
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

// AddRole appends a role with given groups.
func (a *AuthConfigMap) AddRole(arn string, groups []string) error {
	roles, err := a.roles()
	if err != nil {
		return err
	}
	roles = append(roles, mapRole{
		"rolearn":  arn,
		"username": "system:node:{{EC2PrivateDNSName}}",
		"groups":   groups,
	})
	return a.setRoles(roles)
}

// RemoveRole removes exactly one entry, even if there are duplicates.
// If it cannot find the role it returns an error.
func (a *AuthConfigMap) RemoveRole(arn string) error {
	if arn == "" {
		return errors.New("nodegroup instance role ARN is not set")
	}
	roles, err := a.roles()
	if err != nil {
		return err
	}

	for i, role := range roles {
		if role["rolearn"] == arn {
			logger.Info("removing role %s from config map", arn)
			roles = append(roles[:i], roles[i+1:]...)
			return a.setRoles(roles)
		}
	}

	return fmt.Errorf("instance role ARN %q not found in config map", arn)
}

func (a *AuthConfigMap) roles() (mapRoles, error) {
	var roles mapRoles
	if err := yaml.Unmarshal([]byte(a.cm.Data[rolesData]), &roles); err != nil {
		return nil, errors.Wrap(err, "unmarshalling mapRoles")
	}
	return roles, nil
}

func (a *AuthConfigMap) setRoles(r mapRoles) error {
	bs, err := yaml.Marshal(r)
	if err != nil {
		return errors.Wrap(err, "marshalling mapRoles")
	}
	a.cm.Data[rolesData] = string(bs)
	return nil
}

// Save persists the configmap to the cluster. It determines
// whether to create or update by looking at the configmap's UID.
func (a *AuthConfigMap) Save(client v1.ConfigMapInterface) (err error) {
	if a.cm.UID == "" {
		a.cm, err = client.Create(a.cm)
		return err
	}

	a.cm, err = client.Update(a.cm)
	return err
}

// ObjectMeta constructs metadata for the configmap
func ObjectMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      objectName,
		Namespace: objectNamespace,
	}
}

// AddNodeGroup creates or adds a nodegroup IAM role in the auth config map for the given nodegroup
func AddNodeGroup(clientSet kubernetes.Interface, ng *api.NodeGroup) error {
	client := clientSet.CoreV1().ConfigMaps(objectNamespace)

	// check if object exists
	cm, err := client.Get(objectName, metav1.GetOptions{})
	if err != nil && !kerr.IsNotFound(err) {
		// something must have gone terribly wrong
		return errors.Wrapf(err, "getting auth ConfigMap")
	}

	acm := New(cm)
	if err := acm.AddRole(ng.IAM.InstanceRoleARN, DefaultNodeGroups); err != nil {
		return errors.Wrap(err, "adding nodegroup to auth ConfigMap")
	}
	if err := acm.Save(client); err != nil {
		return errors.Wrap(err, "saving auth ConfigMap")
	}
	logger.Debug("saved auth ConfigMap for %s", ng.Name)
	return nil
}

// RemoveNodeGroup removes a nodegroup from the config map and does a client update
func RemoveNodeGroup(clientSet kubernetes.Interface, ng *api.NodeGroup) error {
	client := clientSet.CoreV1().ConfigMaps(objectNamespace)

	cm, err := client.Get(objectName, metav1.GetOptions{})
	if err != nil {
		return errors.Wrapf(err, "getting auth ConfigMap")
	}

	acm := New(cm)
	if err := acm.RemoveRole(ng.IAM.InstanceRoleARN); err != nil {
		return errors.Wrap(err, "removing nodegroup from auth ConfigMap")
	}
	if err := acm.Save(client); err != nil {
		return errors.Wrap(err, "updating auth ConfigMap after removing role")
	}
	logger.Debug("updated auth ConfigMap for %s", ng.Name)
	return nil
}
