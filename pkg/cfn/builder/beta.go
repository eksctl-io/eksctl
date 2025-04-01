package builder

import (
	"context"
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/weaveworks/eksctl/pkg/awsapi"
	cft "github.com/weaveworks/eksctl/pkg/cfn/template"
	"github.com/weaveworks/eksctl/pkg/goformation"
	gfn "github.com/weaveworks/eksctl/pkg/goformation/cloudformation"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/cloudformation"
	gfneks "github.com/weaveworks/eksctl/pkg/goformation/cloudformation/eks"
	"github.com/weaveworks/eksctl/pkg/goformation/cloudformation/lambda"
	gfnt "github.com/weaveworks/eksctl/pkg/goformation/cloudformation/types"
)

//go:embed templates/beta-resources.yaml
var betaResourcesTemplate []byte

//go:embed templates/beta.py
var lambdaBetaPy []byte

func addBetaResources(stsAPI awsapi.STS, stackName string, clusterTemplate *gfn.Template, g *gfneks.Cluster) error {

	identity, err := stsAPI.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return fmt.Errorf("unable to get identity: %w", err)
	}
	userArn := *identity.Arn
	baseArn := userArn[:strings.LastIndex(userArn, "/")]
	roleArn := fmt.Sprintf("%s%s", baseArn, "/{{SessionName}}")
	iamARN := strings.Replace(
		strings.Replace(baseArn, "assumed-role", "role", 1),
		"sts", "iam", 1)

	clusterName := "eksctl-" + stackName + "-cluster"

	template, err := goformation.ParseYAML(betaResourcesTemplate)
	if err != nil {
		return err
	}
	for resourceName, resource := range template.Resources {
		clusterTemplate.Resources[resourceName] = resource
	}
	for key, output := range template.Outputs {
		clusterTemplate.Outputs[key] = output
	}
	customResource := clusterTemplate.Resources["ControlPlane"].(*gfn.CustomResource)
	if g.AccessConfig != nil {
		customResource.Properties["AccessConfig"] = g.AccessConfig
	}
	if g.BootstrapSelfManagedAddons != nil {
		customResource.Properties["BootstrapSelfManagedAddons"] = g.BootstrapSelfManagedAddons
	}
	if g.ComputeConfig != nil {
		customResource.Properties["ComputeConfig"] = g.ComputeConfig
	}
	if g.EncryptionConfig != nil {
		customResource.Properties["EncryptionConfig"] = g.EncryptionConfig
	}
	if g.KubernetesNetworkConfig != nil {
		customResource.Properties["KubernetesNetworkConfig"] = g.KubernetesNetworkConfig
	}
	if g.Logging != nil {
		customResource.Properties["Logging"] = g.Logging
	}
	if g.Name != nil {
		customResource.Properties["Name"] = g.Name
	}
	if g.OutpostConfig != nil {
		customResource.Properties["OutpostConfig"] = g.OutpostConfig
	}
	if g.RemoteNetworkConfig != nil {
		customResource.Properties["RemoteNetworkConfig"] = g.RemoteNetworkConfig
	}
	if g.ResourcesVpcConfig != nil {
		customResource.Properties["ResourcesVpcConfig"] = g.ResourcesVpcConfig
	}
	if g.RoleArn != nil {
		customResource.Properties["RoleArn"] = g.RoleArn
	}
	if g.StorageConfig != nil {
		customResource.Properties["StorageConfig"] = g.StorageConfig
	}
	if g.Tags != nil {
		g.Tags = append(g.Tags, cloudformation.Tag{
			Key:   gfnt.NewString("Name"),
			Value: gfnt.NewString(clusterName + "/ControlPlane"),
		})
		customResource.Properties["Tags"] = g.Tags
	} else {
		customResource.Properties["Tags"] = []cloudformation.Tag{
			{
				Key:   gfnt.NewString("Name"),
				Value: gfnt.NewString(clusterName + "/ControlPlane"),
			},
		}
	}
	if g.UpgradePolicy != nil {
		customResource.Properties["UpgradePolicy"] = g.UpgradePolicy
	}
	if g.Version != nil {
		customResource.Properties["Version"] = g.Version
	}
	if g.ZonalShiftConfig != nil {
		customResource.Properties["ZonalShiftConfig"] = g.ZonalShiftConfig
	}

	customResource.Properties["IAMPrincipalArn"] = gfnt.NewString(iamARN)
	customResource.Properties["STSRoleArn"] = gfnt.NewString(roleArn)

	customFunction := clusterTemplate.Resources["CustomEKSFunction"].(*lambda.Function)
	customFunction.Code = &lambda.Function_Code{
		ZipFile: gfnt.NewString(string(lambdaBetaPy)),
	}

	clusterTemplate.Outputs["EKSFunctionArn"] = gfn.Output{
		Value: gfnt.MakeFnGetAttString("CustomEKSFunction", "Arn"),
		Export: &gfn.Export{
			Name: gfnt.MakeFnSubString(fmt.Sprintf("${%s}::EKSFunctionArn", gfnt.StackName)),
		},
	}

	clusterTemplate.Parameters["EksEndpointUrl"] = gfn.Parameter{
		Type:        "String",
		Description: "The endpoint URL for the EKS service",
		Default:     gfnt.NewString(os.Getenv("AWS_ENDPOINT_URL_EKS")),
	}
	return nil
}

