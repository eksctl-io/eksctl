package v1alpha5

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	inst_types "github.com/weaveworks/eksctl/pkg/insttypes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SetClusterConfigDefaults will set defaults for a given cluster
func SetClusterConfigDefaults(cfg *ClusterConfig) {
	if cfg.IAM == nil {
		cfg.IAM = &ClusterIAM{}
	}

	if cfg.IAM.WithOIDC == nil {
		cfg.IAM.WithOIDC = Disabled()
	}

	for _, sa := range cfg.IAM.ServiceAccounts {
		if sa.Namespace == "" {
			sa.Namespace = metav1.NamespaceDefault
		}
	}

	if cfg.HasClusterCloudWatchLogging() && len(cfg.CloudWatch.ClusterLogging.EnableTypes) == 1 {
		switch cfg.CloudWatch.ClusterLogging.EnableTypes[0] {
		case "all", "*":
			cfg.CloudWatch.ClusterLogging.EnableTypes = SupportedCloudWatchClusterLogTypes()
		}
	}
}

// SetNodeGroupDefaults will set defaults for a given nodegroup
func SetNodeGroupDefaults(ng *NodeGroup, meta *ClusterMeta) {
	if ng.InstanceType == "" {
		if HasMixedInstances(ng) {
			ng.InstanceType = "mixed"
		} else {
			ng.InstanceType = DefaultNodeType
		}
	}
	if ng.AMIFamily == "" {
		ng.AMIFamily = DefaultNodeImageFamily
	}

	if ng.SecurityGroups == nil {
		ng.SecurityGroups = &NodeGroupSGs{
			AttachIDs: []string{},
		}
	}
	if ng.SecurityGroups.WithLocal == nil {
		ng.SecurityGroups.WithLocal = Enabled()
	}
	if ng.SecurityGroups.WithShared == nil {
		ng.SecurityGroups.WithShared = Enabled()
	}

	if ng.SSH == nil {
		ng.SSH = &NodeGroupSSH{
			Allow: Disabled(),
		}
	}

	setSSHDefaults(ng.SSH)

	if !IsSetAndNonEmptyString(ng.VolumeType) {
		ng.VolumeType = &DefaultNodeVolumeType
	}

	if ng.IAM == nil {
		ng.IAM = &NodeGroupIAM{}
	}

	setIAMDefaults(ng.IAM)

	if ng.Labels == nil {
		ng.Labels = make(map[string]string)
	}
	setDefaultNodeLabels(ng.Labels, meta.Name, ng.Name)

	switch ng.AMIFamily {
	case NodeImageFamilyBottlerocket:
		setBottlerocketNodeGroupDefaults(ng)
	default:
		err := SetKubeletExtraConfigDefaults(ng, meta)
		if err != nil {
			fmt.Printf("Encountered error when setting KubeletConfig defaults: %s\n", err.Error())
			os.Exit(2)
		}
	}
}

// SetManagedNodeGroupDefaults sets default values for a ManagedNodeGroup
func SetManagedNodeGroupDefaults(ng *ManagedNodeGroup, meta *ClusterMeta) {
	if ng.AMIFamily == "" {
		ng.AMIFamily = NodeImageFamilyAmazonLinux2
	}
	if ng.InstanceType == "" {
		ng.InstanceType = DefaultNodeType
	}
	if ng.ScalingConfig == nil {
		ng.ScalingConfig = &ScalingConfig{}
	}
	if ng.SSH == nil {
		ng.SSH = &NodeGroupSSH{
			Allow: Disabled(),
		}
	}
	setSSHDefaults(ng.SSH)

	if ng.IAM == nil {
		ng.IAM = &NodeGroupIAM{}
	}
	setIAMDefaults(ng.IAM)

	if ng.Labels == nil {
		ng.Labels = make(map[string]string)
	}
	setDefaultNodeLabels(ng.Labels, meta.Name, ng.Name)

	if ng.Tags == nil {
		ng.Tags = make(map[string]string)
	}
	ng.Tags[NodeGroupNameTag] = ng.Name
	ng.Tags[NodeGroupTypeTag] = string(NodeGroupTypeManaged)
}

