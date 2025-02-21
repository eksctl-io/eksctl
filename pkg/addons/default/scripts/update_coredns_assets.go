package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	awseks "github.com/aws/aws-sdk-go-v2/service/eks"
	"github.com/blang/semver"

	api "github.com/weaveworks/eksctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/weaveworks/eksctl/pkg/eks"
)

func main() {
	ctx := context.Background()

	clusterProvider, err := eks.New(ctx, &api.ProviderConfig{}, nil)
	if err != nil {
		log.Fatalf("failed to create the AWS provider: %v", err)
	}

	cvm, err := eks.NewClusterVersionsManager(clusterProvider.AWSProvider.EKS())
	if err != nil {
		log.Fatal(err)
	}

	supportedVersions := cvm.SupportedVersions()
	for index, kubernetesVersion := range supportedVersions {
		latestVersion := getLatestCoreDNSVersion(ctx, clusterProvider, kubernetesVersion)
		if latestVersion == "" {
			continue
		}
		filePath := filepath.Join("pkg", "addons", "default", "assets",
			fmt.Sprintf("coredns-%s.json", kubernetesVersion))
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			filePrevVersionPath := filepath.Join("pkg", "addons", "default", "assets",
				fmt.Sprintf("coredns-%s.json", supportedVersions[index-1]))
			copyFile(filePrevVersionPath, filePath)
		}
		replaceCurrentVersionIfOutdated(filePath, latestVersion, kubernetesVersion)
	}

}

func copyFile(src, dst string) {
	content, err := os.ReadFile(src)
	if err != nil {
		log.Fatalf("failed to read %s: %v", src, err)
	}
	if err := os.WriteFile(dst, content, 0644); err != nil {
		log.Fatalf("failed to write %s: %v", dst, err)
	}
}

func getLatestCoreDNSVersion(ctx context.Context, clusterProvider *eks.ClusterProvider, kubernetesVersion string) string {
	output, err := clusterProvider.AWSProvider.EKS().DescribeAddonVersions(ctx, &awseks.DescribeAddonVersionsInput{
		AddonName:         aws.String("coredns"),
		KubernetesVersion: &kubernetesVersion,
	})
	if err != nil {
		log.Fatalf("failed calling EKS::DescribeAddonVersions: %v", err)
	}

	if len(output.Addons[0].AddonVersions) == 0 {
		return ""
	}
	var corednsVersions []string
	regexpVersion := regexp.MustCompile(`v\d+\.\d+\.\d+-eksbuild\.\d+`)
	for _, info := range output.Addons[0].AddonVersions {
		if regexpVersion.MatchString(*info.AddonVersion) {
			corednsVersions = append(corednsVersions, *info.AddonVersion)
		}
	}

	sort.Slice(corednsVersions, func(i, j int) bool {
		vi, err := semver.Parse(trim(corednsVersions[i]))
		if err != nil {
			log.Fatalf("failed to parse coredns version %s: %v", trim(corednsVersions[i]), err)
		}
		vj, err := semver.Parse(trim(corednsVersions[j]))
		if err != nil {
			log.Fatalf("failed to parse coredns version %s: %v", trim(corednsVersions[j]), err)
		}
		if vi.Compare(vj) >= 0 {
			return true
		}
		return false
	})

	return corednsVersions[0]
}

func replaceCurrentVersionIfOutdated(filePath string, latestVersion string, kubernetesVersion string) {
	coreFile, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("failed to read coredns-%s.json: %v", kubernetesVersion, err)
	}

	regexpVersion := regexp.MustCompile(`v\d+\.\d+\.\d+-eksbuild\.\d+`)
	currentVersion := regexpVersion.FindString(string(coreFile))
	if currentVersion == "" {
		log.Fatalf("couldn't find coredns version in coredns-%s.json", kubernetesVersion)
	}

	updatedCoreFile := regexpVersion.ReplaceAllString(string(coreFile), latestVersion)
	if err := os.WriteFile(filePath, []byte(updatedCoreFile), 0644); err != nil {
		log.Fatalf("failed to write coredns-%s.json: %v", kubernetesVersion, err)
	}
}

func trim(version string) string {
	return strings.TrimPrefix(version, "v")
}
