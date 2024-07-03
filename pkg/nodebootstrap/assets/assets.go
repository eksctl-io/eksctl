package assets

import (
	// Import go:embed
	_ "embed"
)

// BootstrapAl2Sh holds the bootstrap.al2.sh contents
//
//go:embed scripts/bootstrap.al2.sh
var BootstrapAl2Sh string

// BootstrapHelperSh holds the bootstrap.helper.sh contents
//
//go:embed scripts/bootstrap.helper.sh
var BootstrapHelperSh string

// BootstrapUbuntuSh holds the bootstrap.ubuntu.sh contents
//
//go:embed scripts/bootstrap.ubuntu.sh
var BootstrapUbuntuSh string

// AL2023XTablesLock holds the contents for creating a lock file for AL2023 AMIs.
//
//go:embed scripts/al2023-xtables.lock.sh
var AL2023XTablesLock string

// EfaAl2Sh holds the efa.al2.sh contents
//
//go:embed scripts/efa.al2.sh
var EfaAl2Sh string

// EfaAl2023Sh holds the efa.al2023.sh contents
//
//go:embed scripts/efa.al2023.sh
var EfaAl2023Sh string

// EfaManagedBoothook holds the efa.managed.boothook contents
//
//go:embed scripts/efa.managed.boothook
var EfaManagedBoothook string

// EfaManagedAL2023Boothook holds the efa.managed.al2023.boothook contents
//
//go:embed scripts/efa.managed.al2023.boothook
var EfaManagedAL2023Boothook string

// InstallSsmAl2Sh holds the install-ssm.al2.sh contents
//
//go:embed scripts/install-ssm.al2.sh
var InstallSsmAl2Sh string

// KubeletYaml holds the kubelet.yaml contents
//
//go:embed scripts/kubelet.yaml
var KubeletYaml string
