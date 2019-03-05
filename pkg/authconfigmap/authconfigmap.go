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
	"k8s.io/client-go/kubernetes"
)

type mapRolesData []map[string]interface{}

const (
	objectName      = "aws-auth"
	objectNamespace = metav1.NamespaceSystem
)

// ObjectMeta constructs metadata for the configmap
func ObjectMeta() metav1.ObjectMeta {
	return metav1.ObjectMeta{
		Name:      objectName,
		Namespace: objectNamespace,
	}
}

func makeMapRolesData() mapRolesData { return []map[string]interface{}{} }

func appendNodeRole(mapRoles *mapRolesData, ngInstanceRoleARN string) {
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
		ObjectMeta: ObjectMeta(),
		Data: map[string]string{
			"mapRoles": string(mapRolesBytes),
		},
	}
	return cm, nil
}

func update(cm *corev1.ConfigMap, mapRoles *mapRolesData) error {
	mapRolesBytes, err := yaml.Marshal(*mapRoles)
	if err != nil {
		return err
	}
	cm.Data["mapRoles"] = string(mapRolesBytes)
	return nil
}

// NewForRole creates ConfigMap with a single role ARN
func NewForRole(arn string) (*corev1.ConfigMap, error) {
	mapRoles := makeMapRolesData()
	appendNodeRole(&mapRoles, arn)
	return new(&mapRoles)
}

// AddRole updates ConfigMap by appending a single nodegroup ARN
func AddRole(cm *corev1.ConfigMap, arn string) error {
	mapRoles := makeMapRolesData()
	if err := yaml.Unmarshal([]byte(cm.Data["mapRoles"]), &mapRoles); err != nil {
		return err
	}
	appendNodeRole(&mapRoles, arn)
	return update(cm, &mapRoles)
}

// AddNodeGroup creates or adds a nodegroup IAM role in the auth config map for the given nodegroup
func AddNodeGroup(clientSet kubernetes.Interface, ng *api.NodeGroup) error {
	cm := &corev1.ConfigMap{}
	client := clientSet.CoreV1().ConfigMaps(objectNamespace)
	create := false

	// check if object exists
	if existing, err := client.Get(objectName, metav1.GetOptions{}); err != nil {
		if kerr.IsNotFound(err) {
			create = true // doesn't exsits, will create
		} else {
			// something must have gone terribly wrong
			return errors.Wrapf(err, "getting auth ConfigMap")
		}
	} else {
		*cm = *existing // use existing object
	}

	if create {
		// build new object with the given role
		cm, err := NewForRole(ng.IAM.InstanceRoleARN)
		if err != nil {
			return errors.Wrap(err, "constructing auth ConfigMap")
		}
		// and create it in the cluster
		if _, err := client.Create(cm); err != nil {
			return errors.Wrap(err, "creating auth ConfigMap")
		}
		logger.Debug("created auth ConfigMap for %s", ng.Name)
		return nil
	}

	// in case we already have an onject, and the given role to it
	if err := AddRole(cm, ng.IAM.InstanceRoleARN); err != nil {
		return errors.Wrap(err, "creating an update for auth ConfigMap")
	}
	// and update it in the cluster
	if _, err := client.Update(cm); err != nil {
		return errors.Wrap(err, "updating auth ConfigMap")
	}
	logger.Debug("updated auth ConfigMap for %s", ng.Name)
	return nil
}

func doRemoveRole(mapRoles *mapRolesData, arn string) (mapRolesData, error) {
	found := false
	mapRolesUpdated := makeMapRolesData()

	for _, role := range *mapRoles {
		if role["rolearn"] == arn && !found {
			found = true
			logger.Info("removing %s from config map", arn)
		} else {
			mapRolesUpdated = append(mapRolesUpdated, role)
		}
	}

	if !found {
		return nil, fmt.Errorf("instance role ARN %s not found in config map", arn)
	}

	return mapRolesUpdated, nil
}

// RemoveRole removes a nodegroup's instance mapped role from the config map
func RemoveRole(cm *corev1.ConfigMap, ngInstanceRoleARN string) error {
	if ngInstanceRoleARN == "" {
		return errors.New("config map is unchanged as the nodegroup instance ARN is not set")
	}
	mapRoles := makeMapRolesData()

	if err := yaml.Unmarshal([]byte(cm.Data["mapRoles"]), &mapRoles); err != nil {
		return err
	}
	mapRolesUpdated, err := doRemoveRole(&mapRoles, ngInstanceRoleARN)
	if err != nil {
		return err
	}
	return update(cm, &mapRolesUpdated)
}

// RemoveNodeGroup removes a nodegroup from the config map and does a client update
func RemoveNodeGroup(clientSet kubernetes.Interface, ng *api.NodeGroup) error {
	client := clientSet.CoreV1().ConfigMaps(objectNamespace)

	cm, err := client.Get(objectName, metav1.GetOptions{})
	if err != nil {
		return errors.Wrapf(err, "getting auth ConfigMap")
	}

	if err := RemoveRole(cm, ng.IAM.InstanceRoleARN); err != nil {
		return errors.Wrapf(err, "removing nodegroup from auth ConfigMap")
	}
	if _, err := client.Update(cm); err != nil {
		return errors.Wrap(err, "updating auth ConfigMap and removing instance role")
	}
	logger.Debug("updated auth ConfigMap for %s", ng.Name)
	return nil
}
