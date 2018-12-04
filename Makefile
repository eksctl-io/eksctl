built_at := $(shell date +%s)
git_commit := $(shell git describe --dirty --always)

version_pkg := github.com/weaveworks/eksctl/pkg/version

EKSCTL_BUILD_IMAGE ?= weaveworks/eksctl:build
EKSCTL_IMAGE ?= weaveworks/eksctl:latest

.DEFAULT_GOAL := help

ifneq ($(TEST_V),)
TEST_ARGS ?= -v -ginkgo.v
endif

##@ Dependencies

.PHONY: install-build-deps
install-build-deps: ## Install dependencies (packages and tools)
	@cd build && dep ensure && ./install.sh

##@ Build

.PHONY: build
build: ## Build eksctl
	@go build -ldflags "-X $(version_pkg).gitCommit=$(git_commit) -X $(version_pkg).builtAt=$(built_at)" ./cmd/eksctl

##@ Testing & CI

.PHONY: test
test: generate ## Run unit test (and re-generate code under test)
	@git diff --exit-code pkg/nodebootstrap/assets.go > /dev/null || (git --no-pager diff; exit 1)
	@git diff --exit-code ./pkg/eks/mocks > /dev/null || (git --no-pager diff; exit 1)
	@$(MAKE) unit-test
	@test -z $(COVERALLS_TOKEN) || $(GOPATH)/bin/goveralls -coverprofile=coverage.out -service=circle-ci

.PHONY: unit-test
unit-test: ## Run unit test only
	@CGO_ENABLED=0 go test -covermode=count -coverprofile=coverage.out ./pkg/... ./cmd/... $(TEST_ARGS)

LINTER ?= gometalinter ./...
.PHONY: lint
lint: ## Run linter over the codebase
	@$(GOPATH)/bin/$(LINTER)

.PHONY: ci
ci: test lint ## Target for CI system to invoke to run tests and linting

TEST_CLUSTER ?= integration-test-dev
.PHONY: integration-test-dev
integration-test-dev: build ## Run the integration tests without cluster teardown. For use when developing integration tests.
	@./eksctl utils write-kubeconfig \
		--auto-kubeconfig \
		--name=$(TEST_CLUSTER)
	@go test -tags integration -timeout 21m ./integration/... \
		$(TEST_ARGS) \
		-args \
		-eksctl.cluster=$(TEST_CLUSTER) \
		-eksctl.create=false \
		-eksctl.delete=false \
		-eksctl.kubeconfig=$(HOME)/.kube/eksctl/clusters/$(TEST_CLUSTER)

create-integration-test-dev-cluster: build ## Create a test cluster for use when developing integration tests
	@./eksctl create cluster --name=integration-test-dev --auto-kubeconfig

delete-integration-test-dev-cluster: build ## Delete the test cluster for use when developing integration tests
	@./eksctl delete cluster --name=integration-test-dev --auto-kubeconfig

.PHONY: integration-test
integration-test: build ## Run the integration tests (with cluster creation and cleanup)
	@go test -tags integration -timeout 60m ./integration/... $(TEST_ARGS)

##@ Code Generation

.PHONY: generate
generate: ## Generate code
	@chmod g-w  ./pkg/nodebootstrap/assets/*
	@go generate ./pkg/nodebootstrap ./pkg/eks/mocks

.PHONY: generate-ami
generate-ami: ## Generate the list of AMIs for use with static resolver. Queries AWS.
	@go generate ./pkg/ami

.PHONY: ami-check
ami-check: generate-ami  ## Check whether the AMIs have been updated and fail if they have. Designed for a automated test
	@git diff --exit-code pkg/ami/static_resolver_ami.go > /dev/null || (git --no-pager diff; exit 1)

##@ Docker

.PHONY: eksctl-build-image
eksctl-build-image: ## Create the the eksctl build docker image
	@-docker pull $(EKSCTL_BUILD_IMAGE)
	@docker build --tag=$(EKSCTL_BUILD_IMAGE) --cache-from=$(EKSCTL_BUILD_IMAGE) ./build

EKSCTL_IMAGE_BUILD_ARGS := --build-arg=EKSCTL_BUILD_IMAGE=$(EKSCTL_BUILD_IMAGE)
ifneq ($(COVERALLS_TOKEN),)
EKSCTL_IMAGE_BUILD_ARGS += --build-arg=COVERALLS_TOKEN=$(COVERALLS_TOKEN)
endif
ifneq ($(JUNIT_REPORT_FOLDER),)
EKSCTL_IMAGE_BUILD_ARGS += --build-arg=JUNIT_REPORT_FOLDER=$(JUNIT_REPORT_FOLDER)
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

.PHONY: server-pages
serve-pages: ## Serve the site locally
	@-docker rm -f eksctl-jekyll
	@$(JEKYLL) jekyll serve -d /_site --watch --force_polling -H 0.0.0.0 -P 4000

.PHONY: build-page
build-pages: ## Generate the site using jekyll
	@-docker rm -f eksctl-jekyll
	@$(JEKYLL) jekyll build --verbose

##@ Utility

.PHONY: help
help:  ## Display this help. Thanks to https://suva.sh/posts/well-documented-makefiles/
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
