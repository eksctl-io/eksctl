// Make sure to run the following commands after changes to this file are made:
// `make generate-all && make lint && make check-all-generated-files-up-to-date`
// you may also need to run `make push-build-image` depending on what has changed
module github.com/weaveworks/eksctl

go 1.20

require (
	github.com/Masterminds/semver/v3 v3.2.1
	github.com/aws/amazon-ec2-instance-selector/v2 v2.4.2-0.20230601180523-74e721cb8c1e
	github.com/aws/aws-sdk-go v1.46.0
	github.com/aws/aws-sdk-go-v2 v1.21.2
	github.com/aws/aws-sdk-go-v2/config v1.19.0
	github.com/aws/aws-sdk-go-v2/credentials v1.13.43
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.31.0
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.35.0
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.29.2
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.24.2
	github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider v1.27.2
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.124.0
	github.com/aws/aws-sdk-go-v2/service/eks v1.29.7
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing v1.17.2
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.21.6
	github.com/aws/aws-sdk-go-v2/service/iam v1.22.7
	github.com/aws/aws-sdk-go-v2/service/kms v1.24.7
	github.com/aws/aws-sdk-go-v2/service/outposts v1.30.2
	github.com/aws/aws-sdk-go-v2/service/ssm v1.39.0
	github.com/aws/aws-sdk-go-v2/service/sts v1.23.2
	github.com/aws/smithy-go v1.15.0
	github.com/benjamintf1/unmarshalledmatchers v1.0.0
	github.com/blang/semver v3.5.1+incompatible
	github.com/bxcodec/faker v2.0.1+incompatible
	github.com/cenk/backoff v2.2.1+incompatible
	github.com/cloudflare/cfssl v1.6.4
	github.com/dave/dst v0.27.3
	github.com/dave/jennifer v1.7.0
	github.com/dlespiau/kube-test-harness v0.0.0-20200915102055-a03579200ae8
	github.com/evanphx/json-patch/v5 v5.7.0
	github.com/fatih/color v1.15.0
	github.com/github-release/github-release v0.10.0
	github.com/gobwas/glob v0.2.3
	github.com/gofrs/flock v0.8.1
	github.com/golangci/golangci-lint v1.54.2
	github.com/google/uuid v1.3.1
	github.com/goreleaser/goreleaser v1.11.5
	github.com/hashicorp/go-version v1.6.0
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/kris-nova/logger v0.2.1
	github.com/kris-nova/lolgopher v0.0.0-20210112022122-73f0047e8b65
	github.com/kubicorn/kubicorn v0.0.0-20180829191017-06f6bce92acc
	github.com/lithammer/dedent v1.1.0
	github.com/maxbrunsfeld/counterfeiter/v6 v6.6.2
	github.com/onsi/ginkgo/v2 v2.12.1
	github.com/onsi/gomega v1.28.0
	github.com/orcaman/concurrent-map v1.0.0
	github.com/otiai10/copy v1.14.0
	github.com/pelletier/go-toml v1.9.5
	github.com/pkg/errors v0.9.1
	github.com/sethvargo/go-password v0.2.0
	github.com/spf13/afero v1.10.0
	github.com/spf13/cobra v1.7.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.4
	github.com/tidwall/gjson v1.17.0
	github.com/tidwall/sjson v1.2.5
	github.com/tj/assert v0.0.3
	github.com/vburenin/ifacemaker v1.2.1
	github.com/vektra/mockery/v2 v2.35.3
	github.com/weaveworks/goformation/v4 v4.10.2-0.20230526082129-5f5eaa9609b8
	github.com/weaveworks/schemer v0.0.0-20230525114451-47139fe25848
	github.com/xgfone/netaddr v0.5.1
	golang.org/x/crypto v0.14.0
	golang.org/x/exp v0.0.0-20231006140011-7918f672742d
	golang.org/x/oauth2 v0.13.0
	golang.org/x/sync v0.4.0
	golang.org/x/tools v0.14.0
	gopkg.in/yaml.v2 v2.4.0
	helm.sh/helm/v3 v3.11.2
	k8s.io/api v0.26.0
	k8s.io/apiextensions-apiserver v0.26.0
	k8s.io/apimachinery v0.26.0
	k8s.io/cli-runtime v0.26.0
	k8s.io/client-go v0.26.0
	k8s.io/cloud-provider v0.25.5
	k8s.io/code-generator v0.25.11
	k8s.io/kops v1.25.4
	k8s.io/kubelet v0.25.5
	k8s.io/utils v0.0.0-20221128185143-99ec85e7a448
	sigs.k8s.io/mdtoc v1.1.0
	sigs.k8s.io/yaml v1.3.0
)

