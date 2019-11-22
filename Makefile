include Makefile.common

version_pkg := github.com/weaveworks/eksctl/pkg/version

gopath := $(shell go env GOPATH)
gocache := $(shell go env GOCACHE)

GOBIN ?= $(gopath)/bin


always_generate_in_packages := ./pkg/nodebootstrap ./pkg/addons/default ./pkg/addons

generated_code_deep_copy_helper := pkg/apis/eksctl.io/v1alpha5/zz_generated.deepcopy.go

generated_code_aws_sdk_mocks := $(wildcard pkg/eks/mocks/*API.go)

conditionally_generated_files := \
  site/content/usage/20-schema.md \
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
	CGO_ENABLED=0 time go build -ldflags "-X $(version_pkg).gitCommit=$(git_commit) -X $(version_pkg).builtAt=$(built_at)" ./cmd/eksctl

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
	time "$(GOBIN)/gometalinter" ./pkg/... ./cmd/... ./integration/...

.PHONY: test
test:
	$(MAKE) lint
	$(MAKE) check-all-generated-files-up-to-date
	$(MAKE) unit-test
	$(MAKE) build-integration-test

.PHONY: unit-test
unit-test: ## Run unit test only
	CGO_ENABLED=0 time go test ./pkg/... ./cmd/... $(UNIT_TEST_ARGS)

.PHONY: unit-test-race
unit-test-race: ## Run unit test with race detection
	CGO_ENABLED=1 time go test -race ./pkg/... ./cmd/... $(UNIT_TEST_ARGS)

.PHONY: build-integration-test
build-integration-test: $(all_generated_code) ## Build integration test binary
	time go test -tags integration ./integration/ -c -o eksctl-integration-test

.PHONY: integration-test
integration-test: build build-integration-test ## Run the integration tests (with cluster creation and cleanup)
	cd integration; ../eksctl-integration-test -test.timeout 120m $(INTEGRATION_TEST_ARGS)

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
		-eksctl.create=false \
		-eksctl.delete=false \
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

site/content/usage/20-schema.md: $(call godeps,cmd/schema/generate.go)
	time go run ./cmd/schema/generate.go $@

deep_copy_helper_input = $(shell $(call godeps_cmd,./pkg/apis/...) | sed 's|$(generated_code_deep_copy_helper)||' )
$(generated_code_deep_copy_helper): $(deep_copy_helper_input) .license-header ## Generate Kubernetes API helpers
	time env GOPATH="$(gopath)" bash "$(gopath)/pkg/mod/k8s.io/code-generator@v0.0.0-20190831074504-732c9ca86353/generate-groups.sh" \
	  deepcopy,defaulter _ ./pkg/apis eksctl.io:v1alpha5 --go-header-file .license-header --output-base="$(git_toplevel)" \
	  || (cat .license-header ; cat $(generated_code_deep_copy_helper); exit 1)

$(generated_code_aws_sdk_mocks): $(call godeps,pkg/eks/mocks/mocks.go)
	mkdir -p vendor/github.com/aws/
	@# Hack for Mockery to find the dependencies handled by `go mod`
	ln -sfn "$(gopath)/pkg/mod/github.com/weaveworks-experiments/aws-sdk-go@v1.25.14-0.20191118174145-b6a927b24aa9" vendor/github.com/aws/aws-sdk-go
	time env GOBIN=$(GOBIN) go generate ./pkg/eks/mocks

##@ Docker

.PHONY: eksctl-image
eksctl-image: ## Build the eksctl image that has release artefacts and no build dependencies
	$(MAKE) -f Makefile.docker $@

##@ Site

HUGO := $(GOBIN)/hugo
HUGO_ARGS ?= --gc --minify

.PHONY: serve-pages
serve-pages: ## Serve the site locally
	cd site/ ; $(HUGO) serve $(HUGO_ARGS)

.PHONY: build-pages
build-pages: ## Generate the site
	cd site/ ; $(HUGO) $(HUGO_ARGS)

##@ Utility

.PHONY: help
help:  ## Display this help. Thanks to https://suva.sh/posts/well-documented-makefiles/
ifeq ($(OS),Windows_NT)
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <target>\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-40s %s\n", $$1, $$2 } /^##@/ { printf "\n%s\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
else
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-40s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
endif