func setIAMDefaults(iamConfig *NodeGroupIAM) {
	if iamConfig.WithAddonPolicies.ImageBuilder == nil {
		iamConfig.WithAddonPolicies.ImageBuilder = Disabled()
	}
	if iamConfig.WithAddonPolicies.AutoScaler == nil {
		iamConfig.WithAddonPolicies.AutoScaler = Disabled()
	}
	if iamConfig.WithAddonPolicies.ExternalDNS == nil {
		iamConfig.WithAddonPolicies.ExternalDNS = Disabled()
	}
	if iamConfig.WithAddonPolicies.CertManager == nil {
		iamConfig.WithAddonPolicies.CertManager = Disabled()
	}
	if iamConfig.WithAddonPolicies.ALBIngress == nil {
		iamConfig.WithAddonPolicies.ALBIngress = Disabled()
	}
	if iamConfig.WithAddonPolicies.XRay == nil {
		iamConfig.WithAddonPolicies.XRay = Disabled()
	}
	if iamConfig.WithAddonPolicies.CloudWatch == nil {
		iamConfig.WithAddonPolicies.CloudWatch = Disabled()
	}
	if iamConfig.WithAddonPolicies.EBS == nil {
		iamConfig.WithAddonPolicies.EBS = Disabled()
	}
	if iamConfig.WithAddonPolicies.FSX == nil {
		iamConfig.WithAddonPolicies.FSX = Disabled()
	}
	if iamConfig.WithAddonPolicies.EFS == nil {
		iamConfig.WithAddonPolicies.EFS = Disabled()
	}
}

func setSSHDefaults(sshConfig *NodeGroupSSH) {
	numSSHFlagsEnabled := countEnabledFields(
		sshConfig.PublicKeyName,
		sshConfig.PublicKeyPath,
		sshConfig.PublicKey)

	if numSSHFlagsEnabled == 0 {
		if IsEnabled(sshConfig.Allow) {
			sshConfig.PublicKeyPath = &DefaultNodeSSHPublicKeyPath
		} else {
			sshConfig.Allow = Disabled()
		}
	} else if !IsDisabled(sshConfig.Allow) {
		// Enable SSH if not explicitly disabled when passing an SSH key
		sshConfig.Allow = Enabled()
	}

}

func setDefaultNodeLabels(labels map[string]string, clusterName, nodeGroupName string) {
	labels[ClusterNameLabel] = clusterName
	labels[NodeGroupNameLabel] = nodeGroupName
}

type getRscDefaultFunc func(string, *ClusterMeta) (string, error)
type setRscDefaultFunc func(*NodeGroup, string, *ClusterMeta, getRscDefaultFunc) error

var rscParams = []struct {
	setFn   setRscDefaultFunc
	getFn   getRscDefaultFunc
	rscType string
}{
	{setCPUReservationDefaults, getCPUReservations, "cpu"},
	{setMemoryResevationDefaults, getMemReservations, "memory"},
	{setEphemeralStorageDefaults, getEphemeralStorageReservations, "ephemeral-storage"},
}

// SetKubeletExtraConfigDefaults adds Kubelet CPU, Mem, and Storage Reservation default values for a nodegroup
func SetKubeletExtraConfigDefaults(ng *NodeGroup, meta *ClusterMeta) error {
	for _, pSet := range rscParams {
		err := pSet.setFn(ng, pSet.rscType, meta, pSet.getFn)
		if err != nil {
			return err
		}
	}
	return nil
}

func setCPUReservationDefaults(ng *NodeGroup, rscType string, meta *ClusterMeta, getFn getRscDefaultFunc) error {
	return setReservationDefault(ng, rscType, meta, getFn)
}

func setMemoryResevationDefaults(ng *NodeGroup, rscType string, meta *ClusterMeta, getFn getRscDefaultFunc) error {
	return setReservationDefault(ng, rscType, meta, getFn)
}

func setEphemeralStorageDefaults(ng *NodeGroup, rscType string, meta *ClusterMeta, getFn getRscDefaultFunc) error {
	return setReservationDefault(ng, rscType, meta, getFn)
}

