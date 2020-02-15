package insttypes

import (
	"fmt"
	"os"
	"reflect"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var DefaultTotalEphemStorageSize = 20 //GB

func isPtr(i interface{}) bool {
	return reflect.ValueOf(i).Type().Kind() == reflect.Ptr
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func deref(i interface{}) interface{} {
	if isNil(i) || !isPtr(i) {
		return i
	}
	ia := reflect.ValueOf(i)
	reflect.ValueOf(i).Elem().Set(ia)
	return ia
}

func toInt(i int64) int {
	return int(i)
}

type EC2InstanceTypeInfo struct {
	InstanceStorageInfo      *EC2InstanceStorageInfo
	InstanceStorageSupported bool
	InstanceType             string
	MemoryInfo               *EC2MemoryInfo
	VCpuInfo                 *EC2VCpuInfo
}

func toEC2InstanceTypeInfo(iti *ec2.InstanceTypeInfo) *EC2InstanceTypeInfo {
	eiti := EC2InstanceTypeInfo{}
	eiti.InstanceStorageInfo = toEC2InstanceStorageInfo(iti.InstanceStorageInfo)
	eiti.InstanceStorageSupported = *iti.InstanceStorageSupported
	eiti.InstanceType = *iti.InstanceType
	eiti.MemoryInfo = toEC2MemoryInfo(iti.MemoryInfo)
	eiti.VCpuInfo = toEC2VCpuInfo(iti.VCpuInfo)
	return &eiti
}

func toEC2InstanceTypeInfos(itis []*ec2.InstanceTypeInfo) []*EC2InstanceTypeInfo {
	itypeInfos := make([]*EC2InstanceTypeInfo, len(itis))
	for i, iti := range itis {
		itypeInfos[i] = toEC2InstanceTypeInfo(iti)
	}
	return itypeInfos
}

func (eiti *EC2InstanceTypeInfo) String() string {
	return fmt.Sprintf("%+v\n%v\n%s\n%v\n%v\n",
		deref(eiti.InstanceStorageInfo),
		deref(eiti.InstanceStorageSupported),
		deref(eiti.InstanceType),
		deref(eiti.MemoryInfo),
		deref(eiti.VCpuInfo),
	)
}

type EC2InstanceStorageInfo struct {
	TotalSizeInGB int
}

func toEC2InstanceStorageInfo(isi *ec2.InstanceStorageInfo) *EC2InstanceStorageInfo {
	eisi := EC2InstanceStorageInfo{}
	if isi == nil {
		eisi.TotalSizeInGB = DefaultTotalEphemStorageSize
	} else {
		eisi.TotalSizeInGB = toInt(*isi.TotalSizeInGB)
	}
	return &eisi
}

func (eisi *EC2InstanceStorageInfo) String() string {
	return fmt.Sprintf("%v\n",
		deref(eisi.TotalSizeInGB),
	)
}

type EC2MemoryInfo struct {
	SizeInMiB int
}

func toEC2MemoryInfo(mi *ec2.MemoryInfo) *EC2MemoryInfo {
	emi := EC2MemoryInfo{}
	emi.SizeInMiB = toInt(*mi.SizeInMiB)
	return &emi
}

func (emi *EC2MemoryInfo) String() string {
	return fmt.Sprintf("%v\n",
		emi.SizeInMiB,
	)
}

type EC2VCpuInfo struct {
	DefaultVCpus int
}

func toEC2VCpuInfo(vci *ec2.VCpuInfo) *EC2VCpuInfo {
	evci := EC2VCpuInfo{}
	evci.DefaultVCpus = toInt(*vci.DefaultVCpus)
	return &evci
}

func (evci *EC2VCpuInfo) String() string {
	return fmt.Sprintf("%v\n",
		evci.DefaultVCpus,
	)
}

var EC2InstanceTypeMapping map[string]*EC2InstanceTypeInfo

// supportedRegions are the regions where EKS is available
func supportedRegions() []string {
	return []string{
		"us-west-2",
		"us-east-1",
		"us-east-2",
		"ca-central-1",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
		"eu-north-1",
		"eu-central-1",
		"ap-northeast-1",
		"ap-northeast-2",
		"ap-southeast-1",
		// "ap-southest-2",
		"ap-south-1",
		"ap-east-1",
		"me-south-1",
		"sa-east-1",
	}
}

//go:generate go run ./static_resolver_itype_generate.go

// Generate InstanceTypes maps (per region)
func GenerateRegionalInstanceTypesMap() map[string]map[string]*EC2InstanceTypeInfo {

	clients := newMultiRegionClient()

	var regionMap map[string]map[string]*EC2InstanceTypeInfo = make(map[string]map[string]*EC2InstanceTypeInfo, len(supportedRegions()))
	for _, region := range supportedRegions() {
		client, ok := clients[region]
		if !ok {
			exitErrorf("unable to get ec2 client for region %s", region)
		}
		regionalInstTypes, err := getInstanceTypes(client)
		if err != nil {
			exitErrorf("unable to get instance types for region %s", err)
		}
		EC2InstanceTypeMapping := make(map[string]*EC2InstanceTypeInfo)
		for _, itype := range regionalInstTypes {
			EC2InstanceTypeMapping[itype.InstanceType] = itype
			regionMap[region] = EC2InstanceTypeMapping
		}
	}
	return regionMap
}

func getInstanceTypes(svc *ec2.EC2) ([]*EC2InstanceTypeInfo, error) {

	//  Returns a list of key/value pairs
	input := &ec2.DescribeInstanceTypesInput{}
	instTypes, err := svc.DescribeInstanceTypes(input)
	if err != nil {
		return nil, err
	}
	var itypes []*ec2.InstanceTypeInfo

	token := instTypes.NextToken
	for token != nil {
		input := &ec2.DescribeInstanceTypesInput{NextToken: token}
		itypesOut, err := svc.DescribeInstanceTypes(input)
		if err != nil {
			return nil, err
		}
		itypes = append(itypes, itypesOut.InstanceTypes...)
		for _, instType := range itypesOut.InstanceTypes {
			itypes = append(itypes, instType)
		}
		token = itypesOut.NextToken
	}
	return toEC2InstanceTypeInfos(itypes), nil
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}

func newSession(region string) *session.Session {
	config := aws.NewConfig()
	config = config.WithRegion(region)
	config = config.WithCredentialsChainVerboseErrors(true)

	// Create the options for the session
	opts := session.Options{
		Config:                  *config,
		SharedConfigState:       session.SharedConfigEnable,
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
	}

	return session.Must(session.NewSessionWithOptions(opts))
}

func newMultiRegionClient() map[string]*ec2.EC2 {
	clients := make(map[string]*ec2.EC2)
	for _, region := range supportedRegions() {
		clients[region] = ec2.New(newSession(region))
	}
	return clients
}
