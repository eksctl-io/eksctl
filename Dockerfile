# Make sure to bump BUILD_IMAGE_VERSION if you make any dependency changes

#
# Dependency cache
#
FROM golang:1.12-alpine3.9 AS dependencies

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
COPY go.mod go.sum install-build-deps.sh /src/

# Download all the dependencies
RUN go mod download

RUN ./install-build-deps.sh
RUN go install github.com/goreleaser/goreleaser

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


#
# Compiled dependencies cache
#
FROM dependencies as compiled_dependencies

COPY go-deps.txt /src/

# Build all the dependencies
RUN go build $(cat go-deps.txt)


#
# Builder image (not cached)
#
FROM compiled_dependencies AS builder

LABEL eksctl.builder=true

COPY . /src

ARG COVERALLS_TOKEN
ENV COVERALLS_TOKEN $COVERALLS_TOKEN

ENV JUNIT_REPORT_DIR /src/test-results/ginkgo
RUN mkdir -p "${JUNIT_REPORT_DIR}"

WORKDIR /src
ARG TEST_TARGET
RUN make $TEST_TARGET
RUN make build \
    && cp ./eksctl /out/usr/local/bin/eksctl
RUN make build-integration-test \
    && mkdir -p /out/usr/local/share/eksctl \
    && cp integration/*.yaml /out/usr/local/share/eksctl \
    && cp ./eksctl-integration-test /out/usr/local/bin/eksctl-integration-test

RUN go build github.com/kubernetes-sigs/aws-iam-authenticator/cmd/aws-iam-authenticator \
    && cp ./aws-iam-authenticator /out/usr/local/bin/aws-iam-authenticator


#
# Final image
#
FROM scratch
CMD eksctl
COPY --from=builder /out /