func setReservationDefault(ng *NodeGroup, resType string, meta *ClusterMeta, setFn getRscDefaultFunc) error {
	kec := (*ng).KubeletExtraConfig
	if kec == nil {
		kec = &InlineDocument{}
	}
	rsrcRes, err := setFn((*ng).InstanceType, meta)
	if err != nil {
		return err
	}
	kubeReserved := getKubeReserved(*kec)
	// only set kubelet reservations for resource types that aren't already set in config
	if _, ok := kubeReserved[resType]; !ok {
		kubeReserved[resType] = rsrcRes
	}
	(*kec)["kubeReserved"] = kubeReserved
	ng.KubeletExtraConfig = kec
	return nil
}

func getKubeReserved(kec InlineDocument) map[string]interface{} {
	kubeReserved, ok := kec["kubeReserved"].(map[string]interface{})
	if !ok {
		kubeReserved = make(map[string]interface{})
	}
	return kubeReserved
}

// See: https://docs.microsoft.com/en-us/azure/aks/concepts-clusters-workloads
var cpuAllocations = map[int]string{
	1:  "60m",
	2:  "100m",  //+40
	4:  "140m",  //+40
	8:  "180m",  //+40
	16: "260m",  //+80
	32: "420m",  //+160
	48: "580m",  //+160
	64: "740m",  //+320
	96: "1040m", //+320
}

func getCPUReservations(it string, meta *ClusterMeta) (string, error) {
	cores, err := getInstanceTypeCores(it, meta)
	if err != nil {
		return "", err
	}

	reservedCores, ok := cpuAllocations[cores]
	if !ok {
		return "", fmt.Errorf("could not find suggested core reservation for instance type: %s", it)
	}
	return reservedCores, nil
}

func getInstanceTypeCores(it string, meta *ClusterMeta) (int, error) {
	instTypeInfos, err := getInstanceTypeInfo(it, meta)
	if err != nil {
		return 0, err
	}
	vCPUInfo := (*instTypeInfos).VCpuInfo
	cpuCores := vCPUInfo.DefaultVCpus
	return cpuCores, nil
}

type memEntry struct {
	max      float64
	fraction float64
}

// See: https://docs.microsoft.com/en-us/azure/aks/concepts-clusters-workloads
var memPercentages = []memEntry{
	{max: 4, fraction: 0.25},
	{max: 8, fraction: 0.20},
	{max: 16, fraction: 0.10},
	{max: 128, fraction: 0.06},
	{max: 65535, fraction: 0.02},
}

func getMemReservations(it string, meta *ClusterMeta) (string, error) {
	instMem, err := getInstanceTypeMem(it, meta)
	if err != nil {
		return "", err
	}
	var lower, reserved float64 = 0.0, 0.0
	for _, memEnt := range memPercentages {
		k, v := memEnt.max, memEnt.fraction
		if instMem <= k {
			reserved += v * (instMem - lower)
			break
		} else {
			reserved += v * (k - lower)
		}
		lower = k
	}
	reservedStr := formatMem(reserved)
	return reservedStr, nil
}

// formatFloat removes duplicate trailing zeros and ensures decimal format
func formatMem(f float64) string {
	ff := strconv.FormatFloat(f, 'f', -1, 32)
	if !strings.Contains(ff, ".") {
		ff += ".0"
	}
	return ff + "Mi"
}

func getInstanceTypeMem(it string, meta *ClusterMeta) (float64, error) {
	instTypeInfo, err := getInstanceTypeInfo(it, meta)
	if err != nil {
		return 0, err
	}
	memInfo := (*instTypeInfo).MemoryInfo
	memSize := float64(memInfo.SizeInMiB)
	memStr := fmt.Sprintf("%.2f", float64(memSize/1024.0))
	return strconv.ParseFloat(memStr, 64)
}

func getEphemeralStorageReservations(it string, meta *ClusterMeta) (string, error) {
	storageSize, err := getInstanceTypeStorage(it, meta)
	if err != nil {
		return "", err
	}
	// at least 1GB but no larger than 15GB
	larger := math.Max(1.0, float64(storageSize)/16.0)
	smaller := math.Min(15.0, larger)
	storSize, storErr := formatStorageSize(smaller)
	return storSize, storErr
}

