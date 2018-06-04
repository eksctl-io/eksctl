builtAt := $(shell date +%s)
gitCommit := $(shell git describe --dirty --always)
gitTag := $(shell git describe --tags --abbrev=0)

build: update-bindata
	go build -ldflags "-X main.gitCommit=$(gitCommit) -X main.builtAt=$(builtAt)" ./cmd/eksctl

update-bindata:
	go generate ./pkg/eks

install-bindata:
	go get -u github.com/jteeuwen/go-bindata/...

eksctl_build_image:
	docker build --tag=eksctl_build ./build

eksctl_image: eksctl_build_image
	docker build --tag=eksctl ./

release: eksctl_build_image
	docker run \
	  --env=GITHUB_TOKEN \
	  --env=CIRCLE_TAG \
	  --volume=$(CURDIR):/go/src/github.com/weaveworks/eksctl \
	  --workdir=/go/src/github.com/weaveworks/eksctl \
	    eksctl_build \
	      make do_release

do_release:
	@if [ $(CIRCLE_TAG) = latest_release ] ; then \
	  git tag -d $(gitTag) ; \
	  github-release info --user weaveworks --repo eksctl --tag latest_release > /dev/null 2>&1 && \
	    github-release delete --user weaveworks --repo eksctl --tag latest_release ; \
	  goreleaser release --skip-validate --config=./.goreleaser.floating.yml ; \
	else \
	  goreleaser release --config=./.goreleaser.permalink.yml ; \
        fi
