built_at := $(shell date +%s)
git_commit := $(shell git describe --dirty --always)

version_pkg := github.com/weaveworks/eksctl/pkg/version

EKSCTL_BUILD_IMAGE ?= weaveworks/eksctl-build:latest
EKSCTL_IMAGE ?= weaveworks/eksctl:latest

GO_BUILD_TAGS ?= netgo

GOBIN ?= $(shell echo `go env GOPATH`/bin)

.DEFAULT_GOAL := help

##@ Dependencies

.PHONY: install-build-deps
install-build-deps: ## Install dependencies (packages and tools)
	@cd build && ./install.sh

##@ Build

.PHONY: build
build: ## Build eksctl
	CGO_ENABLED=0 go build -tags "$(GO_BUILD_TAGS)" -ldflags "-X $(version_pkg).gitCommit=$(git_commit) -X $(version_pkg).builtAt=$(built_at)" ./cmd/eksctl

##@ Testing & CI

ifneq ($(TEST_V),)
UNIT_TEST_ARGS ?= -v -ginkgo.v
INTEGRATION_TEST_ARGS ?= -test.v -ginkgo.v
endif

ifneq ($(INTEGRATION_TEST_FOCUS),)
INTEGRATION_TEST_ARGS ?= -test.v -ginkgo.v -ginkgo.focus "$(INTEGRATION_TEST_FOCUS)"
endif

LINTER ?= gometalinter ./pkg/... ./cmd/... ./integration/...
.PHONY: lint
lint: ## Run linter over the codebase
	@"$(GOBIN)/$(LINTER)"

.PHONY: test
test: generate ## Run unit test (and re-generate code under test)
	@$(MAKE) lint
	@git diff --exit-code pkg/nodebootstrap/assets.go > /dev/null || (git --no-pager diff pkg/nodebootstrap/assets.go; exit 1)
	@git diff --exit-code ./pkg/eks/mocks > /dev/null || (git --no-pager diff ./pkg/eks/mocks; exit 1)
	@git diff --exit-code ./pkg/addons/default > /dev/null || (git --no-pager diff ./pkg/addons/default; exit 1)
	@$(MAKE) unit-test
	@test -z $(COVERALLS_TOKEN) || "$(GOBIN)/goveralls" -coverprofile=coverage.out -service=circle-ci
	@$(MAKE) build-integration-test

.PHONY: unit-test
unit-test: ## Run unit test only
	@CGO_ENABLED=0 go test -covermode=count -coverprofile=coverage.out ./pkg/... ./cmd/... $(UNIT_TEST_ARGS)

.PHONY: unit-test-race
unit-test-race: ## Run unit test with race detection
	@CGO_ENABLED=1 go test -race -covermode=atomic -coverprofile=coverage.out ./pkg/... ./cmd/... $(UNIT_TEST_ARGS)

.PHONY: build-integration-test
build-integration-test: ## Build integration test binary
	@go test -tags integration ./integration/... -c -o ./eksctl-integration-test

.PHONY: integration-test
integration-test: build build-integration-test ## Run the integration tests (with cluster creation and cleanup)
	@cd integration; ../eksctl-integration-test -test.timeout 60m \
		$(INTEGRATION_TEST_ARGS)

.PHONY: integration-test-container
integration-test-container: eksctl-image ## Run the integration tests inside a Docker container
	$(MAKE) integration-test-container-pre-built

.PHONY: integration-test-container-pre-built
integration-test-container-pre-built: ## Run the integration tests inside a Docker container
	@docker run \
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
	@./eksctl utils write-kubeconfig \
		--auto-kubeconfig \
		--name=$(TEST_CLUSTER)
	$(info it is recommended to watch events with "kubectl get events --watch --all-namespaces --kubeconfig=$(HOME)/.kube/eksctl/clusters/$(TEST_CLUSTER)")
	@cd integration ; ../eksctl-integration-test -test.timeout 21m \
		$(INTEGRATION_TEST_ARGS) \
		-eksctl.cluster=$(TEST_CLUSTER) \
		-eksctl.create=false \
		-eksctl.delete=false \
		-eksctl.kubeconfig=$(HOME)/.kube/eksctl/clusters/$(TEST_CLUSTER)

