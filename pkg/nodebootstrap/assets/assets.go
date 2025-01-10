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

// KubeletYaml holds the kubelet.yaml contents
//
//go:embed scripts/kubelet.yaml
var KubeletYaml string
