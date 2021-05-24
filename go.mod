// Make sure to run the following commands after changes to this file are made:
// ` make -f Makefile.docker update-build-image-tag && make -f Makefile.docker push-build-image`
module github.com/weaveworks/eksctl

go 1.16

require (
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/aws/amazon-ec2-instance-selector/v2 v2.0.3-0.20210303155736-3e43512d88f8
	github.com/aws/aws-sdk-go v1.38.38
	github.com/benjamintf1/unmarshalledmatchers v0.0.0-20190408201839-bb1c1f34eaea
	github.com/blang/semver v3.5.1+incompatible
	github.com/bxcodec/faker v2.0.1+incompatible
	github.com/cloudflare/cfssl v1.5.0
	github.com/dave/jennifer v1.4.1
	github.com/dlespiau/kube-test-harness v0.0.0-20200915102055-a03579200ae8
	github.com/evanphx/json-patch/v5 v5.2.0
	github.com/fatih/color v1.10.0
	github.com/fluxcd/flux/pkg/install v0.0.0-20201001122558-cb08da1b356a // flux 1.21.0
	github.com/fluxcd/go-git-providers v0.0.3
	github.com/fluxcd/helm-operator/pkg/install v0.0.0-20200729150005-1467489f7ee4 // helm-operator 1.2.0
	github.com/github-release/github-release v0.10.0
	github.com/gobwas/glob v0.2.3
	github.com/gofrs/flock v0.8.0
	github.com/golangci/golangci-lint v1.39.0
	github.com/gomarkdown/markdown v0.0.0-20201113031856-722100d81a8e // indirect
	github.com/google/certificate-transparency-go v1.1.1 // indirect
	github.com/google/uuid v1.2.0
	github.com/goreleaser/goreleaser v0.162.0
	github.com/hashicorp/go-version v1.3.0
	github.com/instrumenta/kubeval v0.0.0-20190918223246-8d013ec9fc56
	github.com/justinbarrick/go-k8s-portforward v1.0.4-0.20200904152830-b575325c1855
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/kevinburke/go-bindata v3.22.0+incompatible
	github.com/kevinburke/rest v0.0.0-20210106114233-22cd0577e450 // indirect
	github.com/kris-nova/logger v0.2.1
	github.com/kris-nova/lolgopher v0.0.0-20210112022122-73f0047e8b65
	github.com/kris-nova/novaarchive v0.0.0-20210219195539-c7c1cabb2577 // indirect
	github.com/kubicorn/kubicorn v0.0.0-20180829191017-06f6bce92acc
	github.com/lithammer/dedent v1.1.0
	github.com/maxbrunsfeld/counterfeiter/v6 v6.4.1
	github.com/onsi/ginkgo v1.16.1
	github.com/onsi/gomega v1.11.0
	github.com/pelletier/go-toml v1.9.0
	github.com/pkg/errors v0.9.1
	github.com/russross/blackfriday/v2 v2.1.0
	github.com/spf13/afero v1.6.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/gjson v1.7.5
	github.com/tidwall/sjson v1.1.6
	github.com/tj/assert v0.0.3
	github.com/tomnomnom/linkheader v0.0.0-20180905144013-02ca5825eb80 // indirect
	github.com/vektra/mockery v1.1.2
	github.com/voxelbrain/goptions v0.0.0-20180630082107-58cddc247ea2 // indirect
	github.com/weaveworks/goformation/v4 v4.10.2-0.20210524152715-0063d430cbd7
	github.com/weaveworks/launcher v0.0.2-0.20200715141516-1ca323f1de15
	github.com/whilp/git-urls v0.0.0-20191001220047-6db9661140c0
	golang.org/x/tools v0.1.0
	k8s.io/api v0.19.5
	k8s.io/apiextensions-apiserver v0.19.5
	k8s.io/apimachinery v0.19.5
	k8s.io/cli-runtime v0.19.5
	k8s.io/client-go v0.19.5
	k8s.io/cloud-provider v0.19.5
	k8s.io/code-generator v0.19.5
	k8s.io/kops v1.19.0
	k8s.io/kubelet v0.19.5
	k8s.io/kubernetes v1.19.5
	k8s.io/legacy-cloud-providers v0.19.5
	sigs.k8s.io/aws-iam-authenticator v0.5.2
	sigs.k8s.io/mdtoc v1.0.1
	sigs.k8s.io/yaml v1.2.0
)

replace (
	// Used to get around some weird etcd/grpc incompatibilty
	google.golang.org/grpc => google.golang.org/grpc v1.29.0
	// Used to pin the k8s library versions regardless of what other dependencies enforce
	k8s.io/api => k8s.io/api v0.19.5
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.19.5
	k8s.io/apimachinery => k8s.io/apimachinery v0.19.5
	k8s.io/apiserver => k8s.io/apiserver v0.19.5
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.19.5
	k8s.io/client-go => k8s.io/client-go v0.19.5
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.19.5
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.19.5
	k8s.io/code-generator => k8s.io/code-generator v0.19.5
	k8s.io/component-base => k8s.io/component-base v0.19.5
	k8s.io/cri-api => k8s.io/cri-api v0.19.5
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.19.5
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.19.5
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.19.5
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.19.5
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.19.5
	k8s.io/kubectl => k8s.io/kubectl v0.19.5
	k8s.io/kubelet => k8s.io/kubelet v0.19.5
	k8s.io/kubernetes => k8s.io/kubernetes v1.19.5
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.19.5
	k8s.io/metrics => k8s.io/metrics v0.19.5
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.19.5
)
