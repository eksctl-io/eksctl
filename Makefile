built_at := $(shell date +%s)
git_commit := $(shell git describe --dirty --always)

EKSCTL_BUILD_IMAGE ?= weaveworks/eksctl:build
EKSCTL_IMAGE ?= weaveworks/eksctl:latest

.PHONY: build
build:
	go build -ldflags "-X main.gitCommit=$(git_commit) -X main.builtAt=$(built_at)" ./cmd/eksctl

.PHONY: test
test:
	go test $(shell go list ./... | grep -v "^github.com/weaveworks/eksctl/vendor" | sort -u)


.PHONY: update-bindata
update-bindata:
	go generate ./pkg/eks

.PHONY: install-bindata
install-bindata:
	go get -u github.com/jteeuwen/go-bindata/...

.PHONY: eksctl-build-image
eksctl-build-image:
	@-docker pull $(EKSCTL_BUILD_IMAGE)
	@docker build --tag=$(EKSCTL_BUILD_IMAGE) --cache-from=$(EKSCTL_BUILD_IMAGE) ./build

.PHONY: eksctl-image
eksctl-image: eksctl-build-image
	@docker build --tag=$(EKSCTL_IMAGE) --build-arg=EKSCTL_BUILD_IMAGE=$(EKSCTL_BUILD_IMAGE) ./

.PHONY: release
release: eksctl-build-image
	docker run \
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
