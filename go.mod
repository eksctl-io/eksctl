// Make sure to run the following commands after changes to this file are made:
// `make generate-all && make lint && make check-all-generated-files-up-to-date`
// you may also need to run `make push-build-image` depending on what has changed
module github.com/weaveworks/eksctl

go 1.18

require (
	github.com/Masterminds/semver/v3 v3.1.1
	github.com/aws/amazon-ec2-instance-selector/v2 v2.3.1
	github.com/aws/aws-sdk-go v1.44.49
	github.com/aws/aws-sdk-go-v2 v1.16.7
	github.com/aws/aws-sdk-go-v2/config v1.15.13
	github.com/aws/aws-sdk-go-v2/credentials v1.12.8
	github.com/aws/aws-sdk-go-v2/service/autoscaling v1.23.5
	github.com/aws/aws-sdk-go-v2/service/cloudformation v1.22.0
	github.com/aws/aws-sdk-go-v2/service/cloudtrail v1.16.4
	github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs v1.15.10
	github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider v1.17.0
	github.com/aws/aws-sdk-go-v2/service/ec2 v1.47.2
	github.com/aws/aws-sdk-go-v2/service/eks v1.21.4
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing v1.14.8
	github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2 v1.18.8
	github.com/aws/aws-sdk-go-v2/service/iam v1.18.9
	github.com/aws/aws-sdk-go-v2/service/kms v1.17.2
	github.com/aws/aws-sdk-go-v2/service/ssm v1.27.4
	github.com/aws/aws-sdk-go-v2/service/sts v1.16.9
	github.com/aws/smithy-go v1.12.0
	github.com/benjamintf1/unmarshalledmatchers v1.0.0
	github.com/blang/semver v3.5.1+incompatible
	github.com/bxcodec/faker v2.0.1+incompatible
	github.com/cenk/backoff v2.2.1+incompatible
	github.com/cloudflare/cfssl v1.6.1
	github.com/dave/dst v0.27.0
	github.com/dave/jennifer v1.5.0
	github.com/dlespiau/kube-test-harness v0.0.0-20200915102055-a03579200ae8
	github.com/evanphx/json-patch/v5 v5.6.0
	github.com/fatih/color v1.13.0
	github.com/github-release/github-release v0.10.0
	github.com/gobwas/glob v0.2.3
	github.com/gofrs/flock v0.8.1
	github.com/golangci/golangci-lint v1.45.2
	github.com/google/uuid v1.3.0
	github.com/goreleaser/goreleaser v1.6.1
	github.com/hashicorp/go-version v1.6.0
	github.com/kballard/go-shellquote v0.0.0-20180428030007-95032a82bc51
	github.com/kris-nova/logger v0.2.2
	github.com/kris-nova/lolgopher v0.0.0-20210112022122-73f0047e8b65
	github.com/kubicorn/kubicorn v0.0.0-20191114212505-a2c64ce430b9
	github.com/lithammer/dedent v1.1.0
	github.com/maxbrunsfeld/counterfeiter/v6 v6.4.1
	github.com/onsi/ginkgo/v2 v2.1.4
	github.com/onsi/gomega v1.19.0
	github.com/orcaman/concurrent-map v1.0.0
	github.com/otiai10/copy v1.7.0
	github.com/pelletier/go-toml v1.9.5
	github.com/pkg/errors v0.9.1
	github.com/sethvargo/go-password v0.2.0
	github.com/spf13/afero v1.8.2
	github.com/spf13/cobra v1.5.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.0
	github.com/tidwall/gjson v1.14.1
	github.com/tidwall/sjson v1.2.4
	github.com/tj/assert v0.0.3
	github.com/vburenin/ifacemaker v1.2.1-0.20220402060045-b2018d8549dc
	github.com/vektra/mockery v1.1.2
	github.com/weaveworks/goformation/v4 v4.10.2-0.20211208101807-d5ec5126726c
	github.com/weaveworks/launcher v0.0.2-0.20200715141516-1ca323f1de15
	github.com/weaveworks/schemer v0.0.0-20210802122110-338b258ad2ca
	github.com/xgfone/netaddr v0.5.1
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f
	golang.org/x/tools v0.1.11
	gopkg.in/yaml.v1 v1.0.0-20140924161607-9f9df34309c0
	gopkg.in/yaml.v2 v2.4.0
	helm.sh/helm/v3 v3.9.0
	k8s.io/api v0.24.2
	k8s.io/apiextensions-apiserver v0.24.2
	k8s.io/apimachinery v0.24.2
	k8s.io/cli-runtime v0.24.2
	k8s.io/client-go v1.5.2
	k8s.io/cloud-provider v0.24.2
	k8s.io/code-generator v0.22.1
	k8s.io/kops v1.23.2
	k8s.io/kubelet v0.24.2
	k8s.io/legacy-cloud-providers v0.24.2
	sigs.k8s.io/mdtoc v1.1.0
	sigs.k8s.io/yaml v1.3.0
)

