package authconfigmap

import (
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha4"

	yaml "gopkg.in/yaml.v2"

	corev1 "k8s.io/api/core/v1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
)

type mapRolesData []map[string]interface{}

const (
	// AuthConfigMapName is the name of the ConfigMap
	AuthConfigMapName = "aws-auth"
	// AuthConfigMapNamespace is the namespace of the ConfigMap
	AuthConfigMapNamespace = "kube-system"
)

func makeMapRolesData() mapRolesData { return []map[string]interface{}{} }

func appendNodeRoleToMapRoles(mapRoles *mapRolesData, ngInstanceRoleARN string) {
	newEntry := map[string]interface{}{
		"rolearn":  ngInstanceRoleARN,
		"username": "system:node:{{EC2PrivateDNSName}}",
		"groups": []string{
			"system:bootstrappers",
			"system:nodes",
		},
	}
	*mapRoles = append(*mapRoles, newEntry)
}

func new(mapRoles *mapRolesData) (*corev1.ConfigMap, error) {
	mapRolesBytes, err := yaml.Marshal(*mapRoles)
	if err != nil {
		return nil, err
	}
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      AuthConfigMapName,
			Namespace: AuthConfigMapNamespace,
		},
		Data: map[string]string{
			"mapRoles": string(mapRolesBytes),
		},
	}
	return cm, nil
}

// New creates ConfigMap with a single nodegroup ARN
func New(ngInstanceRoleARN string) (*corev1.ConfigMap, error) {
	mapRoles := makeMapRolesData()
	appendNodeRoleToMapRoles(&mapRoles, ngInstanceRoleARN)
	return new(&mapRoles)
}

func updateAuthConfigMap(cm *corev1.ConfigMap, mapRoles *mapRolesData) error {
	mapRolesBytes, err := yaml.Marshal(*mapRoles)
	if err != nil {
		return err
	}
	cm.Data["mapRoles"] = string(mapRolesBytes)
	return nil
}

// AddNodeRoleToAuthConfigMap updates ConfigMap by appending a single nodegroup ARN
func AddNodeRoleToAuthConfigMap(cm *corev1.ConfigMap, ngInstanceRoleARN string) error {
	mapRoles := makeMapRolesData()
	if err := yaml.Unmarshal([]byte(cm.Data["mapRoles"]), &mapRoles); err != nil {
		return err
	}
	appendNodeRoleToMapRoles(&mapRoles, ngInstanceRoleARN)
	return updateAuthConfigMap(cm, &mapRoles)
}

// CreateOrAddNodeGroupRole creates or adds a node group IAM role the auth config map for the given nodegroup
func CreateOrAddNodeGroupRole(clientSet *clientset.Clientset, ng *api.NodeGroup) error {
	cm := &corev1.ConfigMap{}
	client := clientSet.CoreV1().ConfigMaps(AuthConfigMapNamespace)
	create := false

	if existing, err := client.Get(AuthConfigMapName, metav1.GetOptions{}); err != nil {
		if kerr.IsNotFound(err) {
			create = true
		} else {
			return errors.Wrapf(err, "getting auth ConfigMap")
		}
	} else {
		*cm = *existing
	}

	if create {
		cm, err := New(ng.IAM.InstanceRoleARN)
		if err != nil {
			return errors.Wrap(err, "constructing auth ConfigMap")
		}
		if _, err := client.Create(cm); err != nil {
			return errors.Wrap(err, "creating auth ConfigMap")
		}
		logger.Debug("created auth ConfigMap for %s", ng.Name)
		return nil
	}

	if err := AddNodeRoleToAuthConfigMap(cm, ng.IAM.InstanceRoleARN); err != nil {
		return errors.Wrap(err, "creating an update for auth ConfigMap")
	}
	if _, err := client.Update(cm); err != nil {
		return errors.Wrap(err, "updating auth ConfigMap")
	}
	logger.Debug("updated auth ConfigMap for %s", ng.Name)
	return nil
}

func removeNodeRoleFromMapRoles(mapRoles *mapRolesData, ngInstanceRoleARN string) (mapRolesData, error) {
	found := false
	mapRolesUpdated := makeMapRolesData()

	for _, role := range *mapRoles {
		if role["rolearn"] == ngInstanceRoleARN && !found {
			found = true
			logger.Info("removing %s from config map", ngInstanceRoleARN)
		} else {
			mapRolesUpdated = append(mapRolesUpdated, role)
		}
	}

	if !found {
		return nil, fmt.Errorf("instance role ARN %s not found in config map", ngInstanceRoleARN)
	}

	return mapRolesUpdated, nil
}

// RemoveNodeRoleFromAuthConfigMap removes a node group's instance mapped role from the config map
func RemoveNodeRoleFromAuthConfigMap(cm *corev1.ConfigMap, ngInstanceRoleARN string) error {
	if ngInstanceRoleARN == "" {
		return errors.New("config map is unchanged as the node group instance ARN is not set")
	}
	mapRoles := makeMapRolesData()

	if err := yaml.Unmarshal([]byte(cm.Data["mapRoles"]), &mapRoles); err != nil {
		return err
	}
	mapRolesUpdated, err := removeNodeRoleFromMapRoles(&mapRoles, ngInstanceRoleARN)
	if err != nil {
		return err
	}
	return updateAuthConfigMap(cm, &mapRolesUpdated)
}

// RemoveNodeGroupRole removes a nodegroup from the config map and does a client update
func RemoveNodeGroupRole(clientSet *clientset.Clientset, ng *api.NodeGroup) error {
	client := clientSet.CoreV1().ConfigMaps(AuthConfigMapNamespace)

	cm, err := client.Get(AuthConfigMapName, metav1.GetOptions{})
	if err != nil {
		return errors.Wrapf(err, "getting auth ConfigMap")
	}

	if err := RemoveNodeRoleFromAuthConfigMap(cm, ng.IAM.InstanceRoleARN); err != nil {
		return errors.Wrapf(err, "removing node group from auth ConfigMap")
	}
	if _, err := client.Update(cm); err != nil {
		return errors.Wrap(err, "updating auth ConfigMap and removing instance role")
	}
	logger.Debug("updated auth ConfigMap for %s", ng.Name)
	return nil
}
