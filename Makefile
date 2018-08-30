built_at := $(shell date +%s)
git_commit := $(shell git describe --dirty --always)

EKSCTL_BUILD_IMAGE ?= weaveworks/eksctl:build
EKSCTL_IMAGE ?= weaveworks/eksctl:latest

.PHONY: build
build:
	@go build -ldflags "-X main.gitCommit=$(git_commit) -X main.builtAt=$(built_at)" ./cmd/eksctl

.PHONY: install-build-deps
install-build-deps:
	@cd build && dep ensure && ./install.sh

.PHONY: test
test:
	@go test -v -covermode=count -coverprofile=coverage.out ./pkg/... ./cmd/...
	@test -z $(COVERALLS_TOKEN) || goveralls -coverprofile=coverage.out -service=circle-ci

.PHONY: integration-test-dev
integration-test-dev: build
	@go test -tags integration -v -timeout 21m ./tests/integration/... \
		-args \
		-eksctl.cluster=integration-test-dev \
		-eksctl.create=false \
		-eksctl.delete=false \
		-eksctl.kubeconfig=$(HOME)/.kube/eksctl/clusters/integration-test-dev

.PHONY: integration-test
integration-test: build
	@go test -tags integration -v -timeout 21m ./tests/integration/...

.PHONY: generated
generate:
	@go generate ./pkg/eks ./pkg/eks/mocks

.PHONY: eksctl-build-image
eksctl-build-image:
	@-docker pull $(EKSCTL_BUILD_IMAGE)
	@docker build --tag=$(EKSCTL_BUILD_IMAGE) --cache-from=$(EKSCTL_BUILD_IMAGE) ./build

EKSCTL_IMAGE_BUILD_ARGS := --build-arg=EKSCTL_BUILD_IMAGE=$(EKSCTL_BUILD_IMAGE)
ifneq ($(COVERALLS_TOKEN),)
EKSCTL_IMAGE_BUILD_ARGS += --build-arg=COVERALLS_TOKEN=$(COVERALLS_TOKEN)
endif

.PHONY: eksctl-image
eksctl-image: eksctl-build-image
	@docker build --tag=$(EKSCTL_IMAGE) $(EKSCTL_IMAGE_BUILD_ARGS) ./

.PHONY: release
release: eksctl-build-image
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

.PHONY: server-pages
serve-pages:
	@-docker rm -f eksctl-jekyll
	@$(JEKYLL) jekyll serve -d /_site --watch --force_polling -H 0.0.0.0 -P 4000

.PHONY: build-page
build-pages:
	@-docker rm -f eksctl-jekyll
	@$(JEKYLL) jekyll build --verbose
