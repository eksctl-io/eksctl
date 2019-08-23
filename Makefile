built_at := $(shell date +%s)
git_commit = $(shell git describe --dirty --always)

git_toplevel := $(shell git rev-parse --show-toplevel)
version_pkg := github.com/weaveworks/eksctl/pkg/version

build_image_input := Dockerfile.base Dockerfile install-build-deps.sh go.mod go.sum

# We use git object hashes for determining the input files used for any given build image
# E.g. given `weaveworks/eksctl-build:e6e8800773d3adf8e7999a23dcdb07414c66a4da` one can
# run `git show e6e8800773d3adf8e7999a23dcdb07414c66a4da` to get contents of `.build_image_manifest`,
# and `git show <hash>` for each of the hashes in the manifest to determine contents of each of the
# files used in `$(build_image_input)` at the time.
base_build_image_tag = base-$(shell git ls-tree --full-tree @ -- Dockerfile.base | awk '{ print $$3 }')
build_image_tag = $(shell git ls-tree --full-tree @ -- .build_image_manifest | awk '{ print $$3 }')

base_build_image_name = weaveworks/eksctl-build:$(base_build_image_tag)
build_image_name = weaveworks/eksctl-build:$(build_image_tag)

unique_tag = $(shell printf "%s-%s-%s" `git rev-parse @` $(build_image_tag) $(built_at))

intermidiate_container_name = eksctl-build-$(unique_tag)
intermediate_image_name = weaveworks/eksctl-intermediate:$(unique_tag)

eksctl_image_name ?= weaveworks/eksctl:latest

gopath := $(shell go env GOPATH)
gocache := $(shell go env GOCACHE)

docker_build := time docker build
# We should eventually switch to buildkit, as it has a many feature and cleaner output with timing info,
# but right now 'docker build' doesn't allow us to export build cache images, so we cannot use it yet
# docker_build := env DOCKER_BUILDKIT=1 $(docker_build)

GOBIN ?= $(gopath)/bin

ifeq ($(OS),Windows_NT)
TEST_TARGET=unit-test
else
TEST_TARGET=test
endif

