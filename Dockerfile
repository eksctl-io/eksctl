# Make sure to bump the version of EKSCTL_BUILD_IMAGE if you make any changes in the buildcache
FROM golang:1.12-alpine3.9 AS buildcache

RUN apk add --no-cache \
      curl \
      git \
      make \
      bash \
      gcc \
      musl-dev \
      && true

ENV CGO_ENABLED=0

WORKDIR /src
# We intentionally don't copy the full source tree to prevent code-changes
# from invalidating the dependency build cache image
COPY go.mod go.sum go-deps.txt install-build-deps.sh /src/

ARG GO_BUILD_TAGS

# Download all the dependencies and build them
RUN go mod download
RUN go build -tags "${GO_BUILD_TAGS}" $(cat go-deps.txt)

RUN GO_BUILD_TAGS="${GO_BUILD_TAGS}" ./install-build-deps.sh
RUN go install -tags "${GO_BUILD_TAGS}" github.com/goreleaser/goreleaser

RUN mkdir -p /out/etc/apk && cp -r /etc/apk/* /out/etc/apk/
RUN apk add --no-cache --initdb --root /out \
    alpine-baselayout \
    busybox \
    ca-certificates \
    coreutils \
    git \
    libc6-compat \
    && true

ENV KUBECTL_VERSION v1.11.5
RUN curl --silent --location "https://dl.k8s.io/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" --output /out/usr/local/bin/kubectl \
    && chmod +x /out/usr/local/bin/kubectl

FROM buildcache

LABEL eksctl.builder=true

COPY . /src

ARG COVERALLS_TOKEN
ENV COVERALLS_TOKEN $COVERALLS_TOKEN

ARG TEST_TARGET
ENV JUNIT_REPORT_DIR /src/test-results/ginkgo
RUN mkdir -p "${JUNIT_REPORT_DIR}"

WORKDIR /src
RUN make $TEST_TARGET
RUN make build GO_BUILD_TAGS="${GO_BUILD_TAGS}" \
    && cp ./eksctl /out/usr/local/bin/eksctl
RUN make build-integration-test \
    && mkdir -p /out/usr/local/share/eksctl \
    && cp integration/*.yaml /out/usr/local/share/eksctl \
    && cp ./eksctl-integration-test /out/usr/local/bin/eksctl-integration-test

RUN go build github.com/kubernetes-sigs/aws-iam-authenticator/cmd/aws-iam-authenticator \
    && cp ./aws-iam-authenticator /out/usr/local/bin/aws-iam-authenticator

FROM scratch
CMD eksctl
COPY --from=buildcache  /out /
