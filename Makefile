built_at := $(shell date +%s)
git_commit := $(shell git describe --dirty --always)
version_pkg := github.com/weaveworks/eksctl/pkg/version

# The dependencies version should be bumped every time the build dependencies are updated
EKSCTL_DEPENDENCIES_IMAGE ?= weaveworks/eksctl-build:deps-0.7
EKSCTL_BUILDER_IMAGE ?= weaveworks/eksctl-builder:latest
EKSCTL_IMAGE ?= weaveworks/eksctl:latest

GOBIN ?= $(shell echo `go env GOPATH`/bin)

AWS_SDK_MOCKS := $(wildcard pkg/eks/mocks/*API.go)

DEEP_COPY_HELPER := pkg/apis/eksctl.io/v1alpha5/zz_generated.deepcopy.go
GENERATED_GO_FILES := pkg/addons/default/assets.go \
pkg/nodebootstrap/assets.go \
$(DEEP_COPY_HELPER) \
pkg/ami/static_resolver_ami.go \
$(AWS_SDK_MOCKS)

GENERATED_GO_FILES := $(GENERATED_FILES) site/content/usage/20-schema.md


.DEFAULT_GOAL := help

##@ Dependencies

.PHONY: install-build-deps
install-build-deps: ## Install dependencies (packages and tools)
	./install-build-deps.sh

##@ Build

godeps_cmd = go list -deps -f '{{if not .Standard}}{{ $$dep := . }}{{range .GoFiles}}{{$$dep.Dir}}/{{.}} {{end}}{{end}}' $(1) | sed "s%${PWD}/%%g"
godeps = $(shell $(call godeps_cmd,$(1)))

.PHONY: build
build: $(GENERATED_GO_FILES) ## Build main binary
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

.PHONY: lint
lint: ## Run linter over the codebase
	time "$(GOBIN)/gometalinter" ./pkg/... ./cmd/... ./integration/...

.PHONY: test
test:
	$(MAKE) lint
	$(MAKE) check-generated-sources-up-to-date
	$(MAKE) unit-test
	$(MAKE) build-integration-test

.PHONY: unit-test
unit-test: ## Run unit test only
	CGO_ENABLED=0 time go test ./pkg/... ./cmd/... $(UNIT_TEST_ARGS)

.PHONY: unit-test-race
unit-test-race: ## Run unit test with race detection
	CGO_ENABLED=1 time go test -race ./pkg/... ./cmd/... $(UNIT_TEST_ARGS)

.PHONY: build-integration-test
build-integration-test: $(GENERATED_GO_FILES) ## Build integration test binary
	time go test -tags integration ./integration/... -c -o eksctl-integration-test

.PHONY: integration-test
integration-test: eksctl-integration-test ## Run the integration tests (with cluster creation and cleanup)
	cd integration; ../eksctl-integration-test -test.timeout 60m $(INTEGRATION_TEST_ARGS)

.PHONY: integration-test-container
integration-test-container: eksctl-image ## Run the integration tests inside a Docker container
	$(MAKE) integration-test-container-pre-built

.PHONY: integration-test-container-pre-built
integration-test-container-pre-built: ## Run the integration tests inside a Docker container
	docker run \
	  --env=AWS_PROFILE \
	  --volume=$(HOME)/.aws:/root/.aws \
	  --workdir=/usr/local/share/eksctl \
	    $(EKSCTL_IMAGE) \
		  eksctl-integration-test \
		    -eksctl.path=/usr/local/bin/eksctl \
			-eksctl.kubeconfig=/tmp/kubeconfig \
			  $(INTEGRATION_TEST_ARGS)

TEST_CLUSTER ?= integration-test-dev
.PHONY: integration-test-dev
integration-test-dev: eksctl-integration-test ## Run the integration tests without cluster teardown. For use when developing integration tests.
	./eksctl utils write-kubeconfig \
		--auto-kubeconfig \
		--name=$(TEST_CLUSTER)
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

.PHONY: regenerate-sources
# TODO: generate-ami is broken (see https://github.com/weaveworks/eksctl/issues/949 ), include it when fixed
regenerate-sources: $(GENERATED_FILES) # generate-ami ## Re-generate all the automatically-generated source files

.PHONY: check-generated-files-up-to-date
check-generated-sources-up-to-date: regenerate-sources
	git diff --quiet -- $(GENERATED_FILES) || (git --no-pager diff $(GENERATED_FILES); exit 1)

pkg/addons/default/assets.go: pkg/addons/default/assets/*
	env GOBIN=$(GOBIN) time go generate ./$(@D)

pkg	/nodebootstrap/assets.go: pkg/nodebootstrap/assets/*
	chmod g-w $^
	env GOBIN=$(GOBIN) time go generate ./$(@D)

.license-header: LICENSE
	@# generate-groups.sh can't find the lincense header when using Go modules, so we provide one
	printf "/*\n%s\n*/\n" "$$(cat LICENSE)" > $@

DEEP_COPY_DEPS := $(shell $(call godeps_cmd,./pkg/apis/...) | sed 's%$(DEEP_COPY_HELPER)%%' )
$(DEEP_COPY_HELPER): $(DEEP_COPY_DEPS) .license-header ## Generate Kubernetes API helpers
	time go mod download k8s.io/code-generator # make sure the code-generator is present
	time env GOPATH="$$(go env GOPATH)" bash "$$(go env GOPATH)/pkg/mod/k8s.io/code-generator@v0.0.0-20190612205613-18da4a14b22b/generate-groups.sh" \
	  deepcopy,defaulter _ ./pkg/apis eksctl.io:v1alpha5 --go-header-file .license-header --output-base="$${PWD}" \
	  || (cat codegenheader.txt ; cat $(DEEP_COPY_HELPER); exit 1)

# static_resolver_ami.go doesn't only depend on files (it should be refreshed whenever a release is made in AWS)
# so we need to forcicly generate it
.PHONY: generate-ami
generate-ami: ## Generate the list of AMIs for use with static resolver. Queries AWS.
	time go generate ./pkg/ami

site/content/usage/20-schema.md: $(call godeps,cmd/schema/generate.go)
	time go run ./cmd/schema/generate.go $@

$(AWS_SDK_MOCKS): $(call godeps,pkg/eks/mocks/mocks.go)
	mkdir -p vendor/github.com/aws/
	@# Hack for Mockery to find the dependencies handled by `go mod`
	ln -sfn "$$(go env GOPATH)/pkg/mod/github.com/aws/aws-sdk-go@v1.19.18" vendor/github.com/aws/aws-sdk-go
	time env GOBIN=$(GOBIN) go generate ./pkg/eks/mocks

##@ Docker

ifeq ($(OS),Windows_NT)
TEST_TARGET=unit-test
else
TEST_TARGET=test
endif

.PHONY: eksctl-deps-image
eksctl-deps-image: ## Create a cache image with dependencies
	-time docker pull $(EKSCTL_DEPENDENCIES_IMAGE)
	@# Pin dependency file permissions.
	@# Docker uses the file permissions as part of the COPY hash, which can lead to cache misses
	@# in hosts with different default file permissions (umask).
	chmod 0700 install-build-deps.sh
	chmod 0600 go.mod go.sum
	time docker build --cache-from=$(EKSCTL_DEPENDENCIES_IMAGE) --tag=$(EKSCTL_DEPENDENCIES_IMAGE) -f Dockerfile.deps .

.PHONY: eksctl-image
eksctl-image: eksctl-deps-image## Create the eksctl image
	time docker run -t --name eksctl-builder -e TEST_TARGET=$(TEST_TARGET) \
	  -v "${PWD}":/src -v "$$(go env GOCACHE):/root/.cache/go-build" -v "$$(go env GOPATH)/pkg/mod:/go/pkg/mod" \
          $(EKSCTL_DEPENDENCIES_IMAGE) /src/eksctl-image-builder.sh || ( docker rm eksctl-builder; exit 1 )
	time docker commit eksctl-builder $(EKSCTL_BUILDER_IMAGE) && docker rm eksctl-builder
	docker build --tag $(EKSCTL_IMAGE) .


##@ Release

docker_run_release_script = docker run \
  --env=GITHUB_TOKEN \
  --env=CIRCLE_TAG \
  --env=CIRCLE_PROJECT_USERNAME \
  --volume=$(CURDIR):/src \
  --workdir=/src \
    $(EKSCTL_BUILDER_IMAGE)

.PHONY: release-candidate
release-candidate: eksctl-image ## Create a new eksctl release candidate
	$(call docker_run_release_script) ./do-release-candidate.sh

.PHONY: release
release: eksctl-image ## Create a new eksctl release
	$(call docker_run_release_script) ./do-release.sh

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
	awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
