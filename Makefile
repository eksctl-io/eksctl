COMMIT:=$(shell git describe --dirty --always)

build: update-bindata
	go build -ldflags "-X main.commit=$(COMMIT)" ./cmd/eksctl

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
	  --volume=$(CURDIR):/go/src/github.com/weaveworks/eksctl \
	  --workdir=/go/src/github.com/weaveworks/eksctl \
	    eksctl_build \
	      goreleaser release
