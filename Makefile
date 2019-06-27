built_at := $(shell date +%s)
git_commit := $(shell git describe --dirty --always)

version_pkg := github.com/weaveworks/eksctl/pkg/version

# The dependencies version should be bumped every time the build dependencies are updated
EKSCTL_DEPENDENCIES_IMAGE ?= weaveworks/eksctl-build:deps-0.4
EKSCTL_BUILDER_IMAGE ?= weaveworks/eksctl-builder:latest
EKSCTL_IMAGE ?= weaveworks/eksctl:latest

GOBIN ?= $(shell echo `go env GOPATH`/bin)

.DEFAULT_GOAL := help

##@ Dependencies

.PHONY: install-build-deps
install-build-deps: ## Install dependencies (packages and tools)
	./install-build-deps.sh

##@ Build

.PHONY: build
build: generate-bindata-assets generate-kubernetes-types  ## Build eksctl
	CGO_ENABLED=0 time go build -v -ldflags "-X $(version_pkg).gitCommit=$(git_commit) -X $(version_pkg).builtAt=$(built_at)" ./cmd/eksctl

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
test: ## Run unit test (and re-generate code under test)
	$(MAKE) lint
	$(MAKE) generate-aws-mocks-test generate-bindata-assets-test generate-kubernetes-types-test
	$(MAKE) unit-test
	$(MAKE) build-integration-test

.PHONY: unit-test
unit-test: ## Run unit test only
	CGO_ENABLED=0 time go test ./pkg/... ./cmd/... $(UNIT_TEST_ARGS)

.PHONY: unit-test-race
unit-test-race: ## Run unit test with race detection
	CGO_ENABLED=1 time go test -race ./pkg/... ./cmd/... $(UNIT_TEST_ARGS)

.PHONY: build-integration-test
build-integration-test: ## Build integration test binary
	time go test -tags integration ./integration/... -c -o ./eksctl-integration-test

.PHONY: integration-test
integration-test: build build-integration-test ## Run the integration tests (with cluster creation and cleanup)
	cd integration; ../eksctl-integration-test -test.timeout 60m \
		$(INTEGRATION_TEST_ARGS)

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
integration-test-dev: build build-integration-test ## Run the integration tests without cluster teardown. For use when developing integration tests.
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

.PHONY: generate-bindata-assets
generate-bindata-assets: ## Generate bindata assets (node bootstrap config files & add-on manifests)
	chmod g-w  ./pkg/nodebootstrap/assets/*
	env GOBIN=$(GOBIN) time go generate ./pkg/nodebootstrap ./pkg/addons/default

.PHONY: generate-bindata-assets-test
generate-bindata-assets-test: generate-bindata-assets ## Test if generated bindata assets are checked-in
	git diff --exit-code ./pkg/nodebootstrap/assets.go > /dev/null || (git --no-pager diff ./pkg/nodebootstrap/assets.go; exit 1)
	git diff --exit-code ./pkg/addons/default/assets.go > /dev/null || (git --no-pager diff ./pkg/addons/default/assets.go; exit 1)

.license-header: LICENSE
	@# generate-groups.sh can't find the lincense header when using Go modules, so we provide one
	printf "/*\n%s\n*/\n" "$$(cat LICENSE)" > $@

.PHONY: generate-kubernetes-types
generate-kubernetes-types: .license-header ## Generate Kubernetes API helpers
	time go mod download k8s.io/code-generator # make sure the code-generator is present
	time env GOPATH="$$(go env GOPATH)" bash "$$(go env GOPATH)/pkg/mod/k8s.io/code-generator@v0.0.0-20190612205613-18da4a14b22b/generate-groups.sh" \
	  deepcopy,defaulter pkg/apis ./pkg/apis eksctl.io:v1alpha5 --go-header-file .license-header --output-base="$${PWD}" \
	  || (cat codegenheader.txt ; cat pkg/apis/eksctl.io/v1alpha5/zz_generated.deepcopy.go ; exit 1)

.PHONY: generate-kubernetes-types-test
generate-kubernetes-types-test: generate-kubernetes-types ## Test if generated Kubernetes API helpers are checked-in
	git diff --exit-code ./pkg/nodebootstrap/assets.go > /dev/null || (git --no-pager diff ./pkg/nodebootstrap/assets.go; exit 1)

.PHONY: generate-ami
generate-ami: ## Generate the list of AMIs for use with static resolver. Queries AWS.
	time go generate ./pkg/ami

.PHONY: generate-schema
generate-schema: ## Generate the schema file in the documentation site.
	@go run ./cmd/schema/generate.go

.PHONY: ami-check
ami-check: generate-ami ## Check whether the AMIs have been updated and fail if they have. Designed for a automated test
	git diff --exit-code pkg/ami/static_resolver_ami.go > /dev/null || (git --no-pager diff; exit 1)

.PHONY: generate-aws-mocks
generate-aws-mocks: ## Generate mocks for AWS SDK
	mkdir -p vendor/github.com/aws/
	@# Hack for Mockery to find the dependencies handled by `go mod`
	ln -sfn "$$(go env GOPATH)/pkg/mod/github.com/aws/aws-sdk-go@v1.19.18" vendor/github.com/aws/aws-sdk-go
	time env GOBIN=$(GOBIN) go generate ./pkg/eks/mocks

.PHONY: generate-aws-mocks-test
generate-aws-mocks-test: generate-aws-mocks ## Test if generated mocks for AWS SDK are checked-in
	git diff --exit-code ./pkg/eks/mocks > /dev/null || (git --no-pager diff ./pkg/eks/mocks; exit 1)

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

.PHONY: release
release: eksctl-build-image ## Create a new eksctl release
	docker run \
	  --env=GITHUB_TOKEN \
	  --env=CIRCLE_TAG \
	  --env=CIRCLE_PROJECT_USERNAME \
	  --volume=$(CURDIR):/src \
	  --workdir=/src \
	    $(EKSCTL_BUILDER_IMAGE) \
	      ./do-release.sh

JEKYLL := docker run --tty --rm \
  --name=eksctl-jekyll \
  --volume="$(CURDIR)":/usr/src/app \
  --publish="4000:4000" \
    starefossen/github-pages

##@ Site

.PHONY: serve-pages
serve-pages: ## Serve the site locally
	-docker rm -f eksctl-jekyll
	$(JEKYLL) jekyll serve -d /_site --watch --force_polling -H 0.0.0.0 -P 4000

.PHONY: build-pages
build-pages: ## Generate the site using jekyll
	-docker rm -f eksctl-jekyll
	$(JEKYLL) jekyll build --verbose

##@ Utility

.PHONY: help
help:  ## Display this help. Thanks to https://suva.sh/posts/well-documented-makefiles/
	awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