generated_code_aws_sdk_mocks := $(wildcard pkg/eks/mocks/*API.go)

generated_code_deep_copy_helper := pkg/apis/eksctl.io/v1alpha5/zz_generated.deepcopy.go

all_generated_code := \
  pkg/addons/default/assets.go \
  pkg/nodebootstrap/assets.go \
  pkg/nodebootstrap/maxpods.go \
  pkg/addons/default/assets/aws-node.yaml \
  pkg/ami/static_resolver_ami.go \
  $(generated_code_deep_copy_helper) $(generated_code_aws_sdk_mocks)

all_generated_files := \
  site/content/usage/20-schema.md \
  $(all_generated_code)

.DEFAULT_GOAL := help

##@ Dependencies

.PHONY: install-build-deps
install-build-deps: ## Install dependencies (packages and tools)
	./install-build-deps.sh

##@ Build

godeps_cmd = go list -deps -f '{{if not .Standard}}{{ $$dep := . }}{{range .GoFiles}}{{$$dep.Dir}}/{{.}} {{end}}{{end}}' $(1) | sed "s|$(git_toplevel)/||g"
godeps = $(shell $(call godeps_cmd,$(1)))

.PHONY: build
build: $(all_generated_code) ## Build main binary
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

.PHONY: generate-all
# TODO: generate-ami is broken (see https://github.com/weaveworks/eksctl/issues/949 ), include it when fixed
generate-all: $(all_generated_files) # generate-ami ## Re-generate all the automatically-generated source files

.PHONY: check-all-generated-files-up-to-date
check-all-generated-files-up-to-date: generate-all
	git diff --quiet -- $(all_generated_files) || (git --no-pager diff $(all_generated_files); echo "HINT: to fix this, run 'git commit $(all_generated_files) --message \"Update generated files\"'"; exit 1)

pkg/addons/default/assets.go: pkg/addons/default/assets/*
	env GOBIN=$(GOBIN) time go generate ./$(@D)

pkg/addons/default/assets/aws-node.yaml:
	env GOBIN=$(GOBIN) go generate ./pkg/addons/default

pkg/nodebootstrap/assets.go: pkg/nodebootstrap/assets/*
	env GOBIN=$(GOBIN) time go generate ./$(@D)

.PHONY: pkg/nodebootstrap/maxpods.go
pkg/nodebootstrap/maxpods.go:
	env GOBIN=$(GOBIN) time go generate ./$(@D)

.license-header: LICENSE
	@# generate-groups.sh can't find the lincense header when using Go modules, so we provide one
	printf "/*\n%s\n*/\n" "$$(cat LICENSE)" > $@

deep_copy_helper_input := $(shell $(call godeps_cmd,./pkg/apis/...) | sed 's|$(generated_code_deep_copy_helper)||' )
$(generated_code_deep_copy_helper): $(deep_copy_helper_input) .license-header ## Generate Kubernetes API helpers
	time go mod download k8s.io/code-generator # make sure the code-generator is present
	time env GOPATH="$(gopath)" bash "$(gopath)/pkg/mod/k8s.io/code-generator@v0.0.0-20190808180452-d0071a119380/generate-groups.sh" \
	  deepcopy,defaulter _ ./pkg/apis eksctl.io:v1alpha5 --go-header-file .license-header --output-base="$(git_toplevel)" \
	  || (cat codegenheader.txt ; cat $(generated_code_deep_copy_helper); exit 1)

# static_resolver_ami.go doesn't only depend on files (it should be refreshed whenever a release is made in AWS)
# so we need to forcibly generate it
.PHONY: generate-ami
generate-ami: ## Generate the list of AMIs for use with static resolver. Queries AWS.
	time go generate ./pkg/ami

site/content/usage/20-schema.md: $(call godeps,cmd/schema/generate.go)
	time go run ./cmd/schema/generate.go $@

$(generated_code_aws_sdk_mocks): $(call godeps,pkg/eks/mocks/mocks.go)
	mkdir -p vendor/github.com/aws/
	@# Hack for Mockery to find the dependencies handled by `go mod`
	ln -sfn "$(gopath)/pkg/mod/github.com/aws/aws-sdk-go@v1.23.15" vendor/github.com/aws/aws-sdk-go
	time env GOBIN=$(GOBIN) go generate ./pkg/eks/mocks

##@ Docker

update_build_image_manifest = git ls-tree --full-tree @ -- $(build_image_input) > .build_image_manifest

.PHONY: check-build-image-manifest-up-to-date
check-build-image-manifest-up-to-date: ## Update build image manifest and commits the changes
	git diff --quiet -- $(build_image_input) || (git --no-pager diff $(build_image_input); exit 1)
	$(update_build_image_manifest)
	git diff --quiet -- .build_image_manifest || (git --no-pager diff .build_image_manifest; echo "HINT: to fix this, run 'make update-build-image-manifest'"; exit 1)

.PHONY: update-build-image-manifest
update-build-image-manifest: ## Update build image manifest and commits the changes
	$(update_build_image_manifest)
	git commit --quiet .build_image_manifest --message 'Update build image manifest'
	$(info build image name is $(build_image_name))

.PHONY: base-build-image
base-build-image: ## Build the base build image that has all of external dependencies
	-docker pull $(base_build_image_name)
	cat Dockerfile.base \
	  | $(docker_build) \
	  --iidfile=.base_build_image.iid \
	  --cache-from=$(base_build_image_name) \
	  --tag=$(base_build_image_name) -

.PHONY: build-image
build-image: base-build-image check-build-image-manifest-up-to-date ## Build the build image that has all of external dependencies
	-docker pull $(build_image_name)
	rm -rf .build_image_context
	mkdir .build_image_context
	cp $(build_image_input) .build_image_context/
	@# git only ensures owner's perssions are set, it ignores group and other's permissions,
	@# we need to force permissions here to avoid cache misses
	chmod go-wrx .build_image_context/*
	$(docker_build) \
	  --iidfile=.build_image.iid \
	  --cache-from=$(build_image_name) \
	  --build-arg=BASE_BUILD_IMAGE_ID=$(base_build_image_name) \
	  --tag=$(build_image_name) \
	  .build_image_context/

.PHONY: push-build-image
push-build-image: base-build-image build-image ## Push the build image to the registry
	docker push $(base_build_image_name)
	docker tag $(base_build_image_name) weaveworks/eksctl-build:base-latest
	docker push weaveworks/eksctl-build:base-latest
	docker push $(build_image_name)

.PHONY: intermediate-image
intermediate-image: build-image ## Build the intermediate image that has all artefacts
	time docker run \
	  --tty \
	  --name=$(intermidiate_container_name) \
	  --env=TEST_TARGET=$(TEST_TARGET) \
	  --volume=$(git_toplevel):/src \
	  --volume=$(gocache):/root/.cache/go-build \
	  --volume=$(gopath)/pkg/mod:/go/pkg/mod \
	    $(build_image_name) /src/eksctl-image-builder.sh \
	  || (docker rm $(intermidiate_container_name) ; exit 1)
	time docker commit $(intermidiate_container_name) $(intermediate_image_name) \
	  && docker rm $(intermidiate_container_name)

.PHONY: eksctl-image
eksctl-image: intermediate-image ## Build the eksctl image that has release artefacts and no build dependencies
	printf 'FROM scratch\nCMD eksctl\nCOPY --from=%s /out /' $(intermediate_image_name) \
	  | $(docker_build) \
	      --iidfile=.eksctl_image.iid \
	      --tag="$(eksctl_image_name)" -

##@ Release

docker_run_release_script = docker run \
  --env=GITHUB_TOKEN \
  --env=CIRCLE_TAG \
  --env=CIRCLE_PROJECT_USERNAME \
  --volume=$(CURDIR):/src \
  --workdir=/src \
    $(intermediate_image_name)

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
ifeq ($(OS),Windows_NT)
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make <target>\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  %-40s %s\n", $$1, $$2 } /^##@/ { printf "\n%s\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
else
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-40s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
endif
