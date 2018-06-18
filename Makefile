builtAt := $(shell date +%s)
gitCommit := $(shell git describe --dirty --always)
gitTag := $(shell git describe --tags --abbrev=0 --always)

EKSCTL_BUILD_IMAGE ?= weaveworks/eksctl:build
EKSCTL_IMAGE ?= weaveworks/eksctl:latest

.PHONY: build
build:
	go build -ldflags "-X main.gitCommit=$(gitCommit) -X main.builtAt=$(builtAt)" ./cmd/eksctl

.PHONY: update-bindata
update-bindata:
	go generate ./pkg/eks

.PHONY: install-bindata
install-bindata:
	go get -u github.com/jteeuwen/go-bindata/...

.PHONY: eksctl_build_image
eksctl_build_image:
	@-docker pull $(EKSCTL_BUILD_IMAGE)
	@docker build --tag=$(EKSCTL_BUILD_IMAGE) --cache-from=$(EKSCTL_BUILD_IMAGE) ./build

.PHONY: eksctl_image
eksctl_image: eksctl_build_image
	@docker build --tag=$(EKSCTL_IMAGE) --build-arg=EKSCTL_BUILD_IMAGE=$(EKSCTL_BUILD_IMAGE) ./

.PHONY: release
release: eksctl_build_image
	docker run \
	  --env=GITHUB_TOKEN \
	  --env=CIRCLE_TAG \
	  --volume=$(CURDIR):/go/src/github.com/weaveworks/eksctl \
	  --workdir=/go/src/github.com/weaveworks/eksctl \
	    eksctl_build \
	      make do_release

.PHONY: do_release
do_release:
	@if [ $(CIRCLE_TAG) = latest_release ] ; then \
	  git tag -d $(gitTag) ; \
	  github-release info --user weaveworks --repo eksctl --tag latest_release > /dev/null 2>&1 && \
	    github-release delete --user weaveworks --repo eksctl --tag latest_release ; \
	  goreleaser release --skip-validate --config=./.goreleaser.floating.yml ; \
	else \
	  goreleaser release --config=./.goreleaser.permalink.yml ; \
        fi

JEKYLL := docker run --tty --rm \
  --name=eksctl-jekyll \
  --volume="$(CURDIR)":/usr/src/app \
  --publish="4000:4000" \
    starefossen/github-pages

.PHONY: server_pages
serve_pages:
	@-docker rm -f eksctl-jekyll
	@$(JEKYLL) jekyll serve -d /_site --watch --force_polling -H 0.0.0.0 -P 4000

.PHONY: build_page
build_pages:
	@-docker rm -f eksctl-jekyll
	@$(JEKYLL) jekyll build --verbose