require (
	4d63.com/gocheckcompilerdirectives v1.2.1 // indirect
	4d63.com/gochecknoglobals v0.2.1 // indirect
	cloud.google.com/go v0.110.7 // indirect
	cloud.google.com/go/compute v1.23.0 // indirect
	cloud.google.com/go/compute/metadata v0.2.3 // indirect
	cloud.google.com/go/iam v1.1.1 // indirect
	cloud.google.com/go/kms v1.15.0 // indirect
	cloud.google.com/go/storage v1.30.1 // indirect
	code.gitea.io/sdk/gitea v0.15.1 // indirect
	github.com/4meepo/tagalign v1.3.2 // indirect
	github.com/Abirdcfly/dupword v0.0.12 // indirect
	github.com/AlekSi/pointer v1.2.0 // indirect
	github.com/Antonboom/errname v0.1.12 // indirect
	github.com/Antonboom/nilnil v0.1.7 // indirect
	github.com/Azure/azure-pipeline-go v0.2.3 // indirect
	github.com/Azure/azure-sdk-for-go v67.1.0+incompatible // indirect
	github.com/Azure/azure-storage-blob-go v0.15.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.28 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.21 // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.11 // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.5 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/BurntSushi/toml v1.3.2 // indirect
	github.com/Djarvur/go-err113 v0.0.0-20210108212216-aea10b59be24 // indirect
	github.com/GaijinEntertainment/go-exhaustruct/v3 v3.1.0 // indirect
	github.com/MakeNowJust/heredoc v1.0.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/Masterminds/sprig/v3 v3.2.3 // indirect
	github.com/Masterminds/squirrel v1.5.3 // indirect
	github.com/Microsoft/go-winio v0.6.0 // indirect
	github.com/OpenPeeDeeP/depguard/v2 v2.1.0 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20211112122917-428f8eabeeb3 // indirect
	github.com/acomagu/bufpipe v1.0.4 // indirect
	github.com/alexkohler/nakedret/v2 v2.0.2 // indirect
	github.com/alexkohler/prealloc v1.0.0 // indirect
	github.com/alingse/asasalint v0.0.11 // indirect
	github.com/apparentlymart/go-cidr v1.1.0 // indirect
	github.com/armon/go-metrics v0.4.1 // indirect
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20200428143746-21a406dcc535 // indirect
	github.com/ashanbrown/forbidigo v1.6.0 // indirect
	github.com/ashanbrown/makezero v1.1.1 // indirect
	github.com/atc0005/go-teams-notify/v2 v2.6.1 // indirect
	github.com/atotto/clipboard v0.1.4 // indirect
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.1 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.13.13 // indirect
	github.com/aws/aws-sdk-go-v2/feature/s3/manager v1.11.3 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.43 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.37 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.45 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.1 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.37 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/pricing v1.17.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/s3 v1.26.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.15.2 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.17.3 // indirect
	github.com/awslabs/goformation/v4 v4.19.5 // indirect
	github.com/aymanbagabas/go-osc52/v2 v2.0.1 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bkielbasa/cyclop v1.2.1 // indirect
	github.com/blakesmith/ar v0.0.0-20190502131153-809d4375e1fb // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/blizzy78/varnamelen v0.8.0 // indirect
	github.com/bombsimon/wsl/v3 v3.4.0 // indirect
	github.com/breml/bidichk v0.2.4 // indirect
	github.com/breml/errchkjson v0.3.1 // indirect
	github.com/butuzov/ireturn v0.2.0 // indirect
	github.com/butuzov/mirror v1.1.0 // indirect
	github.com/caarlos0/ctrlc v1.2.0 // indirect
	github.com/caarlos0/env/v6 v6.10.1 // indirect
	github.com/caarlos0/go-reddit/v3 v3.0.1 // indirect
	github.com/caarlos0/go-shellwords v1.0.12 // indirect
	github.com/caarlos0/log v0.1.6 // indirect
	github.com/cavaliergopher/cpio v1.0.1 // indirect
	github.com/ccojocar/zxcvbn-go v1.0.1 // indirect
	github.com/cenkalti/backoff/v3 v3.0.0 // indirect
	github.com/cenkalti/backoff/v4 v4.1.2 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/chai2010/gettext-go v1.0.2 // indirect
	github.com/charithe/durationcheck v0.0.10 // indirect
	github.com/charmbracelet/bubbles v0.15.0 // indirect
	github.com/charmbracelet/bubbletea v0.24.1 // indirect
	github.com/charmbracelet/lipgloss v0.7.1 // indirect
	github.com/chavacava/garif v0.0.0-20230227094218-b8c73b2037b8 // indirect
	github.com/chigopher/pathlib v0.15.0 // indirect
	github.com/containerd/console v1.0.4-0.20230313162750-1ae8d489ac81 // indirect
	github.com/containerd/containerd v1.6.18 // indirect
	github.com/containerd/stargz-snapshotter/estargz v0.12.1 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/curioswitch/go-reassign v0.2.0 // indirect
	github.com/cyphar/filepath-securejoin v0.2.4 // indirect
	github.com/daixiang0/gci v0.11.0 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/denis-tingaikin/go-header v0.4.3 // indirect
	github.com/dghubble/go-twitter v0.0.0-20211115160449-93a8679adecb // indirect
	github.com/dghubble/oauth1 v0.7.1 // indirect
	github.com/dghubble/sling v1.4.0 // indirect
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/disgoorg/disgo v0.13.20 // indirect
	github.com/disgoorg/log v1.2.0 // indirect
	github.com/disgoorg/snowflake/v2 v2.0.0 // indirect
	github.com/docker/cli v20.10.21+incompatible // indirect
	github.com/docker/distribution v2.8.2+incompatible // indirect
	github.com/docker/docker v20.10.24+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.7.0 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/elazarl/goproxy v0.0.0-20230808193330-2592e75ae04a // indirect
	github.com/emicklei/go-restful/v3 v3.11.0 // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/esimonov/ifshort v1.0.4 // indirect
	github.com/ettle/strcase v0.1.1 // indirect
	github.com/evanphx/json-patch v5.6.0+incompatible // indirect
	github.com/evertras/bubble-table v0.15.2 // indirect
	github.com/exponent-io/jsonpath v0.0.0-20151013193312-d6023ce2651d // indirect
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/firefart/nonamedreturns v1.0.4 // indirect
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	github.com/fzipp/gocyclo v0.6.0 // indirect
	github.com/go-critic/go-critic v0.9.0 // indirect
	github.com/go-errors/errors v1.0.1 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.3.1 // indirect
	github.com/go-git/go-git/v5 v5.4.2 // indirect
	github.com/go-gorp/gorp/v3 v3.0.5 // indirect
	github.com/go-ini/ini v1.67.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible // indirect
	github.com/go-toolsmith/astcast v1.1.0 // indirect
	github.com/go-toolsmith/astcopy v1.1.0 // indirect
	github.com/go-toolsmith/astequal v1.1.0 // indirect
	github.com/go-toolsmith/astfmt v1.1.0 // indirect
	github.com/go-toolsmith/astp v1.1.0 // indirect
	github.com/go-toolsmith/strparse v1.1.0 // indirect
	github.com/go-toolsmith/typep v1.1.0 // indirect
	github.com/go-xmlfmt/xmlfmt v1.1.2 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.4.1 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/golangci/check v0.0.0-20180506172741-cfe4005ccda2 // indirect
	github.com/golangci/dupl v0.0.0-20180902072040-3e9179ac440a // indirect
	github.com/golangci/go-misc v0.0.0-20220329215616-d24fe342adfe // indirect
	github.com/golangci/gofmt v0.0.0-20220901101216-f2edd75033f2 // indirect
	github.com/golangci/lint-1 v0.0.0-20191013205115-297bf364a8e0 // indirect
	github.com/golangci/maligned v0.0.0-20180506175553-b1d89398deca // indirect
	github.com/golangci/misspell v0.4.1 // indirect
	github.com/golangci/revgrep v0.0.0-20220804021717-745bb2f7c2e6 // indirect
	github.com/golangci/unconvert v0.0.0-20180507085042-28b1c447d1f4 // indirect
	github.com/gomarkdown/markdown v0.0.0-20230922112808-5421fefb8386 // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/google/certificate-transparency-go v1.1.4 // indirect
	github.com/google/gnostic v0.6.9 // indirect
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/google/go-containerregistry v0.12.1 // indirect
	github.com/google/go-github/v47 v47.1.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/pprof v0.0.0-20230926050212-f7f687d19a98 // indirect
	github.com/google/s2a-go v0.1.7 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/wire v0.5.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.1 // indirect
	github.com/googleapis/gax-go/v2 v2.12.0 // indirect
	github.com/gophercloud/gophercloud v1.1.1 // indirect
	github.com/gordonklaus/ineffassign v0.0.0-20230610083614-0e73809eb601 // indirect
	github.com/goreleaser/chglog v0.4.2 // indirect
	github.com/goreleaser/fileglob v1.3.0 // indirect
	github.com/goreleaser/nfpm/v2 v2.30.1 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/gostaticanalysis/analysisutil v0.7.1 // indirect
	github.com/gostaticanalysis/comment v1.4.2 // indirect
	github.com/gostaticanalysis/forcetypeassert v0.1.0 // indirect
	github.com/gostaticanalysis/nilerr v0.1.1 // indirect
	github.com/gosuri/uitable v0.0.4 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.5.0 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-plugin v1.4.5 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-secure-stdlib/mlock v0.1.1 // indirect
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.6 // indirect
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/hashicorp/go-uuid v1.0.2 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/vault/api v1.8.2 // indirect
	github.com/hashicorp/vault/sdk v0.6.0 // indirect
	github.com/hashicorp/yamux v0.0.0-20180604194846-3520598351bb // indirect
	github.com/hexops/gotextdiff v1.0.3 // indirect
	github.com/huandu/xstrings v1.4.0 // indirect
	github.com/iancoleman/orderedmap v0.2.0 // indirect
	github.com/iancoleman/strcase v0.3.0 // indirect
	github.com/imdario/mergo v0.3.16 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/invopop/jsonschema v0.7.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jessevdk/go-flags v1.5.0 // indirect
	github.com/jgautheron/goconst v1.5.1 // indirect
	github.com/jingyugao/rowserrcheck v1.1.1 // indirect
	github.com/jinzhu/copier v0.4.0 // indirect
	github.com/jirfag/go-printf-func-name v0.0.0-20200119135958-7558a9eaa5af // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jmoiron/sqlx v1.3.5 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/julz/importas v0.1.0 // indirect
	github.com/kevinburke/rest v0.0.0-20210106114233-22cd0577e450 // indirect
	github.com/kevinburke/ssh_config v1.1.0 // indirect
	github.com/kisielk/errcheck v1.6.3 // indirect
	github.com/kisielk/gotool v1.0.0 // indirect
	github.com/kkHAIKE/contextcheck v1.1.4 // indirect
	github.com/klauspost/compress v1.17.0 // indirect
	github.com/klauspost/pgzip v1.2.6 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/kris-nova/novaarchive v0.0.0-20210219195539-c7c1cabb2577 // indirect
	github.com/kulti/thelper v0.6.3 // indirect
	github.com/kunwardeep/paralleltest v1.0.8 // indirect
	github.com/kyoh86/exportloopref v0.1.11 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/ldez/gomoddirectives v0.2.3 // indirect
	github.com/ldez/tagliatelle v0.5.0 // indirect
	github.com/leonklingele/grouper v1.1.1 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/lufeee/execinquery v1.2.1 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/maratori/testableexamples v1.0.0 // indirect
	github.com/maratori/testpackage v1.1.1 // indirect
	github.com/matoous/godox v0.0.0-20230222163458-006bad1f9d26 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-ieproxy v0.0.1 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-localereader v0.0.1 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/mbilski/exhaustivestruct v1.2.0 // indirect
	github.com/mgechev/revive v1.3.2 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v1.0.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/mmarkdown/mmark v2.0.40+incompatible // indirect
	github.com/moby/locker v1.0.1 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/moby/term v0.0.0-20221205130635-1aeaba878587 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00 // indirect
	github.com/moricho/tparallel v0.3.1 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/muesli/ansi v0.0.0-20211018074035-2e021307bc4b // indirect
	github.com/muesli/cancelreader v0.2.2 // indirect
	github.com/muesli/mango v0.1.0 // indirect
	github.com/muesli/mango-cobra v1.2.0 // indirect
	github.com/muesli/mango-pflag v0.1.0 // indirect
	github.com/muesli/reflow v0.3.0 // indirect
	github.com/muesli/roff v0.1.0 // indirect
	github.com/muesli/termenv v0.15.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/nakabonne/nestif v0.3.1 // indirect
	github.com/nishanths/exhaustive v0.11.0 // indirect
	github.com/nishanths/predeclared v0.2.2 // indirect
	github.com/nunnatsa/ginkgolinter v0.13.5 // indirect
	github.com/oklog/run v1.0.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/oliveagle/jsonpath v0.0.0-20180606110733-2e52cf6e6852 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0-rc3 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pelletier/go-toml/v2 v2.1.0 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pierrec/lz4 v2.5.2+incompatible // indirect
	github.com/pkg/sftp v1.13.5 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/polyfloyd/go-errorlint v1.4.4 // indirect
	github.com/prometheus/client_golang v1.14.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/quasilyte/go-ruleguard v0.4.0 // indirect
	github.com/quasilyte/gogrep v0.5.0 // indirect
	github.com/quasilyte/regex/syntax v0.0.0-20210819130434-b3f0c404a727 // indirect
	github.com/quasilyte/stdinfo v0.0.0-20220114132959-f7386bf02567 // indirect
	github.com/rivo/uniseg v0.4.2 // indirect
	github.com/rs/zerolog v1.31.0 // indirect
	github.com/rubenv/sql-migrate v1.3.1 // indirect
	github.com/russross/blackfriday v1.6.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/ryancurrah/gomodguard v1.3.0 // indirect
	github.com/ryanrolds/sqlclosecheck v0.4.0 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/sagikazarmark/locafero v0.3.0 // indirect
	github.com/sagikazarmark/slog-shim v0.1.0 // indirect
	github.com/sahilm/fuzzy v0.1.0 // indirect
	github.com/sanathkr/go-yaml v0.0.0-20170819195128-ed9d249f429b // indirect
	github.com/sanathkr/yaml v0.0.0-20170819201035-0056894fa522 // indirect
	github.com/sanposhiho/wastedassign/v2 v2.0.7 // indirect
	github.com/sasha-s/go-csync v0.0.0-20210812194225-61421b77c44b // indirect
	github.com/sashamelentyev/interfacebloat v1.1.0 // indirect
	github.com/sashamelentyev/usestdlibvars v1.24.0 // indirect
	github.com/securego/gosec/v2 v2.17.0 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/shazow/go-diff v0.0.0-20160112020656-b6b7b6733b8c // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/sivchari/containedctx v1.0.3 // indirect
	github.com/sivchari/nosnakecase v1.7.0 // indirect
	github.com/sivchari/tenv v1.7.1 // indirect
	github.com/slack-go/slack v0.11.3 // indirect
	github.com/sonatard/noctx v0.0.2 // indirect
	github.com/sourcegraph/conc v0.3.0 // indirect
	github.com/sourcegraph/go-diff v0.7.1-0.20230316160316-1b4d09c1adcb // indirect
	github.com/spf13/cast v1.5.1 // indirect
	github.com/spf13/viper v1.17.0 // indirect
	github.com/spotinst/spotinst-sdk-go v1.133.0 // indirect
	github.com/ssgreg/nlreturn/v2 v2.2.1 // indirect
	github.com/stbenjam/no-sprintf-host-port v0.1.1 // indirect
	github.com/stretchr/objx v0.5.0 // indirect
	github.com/subosito/gotenv v1.6.0 // indirect
	github.com/t-yuki/gocover-cobertura v0.0.0-20180217150009-aaee18c8195c // indirect
	github.com/tdakkota/asciicheck v0.2.0 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	github.com/tetafro/godot v1.4.14 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/timakin/bodyclose v0.0.0-20230421092635-574207250966 // indirect
	github.com/timonwong/loggercheck v0.9.4 // indirect
	github.com/tomarrell/wrapcheck/v2 v2.8.1 // indirect
	github.com/tommy-muehle/go-mnd/v2 v2.5.1 // indirect
	github.com/tomnomnom/linkheader v0.0.0-20180905144013-02ca5825eb80 // indirect
	github.com/ulikunitz/xz v0.5.11 // indirect
	github.com/ultraware/funlen v0.1.0 // indirect
	github.com/ultraware/whitespace v0.0.5 // indirect
	github.com/uudashr/gocognit v1.0.7 // indirect
	github.com/vbatts/tar-split v0.11.2 // indirect
	github.com/voxelbrain/goptions v0.0.0-20180630082107-58cddc247ea2 // indirect
	github.com/withfig/autocomplete-tools/integrations/cobra v1.2.1 // indirect
	github.com/xanzy/go-gitlab v0.73.1 // indirect
	github.com/xanzy/ssh-agent v0.3.1 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/xen0n/gosmopolitan v1.2.1 // indirect
	github.com/xlab/treeprint v1.1.0 // indirect
	github.com/yagipy/maintidx v1.0.0 // indirect
	github.com/yeya24/promlinter v0.2.0 // indirect
	github.com/ykadowak/zerologlint v0.1.3 // indirect
	gitlab.com/bosi/decorder v0.4.0 // indirect
	gitlab.com/digitalxero/go-conventional-commit v1.0.7 // indirect
	go.opencensus.io v0.24.0 // indirect
	go.starlark.net v0.0.0-20200306205701-8dd3e2ee1dd5 // indirect
	go.tmz.dev/musttag v0.7.2 // indirect
	go.uber.org/atomic v1.11.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.24.0 // indirect
	gocloud.dev v0.26.0 // indirect
	golang.org/x/exp/typeparams v0.0.0-20230307190834-24139beb5833 // indirect
	golang.org/x/mod v0.13.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/sys v0.13.0 // indirect
	golang.org/x/term v0.13.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	golang.org/x/time v0.3.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/api v0.143.0 // indirect
	google.golang.org/appengine v1.6.8 // indirect
	google.golang.org/genproto v0.0.0-20230913181813-007df8e322eb // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20230913181813-007df8e322eb // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20230920204549-e6e6cdab5c13 // indirect
	google.golang.org/grpc v1.58.2 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/gcfg.v1 v1.2.3 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.67.0 // indirect
	gopkg.in/mail.v2 v2.3.1 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	honnef.co/go/tools v0.4.5 // indirect
	k8s.io/apiserver v0.26.0 // indirect
	k8s.io/cloud-provider-aws v1.25.1 // indirect
	k8s.io/component-base v0.26.0 // indirect
	k8s.io/csi-translation-lib v0.25.5 // indirect
	k8s.io/gengo v0.0.0-20221011193443-fad74ee6edd9 // indirect
	k8s.io/klog/v2 v2.80.1 // indirect
	k8s.io/kube-openapi v0.0.0-20221012153701-172d655c2280 // indirect
	k8s.io/kubectl v0.26.0 // indirect
	mvdan.cc/gofumpt v0.5.0 // indirect
	mvdan.cc/interfacer v0.0.0-20180901003855-c20040233aed // indirect
	mvdan.cc/lint v0.0.0-20170908181259-adc824a0674b // indirect
	mvdan.cc/unparam v0.0.0-20221223090309-7455f1af531d // indirect
	oras.land/oras-go v1.2.2 // indirect
	sigs.k8s.io/json v0.0.0-20220713155537-f223a00ba0e2 // indirect
	sigs.k8s.io/kustomize/api v0.12.1 // indirect
	sigs.k8s.io/kustomize/kyaml v0.13.9 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)

