package eks

import (
	"github.com/ghodss/yaml"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/api/meta/v1"
)

func (c *Config) nodeAuthConfigMap() (*corev1.ConfigMap, error) {

	/*
		apiVersion: v1
		kind: ConfigMap
		metadata:
		  name: aws-auth
		  namespace: default
		data:
		  mapRoles: |
		    - rolearn: "${nodeInstanceRoleARN}"
		      username: system:node:{{EC2PrivateDNSName}}
		      groups:
		        - system:bootstrappers
		        - system:nodes
		        - system:node-proxier
	*/

	mapRoles := make([]map[string]interface{}, 1)
	mapRoles[0] = make(map[string]interface{})

	mapRoles[0]["rolearn"] = c.nodeInstanceRoleARN
	mapRoles[0]["username"] = "system:node:{{EC2PrivateDNSName}}"
	mapRoles[0]["groups"] = []string{
		"system:bootstrappers",
		"system:nodes",
		"system:nodes",
	}

	mapRolesBytes, err := yaml.Marshal(mapRoles)
	if err != nil {
		return nil, err
	}

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "aws-auth",
			Namespace: "default",
		},
		BinaryData: map[string][]byte{
			"mapRoles": mapRolesBytes,
		},
	}

	return cm, nil
}

// def generate_sts_token(name):
//     sts = setupSTSBoto()
//     prefix = "k8s-aws-v1."

//     signedURL = sts.generate_presigned_url(ClientMethod='get_caller_identity',  Params={}, ExpiresIn=60)
//     encodedURL = base64.b64encode(signedURL)

//     return prefix+encodedURL```

// the issue with boto is it doesn't allow you to generate a signed url with additional headers like the golang package so I have to rewrite this to manually sign the url
