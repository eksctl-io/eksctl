package utils

import (
	"fmt"

	"github.com/kris-nova/logger"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type mapRolesData []map[string]interface{}

const (
	// AuthConfigMapName is the name of the ConfigMap
	AuthConfigMapName = "aws-auth"
	// AuthConfigMapNamespace is the namespace of the ConfigMap
	AuthConfigMapNamespace = "kube-system"
)

func makeMapRolesData() mapRolesData { return []map[string]interface{}{} }

func appendNodeGroupToAuthConfigMap(mapRoles *mapRolesData, ngInstanceRoleARN string) {
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

func newAuthConfigMap(mapRoles mapRolesData) (*corev1.ConfigMap, error) {
	mapRolesBytes, err := yaml.Marshal(mapRoles)
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

// NewAuthConfigMap creates ConfigMap with a single nodegroup ARN
func NewAuthConfigMap(ngInstanceRoleARN string) (*corev1.ConfigMap, error) {
	mapRoles := makeMapRolesData()
	appendNodeGroupToAuthConfigMap(&mapRoles, ngInstanceRoleARN)
	return newAuthConfigMap(mapRoles)
}

func updateAuthConfigMap(cm *corev1.ConfigMap, mapRoles mapRolesData) error {
	mapRolesBytes, err := yaml.Marshal(mapRoles)
	if err != nil {
		return err
	}
	cm.Data["mapRoles"] = string(mapRolesBytes)
	return nil
}

// UpdateAuthConfigMap updates ConfigMap by appending a single nodegroup ARN
func UpdateAuthConfigMap(cm *corev1.ConfigMap, ngInstanceRoleARN string) error {
	mapRoles := makeMapRolesData()
	if err := yaml.Unmarshal([]byte(cm.Data["mapRoles"]), &mapRoles); err != nil {
		return err
	}
	appendNodeGroupToAuthConfigMap(&mapRoles, ngInstanceRoleARN)
	return updateAuthConfigMap(cm, mapRoles)
}

// RemoveNodeGroupFromAuthConfigMap removes a node group's instance mapped role from the ConfigMap
func RemoveNodeGroupFromAuthConfigMap(cm *corev1.ConfigMap, ngInstanceRoleARN string) error {
	if ngInstanceRoleARN == "" {
		return errors.New("config map is unchanged as the node group instance ARN is not set")
	}
	mapRoles := makeMapRolesData()
	found := false
	var mapRolesUpdated mapRolesData
	if err := yaml.Unmarshal([]byte(cm.Data["mapRoles"]), &mapRoles); err != nil {
		return err
	}

	for _, role := range mapRoles {
		if role["rolearn"] == ngInstanceRoleARN {
			found = true
			logger.Info("removing %s from config map", ngInstanceRoleARN)
		} else {
			mapRolesUpdated = append(mapRolesUpdated, role)
		}
	}

	if !found {
		return fmt.Errorf("instance role ARN %s not found in config map", ngInstanceRoleARN)
	}

	return updateAuthConfigMap(cm, mapRolesUpdated)
}