func addBetaManagedNodeGroupResources(managedResource *gfneks.Nodegroup, stackName string) *gfn.CustomResource {
	customResource := &gfn.CustomResource{
		Type: "Custom::EksManagedNodeGroup",
	}
	customResource.Properties = make(map[string]interface{})
	functionArn := gfnt.MakeFnImportValueString(fmt.Sprintf("eksctl-%s-cluster::EKSFunctionArn", stackName))
	customResource.Properties["ServiceToken"] = functionArn

	if managedResource.AmiType != nil {
		customResource.Properties["AmiType"] = managedResource.AmiType
	}
	if managedResource.CapacityType != nil {
		customResource.Properties["CapacityType"] = managedResource.CapacityType
	}
	if managedResource.ClusterName != nil {
		customResource.Properties["ClusterName"] = managedResource.ClusterName
	}
	if managedResource.DiskSize != nil {
		customResource.Properties["DiskSize"] = managedResource.DiskSize
	}
	if managedResource.ForceUpdateEnabled != nil {
		customResource.Properties["ForceUpdateEnabled"] = managedResource.ForceUpdateEnabled
	}
	if managedResource.InstanceTypes != nil {
		customResource.Properties["InstanceTypes"] = managedResource.InstanceTypes
	}
	if managedResource.Labels != nil {
		customResource.Properties["Labels"] = managedResource.Labels
	}
	if managedResource.LaunchTemplate != nil {
		customResource.Properties["LaunchTemplate"] = managedResource.LaunchTemplate
	}
	if managedResource.NodeRepairConfig != nil {
		customResource.Properties["NodeRepairConfig"] = managedResource.NodeRepairConfig
	}
	if managedResource.NodeRole != nil {
		customResource.Properties["NodeRole"] = managedResource.NodeRole
	}
	if managedResource.NodegroupName != nil {
		customResource.Properties["NodegroupName"] = managedResource.NodegroupName
	}
	if managedResource.ReleaseVersion != nil {
		customResource.Properties["ReleaseVersion"] = managedResource.ReleaseVersion
	}
	if managedResource.RemoteAccess != nil {
		customResource.Properties["RemoteAccess"] = managedResource.RemoteAccess
	}
	if managedResource.ScalingConfig != nil {
		customResource.Properties["ScalingConfig"] = managedResource.ScalingConfig
	}
	if managedResource.Subnets != nil {
		customResource.Properties["Subnets"] = managedResource.Subnets
	}
	if managedResource.Tags != nil {
		customResource.Properties["Tags"] = managedResource.Tags
	}
	if managedResource.Taints != nil {
		customResource.Properties["Taints"] = managedResource.Taints
	}
	if managedResource.UpdateConfig != nil {
		customResource.Properties["UpdateConfig"] = managedResource.UpdateConfig
	}
	if managedResource.Version != nil {
		customResource.Properties["Version"] = managedResource.Version
	}

	return customResource
}

func createBetaAssumeRolePolicy() interface{} {
	statements := []cft.MapOfInterfaces{
		{
			"Effect": "Allow",
			"Principal": cft.MapOfInterfaces{
				"Service": "eks.amazonaws.com",
			},
			"Action": []string{
				"sts:AssumeRole",
				"sts:TagSession",
			},
		},
		{
			"Effect": "Allow",
			"Principal": cft.MapOfInterfaces{
				"Service": "eks-beta-pdx.aws.internal",
			},
			"Action": []string{
				"sts:AssumeRole",
				"sts:TagSession",
			},
		},
		{
			"Effect": "Allow",
			"Principal": cft.MapOfInterfaces{
				"Service": "eks-gamma.aws.internal",
			},
			"Action": []string{
				"sts:AssumeRole",
				"sts:TagSession",
			},
		},
	}
	return cft.MakePolicyDocument(statements...)
}

func addBetaAccessEntry(stackName string, accessEntryType string) *gfn.CustomResource {
	customResource := &gfn.CustomResource{
		Type: "Custom::EksAccessEntry",
	}
	customResource.Properties = make(map[string]interface{})
	functionArn := gfnt.MakeFnImportValueString(fmt.Sprintf("eksctl-%s-cluster::EKSFunctionArn", stackName))
	customResource.Properties["ServiceToken"] = functionArn
	customResource.Properties["PrincipalArn"] = gfnt.MakeFnGetAttString(cfnIAMInstanceRoleName, "Arn")
	customResource.Properties["ClusterName"] = gfnt.NewString(stackName)
	customResource.Properties["Type"] = gfnt.NewString(accessEntryType)
	return customResource
}