func formatStorageSize(f float64) (string, error) {
	// set precision to 2 decimal points
	fstr := fmt.Sprintf("%.2f", f)
	f64, err := strconv.ParseFloat(fstr, 64)
	if err != nil {
		return "", err
	}
	// remove any trailing zeros and convert to string
	return strconv.FormatFloat(f64, 'f', -1, 64) + "Gi", nil
}

func getInstanceTypeStorage(it string, meta *ClusterMeta) (int, error) {
	defaultInstanceTypeStorage := 20 //GB
	instTypeInfo, err := getInstanceTypeInfo(it, meta)
	if err != nil {
		return 0, err
	}
	// If no default instance storage defined in instance type
	if !instTypeInfo.InstanceStorageSupported {
		return defaultInstanceTypeStorage, nil
	}
	storageSize := (*instTypeInfo).InstanceStorageInfo.TotalSizeInGB
	return storageSize, nil
}

func getInstanceTypeInfo(it string, meta *ClusterMeta) (*inst_types.EC2InstanceTypeInfo, error) {
	if meta == nil {
		meta = &ClusterMeta{}
	}
	if meta.Region == "" {
		meta.Region = DefaultRegion
	}
	regionMap, ok := inst_types.StaticInstanceTypes[meta.Region]
	if !ok {
		return nil, fmt.Errorf("unable to find region \"%s\" in region map", meta.Region)
	}
	instTypeInfo, ok := regionMap[it]
	if !ok {
		return nil, fmt.Errorf("unable to find instance type: \"%s\" in region: \"%s\"", it, meta.Region)
	}
	return instTypeInfo, nil
}

func setBottlerocketNodeGroupDefaults(ng *NodeGroup) {
	// Initialize config object if not present.
	if ng.Bottlerocket == nil {
		ng.Bottlerocket = &NodeGroupBottlerocket{}
	}

	// Default to resolving Bottlerocket images using SSM if not specified by
	// the user.
	if ng.AMI == "" {
		ng.AMI = NodeImageResolverAutoSSM
	}

	// Use the SSH settings if the user hasn't explicitly configured the Admin
	// Container. If SSH was enabled, the user will be able to ssh into the
	// Bottlerocket node via the admin container.
	if ng.Bottlerocket.EnableAdminContainer == nil && ng.SSH != nil {
		ng.Bottlerocket.EnableAdminContainer = ng.SSH.Allow
	}
}

// DefaultClusterNAT will set the default value for Cluster NAT mode
func DefaultClusterNAT() *ClusterNAT {
	single := ClusterSingleNAT
	return &ClusterNAT{
		Gateway: &single,
	}
}

// SetClusterEndpointAccessDefaults sets the default values for cluster endpoint access
func SetClusterEndpointAccessDefaults(vpc *ClusterVPC) {
	if vpc.ClusterEndpoints == nil {
		vpc.ClusterEndpoints = ClusterEndpointAccessDefaults()
	}

	if vpc.ClusterEndpoints.PublicAccess == nil {
		vpc.ClusterEndpoints.PublicAccess = Enabled()
	}

	if vpc.ClusterEndpoints.PrivateAccess == nil {
		vpc.ClusterEndpoints.PrivateAccess = Disabled()
	}
}

// ClusterEndpointAccessDefaults returns a ClusterEndpoints pointer with default values set.
func ClusterEndpointAccessDefaults() *ClusterEndpoints {
	return &ClusterEndpoints{
		PrivateAccess: Disabled(),
		PublicAccess:  Enabled(),
	}
}

// SetDefaultFargateProfile configures this ClusterConfig to have a single
// Fargate profile called "default", with two selectors matching respectively
// the "default" and "kube-system" Kubernetes namespaces.
func (c *ClusterConfig) SetDefaultFargateProfile() {
	c.FargateProfiles = []*FargateProfile{
		{
			Name: "fp-default",
			Selectors: []FargateProfileSelector{
				{Namespace: "default"},
				{Namespace: "kube-system"},
			},
		},
	}
}