create-integration-test-dev-cluster: build ## Create a test cluster for use when developing integration tests
	@./eksctl create cluster --name=integration-test-dev --auto-kubeconfig --nodes=1 --nodegroup-name=ng-0

delete-integration-test-dev-cluster: build ## Delete the test cluster for use when developing integration tests
	@./eksctl delete cluster --name=integration-test-dev --auto-kubeconfig

##@ Code Generation

.PHONY: generate
generate: ## Generate code
	@chmod g-w  ./pkg/nodebootstrap/assets/*
	@go generate ./pkg/nodebootstrap ./pkg/eks/mocks ./pkg/addons/default

.PHONY: generate-ami
generate-ami: ## Generate the list of AMIs for use with static resolver. Queries AWS.
	@go generate ./pkg/ami

.PHONY: ami-check
ami-check: generate-ami ## Check whether the AMIs have been updated and fail if they have. Designed for a automated test
	@git diff --exit-code pkg/ami/static_resolver_ami.go > /dev/null || (git --no-pager diff; exit 1)


.PHONY: generate-kubernetes-types
generate-kubernetes-types: ## Generate Kubernetes API helpers
	@build/vendor/k8s.io/code-generator/generate-groups.sh deepcopy,defaulter _ github.com/weaveworks/eksctl/pkg/apis eksctl.io:v1alpha5

##@ Docker

.PHONY: eksctl-build-image
eksctl-build-image: ## Create the the eksctl build docker image
	@-docker pull $(EKSCTL_BUILD_IMAGE)
	@docker build --tag=$(EKSCTL_BUILD_IMAGE) --cache-from=$(EKSCTL_BUILD_IMAGE) ./build

EKSCTL_IMAGE_BUILD_ARGS := \
  --build-arg=EKSCTL_BUILD_IMAGE=$(EKSCTL_BUILD_IMAGE) \
  --build-arg=GO_BUILD_TAGS=$(GO_BUILD_TAGS)

ifneq ($(COVERALLS_TOKEN),)
EKSCTL_IMAGE_BUILD_ARGS += --build-arg=COVERALLS_TOKEN=$(COVERALLS_TOKEN)
endif
ifneq ($(JUNIT_REPORT_DIR),)
EKSCTL_IMAGE_BUILD_ARGS += --build-arg=JUNIT_REPORT_DIR=$(JUNIT_REPORT_DIR)
endif
ifeq ($(OS),Windows_NT)
EKSCTL_IMAGE_BUILD_ARGS += --build-arg=TEST_TARGET=unit-test
else
EKSCTL_IMAGE_BUILD_ARGS += --build-arg=TEST_TARGET=test
endif

.PHONY: eksctl-image
eksctl-image: eksctl-build-image ## Create the eksctl image
	@docker build --tag=$(EKSCTL_IMAGE) $(EKSCTL_IMAGE_BUILD_ARGS) ./
	./get-testresults.sh

##@ Release

.PHONY: release
release: eksctl-build-image ## Create a new eksctl release
	@docker run \
	  --env=GITHUB_TOKEN \
	  --env=CIRCLE_TAG \
	  --env=CIRCLE_PROJECT_USERNAME \
	  --volume=$(CURDIR):/go/src/github.com/weaveworks/eksctl \
	  --workdir=/go/src/github.com/weaveworks/eksctl \
	    $(EKSCTL_BUILD_IMAGE) \
	      ./do-release.sh

JEKYLL := docker run --tty --rm \
  --name=eksctl-jekyll \
  --volume="$(CURDIR)":/usr/src/app \
  --publish="4000:4000" \
    starefossen/github-pages

##@ Site

.PHONY: serve-pages
serve-pages: ## Serve the site locally
	@-docker rm -f eksctl-jekyll
	@$(JEKYLL) jekyll serve -d /_site --watch --force_polling -H 0.0.0.0 -P 4000

.PHONY: build-pages
build-pages: ## Generate the site using jekyll
	@-docker rm -f eksctl-jekyll
	@$(JEKYLL) jekyll build --verbose

##@ Utility

.PHONY: help
help:  ## Display this help. Thanks to https://suva.sh/posts/well-documented-makefiles/
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