require (
	4d63.com/gochecknoglobals v0.1.0 // indirect
	bitbucket.org/creachadair/shell v0.0.7 // indirect
	cloud.google.com/go v0.102.0 // indirect
	cloud.google.com/go/compute v1.7.0 // indirect
	cloud.google.com/go/kms v1.1.0 // indirect
	cloud.google.com/go/storage v1.22.1 // indirect
	code.gitea.io/sdk/gitea v0.15.1 // indirect
	github.com/AlekSi/pointer v1.2.0 // indirect
	github.com/Antonboom/errname v0.1.5 // indirect
	github.com/Antonboom/nilnil v0.1.0 // indirect
	github.com/Azure/azure-pipeline-go v0.2.3 // indirect
	github.com/Azure/azure-sdk-for-go v60.2.0+incompatible // indirect
	github.com/Azure/azure-storage-blob-go v0.15.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/Azure/go-autorest v14.2.0+incompatible // indirect
	github.com/Azure/go-autorest/autorest v0.11.23 // indirect
	github.com/Azure/go-autorest/autorest/adal v0.9.18 // indirect
	github.com/Azure/go-autorest/autorest/azure/auth v0.5.10 // indirect
	github.com/Azure/go-autorest/autorest/azure/cli v0.4.4 // indirect
	github.com/Azure/go-autorest/autorest/date v0.3.0 // indirect
	github.com/Azure/go-autorest/autorest/to v0.4.0 // indirect
	github.com/Azure/go-autorest/autorest/validation v0.3.1 // indirect
	github.com/Azure/go-autorest/logger v0.2.1 // indirect
	github.com/Azure/go-autorest/tracing v0.6.0 // indirect
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/DisgoOrg/disgohook v1.4.4 // indirect
	github.com/DisgoOrg/log v1.1.2 // indirect
	github.com/DisgoOrg/restclient v1.2.8 // indirect
	github.com/Djarvur/go-err113 v0.0.0-20210108212216-aea10b59be24 // indirect
	github.com/MakeNowJust/heredoc v1.0.0 // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/Masterminds/sprig/v3 v3.2.2 // indirect
	github.com/Masterminds/squirrel v1.5.3 // indirect
	github.com/Microsoft/go-winio v0.5.2 // indirect
	github.com/OpenPeeDeeP/depguard v1.1.0 // indirect
	github.com/ProtonMail/go-crypto v0.0.0-20211112122917-428f8eabeeb3 // indirect
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/acomagu/bufpipe v1.0.3 // indirect
	github.com/alecthomas/jsonschema v0.0.0-20211209230136-e2b41affa5c1 // indirect
	github.com/alexkohler/prealloc v1.0.0 // indirect
	github.com/apex/log v1.9.0 // indirect
	github.com/apparentlymart/go-cidr v1.1.0 // indirect
	github.com/armon/go-metrics v0.4.0 // indirect
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d // indirect
	github.com/ashanbrown/forbidigo v1.3.0 // indirect
	github.com/ashanbrown/makezero v1.1.1 // indirect
	github.com/atc0005/go-teams-notify/v2 v2.6.1 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.8 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.14 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.8 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.15 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.8 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.11 // indirect
	github.com/awslabs/goformation/v4 v4.19.5 // indirect
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bgentry/speakeasy v0.1.0 // indirect
	github.com/bkielbasa/cyclop v1.2.0 // indirect
	github.com/blakesmith/ar v0.0.0-20190502131153-809d4375e1fb // indirect
	github.com/blang/semver/v4 v4.0.0 // indirect
	github.com/blizzy78/varnamelen v0.6.1 // indirect
	github.com/bombsimon/wsl/v3 v3.3.0 // indirect
	github.com/breml/bidichk v0.2.2 // indirect
	github.com/breml/errchkjson v0.2.3 // indirect
	github.com/butuzov/ireturn v0.1.1 // indirect
	github.com/caarlos0/ctrlc v1.0.0 // indirect
	github.com/caarlos0/env/v6 v6.9.1 // indirect
	github.com/caarlos0/go-reddit/v3 v3.0.1 // indirect
	github.com/caarlos0/go-shellwords v1.0.12 // indirect
	github.com/cavaliergopher/cpio v1.0.1 // indirect
	github.com/cenkalti/backoff/v3 v3.2.2 // indirect
	github.com/cenkalti/backoff/v4 v4.1.2 // indirect
	github.com/census-instrumentation/opencensus-proto v0.3.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/charithe/durationcheck v0.0.9 // indirect
	github.com/chavacava/garif v0.0.0-20210405164556-e8a0a408d6af // indirect
	github.com/cncf/udpa/go v0.0.0-20220112060539-c52dc94e7fbe // indirect
	github.com/cncf/xds/go v0.0.0-20220520190051-1e77728a1eaa // indirect
	github.com/containerd/containerd v1.6.6 // indirect
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd/v22 v22.3.2 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/cyphar/filepath-securejoin v0.2.3 // indirect
	github.com/daixiang0/gci v0.3.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/denis-tingaikin/go-header v0.4.3 // indirect
	github.com/denverdino/aliyungo v0.0.0-20220610083100-ab5f747cb559 // indirect
	github.com/dghubble/go-twitter v0.0.0-20211115160449-93a8679adecb // indirect
	github.com/dghubble/oauth1 v0.7.1 // indirect
	github.com/dghubble/sling v1.4.0 // indirect
	github.com/dimchansky/utfbom v1.1.1 // indirect
	github.com/docker/cli v20.10.17+incompatible // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.17+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.6.4 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/elazarl/goproxy v0.0.0-20190421051319-9d40249d3c2f // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/envoyproxy/go-control-plane v0.10.3 // indirect
	github.com/envoyproxy/protoc-gen-validate v0.6.7 // indirect
	github.com/esimonov/ifshort v1.0.4 // indirect
	github.com/ettle/strcase v0.1.1 // indirect
	github.com/evanphx/json-patch v5.6.0+incompatible // indirect
	github.com/exponent-io/jsonpath v0.0.0-20210407135951-1de76d718b3f // indirect
	github.com/fatih/structtag v1.2.0 // indirect
	github.com/form3tech-oss/jwt-go v3.2.5+incompatible // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/fullstorydev/grpcurl v1.8.6 // indirect
	github.com/fzipp/gocyclo v0.4.0 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-critic/go-critic v0.6.2 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/go-git/gcfg v1.5.0 // indirect
	github.com/go-git/go-billy/v5 v5.3.1 // indirect
	github.com/go-git/go-git/v5 v5.4.2 // indirect
	github.com/go-gorp/gorp/v3 v3.0.2 // indirect
	github.com/go-ini/ini v1.66.6 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/spec v0.20.6 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible // indirect
	github.com/go-toolsmith/astcast v1.0.0 // indirect
	github.com/go-toolsmith/astcopy v1.0.0 // indirect
	github.com/go-toolsmith/astequal v1.0.1 // indirect
	github.com/go-toolsmith/astfmt v1.0.0 // indirect
	github.com/go-toolsmith/astp v1.0.0 // indirect
	github.com/go-toolsmith/strparse v1.0.0 // indirect
	github.com/go-toolsmith/typep v1.0.2 // indirect
	github.com/go-xmlfmt/xmlfmt v0.0.0-20191208150333-d5b6f63a941b // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.2.0 // indirect
	github.com/golang/glog v1.0.0 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/mock v1.6.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/golangci/check v0.0.0-20180506172741-cfe4005ccda2 // indirect
	github.com/golangci/dupl v0.0.0-20180902072040-3e9179ac440a // indirect
	github.com/golangci/go-misc v0.0.0-20180628070357-927a3d87b613 // indirect
	github.com/golangci/gofmt v0.0.0-20190930125516-244bba706f1a // indirect
	github.com/golangci/lint-1 v0.0.0-20191013205115-297bf364a8e0 // indirect
	github.com/golangci/maligned v0.0.0-20180506175553-b1d89398deca // indirect
	github.com/golangci/misspell v0.3.5 // indirect
	github.com/golangci/revgrep v0.0.0-20210930125155-c22e5001d4f2 // indirect
	github.com/golangci/unconvert v0.0.0-20180507085042-28b1c447d1f4 // indirect
	github.com/gomarkdown/markdown v0.0.0-20210514010506-3b9f47219fe7 // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/google/certificate-transparency-go v1.1.3 // indirect
	github.com/google/gnostic v0.6.9 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/google/go-containerregistry v0.10.0 // indirect
	github.com/google/go-github/v43 v43.0.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/rpmpack v0.0.0-20211125064518-d0ed9b1b61b9 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/trillian v1.4.1 // indirect
	github.com/google/wire v0.5.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.1.0 // indirect
	github.com/googleapis/gax-go/v2 v2.4.0 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/gophercloud/gophercloud v0.25.0 // indirect
	github.com/gordonklaus/ineffassign v0.0.0-20210914165742-4cc7213b9bc8 // indirect
	github.com/goreleaser/chglog v0.1.2 // indirect
	github.com/goreleaser/fileglob v1.3.0 // indirect
	github.com/goreleaser/nfpm/v2 v2.14.0 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/gostaticanalysis/analysisutil v0.7.1 // indirect
	github.com/gostaticanalysis/comment v1.4.2 // indirect
	github.com/gostaticanalysis/forcetypeassert v0.1.0 // indirect
	github.com/gostaticanalysis/nilerr v0.1.1 // indirect
	github.com/gosuri/uitable v0.0.4 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0 // indirect
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.10.3 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-hclog v1.2.1 // indirect
	github.com/hashicorp/go-immutable-radix v1.3.1 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-plugin v1.4.4 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.1 // indirect
	github.com/hashicorp/go-rootcerts v1.0.2 // indirect
	github.com/hashicorp/go-secure-stdlib/mlock v0.1.2 // indirect
	github.com/hashicorp/go-secure-stdlib/parseutil v0.1.6 // indirect
	github.com/hashicorp/go-secure-stdlib/strutil v0.1.2 // indirect
	github.com/hashicorp/go-sockaddr v1.0.2 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hashicorp/golang-lru v0.5.4 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/hashicorp/vault/api v1.7.2 // indirect
	github.com/hashicorp/vault/sdk v0.5.2 // indirect
	github.com/hashicorp/yamux v0.0.0-20211028200310-0bc27b27de87 // indirect
	github.com/hexops/gotextdiff v1.0.3 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/iancoleman/orderedmap v0.2.0 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/jessevdk/go-flags v1.5.0 // indirect
	github.com/jgautheron/goconst v1.5.1 // indirect
	github.com/jhump/protoreflect v1.12.0 // indirect
	github.com/jingyugao/rowserrcheck v1.1.1 // indirect
	github.com/jirfag/go-printf-func-name v0.0.0-20200119135958-7558a9eaa5af // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/jmoiron/sqlx v1.3.5 // indirect
	github.com/jonboulle/clockwork v0.3.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/julz/importas v0.1.0 // indirect
	github.com/kevinburke/rest v0.0.0-20210106114233-22cd0577e450 // indirect
	github.com/kevinburke/ssh_config v1.1.0 // indirect
	github.com/kisielk/errcheck v1.6.0 // indirect
	github.com/kisielk/gotool v1.0.0 // indirect
	github.com/klauspost/compress v1.15.7 // indirect
	github.com/klauspost/pgzip v1.2.5 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/kris-nova/novaarchive v0.0.0-20210219195539-c7c1cabb2577 // indirect
	github.com/kulti/thelper v0.5.1 // indirect
	github.com/kunwardeep/paralleltest v1.0.3 // indirect
	github.com/kyoh86/exportloopref v0.1.8 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/ldez/gomoddirectives v0.2.2 // indirect
	github.com/ldez/tagliatelle v0.3.1 // indirect
	github.com/leonklingele/grouper v1.1.0 // indirect
	github.com/lib/pq v1.10.6 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/magiconair/properties v1.8.6 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/maratori/testpackage v1.0.1 // indirect
	github.com/matoous/godox v0.0.0-20210227103229-6504466cf951 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-ieproxy v0.0.7 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mbilski/exhaustivestruct v1.2.0 // indirect
	github.com/mgechev/revive v1.1.4 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-homedir v1.1.0 // indirect
	github.com/mitchellh/go-testing-interface v1.14.1 // indirect
	github.com/mitchellh/go-wordwrap v1.0.1 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/mmarkdown/mmark v2.0.40+incompatible // indirect
	github.com/moby/locker v1.0.1 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00 // indirect
	github.com/moricho/tparallel v0.2.1 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/muesli/coral v1.0.0 // indirect
	github.com/muesli/mango v0.1.0 // indirect
	github.com/muesli/mango-coral v1.0.1 // indirect
	github.com/muesli/mango-pflag v0.1.0 // indirect
	github.com/muesli/roff v0.1.0 // indirect
	github.com/nakabonne/nestif v0.3.1 // indirect
	github.com/nbutton23/zxcvbn-go v0.0.0-20210217022336-fa2cb2858354 // indirect
	github.com/nishanths/exhaustive v0.7.11 // indirect
	github.com/nishanths/predeclared v0.2.1 // indirect
	github.com/oklog/run v1.1.0 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.3-0.20220114050600-8b9d41f48198 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/phayes/checkstyle v0.0.0-20170904204023-bfd46e6a821d // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/pkg/sftp v1.13.5 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/polyfloyd/go-errorlint v0.0.0-20211125173453-6d6d39c5bb8b // indirect
	github.com/prometheus/client_golang v1.12.2 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.35.0 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	github.com/quasilyte/go-ruleguard v0.3.15 // indirect
	github.com/quasilyte/gogrep v0.0.0-20220103110004-ffaa07af02e3 // indirect
	github.com/quasilyte/regex/syntax v0.0.0-20200407221936-30656e2c4a95 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/rubenv/sql-migrate v1.1.2 // indirect
	github.com/russross/blackfriday v1.6.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/ryancurrah/gomodguard v1.2.3 // indirect
	github.com/ryanrolds/sqlclosecheck v0.3.0 // indirect
	github.com/ryanuber/go-glob v1.0.0 // indirect
	github.com/sanathkr/go-yaml v0.0.0-20170819195128-ed9d249f429b // indirect
	github.com/sanathkr/yaml v0.0.0-20170819201035-0056894fa522 // indirect
	github.com/sanposhiho/wastedassign/v2 v2.0.6 // indirect
	github.com/securego/gosec/v2 v2.10.0 // indirect
	github.com/sergi/go-diff v1.2.0 // indirect
	github.com/shazow/go-diff v0.0.0-20160112020656-b6b7b6733b8c // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/sivchari/containedctx v1.0.2 // indirect
	github.com/sivchari/tenv v1.4.7 // indirect
	github.com/slack-go/slack v0.10.2 // indirect
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/sonatard/noctx v0.0.1 // indirect
	github.com/sourcegraph/go-diff v0.6.1 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.10.1 // indirect
	github.com/spotinst/spotinst-sdk-go v1.123.0 // indirect
	github.com/ssgreg/nlreturn/v2 v2.2.1 // indirect
	github.com/stretchr/objx v0.4.0 // indirect
	github.com/subosito/gotenv v1.2.0 // indirect
	github.com/sylvia7788/contextcheck v1.0.4 // indirect
	github.com/tdakkota/asciicheck v0.1.1 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
	github.com/tetafro/godot v1.4.11 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/timakin/bodyclose v0.0.0-20210704033933-f49887972144 // indirect
	github.com/tmc/grpc-websocket-proxy v0.0.0-20220101234140-673ab2c3ae75 // indirect
	github.com/tomarrell/wrapcheck/v2 v2.5.0 // indirect
	github.com/tommy-muehle/go-mnd/v2 v2.5.0 // indirect
	github.com/tomnomnom/linkheader v0.0.0-20180905144013-02ca5825eb80 // indirect
	github.com/transparency-dev/merkle v0.0.1 // indirect
	github.com/ulikunitz/xz v0.5.10 // indirect
	github.com/ultraware/funlen v0.0.3 // indirect
	github.com/ultraware/whitespace v0.0.5 // indirect
	github.com/urfave/cli v1.22.9 // indirect
	github.com/uudashr/gocognit v1.0.5 // indirect
	github.com/voxelbrain/goptions v0.0.0-20180630082107-58cddc247ea2 // indirect
	github.com/xanzy/go-gitlab v0.56.0 // indirect
	github.com/xanzy/ssh-agent v0.3.1 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/xiang90/probing v0.0.0-20190116061207-43a291ad63a2 // indirect
	github.com/xlab/treeprint v1.1.0 // indirect
	github.com/yagipy/maintidx v1.0.0 // indirect
	github.com/yeya24/promlinter v0.1.1-0.20210918184747-d757024714a1 // indirect
	gitlab.com/bosi/decorder v0.2.1 // indirect
	go.etcd.io/bbolt v1.3.6 // indirect
	go.etcd.io/etcd/api/v3 v3.5.4 // indirect
	go.etcd.io/etcd/client/pkg/v3 v3.5.4 // indirect
	go.etcd.io/etcd/client/v2 v2.305.4 // indirect
	go.etcd.io/etcd/client/v3 v3.5.4 // indirect
	go.etcd.io/etcd/etcdctl/v3 v3.5.4 // indirect
	go.etcd.io/etcd/etcdutl/v3 v3.5.4 // indirect
	go.etcd.io/etcd/pkg/v3 v3.5.4 // indirect
	go.etcd.io/etcd/raft/v3 v3.5.4 // indirect
	go.etcd.io/etcd/server/v3 v3.5.4 // indirect
	go.etcd.io/etcd/tests/v3 v3.5.4 // indirect
	go.etcd.io/etcd/v3 v3.5.4 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.32.0 // indirect
	go.opentelemetry.io/otel v1.7.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp v0.20.0 // indirect
	go.opentelemetry.io/otel/internal/metric v0.27.0 // indirect
	go.opentelemetry.io/otel/metric v0.30.0 // indirect
	go.opentelemetry.io/otel/sdk v1.7.0 // indirect
	go.opentelemetry.io/otel/sdk/export/metric v0.28.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v0.30.0 // indirect
	go.opentelemetry.io/otel/trace v1.7.0 // indirect
	go.opentelemetry.io/proto/otlp v0.18.0 // indirect
	go.starlark.net v0.0.0-20220328144851-d1966c6b9fcd // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	go.uber.org/zap v1.21.0 // indirect
	gocloud.dev v0.24.0 // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/net v0.0.0-20220706163947-c90051bbdb60 // indirect
	golang.org/x/oauth2 v0.0.0-20220630143837-2104d58473e0 // indirect
	golang.org/x/sys v0.0.0-20220704084225-05e143d24a9e // indirect
	golang.org/x/term v0.0.0-20220526004731-065cf7ba2467 // indirect
	golang.org/x/text v0.3.7 // indirect
	golang.org/x/time v0.0.0-20220609170525-579cf78fd858 // indirect
	golang.org/x/xerrors v0.0.0-20220609144429-65e65417b02f // indirect
	google.golang.org/api v0.86.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20220706185917-7780775163c4 // indirect
	google.golang.org/grpc v1.47.0 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/cheggaaa/pb.v1 v1.0.28 // indirect
	gopkg.in/gcfg.v1 v1.2.3 // indirect
	gopkg.in/gorp.v1 v1.7.2 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/ini.v1 v1.66.6 // indirect
	gopkg.in/mail.v2 v2.3.1 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
	gopkg.in/warnings.v0 v0.1.2 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	honnef.co/go/tools v0.2.2 // indirect
	k8s.io/apiserver v0.24.2 // indirect
	k8s.io/component-base v0.24.2 // indirect
	k8s.io/csi-translation-lib v0.24.2 // indirect
	k8s.io/gengo v0.0.0-20210813121822-485abfe95c7c // indirect
	k8s.io/klog/v2 v2.70.1 // indirect
	k8s.io/kube-openapi v0.0.0-20220627174259-011e075b9cb8 // indirect
	k8s.io/kubectl v0.24.2 // indirect
	k8s.io/utils v0.0.0-20220706174534-f6158b442e7c // indirect
	mvdan.cc/gofumpt v0.3.0 // indirect
	mvdan.cc/interfacer v0.0.0-20180901003855-c20040233aed // indirect
	mvdan.cc/lint v0.0.0-20170908181259-adc824a0674b // indirect
	mvdan.cc/unparam v0.0.0-20211214103731-d0ef000c54e5 // indirect
	oras.land/oras-go v1.2.0 // indirect
	sigs.k8s.io/kustomize/api v0.11.5 // indirect
	sigs.k8s.io/kustomize/kyaml v0.13.7 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
)

replace (
	// Used to pin the k8s library versions regardless of what other dependencies enforce
	k8s.io/api => k8s.io/api v0.21.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.21.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.21.2
	k8s.io/apiserver => k8s.io/apiserver v0.21.2
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.21.2
	k8s.io/client-go => k8s.io/client-go v0.21.2
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.21.2
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.21.2
	k8s.io/code-generator => k8s.io/code-generator v0.21.2
	k8s.io/component-base => k8s.io/component-base v0.21.2
	k8s.io/cri-api => k8s.io/cri-api v0.21.2
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.21.2
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.21.2
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.21.2
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.21.2
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.21.2
	k8s.io/kubectl => k8s.io/kubectl v0.21.2
	k8s.io/kubelet => k8s.io/kubelet v0.21.2
	k8s.io/kubernetes => k8s.io/kubernetes v1.19.5
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.21.2
	k8s.io/metrics => k8s.io/metrics v0.21.2
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.21.2
)
