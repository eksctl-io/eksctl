package delete

// if we want to put it into package delete_test we would have to create another *_test.go file that exports
// the unexported function or make the function itself exportable

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/authconfigmap"
	"github.com/weaveworks/eksctl/pkg/cfn/manager"
	"github.com/weaveworks/eksctl/pkg/eks"
	"github.com/weaveworks/eksctl/pkg/testutils/mockprovider"
)

var (
	ng1 = &api.NodeGroup{
		Name: "ng1",
		IAM: &api.NodeGroupIAM{
			InstanceRoleARN:    "arn:aws:iam::122333:role/eksctl-cluster-ng1-NodeInstanceRole-ASDF",
			InstanceProfileARN: "arn:aws:iam::122333:instance-profile/eksctl-cluster-ng1-instance-profile",
		},
	}

	ng1b = &api.NodeGroup{
		Name: "ng1b",
		IAM: &api.NodeGroupIAM{
			InstanceRoleARN:    "arn:aws:iam::122333:role/eksctl-cluster-ng1-NodeInstanceRole-ASDF",
			InstanceProfileARN: "arn:aws:iam::122333:instance-profile/eksctl-cluster-ng1-instance-profile",
		},
	}

	ng2 = &api.NodeGroup{
		Name: "ng2",
		IAM: &api.NodeGroupIAM{
			InstanceRoleARN:    "arn:aws:iam::122333:role/eksctl-cluster-ng2-NodeInstanceRole-ZXCV",
			InstanceProfileARN: "arn:aws:iam::122333:instance-profile/eksctl-cluster-ng2-instance-profile",
		},
	}

	mng1 = &api.ManagedNodeGroup{
		Name: "mng1",
		IAM: &api.NodeGroupIAM{
			InstanceRoleARN:    "arn:aws:iam::122333:role/eksctl-cluster-mng1-NodeInstanceRole-BLAH",
			InstanceProfileARN: "arn:aws:iam::122333:instance-profile/eksctl-cluster-mng1-instance-profile",
		},
	}

	mng2 = &api.ManagedNodeGroup{
		Name: "mng2",
		IAM: &api.NodeGroupIAM{
			InstanceRoleARN:    "arn:aws:iam::122333:role/eksctl-cluster-ng1-NodeInstanceRole-ASDF",
			InstanceProfileARN: "arn:aws:iam::122333:instance-profile/eksctl-cluster-mng1-instance-profile",
		},
	}

	//fg1 = &api.FargateProfile{
	//	Name:                "fg1",
	//	PodExecutionRoleARN: "arn:aws:iam::122333:role/eksctl-cluster-ServiceRole-HELLOTHERE",
	//}
)

const InstanceRoleFmt = `
- groups:
  - system:bootstrappers
  - system:nodes
  rolearn: %s
  username: system:node:{{EC2PrivateDNSName}}
`

const ServiceRoleFmt = `
- groups:
  - system:bootstrappers
  - system:nodes
  - system:node-proxier
  rolearn: %s
  username: system:node:{{SessionName}}
`

const ManagedNodeGroupTemplateBodyFmt = `
{
	"Resources": {
		"ManagedNodeGroup": {
			"Properties": {
				"InstanceTypes": ["t3.micro"],
				"DesiredCapacity": "1",
				"LaunchTemplate": {
					"LaunchTemplateName": "%s"
				},
				"MaxSize": "1",
				"MinSize": "1",
				"Tags": [
					{
						"Key": "Name",
						"PropagateAtLaunch": "true",
						"Value": "%s-Node"
					},
					{
						"Key": "kubernetes.io/cluster/%s",
						"PropagateAtLaunch": "true",
						"Value": "owned"
					}
				]
			}
		}
	},
	"Outputs": {
	    "InstanceRoleARN": {
			"Value": "%s"
    	}
  	}
}
`

const UnmanagedNodeGroupTemplateBodyFmt = `
{
	"Resources": {
		"NodeGroup": {
			"Type": "AWS::AutoScaling::AutoScalingGroup",
			"Properties": {
				"DesiredCapacity": "1",
				"LaunchTemplate": {
					"LaunchTemplateName": "%s"
				},
				"MaxSize": "1",
				"MinSize": "1",
				"Tags": [
					{
						"Key": "Name",
						"PropagateAtLaunch": "true",
						"Value": "%s-Node"
					},
					{
						"Key": "kubernetes.io/cluster/%s",
						"PropagateAtLaunch": "true",
						"Value": "owned"
					}
				]
			}
		},
		"NodeGroupLaunchTemplate": {
			"Type": "AWS::EC2::LaunchTemplate",
			"Properties": {
				"LaunchTemplateData": {
					"InstanceType": "t3.micro",
				},
				"LaunchTemplateName": "%s"
			}
		}
	},
	"Outputs": {
	    "InstanceRoleARN": {
			"Value": "%s"
    	}
  	}
}
`

