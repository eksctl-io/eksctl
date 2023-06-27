package assets

import (
	// Import go:embed
	_ "embed"
)

//BootstrapAl2Sh holds the bootstrap.al2.sh contents
//go:embed scripts/bootstrap.al2.sh
var BootstrapAl2Sh string

//BootstrapHelperSh holds the bootstrap.helper.sh contents
//go:embed scripts/bootstrap.helper.sh
var BootstrapHelperSh string

//BootstrapUbuntuSh holds the bootstrap.ubuntu.sh contents
//go:embed scripts/bootstrap.ubuntu.sh
var BootstrapUbuntuSh string

//EfaAl2Sh holds the efa.al2.sh contents
//go:embed scripts/efa.al2.sh
var EfaAl2Sh string

//EfaManagedBoothook holds the efa.managed.boothook contents
//go:embed scripts/efa.managed.boothook
var EfaManagedBoothook string

//InstallSsmAl2Sh holds the install-ssm.al2.sh contents
//go:embed scripts/install-ssm.al2.sh
var InstallSsmAl2Sh string

//KubeletYaml holds the kubelet.yaml contents
//go:embed scripts/kubelet.yaml
var KubeletYaml string
