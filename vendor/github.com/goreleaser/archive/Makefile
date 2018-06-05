SOURCE_FILES?=./...
TEST_PATTERN?=.
TEST_OPTIONS?=

setup:
	go get -u github.com/alecthomas/gometalinter
	go get -u golang.org/x/tools/cmd/cover
	go get -u -t ./...
	gometalinter --install

test:
	go test $(TEST_OPTIONS) -v -coverpkg=./... -race -covermode=atomic -coverprofile=coverage.txt $(SOURCE_FILES) -run $(TEST_PATTERN) -timeout=2m

cover: test
	go tool cover -html=coverage.txt

fmt:
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do gofmt -w -s "$$file"; goimports -w "$$file"; done

lint:
	gometalinter --vendor ./...

ci: lint test

.DEFAULT_GOAL := build
