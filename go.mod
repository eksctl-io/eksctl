// Make sure to bump the version of EKSCTL_DEPENDENCIES_IMAGE if you make any changes here
module github.com/weaveworks/eksctl

go 1.12

require (
	github.com/MakeNowJust/heredoc v0.0.0-20171113091838-e9091a26100e // indirect
	github.com/alecthomas/jsonschema v0.0.0-20190530235721-fd8d96416671
	github.com/aws/aws-sdk-go v1.23.15
	github.com/awslabs/goformation v0.0.0-20190320125420-ac0a17860cf1
	github.com/blang/semver v3.5.1+incompatible
	github.com/chai2010/gettext-go v0.0.0-20170215093142-bf70f2a70fb1 // indirect
	github.com/christopherhein/go-version v0.0.0-20180807222509-fee8dd1f7c24
	github.com/cloudflare/cfssl v0.0.0-20190726000631-633726f6bcb7
	github.com/coredns/coredns v0.0.0-20170910182647-1b60688dc8f7 // indirect
	github.com/coreos/bbolt v1.3.3 // indirect
	github.com/coreos/etcd v3.3.13+incompatible // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20190620071333-e64a0ec8b42a // indirect
	github.com/dave/jennifer v1.3.0
	github.com/dlespiau/kube-test-harness v0.0.0-20190110151726-c51c87635b61
	github.com/docker/distribution v2.7.1+incompatible // indirect
	github.com/docker/docker v1.13.1 // indirect
	github.com/evanphx/json-patch v4.2.0+incompatible
	github.com/exponent-io/jsonpath v0.0.0-20151013193312-d6023ce2651d // indirect
	github.com/fluxcd/helm-operator v1.0.0-rc1
	github.com/go-ini/ini v1.37.0 // indirect
	github.com/gobuffalo/envy v1.7.0 // indirect
	github.com/gobwas/glob v0.2.3
	github.com/gofrs/flock v0.7.1 // indirect
	github.com/gohugoio/hugo v0.55.6
	github.com/google/btree v1.0.0 // indirect
	github.com/google/certificate-transparency-go v1.0.21 // indirect
	github.com/goreleaser/goreleaser v0.110.0
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.9.2 // indirect
	github.com/instrumenta/kubeval v0.0.0-20190804145309-805845b47dfc
	github.com/jmoiron/sqlx v1.2.0 // indirect
	github.com/jteeuwen/go-bindata v3.0.8-0.20180305030458-6025e8de665b+incompatible
	github.com/justinbarrick/go-k8s-portforward v1.0.4-0.20190722134107-d79fe1b9d79d
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/kr/fs v0.1.0 // indirect
	github.com/kris-nova/logger v0.0.0-20181127235838-fd0d87064b06
	github.com/kris-nova/lolgopher v0.0.0-20180124180951-14d43f83481a // indirect
	github.com/kubernetes-sigs/aws-iam-authenticator v0.4.0
	github.com/kubicorn/kubicorn v0.0.0-20180829191017-06f6bce92acc
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/lithammer/dedent v1.1.0
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mattn/go-zglob v0.0.1 // indirect
	github.com/miekg/coredns v0.0.0-20170910182647-1b60688dc8f7 // indirect
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/pelletier/go-toml v1.4.0 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pkg/errors v0.8.1
	github.com/pkg/sftp v1.8.3 // indirect
	github.com/prometheus/client_golang v1.0.0 // indirect
	github.com/riywo/loginshell v0.0.0-20190610082906-2ed199a032f6
	github.com/rubenv/sql-migrate v0.0.0-20190902133344-8926f37f0bc1 // indirect
	github.com/sanathkr/yaml v1.0.0 // indirect
	github.com/sirupsen/logrus v1.4.2 // indirect

	github.com/spf13/afero v1.2.2

	github.com/spf13/cobra v0.0.4

	github.com/spf13/pflag v1.0.3
	github.com/spf13/viper v1.3.2
	github.com/spotinst/spotinst-sdk-go v0.0.0-20181012192533-fed4677dbf8f // indirect
	github.com/stretchr/testify v1.3.0
	github.com/tidwall/gjson v1.1.3
	github.com/tidwall/match v1.0.0 // indirect
	github.com/tidwall/sjson v1.0.2
	github.com/vektra/mockery v0.0.0-20181123154057-e78b021dcbb5
	github.com/vmware/govmomi v0.19.0 // indirect
	github.com/voxelbrain/goptions v0.0.0-20180630082107-58cddc247ea2 // indirect
	github.com/weaveworks/flux v0.0.0-20190729133003-c78ccd3706b5
	github.com/weaveworks/github-release v0.6.2
	github.com/weaveworks/launcher v0.0.0-20180711153254-f1b2830d4f2d
	github.com/whilp/git-urls v0.0.0-20160530060445-31bac0d230fa
	github.com/zmap/zlint v0.0.0-20190806182416-88c3f6b6f2f5 // indirect
	go.etcd.io/bbolt v1.3.3 // indirect
	go.uber.org/atomic v1.4.0 // indirect
	go.uber.org/zap v1.10.0 // indirect
	golang.org/x/tools v0.0.0-20190614205625-5aca471b1d59
	google.golang.org/grpc v1.21.1 // indirect
	gopkg.in/gcfg.v1 v1.2.3 // indirect
	gopkg.in/gorp.v1 v1.7.2 // indirect
	gopkg.in/ini.v1 v1.42.0 // indirect
	k8s.io/api v0.0.0-20190927115716-5d581ce610b0
	k8s.io/apiextensions-apiserver v0.0.0-20190918224502-6154570c2037
	k8s.io/apimachinery v0.0.0-20190927035529-0104e33c351d
	k8s.io/apiserver v0.0.0-20190918223255-26459790ef01 // indirect
	k8s.io/cli-runtime v0.0.0-20190918224932-e56234cc6b3d
	k8s.io/client-go v11.0.1-0.20190918222721-c0e3722d5cf0+incompatible
	k8s.io/cloud-provider v0.0.0-20190918225840-7f3416179ad8
	k8s.io/code-generator v0.0.0-20190808180452-d0071a119380
	k8s.io/csi-translation-lib v0.0.0-20190919022632-f5fbd7244482 // indirect
	k8s.io/helm v2.14.3+incompatible
	k8s.io/klog v1.0.0 // indirect
	k8s.io/kops v0.0.0-20190222135932-278e6606534e
	k8s.io/kubelet v0.0.0-20190313123811-3556bcde9670
	k8s.io/kubernetes v1.14.7
	k8s.io/utils v0.0.0-20190920012459-5008bf6f8cd6 // indirect
	sigs.k8s.io/kustomize v2.0.3+incompatible // indirect
	sigs.k8s.io/yaml v1.1.0
	vbom.ml/util v0.0.0-20180919145318-efcd4e0f9787 // indirect
)