func createTemplateBody(ngType api.NodeGroupType, stackName, clusterName, roleARN string) string {
	templateBodyFmt := ""
	if ngType == api.NodeGroupTypeUnmanaged {
		templateBodyFmt = UnmanagedNodeGroupTemplateBodyFmt
	} else if ngType == api.NodeGroupTypeManaged {
		templateBodyFmt = ManagedNodeGroupTemplateBodyFmt
	}
	return fmt.Sprintf(templateBodyFmt, stackName, clusterName, clusterName, stackName, roleARN)
}

func createAuthConfigMap(resources ...interface{}) (kubernetes.Interface, *v1.ConfigMap) {
	clientSet := fake.NewSimpleClientset()
	mapRoles := []string{} // we use the literal declaration as we want "[]"
	var mapRole string
	for _, r := range resources {
		switch r.(type) {
		case *api.NodeGroup:
			mapRole = fmt.Sprintf(InstanceRoleFmt, r.(*api.NodeGroup).IAM.InstanceRoleARN)
		case *api.ManagedNodeGroup:
			mapRole = fmt.Sprintf(InstanceRoleFmt, r.(*api.ManagedNodeGroup).IAM.InstanceRoleARN)
		case *api.FargateProfile:
			mapRole = fmt.Sprintf(InstanceRoleFmt, r.(*api.FargateProfile).PodExecutionRoleARN)
		default:
			continue
		}
		mapRoles = append(mapRoles, mapRole)
	}
	acmData := map[string]string{
		"mapRoles": strings.Join(mapRoles, "\n"),
	}
	acm := &v1.ConfigMap{
		ObjectMeta: authconfigmap.ObjectMeta(),
		Data:       acmData,
	}
	acm.UID = "12345" // required to set updates and is not provided to us by the fake clientset
	acm, _ = clientSet.CoreV1().ConfigMaps(authconfigmap.ObjectNamespace).Create(acm)
	return clientSet, acm
}

