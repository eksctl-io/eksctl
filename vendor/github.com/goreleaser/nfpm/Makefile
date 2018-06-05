SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=
OS=$(shell uname -s)

export PATH := ./bin:$(PATH)

# Install all the build and lint dependencies
setup:
	curl -sfL https://install.goreleaser.com/github.com/caarlos0/bandep.sh | sh
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh
ifeq ($(OS), Darwin)
	brew install dep
else
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif
	dep ensure -vendor-only
	echo "make check" > .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
.PHONY: setup

check:
	bandep --ban github.com/tj/assert
.PHONY: check

pull_test_imgs:
	grep FROM ./acceptance/testdata/*.dockerfile | cut -f2 -d' ' | sort | uniq | while read -r img; do docker pull "$$img"; done
.PHONY: pull_test_imgs

test: pull_test_imgs
	go test $(TEST_OPTIONS) -v -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.out $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=5m
.PHONY: test

cover: test
	go tool cover -html=coverage.out
.PHONY: cover

fmt:
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done
.PHONY: fmt

lint: check
	golangci-lint run --enable-all ./...
.PHONY: check

ci: build lint test
.PHONY: ci

build:
	go build -o nfpm ./cmd/nfpm/main.go
.PHONY: build

todo:
	@grep \
		--exclude-dir=vendor \
		--exclude-dir=node_modules \
		--exclude-dir=bin \
		--exclude=Makefile \
		--text \
		--color \
		-nRo -E ' TODO:.*|SkipNow' .
.PHONY: todo

.DEFAULT_GOAL := build