replace (
	// Override since git.apache.org is down.  The docs say to fetch from github.
	git.apache.org/thrift.git => github.com/apache/thrift v0.0.0-20180902110319-2566ecd5d999
	// client-go@v11.0.1 Azure plugin requires go-autorest v11.1.0
	// This is derived from:
	//   https://github.com/kubernetes/client-go/blob/c0e3722d5cf089299130492c554bc2b4b8eeb2bb/Godeps/Godeps.json#L19
	//   https://github.com/Azure/go-autorest/commit/ea233b6412b0421a65dc6160e16c893364664a95
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v11.1.0+incompatible
	// Needed due to to Sirupsen/sirupsen case clash
	github.com/Sirupsen/logrus => github.com/sirupsen/logrus v1.4.2
	github.com/awslabs/goformation => github.com/errordeveloper/goformation v0.0.0-20190507151947-a31eae35e596
	// go mod appears to pick wrong version of github.com/docker/distribution automatically, which breaks k8s.io/kubernetes@v1.12.6
	github.com/docker/distribution => github.com/docker/distribution v0.0.0-20190619192407-5223c27422cc
	// Used to pin the k8s library versions to 1.14.7 regardless of what other dependencies enforce
	k8s.io/api => k8s.io/api v0.0.0-20190816222004-e3a6b8045b0b
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190918224502-6154570c2037
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190816221834-a9f1d8a9c101
	k8s.io/client-go => k8s.io/client-go v11.0.1-0.20190918222721-c0e3722d5cf0+incompatible
)
