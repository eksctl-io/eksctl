include Makefile.common

version_pkg := github.com/weaveworks/eksctl/pkg/version

gopath := $(shell go env GOPATH)
gocache := $(shell go env GOCACHE)

GOBIN ?= $(gopath)/bin


always_generate_in_packages := ./pkg/nodebootstrap ./pkg/addons/default ./pkg/addons

generated_code_deep_copy_helper := pkg/apis/eksctl.io/v1alpha5/zz_generated.deepcopy.go

generated_code_aws_sdk_mocks := $(wildcard pkg/eks/mocks/*API.go)

conditionally_generated_files := \
  userdocs/src/usage/schema.md \
  $(generated_code_deep_copy_helper) $(generated_code_aws_sdk_mocks)

all_generated_files := \
  pkg/nodebootstrap/assets.go \
  pkg/nodebootstrap/maxpods.go \
  pkg/addons/default/assets.go \
  pkg/addons/default/assets/aws-node.yaml \
  pkg/addons/assets.go \
  pkg/ami/static_resolver_ami.go \
  $(conditionally_generated_files)

.DEFAULT_GOAL := help

##@ Dependencies

.PHONY: install-build-deps
install-build-deps: ## Install dependencies (packages and tools)
	./install-build-deps.sh

##@ Build

godeps_cmd = go list -deps -f '{{if not .Standard}}{{ $$dep := . }}{{range .GoFiles}}{{$$dep.Dir}}/{{.}} {{end}}{{end}}' $(1) | sed "s|$(git_toplevel)/||g"
godeps = $(shell $(call godeps_cmd,$(1)))

.PHONY: build
build: generate-always ## Build main binary
	CGO_ENABLED=0 time go build -ldflags "-X $(version_pkg).gitCommit=$(git_commit) -X $(version_pkg).buildDate=$(build_date)" ./cmd/eksctl

# Build binaries for Linux, Windows and Mac and place them in dist/
.PHONY: build-all
build-all: generate-always
	goreleaser --config=.goreleaser-local.yaml --snapshot --skip-publish --rm-dist

##@ Testing & CI

ifneq ($(TEST_V),)
UNIT_TEST_ARGS ?= -v -ginkgo.v
INTEGRATION_TEST_ARGS ?= -test.v -ginkgo.v
endif

ifneq ($(INTEGRATION_TEST_FOCUS),)
INTEGRATION_TEST_ARGS ?= -test.v -ginkgo.v -ginkgo.focus "$(INTEGRATION_TEST_FOCUS)"
endif

ifneq ($(INTEGRATION_TEST_REGION),)
INTEGRATION_TEST_ARGS += -eksctl.region=$(INTEGRATION_TEST_REGION)
$(info will launch integration tests in region $(INTEGRATION_TEST_REGION))
endif

ifneq ($(INTEGRATION_TEST_VERSION),)
INTEGRATION_TEST_ARGS += -eksctl.version=$(INTEGRATION_TEST_VERSION)
$(info will launch integration tests for Kubernetes version $(INTEGRATION_TEST_VERSION))
endif

ifneq ($(SSH_KEY_PATH),)
INTEGRATION_TEST_ARGS += -eksctl.git.sshkeypath=$(SSH_KEY_PATH)
$(info will launch integration tests with ssh key path $(SSH_KEY_PATH))
endif

.PHONY: lint
lint: ## Run linter over the codebase
	time "$(GOBIN)/golangci-lint" run

.PHONY: test
test:
	$(MAKE) lint
	$(MAKE) check-all-generated-files-up-to-date
	$(MAKE) unit-test
	$(MAKE) build-integration-test

.PHONY: unit-test
unit-test: ## Run unit test only
	CGO_ENABLED=0 time go test  -tags=release ./pkg/... ./cmd/... $(UNIT_TEST_ARGS)

.PHONY: unit-test-race
unit-test-race: ## Run unit test with race detection
	CGO_ENABLED=1 time go test -race ./pkg/... ./cmd/... $(UNIT_TEST_ARGS)

.PHONY: build-integration-test
build-integration-test: $(all_generated_code)
	@# Compile integration test binary without running any.
	@# Required as build failure aren't listed when running go build below. See also: https://github.com/golang/go/issues/15513
	go test -tags integration -run=^$$ ./integration/...
	@# Build integration test binary:
	go build -tags integration -o ./eksctl-integration-test ./integration/main.go

.PHONY: integration-test
integration-test: build build-integration-test ## Run the integration tests (with cluster creation and cleanup)
	JUNIT_REPORT_DIR=$(git_toplevel)/test-results ./eksctl-integration-test $(INTEGRATION_TEST_ARGS)

.PHONY: integration-test-container
integration-test-container: eksctl-image ## Run the integration tests inside a Docker container
	$(MAKE) integration-test-container-pre-built

.PHONY: integration-test-container-pre-built
integration-test-container-pre-built: ## Run the integration tests inside a Docker container
	docker run \
	  --env=AWS_PROFILE \
	  --volume=$(HOME)/.aws:/root/.aws \
	  --volume=$(HOME)/.ssh:/root/.ssh \
	  --workdir=/usr/local/share/eksctl \
	    $(eksctl_image_name) \
		  eksctl-integration-test \
		    -eksctl.path=/usr/local/bin/eksctl \
			-eksctl.kubeconfig=/tmp/kubeconfig \
			  $(INTEGRATION_TEST_ARGS)

TEST_CLUSTER ?= integration-test-dev
.PHONY: integration-test-dev
integration-test-dev: build-integration-test ## Run the integration tests without cluster teardown. For use when developing integration tests.
	./eksctl utils write-kubeconfig \
		--auto-kubeconfig \
		--cluster=$(TEST_CLUSTER)
	$(info it is recommended to watch events with "kubectl get events --watch --all-namespaces --kubeconfig=$(HOME)/.kube/eksctl/clusters/$(TEST_CLUSTER)")
	cd integration ; ../eksctl-integration-test -test.timeout 21m \
		$(INTEGRATION_TEST_ARGS) \
		-eksctl.cluster=$(TEST_CLUSTER) \
		-eksctl.skip.create=true \
		-eksctl.skip.delete=true \
		-eksctl.kubeconfig=$(HOME)/.kube/eksctl/clusters/$(TEST_CLUSTER)

create-integration-test-dev-cluster: build ## Create a test cluster for use when developing integration tests
	./eksctl create cluster --name=integration-test-dev --auto-kubeconfig --nodes=1 --nodegroup-name=ng-0

delete-integration-test-dev-cluster: build ## Delete the test cluster for use when developing integration tests
	./eksctl delete cluster --name=integration-test-dev --auto-kubeconfig

##@ Code Generation

.PHONY: generate-always
generate-always: ## Generate code (required for every build)
	@# go-bindata targets must run every time, as dependencies are too complex to declare in make:
	@# - deleting an asset is breaks the dependencies
	@# - different version of go-bindata generate different code
	@$(GOBIN)/go-bindata -v
	env GOBIN=$(GOBIN) time go generate $(always_generate_in_packages)

.PHONY: generate-all
generate-all: generate-always $(conditionally_generated_files) ## Re-generate all the automatically-generated source files

.PHONY: check-all-generated-files-up-to-date
check-all-generated-files-up-to-date: generate-all
	git diff --quiet -- $(all_generated_files) || (git --no-pager diff $(all_generated_files); echo "HINT: to fix this, run 'git commit $(all_generated_files) --message \"Update generated files\"'"; exit 1)

.license-header: LICENSE
	@# generate-groups.sh can't find the lincense header when using Go modules, so we provide one
	printf "/*\n%s\n*/\n" "$$(cat LICENSE)" > $@

.PHONY: generate-ami
generate-ami: ## Generate the list of AMIs for use with static resolver. Queries AWS.
	time go generate ./pkg/ami

userdocs/src/usage/schema.md: $(call godeps,cmd/schema/generate.go)
	time go run ./cmd/schema/generate.go $@

deep_copy_helper_input = $(shell $(call godeps_cmd,./pkg/apis/...) | sed 's|$(generated_code_deep_copy_helper)||' )
$(generated_code_deep_copy_helper): $(deep_copy_helper_input) .license-header ## Generate Kubernetes API helpers
	./tools/update-codegen.sh

$(generated_code_aws_sdk_mocks): $(call godeps,pkg/eks/mocks/mocks.go)
	mkdir -p vendor/github.com/aws/
	@# Hack for Mockery to find the dependencies handled by `go mod`
	ln -sfn "$(gopath)/pkg/mod/github.com/weaveworks/aws-sdk-go@v1.25.14-0.20191218135223-757eeed07291" vendor/github.com/aws/aws-sdk-go
	time env GOBIN=$(GOBIN) go generate ./pkg/eks/mocks

.PHONY: generate-kube-reserved
generate-kube-reserved: ## Update instance list with respective specs
	@cd ./pkg/nodebootstrap/ && go run reserved_generate.go

##@ Release
.PHONY: prepare-release
prepare-release:
	./tag-release.sh

.PHONY: prepare-release-candidate
prepare-release-candidate:
	./tag-release-candidate.sh

.PHONY: print-version
print-version:
	@go run pkg/version/generate/release_generate.go print-version

.PHONY: upload-github
upload-github:
	@echo "Releasing version $(eksctl_version) in $(git_org)/$(git_repo)"
	@echo "Check draft exists..." && github-release info --user $(git_org) --repo $(git_repo) --tag $(eksctl_version)
	github-release upload --user $(git_org) --repo $(git_repo) --tag $(eksctl_version) --file dist/eksctl_Windows_amd64.zip --name eksctl_Windows_amd64.zip
	github-release upload --user $(git_org) --repo $(git_repo) --tag $(eksctl_version) --file dist/eksctl_Darwin_amd64.tar.gz --name eksctl_Darwin_amd64.tar.gz
	github-release upload --user $(git_org) --repo $(git_repo) --tag $(eksctl_version) --file dist/eksctl_Linux_amd64.tar.gz --name eksctl_Linux_amd64.tar.gz

.PHONY: publish-github
publish-github: upload-github
	github-release publish --user $(git_org) --repo $(git_repo) --tag $(eksctl_version)

.PHONY: publish-rc-github
publish-rc-github: upload-github
	github-release release --user $(git_org) --repo $(git_repo) --tag $(shell eksctl version) --pre-release

.PHONY: publish-homebrew
publish-homebrew:
	@echo "Publishing to weaveworks/homebrew-tap"
	git clone --depth 1 --branch master git@github.com:weaveworks/homebrew-tap.git
	@go run tools/brew/update_formula.go \
		-template tools/brew/formula.tmpl \
		-outputPath homebrew-tap/Formula/$(git_repo).rb \
		-version $(eksctl_version) \
		-linux-url https://github.com/$(git_org)/$(git_repo)/releases/download/$(eksctl_version)/eksctl_Linux_amd64.tar.gz \
		-mac-url https://github.com/$(git_org)/$(git_repo)/releases/download/$(eksctl_version)/eksctl_Darwin_amd64.tar.gz
	cd homebrew-tap; git commit --message "Brew formula update for $(git_repo) version $(eksctl_version)" -- Formula/$(git_repo).rb
	cd homebrew-tap; git push origin master

##@ Docker

.PHONY: eksctl-image
eksctl-image: ## Build the eksctl image that has release artefacts and no build dependencies
	$(MAKE) -f Makefile.docker $@

##@ Site

docs-deps:
	pip3 install -r userdocs/requirements.txt

.PHONY: serve-pages
serve-pages: docs-deps ## Serve the site locally
	cd userdocs/ ; mkdocs serve

.PHONY: build-pages
build-pages: docs-deps ## Generate the site
	cd userdocs/ ; mkdocs build

##@ Utility

.PHONY: help
help:  ## Display this help. Thanks to https://suva.sh/posts/well-documented-makefiles/
ifeq ($(OS),Windows_NT)
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <target>\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-40s %s\n", $$1, $$2 } /^##@/ { printf "\n%s\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
else
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-40s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
endif
