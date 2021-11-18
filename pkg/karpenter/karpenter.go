package karpenter

import "github.com/weaveworks/eksctl/pkg/karpenter/providers"

const (
	karpenterHelmRepo      = "https://charts.karpenter.sh"
	karpenterHelmChartName = "karpenter/karpenter"
	karpenterReleaseName   = "karpenter"
)

// KarpenterInstaller defines an installer for Karpenter.
type KarpenterInstaller interface {
	InstallKarpenter()
	UninstallKarpenter()
}

// Installer implements the Karpenter installer using a HelmInstaller.
type Installer struct {
	HelmInstaller providers.HelmInstaller
}

func NewKarpenterInstaller(installer providers.HelmInstaller) *Installer {
	return &Installer{
		HelmInstaller: installer,
	}
}

func (k *Installer) InstallKarpenter() {

}

func (k *Installer) UninstallKarpenter() {

}