replace (
	// Used to pin the k8s library versions regardless of what other dependencies enforce
	k8s.io/api => k8s.io/api v0.25.11
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.25.11
	k8s.io/apimachinery => k8s.io/apimachinery v0.25.11
	k8s.io/apiserver => k8s.io/apiserver v0.25.11
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.25.11
	k8s.io/client-go => k8s.io/client-go v0.25.11
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.25.11
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.25.11
	k8s.io/code-generator => k8s.io/code-generator v0.25.11
	k8s.io/component-base => k8s.io/component-base v0.25.11
	k8s.io/component-helpers => k8s.io/component-helpers v0.25.11
	k8s.io/controller-manager => k8s.io/controller-manager v0.25.11
	k8s.io/cri-api => k8s.io/cri-api v0.25.11
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.25.11
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.25.11
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.25.11
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.25.11
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.25.11
	k8s.io/kubectl => k8s.io/kubectl v0.25.5
	k8s.io/kubelet => k8s.io/kubelet v0.25.11
	k8s.io/kubernetes => k8s.io/kubernetes v1.25.11
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.25.11
	k8s.io/metrics => k8s.io/metrics v0.25.11
	k8s.io/mount-utils => k8s.io/mount-utils v0.25.11
	k8s.io/pod-security-admission => k8s.io/pod-security-admission v0.25.11
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.25.11
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.25.11
	k8s.io/sample-controller => k8s.io/sample-controller v0.25.11
)

// Ensure k8s dependencies are also pinned accordingly
replace github.com/acomagu/bufpipe => github.com/acomagu/bufpipe v1.0.4