var _ = Describe("removeARN()", func() {

	p := mockprovider.NewMockProvider()

	mockCFN := p.MockCloudFormation()

	mockCFN.On("GetTemplate",
		mock.MatchedBy(func(input *cloudformation.GetTemplateInput) bool {
			return input.StackName != nil && *input.StackName == "ng1"
		})).
		Return(&cloudformation.GetTemplateOutput{
			TemplateBody: aws.String(createTemplateBody(api.NodeGroupTypeUnmanaged, "ng1", "my-cluster", ng1.IAM.InstanceRoleARN)),
		}, nil)

	mockCFN.On("GetTemplate",
		mock.MatchedBy(func(input *cloudformation.GetTemplateInput) bool {
			return input.StackName != nil && *input.StackName == "ng1b"
		})).
		Return(&cloudformation.GetTemplateOutput{
			TemplateBody: aws.String(createTemplateBody(api.NodeGroupTypeUnmanaged, "ng1b", "my-cluster", ng1b.IAM.InstanceRoleARN)),
		}, nil)

	mockCFN.On("GetTemplate",
		mock.MatchedBy(func(input *cloudformation.GetTemplateInput) bool {
			return input.StackName != nil && *input.StackName == "ng2"
		})).
		Return(&cloudformation.GetTemplateOutput{
			TemplateBody: aws.String(createTemplateBody(api.NodeGroupTypeUnmanaged, "ng2", "my-cluster", ng2.IAM.InstanceRoleARN)),
		}, nil)

	mockCFN.On("GetTemplate",
		mock.MatchedBy(func(input *cloudformation.GetTemplateInput) bool {
			return input.StackName != nil && *input.StackName == "mng1"
		})).
		Return(&cloudformation.GetTemplateOutput{
			TemplateBody: aws.String(createTemplateBody(api.NodeGroupTypeManaged, "mng1", "my-cluster", ng1.IAM.InstanceRoleARN)),
		}, nil)

	ctl := &eks.ClusterProvider{
		Provider: p,
		Status:   &eks.ProviderStatus{},
	}

	cfg := &api.ClusterConfig{
		Metadata: &api.ClusterMeta{
			Name:   "my-cluster",
			Region: "us-east-1",
		},
	}

	cfnDescriptions := []*manager.Stack{
		{
			StackName: aws.String(ng1.Name),
			RoleARN:   aws.String(ng1.IAM.InstanceRoleARN),
			Tags: []*cloudformation.Tag{
				{Key: aws.String(api.NodeGroupNameTag), Value: aws.String(ng1.Name)},
				{Key: aws.String(api.NodeGroupTypeTag), Value: aws.String(string(api.NodeGroupTypeUnmanaged))},
			},
			Outputs: []*cloudformation.Output{
				{OutputKey: aws.String("InstanceRoleARN"), OutputValue: aws.String(ng1.IAM.InstanceRoleARN)},
			},
		},
		{
			StackName: aws.String(ng1b.Name),
			RoleARN:   aws.String(ng1b.IAM.InstanceRoleARN),
			Tags: []*cloudformation.Tag{
				{Key: aws.String(api.NodeGroupNameTag), Value: aws.String(ng1b.Name)},
				{Key: aws.String(api.NodeGroupTypeTag), Value: aws.String(string(api.NodeGroupTypeUnmanaged))},
			},
			Outputs: []*cloudformation.Output{
				{OutputKey: aws.String("InstanceRoleARN"), OutputValue: aws.String(ng1b.IAM.InstanceRoleARN)},
			},
		},
		{
			StackName: aws.String(ng2.Name),
			RoleARN:   aws.String(ng2.IAM.InstanceRoleARN),
			Tags: []*cloudformation.Tag{
				{Key: aws.String(api.NodeGroupNameTag), Value: aws.String(ng2.Name)},
				{Key: aws.String(api.NodeGroupTypeTag), Value: aws.String(string(api.NodeGroupTypeUnmanaged))},
			},
			Outputs: []*cloudformation.Output{
				{OutputKey: aws.String("InstanceRoleARN"), OutputValue: aws.String(ng2.IAM.InstanceRoleARN)},
			},
		},
		{
			StackName: aws.String(mng1.Name),
			RoleARN:   aws.String(mng1.IAM.InstanceRoleARN),
			Tags: []*cloudformation.Tag{
				{Key: aws.String(api.NodeGroupNameTag), Value: aws.String(mng1.Name)},
				{Key: aws.String(api.NodeGroupTypeTag), Value: aws.String(string(api.NodeGroupTypeManaged))},
			},
			Outputs: []*cloudformation.Output{
				{OutputKey: aws.String("InstanceRoleARN"), OutputValue: aws.String(mng1.IAM.InstanceRoleARN)},
			},
		},
	}

	stackManager := manager.NewStackCollection(p, cfg)

	It("should remove identity from auth configmap via name when removing only node group", func() {
		cfg.NodeGroups = []*api.NodeGroup{ng1} // to mark for deletion
		clientSet, acm := createAuthConfigMap(ng1)
		Expect(acm.Data["mapRoles"]).NotTo(BeEmpty())

		err := removeARN(cfnDescriptions, stackManager, cfg, ctl, false, clientSet)
		Expect(err).NotTo(HaveOccurred())

		acm, _ = clientSet.CoreV1().ConfigMaps("kube-system").Get("aws-auth", metav1.GetOptions{})
		Expect(acm.Data["mapRoles"]).To(Equal("[]\n"))
	})

	It("should remove identity from auth configmap when no one else is using it", func() {
		cfg.NodeGroups = []*api.NodeGroup{ng1} // to mark for deletion
		clientSet, acm := createAuthConfigMap(ng1, ng2)
		Expect(acm.Data["mapRoles"]).NotTo(BeEmpty())

		err := removeARN(cfnDescriptions, stackManager, cfg, ctl, false, clientSet)
		Expect(err).NotTo(HaveOccurred())

		acm, _ = clientSet.CoreV1().ConfigMaps("kube-system").Get("aws-auth", metav1.GetOptions{})
		Expect(acm.Data["mapRoles"]).NotTo(Equal("[]\n"))
	})

	It("should only remove one identity entry from auth configmap", func() {
		cfg.NodeGroups = []*api.NodeGroup{ng1} // to mark for deletion
		clientSet, acm := createAuthConfigMap(ng1, ng1b)
		Expect(acm.Data["mapRoles"]).NotTo(BeEmpty())

		err := removeARN(cfnDescriptions, stackManager, cfg, ctl, false, clientSet)
		Expect(err).NotTo(HaveOccurred())

		acm, _ = clientSet.CoreV1().ConfigMaps("kube-system").Get("aws-auth", metav1.GetOptions{})
		Expect(acm.Data["mapRoles"]).NotTo(Equal("[]\n"))
	})

	It("should only remove one identity entry from auth configmap when deleting managed node group", func() {
		cfg.NodeGroups = []*api.NodeGroup{}
		cfg.ManagedNodeGroups = []*api.ManagedNodeGroup{mng2} // to mark for deletion
		clientSet, acm := createAuthConfigMap(mng2, mng1, ng1, ng1b, ng2)
		Expect(acm.Data["mapRoles"]).NotTo(BeEmpty())

		err := removeARN(cfnDescriptions, stackManager, cfg, ctl, false, clientSet)
		Expect(err).NotTo(HaveOccurred())

		acm, _ = clientSet.CoreV1().ConfigMaps("kube-system").Get("aws-auth", metav1.GetOptions{})

		roles := []string{
			fmt.Sprintf(InstanceRoleFmt, mng1.IAM.InstanceRoleARN),
			fmt.Sprintf(InstanceRoleFmt, ng1.IAM.InstanceRoleARN),
			fmt.Sprintf(InstanceRoleFmt, ng1b.IAM.InstanceRoleARN),
			fmt.Sprintf(InstanceRoleFmt, ng2.IAM.InstanceRoleARN),
		}

		Expect(acm.Data["mapRoles"]).To(MatchYAML(strings.Join(roles, "\n")))
	})
})
