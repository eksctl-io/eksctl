package clusterautoscaler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	. "github.com/weaveworks/eksctl/pkg/addons/clusterautoscaler"
)

func TestValidate(t *testing.T) {
	invalidParams := TemplateParameters{}
	assert.EqualError(t, invalidParams.Validate(), "blank cluster name")
	assert.Equal(t, "", invalidParams.ClusterName)
	assert.Equal(t, "", invalidParams.Namespace)
	assert.Equal(t, "", invalidParams.ImageVersion)

	params := TemplateParameters{ClusterName: "cluster-name-test"}
	assert.NoError(t, params.Validate())
	assert.Equal(t, "cluster-name-test", params.ClusterName)
	assert.Equal(t, "kube-system", params.Namespace)
	assert.Equal(t, "v1.12.3", params.ImageVersion)
}

func TestGenerateManifests(t *testing.T) {
	invalidManifest, err := GenerateManifests(TemplateParameters{})
	assert.EqualError(t, err, "blank cluster name")
	assert.Len(t, invalidManifest, 0)

	defaultManifest, err := GenerateManifests(TemplateParameters{ClusterName: "cluster-name-test"})
	assert.NoError(t, err)
	manifest := string(defaultManifest)
	assert.Contains(t, manifest, "  namespace: kube-system")
	assert.Contains(t, manifest, "        - image: k8s.gcr.io/cluster-autoscaler:v1.12.3")
	assert.Contains(t, manifest, "            - --node-group-auto-discovery=asg:tag=k8s.io/cluster-autoscaler/enabled,k8s.io/cluster-autoscaler/cluster-name-test")

	manifestBytes, err := GenerateManifests(TemplateParameters{
		ClusterName:  "cluster-name2-test",
		ImageVersion: "v2.0.0",
		Namespace:    "custom-ns",
	})
	assert.NoError(t, err)
	manifest = string(manifestBytes)
	assert.NotContains(t, manifest, "  namespace: kube-system")
	assert.Contains(t, manifest, "  namespace: custom-ns")
	assert.NotContains(t, manifest, "        - image: k8s.gcr.io/cluster-autoscaler:v1.12.3")
	assert.Contains(t, manifest, "        - image: k8s.gcr.io/cluster-autoscaler:v2.0.0")
	assert.NotContains(t, manifest, "            - --node-group-auto-discovery=asg:tag=k8s.io/cluster-autoscaler/enabled,k8s.io/cluster-autoscaler/cluster-name-test")
	assert.Contains(t, manifest, "            - --node-group-auto-discovery=asg:tag=k8s.io/cluster-autoscaler/enabled,k8s.io/cluster-autoscaler/cluster-name2-test")
}
